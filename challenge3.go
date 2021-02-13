package main

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
)

const APP_KEY = "app=MTYxMzExNjk5N3xEdi1CQkFFQ180SUFBUkFCRUFBQV9fSF9nZ0FDQm5OMGNtbHVad3dSQUE5emFXZHVYM1Z3WDJSbGRHRnBiSE1HYzNSeWFXNW5EUC1sQVAtaWV5SmxiV0ZwYkY5aFpHUnlaWE56SWpvaVozVjVPVEl5UUdkdFlXbHNMbU52YlNJc0ltWnBjbk4wWDI1aGJXVWlPaUpIZFhraUxDSnNZWE4wWDI1aGJXVWlPaUpFZFc1emEza2lMQ0p2WTJOMWNHRjBhVzl1SWpvaUlpd2lkR1ZzSWpvaUlpd2lkVzVwZG1WeWMybDBlU0k2SWlJc0luZHZjbXR3YkdGalpTSTZJaUlzSW05bVptVnlhVzVuY3lJNmRISjFaU3dpZEdWeWJYTWlPblJ5ZFdWOUJuTjBjbWx1Wnd3S0FBaDFjMlZ5Ym1GdFpRWnpkSEpwYm1jTUNBQUdSM1Y1T1RJeXywI4NfyRdkcEk6YyhBlU3tt4OFJ3yocjuf8MRqZ_urGA=="
const BASE_URL = "https://welcome.cfapps.us10.hana.ondemand.com/device"

type Device struct{
	DeviceType string 	`json:"device_type"`
	Md5Id 	   string 	`json:"id"`
	Seq 	   int 		`json:"seq"`
	State      string 	`json:"state"`
}

type DeviceList struct {
	NextPage   string 	`json:"next_page"`
	Pages 	   int 		`json:"pages"`
	NumItems   int 		`json:"num_items"`
	Items 	   []Device `json:"items"`
}

type DeviceStatus struct{
	ActiveTime  int 	`json:"active_time"`
	Cpu         int 	`json:"cpu"`
	DeviceType  string  `json:"device_type"`
	Id 			string  `json:"id"`
	Seq 		int 	`json:"seq"`
	State 		string  `json:"state"`
	Temperature int 	`json:"temperature"`
}

func DoRequest(method, URL string) []byte{
	req, _ := http.NewRequest(method, URL, nil)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Cookie", APP_KEY)
	resp, _ := http.DefaultClient.Do(req)
	defer resp.Body.Close()
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	return bodyBytes
}

func GetPageNumber() int {
	respBytes := DoRequest("GET", BASE_URL)
	var deviceList DeviceList
	json.Unmarshal(respBytes, &deviceList)
	return deviceList.Pages
}

func GetOnlineDevices() []Device{
	var result []Device
	pagesNum := GetPageNumber()
	for i := 0; i <= pagesNum; i++{
		print(i)
		print(" ")
		respBytes := DoRequest("GET", BASE_URL + "?next=" + strconv.Itoa(i))
		var deviceList DeviceList
		json.Unmarshal(respBytes, &deviceList)
		for _, device := range deviceList.Items {
			if device.State == "Online"{
				result = append(result, device)
			}
		}
	}
	return result
}

func IsDeviceHot(deviceId string) bool{
	respBytes := DoRequest("GET", BASE_URL + "/" + deviceId + "/status")
	var deviceStatus DeviceStatus
	json.Unmarshal(respBytes, &deviceStatus)
	println(" Temp: " + strconv.Itoa(deviceStatus.Temperature))
	return deviceStatus.Temperature > 29
}


func GetMD5Hash(text string) string {
	hash := md5.Sum([]byte(text))
	return hex.EncodeToString(hash[:])
}

func main() {
	devices := GetOnlineDevices()
	goodDevicesStr := ""
	for i, device := range devices{
		print(i)
		print("/")
		print(len(devices))
		if !IsDeviceHot(device.Md5Id){
			println(device.Md5Id)
			goodDevicesStr += device.Md5Id
		}
	}
	println(GetMD5Hash(goodDevicesStr))
}