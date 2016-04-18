package main

import (
	"encoding/xml"
	"fmt"
	"github.com/codegangsta/cli"
	"log"
	"os/exec"
	"strings"
	//"io/ioutil"
	//"regexp"
)

var diskcodetotype = map[int]string{1: "HDD", 2: "SSD"}

type Device struct {
	DeviceName string
	Vd         *VirtualDisk
	FsLabel    string
	MountPoint string
	FstabEntry string
	partition  bool
}

type OMAControllerIDs struct {
	Cli           bool  `xml:"cli, attr"`
	ControllerIDs []int `xml:"Controllers>DCStorageObject>ControllerNum"`
}

type OMAVirtualDisks struct {
	Cli bool          `xml:"cli, attr"`
	VDs []VirtualDisk `xml:"VirtualDisks>DCStorageObject"`
}

type VirtualDisk struct {
	ControllerNum int    `xml:"ControllerNum"`
	MediaType     int    `xml:"MediaType"`
	DeviceName    string `xml:"DeviceName"`
	Name          string `xml:"Name"`
}

func (vd VirtualDisk) String() string {
	return fmt.Sprintf("%s  type=%s  device: %s", vd.Name, diskcodetotype[vd.MediaType], vd.DeviceName)
}

//type VDisksByController map[int][]VirtualDisk

func omreport(args string) ([]byte, error) {
	out, err := exec.Command("omreport", strings.Split(args, " ")...).Output()
	return out, err
}

func parseControllerIDs(xmldata []byte) (OMAControllerIDs, error) {
	var oma OMAControllerIDs
	err := xml.Unmarshal(xmldata, &oma)
	return oma, err
}

func parseVdisks(xmldata []byte) (OMAVirtualDisks, error) {
	var oma OMAVirtualDisks
	err := xml.Unmarshal(xmldata, &oma)
	return oma, err
}

func getmounts() {
	out, _ := exec.Command("mount -l -t xfs -l").Output()
	fmt.Print(string(out))
}

func Status(ctx *cli.Context) {
	// xmlcont, err := omreport("omreport storage controller -fmt xml")
	// if err != nil {
	//   log.Fatal(err.Error())
	// }
	// controllers, err := parseControllerIDs(xmlcont)
	// if err != nil {
	//   log.Fatal(err.Error())
	// }

	xmlvdisks, err := omreport("storage vdisk -fmt xml")
	if err != nil {
		log.Fatal(err.Error())
	}

	allvdisks, err := parseVdisks(xmlvdisks)
	if err != nil {
		log.Fatal(err.Error())
	}
	allDevices := map[string]Device{}
	for _, vdisk := range allvdisks.VDs {
		allDevices[vdisk.DeviceName] = Device{vdisk.DeviceName, &vdisk, "", "", "", false}
	}

	for _, vd := range allvdisks.VDs {
		fmt.Println(vd)
	}
}
