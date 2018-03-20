package libDatabox

import (
	"encoding/json"
	"errors"
	"strconv"
	"strings"

	zest "github.com/toshbrown/goZestClient"
)

type AggregationType string

type FilterType string

//Allowed values for FilterType and AggregationFunction
const (
	Equals            FilterType      = "equals"
	Contains          FilterType      = "contains"
	Sum               AggregationType = "sum"
	Count             AggregationType = "count"
	Min               AggregationType = "min"
	Max               AggregationType = "max"
	Mean              AggregationType = "mean"
	Median            AggregationType = "median"
	StandardDeviation AggregationType = "sd"
)

// Filter types to hold the required data to apply the filtering functions of the structured json API
type Filter struct {
	TagName    string
	FilterType FilterType
	Value      string
}

// JSONTimeSeriesQueryOptions described the options for the structured json API
type JSONTimeSeriesQueryOptions struct {
	AggregationFunction AggregationType
	Filter              *Filter
}

// JSONTimeSeries_0_3_0 described the the structured json timeseries API
type JSONTimeSeries_0_3_0 interface {
	// Write  will be timestamped with write time in ms since the unix epoch by the store
	Write(dataSourceID string, payload []byte) error
	// WriteAt will be timestamped with timestamp provided in ms since the unix epoch
	WriteAt(dataSourceID string, timestamp int64, payload []byte) error
	// Read the latest value.
	// return data is a JSON object of the format {"timestamp":213123123,"data":[data-written-by-driver]}
	Latest(dataSourceID string) ([]byte, error)
	// Read the earliest value.
	// return data is a JSON object of the format {"timestamp":213123123,"data":[data-written-by-driver]}
	Earliest(dataSourceID string) ([]byte, error)
	// Read the last N values.
	// return data is an array of JSON objects of the format {"timestamp":213123123,"data":[data-written-by-driver]}
	LastN(dataSourceID string, n int, opt JSONTimeSeriesQueryOptions) ([]byte, error)
	// Read the first N values.
	// return data is an array of JSON objects of the format {"timestamp":213123123,"data":[data-written-by-driver]}
	FirstN(dataSourceID string, n int, opt JSONTimeSeriesQueryOptions) ([]byte, error)
	// Read values written after the provided timestamp in in ms since the unix epoch.
	// return data is an array of JSON objects of the format {"timestamp":213123123,"data":[data-written-by-driver]}
	Since(dataSourceID string, sinceTimeStamp int64, opt JSONTimeSeriesQueryOptions) ([]byte, error)
	// Read values written between the start timestamp and end timestamp in in ms since the unix epoch.
	// return data is an array of JSON objects of the format {"timestamp":213123123,"data":[data-written-by-driver]}
	Range(dataSourceID string, formTimeStamp int64, toTimeStamp int64, opt JSONTimeSeriesQueryOptions) ([]byte, error)
	//Length retruns the number of records stored for that dataSourceID
	Length(dataSourceID string) (int, error)
	// Get notifications when a new value is written
	// the returned chan receives JsonObserveResponse of the form {"TimestampMS":213123123,"Json":byte[]}
	Observe(dataSourceID string) (<-chan JsonObserveResponse, error)
	// registerDatasource is used by apps and drivers to register data sources in stores they own.
	RegisterDatasource(metadata DataSourceMetadata) error
	// GetDatasourceCatalogue is used by drivers to get a list of registered data sources in stores they own.
	GetDatasourceCatalogue() ([]byte, error)
}

type jSONTimeSeriesClient struct {
	zestC     zest.ZestClient
	zEndpoint string
	dEndpoint string
}

// NewJSONTimeSeriesClient returns a new jSONTimeSeriesClient to enable interaction with a structured timeseries data store in JSON format.
// The data written must contain at least {"value":[any numeric value]}. This is used in the aggregation functions. Other data can be store and used at KV pairs to filter the data but it can not be processed.
// reqEndpoint is provided in the DATABOX_ZMQ_ENDPOINT environment varable to databox apps and drivers.
func NewJSONTimeSeriesClient(reqEndpoint string, enableLogging bool) (JSONTimeSeries_0_3_0, error) {

	var err error

	tsc := jSONTimeSeriesClient{}
	tsc.zEndpoint = reqEndpoint
	tsc.dEndpoint = strings.Replace(reqEndpoint, ":5555", ":5556", 1)
	tsc.zestC, err = zest.New(tsc.zEndpoint, tsc.dEndpoint, getServerKey(), enableLogging)

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
		return []byte(""), errors.New("Error getting latest data: " + getErr.Error())
	}

	return resp, nil

}

// Earliest will retrieve the first entry stored at the requested datasource ID
// return data is a JSON object of the format {"timestamp":213123123,"data":[data-written-by-driver]}
func (tsc jSONTimeSeriesClient) Earliest(dataSourceID string) ([]byte, error) {

	path := "/ts/" + dataSourceID + "/earliest"

	token, err := requestToken(tsc.zEndpoint+path, "GET")
	if err != nil {
		return []byte(""), err
	}

	resp, getErr := tsc.zestC.Get(token, path, "JSON")
	if getErr != nil {
		invalidateCache(tsc.zEndpoint+path, "GET")
		return []byte(""), errors.New("Error getting earliest data: " + getErr.Error())
	}

	return resp, nil

}

// LastN will retrieve the last N entries stored at the requested datasource ID
// return data is an array of JSON objects of the format {"timestamp":213123123,"data":[data-written-by-driver]}
func (tsc jSONTimeSeriesClient) LastN(dataSourceID string, n int, opt JSONTimeSeriesQueryOptions) ([]byte, error) {

	path := "/ts/" + dataSourceID + "/last/" + strconv.Itoa(n) + tsc.calculatePath(opt)

	token, err := requestToken(tsc.zEndpoint+path, "GET")
	if err != nil {
		return []byte(""), err
	}

	resp, getErr := tsc.zestC.Get(token, path, "JSON")
	if getErr != nil {
		invalidateCache(tsc.zEndpoint+path, "GET")
		return []byte(""), errors.New("Error getting latest data: " + getErr.Error())
	}

	return resp, nil

}

// FirstN will retrieve the first N entries stored at the requested datasource ID
// return data is an array of JSON objects of the format {"timestamp":213123123,"data":[data-written-by-driver]}
func (tsc jSONTimeSeriesClient) FirstN(dataSourceID string, n int, opt JSONTimeSeriesQueryOptions) ([]byte, error) {

	path := "/ts/" + dataSourceID + "/first/" + strconv.Itoa(n) + tsc.calculatePath(opt)

	token, err := requestToken(tsc.zEndpoint+path, "GET")
	if err != nil {
		return []byte(""), err
	}

	resp, getErr := tsc.zestC.Get(token, path, "JSON")
	if getErr != nil {
		invalidateCache(tsc.zEndpoint+path, "GET")
		return []byte(""), errors.New("Error getting latest data: " + getErr.Error())
	}

	return resp, nil

}

//Since will retrieve all entries since the requested timestamp (ms since unix epoch)
// return data is a JSON object of the format {"timestamp":213123123,"data":[data-written-by-driver]}
func (tsc jSONTimeSeriesClient) Since(dataSourceID string, sinceTimeStamp int64, opt JSONTimeSeriesQueryOptions) ([]byte, error) {

	path := "/ts/" + dataSourceID + "/since/" + strconv.FormatInt(sinceTimeStamp, 10) + tsc.calculatePath(opt)

	token, err := requestToken(tsc.zEndpoint+path, "GET")
	if err != nil {
		return []byte(""), err
	}

	resp, getErr := tsc.zestC.Get(token, path, "JSON")
	if getErr != nil {
		invalidateCache(tsc.zEndpoint+path, "GET")
		return []byte(""), errors.New("Error getting latest data: " + getErr.Error())
	}

	return resp, nil

}

// Range will retrieve all entries between  formTimeStamp and toTimeStamp timestamp in ms since unix epoch
// return data is a JSON object of the format {"timestamp":213123123,"data":[data-written-by-driver]}
func (tsc jSONTimeSeriesClient) Range(dataSourceID string, formTimeStamp int64, toTimeStamp int64, opt JSONTimeSeriesQueryOptions) ([]byte, error) {

	path := "/ts/" + dataSourceID + "/range/" + strconv.FormatInt(formTimeStamp, 10) + "/" + strconv.FormatInt(toTimeStamp, 10) + tsc.calculatePath(opt)

	token, err := requestToken(tsc.zEndpoint+path, "GET")
	if err != nil {
		return []byte(""), err
	}

	resp, getErr := tsc.zestC.Get(token, path, "JSON")
	if getErr != nil {
		invalidateCache(tsc.zEndpoint+path, "GET")
		return []byte(""), errors.New("Error getting latest data: " + getErr.Error())
	}

	return resp, nil

}

//Length retruns the number of records stored for that dataSourceID
func (tsc jSONTimeSeriesClient) Length(dataSourceID string) (int, error) {

	path := "/ts/" + dataSourceID + "/length"

	token, err := requestToken(tsc.zEndpoint+path, "GET")
	if err != nil {
		return 0, err
	}

	resp, getErr := tsc.zestC.Get(token, path, "JSON")
	if getErr != nil {
		invalidateCache(tsc.zEndpoint+path, "GET")
		return 0, errors.New("Error getting latest data: " + getErr.Error())
	}

	type legnthResult struct {
		Length int `json:"length"`
	}

	var val legnthResult
	err = json.Unmarshal(resp, &val)
	if err != nil {
		return 0, err
	}

	return val.Length, nil
}

func (tsc jSONTimeSeriesClient) Observe(dataSourceID string) (<-chan JsonObserveResponse, error) {

	path := "/ts/" + dataSourceID

	token, err := requestToken(tsc.zEndpoint+path, "GET")
	if err != nil {
		return nil, err
	}

	payloadChan, getErr := tsc.zestC.Observe(token, path, "JSON", 0)
	if getErr != nil {
		invalidateCache(tsc.zEndpoint+path, "GET")
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

// RegisterDatasource is used by apps and drivers to register datasource in stores they
// own.
func (tsc jSONTimeSeriesClient) RegisterDatasource(metadata DataSourceMetadata) error {

	path := "/cat"

	token, err := requestToken(tsc.zEndpoint+path, "POST")
	if err != nil {
		return err
	}
	hypercatJSON, err := dataSourceMetadataToHypercat(metadata, tsc.zEndpoint+"/ts/")

	writeErr := tsc.zestC.Post(token, path, hypercatJSON, "JSON")
	if writeErr != nil {
		invalidateCache(tsc.zEndpoint+path, "POST")
		return errors.New("Error writing: " + writeErr.Error())
	}

	return nil
}

func (tsc jSONTimeSeriesClient) GetDatasourceCatalogue() ([]byte, error) {
	path := "/cat"

	token, err := requestToken(tsc.zEndpoint+path, "GET")
	if err != nil {
		return nil, err
	}

	hypercatJSON, getErr := tsc.zestC.Get(token, path, "JSON")
	if getErr != nil {
		invalidateCache(tsc.zEndpoint+path, "GET")
		return []byte{}, errors.New("Error reading: " + getErr.Error())
	}

	return hypercatJSON, nil
}

func (tsc jSONTimeSeriesClient) calculatePath(opt JSONTimeSeriesQueryOptions) string {
	aggregationPath := ""
	if opt.AggregationFunction != "" {
		aggregationPath = "/" + string(opt.AggregationFunction)
	}

	if opt.Filter == nil || opt.Filter.TagName == "" || opt.Filter.FilterType == "" || opt.Filter.Value == "" {
		return aggregationPath
	}

	return "/filter/" + string(opt.Filter.TagName) + "/" + string(opt.Filter.FilterType) + "/" + string(opt.Filter.Value) + aggregationPath
}
