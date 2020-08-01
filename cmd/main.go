package main

import (
	"CollocationPlugin/pkg/collocation"
	"math/rand"
	"os"
	"time"

	"k8s.io/kubernetes/cmd/kube-scheduler/app"
)

func main() {
	rand.Seed(time.Now().UnixNano())
	// Register custom plugins to the scheduler framework.
	// Later they can consist of scheduler profile(s) and hence
	// used by various kinds of workloads.
	command := app.NewSchedulerCommand(
		app.WithPlugin(collocation.Name, collocation.New),
	)
	if err := command.Execute(); err != nil {
		os.Exit(1)
	}
}
