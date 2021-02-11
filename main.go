package main

import (
	"strings"

	"github.com/LambdaTest/real-device-db-insertion/dbconn"
	"github.com/LambdaTest/real-device-db-insertion/gigafox"
	"github.com/LambdaTest/real-device-db-insertion/logging"
	"github.com/LambdaTest/real-device-db-insertion/models"
)

// var (
// 	// Logger is a configured logrus.Logger.
// 	log *logrus.Logger
// )

func init() {
	dbconn.InitLTMSConn()
	dbconn.InitLVMSConn()
	logging.CreateLog()
}

func main() {

	devices := gigafox.GetDeviceData()
	// s, _ := json.MarshalIndent(devices, "", "\t")
	// fmt.Println(string(s))
	//logging.Log.Info(devices[0])

	CheckForIosDevice(devices)
}

func CheckForIosDevice(devices []gigafox.Device) {
	for _, device := range devices {
		if strings.ToUpper(device.OperatingSystem) == "IOS" {
			ltmsDB := dbconn.GetLTMSConn()
			lvmsDB := dbconn.GetLVMSConn()
			logging.Log.Info(device)
			models.InsertInLTMSNew(ltmsDB, device, logging.Log)
			models.InsertInLVMSNew(lvmsDB, device, logging.Log)
			break
		}
	}
}
