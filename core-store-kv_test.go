package libDatabox

import (
	s "strings"
	"testing"
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
