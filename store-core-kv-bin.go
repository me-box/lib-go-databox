package libDatabox

import (
	"errors"
	"io/ioutil"
	"strings"

	zest "github.com/toshbrown/goZestClient"
)

type BinaryKeyValue_0_2_0 interface {
	// Write text value
	Write(dataSourceID string, payload []byte) error
	// Read text values.
	Read(dataSourceID string) ([]byte, error)
	// Read JSON values.
	// Get notifications of updated values
	Observe(dataSourceID string) (<-chan []byte, error)
	// Get notifications of updated values
	RegisterDatasource(metadata DataSourceMetadata) error
}

type binaryKeyValueClient struct {
	zestClient         zest.ZestClient
	zestEndpoint       string
	zestDealerEndpoint string
}

// NewBinaryKeyValueClient returns a new NewBinaryKeyValueClient to enable reading and writing of binary data key value to the store
// reqEndpoint is provided in the DATABOX_ZMQ_ENDPOINT environment varable to databox apps and drivers.
func NewBinaryKeyValueClient(reqEndpoint string, enableLogging bool) (BinaryKeyValue_0_2_0, error) {

	serverKey, err := ioutil.ReadFile("/run/secrets/ZMQ_PUBLIC_KEY")
	if err != nil {
		return binaryKeyValueClient{}, err
	}

	kvc := binaryKeyValueClient{}
	kvc.zestEndpoint = reqEndpoint
	kvc.zestDealerEndpoint = strings.Replace(reqEndpoint, ":5555", ":5556", 1)
	kvc.zestClient, err = zest.New(kvc.zestEndpoint, kvc.zestDealerEndpoint, string(serverKey), enableLogging)
	return kvc, err
}

func (kvc binaryKeyValueClient) Write(dataSourceID string, payload []byte) error {
	path := "/kv/" + dataSourceID

	token, err := requestToken(kvc.zestEndpoint+path, "POST")
	if err != nil {
		return err
	}

	err = kvc.zestClient.Post(token, path, payload, "BINARY")
	if err != nil {
		invalidateCache(kvc.zestEndpoint+path, "POST")
		return errors.New("Error writing data: " + err.Error())
	}

	return nil
}

func (kvc binaryKeyValueClient) Read(dataSourceID string) ([]byte, error) {
	path := "/kv/" + dataSourceID

	token, err := requestToken(kvc.zestEndpoint+path, "POST")
	if err != nil {
		return []byte(""), err
	}

	data, getErr := kvc.zestClient.Get(token, path, "BINARY")
	if getErr != nil {
		invalidateCache(kvc.zestEndpoint+path, "POST")
		return []byte(""), errors.New("Error reading data: " + err.Error())
	}

	return data, nil
}

func (kvc binaryKeyValueClient) Observe(dataSourceID string) (<-chan []byte, error) {
	path := "/kv/" + dataSourceID

	token, err := requestToken(kvc.zestEndpoint+path, "GET")
	if err != nil {
		return nil, err
	}

	payloadChan, getErr := kvc.zestClient.Observe(token, path, "BINARY")
	if getErr != nil {
		invalidateCache(kvc.zestEndpoint+path, "GET")
		return nil, errors.New("Error observing: " + err.Error())
	}

	return payloadChan, nil
}

func (kvc binaryKeyValueClient) RegisterDatasource(metadata DataSourceMetadata) error {

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
