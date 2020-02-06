package cmd

import (
	"github.com/spf13/cobra"
	"os"
)

import "fmt"

func runFunc(cmd *cobra.Command, args []string) {
	fmt.Println("config file:", cfgFilePath, "is debug :", debug, "port is :", port)
	if cfgFilePath == "" {
		cmd.Println("Please provide the path to your melody config file ")
	}
	//Parse config file
	serviceConfig, err := parser.Parse(cfgFilePath)
	if err != nil {
		//Show config file parse error and exit
		cmd.Printf("ERROR parsing the melody config file: %s\n", err.Error())
		os.Exit(-1)
	}
	//Judge is debug
	serviceConfig.Debug = serviceConfig.Debug || debug
	//Run with service config
	run(serviceConfig)
}
