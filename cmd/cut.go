/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// cutCmd represents the cut command
var cutCmd = &cobra.Command{
	Use:   "cut",
	Short: "Cut metainfo between pages with step",
	Long:  `In some type of memory dumps we can meet additional meta-information about pages. You can use this command for cut it.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("cut called")
	},
}

func init() {
	rootCmd.AddCommand(cutCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// cutCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// cutCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
