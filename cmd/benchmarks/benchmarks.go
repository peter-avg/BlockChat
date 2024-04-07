package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
)

func main() {
	numberOfClients := flag.Int("numberOfClients", 0, "Number of clients")

	flag.Parse()

	if *numberOfClients == 0 || (*numberOfClients != 5 && *numberOfClients != 10) {
		fmt.Println("Usage: ./benchmark -numberOfClients <number> ,(<number> equals to 5 or 10)")
		os.Exit(1)
	}

	var dirName string
	if *numberOfClients == 5 {
		dirName = "5-nodes"
	} else if *numberOfClients == 10 {
		dirName = "10-nodes"
	}
	ports := make([]int, *numberOfClients)

	// Populate the ports array with ports starting from 5000
	for i := 0; i < *numberOfClients; i++ {
		ports[i] = 5000 + i
	}
	wd, err := os.Getwd()
	if err != nil {
		fmt.Println("Error getting current working directory:", err)
		return
	}
	dir := filepath.Join(wd, "..", "..", "resources", dirName)
	var wg sync.WaitGroup

	for i := 0; i < *numberOfClients; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			fileName := fmt.Sprintf("trans%d.txt", i)
			filePath := filepath.Join(dir, fileName)

			// Open the file for reading
			file, err := os.Open(filePath)
			if err != nil {
				fmt.Printf("Error opening file %s: %s\n", fileName, err)
				return
			}
			defer file.Close()
			if i == 0 {
				scanner := bufio.NewScanner(file)

				for scanner.Scan() {
					line := scanner.Text()
					lineValue := line
					parts := strings.SplitN(lineValue, " ", 2)

					var command = exec.Command("cli", "--port", strconv.Itoa(ports[i]), "-t", parts[0], parts[1])
					//output, err := command.Output()
					//if err != nil {
					//	fmt.Printf("Error getting output for file %s, line %s: %s\n", fileName, line, err)
					//}
					//fmt.Printf("Output for file %s, line %s: %s\n", fileName, line, string(output))
					command.Stdout = os.Stdout
					command.Stderr = os.Stderr
					err := command.Run()
					if err != nil {
						log.Fatalf("cmd.Run() failed with %s\n", err)
					}
				}

				if err := scanner.Err(); err != nil {
					fmt.Printf("Error scanning file %s: %s\n", fileName, err)
					return
				}
			}
		}(i)
	}
	wg.Wait()
}
