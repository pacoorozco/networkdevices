package models

import (
	"reflect"
	"testing"
)

func TestDevice_Validate(t *testing.T) {
	testCases := []struct {
		name        string
		input       Device
		errExpected bool
	}{
		{
			name: "Correct Device",
			input: Device{
				FQDN:    "hostname.domain.com.",
				Model:   "ios-xr",
				Version: "11.2STY",
			},
			errExpected: false,
		},
		{
			name: "Correct Device wo/ Version",
			input: Device{
				FQDN:  "hostname.domain.com.",
				Model: "ios-xr",
			},
			errExpected: false,
		},
		{
			name:        "Empty Device",
			input:       Device{},
			errExpected: true,
		},
		{
			name: "Device w/ invalid Model",
			input: Device{
				FQDN:  "hostname.domain.com.",
				Model: "non-valid",
			},
			errExpected: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.input.Validate()
			if err == nil && tc.errExpected {
				t.Fatalf("error was expected")
			}
			if err != nil && !tc.errExpected {
				t.Fatalf("error was not expected: err: %s", err)
			}
		})
	}
}

func TestDevice_Presenter(t *testing.T) {
	testCases := []struct {
		name  string
		input Device
		want  Device
	}{
		{
			name: "Without Changes",
			input: Device{
				FQDN:    "hostname.domain.com.",
				Model:   "ios-xr",
				Version: "11.2STY",
			},
			want: Device{
				FQDN:    "hostname.domain.com.",
				Model:   "ios-xr",
				Version: "11.2STY",
			},
		},
		{
			name: "Formatted Version",
			input: Device{
				FQDN:  "hostname.domain.com.",
				Model: "ios-xr",
			},
			want: Device{
				FQDN:  "hostname.domain.com.",
				Model: "ios-xr",
				Version: "unknown",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := tc.input.Presenter()
			if !reflect.DeepEqual(tc.want, got) {
				t.Errorf("want: %v, got: %v", tc.want, got)
			}
		})
	}
}
