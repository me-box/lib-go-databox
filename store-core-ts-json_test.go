package libDatabox

import (
	"os"
	"strconv"
	s "strings"
	"testing"
	"time"
)

func TestMain(m *testing.M) {
	Setup()
	retCode := m.Run()
	Teardown()
	os.Exit(retCode)
}

var tsc JSONTimeSeries_0_2_0
var dsID string

func Setup() {

	var err error
	tsc, err = NewJSONTimeSeriesClient("tcp://127.0.0.1:5555", false)
	if err != nil {
		panic("Cant connect to Zest server. Did you start one? " + err.Error())
	}

	dsID = "test" + strconv.Itoa(int(time.Now().UnixNano()/int64(time.Millisecond)))
}

func Teardown() {
	//todo
}
func TestWrite(t *testing.T) {
	err := tsc.Write(dsID, []byte("{\"test\":\"data\"}"))
	if err != nil {
		t.Errorf("Write to %s failed expected err to be nil got %s", dsID, err.Error())
	}
}

func TestLatest(t *testing.T) {
	result, err := tsc.Latest(dsID)
	if err != nil {
		t.Errorf("Call to Latest failed with error %s", err.Error())
	}
	expected := []byte(`{"test":"data"}`)
	cont := s.Contains(string(result), string(expected))
	if cont != true {
		t.Errorf("Write to %s failed expected %s but got %s", dsID, expected, result)
	}
}
