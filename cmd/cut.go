/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"errors"
	"log"

	"github.com/spf13/cobra"
	"golang.org/x/sync/errgroup"

	"github.com/Nexadis/fw-tools/internal/cut"
)

// cutCmd represents the cut command
var cutCmd = &cobra.Command{
	Use:   "cut filename",
	Short: "Cut metainfo between pages with step",
	Long:  `In some type of memory dumps we can meet additional meta-information about pages. You can use this command for cut it.`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("set filename")
		}
		cfg.Inputs = append(cfg.Inputs, args...)
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		if len(cfg.Inputs) == 1 {
			c := cut.New(cfg.Cut)
			err := c.Open(cfg.Inputs[0], cfg.Output)
			if err != nil {
				log.Fatal(err)
			}
			defer c.Close()
			err = c.Run(context.Background())
			if err != nil {
				log.Fatal(err)
			}
			return
		}
		errgrp, ctx := errgroup.WithContext(context.Background())
		for _, input := range cfg.Inputs {
			input := input
			errgrp.Go(func() error {
				c := cut.New(cfg.Cut)
				err := c.Open(input, "")
				if err != nil {
					return err
				}
				defer c.Close()
				err = c.Run(ctx)
				if err != nil {
					return err
				}
				return nil

			})

		}
		err := errgrp.Wait()
		if err != nil {
			log.Fatal(err)
		}

	},
}

func init() {
	cutCmd.Flags().IntVarP(&cfg.Cut.Page, "page", "p", 0x400, "Page size, which will writed")
	cutCmd.Flags().IntVarP(&cfg.Cut.Skip, "skip", "s", 0x20, "Metainfo size, which will skipped")
	rootCmd.AddCommand(cutCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// cutCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// cutCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
