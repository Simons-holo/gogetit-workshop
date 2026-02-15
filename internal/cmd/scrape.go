package cmd

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/anxkhn/gogetit-workshop/internal/scraper"
	"github.com/spf13/cobra"
)

var jsonOutput bool

var scrapeCmd = &cobra.Command{
	Use:   "scrape <url>",
	Short: "Scrape website metadata",
	Long: `Scrape a website and extract metadata including title, description,
links, images, and other Open Graph data.`,
	Args: cobra.ExactArgs(1),
	Run:  runScrape,
}

func init() {
	rootCmd.AddCommand(scrapeCmd)
	scrapeCmd.Flags().BoolVar(&jsonOutput, "json", false, "Output in JSON format")
}

func runScrape(cmd *cobra.Command, args []string) {
	ctx := context.Background()

	s := scraper.New()
	metadata, err := s.Scrape(ctx, args[0])
	if err != nil {
		fmt.Fprintf(cmd.ErrOrStderr(), "Error scraping: %v\n", err)
		return
	}

	if jsonOutput {
		enc := json.NewEncoder(cmd.OutOrStdout())
		enc.SetIndent("", "  ")
		enc.Encode(metadata)
		return
	}

	fmt.Printf("Title: %s\n", metadata.Title)
	fmt.Printf("Description: %s\n", metadata.Description)
	fmt.Printf("Links (%d):\n", len(metadata.Links))
	for _, link := range metadata.Links {
		fmt.Printf("  - %s\n", link)
	}
	fmt.Printf("Images (%d):\n", len(metadata.Images))
	for _, img := range metadata.Images {
		fmt.Printf("  - %s\n", img)
	}
}
