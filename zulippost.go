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

//const (
//zulipURL     = "https://zulip.tok.access-company.com/api/v1/messages"
//emailAddress = "feed-feed-bot@access-company.com"
//apiKey       = "jrMQg9AUAQPg4lRb2mi6NCBgdUr5BmUR"

//messageType = "private"
//to          = "yosuke.akatsuka@access-company.com"

//messageType = "stream"
//to          = "EngineerAll"
//subject     = "RSS"
//)

const (
	envZulipURL     = "ZULIPPOST_URL"
	envEmailAddress = "ZULIPPOST_EMAIL"
	envAPIKey       = "ZULIPPOST_APIKEY"
	envMessageType  = "ZULIPPOST_MESSAGETYPE" // "private" or "stream"
	envTo           = "ZULIPPOST_TO"          // "email" or "stream name"
	envSubject      = "ZULIPPOST_SUBJECT"     // "subject" valid if type is stream
)

func main() {
	zulipURL := getEnvVar(envZulipURL)
	if zulipURL == "" {
		fmt.Println(envZulipURL, "is not specified.")
		os.Exit(1)
	}

	emailAddress := getEnvVar(envEmailAddress)
	if emailAddress == "" {
		fmt.Println(envEmailAddress, "is not specified.")
		os.Exit(1)
	}

	apiKey := getEnvVar(envAPIKey)
	if apiKey == "" {
		fmt.Println(envAPIKey, "is not specified.")
		os.Exit(1)
	}

	messageType := getEnvVar(envMessageType)
	if messageType == "" {
		fmt.Println(envMessageType, "is not specified.")
		os.Exit(1)
	}

	to := getEnvVar(envTo)
	if to == "" {
		fmt.Println(envTo, "is not specified.")
		os.Exit(1)
	}

	subject := getEnvVar(envSubject)
	if messageType == "stream" && subject == "" {
		fmt.Println(envSubject, "is not specified.")
		os.Exit(1)
	}

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
		"type":    {messageType},
		"to":      {to},
		"subject": {subject},
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
