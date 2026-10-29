package main

import (
	"testing"
	"time"
)

func TestArchiveWindow(t *testing.T) {
	a, b, e := archiveWindow("BTCUSDT-5m-2024-02.zip")
	if e != nil || time.UnixMilli(a).UTC().Day() != 1 || time.UnixMilli(b).UTC().Month() != time.March {
		t.Fatal(a, b, e)
	}
	a, b, e = archiveWindow("BTCUSDT-5m-2024-12-29.zip")
	if e != nil || b-a != 86400000 {
		t.Fatal(a, b, e)
	}
	if _, _, e = archiveWindow("ETHUSDT-5m-2024-02.zip"); e == nil {
		t.Fatal("accepted wrong instrument")
	}
}
