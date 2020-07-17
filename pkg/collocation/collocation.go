package collocation

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
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

func (s CollocationScore) Score(ctx context.Context, state *framework.CycleState, p *v1.Pod, nodeName string) (int64, *framework.Status) {
	namespace := p.Namespace
	typeA := p.Annotations[TypeSelector]

	pods, err := s.handle.ClientSet().CoreV1().Pods(namespace).List(ctx, metav1.ListOptions{
		FieldSelector: "spec.nodeName=" + nodeName,
	})

	score := float64(0)

	if err == nil && pods != nil {
		coefficientsFile, coefficientsFileErr := os.Open("coefficients.json")

		if coefficientsFileErr == nil {
			defer coefficientsFile.Close()
			byteValue, _ := ioutil.ReadAll(coefficientsFile)
			var coefficients map[string]interface{}
			json.Unmarshal(byteValue, &coefficients)

			for _, pood := range pods.Items {
				typeB := pood.Annotations[TypeSelector]
				sizeB := pood.Annotations[SizeSelector]
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
					return int64(0), framework.NewStatus(framework.Error, "Could not parse size for pod "+pood.Name)
				}
			}
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
	highest := int64(math.MinInt64)
	lowest := int64(math.MaxInt64)
	for _, nodeScore := range nodeScores {
		highest = max(highest, nodeScore.Score)
		lowest = min(lowest, nodeScore.Score)
	}

	minimum := float64(lowest)
	maximum := float64(highest)

	for i, nodeScore := range nodeScores {
		nodeScores[i].Score = int64((1 - ((float64(nodeScore.Score) - minimum) / (maximum - minimum))) * float64(framework.MaxNodeScore))
	}
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
