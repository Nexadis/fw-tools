/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"errors"
	"log"

	"github.com/spf13/cobra"

	"github.com/Nexadis/fw-tools/internal/config"
	"github.com/Nexadis/fw-tools/internal/swap"
)

var swapCfg config.Swap

// swapCmd represents the swap command
var swapCmd = &cobra.Command{
	Use:   "swap [filename] [filename2]...",
	Short: "Swap bits in byte, or bytes in word, or words in dword",
	Long: `Swap bits in byte, or bytes in word, or words in dword. For example:

	1011 0110 -> 0110 1101 	# inverse bits
	ABCD 			-> BADC 			# swap halfs
	ABCD 			-> CDAB 			# swap bytes
	ABCD1234 	-> 1234ABCD		# swap words
	ABCD1234567890EF -> 0x567890EFABCD1234 # swap dwords
`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("set filename")
		}
		cfg.Inputs = append(cfg.Inputs, args...)
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		for _, input := range cfg.Inputs {
			s := swap.Swapper{
				Input:  input,
				Output: cfg.Output,
				Config: swapCfg,
			}
			err := s.Run()
			if err != nil {
				log.Fatal(err)
			}
		}
	},
}

func init() {
	swapCmd.Flags().BoolVarP(&swapCfg.Bits, "bits", "ib", false, "Inverse bits in byte")
	swapCmd.Flags().BoolVarP(&swapCfg.Halfs, "halfs", "h", false, "Swap halfs of byte")
	swapCmd.Flags().BoolVarP(&swapCfg.Bytes, "bytes", "sb", false, "Swap neighbors bytes")
	swapCmd.Flags().BoolVarP(&swapCfg.Words, "words", "w", false, "Swap neighbors words")
	swapCmd.Flags().BoolVarP(&swapCfg.Dwords, "dwords", "dw", false, "Swap neighbors dwords")

	rootCmd.AddCommand(swapCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// swapCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// swapCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
