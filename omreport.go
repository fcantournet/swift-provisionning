package main

import (
  "github.com/codegangsta/cli"
	"encoding/xml"
  "log"
  "os/exec"
  "fmt"
  "io/ioutil"
)

var diskcodetotype map[int]string = map[int]string{1: "HDD", 2: "SSD"}

type OMAControllerIDs struct {
	Cli           bool  `xml:"cli, attr"`
	ControllerIDs []int `xml:"Controllers>DCStorageObject>ControllerNum"`
}

type OMAVirtualDisks struct {
	Cli bool          `xml:"cli, attr"`
	VDs []VirtualDisk `xml:"VirtualDisks>DCStorageObject"`
}

type VirtualDisk struct {
	ControllerNum int `xml:"ControllerNum"`
	MediaType     int    `xml:"MediaType"`
	DeviceName    string `xml:"DeviceName"`
	Name          string `xml:"Name"`
}

func (vd VirtualDisk)String() string {
  return fmt.Sprintf("%s  type=%s  device: %s", vd.Name, diskcodetotype[vd.MediaType], vd.DeviceName)
}

//type VDisksByController map[int][]VirtualDisk


func omreport(command string) ([]byte, error) {
  out, err := exec.Command(command).Output()
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

  // xmlvdisks, err := omreport("omreport storage vdisk -fmt xml")
  // if err != nil {
  //   log.Fatal(err.Error())
  // }
  xmlvdisks, err := ioutil.ReadFile("./testdata.xml")

  allvdisks, err := parseVdisks(xmlvdisks)
  if err != nil {
    log.Fatal(err.Error())
  }

  for _, vd := range(allvdisks.VDs) {
    fmt.Println(vd)
  }
}
