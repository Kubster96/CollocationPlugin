package collocation

import (
	"k8s.io/kubernetes/pkg/scheduler/framework/v1alpha1"
	"reflect"
	"testing"
)

func TestNormalizeZeroScores(t *testing.T) {
	collocationScore := &CollocationScore{}
	nodeScores := []v1alpha1.NodeScore{
		{
			Name:  "first-node",
			Score: 0,
		},
		{
			Name:  "second-node",
			Score: 0,
		},
	}
	expectedNodeScores := []v1alpha1.NodeScore{
		{
			Name:  "first-node",
			Score: 50,
		},
		{
			Name:  "second-node",
			Score: 50,
		},
	}

	nodeScoreList := v1alpha1.NodeScoreList(nodeScores)
	expectedNodeScoreList := v1alpha1.NodeScoreList(expectedNodeScores)
	collocationScore.NormalizeScore(nil, nil, nil, nodeScoreList)

	if !reflect.DeepEqual(nodeScoreList, expectedNodeScoreList) {
		t.Errorf("Node score List %v is not equals to expected node score list (%v)", nodeScoreList, expectedNodeScoreList)
	}
}

func TestNormalizeRevertScores(t *testing.T) {
	collocationScore := &CollocationScore{}
	nodeScores := []v1alpha1.NodeScore{
		{
			Name:  "first-node",
			Score: 0,
		},
		{
			Name:  "second-node",
			Score: 123123,
		},
	}
	expectedNodeScores := []v1alpha1.NodeScore{
		{
			Name:  "first-node",
			Score: 100,
		},
		{
			Name:  "second-node",
			Score: 0,
		},
	}

	nodeScoreList := v1alpha1.NodeScoreList(nodeScores)
	expectedNodeScoreList := v1alpha1.NodeScoreList(expectedNodeScores)
	collocationScore.NormalizeScore(nil, nil, nil, nodeScoreList)

	if !reflect.DeepEqual(nodeScoreList, expectedNodeScoreList) {
		t.Errorf("Node score List %v is not equals to expected node score list (%v)", nodeScoreList, expectedNodeScoreList)
	}
}

func TestNormalizeDifferentScores(t *testing.T) {
	collocationScore := &CollocationScore{}
	nodeScores := []v1alpha1.NodeScore{
		{
			Name:  "first-node",
			Score: 67,
		},
		{
			Name:  "second-node",
			Score: 123123,
		},
		{
			Name:  "third-node",
			Score: 48764,
		},
	}
	expectedNodeScores := []v1alpha1.NodeScore{
		{
			Name:  "first-node",
			Score: 100,
		},
		{
			Name:  "second-node",
			Score: 0,
		},
		{
			Name:  "third-node",
			Score: 60,
		},
	}

	nodeScoreList := v1alpha1.NodeScoreList(nodeScores)
	expectedNodeScoreList := v1alpha1.NodeScoreList(expectedNodeScores)
	collocationScore.NormalizeScore(nil, nil, nil, nodeScoreList)

	if !reflect.DeepEqual(nodeScoreList, expectedNodeScoreList) {
		t.Errorf("Node score List %v is not equals to expected node score list (%v)", nodeScoreList, expectedNodeScoreList)
	}
}

func TestMaxSame(t *testing.T) {
	a, b := int64(2), int64(2)
	result := max(a, b)
	expectedResult := int64(2)

	if result != expectedResult {
		t.Errorf("Max function returned %d from values (%d, %d) but should return %d", result, a, b, expectedResult)
	}
}

func TestMaxDifferent(t *testing.T) {
	a, b := int64(10), int64(2)
	result := max(a, b)
	expectedResult := int64(10)

	if result != expectedResult {
		t.Errorf("Max function returned %d from values (%d, %d) but should return %d", result, a, b, expectedResult)
	}
}

func TestMinSame(t *testing.T) {
	a, b := int64(2), int64(2)
	result := min(a, b)
	expectedResult := int64(2)

	if result != expectedResult {
		t.Errorf("Min function returned %d from values (%d, %d) but should return %d", result, a, b, expectedResult)
	}
}

func TestMinDifferent(t *testing.T) {
	a, b := int64(10), int64(2)
	result := min(a, b)
	expectedResult := int64(2)

	if result != expectedResult {
		t.Errorf("Min function returned %d from values (%d, %d) but should return %d", result, a, b, expectedResult)
	}
}

func TestScoreExtensions(t *testing.T) {
	collocationScore := &CollocationScore{}
	result := collocationScore.ScoreExtensions()

	if collocationScore != result {
		t.Errorf("ScoreExtenstions function returned %v but should return %v", result, collocationScore)
	}
}

func TestName(t *testing.T) {
	collocationScore := &CollocationScore{}
	result := collocationScore.Name()
	expected := "CollocationScore"

	if result != expected {
		t.Errorf("Name function returned %v but should return %v", result, expected)
	}
}
