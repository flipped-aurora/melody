package cmd

import (
	"encoding/base64"
	"melody/config"
	"melody/core"

	"github.com/spf13/cobra"
)

var (
	cfgFilePath string
	debug       bool
	port        int
	parser      config.Parser
	run         Executor
	rootCmd     = &cobra.Command{
		Use:   "melody",
		Short: "melody help you to sort out your complex api.",
	}

	//添加校验配置文件
	checkCmd = &cobra.Command{
		Use:     "check",
		Short:   "check that the config",
		Long:    "Validates that the active configuration file has a valid syntax to run the service.\nChange the configuration file by using the --config flag",
		Run:     checkFunc,
		Aliases: []string{"validate"},
		Example: "melody check -d -c config.json",
	}
	runCmd = &cobra.Command{
		Use:     "run ",
		Short:   "run the Melody server",
		Long:    "run the Melody server",
		Run:     runFunc,
		Example: "melody run -d -c melody.json",
	}
)

func init() {
	logo, _ := base64.StdEncoding.DecodeString(encodedLogo)
	rootCmd.SetHelpTemplate("\n" + string(logo) + "\nVersion:" + core.MelodyVersion + "\n\n" + rootCmd.HelpTemplate())
	rootCmd.PersistentFlags().StringVarP(&cfgFilePath, "config", "c", "", "Path of the melody.json")
	rootCmd.PersistentFlags().BoolVarP(&debug, "debug", "d", false, "Enable the Melody debug")
	rootCmd.AddCommand(runCmd)
	rootCmd.AddCommand(checkCmd)
	runCmd.PersistentFlags().IntVarP(&port, "port", "p", 7777, "Listening port for Melody server")
}

const encodedLogo = "4paI4paI4paI4pWXICAg4paI4paI4paI4pWX4paI4paI4paI4paI4paI4paI4paI4pWX4paI4paI4pWXICAgICAg4paI4paI4paI4paI4paI4paI4pWXIOKWiOKWiOKWiOKWiOKWiOKWiOKVlyDilojilojilZcgICDilojilojilZcK4paI4paI4paI4paI4pWXIOKWiOKWiOKWiOKWiOKVkeKWiOKWiOKVlOKVkOKVkOKVkOKVkOKVneKWiOKWiOKVkSAgICAg4paI4paI4pWU4pWQ4pWQ4pWQ4paI4paI4pWX4paI4paI4pWU4pWQ4pWQ4paI4paI4pWX4pWa4paI4paI4pWXIOKWiOKWiOKVlOKVnQrilojilojilZTilojilojilojilojilZTilojilojilZHilojilojilojilojilojilZcgIOKWiOKWiOKVkSAgICAg4paI4paI4pWRICAg4paI4paI4pWR4paI4paI4pWRICDilojilojilZEg4pWa4paI4paI4paI4paI4pWU4pWdIArilojilojilZHilZrilojilojilZTilZ3ilojilojilZHilojilojilZTilZDilZDilZ0gIOKWiOKWiOKVkSAgICAg4paI4paI4pWRICAg4paI4paI4pWR4paI4paI4pWRICDilojilojilZEgIOKVmuKWiOKWiOKVlOKVnSAgCuKWiOKWiOKVkSDilZrilZDilZ0g4paI4paI4pWR4paI4paI4paI4paI4paI4paI4paI4pWX4paI4paI4paI4paI4paI4paI4paI4pWX4pWa4paI4paI4paI4paI4paI4paI4pWU4pWd4paI4paI4paI4paI4paI4paI4pWU4pWdICAg4paI4paI4pWRICAgCuKVmuKVkOKVnSAgICAg4pWa4pWQ4pWd4pWa4pWQ4pWQ4pWQ4pWQ4pWQ4pWQ4pWd4pWa4pWQ4pWQ4pWQ4pWQ4pWQ4pWQ4pWdIOKVmuKVkOKVkOKVkOKVkOKVkOKVnSDilZrilZDilZDilZDilZDilZDilZ0gICAg4pWa4pWQ4pWdICAgCiAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAg"
