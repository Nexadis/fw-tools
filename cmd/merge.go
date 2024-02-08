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

	"github.com/Nexadis/fw-tools/internal/merge"
)

// mergeCmd represents the merge command
var mergeCmd = &cobra.Command{
	Use:   "merge filename filename2 ...",
	Short: "Merge some dumps in defined order of bytes.",
	Long: `For example you can merge dump1 and dump2 byte by byte or word by word or something else.
	If you need you can merge 2 and more dumps in one. Example:

	dump1 					: 0A 		0B 		01 		02
	dump2 					: 	0C 		0D 		03 		04
	dump3 					: 		0E 		0F 		05 		06

	merged by byte 	: 0A0C0E0B0D0F010305020406
	`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 2 {
			return errors.New("set filenames")
		}
		cfg.Inputs = append(cfg.Inputs, args...)
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		m := merge.New(cfg.Merge)
		err := m.Open(cfg.Inputs, cfg.Merge.Output)
		if err != nil {
			log.Fatal(err)
		}
		defer m.Close()
		ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
		defer cancel()
		err = m.Run(ctx)
		if err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	bits := "bits"
	b := "bytes"
	w := "words"
	d := "dwords"
	o := "output"
	mergeCmd.Flags().BoolVarP(&cfg.Merge.ByBit, bits, "", false, "Merge by bits in byte")
	mergeCmd.Flags().BoolVarP(&cfg.Merge.ByByte, b, "b", false, "Merge by bytes")
	mergeCmd.Flags().BoolVarP(&cfg.Merge.ByWord, w, "w", false, "Merge by word")
	mergeCmd.Flags().BoolVarP(&cfg.Merge.ByDword, d, "d", false, "Merge by dwords")
	mergeCmd.Flags().StringVarP(&cfg.Merge.Output, o, "o", "merged.bin", "Merge by dwords")
	// you should choose only one flag
	mergeCmd.MarkFlagsMutuallyExclusive(bits, b, w, d)
	rootCmd.AddCommand(mergeCmd)
}
