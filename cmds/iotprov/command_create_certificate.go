package main

import (
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/iot"
	"github.com/spf13/cobra"
	"github.com/wolfeidau/aws-iot-go/pkg/provision"
)

var (
	cmdCreateCertificate = &cobra.Command{
		Use:   "certificate",
		Short: "Create certificate and key for a thing.",
		Long:  ``,
		Run:   runCmdCreateCertificate,
	}
	certificateThingName string

	pubSubToAnyTopicPolicy = `{
    "Version": "2012-10-17",
    "Statement": [{
        "Effect": "Allow",
        "Action":["iot:*"],
        "Resource": ["*"]
    }]
}`
)

func init() {
	cmdCreateCertificate.Flags().StringVar(&certificateThingName, "name", "", "The name of the thing.")

	cmdRoot.AddCommand(cmdCreateCertificate)
}

func runCmdCreateCertificate(cmd *cobra.Command, args []string) {
	svc = iot.New(session.New(), newAWSConfig())

	err := checkExists(certificateThingName)

	if err != nil {
		stderr("Failed to create cerficate: %v", err)
		os.Exit(1)
	}

	resp, err := svc.CreateKeysAndCertificate(&iot.CreateKeysAndCertificateInput{
		SetAsActive: aws.Bool(true),
	})

	if err != nil {
		stderr("Failed to create cerficate: %v", err)
		os.Exit(1)
	}

	tconfig := provision.NewThingConfig(resp)

	err = tconfig.Save(certificateThingName)

	if err != nil {
		stderr("Failed to create cerficate: %v", err)
		os.Exit(1)
	}

	presp, err := svc.CreatePolicy(&iot.CreatePolicyInput{
		PolicyDocument: aws.String(pubSubToAnyTopicPolicy),
		PolicyName:     aws.String("PubSubToAnyTopic"),
	})

	if err != nil {
		stderr("Failed to create cerficate: %v", err)
		os.Exit(1)
	}

	_, err = svc.AttachPrincipalPolicy(&iot.AttachPrincipalPolicyInput{
		PolicyName: presp.PolicyName,
		Principal:  resp.CertificateArn,
	})

	if err != nil {
		stderr("Failed to create cerficate: %v", err)
		os.Exit(1)
	}

}

func checkExists(name string) error {
	if _, err := svc.DescribeThing(&iot.DescribeThingInput{ThingName: aws.String(name)}); err != nil {
		return err
	}

	if _, err := os.Stat(fmt.Sprintf("%s.yaml", name)); err == nil {
		return fmt.Errorf("config for %s already exists", name)
	}
	return nil
}
