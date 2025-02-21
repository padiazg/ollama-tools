/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	listmodels "github.com/padiazg/ollama-tools/internals/list-models"
	"github.com/spf13/cobra"
)

// memRequirementCmd represents the memRequirement command
var memRequirementCmd = &cobra.Command{
	Use:   "memory-requirement",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		// fmt.Println("memory-requirement called")
		listmodels.List(s)
	},
}

func init() {
	rootCmd.AddCommand(memRequirementCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// memRequirementCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// memRequirementCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
