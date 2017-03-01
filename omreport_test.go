package main

import (
	"io/ioutil"
	"testing"
)

func TestParseControllerXML(t *testing.T) {
	expected := []int{0, 1}
	xmldata, err := ioutil.ReadFile("./testdata_controllers.xml")
	if err != nil {
		t.Fatal(err.Error())
	}
	temp, err := parseControllerIDs(xmldata)
	actual := temp.ControllerIDs
	if len(actual) != len(expected) {
		t.Fatalf("Test failed, expected: '%#v', got:  '%#v'", expected, actual)
	}
	for i := range expected {
		if actual[i] != expected[i] {
			t.Fatalf("Test failed, expected: '%#v', got:  '%#v'", expected, actual)
		}
	}
}

func TestParseVdiskXML(t *testing.T) {
	xmldata, err := ioutil.ReadFile("./testdata.xml")
	if err != nil {
		t.Fatal(err.Error())
	}
	actual, err := parseVdisks(xmldata)
	if err != nil {
		t.FailNow()
	}

	expected := []VirtualDisk{VirtualDisk{ControllerNum: 0, MediaType: 1, DeviceName: "/dev/sda", Name: "system00", ID: 0},
		VirtualDisk{ControllerNum: 0, MediaType: 1, DeviceName: "/dev/sdb", Name: "data00", ID: 1},
		VirtualDisk{ControllerNum: 0, MediaType: 1, DeviceName: "/dev/sdc", Name: "data01", ID: 2},
		VirtualDisk{ControllerNum: 0, MediaType: 1, DeviceName: "/dev/sdd", Name: "data02", ID: 3},
		VirtualDisk{ControllerNum: 0, MediaType: 1, DeviceName: "/dev/sde", Name: "data03", ID: 4},
		VirtualDisk{ControllerNum: 0, MediaType: 1, DeviceName: "/dev/sdf", Name: "data04", ID: 5},
		VirtualDisk{ControllerNum: 0, MediaType: 1, DeviceName: "/dev/sdg", Name: "data05", ID: 6},
		VirtualDisk{ControllerNum: 0, MediaType: 1, DeviceName: "/dev/sdh", Name: "data06", ID: 7},
		VirtualDisk{ControllerNum: 0, MediaType: 1, DeviceName: "/dev/sdi", Name: "data07", ID: 8},
		VirtualDisk{ControllerNum: 0, MediaType: 1, DeviceName: "/dev/sdj", Name: "data08", ID: 9},
		VirtualDisk{ControllerNum: 0, MediaType: 1, DeviceName: "/dev/sdk", Name: "data09", ID: 10},
		VirtualDisk{ControllerNum: 0, MediaType: 1, DeviceName: "/dev/sdl", Name: "data10", ID: 11},
		VirtualDisk{ControllerNum: 0, MediaType: 1, DeviceName: "/dev/sdm", Name: "data11", ID: 12},
		VirtualDisk{ControllerNum: 0, MediaType: 1, DeviceName: "/dev/sdn", Name: "data12", ID: 13},
		VirtualDisk{ControllerNum: 0, MediaType: 1, DeviceName: "/dev/sdo", Name: "data13", ID: 14},
		VirtualDisk{ControllerNum: 0, MediaType: 2, DeviceName: "/dev/sdp", Name: "data26", ID: 15},
		VirtualDisk{ControllerNum: 0, MediaType: 2, DeviceName: "/dev/sdac", Name: "data27", ID: 16},
	}
	failed := false
	for i, vd := range actual.VDs {
		if vd != expected[i] {
			t.Logf("actual : %v != expected : %v", vd, expected[i])
			failed = true
		}
	}
	if failed {
		t.Fail()
	}
}

func TestParsePdiskXML(t *testing.T) {
	xmldata, err := ioutil.ReadFile("./testdata_pdisks.xml")
	if err != nil {
		t.Fatal(err.Error())
	}
	actual, err := parsePdisks(xmldata)
	if err != nil {
		t.FailNow()
	}

	expected := []PhysicalDisk{
		PhysicalDisk{ControllerNum: 0, MediaType: 1, TargetID: 2, ID: 0},
		PhysicalDisk{ControllerNum: 0, MediaType: 1, TargetID: 2, ID: 1},
		PhysicalDisk{ControllerNum: 0, MediaType: 1, TargetID: 2, ID: 2},
		PhysicalDisk{ControllerNum: 0, MediaType: 1, TargetID: 2, ID: 3},
		PhysicalDisk{ControllerNum: 0, MediaType: 1, TargetID: 2, ID: 4},
		PhysicalDisk{ControllerNum: 0, MediaType: 1, TargetID: 2, ID: 5},
		PhysicalDisk{ControllerNum: 0, MediaType: 1, TargetID: 2, ID: 6},
		PhysicalDisk{ControllerNum: 0, MediaType: 1, TargetID: 2, ID: 7},
		PhysicalDisk{ControllerNum: 0, MediaType: 1, TargetID: 2, ID: 8},
		PhysicalDisk{ControllerNum: 0, MediaType: 1, TargetID: 2, ID: 9},
		PhysicalDisk{ControllerNum: 0, MediaType: 1, TargetID: 2, ID: 10},
		PhysicalDisk{ControllerNum: 0, MediaType: 1, TargetID: 2, ID: 11},
		PhysicalDisk{ControllerNum: 0, MediaType: 1, TargetID: 2, ID: 12},
		PhysicalDisk{ControllerNum: 0, MediaType: 1, TargetID: 2, ID: 13},
		PhysicalDisk{ControllerNum: 0, MediaType: 1, TargetID: 2, ID: 14},
		PhysicalDisk{ControllerNum: 0, MediaType: 2, TargetID: 2, ID: 15},
		PhysicalDisk{ControllerNum: 0, MediaType: 2, TargetID: 2, ID: 16},
	}
	failed := false
	for i, pd := range actual.PDs {
		if pd != expected[i] {
			t.Logf("actual : %v != expected : %v", pd, expected[i])
			failed = true
		}
	}
	if failed {
		t.Fail()
	}
}
