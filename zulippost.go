package main

import (
	"bufio"
	"crypto/tls"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
)

func getEnvVar(varName string) (result string) {
	for _, e := range os.Environ() {
		pair := strings.Split(e, "=")
		if pair[0] == varName {
			return pair[1]
		}
	}
	return ""
}

// TODO: fetch from environment variable
const zulipURL = "https://zulip.tok.access-company.com/api/v1/messages"
const emailAddress = "feedfeed-bot@access-company.com"
const apiKey = "jrMQg9AUAQPg4lRb2mi6NCBgdUr5BmUR"

func main() {
	//slackpostWebhookUrl := getEnvVar(ENVKEY_SLACKPOST_WEBHOOK_URL)
	//if slackpostWebhookUrl == "" {
	//	fmt.Println(ENVKEY_SLACKPOST_WEBHOOK_URL, "is not specified.")
	//	os.Exit(1)
	//}

	//slackpostUserName := getEnvVar(ENVKEY_SLACKPOST_USERNAME)
	//if slackpostUserName == "" {
	//	fmt.Println(ENVKEY_SLACKPOST_USERNAME, "is not specified.")
	//	os.Exit(1)
	//}

	//slackpostChannelToPost := getEnvVar(ENVKEY_SLACKPOST_CHANNEL_TO_POST)
	//if slackpostChannelToPost == "" {
	//	fmt.Println(ENVKEY_SLACKPOST_CHANNEL_TO_POST, "is not specified.")
	//	os.Exit(1)
	//}

	in := os.Stdin
	var msg string
	reader := bufio.NewReaderSize(in, 4096)
	for {
		line, _, err := reader.ReadLine()
		if err == io.EOF {
			break
		} else if err != nil {
			fmt.Println("failed to read from stdin. err =", err)
			os.Exit(1)
		}
		msg += string(line) + "\n"
	}

	if len(msg) <= 0 {
		// msg is empty. exit.
		return
	}

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}

	values := url.Values{
		"type":    {"private"},
		"to":      {"yosuke.akatsuka@access-company.com"},
		"content": {msg},
	}

	req, err := http.NewRequest("POST", zulipURL, strings.NewReader(values.Encode()))

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.SetBasicAuth(emailAddress, apiKey)
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	body, _ := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()

	fmt.Println(string(body))
}
