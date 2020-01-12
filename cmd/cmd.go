package cmd

import (
	"fmt"
	"melody/config"
	"os"
)

//Executor defines the func that contains some prepration handles of start server.
type Executor func(config.ServiceConfig)

//Execute for other method to call.
func Execute(configParser config.Parser, executor Executor) {
	parser = configParser
	run = executor
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
