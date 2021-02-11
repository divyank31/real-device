package models

import (
	"database/sql"
	"strings"

	"github.com/LambdaTest/real-device-db-insertion/gigafox"
	"github.com/Sirupsen/logrus"
	_ "github.com/go-sql-driver/mysql"
)

type RealDevice struct {
	DeviceId    string
	OSVersionId string
	DeviceName  string
}

type RealDeviceAlias struct {
	DeviceId string
	Alias    string
}

func InsertInLTMSNew(ltmsDB *sql.DB, device gigafox.Device, Logger *logrus.Logger) {
	var realDevice RealDevice
	var realDeviceAlias RealDeviceAlias

	//preparing data for real_device table
	parts := strings.Split(device.OperatingSystemVersion, ".")

	if parts[1] != "0" {
		realDevice.DeviceId = strings.ReplaceAll(strings.ToLower(device.FriendlyModel), " ", "-") + "-" + parts[0] + "." + parts[1]
		realDevice.OSVersionId = strings.ToLower(device.OperatingSystem) + "-" + parts[0] + "." + parts[1]
	} else {
		realDevice.DeviceId = strings.ReplaceAll(strings.ToLower(device.FriendlyModel), " ", "-") + "-" + parts[0]
		realDevice.OSVersionId = strings.ToLower(device.OperatingSystem) + "-" + parts[0]
	}
	realDevice.DeviceName = device.FriendlyModel

	Logger.WithFields(logrus.Fields{
		"device_id":     realDevice.DeviceId,
		"os_version_id": realDevice.OSVersionId,
		"device_name":   realDevice.DeviceName,
	}).Info("Real_Device Table Data")

	//Inserting in real_device table
	sqlRD := `INSERT INTO real_device (device_id, os_version_id, device_name, browser_id, resolution, viewport, screen_size, status_ind, display_order)
				SELECT ?, ?, ?, 'Safari', '', '', 6, 'active', 2
					WHERE NOT EXISTS 
						(SELECT * FROM real_device
							WHERE device_id=? AND os_version_id = ? AND device_name= ?)`
	resRD, err := ltmsDB.Query(sqlRD, realDevice.DeviceId, realDevice.OSVersionId, realDevice.DeviceName, realDevice.DeviceId, realDevice.OSVersionId, realDevice.DeviceName)
	defer resRD.Close()

	if err != nil {
		Logger.WithFields(logrus.Fields{
			"err": err.Error(),
		}).Error("Error in Insertion in real_device")
		panic(err.Error())
	}

	//preparing data for real_device_alias table
	realDeviceAlias.DeviceId = realDevice.DeviceId
	realDeviceAlias.Alias = realDevice.DeviceName

	Logger.WithFields(logrus.Fields{
		"device_id": realDeviceAlias.DeviceId,
		"alias":     realDeviceAlias.Alias,
	}).Info("Real_device_alias Table Data")

	//Inserting in real_device_alias table
	sqlRDA := `INSERT INTO real_device_alias (device_id, alias)
				SELECT  ?, ?
					WHERE NOT EXISTS 
						(SELECT * FROM real_device_alias 
							WHERE device_id= ? AND alias= ?)`

	resRDA, err := ltmsDB.Query(sqlRDA, realDeviceAlias.DeviceId, realDeviceAlias.Alias, realDeviceAlias.DeviceId, realDeviceAlias.Alias)
	defer resRDA.Close()

	if err != nil {
		Logger.WithFields(logrus.Fields{
			"err": err.Error(),
		}).Error("Error in Insertion in real_device_alias")
		panic(err.Error())
	}

}
