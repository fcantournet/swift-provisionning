package main

import (
	"io/ioutil"
	"log"
	"testing"
)

func TestParseControllerXML(t *testing.T) {
	expected := []int{0, 1}
	xmldata, err := ioutil.ReadFile("./testdata_controllers.xml")
	if err != nil {
		log.Fatal(err.Error())
	}
	temp, err := parseControllerIDs(xmldata)
	actual := temp.ControllerIDs
	if len(actual) != len(expected) {
		t.Errorf("Test failed, expected: '%#v', got:  '%#v'", expected, actual)
	}
	for i := range expected {
		if actual[i] != expected[i] {
			t.Errorf("Test failed, expected: '%#v', got:  '%#v'", expected, actual)
		}
	}
}
