package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

var (
	imgPath = flag.String("img-path", "", "Path to the image to be processed")
	jwt     = flag.String("jwt", "", "JWT for authentication")
	lang    = flag.String("lang", "English", "Language for the haiku")
	tags    = flag.String("tags", "", "Comma-separated list of tags for the haiku")
	port    = flag.String("port", "8080", "Port to run the server on")
)

func main() {
	parseInput()

	encoded, err := base64Img(*imgPath)
	if err != nil {
		log.Fatalln("Error encoding image:", err)
	}

	tagList := tagList()

	bodyBytes, err := reqBody(encoded, tagList)
	if err != nil {
		log.Fatalln("Error creating request body:", err)
	}

	req, err := req(bodyBytes)
	if err != nil {
		log.Fatalln("Error creating request:", err)
	}

	fmt.Println("Sending request to server...")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalln("Error making request:", err)
	}

	handleResponse(resp)
}

func parseInput() {
	flag.Parse()

	if *imgPath == "" {
		flag.Usage()
		os.Exit(1)
	}

	if *jwt == "" {
		flag.Usage()
		os.Exit(1)
	}
}

func base64Img(imgPath string) (string, error) {
	// Read the image file
	data, err := os.ReadFile(imgPath)
	if err != nil {
		return "", err
	}
	// Encode the image to base64
	encoded := base64.StdEncoding.EncodeToString(data)
	return encoded, nil
}

func tagList() []string {
	var list []string
	if *tags != "" {
		list = strings.Split(*tags, ",")
	}

	return list
}

func reqBody(base64Image string, tagList []string) ([]byte, error) {
	body := struct {
		Base64Image string   `json:"base64Image"`
		Language    string   `json:"language"`
		Tags        []string `json:"tags"`
	}{
		Base64Image: base64Image,
		Language:    *lang,
		Tags:        tagList,
	}

	return json.Marshal(body)
}

func req(bodyBytes []byte) (*http.Request, error) {
	req, err := http.NewRequest("POST", "http://127.0.0.1:"+*port, bytes.NewBuffer(bodyBytes))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+*jwt)

	return req, nil
}

func handleResponse(resp *http.Response) {
	defer resp.Body.Close()

	println("Response status:", resp.Status)

	// Pretty-print the response body.
	respBodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln("Error reading response body:", err)
	}

	var prettyJSON bytes.Buffer
	if err := json.Indent(&prettyJSON, respBodyBytes, "", "  "); err != nil {
		log.Fatalln("Error formatting JSON:", err, string(respBodyBytes))
	} else {
		fmt.Println("Response body:", prettyJSON.String())
	}
}
