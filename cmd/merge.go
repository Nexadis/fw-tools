/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"errors"
	"log"

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
		err, closeFunc := m.Open(cfg.Inputs, cfg.Output)
		if err != nil {
			log.Fatal(err)
		}
		defer closeFunc()
		err = m.Run(context.TODO())
		if err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	mergeCmd.Flags().BoolVarP(&cfg.Merge.ByBit, "bits", "", false, "Merge by bits in byte")
	mergeCmd.Flags().BoolVarP(&cfg.Merge.ByByte, "bytes", "b", false, "Merge by bytes")
	mergeCmd.Flags().BoolVarP(&cfg.Merge.ByWord, "words", "w", false, "Merge by word")
	mergeCmd.Flags().BoolVarP(&cfg.Merge.ByDword, "dwords", "d", false, "Merge by dwords")
	mergeCmd.MarkFlagsMutuallyExclusive(mergeCmd.Flags().Args()...)
	rootCmd.AddCommand(mergeCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// mergeCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// mergeCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
