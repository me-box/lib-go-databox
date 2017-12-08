package libDatabox

import (
	"errors"
	"io/ioutil"
	"strings"

	zest "github.com/toshbrown/goZestClient"
)

type JSONKeyValue_0_2_0 interface {
	// Write JSON value
	Write(dataSourceID string, payload []byte) error
	// Read JSON values. Returns a []bytes containing a JSON string.
	Read(dataSourceID string) ([]byte, error)
	// Get notifications of updated values Returns a channel that receives []bytes containing a JSON string when a new value is added.
	Observe(dataSourceID string) (<-chan []byte, error)
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
func NewJSONKeyValueClient(reqEndpoint string, enableLogging bool) (JSONKeyValue_0_2_0, error) {

	serverKey, err := ioutil.ReadFile("/run/secrets/ZMQ_PUBLIC_KEY")
	if err != nil {
		return jsonKeyValueClient{}, err
	}

	kvc := jsonKeyValueClient{}
	kvc.zestEndpoint = reqEndpoint
	kvc.zestDealerEndpoint = strings.Replace(reqEndpoint, ":5555", ":5556", 1)
	thisKVC, err := zest.New(kvc.zestEndpoint, kvc.zestDealerEndpoint, string(serverKey), enableLogging)
	kvc.zestClient = thisKVC
	return kvc, err
}

func (kvc jsonKeyValueClient) Write(dataSourceID string, payload []byte) error {
	path := "/kv/" + dataSourceID

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

func (kvc jsonKeyValueClient) Read(dataSourceID string) ([]byte, error) {
	path := "/kv/" + dataSourceID

	token, err := requestToken(kvc.zestEndpoint+path, "POST")
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

func (kvc jsonKeyValueClient) Observe(dataSourceID string) (<-chan []byte, error) {
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

	return payloadChan, nil
}

func (kvc jsonKeyValueClient) RegisterDatasource(metadata DataSourceMetadata) error {

	path := "/cat"

	token, err := requestToken(kvc.zestEndpoint+path, "POST")
	if err != nil {
		return errors.New("Error getting token: " + err.Error())
	}
	hypercatJSON, err := dataSourceMetadataToHypercat(metadata)

	writeErr := kvc.zestClient.Post(token, path, hypercatJSON, "JSON")
	if writeErr != nil {
		invalidateCache(kvc.zestEndpoint+path, "POST")
		return errors.New("Error writing: " + writeErr.Error())
	}

	return nil
}
