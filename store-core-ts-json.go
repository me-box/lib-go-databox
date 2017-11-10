package libDatabox

import (
	"errors"
	"io/ioutil"
	"strconv"
	"strings"

	zest "github.com/toshbrown/goZestClient"
)

type JSONTimeSeries_0_2_0 interface {
	// Write  will be timestamped with write time in ms since the unix epoch by the store
	Write(dataSourceID string, payload []byte) error
	// WriteAt will be timestamped with timestamp provided in ms since the unix epoch
	WriteAt(dataSourceID string, timestamp int64, payload []byte) error
	// Read the latest value.
	// return data is a JSON object of the format {"timestamp":213123123,"data":[data-written-by-driver]}
	Latest(dataSourceID string) ([]byte, error)
	// Read the last N values.
	// return data is an array of JSON objects of the format {"timestamp":213123123,"data":[data-written-by-driver]}
	LastN(dataSourceID string, n int) ([]byte, error)
	// Read values written after the provided timestamp in in ms since the unix epoch.
	// return data is an array of JSON objects of the format {"timestamp":213123123,"data":[data-written-by-driver]}
	Since(dataSourceID string, sinceTimeStamp int64) ([]byte, error)
	// Read values written between the start timestamp and end timestamp in in ms since the unix epoch.
	// return data is an array of JSON objects of the format {"timestamp":213123123,"data":[data-written-by-driver]}
	Range(dataSourceID string, formTimeStamp int64, toTimeStamp int64) ([]byte, error)
	// Get notifications when a new value is written
	// the returned chan receives valuse of the form {"timestamp":213123123,"data":[data-written-by-driver]}
	Observe(dataSourceID string) (<-chan []byte, error)
	// registerDatasource is used by apps and drivers to register data sources in stores they own.
	RegisterDatasource(metadata DataSourceMetadata) error
}

type jSONTimeSeriesClient struct {
	zestC     zest.ZestClient
	zEndpoint string
	dEndpoint string
}

// NewJSONTimeSeriesClient returns a new jSONTimeSeriesClient to enable interaction with a time series data store in JSON format
// reqEndpoint is provided in the DATABOX_ZMQ_ENDPOINT environment varable to databox apps and drivers.
func NewJSONTimeSeriesClient(reqEndpoint string, enableLogging bool) (JSONTimeSeries_0_2_0, error) {

	serverKey, err := ioutil.ReadFile("/run/secrets/ZMQ_PUBLIC_KEY")
	if err != nil {
		return jSONTimeSeriesClient{}, err
	}

	tsc := jSONTimeSeriesClient{}
	tsc.zEndpoint = reqEndpoint
	tsc.dEndpoint = strings.Replace(reqEndpoint, ":5555", ":5556", 1)
	tsc.zestC, err = zest.New(tsc.zEndpoint, tsc.dEndpoint, string(serverKey), enableLogging)

	return tsc, err
}

// Write will add data to the times series data store. Data will be time stamped at insertion (format ms since 1970)
func (tsc jSONTimeSeriesClient) Write(dataSourceID string, payload []byte) error {

	path := "/ts/" + dataSourceID

	token, err := requestToken(tsc.zEndpoint+path, "POST")
	if err != nil {
		return err
	}

	err = tsc.zestC.Post(token, path, payload, "JSON")
	if err != nil {
		invalidateCache(tsc.zEndpoint+path, "POST")
		return errors.New("Error writing: " + err.Error())
	}

	return nil

}

// WriteAt will add data to the times series data store. Data will be time stamped with the timstamp provided in the
// timstamp paramiter (format ms since 1970)
func (tsc jSONTimeSeriesClient) WriteAt(dataSourceID string, timstamp int64, payload []byte) error {

	path := "/ts/" + dataSourceID + "/at/"

	token, err := requestToken(tsc.zEndpoint+path+"*", "POST")
	if err != nil {
		return err
	}

	path = path + strconv.FormatInt(timstamp, 10)

	err = tsc.zestC.Post(token, path, payload, "JSON")
	if err != nil {
		invalidateCache(tsc.zEndpoint+path+"*", "POST")
		return errors.New("Error writing: " + err.Error())
	}

	return nil

}

//Latest will retrieve the last entry stored at the requested datasource ID
// return data is a JSON object of the format {"timestamp":213123123,"data":[data-written-by-driver]}
func (tsc jSONTimeSeriesClient) Latest(dataSourceID string) ([]byte, error) {

	path := "/ts/" + dataSourceID + "/latest"

	token, err := requestToken(tsc.zEndpoint+path, "GET")
	if err != nil {
		return []byte(""), err
	}

	resp, getErr := tsc.zestC.Get(token, path, "JSON")
	if getErr != nil {
		invalidateCache(tsc.zEndpoint+path, "GET")
		return []byte(""), errors.New("Error getting latest data: " + err.Error())
	}

	return resp, nil

}

// LastN will retrieve the last N entries stored at the requested datasource ID
// return data is an array of JSON objects of the format {"timestamp":213123123,"data":[data-written-by-driver]}
func (tsc jSONTimeSeriesClient) LastN(dataSourceID string, n int) ([]byte, error) {

	path := "/ts/" + dataSourceID + "/last/" + strconv.Itoa(n)

	token, err := requestToken(tsc.zEndpoint+path, "GET")
	if err != nil {
		return []byte(""), err
	}

	resp, getErr := tsc.zestC.Get(token, path, "JSON")
	if getErr != nil {
		invalidateCache(tsc.zEndpoint+path, "GET")
		return []byte(""), errors.New("Error getting latest data: " + err.Error())
	}

	return resp, nil

}

//Since will retrieve all entries since the requested timestamp (ms since unix epoch)
// return data is a JSON object of the format {"timestamp":213123123,"data":[data-written-by-driver]}
func (tsc jSONTimeSeriesClient) Since(dataSourceID string, sinceTimeStamp int64) ([]byte, error) {

	path := "/ts/" + dataSourceID + "/since/" + strconv.FormatInt(sinceTimeStamp, 10)

	token, err := requestToken(tsc.zEndpoint+path, "GET")
	if err != nil {
		return []byte(""), err
	}

	resp, getErr := tsc.zestC.Get(token, path, "JSON")
	if getErr != nil {
		invalidateCache(tsc.zEndpoint+path, "GET")
		return []byte(""), errors.New("Error getting latest data: " + err.Error())
	}

	return resp, nil

}

// Range will retrieve all entries between  formTimeStamp and toTimeStamp timestamp in ms since unix epoch
// return data is a JSON object of the format {"timestamp":213123123,"data":[data-written-by-driver]}
func (tsc jSONTimeSeriesClient) Range(dataSourceID string, formTimeStamp int64, toTimeStamp int64) ([]byte, error) {

	path := "/ts/" + dataSourceID + "/range/" + strconv.FormatInt(formTimeStamp, 10) + "/" + strconv.FormatInt(toTimeStamp, 10)

	token, err := requestToken(tsc.zEndpoint+path, "GET")
	if err != nil {
		return []byte(""), err
	}

	resp, getErr := tsc.zestC.Get(token, path, "JSON")
	if getErr != nil {
		invalidateCache(tsc.zEndpoint+path, "GET")
		return []byte(""), errors.New("Error getting latest data: " + err.Error())
	}

	return resp, nil

}

func (tsc jSONTimeSeriesClient) Observe(dataSourceID string) (<-chan []byte, error) {

	path := "/ts/" + dataSourceID

	token, err := requestToken(tsc.zEndpoint+path, "GET")
	if err != nil {
		return nil, err
	}

	payloadChan, getErr := tsc.zestC.Observe(token, path, "JSON")
	if getErr != nil {
		invalidateCache(tsc.zEndpoint+path, "GET")
		return nil, errors.New("Error observing: " + err.Error())
	}

	return payloadChan, nil

}

// RegisterDatasource is used by apps and drivers to register datasource in stores they
// own.
func (tsc jSONTimeSeriesClient) RegisterDatasource(metadata DataSourceMetadata) error {

	path := "/cat"

	token, err := requestToken(tsc.zEndpoint+path, "POST")
	if err != nil {
		return err
	}
	hypercatJSON, err := dataSourceMetadataToHypercat(metadata)

	writeErr := tsc.zestC.Post(token, path, hypercatJSON, "JSON")
	if writeErr != nil {
		invalidateCache(tsc.zEndpoint+path, "POST")
		return errors.New("Error writing: " + writeErr.Error())
	}

	return nil
}
