Wget Clone in Go

This project aims to recreate some functionalities of wget using the Go programming language with the Cobra library for handling command-line flags and arguments. The functionalities include downloading files from the web, renaming files, saving to specific directories, limiting download speed, downloading files in the background, handling multiple downloads asynchronously, and mirroring entire websites.
Objectives

The main objectives of this project are:

    Download a file given a URL:

    sh

$ go run . https://some_url.org/file.zip

Download a file and save it under a different name:

sh

$ go run . -O newname.zip https://some_url.org/file.zip

Download and save the file in a specific directory:

sh

$ go run . -P /path/to/save/ https://some_url.org/file.zip

Set the download speed, limiting the rate of a download:

sh

$ go run . --rate-limit 400k https://some_url.org/file.zip

Download a file in the background:

sh

$ go run . -B https://some_url.org/file.zip

Download multiple files at the same time from a list:

sh

$ go run . -i download_list.txt

Download an entire website and mirror it locally:

sh

    $ go run . --mirror https://www.example.com

Installation

    Clone the repository:

    sh

git clone https://github.com/yourusername/wget-clone.git
cd wget-clone

Install dependencies:

sh

    go mod tidy

Usage
Basic Usage

To download a file from a URL:

sh

$ go run . https://some_url.org/file.zip

Flags

    -B : Download a file in the background and redirect output to a log file.
    -O : Save the downloaded file under a different name.
    -P : Specify the directory to save the downloaded file.
    --rate-limit : Limit the download speed (e.g., 400k or 2M).
    -i : Download multiple files from a list.
    --mirror : Mirror an entire website.
    -R : Reject certain file types when mirroring a website.
    -X : Exclude specific directories when mirroring a website.
    --convert-links : Convert links for offline viewing when mirroring a website.

Examples

Download a file and save it with a different name:

sh

$ go run . -O meme.jpg https://pbs.twimg.com/media/EMtmPFLWkAA8CIS.jpg

Save the downloaded file in a specific directory:

sh

$ go run . -P ~/Downloads/ -O meme.jpg https://pbs.twimg.com/media/EMtmPFLWkAA8CIS.jpg

Limit the download speed:

sh

$ go run . --rate-limit 400k https://pbs.twimg.com/media/EMtmPFLWkAA8CIS.jpg

Download multiple files asynchronously:

sh

$ go run . -i download_list.txt

Mirror an entire website:

sh

$ go run . --mirror https://www.example.com

Mirror a website and reject specific file types:

sh

$ go run . --mirror -R jpg,gif https://example.com

Mirror a website and exclude specific directories:

sh

$ go run . --mirror -X /assets,/css https://example.com

Mirror a website and convert links for offline viewing:

sh

$ go run . --mirror --convert-links https://example.com

Implementation
Setting up Cobra

    Install Cobra:

    sh

go get -u github.com/spf13/cobra@latest

Initialize Cobra in your project:

sh

cobra init

Create commands for each functionality. For example, to download a file:

sh

    cobra add download

Example Command

Here's an example of how you might implement the download command in Cobra:

go

// cmd/download.go
package cmd

import (
    "fmt"
    "net/http"
    "os"
    "time"

    "github.com/spf13/cobra"
)

// downloadCmd represents the download command
var downloadCmd = &cobra.Command{
    Use:   "download [url]",
    Short: "Download a file from a URL",
    Args:  cobra.ExactArgs(1),
    Run: func(cmd *cobra.Command, args []string) {
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

        outFile, err := os.Create("downloaded_file")
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

func init() {
    rootCmd.AddCommand(downloadCmd)

    // Add flags here
    downloadCmd.Flags().StringP("output", "O", "", "Save the downloaded file under a different name")
    downloadCmd.Flags().StringP("path", "P", "", "Specify the directory to save the downloaded file")
    downloadCmd.Flags().String("rate-limit", "", "Limit the download speed (e.g., 400k or 2M)")
}

Project Requirements

    At least 4 members are required to start the project.
    Each member must have a minimum audit ratio of 0.5.

Learning Outcomes

This project will help you learn about:

    GNU Wget
    HTTP and FTP protocols
    Algorithms (recursion)
    Website mirroring
    File system operations
    Asynchronous programming in Go

Contributions

If you find any issues or have suggestions for improvements, feel free to submit an issue or propose a change.

This README provides a comprehensive guide to the project, including objectives, installation instructions, usage examples, implementation details using Cobra, project requirements, learning outcomes, and contribution guidelines. Adjust the details as needed to fit your specific implementation.