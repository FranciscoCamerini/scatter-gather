package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"net"
	"os"
	"server/server"
	"strings"
)

var (
	worker server.Server
)

func handleConnection(conn net.Conn) {
	defer conn.Close()

	reader := bufio.NewReader(conn)
	message, err := reader.ReadString('\n')
	if err != nil {
		worker.Log("error reading from connection: %s", err.Error())
		return
	}
	message = strings.TrimSuffix(message, "\n")

	worker.Log("received message: \"%s\"", message)

	words := strings.Split(message, ",")
	response := parseFiles(words)
	if response == "" {
		return
	}
	conn.Write([]byte(fmt.Sprintf("%s\n", response)))
}

func parseFiles(words []string) string {
	entries, err := os.ReadDir("./public")
	if err != nil {
		worker.Log("error reading from public dir: %s", err.Error())
		return ""
	}

	wordMap := make(map[string]map[string]int)
	for i := 0; i < len(words); i++ {
		wordMap[strings.ToLower(words[i])] = make(map[string]int)
	}

	for _, e := range entries {
		file, err := os.Open(fmt.Sprintf("./public/%s", e.Name()))
		if err != nil {
			worker.Log("error opening file: %s", err.Error())
			continue
		}

		scanner := bufio.NewScanner(file)
		scanner.Split(bufio.ScanWords)

		for scanner.Scan() {
			word := strings.ToLower(scanner.Text())
			if wordMap[word] != nil {
				wordMap[word][e.Name()]++
			}
		}
		file.Close()
	}

	json, err := json.Marshal(wordMap)
	if err != nil {
		worker.Log("error converting to JSON: %s", err.Error())
		return ""
	}

	return string(json)
}

func main() {
	var (
		port    int
		pidFile string
	)

	flag.IntVar(&port, "port", 8081, "Sets the port number to listen to")
	flag.StringVar(&pidFile, "pidfile", "", "Sets the pidfile to write to")
	flag.Parse()

	worker = server.Server{
		Port:     port,
		PIDFile:  pidFile,
		Name:     fmt.Sprintf("worker-%d", port),
		LogColor: "\u001B[34m",
	}
	worker.Run(handleConnection)
}
