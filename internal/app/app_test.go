package app

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"testing"

	"github.com/pacoorozco/networkdevices/internal/models"
)

var a App

func TestMain(m *testing.M) {
	a.Initialize()
	os.Exit(m.Run())
}

func TestGetWhenEmptyDeviceStore(t *testing.T) {
	req, _ := http.NewRequest("GET", "/devices", nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)

	if body := response.Body.String(); body != "[]" {
		t.Errorf("Expected an empty array. Got %s", body)
	}
}

func TestGetDevice(t *testing.T) {
	t.Run("Non Existent Device", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/devices/non-existent.domain.com.", nil)
		response := executeRequest(req)

		checkResponseCode(t, http.StatusNotFound, response.Code)

		var m map[string]string
		if err := json.Unmarshal(response.Body.Bytes(), &m); err != nil {
			t.Fatalf("error was not expected at this point")
		}
		if m["error"] != "Device not found" {
			t.Errorf("Expected the 'error' key of the response to be set to 'Device not found'. Got '%s'", m["error"])
		}
	})

	t.Run("Existent Device", func(t *testing.T) {
		want := models.Device{
			FQDN:    "host-1234.domain.com.",
			Model:   "ios-xe",
			Version: "11.2DS",
		}
		if err := a.Storer.SetDevice(want); err != nil {
			t.Fatalf("error was not expected at this point. err: %s", err)
		}

		req, _ := http.NewRequest("GET", "/devices/host-1234.domain.com.", nil)
		response := executeRequest(req)

		checkResponseCode(t, http.StatusOK, response.Code)
		checkResponseBodyContainsDevice(t, want, response.Body.Bytes())
	})

	t.Run("Existent Device wo/ Version", func(t *testing.T) {
		if err := a.Storer.SetDevice(models.Device{
			FQDN:    "host-1234.domain.com.",
			Model:   "ios-xe",
			Version: "",
		}); err != nil {
			t.Fatalf("error was not expected at this point. err: %s", err)
		}

		want := models.Device{
			FQDN:    "host-1234.domain.com.",
			Model:   "ios-xe",
			Version: "unknown",
		}

		req, _ := http.NewRequest("GET", "/devices/host-1234.domain.com.", nil)
		response := executeRequest(req)

		checkResponseCode(t, http.StatusOK, response.Code)
		checkResponseBodyContainsDevice(t, want, response.Body.Bytes())
	})
}

func TestCreateDevice(t *testing.T) {
	testCases := []struct {
		name         string
		input        string
		wantHttpCode int
	}{
		{
			name:         "Correct Device",
			input:        `{"fqdn":"test.domain.com.", "model": "ios-xe", "version": "1.2"}`,
			wantHttpCode: http.StatusCreated,
		},
		{
			name:         "Correct Device wo/ Version",
			input:        `{"fqdn":"test.domain.com.", "model": "ios-xe"}`,
			wantHttpCode: http.StatusCreated,
		},
		{
			name:         "Invalid Model Should Fail",
			input:        `{"fqdn":"test.domain.com.", "model": "invalid", "version": "1.2"}`,
			wantHttpCode: http.StatusBadRequest,
		},
		{
			name:         "Invalid FQDN Should Fail",
			input:        `{"fqdn":"", "model": "ios-xe", "version": "1.2"}`,
			wantHttpCode: http.StatusBadRequest,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var jsonStr = []byte(tc.input)
			req, _ := http.NewRequest("POST", "/devices", bytes.NewBuffer(jsonStr))
			req.Header.Set("Content-Type", "application/json")

			response := executeRequest(req)
			checkResponseCode(t, tc.wantHttpCode, response.Code)
		})
	}
}

func TestUpdateDevice(t *testing.T) {
	_ = addDevices(1)

	req, _ := http.NewRequest("GET", "/devices/host-0.domain.com.", nil)
	response := executeRequest(req)
	var originalDevice map[string]interface{}
	if err := json.Unmarshal(response.Body.Bytes(), &originalDevice); err != nil {
		t.Fatalf("error was not expected at this point")
	}

	var jsonStr = []byte(`{"fqdn":"host-0.domain.com.", "model": "nx-os", "version": "1.2"}`)
	req, _ = http.NewRequest("PUT", "/devices", bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")

	response = executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)

	var m map[string]interface{}
	if err := json.Unmarshal(response.Body.Bytes(), &m); err != nil {
		t.Fatalf("error was not expected at this point")
	}

	if m["fqdn"] != originalDevice["fqdn"] {
		t.Errorf("Expected the FQDN to remain the same (%v). Got %v", originalDevice["fqdn"], m["fqdn"])
	}

	if m["model"] == originalDevice["model"] {
		t.Errorf("Expected the model to change from '%v' to 'nx-os'. Got '%v'", originalDevice["model"], m["model"])
	}

	if m["version"] == originalDevice["version"] {
		t.Errorf("Expected the version to change from '%v' to '1.2'. Got '%v'", originalDevice["version"], m["version"])
	}
}

func TestDeleteDevice(t *testing.T) {
	_ = addDevices(1)

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

func checkResponseBodyContainsDevice(t *testing.T, expected models.Device, body []byte) {
	var got map[string]interface{}
	if err := json.Unmarshal(body, &got); err != nil {
		t.Fatalf("error was not expected at this point")
	}


	if expected.FQDN != got["fqdn"] {
		t.Errorf("FQDN: want: '%s', got: '%v'", expected.FQDN, got["fqdn"])
	}

	if expected.Model != got["model"] {
		t.Errorf("Model: want: '%s', got: '%v'", expected.Model, got["model"])
	}

	if expected.Version != got["version"] {
		t.Errorf("Version: want: '%s', got: '%v'", expected.Version, got["version"])
	}
}

func addDevices(count int) []models.Device {
	if count < 1 {
		count = 1
	}

	devices := make([]models.Device, 0)
	for i := 0; i < count; i++ {
		hostname := "host-" + strconv.Itoa(i) + ".domain.com."
		d := models.Device{
			FQDN:    hostname,
			Model:   "ios-xe",
			Version: "v1." + strconv.Itoa(i),
		}
		_ = a.Storer.SetDevice(d)
		devices = append(devices, d)
	}

	return devices
}
