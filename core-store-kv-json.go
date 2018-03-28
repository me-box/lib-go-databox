package libDatabox

import (
	"bytes"
	"encoding/json"
	"errors"
	"strconv"
	"strings"

	zest "github.com/toshbrown/goZestClient"
)

type JSONKeyValue_0_3_0 interface {
	// Write JSON value
	Write(dataSourceID string, key string, payload []byte) error
	// Read JSON values. Returns a []bytes containing a JSON string.
	Read(dataSourceID string, key string) ([]byte, error)
	//ListKeys returns an array of key registed under the dataSourceID
	ListKeys(dataSourceID string) ([]string, error)
	// Get notifications of updated values for a key. Returns a channel that receives JsonObserveResponse containing a JSON string when a new value is added.
	ObserveKey(dataSourceID string, key string) (<-chan JsonObserveResponse, error)
	// Get notifications of updated values for any key. Returns a channel that receives JsonObserveResponse containing a JSON string when a new value is added.
	Observe(dataSourceID string) (<-chan JsonObserveResponse, error)
	// RegisterDatasource make a new data source for available to the rest of datbox. This can only be used on stores that you have requested in your manifest.
	RegisterDatasource(metadata DataSourceMetadata) error
}

type jsonKeyValueClient struct {
	zestClient         zest.ZestClient
	zestEndpoint       string
	zestDealerEndpoint string
}

// NewJSONKeyValueClient returns a new NewJSONKeyValueClient to enable reading and writing of JSON data key value to the store
// reqEndpoint is provided in the DATABOX_ZMQ_ENDPOINT environment varable to databox apps and drivers.
func NewJSONKeyValueClient(reqEndpoint string, enableLogging bool) (JSONKeyValue_0_3_0, error) {

	kvc := jsonKeyValueClient{}
	kvc.zestEndpoint = reqEndpoint
	kvc.zestDealerEndpoint = strings.Replace(reqEndpoint, ":5555", ":5556", 1)
	thisKVC, err := zest.New(kvc.zestEndpoint, kvc.zestDealerEndpoint, getServerKey(), enableLogging)
	kvc.zestClient = thisKVC
	return kvc, err
}

func (kvc jsonKeyValueClient) Write(dataSourceID string, key string, payload []byte) error {
	path := "/kv/" + dataSourceID + "/" + key

	token, err := requestToken(kvc.zestEndpoint+path, "POST")
	if err != nil {
		return err
	}

	err = kvc.zestClient.Post(token, path, payload, "JSON")
	if err != nil {
		invalidateCache(kvc.zestEndpoint+path, "POST")
		return errors.New("Error writing data: " + err.Error())
	}

	return nil
}

func (kvc jsonKeyValueClient) Read(dataSourceID string, key string) ([]byte, error) {
	path := "/kv/" + dataSourceID + "/" + key

	token, err := requestToken(kvc.zestEndpoint+path, "GET")
	if err != nil {
		return []byte(""), err
	}

	data, getErr := kvc.zestClient.Get(token, path, "JSON")
	if getErr != nil {
		invalidateCache(kvc.zestEndpoint+path, "POST")
		return []byte(""), errors.New("Error reading data: " + getErr.Error())
	}

	return data, nil
}

func (kvc jsonKeyValueClient) ListKeys(dataSourceID string) ([]string, error) {
	path := "/kv/" + dataSourceID + "/keys"

	token, err := requestToken(kvc.zestEndpoint+path, "GET")
	if err != nil {
		return []string{}, err
	}

	data, getErr := kvc.zestClient.Get(token, path, "JSON")
	if getErr != nil {
		invalidateCache(kvc.zestEndpoint+path, "POST")
		return []string{}, errors.New("Error reading data: " + getErr.Error())
	}

	var keysArray []string

	err = json.Unmarshal(data, &keysArray)
	if err != nil {
		return []string{}, errors.New("Error decoding data: " + err.Error())
	}
	return keysArray, nil
}

func (kvc jsonKeyValueClient) ObserveKey(dataSourceID string, key string) (<-chan JsonObserveResponse, error) {
	path := "/kv/" + dataSourceID + "/" + key

	token, err := requestToken(kvc.zestEndpoint+path, "GET")
	if err != nil {
		return nil, err
	}

	payloadChan, getErr := kvc.zestClient.Observe(token, path, "JSON", 0)
	if getErr != nil {
		invalidateCache(kvc.zestEndpoint+path, "GET")
		return nil, errors.New("Error observing: " + getErr.Error())
	}

	objectChan := make(chan JsonObserveResponse)

	go func() {
		for data := range payloadChan {

			ts, dsid, key, payload := parseRawObserveResponse(data)
			resp := JsonObserveResponse{
				TimestampMS:  ts,
				DataSourceID: dsid,
				Key:          key,
				Json:         payload,
			}

			objectChan <- resp
		}

		//if we get here then payloadChan has been closed so close objectChan
		close(objectChan)
	}()

	return objectChan, nil
}

func (kvc jsonKeyValueClient) Observe(dataSourceID string) (<-chan JsonObserveResponse, error) {
	path := "/kv/" + dataSourceID + "/*"

	token, err := requestToken(kvc.zestEndpoint+path, "GET")
	if err != nil {
		return nil, err
	}

	payloadChan, getErr := kvc.zestClient.Observe(token, path, "JSON", 0)
	if getErr != nil {
		invalidateCache(kvc.zestEndpoint+path, "GET")
		return nil, errors.New("Error observing: " + getErr.Error())
	}

	objectChan := make(chan JsonObserveResponse)

	go func() {
		for data := range payloadChan {

			ts, dsid, key, payload := parseRawObserveResponse(data)
			resp := JsonObserveResponse{
				TimestampMS:  ts,
				DataSourceID: dsid,
				Key:          key,
				Json:         payload,
			}

			objectChan <- resp
		}

		//if we get here then payloadChan has been closed so close objectChan
		close(objectChan)
	}()

	return objectChan, nil
}

func parseRawObserveResponse(data []byte) (int64, string, string, []byte) {

	parts := bytes.SplitN(data, []byte(" "), 4)

	_timestamp, _ := strconv.ParseInt(string(parts[0]), 10, 64)

	parts2 := bytes.Split(parts[1], []byte("/"))

	_dataSourceID := string(parts2[2])

	_key := ""
	if len(parts2) > 3 {
		_key = string(parts2[3])
	}

	_data := parts[3]

	return _timestamp, _dataSourceID, _key, _data
}

func (kvc jsonKeyValueClient) RegisterDatasource(metadata DataSourceMetadata) error {

	path := "/cat"

	token, err := requestToken(kvc.zestEndpoint+path, "POST")
	if err != nil {
		return errors.New("Error getting token: " + err.Error())
	}
	hypercatJSON, err := dataSourceMetadataToHypercat(metadata, kvc.zestEndpoint+"/kv/")

	writeErr := kvc.zestClient.Post(token, path, hypercatJSON, "JSON")
	if writeErr != nil {
		invalidateCache(kvc.zestEndpoint+path, "POST")
		return errors.New("Error writing: " + writeErr.Error())
	}

	return nil
}
