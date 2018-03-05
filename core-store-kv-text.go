package libDatabox

import (
	"errors"
	"strings"

	zest "github.com/toshbrown/goZestClient"
)

type TextKeyValue_0_3_0 interface {
	// Write text value
	Write(dataSourceID string, key string, payload string) error
	// Read text values. Returns a string containing the text written to the key.
	Read(dataSourceID string, key string) (string, error)
	// Get notifications of updated values for a key. Returns a channel that receives []bytes containing a JSON string when a new value is added.
	ObserveKey(dataSourceID string, key string) (<-chan string, error)
	// Get notifications of updated values for any key. Returns a channel that receives []bytes containing a JSON string when a new value is added.
	Observe(dataSourceID string) (<-chan string, error)
	// RegisterDatasource make a new data source for available to the rest of datbox. This can only be used on stores that you have requested in your manifest.
	RegisterDatasource(metadata DataSourceMetadata) error
}

type textKeyValueClient struct {
	zestClient         zest.ZestClient
	zestEndpoint       string
	zestDealerEndpoint string
}

// NewTextKeyValueClient returns a new TextKeyValue_0_3_0 to enable reading and writing of string data key value to the store
// reqEndpoint is provided in the DATABOX_ZMQ_ENDPOINT environment varable to databox apps and drivers.
func NewTextKeyValueClient(reqEndpoint string, enableLogging bool) (TextKeyValue_0_3_0, error) {

	kvc := textKeyValueClient{}
	kvc.zestEndpoint = reqEndpoint
	kvc.zestDealerEndpoint = strings.Replace(reqEndpoint, ":5555", ":5556", 1)
	thisKVC, err := zest.New(kvc.zestEndpoint, kvc.zestDealerEndpoint, getServerKey(), enableLogging)
	kvc.zestClient = thisKVC
	return kvc, err
}

func (kvc textKeyValueClient) Write(dataSourceID string, key string, payload string) error {
	path := "/kv/" + dataSourceID + "/" + key

	token, err := requestToken(kvc.zestEndpoint+path, "POST")
	if err != nil {
		return err
	}

	err = kvc.zestClient.Post(token, path, []byte(payload), "TEXT")
	if err != nil {
		invalidateCache(kvc.zestEndpoint+path, "POST")
		return errors.New("Error writing data: " + err.Error())
	}

	return nil
}

func (kvc textKeyValueClient) Read(dataSourceID string, key string) (string, error) {
	path := "/kv/" + dataSourceID + "/" + key

	token, err := requestToken(kvc.zestEndpoint+path, "POST")
	if err != nil {
		return "", err
	}

	data, getErr := kvc.zestClient.Get(token, path, "TEXT")
	if getErr != nil {
		invalidateCache(kvc.zestEndpoint+path, "POST")
		return "", errors.New("Error reading data: " + getErr.Error())
	}

	return string(data), nil
}

func (kvc textKeyValueClient) Observe(dataSourceID string) (<-chan string, error) {
	path := "/kv/" + dataSourceID

	token, err := requestToken(kvc.zestEndpoint+path, "GET")
	if err != nil {
		return nil, err
	}

	payloadChan, getErr := kvc.zestClient.Observe(token, path, "JSON", 0)
	if getErr != nil {
		invalidateCache(kvc.zestEndpoint+path, "GET")
		return nil, errors.New("Error observing: " + getErr.Error())
	}

	stringChan := make(chan string)
	go func(byteChan <-chan []byte, outputChan chan string) {
		for {
			outputChan <- string(<-byteChan)
		}
	}(payloadChan, stringChan)
	return stringChan, nil
}

func (kvc textKeyValueClient) ObserveKey(dataSourceID string, key string) (<-chan string, error) {
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

	stringChan := make(chan string)
	go func(byteChan <-chan []byte, outputChan chan string) {
		for {
			outputChan <- string(<-byteChan)
		}
	}(payloadChan, stringChan)
	return stringChan, nil
}

func (kvc textKeyValueClient) RegisterDatasource(metadata DataSourceMetadata) error {

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
