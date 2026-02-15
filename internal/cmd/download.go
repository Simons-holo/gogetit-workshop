package cmd

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/anxkhn/gogetit-workshop/internal/config"
	"github.com/anxkhn/gogetit-workshop/internal/downloader"
	"github.com/anxkhn/gogetit-workshop/internal/progress"
	"github.com/spf13/cobra"
)

var (
	outputDir   string
	concurrency int
	timeout     int
	retry       int
)

var downloadCmd = &cobra.Command{
	Use:   "download <url> [url...]",
	Short: "Download files concurrently",
	Long: `Download one or more files concurrently with progress tracking.
Supports resumable downloads and automatic retries.`,
	Args: cobra.MinimumNArgs(1),
	Run:  runDownload,
}

func init() {
	rootCmd.AddCommand(downloadCmd)
	downloadCmd.Flags().StringVarP(&outputDir, "output-dir", "o", ".", "Output directory for downloaded files")
	downloadCmd.Flags().IntVarP(&concurrency, "concurrency", "c", 3, "Number of concurrent downloads")
	downloadCmd.Flags().IntVarP(&timeout, "timeout", "t", 30, "Timeout in seconds for each download")
	downloadCmd.Flags().IntVarP(&retry, "retry", "r", 3, "Number of retry attempts")
}

func runDownload(cmd *cobra.Command, args []string) {
	cfg := config.Get()
	cfg.OutputDir = outputDir
	cfg.Concurrency = concurrency
	cfg.Timeout = timeout
	cfg.Retry = retry

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-sigChan
		fmt.Println("\nReceived interrupt signal, cancelling...")
		cancel()
	}()

	downloaderCfg := &downloader.Config{
		OutputDir:   cfg.OutputDir,
		Concurrency: cfg.Concurrency,
		Timeout:     cfg.Timeout,
		Retry:       cfg.Retry,
		UserAgent:   cfg.UserAgent,
	}
	pool := downloader.NewPool(downloaderCfg)
	prog := progress.New(len(args))

	go prog.Start()

	results := pool.Download(ctx, args)

	prog.Stop()

	successCount := 0
	for _, result := range results {
		if result.Error == nil {
			successCount++
			fmt.Printf("Downloaded: %s -> %s\n", result.URL, result.FilePath)
		} else {
			fmt.Printf("Failed: %s: %v\n", result.URL, result.Error)
		}
	}

	fmt.Printf("\nCompleted: %d/%d files downloaded successfully\n", successCount, len(args))
}
