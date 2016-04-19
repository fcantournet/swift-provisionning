package main

import (
	"encoding/xml"
	"fmt"
	"github.com/codegangsta/cli"
	"log"
	"os/exec"
	"regexp"
	"strings"
)

var diskcodetotype = map[int]string{1: "HDD", 2: "SSD"}

const maxHDD = 23
const maxSSD = 2

// Device represents a swift Device on this node
type Device struct {
	DeviceName string
	Vd         VirtualDisk
	FsLabel    string
	MountPoint string
	FstabEntry string
	partition  bool
}

// OMAControllerIDs represents the list of controllers ID
// Deserialized from omreport output
type OMAControllerIDs struct {
	Cli           bool  `xml:"cli, attr"`
	ControllerIDs []int `xml:"Controllers>DCStorageObject>ControllerNum"`
}

// OMAVirtualDisks represents the list of Virtula Disks
// Deserialized from omreport output
type OMAVirtualDisks struct {
	Cli bool          `xml:"cli, attr"`
	VDs []VirtualDisk `xml:"VirtualDisks>DCStorageObject"`
}

// VirtualDisk represents a Virtula Disk
// Deserialized from omreport output
type VirtualDisk struct {
	ControllerNum int    `xml:"ControllerNum"`
	MediaType     int    `xml:"MediaType"`
	DeviceName    string `xml:"DeviceName"`
	Name          string `xml:"Name"`
	ID            int    `xml:"DeviceID"`
}

func (vd VirtualDisk) String() string {
	return fmt.Sprintf("%s  type=%s  device: %s", vd.Name, diskcodetotype[vd.MediaType], vd.DeviceName)
}

//type VDisksByController map[int][]VirtualDisk

func omreport(args string) ([]byte, error) {
	out, err := exec.Command("omreport", strings.Split(args, " ")...).Output()
	return out, err
}

func omconfig(args string) ([]byte, error) {
	out, err := exec.Command("omconfig", strings.Split(args, " ")...).Output()
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

func renameVdisk(newName string, controllerNum, vdiskID int, dry bool) (string, error) {
	fmt.Printf("omconfig storage vdisk action=rename controller=%v vdisk=%v name=%v\n", controllerNum, vdiskID, newName)
	if !dry {
		out, err := omconfig(fmt.Sprintf("storage vdisk action=rename controller=%v vdisk=%v name=%v", controllerNum, vdiskID, newName))
		log.Println(string(out))
		if err != nil {
			log.Fatal(err.Error())
		}
		return string(out), nil
	}
	return "", nil
}

func vdiskNameAvailable(name string, vdisks []VirtualDisk) bool {
	for _, vdisk := range vdisks {
		if vdisk.Name == name {
			return false
		}
	}
	return true
}

func getAvailableDiskNames(vdisks []VirtualDisk) (chan string, chan string) {
	availableHDD, availableSSD := make(chan string, maxHDD), make(chan string, maxSSD)

	for i := 0; i < maxHDD; i++ {
		name := fmt.Sprintf("HDD-%v", i)
		if vdiskNameAvailable(name, vdisks) {
			availableHDD <- name
		}
	}

	for i := 0; i < maxSSD; i++ {
		name := fmt.Sprintf("SSD-%v", i)
		if vdiskNameAvailable(name, vdisks) {
			availableSSD <- name
		}
	}
	close(availableHDD)
	close(availableSSD)

	return availableHDD, availableSSD
}

// RenameVdisks renames the vdisks already created following the (HDD|SSD)-x pattern
func RenameVdisks(ctx *cli.Context) {
	allvdisks, _ := getAllVdisks()
	availHDD, availSSD := getAvailableDiskNames(allvdisks.VDs)

	for _, vdisk := range allvdisks.VDs {
		match, err := regexp.Match("SSD|HDD|system", []byte(vdisk.Name))
		if err != nil {
			log.Fatal(err.Error())
		}
		if match {
			continue
		} else {
			switch vdisk.MediaType {
			case 1:
				if name, ok := <-availHDD; ok {
					renameVdisk(name, vdisk.ControllerNum, vdisk.ID, ctx.BoolT("dry"))
				}
			case 2:
				if name, ok := <-availSSD; ok {
					renameVdisk(name, vdisk.ControllerNum, vdisk.ID, ctx.BoolT("dry"))
				}
			default:
				log.Fatalf("Wrong MediaType")
			}
		}
	}
}

func getAllVdisks() (OMAVirtualDisks, error) {
	xmlvdisks, err := omreport("storage vdisk -fmt xml")
	if err != nil {
		log.Fatal(err.Error())
	}

	allvdisks, err := parseVdisks(xmlvdisks)
	if err != nil {
		log.Fatal(err.Error())
	}
	return allvdisks, nil
}

// Status print the current status of the Nodes Devices. Called from main.
func Status(ctx *cli.Context) {

	allvdisks, _ := getAllVdisks()
	allDevices := map[string]Device{}
	for _, vdisk := range allvdisks.VDs {
		allDevices[vdisk.DeviceName] = Device{vdisk.DeviceName, vdisk, "", "", "", false}
	}

	for _, vd := range allvdisks.VDs {
		fmt.Println(vd)
	}
}
