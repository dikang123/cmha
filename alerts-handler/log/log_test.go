package log

import (
	"os"
	"testing"
)

func TestLog(t *testing.T) {
	logfile, err := os.OpenFile("./logtest.log", os.O_CREATE|os.O_RDWR|os.O_APPEND, 0666)
	if err != nil {
		t.Fatalf("open logfile fail", err)
	}
	defer logfile.Close()
	if err := LogInit("warn", logfile); err != nil {
		t.Fatalf("log init fail", err)
	}
	Info("info", "testing")
	Warn("warning", "testing")
	Pannic("pannic", "testing")

}
