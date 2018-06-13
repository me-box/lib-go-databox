package libDatabox

import (
	"encoding/json"
	"errors"
)

// KVTextWrite Write will add data to the key value data store.
func (csc *CoreStoreClient) KVTextWrite(dataSourceID string, key string, payload []byte) error {

	path := "/kv/" + dataSourceID + "/" + key

	return csc.write(path, payload, ContentTypeTEXT)

}

// KVTextRead will read the vale store at under tha key
// return data is a Text object of the format {"timestamp":213123123,"data":[data-written-by-driver]}
func (csc *CoreStoreClient) KVTextRead(dataSourceID string, key string) ([]byte, error) {

	path := "/kv/" + dataSourceID + "/" + key

	return csc.read(path, ContentTypeTEXT)

}

// KVTextDelete deletes data under the key.
func (csc *CoreStoreClient) KVTextDelete(dataSourceID string, key string) error {

	path := "/kv/" + dataSourceID + "/" + key

	return csc.delete(path, ContentTypeTEXT)

}

// KVTextDeleteAll deletes all keys and data from the datasource.
func (csc *CoreStoreClient) KVTextDeleteAll(dataSourceID string) error {

	path := "/kv/" + dataSourceID

	return csc.delete(path, ContentTypeTEXT)

}

// KVTextListKeys returns an array of key registed under the dataSourceID
func (csc *CoreStoreClient) KVTextListKeys(dataSourceID string) ([]string, error) {

	path := "/kv/" + dataSourceID + "/keys"

	data, err := csc.read(path, ContentTypeTEXT)
	if err != nil {
		return []string{}, err
	}

	var keysArray []string

	err = json.Unmarshal(data, &keysArray)
	if err != nil {
		return []string{}, errors.New("KVTextListKeys: Error decoding data. " + err.Error())
	}
	return keysArray, nil

}

func (csc *CoreStoreClient) KVTextObserve(dataSourceID string) (<-chan []byte, error) {

	path := "/kv/" + dataSourceID + "/*"

	return csc.observe(path, ContentTypeTEXT)

}

func (csc *CoreStoreClient) KVTextObserveKey(dataSourceID string, key string) (<-chan []byte, error) {

	path := "/kv/" + dataSourceID + "/" + key

	return csc.observe(path, ContentTypeTEXT)

}
