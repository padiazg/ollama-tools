/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/padiazg/ollama-tools/internals/models"
	"github.com/spf13/cobra"
)

// listModels represents the memRequirement command
var listModels = &cobra.Command{
	Use:   "list-models [model-name]",
	Short: "List models using the Ollama api",
	Long: `List models using the Ollama api
	
If no model-name is espified all models will be retieved and listed. 
You can pass the model-name as an argument or using the --model-name flag`,
	// Args: cobra.ArbitraryArgs,
	Run: func(cmd *cobra.Command, args []string) {
		var (
			model_name string
			as_table   bool
			err        error
		)

		model_name, err = cmd.Flags().GetString("model-name")
		if err != nil {
			fmt.Printf("getting quantization-level flag: %+v", err)
			return
		}

		as_table, err = cmd.Flags().GetBool("table")
		if err != nil {
			fmt.Printf("getting table flag: %+v", err)
			return
		}

		if len(args) > 0 {
			model_name = args[0]
		}

		models.List(s, model_name, as_table)

		// if as_table {
		// 	models.ListTable(s, model_name)
		// } else {

		// }
	},
}

func init() {
	rootCmd.AddCommand(listModels)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// listModels.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// listModels.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	listModels.Flags().StringP("model-name", "m", "", "Model to list")
	listModels.Flags().BoolP("table", "t", false, "Print as table")
}
