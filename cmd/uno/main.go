package main

import (
	"fmt"
	"log"
	"os"

	"github.com/appetito/uno/gen"

	"github.com/spf13/cobra"
)

func main() {
	// var file string
	// flag.StringVar(&file, "f", "", "YAML file path")
	// flag.Parse()

	// gen.Gen(file)
	// Создаем корневую команду
	var rootCmd = &cobra.Command{
		Use:   "uno",
		Short: "Uno framework CLI - API generator, project scaffolder",
		Long:  "Uno is a framework for building microservices in Go. This CLI tool helps to generate API from YAML file and scaffold project structure.",
	}

	// Добавляем подкоманды
	rootCmd.AddCommand(apigenCommand())
	rootCmd.AddCommand(initCommand())

	// Выполняем корневую команду
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func apigenCommand() *cobra.Command {
	var file string
	var cmd = &cobra.Command{
		Use:   "apigen",
		Short: "Generate API from YAML file",
		Long:  "This command generates (overwrites) api package from YAML file",
		Run: func(cmd *cobra.Command, args []string) {
			gen.GenAPI(file)
		},
	}
	cmd.Flags().StringVarP(&file, "file", "f", "", "YAML file path")	
	return cmd
}

func initCommand() *cobra.Command {
	var file string
	var cmd = &cobra.Command{
		Use:   "init",
		Short: "Scaffold project",
		Long:  "This command initializes project - generates API package and initial project structure, including service and handlers",
		Run: func(cmd *cobra.Command, args []string) {
			gen.GenProject(file)
			log.Println("Project initialized")

		},
	}
	cmd.Flags().StringVarP(&file, "file", "f", "", "YAML file path")
	return cmd
}