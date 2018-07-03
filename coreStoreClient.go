package libDatabox

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"

	zest "github.com/me-box/goZestClient"
)

type CoreStoreClient struct {
	ZestC      zest.ZestClient
	Arbiter    *ArbiterClient
	ZEndpoint  string
	DEndpoint  string
	KVJSON     *KVStore
	KVText     *KVStore
	KVBin      *KVStore
	TSBlobJSON *TSBlobStore
	TSBlobText *TSBlobStore
	TSBlobBin  *TSBlobStore
	TSJSON     *TSStore
}

func NewCoreStoreClient(arbiterClient *ArbiterClient, zmqPublicKeyPath string, storeEndPoint string, enableLogging bool) *CoreStoreClient {
	csc := &CoreStoreClient{
		Arbiter: arbiterClient,
	}

	//get the server key
	serverKey, err := ioutil.ReadFile(zmqPublicKeyPath)
	if err != nil {
		fmt.Println("Warning:: failed to read ZMQ_PUBLIC_KEY using default value")
		serverKey = []byte("vl6wu0A@XP?}Or/&BR#LSxn>A+}L)p44/W[wXL3<")
	}

	csc.ZEndpoint = storeEndPoint
	csc.DEndpoint = strings.Replace(storeEndPoint, ":5555", ":5556", 1)
	csc.ZestC, err = zest.New(csc.ZEndpoint, csc.DEndpoint, string(serverKey), enableLogging)
	if err != nil {
		fmt.Println("[NewCoreStoreClient] Error zest.New ", err.Error())
	}

	csc.KVJSON = newKVStore(csc, ContentTypeJSON)
	csc.KVText = newKVStore(csc, ContentTypeTEXT)
	csc.KVBin = newKVStore(csc, ContentTypeBINARY)
	csc.TSBlobJSON = newTSBlobStore(csc, ContentTypeJSON)
	csc.TSBlobText = newTSBlobStore(csc, ContentTypeTEXT)
	csc.TSBlobBin = newTSBlobStore(csc, ContentTypeBINARY)
	csc.TSJSON = newTSStore(csc, ContentTypeBINARY)
	return csc
}

func (csc *CoreStoreClient) GetStoreDataSourceCatalogue(href string) (HypercatRoot, error) {

	target := href + "/cat"
	method := "GET"

	token, err := csc.Arbiter.RequestToken(target, method)
	if err != nil {
		return HypercatRoot{}, err
	}
	//log.Debug("[GetStoreDataSourceCatalogue] got Token: " + string(token))

	hypercatJSON, getErr := csc.ZestC.Get(string(token), "/cat", "JSON")
	if getErr != nil {
		return HypercatRoot{}, err
	}
	//log.Debug("[GetStoreDataSourceCatalogue] got store cat: " + string(hypercatJSON))
	cat := HypercatRoot{}
	json.Unmarshal(hypercatJSON, &cat)

	return cat, nil

}

// RegisterDatasource is used by apps and drivers to register datasource in stores they
// own.
func (csc *CoreStoreClient) RegisterDatasource(metadata DataSourceMetadata) error {

	path := "/cat"

	token, err := csc.Arbiter.RequestToken(csc.ZEndpoint+path, "POST")
	if err != nil {
		return err
	}
	hypercatJSON, err := csc.dataSourceMetadataToHypercat(metadata, csc.ZEndpoint)

	_, writeErr := csc.ZestC.Post(string(token), path, hypercatJSON, "JSON")
	if writeErr != nil {
		csc.Arbiter.InvalidateCache(csc.ZEndpoint+path, "POST")
		return errors.New("Error writing: " + writeErr.Error())
	}

	return nil
}

//dataSourceMetadataToHypercat converts a DataSourceMetadata instance to json for registering a data source
func (csc *CoreStoreClient) dataSourceMetadataToHypercat(metadata DataSourceMetadata, endPoint string) ([]byte, error) {

	if metadata.Description == "" ||
		metadata.ContentType == "" ||
		metadata.Vendor == "" ||
		metadata.DataSourceType == "" ||
		metadata.DataSourceID == "" ||
		metadata.StoreType == "" {

		return nil, errors.New("Missing required metadata")
	}

	cat := HypercatItem{}
	cat.ItemMetadata = append(cat.ItemMetadata, RelValPair{Rel: "urn:X-hypercat:rels:hasDescription:en", Val: metadata.Description})
	cat.ItemMetadata = append(cat.ItemMetadata, RelValPair{Rel: "urn:X-hypercat:rels:isContentType", Val: metadata.ContentType})
	cat.ItemMetadata = append(cat.ItemMetadata, RelValPair{Rel: "urn:X-databox:rels:hasVendor", Val: metadata.Vendor})
	cat.ItemMetadata = append(cat.ItemMetadata, RelValPair{Rel: "urn:X-databox:rels:hasType", Val: metadata.DataSourceType})
	cat.ItemMetadata = append(cat.ItemMetadata, RelValPair{Rel: "urn:X-databox:rels:hasDatasourceid", Val: metadata.DataSourceID})
	cat.ItemMetadata = append(cat.ItemMetadata, RelValPair{Rel: "urn:X-databox:rels:hasStoreType", Val: string(metadata.StoreType)})

	if metadata.IsActuator {
		cat.ItemMetadata = append(cat.ItemMetadata, RelValPairBool{Rel: "urn:X-databox:rels:isActuator", Val: true})
	}

	if metadata.Location != "" {
		cat.ItemMetadata = append(cat.ItemMetadata, RelValPair{Rel: "urn:X-databox:rels:hasLocation", Val: metadata.Location})
	}

	if metadata.Unit != "" {
		cat.ItemMetadata = append(cat.ItemMetadata, RelValPair{Rel: "urn:X-databox:rels:hasUnit", Val: metadata.Unit})
	}

	cat.Href = endPoint + "/" + string(metadata.StoreType) + "/" + metadata.DataSourceID

	return json.Marshal(cat)

}

func (csc *CoreStoreClient) delete(path string, contentType StoreContentType) error {

	token, err := csc.Arbiter.RequestToken(csc.ZEndpoint+path, "DELETE")
	if err != nil {
		return errors.New("Error getting Arbiter Token: " + err.Error())
	}

	err = csc.ZestC.Delete(string(token), path, string(contentType))
	if err != nil {
		csc.Arbiter.InvalidateCache(csc.ZEndpoint+path, "DELETE")
		return errors.New("Error writing: " + err.Error())
	}

	return nil
}

func (csc *CoreStoreClient) read(path string, contentType StoreContentType) ([]byte, error) {

	token, err := csc.Arbiter.RequestToken(csc.ZEndpoint+path, "GET")
	if err != nil {
		return []byte(""), errors.New("Error getting Arbiter Token: " + err.Error())

	}

	resp, getErr := csc.ZestC.Get(string(token), path, string(contentType))
	if getErr != nil {
		csc.Arbiter.InvalidateCache(csc.ZEndpoint+path, "GET")
		return []byte(""), errors.New("Error getting latest data: " + getErr.Error())
	}

	return resp, nil
}

func (csc *CoreStoreClient) observe(path string, contentType StoreContentType) (<-chan ObserveResponse, error) {

	token, err := csc.Arbiter.RequestToken(csc.ZEndpoint+path, "GET")
	if err != nil {
		return nil, errors.New("Error getting Arbiter Token: " + err.Error())

	}

	payloadChan, getErr := csc.ZestC.Observe(string(token), path, string(contentType), 0)
	if getErr != nil {
		csc.Arbiter.InvalidateCache(csc.ZEndpoint+path, "GET")
		return nil, errors.New("Error observing: " + getErr.Error())
	}

	objectChan := make(chan ObserveResponse)

	go func() {
		for data := range payloadChan {
			objectChan <- csc.parseRawObserveResponse(data)
		}

		//if we get here then payloadChan has been closed so close objectChan
		close(objectChan)
	}()

	return objectChan, err
}

func (csc *CoreStoreClient) write(path string, payload []byte, contentType StoreContentType) error {

	token, err := csc.Arbiter.RequestToken(csc.ZEndpoint+path, "POST")
	if err != nil {
		return errors.New("Error getting Arbiter Token: " + err.Error())
	}

	_, err = csc.ZestC.Post(string(token), path, payload, string(contentType))
	if err != nil {
		csc.Arbiter.InvalidateCache(csc.ZEndpoint+path, "POST")
		return errors.New("Error writing: " + err.Error())
	}

	return nil
}

func (csc *CoreStoreClient) parseRawObserveResponse(data []byte) ObserveResponse {

	parts := bytes.SplitN(data, []byte(" "), 4)

	_timestamp, _ := strconv.ParseInt(string(parts[0]), 10, 64)

	parts2 := bytes.Split(parts[1], []byte("/"))

	_dataSourceID := string(parts2[2])

	_key := ""
	if len(parts2) > 3 {
		_key = string(parts2[3])
	}

	_data := parts[3]

	return ObserveResponse{_timestamp, _dataSourceID, _key, _data}
}
