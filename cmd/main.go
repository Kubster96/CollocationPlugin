package main

import (
	"CollocationPlugin/pkg/collocation"
	"k8s.io/kubernetes/cmd/kube-scheduler/app"
	"os"
)

func main() {
	command := app.NewSchedulerCommand(
		app.WithPlugin(collocation.Name, collocation.New),
	)
	if err := command.Execute(); err != nil {
		os.Exit(1)
	}
}
