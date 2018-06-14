package libDatabox

import (
	"encoding/json"
	"errors"
	"strconv"
)

type TSBlobStore struct {
	csc         *CoreStoreClient
	contentType StoreContentType
}

func newTSBlobStore(csc *CoreStoreClient, contentType StoreContentType) *TSBlobStore {
	return &TSBlobStore{
		csc:         csc,
		contentType: contentType,
	}
}

// Write will add data to the times series data store. Data will be time stamped at insertion (format ms since 1970)
func (tbs *TSBlobStore) Write(dataSourceID string, payload []byte) error {

	path := "/ts/blob/" + dataSourceID

	return tbs.csc.write(path, payload, tbs.contentType)

}

// WriteAt will add data to the times series data store. Data will be time stamped with the timstamp provided in the
// timstamp paramiter (format ms since 1970)
func (tbs *TSBlobStore) WriteAt(dataSourceID string, timstamp int64, payload []byte) error {

	path := "/ts/blob/" + dataSourceID + "/at/"

	token, err := tbs.csc.Arbiter.RequestToken(tbs.csc.ZEndpoint+path+"*", "POST")
	if err != nil {
		return err
	}

	path = path + strconv.FormatInt(timstamp, 10)

	err = tbs.csc.ZestC.Post(string(token), path, payload, string(tbs.contentType))
	if err != nil {
		tbs.csc.Arbiter.InvalidateCache(tbs.csc.ZEndpoint+path+"*", "POST")
		return errors.New("Error writing: " + err.Error())
	}

	return nil

}

//TSBlobLatest will retrieve the last entry stored at the requested datasource ID
// return data is a byte array contingin  of the format
// {"timestamp":213123123,"data":[data-written-by-driver]}
func (tbs *TSBlobStore) Latest(dataSourceID string) ([]byte, error) {

	path := "/ts/blob/" + dataSourceID + "/latest"

	return tbs.csc.read(path, tbs.contentType)

}

// Earliest will retrieve the first entry stored at the requested datasource ID
// return data is a byte array contingin  of the format
// {"timestamp":213123123,"data":[data-written-by-driver]}
func (tbs *TSBlobStore) Earliest(dataSourceID string) ([]byte, error) {

	path := "/ts/blob/" + dataSourceID + "/earliest"

	return tbs.csc.read(path, tbs.contentType)

}

// LastN will retrieve the last N entries stored at the requested datasource ID
// return data is a byte array contingin  of the format
// {"timestamp":213123123,"data":[data-written-by-driver]}
func (tbs *TSBlobStore) LastN(dataSourceID string, n int) ([]byte, error) {

	path := "/ts/blob/" + dataSourceID + "/last/" + strconv.Itoa(n)

	return tbs.csc.read(path, tbs.contentType)

}

// FirstN will retrieve the first N entries stored at the requested datasource ID
// return data is a byte array contingin  of the format
// {"timestamp":213123123,"data":[data-written-by-driver]}
func (tbs *TSBlobStore) FirstN(dataSourceID string, n int) ([]byte, error) {

	path := "/ts/blob/" + dataSourceID + "/first/" + strconv.Itoa(n)

	return tbs.csc.read(path, tbs.contentType)

}

// Since will retrieve all entries since the requested timestamp (ms since unix epoch)
// return data is a byte array contingin  of the format
// {"timestamp":213123123,"data":[data-written-by-driver]}
func (tbs *TSBlobStore) Since(dataSourceID string, sinceTimeStamp int64) ([]byte, error) {

	path := "/ts/blob/" + dataSourceID + "/since/" + strconv.FormatInt(sinceTimeStamp, 10)

	return tbs.csc.read(path, tbs.contentType)

}

// Range will retrieve all entries between  formTimeStamp and toTimeStamp timestamp in ms since unix epoch
// return data is a byte array contingin  of the format
// {"timestamp":213123123,"data":[data-written-by-driver]}
func (tbs *TSBlobStore) Range(dataSourceID string, formTimeStamp int64, toTimeStamp int64) ([]byte, error) {

	path := "/ts/blob/" + dataSourceID + "/range/" + strconv.FormatInt(formTimeStamp, 10) + "/" + strconv.FormatInt(toTimeStamp, 10)

	return tbs.csc.read(path, tbs.contentType)

}

//TSBlobLength returns then number of items stored in the timeseries
func (tbs *TSBlobStore) Length(dataSourceID string) (int, error) {
	path := "/ts/blob/" + dataSourceID + "/length"

	resp, getErr := tbs.csc.read(path, tbs.contentType)
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

// Observe allows you to get notifications when a new value is written by a driver
// the returned chan receives chan []byte continuing json of the
// form {"TimestampMS":213123123,"Json":byte[]}
func (tbs *TSBlobStore) Observe(dataSourceID string) (<-chan []byte, error) {

	path := "/ts/blob/" + dataSourceID

	return tbs.csc.observe(path, tbs.contentType)

}
