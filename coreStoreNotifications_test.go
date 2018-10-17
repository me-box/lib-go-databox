package libDatabox

import (
	"bytes"
	"encoding/json"
	"errors"
	"strconv"
	s "strings"
	"testing"
)

func TestFuncRegistration(t *testing.T) {

	//Test function registration
	testFunc := func(contentType StoreContentType, payload []byte) ([]byte, error) {

		return []byte("Testingtesting132"), nil
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
	testFunc := func(contentType StoreContentType, payload []byte) ([]byte, error) {

		return []byte("Testingtesting132" + dsID), nil
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
	testFunc := func(contentType StoreContentType, payload []byte) ([]byte, error) {

		return payload, nil
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
	testFunc := func(contentType StoreContentType, payload []byte) ([]byte, error) {

		return payload, nil
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

func TestFuncCallWithError(t *testing.T) {

	//Test function registration
	testFunc := func(contentType StoreContentType, payload []byte) ([]byte, error) {
		return payload, errors.New("Test Error")
	}

	err := StoreClient2.FUNC.Register("databox", "TestFuncCallWithError", ContentTypeJSON, testFunc)
	if err != nil {
		t.Errorf("FUNC.Register failed expected err to be nil got %s", err.Error())
	}

	funcResponseChan, err := StoreClient2.FUNC.Call("TestFuncCallWithError", []byte("This is a test"), ContentTypeJSON)
	if err != nil {
		t.Errorf("TestFuncCallWithError Call failed expected err to be nil got %s", err.Error())
	}

	response := <-funcResponseChan

	if response.Status != FuncStatusError {
		t.Errorf("TestFuncCallWithError Call failed expected status to be  to be FuncStatusError got %d with the message %s", response.Status, response.Response)
	}

	if !bytes.Equal(response.Response, []byte("Test Error")) {
		t.Errorf("TestFuncCallWithError Call failed Response to be 'Test Error' got %s ", response.Response)
	}

}

func TestFuncCallWith2FuncsAndOneClient(t *testing.T) {

	//Test function registration
	testFunc := func(contentType StoreContentType, payload []byte) ([]byte, error) {

		return payload, nil
	}

	err := StoreClient2.FUNC.Register("databox", "TestFuncCallWith2FuncsAndOneClient1", ContentTypeJSON, testFunc)
	if err != nil {
		t.Errorf("FUNC.Register failed expected err to be nil got %s", err.Error())
	}

	err = StoreClient2.FUNC.Register("databox", "TestFuncCallWith2FuncsAndOneClient2", ContentTypeJSON, testFunc)
	if err != nil {
		t.Errorf("FUNC.Register failed expected err to be nil got %s", err.Error())
	}

	//call the function
	for i := 0; i < 30; i++ {
		funcResponseChan1, err := StoreClient2.FUNC.Call("TestFuncCallWith2FuncsAndOneClient1", []byte("This is test "+strconv.Itoa(i)), ContentTypeJSON)
		if err != nil {
			t.Errorf("TestFunc Call failed expected err to be nil got %s", err.Error())
		}

		funcResponseChan2, err := StoreClient2.FUNC.Call("TestFuncCallWith2FuncsAndOneClient2", []byte("This is test "+strconv.Itoa(i)), ContentTypeJSON)
		if err != nil {
			t.Errorf("TestFunc Call failed expected err to be nil got %s", err.Error())
		}

		select {
		case response1 := <-funcResponseChan1:
			if response1.Status != FuncStatusOK {
				t.Errorf("TestFuncCallWith2FuncsAndOneClient1 Call failed expected status to be  to be FuncStatusOK got %d with the message %s", response1.Status, response1.Response)
			}

			if !bytes.Equal(response1.Response, []byte("This is test "+strconv.Itoa(i))) {
				t.Errorf("TestFuncCallWith2FuncsAndOneClient1 Call failed Response to be 'This is test %d' got '%s' ", i, response1.Response)
			}
		case response2 := <-funcResponseChan2:
			if response2.Status != FuncStatusOK {
				t.Errorf("TestFuncCallWith2FuncsAndOneClient2 Call failed expected status to be  to be FuncStatusOK got %d with the message %s", response2.Status, response2.Response)
			}

			if !bytes.Equal(response2.Response, []byte("This is test "+strconv.Itoa(i))) {
				t.Errorf("TestFuncCallWith2FuncsAndOneClient2 Call failed Response to be 'This is test %d' got '%s' ", i, response2.Response)
			}

		}

	}

}
