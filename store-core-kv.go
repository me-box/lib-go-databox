package libDatabox

import (
	"errors"
	"io/ioutil"
	"strings"

	zest "github.com/toshbrown/goZestClient"
)

type KeyValue_0_2_0 interface {
	// Write value.
	Write(payload string) error
	// Read values.
	Read() (string, error)
	// Get notifications of updated values
	Observe() (<-chan string, error)
	// registerDatasource is used by apps and drivers to register data sources in stores they
	// own.
	RegisterDatasource(metadata DataSourceMetadata) error
}

type keyValueClient struct {
	zestC        zest.ZestClient
	zEndpoint    string
	dEndpoint    string
	dataSourceID string
}

// NewKeyValueClient returns a new KeyValueClient to enable reading and writing of key value data to the store
// reqEndpoint is provided in the DATABOX_ZMQ_ENDPOINT environment varable to databox apps and drivers.
// dataSourceID is passed in the to apps in the environment varable and can be extracted from the hypercat
// drivers are responsible for managing their dataSourceIDs
func NewKeyValueClient(reqEndpoint string, dataSourceID string, enableLogging bool) (KeyValue_0_2_0, error) {

	serverKey, err := ioutil.ReadFile("/run/secrets/ZMQ_PUBLIC_KEY")
	if err != nil {
		return keyValueClient{}, err
	}

	kvc := keyValueClient{}
	kvc.zEndpoint = reqEndpoint
	kvc.dEndpoint = strings.Replace(reqEndpoint, ":5555", ":5556", 1)
	kvc.zestC, err = zest.New(kvc.zEndpoint, kvc.dEndpoint, string(serverKey), enableLogging)
	kvc.dataSourceID = dataSourceID

	return kvc, err
}

func (kvc keyValueClient) Write(payload string) error {

	path := "/kv/" + kvc.dataSourceID

	token, err := requestToken(kvc.zEndpoint+path, "POST")
	if err != nil {
		return err
	}

	err = kvc.zestC.Post(token, path, payload)
	if err != nil {
		return errors.New("Error posting data: " + err.Error())
	}

	return nil

}

func (kvc keyValueClient) Read() (string, error) {

	path := "/kv/" + kvc.dataSourceID

	token, err := requestToken(kvc.zEndpoint+path, "GET")
	if err != nil {
		return "", err
	}

	resp, getErr := kvc.zestC.Get(token, path)
	if getErr != nil {
		return "", errors.New("Error getting data: " + err.Error())
	}

	return resp, nil

}

func (kvc keyValueClient) Observe() (<-chan string, error) {

	path := "/kv/" + kvc.dataSourceID

	token, err := requestToken(kvc.zEndpoint+path, "GET")
	if err != nil {
		return nil, err
	}

	payloadChan, getErr := kvc.zestC.Observe(token, path)
	if getErr != nil {
		return nil, errors.New("Error observing: " + err.Error())
	}

	return payloadChan, nil

}

func (kvc keyValueClient) RegisterDatasource(metadata DataSourceMetadata) error {

	path := "/cat"

	token, err := requestToken(kvc.zEndpoint+path, "POST")
	if err != nil {
		return errors.New("Error getting token: " + err.Error())
	}
	hypercatJSON, err := dataSourceMetadataToHypercat(metadata)

	writeErr := kvc.zestC.Post(token, path, string(hypercatJSON))
	if writeErr != nil {
		return errors.New("Error writing: " + writeErr.Error())
	}

	return nil
}
