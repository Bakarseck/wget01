package cmd

import (
    "fmt"
    "io"
    "net/http"
    "os"
	"os/exec"
    "strings"
    "time"

    "github.com/spf13/cobra"
)

var (
    output     string
    filePath   string
    rateLimit  string
    background bool
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
    Use:   "wget",
    Short: "A wget clone implemented in Go",
    Long:  `This project aims to recreate some functionalities of wget using the Go programming language.`,
    Run: func(cmd *cobra.Command, args []string) {
        handleArguments(args)
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
    rootCmd.Flags().BoolVarP(&background, "background", "B", false, "Download the file in the background")
}

func handleArguments(args []string) {
    if len(args) < 1 {
        fmt.Println("URL is required")
        os.Exit(1)
    }

    url := args[0]

    if background {
        downloadInBackground(url)
    } else {
        downloadFile(url)
    }
}

func downloadFile(url string) {
    startTime := time.Now().Format("2006-01-02 15:04:05")
    fmt.Printf("--%s--  %s\n", startTime, url)

    urlParts := strings.Split(url, "/")
    host := urlParts[2]
    fmt.Printf("Resolving %s... ", host)

    response, err := http.Get(url)
    if err != nil {
        fmt.Println("Error:", err)
        os.Exit(1)
    }
    defer response.Body.Close()

    fmt.Printf("Connecting to %s... connected.\n", host)
    fmt.Printf("HTTP request sent, awaiting response... %s\n", response.Status)

    if response.StatusCode != http.StatusOK {
        fmt.Printf("Error: status %s\n", response.Status)
        os.Exit(1)
    }

    contentLength := response.Header.Get("Content-Length")
    length := "-"
    if contentLength != "" {
        length = contentLength
    }
    contentType := response.Header.Get("Content-Type")
    fmt.Printf("Length: %s [%s]\n", length, contentType)

    fileName := output
    if fileName == "" {
        fileName = urlParts[len(urlParts)-1]
        if fileName == "" {
            fileName = "download"
        }
    }
    if filePath != "" {
        fileName = filePath + "/" + fileName
    }
    fmt.Printf("Saving to: ‘%s’\n\n", fileName)

    outFile, err := os.Create(fileName)
    if err != nil {
        fmt.Println("Error:", err)
        os.Exit(1)
    }
    defer outFile.Close()

    buf := make([]byte, 1024)
    var total int64
    for {
        n, err := response.Body.Read(buf)
        if n > 0 {
            outFile.Write(buf[:n])
            total += int64(n)
            fmt.Printf("\r%s  %3d%%", fileName, int(float64(total)/float64(response.ContentLength)*100))
        }
        if err == io.EOF {
            break
        }
        if err != nil {
            fmt.Println("Error:", err)
            os.Exit(1)
        }
    }

    endTime := time.Now().Format("2006-01-02 15:04:05")
    fmt.Printf("\n\n%s (%d bytes) saved [saved %d/%d]\n", fileName, total, total, response.ContentLength)
    fmt.Printf("%s - ‘%s’ saved [%d/%d]\n", endTime, fileName, total, response.ContentLength)
}

func downloadInBackground(url string) {
    cmd := exec.Command(os.Args[0], append(os.Args[1:], "-B")...)
    cmd.Stdout = nil
    cmd.Stderr = nil
    err := cmd.Start()
    if err != nil {
        fmt.Println("Failed to start background download:", err)
        os.Exit(1)
    }
    fmt.Println("Output will be written to \"wget-log\".")
    go func() {
        cmd.Wait()
    }()
}
