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

	err := StoreClient.FUNC.Register("databox", "testFunc"+dsID, ContentTypeTEXT, testFunc)
	if err != nil {
		t.Errorf("FUNC.Register failed expected err to be nil got %s", err.Error())
	}

	hypercatRoot, getErr := StoreClient.GetStoreDataSourceCatalogue(StoreURL)
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

	err := StoreClient.FUNC.Register("databox", "TestFuncCall"+dsID, ContentTypeTEXT, testFunc)
	if err != nil {
		t.Errorf("FUNC.Register failed expected err to be nil got %s", err.Error())
	}

	funcResponseChan, err := StoreClient.FUNC.Call("TestFuncCall"+dsID, []byte{}, ContentTypeTEXT)
	if err != nil {
		t.Errorf("TestFunc Call failed expected err to be nil got %s", err.Error())
	}

	response := <-funcResponseChan

	if response.Status != FuncStatusOK {
		t.Errorf("TestFunc Call failed expected status to be  to be FuncStatusOK got %d with the message %s", response.Status, response.Response)
	}

	if bytes.Equal(response.Response, []byte("Testingtesting132"+dsID)) {
		t.Errorf("TestFunc Call failed expected status to be  to be Testingtesting132"+dsID+" got %d ", response.Response)
	}

}
