package main

import (
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/iot"
	"github.com/spf13/cobra"
)

var (
	cmdCreateThing = &cobra.Command{
		Use:   "create",
		Short: "Create thing",
		Long:  ``,
		Run:   runCmdCreateThing,
	}
	name      string
	thingType string
)

func init() {
	cmdCreateThing.Flags().StringVar(&thingType, "type", "", "The type of the thing.")
	cmdCreateThing.Flags().StringVar(&name, "name", "", "The name of the thing.")

	cmdRoot.AddCommand(cmdCreateThing)
}

func runCmdCreateThing(cmd *cobra.Command, args []string) {
	svc = iot.New(session.New(), newAWSConfig())

	resp, err := svc.CreateThing(
		&iot.CreateThingInput{
			ThingName:        aws.String(name),
			AttributePayload: buildAttributes(),
		})

	if err != nil {
		stderr("Failed to create thing: %v", err)
		os.Exit(1)
	}

	fmt.Println("Thing ARN: ", *resp.ThingArn)
}

func buildAttributes() *iot.AttributePayload {

	ap := &iot.AttributePayload{
		Attributes: map[string]*string{},
	}

	if thingType != "" {
		ap.Attributes["Type"] = aws.String(thingType)
	}

	return ap
}
