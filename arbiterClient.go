package libDatabox

import (
	"bytes"
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
)

type ArbiterClient struct {
	request         *http.Client
	arbiterURL      string
	arbiterToken    string
	tokenCache      map[string][]byte
	tokenCacheMutex *sync.Mutex
}

func NewArbiterClient(arbiterTokenPath string, databoxRequest *http.Client, arbiterURL string) ArbiterClient {

	ac := ArbiterClient{
		arbiterURL:      arbiterURL,
		request:         databoxRequest,
		tokenCache:      make(map[string][]byte),
		tokenCacheMutex: &sync.Mutex{},
	}

	arbToken, err := ioutil.ReadFile(arbiterTokenPath)
	if err != nil {
		fmt.Println("Warning:: failed to read ARBITER_TOKEN using empty string")
		ac.arbiterToken = ""
	} else {
		ac.arbiterToken = b64.StdEncoding.EncodeToString([]byte(arbToken))
	}

	return ac
}

func (arb *ArbiterClient) GetRootDataSourceCatalogue() (HypercatRoot, error) {

	cat, status := arb.makeArbiterRequest("GET", "/cat", arb.arbiterURL, "/cat", "GET")
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

	req, err := http.NewRequest("POST", arb.arbiterURL+"/cm/upsert-container-info", bytes.NewBuffer(jsonPostData))
	if err != nil {
		fmt.Println("[UpdateArbiter] Error:: ", err)
		return err
	}
	req.Header.Set("X-Api-Key", arb.arbiterToken)
	req.Header.Set("Content-Type", "application/json")
	req.Close = true

	//TODO check response
	_, err = arb.request.Do(req)
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

	req, err := http.NewRequest("POST", arb.arbiterURL+"/cm/grant-container-permissions", bytes.NewBuffer(jsonPostData))
	if err != nil {
		fmt.Println("[GrantContainerPermissions] Error:: ", err)
		return err
	}
	req.Header.Set("X-Api-Key", arb.arbiterToken)
	req.Header.Set("Content-Type", "application/json")
	req.Close = true

	//TODO check response
	_, err = arb.request.Do(req)
	if err != nil {
		return err
	}

	return nil
}

func (arb *ArbiterClient) makeArbiterRequest(arbMethod string, path string, hostname string, endpoint string, method string) ([]byte, int) {

	if arb.arbiterURL == "" || arb.arbiterToken == "" {
		//Arbiter not configured
		fmt.Println("makeArbiterRequest Warning:: Arbiter not configured")
		return []byte{}, 200
	}
	var jsonStr = []byte(`{"target":"` + hostname + `","path":"` + endpoint + `","method":"` + method + `"}`)

	url := arb.arbiterURL + path

	req, err := http.NewRequest(arbMethod, url, bytes.NewBuffer(jsonStr))
	if err != nil {
		fmt.Println("makeArbiterRequest Error:: ", err)
		return []byte{}, 503
	}
	req.Header.Set("X-Api-Key", arb.arbiterToken)
	req.Header.Set("Content-Type", "application/json")
	req.Close = true
	resp, err := arb.request.Do(req)
	if err != nil {
		fmt.Println("makeArbiterRequest Error:: ", err)
		return []byte{}, 503
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)

	return body, resp.StatusCode
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
		token, status = arb.makeArbiterRequest("POST", "/token", host, u.Path, method)

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
