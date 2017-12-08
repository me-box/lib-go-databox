package libDatabox

import (
	"fmt"
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

func BenchmarkWrite(b *testing.B) {
	// run the Fib function b.N times
	for n := 0; n < b.N; n++ {
		tsc.Write(dsID, []byte("{\"test\":\"data"+strconv.Itoa(n)+"\"}"))
	}
}

func BenchmarkWriteThenRead(b *testing.B) {

	now := time.Now().UnixNano() / int64(time.Millisecond)

	// run the Fib function b.N times
	for n := 0; n < b.N; n++ {
		tsc.Write(dsID, []byte("{\"test\":\"data"+strconv.Itoa(n)+"\"}"))
	}
	for n := 0; n < b.N-10; n++ {
		tsc.Range(dsID, now, int64(n)+now)
	}
}

func BenchmarkWriteReadMixed(b *testing.B) {

	// run the Fib function b.N times
	for n := 0; n < b.N; n++ {
		tsc.Write(dsID, []byte("{\"test\":\"data"+strconv.Itoa(n)+"\"}"))
		tsc.Latest(dsID)
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

func TestWriteLots(t *testing.T) {

	for i := 1; i <= 10; i++ {
		err := tsc.Write(dsID, []byte("{\"test\":\"data"+strconv.Itoa(i)+"\"}"))
		if err != nil {
			t.Errorf("Write to %s failed expected err to be nil got %s", dsID, err.Error())
		}
	}

	result, err := tsc.Latest(dsID)
	if err != nil {
		t.Errorf("Call to Latest failed with error %s", err.Error())
	}
	expected := []byte(`{"test":"data10"}`)
	cont := s.Contains(string(result), string(expected))
	if cont != true {
		t.Errorf("TestWriteLots failed expected %s but got %s", expected, result)
	}
}

func TestWriteThenWriteAT(t *testing.T) {

	now := time.Now().UnixNano() / int64(time.Millisecond)

	for i := 1; i <= 10; i++ {
		//err := tsc.WriteAt(dsID, now+int64(i), []byte("{\"TestWriteThenWriteAT\":\"data"+strconv.Itoa(i)+"\"}"))
		err := tsc.Write(dsID, []byte("{\"TestWriteThenWriteAT\":\"data"+strconv.Itoa(i)+"\"}"))
		if err != nil {
			t.Errorf("Write to %s failed expected err to be nil got %s", dsID, err.Error())
		}
	}

	//time.Sleep(time.Second * 2)

	now = time.Now().UnixNano() / int64(time.Millisecond)
	fmt.Println(now + 1000)
	err := tsc.WriteAt(dsID, now+1000, []byte("{\"TestWriteThenWriteAT\":\"data11\"}"))
	//err := tsc.Write(dsID, []byte("{\"TestWriteThenWriteAT\":\"data11\"}"))
	if err != nil {
		t.Errorf("Write to %s failed expected err to be nil got %s", dsID, err.Error())
	}

	err = tsc.WriteAt(dsID, now+1001, []byte("{\"TestWriteThenWriteAT\":\"data12\"}"))
	//err = tsc.Write(dsID, []byte("{\"TestWriteThenWriteAT\":\"data12\"}"))
	if err != nil {
		t.Errorf("Write to %s failed expected err to be nil got %s", dsID, err.Error())
	}

	//time.Sleep(time.Second * 5)

	result, err := tsc.LastN(dsID, 2)
	if err != nil {
		t.Errorf("Call to LastN failed with error %s", err.Error())
	}

	expected := []byte(`{"TestWriteThenWriteAT":"data11"}`)
	cont := s.Contains(string(result), string(expected))
	if cont != true {
		t.Errorf("Call to LastN failed expected %s but got %s", expected, result)
	}

	expected = []byte(`{"TestWriteThenWriteAT":"data12"}`)
	cont = s.Contains(string(result), string(expected))
	if cont != true {
		t.Errorf("Call to LastN failed expected %s but got %s", expected, result)
	}

}

func TestLastN(t *testing.T) {

	now := time.Now().UnixNano() / int64(time.Millisecond)

	err := tsc.WriteAt(dsID, now+20, []byte("{\"TestLastN\":\"data11\"}"))
	if err != nil {
		t.Errorf("Write to %s failed expected err to be nil got %s", dsID, err.Error())
	}
	err = tsc.WriteAt(dsID, now+40, []byte("{\"TestLastN\":\"data12\"}"))
	if err != nil {
		t.Errorf("Write to %s failed expected err to be nil got %s", dsID, err.Error())
	}

	result, err := tsc.LastN(dsID, 2)
	if err != nil {
		t.Errorf("Call to LastN failed with error %s", err.Error())
	}

	expected := []byte(`{"TestLastN":"data11"}`)
	cont := s.Contains(string(result), string(expected))
	if cont != true {
		t.Errorf("Call to LastN failed expected %s but got %s", expected, result)
	}

	expected = []byte(`{"TestLastN":"data12"}`)
	cont = s.Contains(string(result), string(expected))
	if cont != true {
		t.Errorf("Call to LastN failed expected %s but got %s", expected, result)
	}
}

func TestEarliest(t *testing.T) {

	for i := 1; i <= 10; i++ {
		err := tsc.Write(dsID+"TestEarliest", []byte("{\"TestEarliest\":\"data"+strconv.Itoa(i)+"\"}"))
		if err != nil {
			t.Errorf("Write to %s failed expected err to be nil got %s", dsID+"TestEarliest", err.Error())
		}
	}

	result, err := tsc.Earliest(dsID + "TestEarliest")
	if err != nil {
		t.Errorf("Call to Earliest failed with error %s", err.Error())
	}

	expected := []byte(`{"TestEarliest":"data1"}`)
	cont := s.Contains(string(result), string(expected))
	if cont != true {
		t.Errorf("Call to Earliest failed expected %s but got %s", expected, result)
	}

}

func TestFirstN(t *testing.T) {

	for i := 1; i <= 10; i++ {
		err := tsc.Write(dsID+"TestFirstN", []byte("{\"TestFirstN\":\"data"+strconv.Itoa(i)+"\"}"))
		if err != nil {
			t.Errorf("Write to %s failed expected err to be nil got %s", dsID+"TestFirstN", err.Error())
		}
	}

	result, err := tsc.FirstN(dsID+"TestFirstN", 2)
	if err != nil {
		t.Errorf("Call to FirstN failed with error %s", err.Error())
	}

	expected := []byte(`{"TestFirstN":"data1"}`)
	cont := s.Contains(string(result), string(expected))
	if cont != true {
		t.Errorf("Call to FirstN failed expected %s but got %s", expected, result)
	}

	expected = []byte(`{"TestFirstN":"data2"}`)
	cont = s.Contains(string(result), string(expected))
	if cont != true {
		t.Errorf("Call to FirstN failed expected %s but got %s", expected, result)
	}
}

func TestWriteAtAndRange(t *testing.T) {

	now := time.Now().UnixNano() / int64(time.Millisecond)
	numRecords := 20
	timeStepMs := 50

	for i := 1; i <= numRecords; i++ {

		err := tsc.WriteAt(dsID, now+int64(timeStepMs*i), []byte("{\"test\":\"data"+strconv.Itoa(i)+"\"}"))
		if err != nil {
			t.Errorf("WriteAt to %s failed expected err to be nil got %s", dsID, err.Error())
		}
	}

	result, err := tsc.Range(dsID, now, now+int64(numRecords*timeStepMs))
	if err != nil {
		t.Errorf("Call to Range failed with error %s", err.Error())
	}

	expected := []byte(`{"test":"data1"}`)
	cont := s.Contains(string(result), string(expected))
	if cont != true {
		t.Errorf("TestWriteAtAndRange failed expected %s but got %s", expected, result)
	}
	expected = []byte(`{"test":"data5"}`)
	cont = s.Contains(string(result), string(expected))
	if cont != true {
		t.Errorf("TestWriteAtAndRange failed expected %s but got %s", expected, result)
	}
	expected = []byte(`{"test":"data` + strconv.Itoa(numRecords) + `"}`)
	cont = s.Contains(string(result), string(expected))
	if cont != true {
		t.Errorf("TestWriteAtAndRange failed expected %s but got %s", expected, result)
	}
}

func TestRegisterDatasource(t *testing.T) {

	dsmd := DataSourceMetadata{
		DataSourceID:   dsID,
		Vendor:         "testing",
		ContentType:    "application/json",
		StoreType:      "ts",
		Description:    "A test DS",
		DataSourceType: "test",
	}

	err := tsc.RegisterDatasource(dsmd)
	if err != nil {
		t.Errorf("RegisterDatasource failed expected err to be nil got %s", err.Error())
	}

	catByteArray, getErr := tsc.GetDatasourceCatalogue()
	if getErr != nil {
		t.Errorf("GetDatasourceCatalogue failed expected err to be nil got %s", getErr.Error())
	}

	dsmdByteArray, _ := dataSourceMetadataToHypercat(dsmd)
	cont := s.Contains(string(catByteArray), string(dsmdByteArray))
	if cont != true {
		t.Errorf("GetDatasourceCatalogue Error '%s' does not contain  %s", string(catByteArray), string(dsmdByteArray))
	}
}

func TestConcurrentWriteAndRead(t *testing.T) {

	doneChanWrite := make(chan int)
	doneChanRead := make(chan int)
	//now := time.Now().UnixNano() / int64(time.Millisecond)
	numRecords := 100

	go func() {
		for i := 1; i <= numRecords; i++ {
			err := tsc.Write(dsID, []byte("{\"TestConcurrentWriteAndRead\":\"data"+strconv.Itoa(i)+"\"}"))
			if err != nil {
				t.Errorf("WriteAt to %s failed expected err to be nil got %s", dsID, err.Error())
			}
			//fmt.Println(string("written:: " + strconv.Itoa(i)))
		}
		doneChanWrite <- 1
	}()

	go func() {
		for i := 1; i <= numRecords; i++ {
			_, err := tsc.Latest(dsID)
			if err != nil {
				t.Errorf("Latest failed expected err to be nil got %s", err.Error())
			}
			//fmt.Println("Got:: ", string(data))
		}
		doneChanRead <- 1
	}()

	<-doneChanWrite
	<-doneChanRead

	result, err := tsc.LastN(dsID, numRecords)
	if err != nil {
		t.Errorf("Call to LastN failed with error %s", err.Error())
	}
	for i := 1; i <= numRecords; i++ {
		expected := []byte("{\"TestConcurrentWriteAndRead\":\"data" + strconv.Itoa(i) + "\"}")
		cont := s.Contains(string(result), string(expected))
		if cont != true {
			t.Errorf("LastN Error '%s' does not contain  %s", string(result), string(expected))
			break
		}
	}
}
