package libDatabox

import (
	"strconv"
)

/*
* href: https://[store-host-name]:[port]/[datasorceID]
* Timestamps in seconds
 */

func StoreTsLatest(href string) (string, error) {

	data, err := makeStoreRequest(href+"/ts/latest", "GET")
	if err != nil {
		return "", err
	}

	return data, nil

}

func StoreTsLast(href string, n int) (string, error) {
	data, err := makeStoreRequest(href+"/ts/last/"+i2s(n), "GET")
	if err != nil {
		return "", err
	}

	return data, nil
}

func StoreTsLastNsince(href string, n int, startTimestamp int) (string, error) {
	data, err := makeStoreRequest(href+"/ts/last/"+i2s(n)+"/since/"+i2s(startTimestamp), "GET")
	if err != nil {
		return "", err
	}

	return data, nil
}

func StoreTsLastNrange(href string, n int, startTimestamp int, endTimestamp int) (string, error) {
	data, err := makeStoreRequest(href+"/ts/last/"+i2s(n)+"/range/"+i2s(startTimestamp)+"/"+i2s(endTimestamp), "GET")
	if err != nil {
		return "", err
	}

	return data, nil
}

func StoreTsSince(href string, startTimestamp int) (string, error) {
	data, err := makeStoreRequest(href+"/ts/since/"+i2s(startTimestamp), "GET")
	if err != nil {
		return "", err
	}

	return data, nil
}

func StoreTsRange(href string, startTimestamp int, endTimestamp int) (string, error) {

	data, err := makeStoreRequest(href+"/ts/range/"+i2s(startTimestamp)+"/"+i2s(endTimestamp), "GET")
	if err != nil {
		return "", err
	}

	return data, nil
}

func StoreTsWrite(href string, data string) (string, error) {

	_, err := makeStoreRequestPOST(href+"/ts", data)
	if err != nil {
		return "", err
	}

	return "", nil
}

func i2s(i int) string {
	return strconv.Itoa(i)
}
