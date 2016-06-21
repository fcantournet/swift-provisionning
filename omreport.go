package main

import (
	"encoding/xml"
	"fmt"
	"github.com/codegangsta/cli"
	"github.com/pkg/errors"
	"log"
	"os/exec"
	"regexp"
	"strings"
)

const (
	cHDD = 1
	cSSD = 2
)

var diskcodetotype = map[int]string{cHDD: "HDD", cSSD: "SSD"}

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

// OMAVirtualDisks represents the list of Virtual Disks
// Deserialized from omreport output
type OMAVirtualDisks struct {
	Cli bool          `xml:"cli, attr"`
	VDs []VirtualDisk `xml:"VirtualDisks>DCStorageObject"`
}

// VirtualDisk represents a Virtual Disk
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
	if err != nil {
		return out, errors.Wrap(err, "omreport failed")
	}
	return out, nil
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

func renameVdisk(newName string, vdisk VirtualDisk, dry bool) (string, error) {
	if newName == vdisk.Name {
		return newName, nil
	}
	fmt.Printf("omconfig storage vdisk action=rename controller=%v vdisk=%v name=%v\n", vdisk.ControllerNum, vdisk.ID, newName)
	if !dry {
		out, err := omconfig(fmt.Sprintf("storage vdisk action=rename controller=%v vdisk=%v name=%v", vdisk.ControllerNum, vdisk.ID, newName))
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

func getAvailableDiskNames(vdisks []VirtualDisk, maxHDD, maxSSD int, yolo bool) (chan string, chan string) {
	availableHDD, availableSSD := make(chan string, maxHDD), make(chan string, maxSSD)

	fmt.Printf("Yolo = %v ", yolo)
	fmt.Printf("maxSSD = %v ", maxSSD)
	fmt.Printf("maxHDD = %v ", maxHDD)

	for i := 0; i < maxHDD; i++ {
		name := fmt.Sprintf("HDD-%v", i)
		if yolo || vdiskNameAvailable(name, vdisks) {
			availableHDD <- name
		}
	}

	for i := 0; i < maxSSD; i++ {
		name := fmt.Sprintf("SSD-%v", i)
		if yolo || vdiskNameAvailable(name, vdisks) {
			availableSSD <- name
		}
	}
	close(availableHDD)
	close(availableSSD)

	return availableHDD, availableSSD
}

func checkMaxDisks(hdd, ssd int, allvd []VirtualDisk) (int, int) {
	allHDD, allSSD := 0, 0
	for _, vd := range allvd {
		switch vd.MediaType {
		case cHDD:
			allHDD++
		case cSSD:
			allSSD++
		default:
			log.Fatalf("Wrong MediaType %v on %v", vd.MediaType, vd)
		}
	}
	if hdd < 0 {
		hdd = allHDD
	}
	if ssd < 0 {
		ssd = allSSD
	}
	return hdd, ssd
}

// RenameVdisks renames the vdisks already created following the (HDD|SSD)-x pattern
func RenameVdisks(ctx *cli.Context) {

	allvdisks, _ := getAllVdisks()

	maxHDD, maxSSD := checkMaxDisks(ctx.Int("maxhdd"), ctx.Int("maxssd"), allvdisks.VDs)
	availHDD, availSSD := getAvailableDiskNames(allvdisks.VDs, maxHDD, maxSSD, ctx.Bool("yolo"))

	for _, vdisk := range allvdisks.VDs {
		match, err := regexp.Match("HDD|SSD|system", []byte(vdisk.Name))
		if err != nil {
			log.Fatal(err.Error())
		}
		log.Printf("matched : %v", match)
		if !match || ctx.Bool("yolo") {
			switch vdisk.MediaType {
			case cHDD:
				if name, ok := <-availHDD; ok {
					renameVdisk(name, vdisk, ctx.BoolT("dry"))
				}
			case cSSD:
				if name, ok := <-availSSD; ok {
					renameVdisk(name, vdisk, ctx.BoolT("dry"))
				}
			default:
				log.Fatalf("Wrong MediaType %v on %v", vdisk.MediaType, vdisk.Name)
			}
		}
	}
}

// Gets all vdisks from omreport
// This will filter the system vdisk
func getAllVdisks() (OMAVirtualDisks, error) {
	xmlvdisks, err := omreport("storage vdisk -fmt xml")
	if err != nil {
		log.Fatal(err.Error())
	}

	allvdisks, err := parseVdisks(xmlvdisks)
	if err != nil {
		log.Fatal(err.Error())
	}

	filtered := make([]VirtualDisk, 0, len(allvdisks.VDs))
	for _, vd := range allvdisks.VDs {
		if strings.Contains(vd.Name, "system") || vd.DeviceName == "/dev/sda" {
			continue
		}
		filtered = append(filtered, vd)
	}
	allvdisks.VDs = filtered
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
