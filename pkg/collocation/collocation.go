package collocation

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/klog"
	framework "k8s.io/kubernetes/pkg/scheduler/framework/v1alpha1"
	"math"
	"os"
	"strconv"
)

const (
	Name         = "CollocationScore"
	TypeSelector = "type"
	SizeSelector = "size"
)

type CollocationScore struct {
	handle framework.FrameworkHandle
}

var _ framework.ScorePlugin = &CollocationScore{}
var _ framework.ScoreExtensions = &CollocationScore{}

func (s CollocationScore) Name() string {
	return Name
}

type PodNameWithType struct {
	name  string
	typee string
}

func (s CollocationScore) Score(ctx context.Context, state *framework.CycleState, p *v1.Pod, nodeName string) (int64, *framework.Status) {
	klog.V(4).Infof("Calculating collocation score for pod %v for node %v", p.Name, nodeName)
	namespace := p.Namespace
	typeA := p.Annotations[TypeSelector]

	var podNamesWithTypes []PodNameWithType

	pods, err := s.handle.ClientSet().CoreV1().Pods(namespace).List(ctx, metav1.ListOptions{
		FieldSelector: "spec.nodeName=" + nodeName,
	})

	for _, pod := range pods.Items {
		podNamesWithTypes = append(podNamesWithTypes, PodNameWithType{
			name:  pod.Name,
			typee: pod.Annotations[TypeSelector],
		})
	}

	klog.V(4).Infof("Pods on node %v are %v", nodeName, podNamesWithTypes)

	score := float64(0)

	if err == nil && pods != nil {
		coefficientsFile, coefficientsFileErr := os.Open("coefficients.json")

		if coefficientsFileErr == nil {
			defer coefficientsFile.Close()
			byteValue, _ := ioutil.ReadAll(coefficientsFile)
			var coefficients map[string]interface{}
			json.Unmarshal(byteValue, &coefficients)

			for _, pod := range pods.Items {
				typeB := pod.Annotations[TypeSelector]
				sizeB := pod.Annotations[SizeSelector]
				sizeBNumber, errSize := strconv.ParseFloat(sizeB, 64)
				if errSize == nil {
					coefficientString := coefficients[typeA+";"+typeB]
					if coefficientString != nil {
						coefficient, errCoefficient := strconv.ParseFloat(fmt.Sprintf("%v", coefficientString), 64)
						if errCoefficient == nil {
							score += coefficient * sizeBNumber
						} else {
							return int64(0), framework.NewStatus(framework.Error, "Could not parse coefficient for types: "+typeA+" and "+typeB)
						}
					} else {
						return int64(0), framework.NewStatus(framework.Error, "Coefficient does not exist for types: "+typeA+" and "+typeB)
					}
				} else {
					return int64(0), framework.NewStatus(framework.Error, "Could not parse size for pod "+pod.Name)
				}
			}
			klog.V(4).Infof("Calculated score for pod %v and node %v - %v", p.Name, nodeName, int64(score))
			return int64(score), framework.NewStatus(framework.Success, "")
		} else {
			return int64(0), framework.NewStatus(framework.Error, "Could not read coefficients.json file")
		}
	} else {
		return int64(0), framework.NewStatus(framework.Error, "Could not get pods from node "+nodeName)
	}
}

func (s CollocationScore) ScoreExtensions() framework.ScoreExtensions {
	return s
}

func (s CollocationScore) NormalizeScore(ctx context.Context, state *framework.CycleState, p *v1.Pod, nodeScores framework.NodeScoreList) *framework.Status {
	maximum := int64(math.MinInt64)
	minimum := int64(math.MaxInt64)
	for _, nodeScore := range nodeScores {
		maximum = max(maximum, nodeScore.Score)
		minimum = min(minimum, nodeScore.Score)
	}

	highest := float64(maximum)
	lowest := float64(minimum)

	if lowest != highest {
		for i, nodeScore := range nodeScores {
			nodeScores[i].Score = int64(((1 - ((float64(nodeScore.Score) - lowest) / (highest - lowest))) * float64(framework.MaxNodeScore-framework.MinNodeScore)) + float64(framework.MinNodeScore))
		}
	} else {
		for i, _ := range nodeScores {
			nodeScores[i].Score = ((framework.MaxNodeScore - framework.MinNodeScore) / 2) + framework.MinNodeScore
		}
	}
	klog.V(4).Infof("Scores after normalization %v", nodeScores)

	return framework.NewStatus(framework.Success, "")
}

func max(a int64, b int64) int64 {
	if a >= b {
		return a
	}
	return b
}

func min(a int64, b int64) int64 {
	if a <= b {
		return a
	}
	return b
}

// New initializes a new plugin and returns it.
func New(_ *runtime.Unknown, handle framework.FrameworkHandle) (framework.Plugin, error) {
	return &CollocationScore{
		handle: handle,
	}, nil
}
