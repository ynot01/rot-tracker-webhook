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

// udpClients["0.0.0.0:7777"] = *a2s.Client
var udpClients map[string]*a2s.Client = make(map[string]*a2s.Client)

// udpClients["0.0.0.0:7777"] = "SIGNAL 13 1 2025 CLUSTER"
var registeredServers map[string]string = make(map[string]string)

// https://discord.com/api/webhooks/<webhook id>/<webhook token>
var myWebhookURL string

// Loops through all official servers and queries them one by one
// If any change in name, consider them rotted and post to Discord
func main() {
	fmt.Printf("[%v] Started rot-tracker-webhook.\n", time.Now().Format(time.RFC850))
	content, fileErr := os.ReadFile("webhook.txt")
	if fileErr != nil {
		fmt.Printf("[%v] fileErr: %v\n", time.Now().Format(time.RFC850), fileErr)
		return
	}
	myWebhookURL = strings.TrimSpace(string(content))
	if myWebhookURL == "Replace this text with your Discord webhook URL" || myWebhookURL == "" {
		fmt.Printf("[%v] Webhook not set! Please set webhook.txt\n", time.Now().Format(time.RFC850))
		return
	}
	for range time.Tick(time.Second * 10) { // Wait a healthy 10 seconds
		official_servers := get_masterlist()
		for official := range official_servers {
			ipAddr := official_servers[official]
			if strings.TrimSpace(ipAddr) == "" {
				continue
			}
			findComment := strings.Index(ipAddr, "//") // Strip comments from masterlist
			if findComment != -1 {
				ipAddr = strings.TrimSpace(ipAddr[:findComment])
			}
			dictKey := ipAddr
			findPort := strings.Index(ipAddr, ":") // Isolate :port (and ignore entries without a port)
			var ipPort string
			if findPort == -1 {
				ipPort = "7777"
			} else {
				ipPort = strings.TrimSpace(ipAddr[findPort:])[1:]
			}
			portInt, atoiErr := strconv.Atoi(ipPort)
			if atoiErr != nil {
				fmt.Printf("[%v] atoiErr: %v\n", time.Now().Format(time.RFC850), atoiErr)
				continue
			}
			if findPort == -1 {
				ipAddr = fmt.Sprintf("%v:%v", ipAddr, portInt+1)
			} else {
				ipAddr = fmt.Sprintf("%v:%v", ipAddr[:findPort], portInt+1) // Add 1 to server port to get the a2s query port
			}
			client, weHaveClient := udpClients[dictKey]
			if !weHaveClient {
				newClient, newClientErr := a2s.NewClient(
					ipAddr,
					a2s.SetAppID(2773280),
				)
				if newClientErr != nil {
					fmt.Printf("[%v] newClientErr: %v\n", time.Now().Format(time.RFC850), newClientErr)
					continue
				}
				client = newClient
				udpClients[dictKey] = newClient
			}
			info, infoErr := client.QueryInfo()
			if infoErr != nil {
				// Don't print on server connection errors- a few of them are down a lot
				// fmt.Printf("%v \"fail\"\n", ipAddr)
				continue
			} else {
				// fmt.Printf("%v \"success\"\n", ipAddr)
			}
			// Servers are stored by server port, not query port!
			oldServerName, serverIsRegistered := registeredServers[dictKey]
			if serverIsRegistered && oldServerName != info.Name { // If the name changed, report it to Discord
				region := strings.ToUpper(get_region_from_keywords(info.ExtendedServerInfo.Keywords))
				send_message_to_discord(dictKey, region, oldServerName, info.Name)
			}
			registeredServers[dictKey] = info.Name
		}
	}
}

// Download updated official server list and split by newline
func get_masterlist() []string {
	resp, masterErr := http.Get(MASTER_URL)
	if masterErr != nil {
		fmt.Printf("[%v] masterErr: %v\n", time.Now().Format(time.RFC850), masterErr)
		return []string{}
	}
	resBody, ioErr := io.ReadAll(resp.Body)
	if ioErr != nil {
		fmt.Printf("[%v] ioErr: %v\n", time.Now().Format(time.RFC850), ioErr)
		return []string{}
	}
	return strings.Split(string(resBody), "\n")
}

// Http POST request to Discord webhook
func send_message_to_discord(ipAddr string, region string, oldServerName string, newServerName string) {
	jsonBody := fmt.Appendf(nil, `{
  "embeds": [
    {
      "title": "An official server has been consumed by the rot!",
      "color": 16749144,
      "fields": [
        {
          "name": "Region",
          "value": "%v"
        },
        {
          "name": "Old name",
          "value": "%v"
        },
        {
          "name": "New name",
          "value": "%v"
        }
      ],
      "footer": {
        "text": "IP Address: %v"
      }
    }
  ],
  "components": [],
  "username": "%v",
  "avatar_url": "%v"
}`, region, oldServerName, newServerName, ipAddr, WEBHOOK_USERNAME, WEBHOOK_AVATAR_URL)
	bodyReader := bytes.NewReader(jsonBody)
	resp, err := http.Post(myWebhookURL, "application/json", bodyReader)
	if err != nil {
		fmt.Printf("[%v] err: %v\n", time.Now().Format(time.RFC850), err)
		return
	}
	if resp.StatusCode > 299 {
		fmt.Printf("[%v] resp.StatusCode: %v\n", time.Now().Format(time.RFC850), resp.StatusCode)
		return
	}
}

// String operations to isolate region
// region:us,uptime:0,protected:1,rot:0.06 = us
// region:au,uptime:30,protected:1,rot:0.05 = au
// region:cn,uptime:30,protected:1,rot:0.06 = cn
func get_region_from_keywords(keywords string) string {
	stringRegion := "Unknown Region"
	keywordStrings := strings.Split(strings.TrimSpace(keywords), ",")
	for n := range keywordStrings {
		keyAndValue := strings.Split(keywordStrings[n], ":")
		key := keyAndValue[0]
		value := keyAndValue[1]
		if key == "region" {
			stringRegion = value
		}
	}
	switch strings.ToUpper(stringRegion) {
	case "US":
		return "🇺🇸"
	case "AU":
		return "🇦🇺"
	case "CN":
		return "🇨🇳"
	case "RU":
		return "🇷🇺"
	case "EUROPEANUNION":
		return "🇪🇺"
	}
	return stringRegion
}
