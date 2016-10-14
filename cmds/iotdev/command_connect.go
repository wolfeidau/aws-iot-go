package main

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/eclipse/paho.mqtt.golang"

	"github.com/spf13/cobra"
	"github.com/wolfeidau/aws-iot-go/pkg/provision"
)

var (
	cmdConnectThing = &cobra.Command{
		Use:   "connect",
		Short: "Connect thing",
		Long:  ``,
		Run:   runCmdThingConnect,
	}

	connectThingName string
)

func init() {
	cmdConnectThing.Flags().StringVar(&connectThingName, "name", "", "The name of the thing.")
	cmdRoot.AddCommand(cmdConnectThing)
}

func runCmdThingConnect(cmd *cobra.Command, args []string) {

	config, err := provision.LoadConfig(connectThingName)

	if err != nil {
		stderr("Failed to connect thing: %v", err)
		os.Exit(1)
	}

	cert, err := tls.X509KeyPair([]byte(config.CertificatePem), []byte(*config.KeyPair.PrivateKey))
	if err != nil {
		stderr("Failed to connect thing: %v", err)
		os.Exit(1)
	}

	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},
	}

	tlsConfig.BuildNameToCertificate()

	serverURL := fmt.Sprintf("ssl://data.iot.%s.amazonaws.com:8883", defaultRegion)

	if err != nil {
		stderr("Failed to connect thing: %v", err)
		os.Exit(1)
	}

	opts := mqtt.NewClientOptions()
	opts.AddBroker(serverURL)
	opts.SetClientID(connectThingName).SetTLSConfig(tlsConfig)
	opts.SetDefaultPublishHandler(handleChange)

	if rootOpts.Debug {
		mqtt.DEBUG = log.New(os.Stdout, "logger: ", log.Lshortfile)
	}

	c := mqtt.NewClient(opts)
	if token := c.Connect(); token.Wait() && token.Error() != nil {
		stderr("Failed to connect thing: %v", token.Error())
		os.Exit(1)
	}

	stdout("connected to %s", serverURL)

	// subscribe to the desired state topic for this thing
	if token := c.Subscribe(fmt.Sprintf("$aws/things/%s/shadow/#", connectThingName), 0, nil); token.Wait() && token.Error() != nil {
		stderr("Failed to subscribe to desired state topic: %v", token.Error())
		os.Exit(1)
	}

	currentState := map[string]interface{}{
		"state": map[string]interface{}{
			"reported": map[string]interface{}{
				"red":   187,
				"green": 114,
				"blue":  222,
			},
		},
	}

	b, err := json.Marshal(&currentState)

	if err != nil {
		stderr("Failed to serialise state: %v", err)
		os.Exit(1)
	}

	if token := c.Publish(fmt.Sprintf("$aws/things/%s/shadow/status", connectThingName), 0, false, b); token.Wait() && token.Error() != nil {
		stderr("Failed to publish thing state to update topic: %v", token.Error())
	}

	quitChannel := make(chan os.Signal)
	signal.Notify(quitChannel, syscall.SIGINT, syscall.SIGTERM)

	<-quitChannel

	stdout("Received quit.")
}

func handleChange(client mqtt.Client, msg mqtt.Message) {
	fmt.Printf("TOPIC: %s\n", msg.Topic())
	fmt.Printf("MSG: %s\n", msg.Payload())
}
