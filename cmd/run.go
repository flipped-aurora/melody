package cmd

import "github.com/spf13/cobra"

import "fmt"

func runFunc(cmd *cobra.Command, args []string) {
	fmt.Println("config file:", cfgFile)
	fmt.Println("is debug :", debug)
	fmt.Println("port is :", port)
	//TODO some preparation for server
}
