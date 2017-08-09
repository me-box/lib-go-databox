package libDatabox

import "os"

var exportServiceURL = os.Getenv("DATABOX_EXPORT_SERVICE_ENDPOINT")

func ExportLongpoll(destination string, payload string) (string, error) {

	var jsonStr = "{\"id\":\"\",\"uri\":" + destination + ",\"data\":" + payload + "\"}"

	res, err := makeStoreRequestPOST(exportServiceURL+"/lp/export", jsonStr)

	return res, err
}

func StoreJSONGetlatest(href string) (string, error) {

	data, err := makeStoreRequest(href+"/ts/latest", "GET")
	if err != nil {
		return "", err
	}

	return data, nil

}

func StoreJSONWriteTS(href string, data string) error {

	_, err := makeStoreRequestPOST(href+"/ts", data)
	if err != nil {
		return err
	}

	return nil

}

func StoreJSONWriteKV(href string, data string) error {

	_, err := makeStoreRequestPOST(href+"/kv", data)
	if err != nil {
		return err
	}

	return nil

}

func StoreJSONreadKV(href string, key string) (string, error) {

	data, err := makeStoreRequest(href+"/"+key+"/kv", "GET")
	if err != nil {
		return "", err
	}

	return data, nil

}
