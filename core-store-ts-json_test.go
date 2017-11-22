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

func TestLastN(t *testing.T) {

	//Using writeAt here cause odd behaviour when executed after TestWriteLots, works if run in isolation. Disabling for now.
	/*now := time.Now().UnixNano() / int64(time.Millisecond)

	err := tsc.WriteAt(dsID, now+20, []byte("{\"test\":\"data11\"}"))
	if err != nil {
		t.Errorf("Write to %s failed expected err to be nil got %s", dsID, err.Error())
	}
	err = tsc.WriteAt(dsID, now+40, []byte("{\"test\":\"data12\"}"))
	if err != nil {
		t.Errorf("Write to %s failed expected err to be nil got %s", dsID, err.Error())
	}*/

	err := tsc.Write(dsID, []byte("{\"test\":\"data11\"}"))
	if err != nil {
		t.Errorf("Write to %s failed expected err to be nil got %s", dsID, err.Error())
	}
	err = tsc.Write(dsID, []byte("{\"test\":\"data12\"}"))
	if err != nil {
		t.Errorf("Write to %s failed expected err to be nil got %s", dsID, err.Error())
	}

	//time.Sleep(time.Millisecond * 100)
	result, err := tsc.LastN(dsID, 2)
	if err != nil {
		t.Errorf("Call to LastN failed with error %s", err.Error())
	}

	expected := []byte(`{"test":"data11"}`)
	cont := s.Contains(string(result), string(expected))
	if cont != true {
		t.Errorf("Call to LastN failed expected %s but got %s", expected, result)
	}

	expected = []byte(`{"test":"data12"}`)
	cont = s.Contains(string(result), string(expected))
	if cont != true {
		t.Errorf("Call to LastN failed expected %s but got %s", expected, result)
	}
}

func TestEarliest(t *testing.T) {

	for i := 1; i <= 10; i++ {
		err := tsc.Write(dsID+"TestEarliest", []byte("{\"test\":\"data"+strconv.Itoa(i)+"\"}"))
		if err != nil {
			t.Errorf("Write to %s failed expected err to be nil got %s", dsID+"TestEarliest", err.Error())
		}
	}

	result, err := tsc.Earliest(dsID + "TestEarliest")
	if err != nil {
		t.Errorf("Call to Earliest failed with error %s", err.Error())
	}

	expected := []byte(`{"test":"data1"}`)
	cont := s.Contains(string(result), string(expected))
	if cont != true {
		t.Errorf("Call to Earliest failed expected %s but got %s", expected, result)
	}

}

func TestFirstN(t *testing.T) {

	for i := 1; i <= 10; i++ {
		err := tsc.Write(dsID+"TestFirstN", []byte("{\"test\":\"data"+strconv.Itoa(i)+"\"}"))
		if err != nil {
			t.Errorf("Write to %s failed expected err to be nil got %s", dsID+"TestFirstN", err.Error())
		}
	}

	result, err := tsc.FirstN(dsID+"TestFirstN", 2)
	if err != nil {
		t.Errorf("Call to FirstN failed with error %s", err.Error())
	}

	expected := []byte(`{"test":"data1"}`)
	cont := s.Contains(string(result), string(expected))
	if cont != true {
		t.Errorf("Call to FirstN failed expected %s but got %s", expected, result)
	}

	expected = []byte(`{"test":"data2"}`)
	cont = s.Contains(string(result), string(expected))
	if cont != true {
		t.Errorf("Call to FirstN failed expected %s but got %s", expected, result)
	}
}
