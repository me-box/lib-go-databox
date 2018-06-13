package libDatabox

import (
	"strconv"
	s "strings"
	"testing"
	"time"
)

func TestKVTextWrite(t *testing.T) {
	err := StoreClient.KVTextWrite(dsID, "key1", []byte("{\"value\":3.1415}"))
	if err != nil {
		t.Errorf("Write to %s failed expected err to be nil got %s", dsID, err.Error())
	}
}

func TestKVTextRead(t *testing.T) {
	err := StoreClient.KVTextWrite(dsID, "key2", []byte("{\"value\":42}"))
	if err != nil {
		t.Errorf("Write to %s failed expected err to be nil got %s", dsID, err.Error())
	}

	result, err := StoreClient.KVTextRead(dsID, "key2")
	if err != nil {
		t.Errorf("Write to %s failed expected err to be nil got %s", dsID, err.Error())
	}

	expected := []byte("{\"value\":42}")
	cont := s.Contains(string(result), string(expected))
	if cont != true {
		t.Errorf("TestWriteLots failed expected %s but got %s", expected, result)
	}
}

func TestKVTextMutiKeys(t *testing.T) {

	err := StoreClient.KVTextWrite(dsID, "key1", []byte("{\"value\":\"some random string\"}"))
	if err != nil {
		t.Errorf("Write to %s failed expected err to be nil got %s", dsID, err.Error())
	}

	err = StoreClient.KVTextWrite(dsID, "key2", []byte("{\"value\":42}"))
	if err != nil {
		t.Errorf("Write to %s failed expected err to be nil got %s", dsID, err.Error())
	}

	result, err := StoreClient.KVTextRead(dsID, "key2")
	if err != nil {
		t.Errorf("Write to %s failed expected err to be nil got %s", dsID, err.Error())
	}

	expected := []byte("{\"value\":42}")
	cont := s.Contains(string(result), string(expected))
	if cont != true {
		t.Errorf("TestWriteLots failed expected %s but got %s", expected, result)
	}

	result, err = StoreClient.KVTextRead(dsID, "key1")
	if err != nil {
		t.Errorf("Write to %s failed expected err to be nil got %s", dsID, err.Error())
	}

	expected = []byte("{\"value\":\"some random string\"}")
	cont = s.Contains(string(result), string(expected))
	if cont != true {
		t.Errorf("TestWriteLots failed expected %s but got %s", expected, result)
	}
}

func TestListKeysKVText(t *testing.T) {
	_dsID := dsID + "TestListKeysKV"
	err := StoreClient.KVTextWrite(_dsID, "key1", []byte("{\"value\":\"some random string\"}"))
	if err != nil {
		t.Errorf("Write to %s failed expected err to be nil got %s", _dsID, err.Error())
	}
	err = StoreClient.KVTextWrite(_dsID, "key2", []byte("{\"value\":\"some random string\"}"))
	if err != nil {
		t.Errorf("Write to %s failed expected err to be nil got %s", _dsID, err.Error())
	}
	err = StoreClient.KVTextWrite(_dsID, "key3", []byte("{\"value\":\"some random string\"}"))
	if err != nil {
		t.Errorf("Write to %s failed expected err to be nil got %s", _dsID, err.Error())
	}

	keys, err := StoreClient.KVTextListKeys(_dsID)
	if err != nil {
		t.Errorf("ListKeys from %s failed expected err to be nil got %s", _dsID, err.Error())
	}

	if keys[0] != "key3" {
		t.Errorf("ListKeys error expected %s got %s", "key1", keys[0])
	}

	if keys[1] != "key2" {
		t.Errorf("ListKeys error expected %s got %s", "key2", keys[1])
	}

	if keys[2] != "key1" {
		t.Errorf("ListKeys error expected %s got %s", "key3", keys[2])
	}

}

func TestKVTextObserveKey(t *testing.T) {

	doneChanWrite := make(chan int)
	//doneChanRead := make(chan int)
	//now := time.Now().UnixNano() / int64(time.Millisecond)
	startAt := 0
	numRecords := 10

	receivedData := [][]byte{}

	go func() {
		dataChan, err := StoreClient.KVTextObserveKey(dsID, "observeTest")
		if err != nil {
			t.Errorf("Observing %s failed expected err to be nil got %s", dsID, err.Error())
		}

		for data := range dataChan {
			receivedData = append(receivedData, data)
			t.Log("received:: ", string(data))
		}

	}()

	//Observe take a bit of time to register we miss some values if we dont wait before writing
	time.Sleep(time.Second * 2)

	go func() {
		for i := startAt; i <= numRecords; i++ {
			err := StoreClient.KVTextWrite(dsID, "observeTest", []byte("{\"value\":"+strconv.Itoa(i)+"}"))
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
		expected := []byte("{\"value\":" + strconv.Itoa(i) + "}")
		cont := s.Contains(string(receivedData[i]), string(expected))
		t.Log(string(receivedData[i]))
		if cont != true {
			t.Errorf("receivedData Error '%s' does not contain  %s", string(receivedData[i]), string(expected))
			break
		}
	}

}

func TestObserveKVText(t *testing.T) {

	doneChanWrite := make(chan int)
	//doneChanRead := make(chan int)
	//now := time.Now().UnixNano() / int64(time.Millisecond)
	startAt := 0
	numRecords := 5

	receivedData := [][]byte{}

	go func() {
		dataChan, err := StoreClient.KVTextObserve(dsID)
		if err != nil {
			t.Errorf("Observing %s failed expected err to be nil got %s", dsID, err.Error())
		}

		for data := range dataChan {
			receivedData = append(receivedData, data)
			t.Log("received:: ", string(data))
		}

	}()

	//Observe take a bit of time to register we miss some values if we dont wait before writing
	time.Sleep(time.Second * 2)

	go func() {
		for i := startAt; i < numRecords; i++ {
			err := StoreClient.KVTextWrite(dsID, "observeTest"+strconv.Itoa(i), []byte("{\"value\":"+strconv.Itoa(i)+"}"))
			if err != nil {
				t.Errorf("WriteAt to %s failed expected err to be nil got %s", dsID, err.Error())
			}
			t.Log(string("written:: " + strconv.Itoa(i)))
			time.Sleep(time.Millisecond * 10)
		}
		// we miss some values if we dont wait before saying we are done!
		time.Sleep(time.Second * 3)
		doneChanWrite <- 1
	}()

	<-doneChanWrite

	if len(receivedData) != numRecords {
		t.Errorf("receivedData Error receivedData does not contain  %d records it only has %d", numRecords, len(receivedData))
		return
	}

	if len(receivedData) < numRecords {
		t.Errorf("receivedData Error:  receivedData should contain '%d' items but contains  %d", numRecords, len(receivedData))
	}
	for i := startAt; i < numRecords; i++ {
		expected := []byte("{\"value\":" + strconv.Itoa(i) + "}")
		cont := s.Contains(string(receivedData[i]), string(expected))
		//t.Log("Data:: ", string(receivedData[i].Json), receivedData[i].Key)
		if cont != true {
			t.Errorf("receivedData Error '%s' does not contain  %s", string(receivedData[i]), string(expected))
			break
		}
	}

}
