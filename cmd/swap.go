/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// swapCmd represents the swap command
var swapCmd = &cobra.Command{
	Use:   "swap",
	Short: "Swap bits in byte, or bytes in word, or words in dword",
	Long: `. For example:

1011 0110 -> 0110 1101 	# swap bits
ABCD 			-> BADC 			# swap half 
ABCD 			-> CDAB 			# swap word
ABCDEFGH 	-> EFGHABCD		# swap dword
`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("swap called")
	},
}

func init() {
	rootCmd.AddCommand(swapCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// swapCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// swapCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
