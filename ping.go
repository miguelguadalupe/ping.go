package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"sync"

	"github.com/logrusorgru/aurora/v3"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run ping.go <host>")
		return
	}

	host := os.Args[1]

	cmd := exec.Command("ping", "-c", "4", host)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()

		scanner := bufio.NewScanner(stdout)
		for scanner.Scan() {
			line := scanner.Text()
			if strings.Contains(line, "icmp_seq") || strings.Contains(line, "rtt") {
				fmt.Println(aurora.Bold(aurora.Green(line)))
			} else if strings.Contains(line, "Destination Host Unreachable") {
				fmt.Println(aurora.Bold(aurora.Red(line)))
			} else if strings.Contains(line, "Request timeout") {
				fmt.Println(aurora.Bold(aurora.BrightBlack(line)))
			} else {
				fmt.Println(line)
			}
		}
	}()

	err = cmd.Start()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	wg.Wait()

	err = cmd.Wait()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
}
