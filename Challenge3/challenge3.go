
package main

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"sort"
	"strconv"
)

const appKey = "app=MTYxMzExNjk5N3xEdi1CQkFFQ180SUFBUkFCRUFBQV9fSF9nZ0FDQm5OMGNtbHVad3dSQUE5emFXZHVYM1Z3WDJSbGRHRnBiSE1HYzNSeWFXNW5EUC1sQVAtaWV5SmxiV0ZwYkY5aFpHUnlaWE56SWpvaVozVjVPVEl5UUdkdFlXbHNMbU52YlNJc0ltWnBjbk4wWDI1aGJXVWlPaUpIZFhraUxDSnNZWE4wWDI1aGJXVWlPaUpFZFc1emEza2lMQ0p2WTJOMWNHRjBhVzl1SWpvaUlpd2lkR1ZzSWpvaUlpd2lkVzVwZG1WeWMybDBlU0k2SWlJc0luZHZjbXR3YkdGalpTSTZJaUlzSW05bVptVnlhVzVuY3lJNmRISjFaU3dpZEdWeWJYTWlPblJ5ZFdWOUJuTjBjbWx1Wnd3S0FBaDFjMlZ5Ym1GdFpRWnpkSEpwYm1jTUNBQUdSM1Y1T1RJeXywI4NfyRdkcEk6YyhBlU3tt4OFJ3yocjuf8MRqZ_urGA=="
const baseURL = "https://welcome.cfapps.us10.hana.ondemand.com/device"

type Device struct {
	DeviceType string `json:"device_type"`
	Id         string `json:"id"`
	Seq        int    `json:"seq"`
	State      string `json:"state"`
}

type DeviceList struct {
	NextPage string   `json:"next_page"`
	Pages    int      `json:"pages"`
	NumItems int      `json:"num_items"`
	Items    []Device `json:"items"`
}

type DeviceStatus struct {
	ActiveTime  int    `json:"active_time"`
	Cpu         int    `json:"cpu"`
	DeviceType  string `json:"device_type"`
	Id          string `json:"id"`
	Seq         int    `json:"seq"`
	State       string `json:"state"`
	Temperature int    `json:"temperature"`
}


// Make an HTTP request to the given URL with the given method. Return the request body.
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

// Call DoRequest Asynchronously
func DoRequestAsync(method, URL string, c chan []byte){
	c <- DoRequest(method, URL)
}

// Get the total number of pages
func GetPagesNumber() int {
	respBytes := DoRequest("GET", baseURL)
	var deviceList DeviceList
	json.Unmarshal(respBytes, &deviceList)
	return deviceList.Pages
}

// Get an array of devices with status Online
func GetOnlineDevices() []Device {
	c := make(chan []byte)
	var result []Device
	pagesNum := GetPagesNumber()
	for i := 0; i <= pagesNum; i++ {
		println("Async req#" + strconv.Itoa(i) + " sent") // DEBUG
		go DoRequestAsync("GET", baseURL+"?next="+strconv.Itoa(i), c)
	}
	for i := 0; i <= pagesNum; i++ {
		var deviceList DeviceList
		json.Unmarshal(<-c, &deviceList)
		for _, device := range deviceList.Items {
			if device.State == "Online" {
				result = append(result, device)
			}
		}
	}
	return result
}

// Get an array of devices with temp <= 29
func GetGoodDevices(devices []Device) []DeviceStatus {
	c := make(chan []byte)

	var goodDevices []DeviceStatus
	for i, device := range devices {
		println("#" + strconv.Itoa(i) +" Async temp req for device # " + device.Id)
		go DoRequestAsync("GET", baseURL+"/"+device.Id+"/status", c)
	}

	for i := 0; i < len(devices); i++ {
		var goodDevice DeviceStatus
		json.Unmarshal(<-c, &goodDevice)
		if goodDevice.Temperature <= 29 {
			goodDevices = append(goodDevices, goodDevice)
		}
	}

	sort.Slice(goodDevices[:], func(i, j int) bool {
		return goodDevices[i].Seq < goodDevices[j].Seq
	})
	return goodDevices
}

func GetMD5Hash(text string) string {
	hash := md5.Sum([]byte(text))
	return hex.EncodeToString(hash[:])
}


func main() {
	devices := GetOnlineDevices()
	goodDevices := GetGoodDevices(devices)
	goodDevicesStr := ""
	for _, device := range goodDevices {
		goodDevicesStr += device.Id
	}
	println(GetMD5Hash(goodDevicesStr))
}
