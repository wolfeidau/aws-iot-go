package main

import (
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/iot/iotiface"
	"github.com/spf13/cobra"
)

var (
	// Version The version of the application (set by make file)
	Version = "UNKNOWN"

	cmdRoot = &cobra.Command{
		Use:   "iotprov",
		Short: "Manage AWS IoT Devices",
		Long:  ``,
	}

	rootOpts struct {
		AWSDebug bool
	}

	svc iotiface.IoTAPI
)

func init() {
	cmdRoot.PersistentFlags().BoolVar(&rootOpts.AWSDebug, "aws-debug", false, "Log debug information from aws-sdk-go library")
}

func newAWSConfig() *aws.Config {
	c := aws.NewConfig()
	if rootOpts.AWSDebug {
		c = c.WithLogLevel(aws.LogDebug)
	}
	return c
}

func main() {
	cmdRoot.Execute()
}

func stderr(msg string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, msg+"\n", args...)
}
