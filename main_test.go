package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"testing"
)

var a App

func TestMain(m *testing.M) {
	a.Initialize()
	os.Exit(m.Run())
}

func TestEmptyTable(t *testing.T) {
	req, _ := http.NewRequest("GET", "/devices", nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)

	if body := response.Body.String(); body != "[]" {
		t.Errorf("Expected an empty array. Got %s", body)
	}
}

func executeRequest(req *http.Request) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	a.Router.ServeHTTP(rr, req)

	return rr
}

func checkResponseCode(t *testing.T, expected, actual int) {
	if expected != actual {
		t.Errorf("Expected response code %d. Got %d\n", expected, actual)
	}
}

func TestGetNonExistentDevice(t *testing.T) {
	req, _ := http.NewRequest("GET", "/devices/non-existent.domain.com.", nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusNotFound, response.Code)

	var m map[string]string
	json.Unmarshal(response.Body.Bytes(), &m)
	if m["error"] != "Device not found" {
		t.Errorf("Expected the 'error' key of the response to be set to 'Device not found'. Got '%s'", m["error"])
	}
}

func TestCreateDevice(t *testing.T) {
	var jsonStr = []byte(`{"fqdn":"test.domain.com.", "model": "os-xe", "version": "1.2"}`)
	req, _ := http.NewRequest("POST", "/devices", bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")

	response := executeRequest(req)
	checkResponseCode(t, http.StatusCreated, response.Code)

	var m map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &m)

	if m["fqdn"] != "test.domain.com." {
		t.Errorf("Expected Device FQDN to be 'test.domain.com.'. Got '%v'", m["fqdn"])
	}

	if m["model"] != "os-xe" {
		t.Errorf("Expected Device model to be 'os-xe'. Got '%v'", m["model"])
	}

	if m["version"] != "1.2" {
		t.Errorf("Expected Device version to be '1.2'. Got '%v'", m["version"])
	}
}

func TestGetDevice(t *testing.T) {
	addDevices(1)

	req, _ := http.NewRequest("GET", "/devices/host-0.domain.com.", nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)
}

func addDevices(count int) {
	if count < 1 {
		count = 1
	}

	for i := 0; i < count; i++ {
		hostname := "host-" + strconv.Itoa(i) + ".domain.com."
		a.DB.Set(hostname, Device{
			FQDN:    hostname,
			Model:   "os-xe",
			Version: "v1." + strconv.Itoa(i),
		})
	}
}

func TestUpdateDevice(t *testing.T) {
	addDevices(1)

	req, _ := http.NewRequest("GET", "/devices/host-0.domain.com.", nil)
	response := executeRequest(req)
	var originalDevice map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &originalDevice)

	var jsonStr = []byte(`{"fqdn":"host-0.domain.com.", "model": "os-nx", "version": "1.2"}`)
	req, _ = http.NewRequest("PUT", "/devices", bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")

	response = executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)

	var m map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &m)

	if m["fqdn"] != originalDevice["fqdn"] {
		t.Errorf("Expected the FQDN to remain the same (%v). Got %v", originalDevice["fqdn"], m["fqdn"])
	}

	if m["model"] == originalDevice["model"] {
		t.Errorf("Expected the model to change from '%v' to 'os-nx'. Got '%v'", originalDevice["model"], m["model"])
	}

	if m["version"] == originalDevice["version"] {
		t.Errorf("Expected the version to change from '%v' to '1.2'. Got '%v'", originalDevice["version"], m["version"])
	}
}

func TestDeleteDevice(t *testing.T) {
	addDevices(1)

	req, _ := http.NewRequest("GET", "/devices/host-0.domain.com.", nil)
	response := executeRequest(req)
	checkResponseCode(t, http.StatusOK, response.Code)

	req, _ = http.NewRequest("DELETE", "/devices/host-0.domain.com.", nil)
	response = executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)

	req, _ = http.NewRequest("GET", "/devices/host-0.domain.com.", nil)
	response = executeRequest(req)
	checkResponseCode(t, http.StatusNotFound, response.Code)
}
