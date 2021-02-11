package models

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/LambdaTest/real-device-db-insertion/gigafox"
	"github.com/Sirupsen/logrus"
	_ "github.com/go-sql-driver/mysql"
)

type Template struct {
	TemplateId string
	Name       string
}

type RealDeviceTemplate struct {
	TemplateId      string
	DeviceId        string
	DeviceName      string
	GigafoxId       string
	GigafoxName     string
	HubUrl          string
	ProxyHost       string
	ProxyPort       uint16
	ProxyApiPort    uint16
	ProxyTslPort    uint16
	Version         string
	GigafoxUsername string
	GigafoxApiKey   string
}

type Vm struct {
	VmId       string
	Name       string
	Ip         string
	TemplateId string
}

func InsertInLVMSNew(lvmsDB *sql.DB, device gigafox.Device, Logger *logrus.Logger) {
	var template Template
	var realDeviceTemplate RealDeviceTemplate
	//var vm Vm

	parts := strings.Split(device.OperatingSystemVersion, ".")
	//Preparing data for template table

	if parts[1] != "0" {
		template.TemplateId = "real-" + strings.ReplaceAll(strings.ToLower(device.FriendlyModel), " ", "-") + "-" + parts[0] + "." + parts[1]
	} else {
		template.TemplateId = "real-" + strings.ReplaceAll(strings.ToLower(device.FriendlyModel), " ", "-") + "-" + parts[0]
	}

	Logger.WithFields(logrus.Fields{
		"template_id": template.TemplateId,
	}).Info("Template Table Data")

	//Inserting in template table
	sqlTemplate := `INSERT INTO lambda_lvms_new.template
					(template_id, name, is_create_enabled, is_delete_enabled, pool_size, idle_vm_size, increment_size, release_version, release_threshold, time_to_live, last_sync, status_ind, update_approach, delete_older_vms, auto_update_release)
					SELECT ?, 'Real IOS', 0, 0, 5, 3, 0, 'v1.0', 50, 24, NULL, 'active', 'delete', 1, 'disabled'
						WHERE NOT EXISTS 
							(SELECT * FROM lambda_lvms_new.template
								WHERE template_id= ?)`
	resTemplate, err := lvmsDB.Query(sqlTemplate, template.TemplateId, template.TemplateId)
	defer resTemplate.Close()

	if err != nil {
		Logger.WithFields(logrus.Fields{
			"err": err.Error(),
		}).Error("Error in Insertion in template")
		panic(err.Error())
	}

	//Preparing data for real_device_template table
	realDeviceTemplate.TemplateId = template.TemplateId
	realDeviceTemplate.DeviceName = device.FriendlyModel
	realDeviceTemplate.GigafoxId = device.Id
	realDeviceTemplate.GigafoxName = device.Name
	realDeviceTemplate.HubUrl = "https://lambdatest5.gigafox.io"

	//Checking if record is already present in real_device_template table
	rows, err := lvmsDB.Query("SELECT device_id FROM real_device_template WHERE template_id = ? ORDER BY device_template_id asc", realDeviceTemplate.TemplateId)
	if err != nil {
		Logger.WithFields(logrus.Fields{
			"err": err.Error(),
		}).Error("Error in Selecting in real_device_template")
		panic(err.Error())
	}
	defer rows.Close()

	var device_id string
	for rows.Next() {
		err := rows.Scan(&device_id)
		if err != nil {
			log.Fatal(err)
		}
	}
	err = rows.Err()
	if err != nil {
		Logger.WithFields(logrus.Fields{
			"err": err.Error(),
		}).Error("Error in Traversing real_device_template data")
		panic(err.Error())
	}

	if len(device_id) <= 0 {
		realDeviceTemplate.DeviceId = realDeviceTemplate.TemplateId + "-1"
	} else {
		lastdeviceExists := device_id[len(device_id)-1:]
		i, err := strconv.Atoi(lastdeviceExists)
		if err != nil {
			// handle error
			fmt.Println(err)
			os.Exit(2)
		}
		i = i + 1
		realDeviceTemplate.DeviceId = realDeviceTemplate.TemplateId + "-" + strconv.Itoa(i)
	}

	//Checking MAX PORT real_device_template table
	rowsPort, err := lvmsDB.Query(`SELECT proxy_port, proxy_api_port , proxy_tls_port FROM lambda_lvms_new.real_device_template where template_id like "%iphone%" or template_id like "%ipad%" ORDER by proxy_port DESC LIMIT 1`)
	if err != nil {
		Logger.WithFields(logrus.Fields{
			"err": err.Error(),
		}).Error("Error in Selecting in real_device_template PORT")
		panic(err.Error())
	}
	defer rowsPort.Close()

	for rowsPort.Next() {
		err := rowsPort.Scan(&realDeviceTemplate.ProxyPort, &realDeviceTemplate.ProxyApiPort, &realDeviceTemplate.ProxyTslPort)
		if err != nil {
			log.Fatal(err)
		}
	}
	err = rowsPort.Err()
	if err != nil {
		Logger.WithFields(logrus.Fields{
			"err": err.Error(),
		}).Error("Error in Traversing real_device_template PORT data")
		panic(err.Error())
	}

	// err = lvmsDB.QueryRow(`SELECT proxy_port, proxy_api_port , proxy_tls_port FROM lambda_lvms_new.real_device_template where template_id like "%iphone%" or template_id like "%ipad%" ORDER by proxy_port DESC LIMIT 1`).Scan(&realDeviceTemplate.ProxyPort, &realDeviceTemplate.ProxyApiPort, &realDeviceTemplate.ProxyTslPort)
	// if err != nil {
	// 	Logger.WithFields(logrus.Fields{
	// 		"err": err.Error(),
	// 	}).Error("Error in Selecting real_device_template port")
	// 	panic(err.Error())
	// }
	if realDeviceTemplate.ProxyPort == 0 && realDeviceTemplate.ProxyApiPort == 0 && realDeviceTemplate.ProxyTslPort == 0 {
		Logger.WithFields(logrus.Fields{
			"err": "Dont know which PORT to allocate",
		}).Warn("PORT allocation issue")
		os.Exit(1)
		//panic("PORT allocation issue")
	} else {
		realDeviceTemplate.ProxyPort += 2
		realDeviceTemplate.ProxyApiPort += 2
		realDeviceTemplate.ProxyTslPort += 2
	}

	Logger.WithFields(logrus.Fields{
		"template_id":    realDeviceTemplate.TemplateId,
		"device_id":      realDeviceTemplate.DeviceId,
		"gigafox_id":     realDeviceTemplate.GigafoxId,
		"gigafox_name":   realDeviceTemplate.GigafoxName,
		"hub_url":        realDeviceTemplate.HubUrl,
		"proxy_port":     realDeviceTemplate.ProxyPort,
		"proxy_api_port": realDeviceTemplate.ProxyApiPort,
		"proxy_tls_port": realDeviceTemplate.ProxyTslPort,
	}).Info("Real_device_template Table Data")

	//Inserting in real_device_template table
	sqlRDT := `INSERT INTO real_device_template
	(template_id, device_id, gigafox_id, gigafox_name, hub_url, status_ind, proxy_host, proxy_port, proxy_api_port, proxy_tls_port, proxy_host_port, video_host, video_host_port, gigafox_username, gigafox_api_key)
		SELECT ?, ?, ?, ?, ?, 'active', 'http://stage-proxy.lambdatest.com', ?, ?, ?, 8080, 'http://stage-video.lambdatest.com', 8080, 'dev@lambdatest.com', '2c083852-d1c8-4f8e-b5b5-8dd876dd6c4c'
				WHERE NOT EXISTS 
				(SELECT * FROM real_device_template
				  WHERE device_id= ?)`
	resRDT, err := lvmsDB.Query(sqlRDT, realDeviceTemplate.TemplateId, realDeviceTemplate.DeviceId, realDeviceTemplate.GigafoxId, realDeviceTemplate.GigafoxName, realDeviceTemplate.HubUrl, realDeviceTemplate.ProxyPort, realDeviceTemplate.ProxyApiPort, &realDeviceTemplate.ProxyTslPort, realDeviceTemplate.DeviceId)
	defer resRDT.Close()

	if err != nil {
		Logger.WithFields(logrus.Fields{
			"err": err.Error(),
		}).Error("Error in Insertion in real_device_template")
		panic(err.Error())
	}

	//Preparing data for vm table

	//Inserting in vm table
	// sqlVm := ``
	// resVm, err := lvmsDB.Query(sqlVm)
	// defer resVm.Close()

	// if err != nil {
	// 	Logger.WithFields(logrus.Fields{
	// 		"err": err.Error(),
	// 	}).Error("Error in Insertion in Vm")
	// 	panic(err.Error())
	// }

}
