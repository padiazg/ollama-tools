/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/padiazg/ollama-tools/internals/tools"
	"github.com/spf13/cobra"
)

// estimateCmd represents the estimate command
var estimateCmd = &cobra.Command{
	Use:   "estimate",
	Short: "Estimates the RAM requirement based on few paramaters ",
	Long:  `Estimates the RAM rwquirement based on few parameters without the need to download any model`,
	Run: func(cmd *cobra.Command, args []string) {
		var (
			parameter_count    int64
			context_length     int
			quantization_level string
			err                error
		)

		parameter_count, err = cmd.Flags().GetInt64("parameter-count")
		if err != nil {
			fmt.Printf("getting parameter-count: %+v", err)
			return
		}

		context_length, err = cmd.Flags().GetInt("context-length")
		if err != nil {
			fmt.Printf("getting context-length: %+v", err)
			return
		}

		quantization_level, err = cmd.Flags().GetString("quantization-level")
		if err != nil {
			fmt.Printf("getting quantization-level: %+v", err)
			return
		}

		mem := tools.EstimateMemory(parameter_count, context_length, quantization_level)
		tools.PrintEstimatedMemoryPlain(mem)
	},
}

func init() {
	rootCmd.AddCommand(estimateCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// estimateCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// estimateCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	estimateCmd.Flags().Int64P("parameter-count", "p", 0, "Parameters count")
	estimateCmd.Flags().IntP("context-length", "c", 0, "Context length")
	estimateCmd.Flags().StringP("quantization-level", "q", "", "Quantization level")
	estimateCmd.MarkFlagRequired("parameter-count")
	estimateCmd.MarkFlagRequired("context-length")
	estimateCmd.MarkFlagRequired("quantization-level")
}
