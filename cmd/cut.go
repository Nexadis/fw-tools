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
		ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
		defer cancel()
		errgrp, ctx := errgroup.WithContext(ctx)
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
	cutCmd.Flags().IntVarP(&cfg.Cut.PageSize, "page", "p", 0x400, "Page size, which will writed")
	cutCmd.Flags().IntVarP(&cfg.Cut.SkipSize, "skip", "s", 0x20, "Metainfo size, which will skipped")
	rootCmd.AddCommand(cutCmd)
}
