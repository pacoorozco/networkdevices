package internal

import "testing"

func TestInMemoryDAL_Set(t *testing.T) {
	dal := NewInMemory()

	want := Device{
		FQDN:    "hostname.domain.com.",
		Model:   "ios-xr",
		Version: "11.2DS",
	}

	dal.Set(want.FQDN, want)

	got := dal.items[want.FQDN]

	t.Errorf("want: %v, got: %v", want, got)


}
