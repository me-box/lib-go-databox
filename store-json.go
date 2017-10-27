package libDatabox

import (
	"strconv"
)

func StoreJSONGetlatest(href string) (string, error) {

	data, err := makeStoreRequest(href+"/ts/latest", "GET")
	if err != nil {
		return "", err
	}

	return data, nil

}

func StoreJSONGetrange(href string, startTimestamp int64, endTimestamp int64) (string, error) {

	params :=  map[string]string {
		"startTimestamp": strconv.FormatInt(startTimestamp, 10),
		"endTimestamp": strconv.FormatInt(endTimestamp, 10),
	}
	data, err := makeStoreRequestForm(href+"/ts/range", "GET",  params)
	if err != nil {
		return "", err
	}

	return data, nil

}

func StoreJSONGetsince(href string, startTimestamp int64) (string, error) {

	params :=  map[string]string {
		"startTimestamp": strconv.FormatInt(startTimestamp, 10),
	}
	data, err := makeStoreRequestForm(href+"/ts/since", "GET",  params)
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
