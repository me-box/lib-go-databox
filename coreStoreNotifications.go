package libDatabox

import (
	"bytes"
	"encoding/json"
	"errors"
	"strings"

	"github.com/google/uuid"
	zest "github.com/me-box/goZestClient"
)

//Func the databox function call, drivers can regiter functions with the Register method. Apps can request access to these in their manifests and call them using the call method Call.
type Func struct {
	csc                   *CoreStoreClient
	funcRequestChan       chan FuncRequest
	registeredFuncHandler map[string]FuncHandler
}

//FuncStatus is an int representing the status of a returned function
type FuncStatus int

const FuncStatusOK FuncStatus = 0
const FuncStatusFailedToGetToken FuncStatus = 97
const FuncStatusInvalidPayload FuncStatus = 98
const FuncStatusError FuncStatus = 99

// FuncResponse describes the response that must be returned by a FuncHandler
type FuncResponse struct {
	Status   FuncStatus
	Response []byte
}

// FuncRequest holds the datareturned from a function call
type FuncRequest struct {
	TimestampMS int64
	Path        string
	ContentType StoreContentType
	Payload     []byte
}

// FuncHandler s are executed when a request is received on the request endpoint of registered functions
type FuncHandler = func(contentType StoreContentType, payload []byte) ([]byte, error)

func newFunc(csc *CoreStoreClient) *Func {
	return &Func{
		csc:                   csc,
		funcRequestChan:       nil,
		registeredFuncHandler: make(map[string]FuncHandler),
	}
}

// Register is used to advertise available functions to other components
// it adds an entry into your stores Hypercat catalogue and starts listening
// for requests. Received requests are routed to the FuncHandler registed
// for the function.
func (f *Func) Register(vendor string, functionName string, contentType StoreContentType, handler FuncHandler) error {

	if _, ok := f.registeredFuncHandler[functionName]; ok {
		//we have already registered this function
		return errors.New("Unable to register function, already exists ")
	}

	//register function
	metadata := DataSourceMetadata{
		ContentType:    contentType,
		Vendor:         vendor,
		DataSourceType: vendor + ":func:" + functionName,
		DataSourceID:   functionName,
		StoreType:      StoreTypeFunc,
		Description:    "A function",
	}

	err := f.csc.RegisterDatasource(metadata)
	if err != nil {
		return errors.New("Unable to register function. " + err.Error())
	}

	//register the FuncHandler
	f.registeredFuncHandler[functionName] = handler

	//create FuncRequestChan by observing /notification/request/functionName/*
	//start go routine to process events, if we have not started one already.
	if f.funcRequestChan == nil {
		rawRequestChan, err := f.csc.observe("/notification/request/*", ContentTypeJSON, zest.ObserveModeNotification)
		if err != nil {
			Err("Could not observe /notification/request/* you will not receive any requests")
		}
		Debug("[Notifications] Setting up Observe on /notification/request/*")
		go f.parseRawFuncRequest(rawRequestChan)
	}
	return nil
}

// Call is used by clients to invoke functions by name. The
// result of the function call is returned via the FuncResponse chan
// only one result will be retuned then the channel will be closed.
func (f Func) Call(functionName string, payload []byte, contentType StoreContentType) (<-chan FuncResponse, error) {

	responseChan := make(chan FuncResponse)
	go f.parseRawFuncResponse(functionName, payload, contentType, responseChan)
	return responseChan, nil

}

func (f *Func) parseRawFuncRequest(rawRequest <-chan ObserveResponse) {
	Debug("[Notifications] Waiting to process function requests")

	for obsResponse := range rawRequest {
		Debug("[Notifications] Got a function request " + string(obsResponse.Data))
		parts := bytes.SplitN(obsResponse.Data, []byte(" "), 5)
		//timestamp, _ := strconv.ParseInt(string(parts[0]), 10, 64)
		responsePath := string(parts[2])
		splitPath := strings.Split(responsePath, "/")
		functionName := splitPath[3]
		ct := string(parts[3])
		payload := []byte{}
		if len(parts) >= 5 {
			payload = parts[4]
		}

		var contentType StoreContentType
		if _, ok := f.registeredFuncHandler[functionName]; ok {

			switch ct {
			case "json":
				contentType = ContentTypeJSON
				break
			case "binary":
				contentType = ContentTypeBINARY
				break
			case "text":
				contentType = ContentTypeTEXT
				break
			default:
				Err("Unknown content type (" + ct + ") in raw function request, not calling function")
				//dont process this data
				continue
			}

			//we have a registered function call it ;-)
			Debug("[Notifications] Calling registered function " + functionName)
			responseData, funcErr := f.registeredFuncHandler[functionName](contentType, payload)

			//Send response to caller
			var resp FuncResponse
			if funcErr == nil {
				resp.Status = FuncStatusOK
				resp.Response = responseData
			} else {
				resp.Status = FuncStatusError
				resp.Response = []byte(funcErr.Error())
			}
			respJson, _ := json.Marshal(resp)
			Debug("[Notifications] Sending response to caller on " + responsePath + " data: " + string(respJson))
			err := f.csc.write(responsePath, respJson, contentType)
			if err != nil {
				Err("Writing request to " + responsePath)
			}
		} else {
			//return an error code
			Err("Unknown/Unregistered function " + functionName)
		}
	}

}

func (f *Func) parseRawFuncResponse(functionName string, payload []byte, contentType StoreContentType, responseChan chan FuncResponse) {

	jobID := uuid.New().String()

	//set up a channel to receive the result
	NotifyResponseChan, err := f.csc.notify("/notification/response/"+functionName+"/"+jobID, contentType)
	if err != nil {
		responseChan <- FuncResponse{
			Status:   FuncStatusError,
			Response: []byte(`[Error] failed setup notification functionName for /notification/response/` + functionName + `/` + jobID + `. ` + err.Error()),
		}
		return
	}
	Debug("[Notifications] Setting up notify on /notification/response/" + functionName + "/" + jobID)

	//call the function
	Debug("[Notifications] Calling /notification/request/" + functionName + "/" + jobID + " with payload: " + string(payload))
	err = f.csc.write("/notification/request/"+functionName+"/"+jobID, payload, contentType)
	if err != nil {
		responseChan <- FuncResponse{
			Status:   FuncStatusError,
			Response: []byte(`[Error] failed to call to ` + functionName + " " + err.Error()),
		}
		return
	}

	//block and await the response
	response := <-NotifyResponseChan
	Debug("response.Data" + string(response.Data))

	var funcResp FuncResponse
	err = json.Unmarshal(response.Data, &funcResp)
	if err != nil {
		responseChan <- FuncResponse{
			Status:   FuncStatusError,
			Response: []byte(`[Error] failed to decode response from ` + functionName + " " + err.Error()),
		}
		return
	}

	Debug("funcResp.Response " + string(funcResp.Response))

	//send the result to the caller
	responseChan <- funcResp

	return
}
