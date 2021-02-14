package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

const appKey = "app=MTYxMzExNjk5N3xEdi1CQkFFQ180SUFBUkFCRUFBQV9fSF9nZ0FDQm5OMGNtbHVad3dSQUE5emFXZHVYM1Z3WDJSbGRHRnBiSE1HYzNSeWFXNW5EUC1sQVAtaWV5SmxiV0ZwYkY5aFpHUnlaWE56SWpvaVozVjVPVEl5UUdkdFlXbHNMbU52YlNJc0ltWnBjbk4wWDI1aGJXVWlPaUpIZFhraUxDSnNZWE4wWDI1aGJXVWlPaUpFZFc1emEza2lMQ0p2WTJOMWNHRjBhVzl1SWpvaUlpd2lkR1ZzSWpvaUlpd2lkVzVwZG1WeWMybDBlU0k2SWlJc0luZHZjbXR3YkdGalpTSTZJaUlzSW05bVptVnlhVzVuY3lJNmRISjFaU3dpZEdWeWJYTWlPblJ5ZFdWOUJuTjBjbWx1Wnd3S0FBaDFjMlZ5Ym1GdFpRWnpkSEpwYm1jTUNBQUdSM1Y1T1RJeXywI4NfyRdkcEk6YyhBlU3tt4OFJ3yocjuf8MRqZ_urGA=="
const baseURL = "https://welcome.cfapps.us10.hana.ondemand.com/node"
const maxNodes = 100

var opts = map[string]string{
	"offleft":  "?toggle=off&go=left",
	"offright": "?toggle=off&go=right",
	"onleft":   "?toggle=on&go=left",
	"onright":  "?toggle=on&go=right",
	"left":     "?go=left",
	"right":    "?go=right",
	"off":      "?toggle=off",
	"on":       "?toggle=on",
}

type Node struct {
	ConnectedDevices int  `json:"connected_iot_devices"`
	State            bool `json:"state"`
}

func DoRequest(method, URL string) []byte {
	req, _ := http.NewRequest(method, URL, nil)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Cookie", appKey)
	resp, _ := http.DefaultClient.Do(req)
	defer resp.Body.Close()
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	return bodyBytes
}

func GetCurNode() Node {
	bodyBytes := DoRequest("GET", baseURL)
	node := JsonToNode(bodyBytes)
	fmt.Println("Node:", node)
	return node
}

func MarkMove(action string) Node {
	bodyBytes := DoRequest("PATCH", baseURL+opts[action])
	node := JsonToNode(bodyBytes)
	fmt.Println("Node:", node)
	return node
}

func JsonToNode(NodeBytes []byte) Node {
	var node Node
	err := json.Unmarshal(NodeBytes, &node)
	if err != nil {
		fmt.Println(err)
	}
	return node
}


func main() {
	var node Node
	for i := 1; i < maxNodes; i++ {
		node = MarkMove("onright")
	}
	sum := 0
	for node.State {
		sum += node.ConnectedDevices
		node = MarkMove("offright")
	}
	println(sum)
}
