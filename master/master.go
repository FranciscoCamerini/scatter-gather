package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"net"
	"os"
	"server/server"
	"strconv"
	"strings"
)

var (
	master       server.Server
	crawlerPorts []int
)

func handleConnection(conn net.Conn) {
	reader := bufio.NewReader(conn)
	defer conn.Close()

	for {
		b, _, err := reader.ReadLine()
		if err != nil {
			master.Log(err.Error())
			return
		}

		msg := string(b)
		master.Log("received message: \"%s\"", msg)

		responseChannel := make(chan map[string]map[string]int, len(crawlerPorts))
		words := strings.Split(msg, " ")
		scatterMessage(words, responseChannel)

		responseCount := 0
		conn.Write([]byte("\n"))
		for response := range responseChannel {
			responseCount++

			for word, appearances := range response {
				if len(appearances) > 0 {
					conn.Write([]byte(fmt.Sprintf("\u001B[32m%s:\u001B[0m\n", word)))
					for file, count := range appearances {
						conn.Write([]byte(fmt.Sprintf("- File: %s. Count: %d\n\n", file, count)))
					}
				} else {
					conn.Write([]byte(fmt.Sprintf("\u001B[31m%s:\u001B[0m\n", word)))
					conn.Write([]byte("- Not found.\n\n"))
				}
			}

			if responseCount == len(crawlerPorts) || responseCount == len(words) {
				break
			}
		}
	}
}

func scatterMessage(words []string, responseChannel chan<- map[string]map[string]int) {
	wordsPerCrawler := len(words) / len(crawlerPorts)
	if wordsPerCrawler == 0 {
		wordsPerCrawler = 1
	}

	wordsProcessed := 0
	for i := 0; i < len(crawlerPorts); i++ {
		if i*wordsPerCrawler >= len(words) {
			break
		}

		startIdx := i * wordsPerCrawler
		endIdx := i*wordsPerCrawler + wordsPerCrawler

		if i == len(crawlerPorts)-1 {
			if wordsPerCrawler < len(words)-wordsProcessed {
				endIdx += len(words) - wordsProcessed - wordsPerCrawler
			}
		}
		wordsProcessed += wordsPerCrawler

		go dialCrawler(crawlerPorts[i], strings.Join(words[startIdx:endIdx][:], ","), responseChannel)
	}
}

func dialCrawler(port int, words string, responseChannel chan<- map[string]map[string]int) {
	master.Log("sending \"%s\" to %d", words, port)

	conn, err := net.Dial("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		master.Log(err.Error())
		return
	}

	_, err = conn.Write([]byte(fmt.Sprintf("%s\n", words)))
	if err != nil {
		master.Log(err.Error())
		return
	}

	reader := bufio.NewReader(conn)
	msg, _ := reader.ReadString('\n')
	response := strings.TrimSuffix(msg, "\n")

	master.Log("response from %d: \"%s\"", port, response)

	var data map[string]map[string]int
	err = json.Unmarshal([]byte(response), &data)
	if err != nil {
		master.Log("Error parsing JSON: %s", err.Error())
		return
	}

	responseChannel <- data
}

func main() {
	var (
		port             int
		crawlerPortsFlag string
		pidFile          string
	)

	flag.IntVar(&port, "port", 8080, "Sets the port number to listen to")
	flag.StringVar(&crawlerPortsFlag, "crawlers", "8081,8081", "Sets ports to be used to spawn crawlers")
	flag.StringVar(&pidFile, "pidfile", "", "Sets the pidfile to write to")
	flag.Parse()

	for _, portString := range strings.Split(crawlerPortsFlag, ",") {
		port, err := strconv.Atoi(portString)
		if err != nil {
			master.Log("%s: %s", err.Error(), portString)
			os.Exit(1)
		}

		crawlerPorts = append(crawlerPorts, port)
	}

	master = server.Server{
		Port:     port,
		PIDFile:  pidFile,
		Name:     "master",
		LogColor: "\u001B[32m",
	}
	master.Run(handleConnection)
}
