package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

func checkFunc(cmd *cobra.Command, args []string) {
	//检查配置文件是否为空
	if cfgFilePath == "" {
		cmd.Println("Please, provide the path to your config file")
		return
	}
	cmd.Printf("Parsing configuration file: %s\n", cfgFilePath)
	v, err := parser.Parse(cfgFilePath)

	//如果开启了debug
	if debug {

		cmd.Printf("Parsed configuration: CacheTTL: %s, Port: %d\n", v.CacheTTL.String(), v.Port)
		// TODO

	}
	//如果错误
	if err != nil {
		cmd.Println("ERROR parsing the configuration file.\n", err.Error())
		os.Exit(1)
		return
	}

	cmd.Println("Syntax OK!")
}
