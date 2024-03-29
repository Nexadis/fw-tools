/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"errors"
	"log"
	"os/signal"
	"syscall"

	"github.com/spf13/cobra"

	"github.com/Nexadis/fw-tools/internal/swap"
)

// swapCmd represents the swap command
var swapCmd = &cobra.Command{
	Use:   "swap filename [filename2]...",
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
		s := swap.Swapper{
			Config: cfg.Swap,
		}
		err := s.Open(cfg.Inputs)
		if err != nil {
			log.Fatal(err)
		}
		defer s.Close()
		ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
		defer cancel()
		err = s.Run(ctx)
		if err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	swapCmd.Flags().BoolVarP(&cfg.Swap.Bits, "bits", "", false, "Inverse bits in byte")
	swapCmd.Flags().BoolVarP(&cfg.Swap.Halfs, "halfs", "", false, "Swap halfs of byte")
	swapCmd.Flags().BoolVarP(&cfg.Swap.Bytes, "bytes", "b", false, "Swap neighbors bytes")
	swapCmd.Flags().BoolVarP(&cfg.Swap.Words, "words", "w", false, "Swap neighbors words")
	swapCmd.Flags().BoolVarP(&cfg.Swap.Dwords, "dwords", "d", false, "Swap neighbors dwords")

	rootCmd.AddCommand(swapCmd)

}
