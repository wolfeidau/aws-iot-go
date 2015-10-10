package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	// Version The version of the application (set by make file)
	Version       = "UNKNOWN"
	defaultRegion = "us-west-2"

	cmdRoot = &cobra.Command{
		Use:   "iotdev",
		Short: "Operate AWS IoT Device",
		Long:  ``,
	}

	rootOpts struct {
		Debug bool
	}
)

func init() {
	cmdRoot.PersistentFlags().BoolVar(&rootOpts.Debug, "debug", false, "Log debug information")
}

func main() {
	cmdRoot.Execute()
}

func stderr(msg string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, msg+"\n", args...)
}

func stdout(msg string, args ...interface{}) {
	fmt.Fprintf(os.Stdout, msg+"\n", args...)
}
