package libDatabox

import (
	"encoding/json"
	"errors"
	"strconv"
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

// TimeSeriesQueryOptions describes the query options for the structured json API
type TimeSeriesQueryOptions struct {
	AggregationFunction AggregationType
	Filter              *Filter
}

type TSStore struct {
	csc         *CoreStoreClient
	contentType StoreContentType
}

func newTSStore(csc *CoreStoreClient, contentType StoreContentType) *TSStore {
	return &TSStore{
		csc:         csc,
		contentType: contentType,
	}
}

// Write will add data to the times series data store. Data will be time stamped at insertion (format ms since 1970)
func (tsc TSStore) Write(dataSourceID string, payload []byte) error {

	path := "/ts/" + dataSourceID

	return tsc.csc.write(path, payload, ContentTypeJSON)

}

// WriteAt will add data to the times series data store. Data will be time stamped with the timstamp provided in the
// timstamp paramiter (format ms since 1970)
func (tsc TSStore) WriteAt(dataSourceID string, timstamp int64, payload []byte) error {

	path := "/ts/" + dataSourceID + "/at/"

	token, err := tsc.csc.Arbiter.RequestToken(tsc.csc.ZEndpoint+path+"*", "POST")

	path = path + strconv.FormatInt(timstamp, 10)

	_, err = tsc.csc.ZestC.Post(string(token), path, payload, string(ContentTypeJSON))
	if err != nil {
		tsc.csc.Arbiter.InvalidateCache(tsc.csc.ZEndpoint+path+"*", "POST")
		return errors.New("Error writing: " + err.Error())
	}

	return nil

}

//Latest will retrieve the last entry stored at the requested datasource ID
// return data is a JSON object of the format {"timestamp":213123123,"data":[data-written-by-driver]}
func (tsc TSStore) Latest(dataSourceID string) ([]byte, error) {

	path := "/ts/" + dataSourceID + "/latest"

	return tsc.csc.read(path, ContentTypeJSON)

}

// Earliest will retrieve the first entry stored at the requested datasource ID
// return data is a JSON object of the format {"timestamp":213123123,"data":[data-written-by-driver]}
func (tsc TSStore) Earliest(dataSourceID string) ([]byte, error) {

	path := "/ts/" + dataSourceID + "/earliest"

	return tsc.csc.read(path, ContentTypeJSON)

}

// LastN will retrieve the last N entries stored at the requested datasource ID
// return data is an array of JSON objects of the format {"timestamp":213123123,"data":[data-written-by-driver]}
func (tsc TSStore) LastN(dataSourceID string, n int, opt TimeSeriesQueryOptions) ([]byte, error) {

	path := "/ts/" + dataSourceID + "/last/" + strconv.Itoa(n) + tsc.calculatePath(opt)

	return tsc.csc.read(path, ContentTypeJSON)

}

// FirstN will retrieve the first N entries stored at the requested datasource ID
// return data is an array of JSON objects of the format {"timestamp":213123123,"data":[data-written-by-driver]}
func (tsc TSStore) FirstN(dataSourceID string, n int, opt TimeSeriesQueryOptions) ([]byte, error) {

	path := "/ts/" + dataSourceID + "/first/" + strconv.Itoa(n) + tsc.calculatePath(opt)

	return tsc.csc.read(path, ContentTypeJSON)

}

//Since will retrieve all entries since the requested timestamp (ms since unix epoch)
// return data is a JSON object of the format {"timestamp":213123123,"data":[data-written-by-driver]}
func (tsc TSStore) Since(dataSourceID string, sinceTimeStamp int64, opt TimeSeriesQueryOptions) ([]byte, error) {

	path := "/ts/" + dataSourceID + "/since/" + strconv.FormatInt(sinceTimeStamp, 10) + tsc.calculatePath(opt)

	return tsc.csc.read(path, ContentTypeJSON)

}

// Range will retrieve all entries between  formTimeStamp and toTimeStamp timestamp in ms since unix epoch
// return data is a JSON object of the format {"timestamp":213123123,"data":[data-written-by-driver]}
func (tsc TSStore) Range(dataSourceID string, formTimeStamp int64, toTimeStamp int64, opt TimeSeriesQueryOptions) ([]byte, error) {

	path := "/ts/" + dataSourceID + "/range/" + strconv.FormatInt(formTimeStamp, 10) + "/" + strconv.FormatInt(toTimeStamp, 10) + tsc.calculatePath(opt)

	return tsc.csc.read(path, ContentTypeJSON)

}

//Length retruns the number of records stored for that dataSourceID
func (tsc TSStore) Length(dataSourceID string) (int, error) {

	path := "/ts/" + dataSourceID + "/length"

	resp, getErr := tsc.csc.read(path, ContentTypeJSON)
	if getErr != nil {
		return 0, getErr
	}

	type legnthResult struct {
		Length int `json:"length"`
	}

	var val legnthResult
	err := json.Unmarshal(resp, &val)
	if err != nil {
		return 0, err
	}

	return val.Length, nil
}

func (tsc TSStore) Observe(dataSourceID string) (<-chan []byte, error) {

	path := "/ts/" + dataSourceID

	return tsc.csc.observe(path, ContentTypeJSON)

}

func (tsc TSStore) calculatePath(opt TimeSeriesQueryOptions) string {
	aggregationPath := ""
	if opt.AggregationFunction != "" {
		aggregationPath = "/" + string(opt.AggregationFunction)
	}

	if opt.Filter == nil || opt.Filter.TagName == "" || opt.Filter.FilterType == "" || opt.Filter.Value == "" {
		return aggregationPath
	}

	return "/filter/" + string(opt.Filter.TagName) + "/" + string(opt.Filter.FilterType) + "/" + string(opt.Filter.Value) + aggregationPath
}
