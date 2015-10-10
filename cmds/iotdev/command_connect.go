package main

import (
	"crypto/tls"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"git.eclipse.org/gitroot/paho/org.eclipse.paho.mqtt.golang.git"

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

	if rootOpts.Debug {
		mqtt.DEBUG = log.New(os.Stdout, "logger: ", log.Lshortfile)
	}

	c := mqtt.NewClient(opts)
	if token := c.Connect(); token.Wait() && token.Error() != nil {
		stderr("Failed to connect thing: %v", token.Error())
		os.Exit(1)
	}

	stdout("connected to %s", serverURL)

	c.Subscribe(fmt.Sprintf("%s/change", connectThingName), 0, handleChange)

	quitChannel := make(chan os.Signal)
	signal.Notify(quitChannel, syscall.SIGINT, syscall.SIGTERM)

	<-quitChannel

	stdout("Received quit.")
}

func handleChange(client *mqtt.Client, msg mqtt.Message) {
	fmt.Printf("TOPIC: %s\n", msg.Topic())
	fmt.Printf("MSG: %s\n", msg.Payload())
}
