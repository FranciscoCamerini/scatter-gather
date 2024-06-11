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
	crawler server.Server
)

func handleConnection(conn net.Conn) {
	defer conn.Close()

	reader := bufio.NewReader(conn)
	msg, err := reader.ReadString('\n')
	msg = strings.TrimSuffix(msg, "\n")

	if err != nil {
		crawler.Log(err.Error())
		return
	}

	crawler.Log("received message: \"%s\"", msg)
	response := parseFiles(strings.Split(msg, ","))
	if response == "" {
		return
	}
	conn.Write([]byte(fmt.Sprintf("%s\n", response)))
}

func parseFiles(words []string) string {
	entries, err := os.ReadDir("./public")
	if err != nil {
		crawler.Log(err.Error())
		return ""
	}

	wordMap := make(map[string]map[string]int)
	for i := 0; i < len(words); i++ {
		wordMap[strings.ToLower(words[i])] = make(map[string]int)
	}

	for _, e := range entries {
		file, err := os.Open(fmt.Sprintf("./public/%s", e.Name()))
		if err != nil {
			crawler.Log("Error opening file: %s", err.Error())
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
		crawler.Log("Error converting to JSON: %s", err.Error())
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

	crawler = server.Server{
		Port:     port,
		PIDFile:  pidFile,
		Name:     fmt.Sprintf("crawler-%d", port),
		LogColor: "\u001B[34m",
	}
	crawler.Run(handleConnection)
}
