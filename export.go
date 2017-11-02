package libDatabox

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

var exportServiceURL = os.Getenv("DATABOX_EXPORT_SERVICE_ENDPOINT")

// ExportLongpoll exports data to external service (payload must be an escaped json string)
// permissions must be requested in the app manifest (drivers dont need to use the export service)
func ExportLongpoll(destination string, payload string) (string, error) {

	//TODO payload must be an escaped json string detect it it is not and error or escape it!!

	var jsonStr = "{\"id\":\"\",\"uri\":\"" + destination + "\",\"data\":" + payload + "}"

	fmt.Println("Sending ", jsonStr)

	res, err := makeStoreRequestPOST(exportServiceURL+"/lp/export", jsonStr)

	return res, err
}

func makeStoreRequestPOST(href string, data string) (string, error) {

	method := "POST"
	token, err := checkTokenCache(href, method)
	if err != nil {
		return "", err
	}

	//perform store request with token
	req, err := http.NewRequest(method, href, bytes.NewBufferString(data))
	req.Header.Set("X-Api-Key", token)
	req.Header.Set("Content-Type", "application/json")
	req.Close = true
	resp, err := databoxClient.Do(req)
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
