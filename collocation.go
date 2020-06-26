package main

import (
	"context"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/runtime"
	corelisters "k8s.io/client-go/listers/core/v1"
	"k8s.io/kubernetes/pkg/api/v1/pod"
	framework "k8s.io/kubernetes/pkg/scheduler/framework/v1alpha1"
)

const Name = "CollocationScore"

type Score struct {
	podLister corelisters.PodLister
}

var _ framework.ScorePlugin = &Score{}

func (s Score) Name() string {
	return Name
}

func (s Score) Score(ctx context.Context, state *framework.CycleState, p *v1.Pod, nodeName string) (int64, *framework.Status) {
	selector := labels.Set{PodGroupName: }.AsSelector();
	namespace := p.Namespace
	pods, err := s.podLister.Pods(namespace).List(selector)


	return 122, framework.NewStatus(framework.Success, "")
}

func (s Score) ScoreExtensions() framework.ScoreExtensions {
	return Score.NormalizeScore
}

func (s Score) NormalizeScore(state *framework.CycleState, _ *v1.Pod, nodeScores framework.NodeScoreList) *framework.Status {
	highest := int64(0)
	for _, nodeScore := range nodeScores {
		highest = max(highest, nodeScore.Score)
	}
	for i, nodeScore := range nodeScores {
		nodeScores[i].Score = nodeScore.Score*framework.MaxNodeScore/highest
	}
	return framework.NewStatus(framework.Success, "")
}

func New(_ *runtime.Unknown, handle framework.FrameworkHandle) (framework.Plugin, error) {
	podLister := handle.SharedInformerFactory().Core().V1().Pods().Lister()
	return &Score{
		podLister: podLister,
	}, nil
}

func max(highest int64, score int64) int64 {
	if highest >= score {
		return highest
	}
	return score
}



