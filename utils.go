//
// A golang library for interfacing with Databox APIs.
//
// Install using go get github.com/me-box/lib-go-databox
//
// Examples can be found in the samples directory
//
package libDatabox

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	b64 "encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"os"
	"strconv"
	s "strings"
	"sync"
	"time"
)

var hostname = os.Getenv("DATABOX_LOCAL_NAME")
var arbiterURL = os.Getenv("DATABOX_ARBITER_ENDPOINT")
var arbiterToken string
var serverKey []byte

var databoxClient *http.Client
var databoxTlsConfig *tls.Config

func init() {

	//get the arbiterToken
	arbToken, err := ioutil.ReadFile("/run/secrets/ARBITER_TOKEN")
	if err != nil {
		fmt.Println("Warning:: failed to read ARBITER_TOKEN using empty string")
		arbiterToken = ""
	} else {
		arbiterToken = b64.StdEncoding.EncodeToString([]byte(arbToken))
	}

	//get the server key
	serverKey, err = ioutil.ReadFile("/run/secrets/ZMQ_PUBLIC_KEY")
	if err != nil {
		fmt.Println("Warning:: failed to read ZMQ_PUBLIC_KEY using default value")
		serverKey = []byte("vl6wu0A@XP?}Or/&BR#LSxn>A+}L)p44/W[wXL3<")
	}

	//setup the https root cert
	CM_HTTPS_CA_ROOT_CERT, err := ioutil.ReadFile("/run/secrets/DATABOX_ROOT_CA")
	var tr *http.Transport
	if err != nil {
		fmt.Println("Warning:: failed to read root certificate certs will not be checked.")
		tr = &http.Transport{
			Dial: (&net.Dialer{
				Timeout: 5 * time.Second,
			}).Dial,
			TLSHandshakeTimeout: 5 * time.Second,
			DisableCompression:  true,
		}

	} else {
		roots := x509.NewCertPool()
		ok := roots.AppendCertsFromPEM([]byte(CM_HTTPS_CA_ROOT_CERT))
		if !ok {
			fmt.Println("Warning:: failed to parse root certificate")
		}

		databoxTlsConfig = &tls.Config{RootCAs: roots}
		tr = &http.Transport{
			TLSClientConfig: databoxTlsConfig,
			Dial: (&net.Dialer{
				Timeout: 5 * time.Second,
			}).Dial,
			TLSHandshakeTimeout: 5 * time.Second,
			DisableCompression:  true,
		}
	}

	databoxClient = &http.Client{
		Transport: tr,
		Timeout:   time.Second * 10,
	}

}

func getDataboxTslConfig() *tls.Config {
	return databoxTlsConfig
}

func getServerKey() string {
	return string(serverKey)
}

//GetHttpsCredentials Returns a string containing the HTTPS credentials to pass to https server when offering an https server.
//These are read form /run/secrets/DATABOX.pem and are generated by the container-manger at run time.
func GetHttpsCredentials() string {
	return string("/run/secrets/DATABOX.pem")
}

func makeArbiterRequest(arbMethod string, path string, hostname string, endpoint string, method string) (string, int) {

	if arbiterURL == "" {
		//Arbiter not configured
		return "", 200
	}
	var jsonStr = []byte(`{"target":"` + hostname + `","path":"` + endpoint + `","method":"` + method + `"}`)

	fmt.Println(string(jsonStr[:]))

	url := arbiterURL + path

	req, err := http.NewRequest(arbMethod, url, bytes.NewBuffer(jsonStr))
	req.Header.Set("X-Api-Key", arbiterToken)
	req.Header.Set("Content-Type", "application/json")
	req.Close = true
	resp, err := databoxClient.Do(req)
	if err != nil {
		return err.Error(), 503
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)

	return string(body[:]), resp.StatusCode
}

var tokenCache = make(map[string]string)
var tokenCacheMutex = &sync.Mutex{}

func requestToken(href string, method string) (string, error) {

	u, err := url.Parse(href)
	if err != nil {
		return href, err
	}

	host, _, err1 := net.SplitHostPort(u.Host)
	if err != nil {
		return href, err1
	}

	routeHash := s.ToUpper(href) + method
	token, exists := tokenCache[routeHash]
	if !exists {
		var status int
		token, status = makeArbiterRequest("POST", "/token", host, u.Path, method)

		if status != 200 {
			err = errors.New(strconv.Itoa(status) + ": " + token)
			return "", err
		}
		tokenCacheMutex.Lock()
		tokenCache[routeHash] = token
		tokenCacheMutex.Unlock()
	}

	return token, err
}

func invalidateCache(href string, method string) {

	tokenCacheMutex.Lock()
	routeHash := s.ToUpper(href) + method
	delete(tokenCache, routeHash)
	tokenCacheMutex.Unlock()

}

func checkTokenCache(href string, method string) (string, error) {

	routeHash := s.ToUpper(href) + method

	_, exists := tokenCache[routeHash]
	if !exists {
		//request a token
		fmt.Println("Token not in cache requesting new one")
		newToken, err := requestToken(href, method)
		if err != nil {
			return "", err
		}
		tokenCacheMutex.Lock()
		tokenCache[routeHash] = newToken
		tokenCacheMutex.Unlock()

	}
	return tokenCache[routeHash], nil
}

type DataSourceMetadata struct {
	Description    string
	ContentType    string
	Vendor         string
	DataSourceType string
	DataSourceID   string
	StoreType      string
	IsActuator     bool
	Unit           string
	Location       string
}

type relValPair struct {
	Rel string `json:"rel"`
	Val string `json:"val"`
}

type relValPairBool struct {
	Rel string `json:"rel"`
	Val bool   `json:"val"`
}

type hypercat struct {
	ItemMetadata []interface{} `json:"item-metadata"`
	Href         string        `json:"href"`
}

//dataSourceMetadataToHypercat converts a DataSourceMetadata instance to json for registering a data source
func dataSourceMetadataToHypercat(metadata DataSourceMetadata, endPoint string) ([]byte, error) {

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

	cat.Href = endPoint + "/" + metadata.DataSourceType + "/" + metadata.DataSourceID

	return json.Marshal(cat)

}

// HypercatToDataSourceMetadata is a helper function to convert the hypercat description of a datasource to a DataSourceMetadata instance
// Also returns the store url for this data source.
func HypercatToDataSourceMetadata(hypercatDataSourceDescription string) (DataSourceMetadata, string, error) {
	dm := DataSourceMetadata{}

	hc := hypercat{}
	err := json.Unmarshal([]byte(hypercatDataSourceDescription), &hc)
	if err != nil {
		return dm, "", err
	}

	for _, pair := range hc.ItemMetadata {
		vals := pair.(map[string]interface{})
		if vals["rel"].(string) == "urn:X-hypercat:rels:hasDescription:en" {
			dm.Description = vals["val"].(string)
			continue
		}
		if vals["rel"].(string) == "urn:X-hypercat:rels:isContentType" {
			dm.ContentType = vals["val"].(string)
			continue
		}
		if vals["rel"].(string) == "urn:X-databox:rels:hasVendor" {
			dm.Vendor = vals["val"].(string)
			continue
		}
		if vals["rel"].(string) == "urn:X-databox:rels:hasType" {
			dm.DataSourceType = vals["val"].(string)
			continue
		}
		if vals["rel"].(string) == "urn:X-databox:rels:hasDatasourceid" {
			dm.DataSourceID = vals["val"].(string)
			continue
		}
		if vals["rel"].(string) == "urn:X-databox:rels:hasStoreType" {
			dm.StoreType = vals["val"].(string)
			continue
		}
		if vals["rel"].(string) == "urn:X-databox:rels:isActuator" {
			dm.IsActuator = vals["val"].(bool)
			continue
		}
		if vals["rel"].(string) == "urn:X-databox:rels:hasLocation" {
			dm.Location = vals["val"].(string)
			continue
		}
		if vals["rel"].(string) == "urn:X-databox:rels:hasUnit" {
			dm.Unit = vals["val"].(string)
			continue
		}

	}

	url, getStoreURLErr := getStoreURLFromDsHref(hc.Href)

	return dm, url, getStoreURLErr
}

// GetStoreURLFromDsHref extracts the base store url from the href provied in the hypercat descriptions.
func getStoreURLFromDsHref(href string) (string, error) {

	u, err := url.Parse(href)
	if err != nil {
		return "", err
	}

	return u.Scheme + "://" + u.Host, nil

}
