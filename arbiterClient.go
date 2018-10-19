package libDatabox

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"

	zest "github.com/me-box/goZestClient"
)

type ArbiterClient struct {
	request         *http.Client
	arbiterZMQURI   string
	ArbiterToken    string
	tokenCache      map[string][]byte
	tokenCacheMutex *sync.Mutex
	ZestC           zest.ZestClient
}

//NewArbiterClient returns an arbiter client for use by components that require conunication with the arbiter
func NewArbiterClient(arbiterTokenPath string, zmqPublicKeyPath string, arbiterZMQURI string) (*ArbiterClient, error) {

	ac := ArbiterClient{
		arbiterZMQURI:   arbiterZMQURI,
		tokenCache:      make(map[string][]byte),
		tokenCacheMutex: &sync.Mutex{},
	}

	arbToken, err := ioutil.ReadFile(arbiterTokenPath)
	if err != nil {
		fmt.Println("Warning:: failed to read ARBITER_TOKEN using default value")
		ac.ArbiterToken = "secret"
	} else {
		ac.ArbiterToken = string(arbToken)
	}

	//get the server public key
	serverKey, err := ioutil.ReadFile(zmqPublicKeyPath)
	if err != nil {
		fmt.Println("Warning:: failed to read ZMQ_PUBLIC_KEY using default value")
		serverKey = []byte("vl6wu0A@XP?}Or/&BR#LSxn>A+}L)p44/W[wXL3<")
	}

	DEndpoint := strings.Replace(arbiterZMQURI, ":4444", ":4445", 1)
	ac.ZestC, err = zest.New(arbiterZMQURI, DEndpoint, string(serverKey), false)
	if err != nil {
		return &ArbiterClient{}, errors.New("Can't connect to Arbiter on " + arbiterZMQURI)
	}

	return &ac, nil
}

// GetRootDataSourceCatalogue is used by the container manager to access the Root hypercat catalogue
func (arb *ArbiterClient) GetRootDataSourceCatalogue() (HypercatRoot, error) {

	cat, status := arb.makeArbiterGETRequest("/cat", arb.arbiterZMQURI, "/cat", "GET")
	if status != 200 {
		err := errors.New(strconv.Itoa(status) + ": " + " GET " + " /cat Failed")
		return HypercatRoot{}, err
	}

	rootCat := HypercatRoot{}

	err := json.Unmarshal(cat, &rootCat)
	if err != nil {
		fmt.Println("[GetRootDataSourceCatalogue] ", err)
	}

	return rootCat, nil
}

// RegesterDataboxComponent allows the container manager to register a new app, driver or store with the arbiter
func (arb *ArbiterClient) RegesterDataboxComponent(name string, tokenString string, databoxType DataboxType) error {

	type JsonPostData struct {
		Name string `json:"name"`
		Key  string `json:"key"`
		Type string `json:"type"`
	}

	postData := JsonPostData{
		Name: name,
		Key:  tokenString,
		Type: string(databoxType),
	}

	jsonPostData, _ := json.Marshal(postData)

	_, err := arb.ZestC.Post(arb.ArbiterToken, "/cm/upsert-container-info", jsonPostData, string(ContentTypeJSON))
	if err != nil {
		fmt.Println("[UpdateArbiter] Error:: ", err)
		return err
	}

	return nil
}

type Route struct {
	Target string `json:"target"`
	Path   string `json:"path"`
	Method string `json:"method"`
}

type ContainerPermissions struct {
	Name    string   `json:"name"`
	Route   Route    `json:"route"`
	Caveats []string `json:"caveats"`
}

// GrantContainerPermissions allows the container manager to grant permissions to an app or driver on a registered store.
func (arb *ArbiterClient) GrantContainerPermissions(permissions ContainerPermissions) error {

	if len(permissions.Caveats) == 0 {
		permissions.Caveats = []string{}
	}

	jsonPostData, _ := json.Marshal(permissions)

	_, err := arb.ZestC.Post(arb.ArbiterToken, "/cm/grant-container-permissions", jsonPostData, string(ContentTypeJSON))
	if err != nil {
		return err
	}

	return nil
}

func (arb *ArbiterClient) makeArbiterGETRequest(path string, hostname string, endpoint string, method string) ([]byte, int) {

	if arb.arbiterZMQURI == "" {
		return []byte{}, 200
	}

	resp, err := arb.ZestC.Get(arb.ArbiterToken, path, string(ContentTypeTEXT))
	if err != nil {
		fmt.Println("makeArbiterGETRequest "+path+" Error:: ", err)
		return []byte{}, 500
	}

	return resp, 200
}

func (arb *ArbiterClient) makeArbiterPostRequest(path string, hostname string, endpoint string, payload []byte) ([]byte, int) {

	if arb.arbiterZMQURI == "" {
		return nil, 200
	}

	resp, err := arb.ZestC.Post(arb.ArbiterToken, path, payload, string(ContentTypeTEXT))
	if err != nil {
		fmt.Println("makeArbiterPostRequest "+path+" Error:: ", err)
		return nil, 500
	}

	return resp, 200
}

// RequestToken is used internally to request a token from the arbiter
func (arb *ArbiterClient) RequestToken(href string, method string) ([]byte, error) {

	u, err := url.Parse(href)
	if err != nil {
		return []byte{}, err
	}

	host, _, err1 := net.SplitHostPort(u.Host)
	if err != nil {
		return []byte{}, err1
	}

	routeHash := strings.ToUpper(href) + method
	arb.tokenCacheMutex.Lock()
	token, exists := arb.tokenCache[routeHash]
	arb.tokenCacheMutex.Unlock()
	if !exists {
		var status int
		payload := []byte(`{"target":"` + host + `","path":"` + u.Path + `","method":"` + method + `","caveats":[]}`)

		token, status = arb.makeArbiterPostRequest("/token", host, u.Path, payload)
		if status != 200 {
			err = errors.New(strconv.Itoa(status) + ": " + string(token))
			return []byte{}, err
		}
		arb.tokenCacheMutex.Lock()
		arb.tokenCache[routeHash] = token
		arb.tokenCacheMutex.Unlock()
	}

	return token, err
}

// InvalidateCache can be used to remove a token from the arbiterClient cache.
// This is done automatically if the token is rejected.
func (arb *ArbiterClient) InvalidateCache(href string, method string) {

	arb.tokenCacheMutex.Lock()
	routeHash := strings.ToUpper(href) + method
	delete(arb.tokenCache, routeHash)
	arb.tokenCacheMutex.Unlock()

}

func (arb *ArbiterClient) checkTokenCache(href string, method string) ([]byte, error) {

	routeHash := strings.ToUpper(href) + method

	arb.tokenCacheMutex.Lock()
	_, exists := arb.tokenCache[routeHash]
	arb.tokenCacheMutex.Unlock()
	if !exists {
		//request a token
		newToken, err := arb.RequestToken(href, method)
		if err != nil {
			return []byte{}, err
		}
		arb.tokenCacheMutex.Lock()
		arb.tokenCache[routeHash] = newToken
		arb.tokenCacheMutex.Unlock()

	}
	return arb.tokenCache[routeHash], nil
}

func (arb *ArbiterClient) RemoveDataboxComponent() {
	//delete-continer-info
}

func (arb *ArbiterClient) GrantComponentPermission() {
	//upsert-continer-info
}

func (arb *ArbiterClient) RevokeComponentPermission() {
	//upsert-continer-info
}
