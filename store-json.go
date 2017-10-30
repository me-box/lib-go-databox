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

func StoreJSONGetrange(href string, startTimestamp float64, endTimestamp float64) (string, error) {

	params :=  "{\"startTimestamp\": "+strconv.FormatFloat(startTimestamp, 'f', -1, 64)+",\"endTimestamp\": "+strconv.FormatFloat(endTimestamp, 'f', -1, 64)+"}"
	data, err := makeStoreRequestJson(href+"/ts/range", "GET",  params)
	if err != nil {
		return "", err
	}

	return data, nil

}

func StoreJSONGetsince(href string, startTimestamp float64) (string, error) {

	params :=  "{\"startTimestamp\": "+strconv.FormatFloat(startTimestamp, 'f', -1, 64)+"}"
	data, err := makeStoreRequestJson(href+"/ts/since", "GET",  params)
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
