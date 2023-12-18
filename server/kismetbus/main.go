package main

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"time"

	"github.com/gorilla/websocket"
)

var (
	apiKey  = os.Getenv("KISMET_SERVER_APIKEY")
	logPath = "/var/log/kismetbus/"

	// Temporary
	// apiKey   = "BB0E469600250E7F999EA5D7DBB7BFD5"
	// url      = "wss://" + "10.133.0.6" + ":" + "2501" + "/eventbus/events.ws"
	// insecure = true
)

type JsonMessage struct {
	Text string `json:"text"`
}

func main() {
	ssl, err := stringToBool(os.Getenv("KISMET_SERVER_SSL"))
	if err != nil {
		log.Fatal(err)
	}
	secure := "ws"
	if ssl {
		secure = "wss"
	}
	url := secure + "://" + os.Getenv("KISMET_SERVER_ADDRESS") + ":" + os.Getenv("KISMET_SERVER_PORT") + "/eventbus/events.ws"
	os.Mkdir("/var/log/kismetbus", 0644)

	authCookie := &http.Cookie{
		Name:  "KISMET",
		Value: apiKey,
	}
	headers := make(http.Header)
	headers.Add("Cookie", authCookie.String())

	insecure, err := stringToBool(os.Getenv("KISMET_SERVER_INSECURE"))
	if err != nil {
		log.Println(err)
	}

	dialer := websocket.DefaultDialer
	if insecure && ssl {
		dialer.TLSClientConfig = &tls.Config{InsecureSkipVerify: insecure}
	}
	log.Println("Connecting to:" + url)
	c, _, err := dialer.Dial(url, headers)
	if err != nil {
		log.Fatal("Dial:", err)
	}
	defer c.Close()

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	done := make(chan struct{})

	go func() {
		defer close(done)
		for {
			currentTime := time.Now()
			date := currentTime.Format(time.DateOnly)
			_, message, err := c.ReadMessage()
			if err != nil {
				log.Println("Read:", err)
				return
			}
			logtofile(filepath.Join(logPath, date+"-"+"message.log"), message)
			log.Println("Received: ", string(message))
		}
	}()

	// Read JSON message from file
	jsonMessage, err := readJSONFromFile("/etc/kismetbus/eventbus.json")
	if err != nil {
		log.Fatal("Error reading JSON from file:", err)
	}
	// Send the JSON message
	err = c.WriteMessage(websocket.TextMessage, jsonMessage)
	if err != nil {
		log.Fatal("Error writing JSON message:", err)
	}
	log.Println("Sending subscribe JSON...")

	for {
		select {
		case <-done:
			return
		case <-interrupt:
			log.Println("interrupt")
			err := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				log.Println("write close:", err)
			}

			// Check if done channel is closed
			_, ok := <-done
			if ok {
				return
			}
		}
	}
}

func logtofile(logfile string, data []byte) {
	f, err := os.OpenFile(logfile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	f.Write(data)
	f.Write([]byte("\n"))
}

func readJSONFromFile(filename string) ([]byte, error) {
	fileContent, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return fileContent, nil
}

func stringToBool(s string) (bool, error) {
	switch strings.ToLower(s) {
	case "true", "t", "yes", "y", "1":
		return true, nil
	case "false", "f", "no", "n", "0":
		return false, nil
	default:
		return false, fmt.Errorf("invalid boolean string: %s", s)
	}
}
