package libDatabox

import (
	"os"
	"strconv"
	"testing"
	"time"
)

func TestMain(m *testing.M) {
	Setup()
	retCode := m.Run()
	Teardown()
	os.Exit(retCode)
}

var tsbc JSONTimeSeriesBlob_0_3_0
var tsc JSONTimeSeries_0_3_0
var kvc JSONKeyValue_0_3_0
var kvcBin BinaryKeyValue_0_3_0
var kvcText TextKeyValue_0_3_0

var dsID string

func Setup() {

	var err error
	tsbc, err = NewJSONTimeSeriesBlobClient("tcp://127.0.0.1:5555", false)
	if err != nil {
		panic("Cant connect to Zest server. Did you start one? " + err.Error())
	}

	tsc, err = NewJSONTimeSeriesClient("tcp://127.0.0.1:5555", false)
	if err != nil {
		panic("Cant connect to Zest server. Did you start one? " + err.Error())
	}

	kvc, err = NewJSONKeyValueClient("tcp://127.0.0.1:5555", false)
	if err != nil {
		panic("Cant connect to Zest server. Did you start one? " + err.Error())
	}

	kvcBin, err = NewBinaryKeyValueClient("tcp://127.0.0.1:5555", false)
	if err != nil {
		panic("Cant connect to Zest server. Did you start one? " + err.Error())
	}

	kvcText, err = NewTextKeyValueClient("tcp://127.0.0.1:5555", false)
	if err != nil {
		panic("Cant connect to Zest server. Did you start one? " + err.Error())
	}

	dsID = "test" + strconv.Itoa(int(time.Now().UnixNano()/int64(time.Millisecond)))
}

func Teardown() {
	//todo
}
