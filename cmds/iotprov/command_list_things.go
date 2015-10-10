package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/iot"
	"github.com/spf13/cobra"
)

var (
	cmdListThings = &cobra.Command{
		Use:   "list",
		Short: "List things",
		Long:  ``,
		Run:   runCmdThingList,
	}

	listThingType string
)

func init() {
	cmdListThings.Flags().StringVar(&listThingType, "type", "", "The type of the thing.")
	cmdRoot.AddCommand(cmdListThings)
}

func runCmdThingList(cmd *cobra.Command, args []string) {

	resp, err := svc.ListThings(buildListThings())

	if err != nil {
		stderr("Failed listing things: %v", err)
		os.Exit(1)
	}

	b, _ := json.MarshalIndent(resp, "", "  ")
	fmt.Fprintf(os.Stdout, "%s\n", string(b))
	//	fmt.Println("> Number of things: ", len(resp.Things))
}

func buildListThings() *iot.ListThingsInput {
	lt := &iot.ListThingsInput{}

	if listThingType != "" {
		lt.AttributeName = aws.String("Type")
		lt.AttributeValue = aws.String(listThingType)
	}

	return lt
}
