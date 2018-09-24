package libDatabox

import (
	"bytes"
	"errors"
	"strings"

	zest "github.com/me-box/goZestClient"
)

type Func struct {
	csc                   *CoreStoreClient
	funcRequestChan       chan FuncRequest
	registeredFuncHandler map[string]FuncHandler
}

type FuncStatus int

const FuncStatusOK = 0
const FuncStatusFailedToGetToken = 97
const FuncStatusInvalidPayload = 98
const FuncStatusError = 99

type FuncResponse struct {
	Status   FuncStatus
	Response []byte
}

type FuncRequest struct {
	TimestampMS int64
	Path        string
	ContentType StoreContentType
	Payload     []byte
}

func newFunc(csc *CoreStoreClient) *Func {
	return &Func{
		csc:                   csc,
		funcRequestChan:       nil,
		registeredFuncHandler: make(map[string]FuncHandler),
	}
}

// FuncHandler s are executed when a request is received on the request endpoint of registered functions
type FuncHandler = func(contentType StoreContentType, payload []byte) []byte

// Register is used to advertise available functions to other components
// it adds an enters into your stores Hypercat catalogue.
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

	//create FuncRequestChan by observing /notification/request/*
	//start go routine to process events, if we have not started one already.
	if f.funcRequestChan == nil {
		rawRequestChan, err := f.csc.observe("/notification/request/*", ContentTypeJSON, zest.ObserveModeNotification)
		if err != nil {
			Err("Could not observe /notification/request/* you will not receive any requests")
		}
		Info("Setting up Observe on /notification/request/*")
		go f.parseRawFuncRequest(rawRequestChan)
	}
	return nil
}

// Call is used by clients to invoke functions by name
func (f Func) Call(functionName string, payload []byte, contentType StoreContentType) (<-chan FuncResponse, error) {

	responseChan := make(chan FuncResponse)
	go f.parseRawFuncResponse(functionName, payload, contentType, responseChan)
	return responseChan, nil

}

func (f *Func) parseRawFuncRequest(rawRequest <-chan ObserveResponse) {

	for obsResponse := range rawRequest {
		parts := bytes.SplitN(obsResponse.Data, []byte(" "), 4)
		//timestamp, _ := strconv.ParseInt(string(parts[0]), 10, 64)
		responsePath := string(parts[1])
		functionName := strings.Replace(string(parts[1]), "/notification/request/", "", 1)
		ct := string(parts[2])
		payload := parts[3]

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
				Err("Unknown content type in raw function request, not calling function")
				//dont process this data
				continue
			}

			//we have a registered function call it ;-)
			Info("Calling registered function " + functionName)
			response := f.registeredFuncHandler[functionName](contentType, payload)

			//Send response to caller
			Info("Sending response to caller on " + responsePath + " data: " + string(response))
			err := f.csc.write(responsePath, response, contentType)
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

	//set up a channel to receive the result
	NotifyResponseChan, err := f.csc.notify("/notification/response/"+functionName, contentType)
	if err != nil {
		responseChan <- FuncResponse{
			Status:   FuncStatusError,
			Response: []byte(`[Error] failed setup notification functionName for /notification/response/` + functionName + `. ` + err.Error()),
		}
		return
	}
	Info("Setting up notify on /notification/response/" + functionName)

	//call the function
	Info("Calling /notification/request/" + functionName + " with payload: " + string(payload))
	err = f.csc.write("/notification/request/"+functionName, payload, contentType)
	if err != nil {
		responseChan <- FuncResponse{
			Status:   FuncStatusInvalidPayload,
			Response: []byte(`[Error] failed to call to ` + functionName + " " + err.Error()),
		}
		return
	}

	//block and await the response
	response := <-NotifyResponseChan

	//send the result to the caller
	responseChan <- FuncResponse{
		Status:   FuncStatusOK,
		Response: response.Data,
	}
	return
}
