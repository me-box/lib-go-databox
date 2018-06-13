package libDatabox

import (
	"strconv"
	s "strings"
	"testing"
	"time"
)

func TestWriteBlob(t *testing.T) {
	err := StoreClient.TSBlobWrite(dsID, []byte("{\"test\":\"data\"}"))
	if err != nil {
		t.Errorf("Write to %s failed expected err to be nil got %s", dsID, err.Error())
	}
}

func BenchmarkWriteBlob(b *testing.B) {
	// run the Fib function b.N times
	for n := 0; n < b.N; n++ {
		StoreClient.TSBlobWrite(dsID, []byte("{\"test\":\"data"+strconv.Itoa(n)+"\"}"))
	}
}

func BenchmarkWriteThenReadBlob(b *testing.B) {

	now := time.Now().UnixNano() / int64(time.Millisecond)

	// run the Fib function b.N times
	for n := 0; n < b.N; n++ {
		StoreClient.TSBlobWrite(dsID, []byte("{\"test\":\"data"+strconv.Itoa(n)+"\"}"))
	}
	for n := 0; n < b.N-10; n++ {
		StoreClient.TSBlobRange(dsID, now, int64(n)+now)
	}
}

func BenchmarkWriteReadMixedBlob(b *testing.B) {

	// run the Fib function b.N times
	for n := 0; n < b.N; n++ {
		StoreClient.TSBlobWrite(dsID, []byte("{\"test\":\"data"+strconv.Itoa(n)+"\"}"))
		StoreClient.TSBlobLatest(dsID)
	}
}

func TestLatestBlob(t *testing.T) {
	result, err := StoreClient.TSBlobLatest(dsID)
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
		err := StoreClient.TSBlobWrite(dsID, []byte("{\"test\":\"data"+strconv.Itoa(i)+"\"}"))
		if err != nil {
			t.Errorf("Write to %s failed expected err to be nil got %s", dsID, err.Error())
		}
	}

	result, err := StoreClient.TSBlobLatest(dsID)
	if err != nil {
		t.Errorf("Call to Latest failed with error %s", err.Error())
	}
	expected := []byte(`{"test":"data10"}`)
	cont := s.Contains(string(result), string(expected))
	if cont != true {
		t.Errorf("TestWriteLots failed expected %s but got %s", expected, result)
	}
}

func TestWriteLengthBlob(t *testing.T) {

	numRecToWrite := 50
	_dsID := dsID + "TestWriteLengthBlob"
	for i := 1; i <= numRecToWrite; i++ {
		err := StoreClient.TSBlobWrite(_dsID, []byte("{\"test\":\"data"+strconv.Itoa(i)+"\"}"))
		if err != nil {
			t.Errorf("Write to %s failed expected err to be nil got %s", _dsID, err.Error())
		}
		time.Sleep(10 * time.Millisecond)
	}

	result, err := StoreClient.TSBlobLength(_dsID)
	if err != nil {
		t.Errorf("Call to Latest failed with error %s", err.Error())
	}
	if numRecToWrite != result {
		t.Errorf("TestWriteLots failed expected %d but got %d", numRecToWrite, result)
	}
}

//TODO this fails looks like a timing thing
/*func TestWriteThenWriteAT(t *testing.T) {

	now := time.Now().UnixNano() / int64(time.Millisecond)

	for i := 1; i <= 10; i++ {
		err := StoreClient.TSBlobWrite(dsID, []byte("{\"TestWriteThenWriteAT\":\"data"+strconv.Itoa(i)+"\"}"))
		if err != nil {
			t.Errorf("Write to %s failed expected err to be nil got %s", dsID, err.Error())
		}
	}

	//time.Sleep(time.Second * 2)

	now = time.Now().UnixNano() / int64(time.Millisecond)
	fmt.Println(now + 1000)
	err := StoreClient.TSBlobWriteAt(dsID, now+1000, []byte("{\"TestWriteThenWriteAT\":\"data11\"}"))
	if err != nil {
		t.Errorf("Write to %s failed expected err to be nil got %s", dsID, err.Error())
	}

	err = StoreClient.TSBlobWriteAt(dsID, now+1001, []byte("{\"TestWriteThenWriteAT\":\"data12\"}"))
	if err != nil {
		t.Errorf("Write to %s failed expected err to be nil got %s", dsID, err.Error())
	}

	//time.Sleep(time.Second * 5)

	result, err := StoreClient.TSBlobLastN(dsID, 2)
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

func TestLastNBlob(t *testing.T) {

	now := time.Now().UnixNano() / int64(time.Millisecond)

	err := StoreClient.TSBlobWriteAt(dsID+"TestLastN", now+20, []byte("{\"TestLastN\":\"data11\"}"))
	if err != nil {
		t.Errorf("Write to %s failed expected err to be nil got %s", dsID, err.Error())
	}
	err = StoreClient.TSBlobWriteAt(dsID+"TestLastN", now+40, []byte("{\"TestLastN\":\"data12\"}"))
	if err != nil {
		t.Errorf("Write to %s failed expected err to be nil got %s", dsID, err.Error())
	}

	result, err := StoreClient.TSBlobLastN(dsID+"TestLastN", 2)
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

func TestEarliestBlob(t *testing.T) {

	for i := 1; i <= 10; i++ {
		err := StoreClient.TSBlobWrite(dsID+"TestEarliestBlob", []byte("{\"TestEarliestBlob\":\"data"+strconv.Itoa(i)+"\"}"))
		if err != nil {
			t.Errorf("Write to %s failed expected err to be nil got %s", dsID+"TestEarliestBlob", err.Error())
		}
		time.Sleep(time.Millisecond * 10)
	}

	result, err := StoreClient.TSBlobEarliest(dsID + "TestEarliestBlob")
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
		err := StoreClient.TSBlobWrite(dsID+"TestFirstNBlob", []byte("{\"TestFirstNBlob\":\"data"+strconv.Itoa(i)+"\"}"))
		if err != nil {
			t.Errorf("Write to %s failed expected err to be nil got %s", dsID+"TestFirstNBlob", err.Error())
		}
		time.Sleep(time.Millisecond * 10)
	}

	result, err := StoreClient.TSBlobFirstN(dsID+"TestFirstNBlob", 2)
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

func TestWriteAtAndRangeBlob(t *testing.T) {

	now := time.Now().UnixNano() / int64(time.Millisecond)
	numRecords := 20
	timeStepMs := 50

	_dsID := dsID + "TestWriteAtAndRangeBlob"

	for i := 1; i <= numRecords; i++ {

		err := StoreClient.TSBlobWriteAt(_dsID, now+int64(timeStepMs*i), []byte("{\"test\":\"data"+strconv.Itoa(i)+"\"}"))
		if err != nil {
			t.Errorf("WriteAt to %s failed expected err to be nil got %s", _dsID, err.Error())
		}
	}

	result, err := StoreClient.TSBlobRange(_dsID, now, now+int64(numRecords*timeStepMs))
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

//TO FIX THIS
/*func TestRegisterDatasourceBlob(t *testing.T) {

	dsmd := DataSourceMetadata{
		DataSourceID:   dsID,
		Vendor:         "testing",
		ContentType:    "application/json",
		StoreType:      "ts",
		Description:    "A test DS",
		DataSourceType: "test",
	}

	err := StoreClient.RegisterDatasource(dsmd)
	if err != nil {
		t.Errorf("RegisterDatasource failed expected err to be nil got %s", err.Error())
	}

	catByteArray, getErr := StoreClient.GetStoreDataSourceCatalogue()
	if getErr != nil {
		t.Errorf("GetDatasourceCatalogue failed expected err to be nil got %s", getErr.Error())
	}

	dsmdByteArray, _ := StoreClient.dataSourceMetadataToHypercat(dsmd, "tcp://127.0.0.1:5555/ts/blob/")
	cont := s.Contains(string(catByteArray), string(dsmdByteArray))
	if cont != true {
		t.Errorf("GetDatasourceCatalogue Error '%s' does not contain  %s", string(catByteArray), string(dsmdByteArray))
	}
}*/

func TestConcurrentWriteAndReadBlob(t *testing.T) {

	doneChanWrite := make(chan int)
	doneChanRead := make(chan int)
	//now := time.Now().UnixNano() / int64(time.Millisecond)
	numRecords := 100

	go func() {
		for i := 1; i <= numRecords; i++ {
			err := StoreClient.TSBlobWrite(dsID, []byte("{\"TestConcurrentWriteAndRead\":\"data"+strconv.Itoa(i)+"\"}"))
			if err != nil {
				t.Errorf("WriteAt to %s failed expected err to be nil got %s", dsID, err.Error())
			}
			//fmt.Println(string("written:: " + strconv.Itoa(i)))
		}
		doneChanWrite <- 1
	}()

	go func() {
		for i := 1; i <= numRecords; i++ {
			_, err := StoreClient.TSBlobLatest(dsID)
			if err != nil {
				t.Errorf("Latest failed expected err to be nil got %s", err.Error())
			}
			//fmt.Println("Got:: ", string(data))
		}
		doneChanRead <- 1
	}()

	<-doneChanWrite
	<-doneChanRead

	result, err := StoreClient.TSBlobLastN(dsID, numRecords)
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

func TestObserveBlob(t *testing.T) {

	doneChanWrite := make(chan int)
	//doneChanRead := make(chan int)
	//now := time.Now().UnixNano() / int64(time.Millisecond)
	startAt := 0
	numRecords := 10

	receivedData := [][]byte{}

	go func() {
		dataChan, err := StoreClient.TSBlobObserve(dsID)
		if err != nil {
			t.Errorf("Observing %s failed expected err to be nil got %s", dsID, err.Error())
		}

		for data := range dataChan {
			receivedData = append(receivedData, data)
			t.Log("received:: " + string(data))
		}

	}()

	//Observe take a bit of time to register we miss some values if we dont wait before writing
	time.Sleep(time.Second)

	go func() {
		for i := startAt; i <= numRecords; i++ {
			err := StoreClient.TSBlobWrite(dsID, []byte("{\"test\":"+strconv.Itoa(i)+", \"data\":\"AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA\"}"))
			if err != nil {
				t.Errorf("WriteAt to %s failed expected err to be nil got %s", dsID, err.Error())
			}
			t.Log(string("written:: " + strconv.Itoa(i)))
		}
		// we miss some values if we dont wait before saying we are done!
		time.Sleep(time.Second * 3)
		doneChanWrite <- 1
	}()

	<-doneChanWrite
	if len(receivedData) < numRecords {
		t.Errorf("receivedData Error:  receivedData should contain '%d' items but contains  %d", numRecords, len(receivedData))
	}
	for i := startAt; i <= numRecords; i++ {
		expected := []byte("{\"test\":" + strconv.Itoa(i) + ", \"data\":\"AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA\"}")
		cont := s.Contains(string(receivedData[i]), string(expected))
		t.Log(string(receivedData[i]))
		if cont != true {
			t.Errorf("receivedData Error '%s' does not contain  %s", string(receivedData[i]), string(expected))
			break
		}
	}

}
