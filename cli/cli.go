package main

import (
	"fmt"
	"os"

	generators "github.com/auditt98/onthego/generators/methods"
	"github.com/auditt98/onthego/utils"
	"github.com/spf13/cobra"
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "app",
		Short: "app",
	}

	genCmd := &cobra.Command{
		Use:   "gen",
		Short: "Generate something",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(len(args))
			if len(args) == 0 {
				fmt.Println("Before os.Exit(1)")
				os.Exit(1)
			} else {
				param := args[0] // Retrieve the parameter
				if param == "handler" {
					utils.ReplaceMarkerInSection("tests", "a:=1", "./tests/tests.go")
					// Run gofmt command on the file to format it

					// reader := bufio.NewReader(os.Stdin)
					// fmt.Print("Enter the name of the handler: ")
					// handlerName, _ := reader.ReadString('\n')
					// fmt.Printf("Generating handler: %s", handlerName)
				}
			}
		},
	}
	rootCmd.AddCommand(genCmd)

	structureCmd := &cobra.Command{
		Use:   "structure",
		Short: "Generate the project structure",
		Run: func(cmd *cobra.Command, args []string) {
			//parse config.yaml to generate the project structure and code
			fmt.Println("Generating project structure...")
			generators.InitGenerator()

		},
	}
	rootCmd.AddCommand(structureCmd)

	// Execute the CLI command
	if err := rootCmd.Execute(); err != nil {
		fmt.Println("Before os.Exit(1)")
		os.Exit(1)
	}
}
