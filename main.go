package main

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/rumblefrog/go-a2s"
)

const WEBHOOK_USERNAME = "RotTracker"
const WEBHOOK_AVATAR_URL = "https://cdn.fastly.steamstatic.com/steamcommunity/public/images/apps/2773280/b1eeb415c1b44677b667de93549594d313e78a8b.jpg"
const MASTER_URL = "https://content.aneurismiv.com/masterlist"

var registeredServers map[string]string = make(map[string]string)
var myWebhookURL string

func main() {
	content, fileErr := os.ReadFile("webhook.txt")
	if fileErr != nil {
		fmt.Printf("fileErr: %v\n", fileErr)
		return
	}
	myWebhookURL = strings.TrimSpace(string(content))
	for range time.Tick(time.Second * 10) { // Wait a healthy 10 seconds
		official_servers := get_masterlist()
		for official := range official_servers {
			ipAddr := official_servers[official]
			findComment := strings.Index(ipAddr, "//") // Strip comments from masterlist
			if findComment != -1 {
				ipAddr = strings.TrimSpace(ipAddr[:findComment])
			}
			findPort := strings.Index(ipAddr, ":") // Isolate :port (and ignore entries without a port)
			var ipPort string
			if findPort == -1 {
				ipPort = "7777"
			} else {
				ipPort = strings.TrimSpace(ipAddr[findPort:])[1:]
			}
			portInt, atoiErr := strconv.Atoi(ipPort)
			if atoiErr != nil {
				fmt.Printf("atoiErr: %v\n", atoiErr)
				continue
			}
			if findPort == -1 {
				ipAddr = fmt.Sprintf("%v:%v", ipAddr, portInt+1)
			} else {
				ipAddr = fmt.Sprintf("%v:%v", ipAddr[:findPort], portInt+1) // Add 1 to server port to get the a2s query port
			}
			client, newClientErr := a2s.NewClient(
				ipAddr,
				a2s.SetAppID(2773280),
			)
			if newClientErr != nil {
				fmt.Printf("newClientErr: %v\n", newClientErr)
				continue
			}
			info, infoErr := client.QueryInfo()
			if infoErr != nil {
				// Don't print on server connection errors- a few of them are down a lot
				// fmt.Printf("%v \"fail\"\n", ipAddr)
				defer client.Close()
				continue
			} else {
				// fmt.Printf("%v \"success\"\n", ipAddr)
			}
			// Servers are stored by server port, not query port!
			oldServerName, serverIsRegistered := registeredServers[official_servers[official]]
			if serverIsRegistered && oldServerName != info.Name { // If the name changed, report it to Discord
				send_message_to_discord(fmt.Sprintf("%v %v just rotted! New name: %v", official_servers[official], oldServerName, info.Name))
			}
			registeredServers[official_servers[official]] = info.Name
			defer client.Close()
		}
	}
}

func get_masterlist() []string {
	resp, masterErr := http.Get(MASTER_URL)
	if masterErr != nil {
		fmt.Printf("masterErr: %v\n", masterErr)
		return []string{}
	}
	resBody, ioErr := io.ReadAll(resp.Body)
	if ioErr != nil {
		fmt.Printf("ioErr: %v\n", ioErr)
		return []string{}
	}
	return strings.Split(string(resBody), "\n")
}

func send_message_to_discord(msg string) {
	jsonBody := fmt.Appendf(nil, `{
		"avatar_url": "%v",
		"username": "%v",
		"content": "%v"
	}`, WEBHOOK_AVATAR_URL, WEBHOOK_USERNAME, msg)
	bodyReader := bytes.NewReader(jsonBody)
	resp, err := http.Post(myWebhookURL, "application/json", bodyReader)
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return
	}
	if resp.StatusCode > 299 {
		fmt.Printf("resp.StatusCode: %v\n", resp.StatusCode)
		return
	}
}
