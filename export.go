package libDatabox

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

//TODO Make exportServiceClient work !!!

var exportServiceURL = os.Getenv("DATABOX_EXPORT_SERVICE_ENDPOINT")

const exportServiceName = "export-service"

type Export struct {
	arb               *ArbiterClient
	databoxHTTPClient *http.Client
}

func newExport(arb *ArbiterClient) *Export {
	return &Export{
		arb:               arb,
		databoxHTTPClient: NewDataboxHTTPsAPI(),
	}
}

// Longpoll exports data to external service (payload must be an escaped json string)
// permissions must be requested in the app manifest (drivers dont need to use the export service)
func (e Export) Longpoll(destination string, payload string) (string, error) {

	//TODO payload must be an escaped json string detect it it is not and error or escape it!!

	var jsonStr = "{\"id\":\"\",\"uri\":\"" + destination + "\",\"data\":" + payload + "}"

	fmt.Println("Sending ", jsonStr)

	res, err := e.makeStoreRequestPOST(exportServiceURL+"/lp/export", destination, jsonStr)

	return res, err
}

func (e Export) makeStoreRequestPOST(href string, destination string, data string) (string, error) {

	method := "POST"
	caveat := "destination = [\"" + destination + "\"]"
	token, err := e.arb.RequestToken(exportServiceName, method, []string{caveat})
	if err != nil {
		return "", err
	}

	//perform store request with token
	req, err := http.NewRequest(method, href, bytes.NewBufferString(data))
	req.Header.Set("X-Api-Key", string(token))
	req.Header.Set("Content-Type", "application/json")
	req.Close = true
	resp, err := e.databoxHTTPClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err1 := ioutil.ReadAll(resp.Body)
	if err1 != nil {
		return "", err1
	}

	return string(body[:]), nil
}
