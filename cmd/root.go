package cmd

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/spf13/cobra"
)

var (
	output     string
	filePath   string
	rateLimit  string
	input      string
	background bool
	content    string
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
	rootCmd.Flags().StringVarP(&input, "input", "i", "", "Downloading different files should be possible asynchronously")
}

func handleArguments(args []string) {
	if input != "" {
		urls := make(map[string]string)
		file, err := os.Open(input)
		if err != nil {
			logEntry(fmt.Sprintf("Error: %s\n", err))
			os.Exit(1)
		}
		scanner := bufio.NewScanner(file)
		lineRead := false
		for scanner.Scan() {
			url := scanner.Text()
			if url != "" {
				urlParts := strings.Split(url, "/")
				file := urlParts[len(urlParts)-1]
				urls[url] = file
				lineRead = true
			}
		}

		if !lineRead {
			logEntry("The file is empty")
			os.Exit(1)
		}

		var wg sync.WaitGroup
		ch := make(chan string)

		for url, filename := range urls {
			wg.Add(1)
			go downloadFilesAsync(url, filename, &wg, ch)
		}

		go func() {
			wg.Wait()
			close(ch)
		}()

		for msg := range ch {
			logEntry(msg)
		}
		return
	}
	if len(args) < 1 {
		fmt.Println("URL is required")
		os.Exit(1)
	}

	url := args[0]

	if background {
		fmt.Println("Output will be written to wget-log")
	}

	downloadFile(url)

	if background {
		writeLog()
	}
}

func writeLog() {
	logFile, err := os.Create("wget-log")
	if err != nil {
		fmt.Println("Failed to create log file:", err)
		os.Exit(1)
	}
	defer logFile.Close()

	_, err = logFile.WriteString(content)
	if err != nil {
		fmt.Println("Failed to write to log file:", err)
		os.Exit(1)
	}
}

func downloadFile(url string) {

	startTime := time.Now().Format("2006-01-02 15:04:05")
	logEntry(fmt.Sprintf("--%s--  %s\n", startTime, url))

	urlParts := strings.Split(url, "/")
	host := urlParts[2]
	logEntry(fmt.Sprintf("Resolving %s... ", host))

	response, err := http.Get(url)
	if err != nil {
		logEntry(fmt.Sprintf("Error: %s\n", err))
		os.Exit(1)
	}
	defer response.Body.Close()

	logEntry(fmt.Sprintf("Connecting to %s... connected.\n", host))
	logEntry(fmt.Sprintf("HTTP request sent, awaiting response... %s\n", response.Status))

	if response.StatusCode != http.StatusOK {
		logEntry(fmt.Sprintf("Error: status %s\n", response.Status))
		os.Exit(1)
	}

	contentLength := response.Header.Get("Content-Length")
	length := "-"
	if contentLength != "" {
		length = contentLength
	}
	contentType := response.Header.Get("Content-Type")
	logEntry(fmt.Sprintf("Length: %s [%s]\n", length, contentType))

	fileName := output
	if fileName == "" {
		fileName = urlParts[len(urlParts)-1]
		if fileName == "" {
			fileName = "download"
		}
	}
	if filePath != "" {
		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			err := os.MkdirAll(filePath, 0755)
			if err != nil {
				logEntry(fmt.Sprintf("Error creating directory: %s\n", err))
				os.Exit(1)
			}
		}
		fileName = filePath + "/" + fileName
	}
	var limitInt int64
	if rateLimit != "" {
		match, _ := regexp.MatchString(`^\d+[kKmM]$`, rateLimit)

		if !match {
			logEntry("Error rate limit value: (e.g.,200K, 400k or 2M)\n")
			os.Exit(1)
		}
		limitStr := rateLimit[:len(rateLimit)-1]
		unit := rateLimit[len(rateLimit)-1:]
		limitInt, _ = strconv.ParseInt(limitStr, 10, 64)
		if unit == "k" || unit == "K" {
			limitInt *= 1000
		} else {
			limitInt *= 10000
		}

		fmt.Println(limitInt, unit)
	}

	logEntry(fmt.Sprintf("Saving to: ‘%s’\n\n", fileName))

	outFile, err := os.Create(fileName)
	if err != nil {
		logEntry(fmt.Sprintf("Error: %s\n", err))
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
			if rateLimit != "" && total > limitInt {
				logEntry(fmt.Sprintf("\nSorry this file exceed octed rate of the bandwidth, limited at: %d octed\n",limitInt ))
				os.Exit(1)
			}
			logEntry(fmt.Sprintf("\r%s  %3d%%", fileName, int(float64(total)/float64(response.ContentLength)*100)))
		}

		if err == io.EOF {
			break
		}
		if err != nil {
			logEntry(fmt.Sprintf("Error: %s\n", err))
			os.Exit(1)
		}
	}

	endTime := time.Now().Format("2006-01-02 15:04:05")
	logEntry(fmt.Sprintf("\n\n%s (%d bytes) saved [saved %d/%d]\n", fileName, total, total, response.ContentLength))
	logEntry(fmt.Sprintf("%s - ‘%s’ saved [%d/%d]\n", endTime, fileName, total, response.ContentLength))
}

func logEntry(entry string) {
	if background {
		content += entry
	} else {
		fmt.Print(entry)
	}
}

func downloadFilesAsync(url string, filename string, wg *sync.WaitGroup, ch chan<- string) {
	defer wg.Done()

	// Envoyer une requête GET
	resp, err := http.Get(url)
	if err != nil {
		ch <- fmt.Sprintf("Failed to download %s\n: %v", url, err)
		return
	}
	defer resp.Body.Close()
	// Créer un fichier pour écrire le contenu

	file, err := os.Create(filename)
	if err != nil {
		ch <- fmt.Sprintf("Failed to create file %s: %v\n", filename, err)
		return
	}
	defer file.Close()

	// Copier le contenu de la réponse HTTP dans le fichier
	_, err1 := io.Copy(file, resp.Body)
	if err1 != nil {
		ch <- fmt.Sprintf("Failed to write to file %s: %v\n", filename, err)
		return
	}

	ch <- fmt.Sprintf("Successfully downloaded %s\n", filename)
}
