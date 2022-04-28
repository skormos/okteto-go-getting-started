package main

import (
	"fmt"

	"k8s.io/client-go/kubernetes"
	v1 "k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/client-go/rest"
)

func k8sInClusterCoreClient() (v1.CoreV1Interface, error) {
	config, err := rest.InClusterConfig()
	if err != nil {
		return nil, fmt.Errorf("while initializing InClusterConfig %w", err)
	}

	client, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("while initializing ClientSet for InClusterConfig %w", err)
	}

	return client.CoreV1(), nil
}
