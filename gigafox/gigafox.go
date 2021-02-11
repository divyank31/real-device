package gigafox

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

const (
	USERNAME = "dev@lambdatest.com"
	PASSWORD = "2c083852-d1c8-4f8e-b5b5-8dd876dd6c4c"
	URL      = "https://lambdatest5.gigafox.io/apiv1/Device"
)

type Device struct {
	Id                     string
	Availability           string
	Enabled                bool
	Name                   string
	OperatingSystem        string
	OperatingSystemVersion string
	FriendlyModel          string
	VendorDeviceName       string
}

func GetDeviceData() []Device {

	var devices []Device

	req, err := http.NewRequest("GET", URL, nil)
	req.SetBasicAuth(USERNAME, PASSWORD)

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")

	res, _ := http.DefaultClient.Do(req)

	body, err := ioutil.ReadAll(res.Body)

	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()

	json.Unmarshal([]byte(body), &devices)
	return devices
	// s, _ := json.MarshalIndent(devices, "", "\t")
	// fmt.Print(string(s))
}
