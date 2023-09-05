package main

import (
	"context"
	"fmt"
	"github.com/AlecAivazis/survey/v2"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/fatih/color"
	"github.com/rodaine/table"
	"log"
)

var NoCompatibleInstanceError = fmt.Errorf("no compatible instance for this action")

type InstanceList struct {
	Instances []Instance
}

type Instance struct {
	InstanceId        string
	InstanceName      string
	InstanceStateName types.InstanceStateName
}

func (i InstanceList) Started() InstanceList {
	var output = InstanceList{}
	for _, value := range i.Instances {
		if value.InstanceStateName == types.InstanceStateNameRunning {
			output.Instances = append(output.Instances, value)
		}
	}
	return output
}

func (i InstanceList) Stopped() InstanceList {
	var output = InstanceList{}
	for _, value := range i.Instances {
		if value.InstanceStateName == types.InstanceStateNameStopped {
			output.Instances = append(output.Instances, value)
		}
	}
	return output
}

func (i InstanceList) ToList() []string {
	var output []string
	for _, value := range i.Instances {
		output = append(output, value.InstanceId)
	}
	return output
}

func (i InstanceList) ToMap() map[string]string {
	var output = make(map[string]string)
	for _, value := range i.Instances {
		output[value.InstanceId] = fmt.Sprintf("%s - (Name: %s)", value.InstanceStateName, value.InstanceName)
	}
	return output
}

type InstanceIdAnswer struct {
	Values []string `survey:"instances"`
}

func InstanceSelect(action ActionType, instances InstanceList) ([]string, error) {
	var answer InstanceIdAnswer
	var qs = []*survey.Question{
		{
			Name: "instances",
			Prompt: &survey.MultiSelect{
				Message: fmt.Sprintf("Which instances did you want to %s", action.String()),
				Options: instances.ToList(),
				Description: func(value string, index int) string {
					return instances.ToMap()[value]
				},
				PageSize: 10,
			},
		},
	}
	if err := survey.Ask(qs, &answer); err != nil {
		return []string{}, NoCompatibleInstanceError
	}
	return answer.Values, nil
}

func StopInstance(client *ec2.Client, instances []string) {
	_, err := client.StopInstances(
		context.TODO(),
		&ec2.StopInstancesInput{
			InstanceIds: instances,
		},
	)
	if err != nil {
		log.Fatalln(err.Error())
	}
	for _, i := range instances {
		fmt.Printf("Stopping %s...", i)
	}
}

func StartInstances(client *ec2.Client, instances []string) {
	_, err := client.StartInstances(
		context.TODO(),
		&ec2.StartInstancesInput{
			InstanceIds: instances,
		},
	)
	if err != nil {
		log.Fatalln(err.Error())
	}

	for _, i := range instances {
		fmt.Printf("Starting %s...", i)
	}
}

func PrintInstances(instances InstanceList) {
	headerFmt := color.New(color.FgGreen, color.Underline).SprintfFunc()
	columnFmt := color.New(color.FgYellow).SprintfFunc()

	output := table.New("ID", "Name", "State")
	output.WithHeaderFormatter(headerFmt).WithFirstColumnFormatter(columnFmt)

	for _, instance := range instances.Instances {
		output.AddRow(instance.InstanceId, instance.InstanceName, instance.InstanceStateName)
	}

	output.Print()
	fmt.Println()
}

func ListInstances(client *ec2.Client) (InstanceList, error) {
	response, err := client.DescribeInstances(context.TODO(), &ec2.DescribeInstancesInput{})
	var result InstanceList
	if err != nil {
		return result, err
	}

	for _, reservation := range response.Reservations {
		for _, instance := range reservation.Instances {
			var instanceName string
			for _, tag := range instance.Tags {
				if aws.ToString(tag.Key) == "Name" {
					instanceName = aws.ToString(tag.Value)
					break
				}
			}
			i := Instance{
				InstanceId:        aws.ToString(instance.InstanceId),
				InstanceName:      instanceName,
				InstanceStateName: instance.State.Name,
			}

			result.Instances = append(result.Instances, i)
		}
	}
	return result, nil
}
