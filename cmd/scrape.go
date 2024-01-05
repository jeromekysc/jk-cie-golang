/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/safetyculture/ci-golang/internal/scraper"
	"github.com/spf13/cobra"
)

// scrapeCmd represents the scrape command
var scrapeCmd = &cobra.Command{
	Use:   "scrape",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	RunE: func(cmd *cobra.Command, args []string) error {

		inc, err := cmd.Flags().GetStringSlice("target")

		if err != nil {
			return fmt.Errorf("error loading targets, %s", err)
		}
		targets := []scraper.Target{}

		for _, v := range inc {
			target, err := scraper.ConvertToTarget(v)

			if err != nil {
				fmt.Printf("error to convert this target")
				continue
			}

			targets = append(targets, *target)
		}

		return scraper.Scrape(os.Stdout, targets)
	},
}

func init() {
	scrapeCmd.Flags().StringSlice("target", []string{}, "List of targets (METHOD|URL)")
	rootCmd.AddCommand(scrapeCmd)
}
