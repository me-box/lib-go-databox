package libDatabox

import (
	"encoding/json"
	"errors"

	zest "github.com/me-box/goZestClient"
)

type KVStore struct {
	csc         *CoreStoreClient
	contentType StoreContentType
}

func newKVStore(csc *CoreStoreClient, contentType StoreContentType) *KVStore {
	return &KVStore{
		csc:         csc,
		contentType: contentType,
	}
}

// Write Write will add data to the key value data store.
func (kvj *KVStore) Write(dataSourceID string, key string, payload []byte) error {

	path := "/kv/" + dataSourceID + "/" + key

	return kvj.csc.write(path, payload, kvj.contentType)

}

// Read will read the vale store at under tha key
// return data is a  object of the format {"timestamp":213123123,"data":[data-written-by-driver]}
func (kvj *KVStore) Read(dataSourceID string, key string) ([]byte, error) {

	path := "/kv/" + dataSourceID + "/" + key

	return kvj.csc.read(path, kvj.contentType)

}

// Delete deletes data under the key.
func (kvj *KVStore) Delete(dataSourceID string, key string) error {

	path := "/kv/" + dataSourceID + "/" + key

	return kvj.csc.delete(path, kvj.contentType)

}

// DeleteAll deletes all keys and data from the datasource.
func (kvj *KVStore) DeleteAll(dataSourceID string) error {

	path := "/kv/" + dataSourceID

	return kvj.csc.delete(path, kvj.contentType)

}

// ListKeys returns an array of key registed under the dataSourceID
func (kvj *KVStore) ListKeys(dataSourceID string) ([]string, error) {

	path := "/kv/" + dataSourceID + "/keys"

	data, err := kvj.csc.read(path, kvj.contentType)
	if err != nil {
		return []string{}, err
	}

	var keysArray []string

	err = json.Unmarshal(data, &keysArray)
	if err != nil {
		return []string{}, errors.New("KVListKeys: Error decoding data. " + err.Error())
	}
	return keysArray, nil

}

func (kvj *KVStore) Observe(dataSourceID string) (<-chan ObserveResponse, error) {

	path := "/kv/" + dataSourceID + "/*"

	return kvj.csc.observe(path, kvj.contentType, zest.ObserveModeData)

}

func (kvj *KVStore) ObserveKey(dataSourceID string, key string) (<-chan ObserveResponse, error) {

	path := "/kv/" + dataSourceID + "/" + key

	return kvj.csc.observe(path, kvj.contentType, zest.ObserveModeData)

}
