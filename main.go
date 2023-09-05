package main

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"log"
)

func main() {
	// Clear the console
	fmt.Printf("\x1bc")

	// Get AWS configuration
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatalln(err.Error())
	}

	// Set EC2 client
	client := ec2.NewFromConfig(cfg)

	// List all instances
	instances, err := ListInstances(client)
	if err != nil {
		log.Fatalln(err.Error())
	}

	// Print all instances
	PrintInstances(instances)

	// Get action selection
	action, err := ActionSelect()
	if err != nil {
		log.Fatalln(err.Error())
	}

	switch action {
	case ActionTypeStart:
		instances = instances.Stopped()
	case ActionTypeStop:
		instances = instances.Started()
	}

	// Get instance selection
	selectedInstances, err := InstanceSelect(action, instances)
	if err != nil {
		log.Fatalln(err.Error())
	}

	// Get user confirmation
	confirm, err := GetConfirm(action, selectedInstances)
	if err != nil {
		log.Fatalln(err.Error())
	}

	// Exit if not confirm
	if !confirm {
		return
	}

	switch action {
	case ActionTypeStart:
		StartInstances(client, selectedInstances)
	case ActionTypeStop:
		StopInstance(client, selectedInstances)
	}
}
