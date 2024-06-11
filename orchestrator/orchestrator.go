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
	orchestrator server.Server
	workerPorts  []int
)

func handleConnection(conn net.Conn) {
	reader := bufio.NewReader(conn)
	defer conn.Close()

	for {
		b, _, err := reader.ReadLine()
		if err != nil {
			orchestrator.Log("error reading from connection: %s", err.Error())
			return
		}

		message := string(b)
		orchestrator.Log("received message: \"%s\"", message)

		responseChannel := make(chan []byte, len(workerPorts))
		words := strings.Split(message, " ")
		scatterMessage(words, responseChannel)

		responseCount := 0
		conn.Write([]byte("\n"))
		for response := range responseChannel {
			responseCount++

			var responseData map[string]map[string]int
			err = json.Unmarshal([]byte(response), &responseData)
			if err != nil {
				orchestrator.Log("error parsing JSON: %s", err.Error())
				break
			}

			for word, appearances := range responseData {
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

			if responseCount == len(workerPorts) || responseCount == len(words) {
				break
			}
		}
	}
}

func scatterMessage(words []string, responseChannel chan<- []byte) {
	wordsPerWorker := len(words) / len(workerPorts)
	if wordsPerWorker == 0 {
		wordsPerWorker = 1
	}

	wordsProcessed := 0
	for i := 0; i < len(workerPorts); i++ {
		if i*wordsPerWorker >= len(words) {
			break
		}

		startIdx := i * wordsPerWorker
		endIdx := i*wordsPerWorker + wordsPerWorker

		if i == len(workerPorts)-1 {
			if wordsPerWorker < len(words)-wordsProcessed {
				endIdx += len(words) - wordsProcessed - wordsPerWorker
			}
		}
		wordsProcessed += wordsPerWorker

		go dialWorker(workerPorts[i], strings.Join(words[startIdx:endIdx][:], ","), responseChannel)
	}
}

func dialWorker(port int, words string, responseChannel chan<- []byte) {
	orchestrator.Log("sending \"%s\" to %d", words, port)

	conn, err := net.Dial("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		orchestrator.Log("error dialing worker: %s", err.Error())
		responseChannel <- nil
		return
	}
	defer conn.Close()

	_, err = conn.Write([]byte(fmt.Sprintf("%s\n", words)))
	if err != nil {
		orchestrator.Log("error writing to worker: %s", err.Error())
		responseChannel <- nil
		return
	}

	reader := bufio.NewReader(conn)
	response, _, err := reader.ReadLine()
	if err != nil {
		orchestrator.Log("error reading response: %s", err.Error())
		responseChannel <- nil
	} else {
		orchestrator.Log("response from %d: \"%s\"", port, string(response))
		responseChannel <- response
	}
}

func main() {
	var (
		port    int
		workers string
		pidFile string
	)

	flag.IntVar(&port, "port", 8080, "Sets the port number to listen to")
	flag.StringVar(&workers, "workers", "8081,8081", "Sets ports to be used to spawn workers. E.g.: 8081,8082")
	flag.StringVar(&pidFile, "pidfile", "", "Sets the pidfile to write to")
	flag.Parse()

	for _, portString := range strings.Split(workers, ",") {
		port, err := strconv.Atoi(portString)
		if err != nil {
			orchestrator.Log("%s: %s", err.Error(), portString)
			os.Exit(1)
		}

		workerPorts = append(workerPorts, port)
	}

	orchestrator = server.Server{
		Port:     port,
		PIDFile:  pidFile,
		Name:     "orchestrator",
		LogColor: "\u001B[32m",
	}
	orchestrator.Run(handleConnection)
}
