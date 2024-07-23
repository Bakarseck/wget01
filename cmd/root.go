package cmd

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/spf13/cobra"
)

var (
	output    string
	filePath  string
	rateLimit string
	mirror    bool
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "wget01",
	Short: "A wget clone implemented in Go",
	Long:  `This project aims to recreate some functionalities of wget using the Go programming language.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			fmt.Println("URL is required")
			os.Exit(1)
		}

		url := args[0]
		startTime := time.Now().Format("2006-01-02 15:04:05")
		fmt.Println("start at", startTime)

		response, err := http.Get(url)
		if err != nil {
			fmt.Println("Error:", err)
			os.Exit(1)
		}
		defer response.Body.Close()

		if response.StatusCode != http.StatusOK {
			fmt.Println("Error: status", response.Status)
			os.Exit(1)
		}

		if output == "" {
			output = "downloaded_file"
		}

		outFile, err := os.Create(output)
		if err != nil {
			fmt.Println("Error:", err)
			os.Exit(1)
		}
		defer outFile.Close()

		_, err = io.Copy(outFile, response.Body)
		if err != nil {
			fmt.Println("Error:", err)
			os.Exit(1)
		}

		endTime := time.Now().Format("2006-01-02 15:04:05")
		fmt.Println("finished at", endTime)
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().StringVarP(&output, "output", "O", "", "Save the downloaded file under a different name")
	rootCmd.Flags().StringVarP(&filePath, "path", "P", "", "Specify the directory to save the downloaded file")
	rootCmd.Flags().StringVarP(&rateLimit, "rate-limit", "r", "", "Limit the download speed (e.g., 400k or 2M)")
	rootCmd.Flags().BoolVarP(&mirror, "mirror", "m", false, "This option should download the entire website being")
}
