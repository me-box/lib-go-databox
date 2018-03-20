package libDatabox

import (
	"strconv"
	s "strings"
	"testing"
	"time"
)

func TestKVWrite(t *testing.T) {
	err := kvc.Write(dsID, "key1", []byte("{\"value\":3.1415}"))
	if err != nil {
		t.Errorf("Write to %s failed expected err to be nil got %s", dsID, err.Error())
	}
}

func TestKVRead(t *testing.T) {
	err := kvc.Write(dsID, "key2", []byte("{\"value\":42}"))
	if err != nil {
		t.Errorf("Write to %s failed expected err to be nil got %s", dsID, err.Error())
	}

	result, err := kvc.Read(dsID, "key2")
	if err != nil {
		t.Errorf("Write to %s failed expected err to be nil got %s", dsID, err.Error())
	}

	expected := []byte("{\"value\":42}")
	cont := s.Contains(string(result), string(expected))
	if cont != true {
		t.Errorf("TestWriteLots failed expected %s but got %s", expected, result)
	}
}

func TestKVMutiKeys(t *testing.T) {

	err := kvc.Write(dsID, "key1", []byte("{\"value\":\"some random string\"}"))
	if err != nil {
		t.Errorf("Write to %s failed expected err to be nil got %s", dsID, err.Error())
	}

	err = kvc.Write(dsID, "key2", []byte("{\"value\":42}"))
	if err != nil {
		t.Errorf("Write to %s failed expected err to be nil got %s", dsID, err.Error())
	}

	result, err := kvc.Read(dsID, "key2")
	if err != nil {
		t.Errorf("Write to %s failed expected err to be nil got %s", dsID, err.Error())
	}

	expected := []byte("{\"value\":42}")
	cont := s.Contains(string(result), string(expected))
	if cont != true {
		t.Errorf("TestWriteLots failed expected %s but got %s", expected, result)
	}

	result, err = kvc.Read(dsID, "key1")
	if err != nil {
		t.Errorf("Write to %s failed expected err to be nil got %s", dsID, err.Error())
	}

	expected = []byte("{\"value\":\"some random string\"}")
	cont = s.Contains(string(result), string(expected))
	if cont != true {
		t.Errorf("TestWriteLots failed expected %s but got %s", expected, result)
	}
}

func TestObserveKeyKV(t *testing.T) {

	doneChanWrite := make(chan int)
	//doneChanRead := make(chan int)
	//now := time.Now().UnixNano() / int64(time.Millisecond)
	startAt := 0
	numRecords := 100

	receivedData := []JsonObserveResponse{}

	go func() {
		dataChan, err := kvc.ObserveKey(dsID, "observeTest")
		if err != nil {
			t.Errorf("Observing %s failed expected err to be nil got %s", dsID, err.Error())
		}

		for data := range dataChan {
			receivedData = append(receivedData, data)
			t.Log("received:: ", string(data.Json), data.TimestampMS)
		}

	}()

	//Observe take a bit of time to register we miss some values if we dont wait before writing
	time.Sleep(time.Second)

	go func() {
		for i := startAt; i <= numRecords; i++ {
			err := kvc.Write(dsID, "observeTest", []byte("{\"value\":"+strconv.Itoa(i)+"}"))
			if err != nil {
				t.Errorf("WriteAt to %s failed expected err to be nil got %s", dsID, err.Error())
			}
			t.Log(string("written:: " + strconv.Itoa(i)))
		}
		doneChanWrite <- 1
	}()

	<-doneChanWrite

	for i := startAt; i <= numRecords; i++ {
		expected := []byte("{\"value\":" + strconv.Itoa(i) + "}")
		cont := s.Contains(string(receivedData[i].Json), string(expected))
		t.Log(string(receivedData[i].Json))
		if cont != true {
			t.Errorf("receivedData Error '%s' does not contain  %s", string(receivedData[i].Json), string(expected))
			break
		}
	}

}

func TestObserveKV(t *testing.T) {

	doneChanWrite := make(chan int)
	//doneChanRead := make(chan int)
	//now := time.Now().UnixNano() / int64(time.Millisecond)
	startAt := 0
	numRecords := 50

	receivedData := []JsonObserveResponse{}

	go func() {
		dataChan, err := kvc.Observe(dsID)
		if err != nil {
			t.Errorf("Observing %s failed expected err to be nil got %s", dsID, err.Error())
		}

		for data := range dataChan {
			receivedData = append(receivedData, data)
			t.Log("received:: ", string(data.Json), data.TimestampMS)
		}

	}()

	//Observe take a bit of time to register we miss some values if we dont wait before writing
	time.Sleep(time.Second / 2)

	go func() {
		for i := startAt; i < numRecords; i++ {
			err := kvc.Write(dsID, "observeTest"+strconv.Itoa(i), []byte("{\"value\":"+strconv.Itoa(i)+"}"))
			if err != nil {
				t.Errorf("WriteAt to %s failed expected err to be nil got %s", dsID, err.Error())
			}
			t.Log(string("written:: " + strconv.Itoa(i)))
			time.Sleep(time.Millisecond * 10)
		}
		doneChanWrite <- 1
	}()

	<-doneChanWrite

	if len(receivedData) != numRecords {
		t.Errorf("receivedData Error receivedData does not contain  %d records it only has %d", numRecords, len(receivedData))
		return
	}

	for i := startAt; i < numRecords; i++ {
		expected := []byte("{\"value\":" + strconv.Itoa(i) + "}")
		cont := s.Contains(string(receivedData[i].Json), string(expected))
		//t.Log("Data:: ", string(receivedData[i].Json), receivedData[i].Key)
		if cont != true && receivedData[i].Key == "observeTest"+strconv.Itoa(i) {
			t.Errorf("receivedData Error '%s' does not contain  %s", string(receivedData[i].Json), string(expected))
			break
		}
	}

}
