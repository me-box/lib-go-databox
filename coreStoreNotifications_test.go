package libDatabox

import (
	"bytes"
	"encoding/json"
	s "strings"
	"testing"
)

func TestFuncRegistration(t *testing.T) {

	//Test function registration
	testFunc := func(contentType StoreContentType, payload []byte) []byte {

		return []byte("Testingtesting132")
	}

	err := StoreClient2.FUNC.Register("databox", "testFunc"+dsID, ContentTypeTEXT, testFunc)
	if err != nil {
		t.Errorf("FUNC.Register failed expected err to be nil got %s", err.Error())
	}

	hypercatRoot, getErr := StoreClient2.GetStoreDataSourceCatalogue(StoreURL)
	if getErr != nil {
		t.Errorf("GetDatasourceCatalogue failed expected err to be nil got %s", getErr.Error())
	}
	catByteArray, _ := json.Marshal(hypercatRoot)

	cont := s.Contains(string(catByteArray), StoreURL+"/notification/request/testFunc"+dsID)
	if cont != true {
		t.Errorf("TestFuncRegistration Error '%s' does not contain  %s", string(catByteArray), StoreURL+"/notification/request/testFunc"+dsID)
	}

}

func TestFuncCall(t *testing.T) {

	//Test function registration
	testFunc := func(contentType StoreContentType, payload []byte) []byte {

		return []byte("Testingtesting132" + dsID)
	}

	err := StoreClient2.FUNC.Register("databox", "TestFunc", ContentTypeJSON, testFunc)
	if err != nil {
		t.Errorf("FUNC.Register failed expected err to be nil got %s", err.Error())
	}

	//call the function
	funcResponseChan, err := StoreClient2.FUNC.Call("TestFunc", []byte{}, ContentTypeJSON)
	if err != nil {
		t.Errorf("TestFunc Call failed expected err to be nil got %s", err.Error())
	}

	response := <-funcResponseChan

	if response.Status != FuncStatusOK {
		t.Errorf("TestFunc Call failed expected status to be  to be FuncStatusOK got %d with the message %s", response.Status, response.Response)
	}

	if !bytes.Equal(response.Response, []byte("Testingtesting132"+dsID)) {
		t.Errorf("TestFunc Call failed expected Response to be Testingtesting132"+dsID+" got %s ", response.Response)
	}

}

func TestFuncCallWithPayload(t *testing.T) {

	//Test function registration
	testFunc := func(contentType StoreContentType, payload []byte) []byte {

		return payload
	}

	err := StoreClient2.FUNC.Register("databox", "TestFuncWithPayload", ContentTypeJSON, testFunc)
	if err != nil {
		t.Errorf("FUNC.Register failed expected err to be nil got %s", err.Error())
	}

	//call the function
	funcResponseChan, err := StoreClient2.FUNC.Call("TestFuncWithPayload", []byte("This is a test"), ContentTypeJSON)
	if err != nil {
		t.Errorf("TestFunc Call failed expected err to be nil got %s", err.Error())
	}

	response := <-funcResponseChan

	if response.Status != FuncStatusOK {
		t.Errorf("TestFunc Call failed expected status to be  to be FuncStatusOK got %d with the message %s", response.Status, response.Response)
	}

	if !bytes.Equal(response.Response, []byte("This is a test")) {
		t.Errorf("TestFunc Call failed expected Response to be 'This is a test' got %s ", response.Response)
	}

}

func TestFuncCallWithPayloadManyCalls(t *testing.T) {

	//Test function registration
	testFunc := func(contentType StoreContentType, payload []byte) []byte {

		return payload
	}

	err := StoreClient2.FUNC.Register("databox", "TestFuncCallWithPayloadManyCalls", ContentTypeJSON, testFunc)
	if err != nil {
		t.Errorf("FUNC.Register failed expected err to be nil got %s", err.Error())
	}

	//call the function
	for i := 0; i < 10; i++ {
		funcResponseChan, err := StoreClient2.FUNC.Call("TestFuncCallWithPayloadManyCalls", []byte("This is a test"), ContentTypeJSON)
		if err != nil {
			t.Errorf("TestFunc Call failed expected err to be nil got %s", err.Error())
		}

		response := <-funcResponseChan

		if response.Status != FuncStatusOK {
			t.Errorf("TestFunc Call failed expected status to be  to be FuncStatusOK got %d with the message %s", response.Status, response.Response)
		}

		if !bytes.Equal(response.Response, []byte("This is a test")) {
			t.Errorf("TestFunc Call failed Response to be 'this is a test' got %s ", response.Response)
		}
	}

}
