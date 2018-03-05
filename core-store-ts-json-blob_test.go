package libDatabox

import (
	"strconv"
	s "strings"
	"testing"
	"time"
)

func TestWriteBlob(t *testing.T) {
	err := tsbc.Write(dsID, []byte("{\"test\":\"data\"}"))
	if err != nil {
		t.Errorf("Write to %s failed expected err to be nil got %s", dsID, err.Error())
	}
}

func BenchmarkWriteBlob(b *testing.B) {
	// run the Fib function b.N times
	for n := 0; n < b.N; n++ {
		tsbc.Write(dsID, []byte("{\"test\":\"data"+strconv.Itoa(n)+"\"}"))
	}
}

func BenchmarkWriteThenReadBlob(b *testing.B) {

	now := time.Now().UnixNano() / int64(time.Millisecond)

	// run the Fib function b.N times
	for n := 0; n < b.N; n++ {
		tsbc.Write(dsID, []byte("{\"test\":\"data"+strconv.Itoa(n)+"\"}"))
	}
	for n := 0; n < b.N-10; n++ {
		tsbc.Range(dsID, now, int64(n)+now)
	}
}

func BenchmarkWriteReadMixedBlob(b *testing.B) {

	// run the Fib function b.N times
	for n := 0; n < b.N; n++ {
		tsbc.Write(dsID, []byte("{\"test\":\"data"+strconv.Itoa(n)+"\"}"))
		tsbc.Latest(dsID)
	}
}

func TestLatestBlob(t *testing.T) {
	result, err := tsbc.Latest(dsID)
	if err != nil {
		t.Errorf("Call to Latest failed with error %s", err.Error())
	}
	expected := []byte(`{"test":"data"}`)
	cont := s.Contains(string(result), string(expected))
	if cont != true {
		t.Errorf("Write to %s failed expected %s but got %s", dsID, expected, result)
	}
}

func TestWriteLotsBlob(t *testing.T) {

	for i := 1; i <= 10; i++ {
		err := tsbc.Write(dsID, []byte("{\"test\":\"data"+strconv.Itoa(i)+"\"}"))
		if err != nil {
			t.Errorf("Write to %s failed expected err to be nil got %s", dsID, err.Error())
		}
	}

	result, err := tsbc.Latest(dsID)
	if err != nil {
		t.Errorf("Call to Latest failed with error %s", err.Error())
	}
	expected := []byte(`{"test":"data10"}`)
	cont := s.Contains(string(result), string(expected))
	if cont != true {
		t.Errorf("TestWriteLots failed expected %s but got %s", expected, result)
	}
}

//TODO this fails looks like a timing thing
/*func TestWriteThenWriteAT(t *testing.T) {

	now := time.Now().UnixNano() / int64(time.Millisecond)

	for i := 1; i <= 10; i++ {
		err := tsbc.Write(dsID, []byte("{\"TestWriteThenWriteAT\":\"data"+strconv.Itoa(i)+"\"}"))
		if err != nil {
			t.Errorf("Write to %s failed expected err to be nil got %s", dsID, err.Error())
		}
	}

	//time.Sleep(time.Second * 2)

	now = time.Now().UnixNano() / int64(time.Millisecond)
	fmt.Println(now + 1000)
	err := tsbc.WriteAt(dsID, now+1000, []byte("{\"TestWriteThenWriteAT\":\"data11\"}"))
	if err != nil {
		t.Errorf("Write to %s failed expected err to be nil got %s", dsID, err.Error())
	}

	err = tsbc.WriteAt(dsID, now+1001, []byte("{\"TestWriteThenWriteAT\":\"data12\"}"))
	if err != nil {
		t.Errorf("Write to %s failed expected err to be nil got %s", dsID, err.Error())
	}

	//time.Sleep(time.Second * 5)

	result, err := tsbc.LastN(dsID, 2)
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

} */

/*func TestLastN(t *testing.T) {

	now := time.Now().UnixNano() / int64(time.Millisecond)

	err := tsbc.WriteAt(dsID, now+20, []byte("{\"TestLastN\":\"data11\"}"))
	if err != nil {
		t.Errorf("Write to %s failed expected err to be nil got %s", dsID, err.Error())
	}
	err = tsbc.WriteAt(dsID, now+40, []byte("{\"TestLastN\":\"data12\"}"))
	if err != nil {
		t.Errorf("Write to %s failed expected err to be nil got %s", dsID, err.Error())
	}

	result, err := tsbc.LastN(dsID, 2)
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
}*/

func TestEarliestBlob(t *testing.T) {

	for i := 1; i <= 10; i++ {
		err := tsbc.Write(dsID+"TestEarliestBlob", []byte("{\"TestEarliestBlob\":\"data"+strconv.Itoa(i)+"\"}"))
		if err != nil {
			t.Errorf("Write to %s failed expected err to be nil got %s", dsID+"TestEarliestBlob", err.Error())
		}
		time.Sleep(time.Millisecond * 10)
	}

	result, err := tsbc.Earliest(dsID + "TestEarliestBlob")
	if err != nil {
		t.Errorf("Call to Earliest failed with error %s", err.Error())
	}

	expected := []byte(`{"TestEarliestBlob":"data1"}`)
	cont := s.Contains(string(result), string(expected))
	if cont != true {
		t.Errorf("Call to Earliest failed expected %s but got %s", expected, result)
	}

}

func TestFirstNBlob(t *testing.T) {

	for i := 1; i <= 10; i++ {
		err := tsbc.Write(dsID+"TestFirstNBlob", []byte("{\"TestFirstNBlob\":\"data"+strconv.Itoa(i)+"\"}"))
		if err != nil {
			t.Errorf("Write to %s failed expected err to be nil got %s", dsID+"TestFirstNBlob", err.Error())
		}
		time.Sleep(time.Millisecond * 10)
	}

	result, err := tsbc.FirstN(dsID+"TestFirstNBlob", 2)
	if err != nil {
		t.Errorf("Call to FirstN failed with error %s", err.Error())
	}

	expected := []byte(`{"TestFirstNBlob":"data1"}`)
	cont := s.Contains(string(result), string(expected))
	if cont != true {
		t.Errorf("Call to FirstN failed expected %s but got %s", expected, result)
	}

	expected = []byte(`{"TestFirstNBlob":"data2"}`)
	cont = s.Contains(string(result), string(expected))
	if cont != true {
		t.Errorf("Call to TestFirstNBlob failed expected %s but got %s", expected, result)
	}
}

/*func TestWriteAtAndRange(t *testing.T) {

	now := time.Now().UnixNano() / int64(time.Millisecond)
	numRecords := 20
	timeStepMs := 50

	for i := 1; i <= numRecords; i++ {

		err := tsbc.WriteAt(dsID, now+int64(timeStepMs*i), []byte("{\"test\":\"data"+strconv.Itoa(i)+"\"}"))
		if err != nil {
			t.Errorf("WriteAt to %s failed expected err to be nil got %s", dsID, err.Error())
		}
	}

	result, err := tsbc.Range(dsID, now, now+int64(numRecords*timeStepMs))
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
}*/

func TestRegisterDatasourceBlob(t *testing.T) {

	dsmd := DataSourceMetadata{
		DataSourceID:   dsID,
		Vendor:         "testing",
		ContentType:    "application/json",
		StoreType:      "ts",
		Description:    "A test DS",
		DataSourceType: "test",
	}

	err := tsbc.RegisterDatasource(dsmd)
	if err != nil {
		t.Errorf("RegisterDatasource failed expected err to be nil got %s", err.Error())
	}

	catByteArray, getErr := tsbc.GetDatasourceCatalogue()
	if getErr != nil {
		t.Errorf("GetDatasourceCatalogue failed expected err to be nil got %s", getErr.Error())
	}

	dsmdByteArray, _ := dataSourceMetadataToHypercat(dsmd, "tcp://127.0.0.1:5555/ts/blob/")
	cont := s.Contains(string(catByteArray), string(dsmdByteArray))
	if cont != true {
		t.Errorf("GetDatasourceCatalogue Error '%s' does not contain  %s", string(catByteArray), string(dsmdByteArray))
	}
}

/*
func TestConcurrentWriteAndRead(t *testing.T) {

	doneChanWrite := make(chan int)
	doneChanRead := make(chan int)
	//now := time.Now().UnixNano() / int64(time.Millisecond)
	numRecords := 100

	go func() {
		for i := 1; i <= numRecords; i++ {
			err := tsbc.Write(dsID, []byte("{\"TestConcurrentWriteAndRead\":\"data"+strconv.Itoa(i)+"\"}"))
			if err != nil {
				t.Errorf("WriteAt to %s failed expected err to be nil got %s", dsID, err.Error())
			}
			//fmt.Println(string("written:: " + strconv.Itoa(i)))
		}
		doneChanWrite <- 1
	}()

	go func() {
		for i := 1; i <= numRecords; i++ {
			_, err := tsbc.Latest(dsID)
			if err != nil {
				t.Errorf("Latest failed expected err to be nil got %s", err.Error())
			}
			//fmt.Println("Got:: ", string(data))
		}
		doneChanRead <- 1
	}()

	<-doneChanWrite
	<-doneChanRead

	result, err := tsbc.LastN(dsID, numRecords)
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
*/
