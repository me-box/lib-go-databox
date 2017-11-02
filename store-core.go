package libDatabox

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"strconv"
	"strings"

	zest "github.com/toshbrown/goZestClient"
)

type KeyValueClient struct {
	zestC     zest.ZestClient
	zEndpoint string
	dEndpoint string
}

func NewKeyValueClient(reqEndpoint string, enableLogging bool) (KeyValueClient, error) {

	serverKey, err := ioutil.ReadFile("/run/secrets/ZMQ_PUBLIC_KEY")
	if err != nil {
		return KeyValueClient{}, err
	}

	kvc := KeyValueClient{}
	kvc.zEndpoint = reqEndpoint
	kvc.dEndpoint = strings.Replace(reqEndpoint, ":5555", ":5556", 1)
	kvc.zestC = zest.New(kvc.zEndpoint, kvc.dEndpoint, string(serverKey), enableLogging)

	return kvc, nil
}

func (kvc KeyValueClient) Write(dataSourceID string, payload string) error {

	path := "/kv/" + dataSourceID

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

func (kvc KeyValueClient) Read(dataSourceID string) (string, error) {

	path := "/kv/" + dataSourceID

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

// RegisterDatasource is used by apps and drivers to register datasource in stores they
// own.
func (kvc KeyValueClient) RegisterDatasource(dataSourceID string, metadata StoreMetadata) error {

	path := "/cat"

	token, err := requestToken(kvc.zEndpoint+path, "POST")
	if err != nil {
		return errors.New("Error getting token: " + err.Error())
	}
	hypercatJSON, err := storeMetadataToJSON(metadata)

	writeErr := kvc.zestC.Post(token, path, string(hypercatJSON))
	if writeErr != nil {
		return errors.New("Error writing: " + writeErr.Error())
	}

	return nil
}

type TimeSeriesClient struct {
	zestC     zest.ZestClient
	zEndpoint string
	dEndpoint string
}

// NewTimeSeriesClient returns a new KeyTimeSeriesClient to enable interaction with a time series data store
// reqEndpoint is provided in the DATABOX_ZMQ_ENDPOINT environment varable to databox app and drivers.
func NewTimeSeriesClient(reqEndpoint string, enableLogging bool) (TimeSeriesClient, error) {

	serverKey, err := ioutil.ReadFile("/run/secrets/ZMQ_PUBLIC_KEY")
	if err != nil {
		return TimeSeriesClient{}, err
	}

	tsc := TimeSeriesClient{}
	tsc.zEndpoint = reqEndpoint
	tsc.dEndpoint = strings.Replace(reqEndpoint, ":5555", ":5556", 1)
	tsc.zestC = zest.New(tsc.zEndpoint, tsc.dEndpoint, string(serverKey), enableLogging)

	return tsc, nil
}

// Write will add data to the times series data store. Data will be time stamped at insertion (format ms since 1970)
func (tsc TimeSeriesClient) Write(dataSourceID string, payload string) error {

	path := "/ts/" + dataSourceID

	token, err := requestToken(tsc.zEndpoint+path, "POST")
	if err != nil {
		return err
	}

	err = tsc.zestC.Post(token, path, payload)
	if err != nil {
		return errors.New("Error writing: " + err.Error())
	}

	return nil

}

// WriteAt will add data to the times series data store. Data will be time stamped with the timstamp provided in the
// timstamp paramiter (format ms since 1970)
func (tsc TimeSeriesClient) WriteAt(dataSourceID string, timstamp int64, payload string) error {

	path := "/ts/" + dataSourceID + "/at/" + strconv.FormatInt(timstamp, 10)

	token, err := requestToken(tsc.zEndpoint+path, "POST")
	if err != nil {
		return err
	}

	err = tsc.zestC.Post(token, path, payload)
	if err != nil {
		return errors.New("Error writing: " + err.Error())
	}

	return nil

}

//Latest will retrieve the last entry stored at the requested datasource ID
// return data is a JSON object of the format {"timestamp":213123123,"data":[data-written-by-driver]}
func (tsc TimeSeriesClient) Latest(dataSourceID string) (string, error) {

	path := "/ts/" + dataSourceID + "/latest"

	token, err := requestToken(tsc.zEndpoint+path, "GET")
	if err != nil {
		return "", err
	}

	resp, getErr := tsc.zestC.Get(token, path)
	if getErr != nil {
		return "", errors.New("Error getting latest data: " + err.Error())
	}

	return resp, nil

}

// RegisterDatasource is used by apps and drivers to register datasource in stores they
// own.
func (tsc TimeSeriesClient) RegisterDatasource(dataSourceID string, metadata StoreMetadata) error {

	path := "/cat"

	token, err := requestToken(tsc.zEndpoint+path, "POST")
	if err != nil {
		return err
	}
	hypercatJSON, err := storeMetadataToJSON(metadata)

	writeErr := tsc.zestC.Post(token, path, string(hypercatJSON))
	if writeErr != nil {
		return errors.New("Error writing: " + writeErr.Error())
	}

	return nil
}

func storeMetadataToJSON(metadata StoreMetadata) ([]byte, error) {

	if metadata.Description == "" ||
		metadata.ContentType == "" ||
		metadata.Vendor == "" ||
		metadata.DataSourceType == "" ||
		metadata.DataSourceID == "" ||
		metadata.StoreType == "" {

		return nil, errors.New("Missing required metadata")
	}

	cat := hypercat{}
	cat.ItemMetadata = append(cat.ItemMetadata, relValPair{Rel: "urn:X-hypercat:rels:hasDescription:en", Val: metadata.Description})
	cat.ItemMetadata = append(cat.ItemMetadata, relValPair{Rel: "urn:X-hypercat:rels:isContentType", Val: metadata.ContentType})
	cat.ItemMetadata = append(cat.ItemMetadata, relValPair{Rel: "urn:X-databox:rels:hasVendor", Val: metadata.Vendor})
	cat.ItemMetadata = append(cat.ItemMetadata, relValPair{Rel: "urn:X-databox:rels:hasType", Val: metadata.DataSourceType})
	cat.ItemMetadata = append(cat.ItemMetadata, relValPair{Rel: "urn:X-databox:rels:hasDatasourceid", Val: metadata.DataSourceID})
	cat.ItemMetadata = append(cat.ItemMetadata, relValPair{Rel: "urn:X-databox:rels:hasStoreType", Val: metadata.StoreType})

	if metadata.IsActuator {
		cat.ItemMetadata = append(cat.ItemMetadata, relValPairBool{Rel: "urn:X-databox:rels:isActuator", Val: true})
	}

	if metadata.Location != "" {
		cat.ItemMetadata = append(cat.ItemMetadata, relValPair{Rel: "urn:X-databox:rels:hasLocation", Val: metadata.Location})
	}

	if metadata.Unit != "" {
		cat.ItemMetadata = append(cat.ItemMetadata, relValPair{Rel: "urn:X-databox:rels:hasUnit", Val: metadata.Unit})
	}

	return json.Marshal(cat)

}
