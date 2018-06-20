package libDatabox

import (
	b64 "encoding/base64"
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
	arbiterToken    string
	tokenCache      map[string][]byte
	tokenCacheMutex *sync.Mutex
	ZestC           zest.ZestClient
}

func NewArbiterClient(arbiterTokenPath string, zmqPublicKeyPath string, arbiterZMQURI string) (*ArbiterClient, error) {

	ac := ArbiterClient{
		arbiterZMQURI:   arbiterZMQURI,
		tokenCache:      make(map[string][]byte),
		tokenCacheMutex: &sync.Mutex{},
	}

	arbToken, err := ioutil.ReadFile(arbiterTokenPath)
	if err != nil {
		fmt.Println("Warning:: failed to read ARBITER_TOKEN using default value")
		ac.arbiterToken = ""
	} else {
		ac.arbiterToken = b64.StdEncoding.EncodeToString([]byte(arbToken))
	}

	//get the server public key
	serverKey, err := ioutil.ReadFile(zmqPublicKeyPath)
	if err != nil {
		fmt.Println("Warning:: failed to read ZMQ_PUBLIC_KEY using default value")
		serverKey = []byte("vl6wu0A@XP?}Or/&BR#LSxn>A+}L)p44/W[wXL3<")
	}

	DEndpoint := strings.Replace(arbiterZMQURI, ":4444", ":4445", 1)
	ac.ZestC, err = zest.New(arbiterZMQURI, DEndpoint, string(serverKey), true)
	if err != nil {
		return &ArbiterClient{}, errors.New("Can't connect to Arbiter on " + arbiterZMQURI)
	}

	return &ac, nil
}

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

func (arb *ArbiterClient) UpdateArbiter(name string, tokenString string, databoxType DataboxType) error {

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

	_, err := arb.ZestC.Post(arb.arbiterToken, "/cm/upsert-container-info", jsonPostData, string(ContentTypeJSON))
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

func (arb *ArbiterClient) GrantContainerPermissions(permissions ContainerPermissions) error {

	if len(permissions.Caveats) == 0 {
		permissions.Caveats = nil
	}

	jsonPostData, _ := json.Marshal(permissions)

	_, err := arb.ZestC.Post(arb.arbiterToken, "/cm/grant-container-permissions", jsonPostData, string(ContentTypeJSON))
	if err != nil {
		return err
	}

	return nil
}

func (arb *ArbiterClient) makeArbiterGETRequest(path string, hostname string, endpoint string, method string) ([]byte, int) {

	if arb.arbiterZMQURI == "" || arb.arbiterToken == "" {
		return []byte{}, 200
	}

	resp, err := arb.ZestC.Get(arb.arbiterToken, path, string(ContentTypeTEXT))
	if err != nil {
		fmt.Println("makeArbiterGETRequest Error:: ", err)
		return []byte{}, 500
	}

	return resp, 200
}

func (arb *ArbiterClient) makeArbiterPostRequest(path string, hostname string, endpoint string, payload []byte) ([]byte, int) {

	if arb.arbiterZMQURI == "" || arb.arbiterToken == "" {
		return nil, 200
	}

	resp, err := arb.ZestC.Post(arb.arbiterToken, path, payload, string(ContentTypeTEXT))
	if err != nil {
		fmt.Println("makeArbiterPostRequest Error:: ", err)
		return nil, 500
	}

	return resp, 200
}

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
		payload := []byte(`{"target":"` + host + `","path":"` + u.Path + `","method":"` + method + `"}`)

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

func (arb *ArbiterClient) RegesterDataboxComponent(componentName, componenttype string) {
	//upsert-continer-info
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
