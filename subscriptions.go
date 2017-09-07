package libDatabox

import (
	"fmt"
	"net/http"
	s "strings"

	"golang.org/x/net/websocket"
)

// WSConnect connects to a target store's notification service
//
// Returns a chan []byte that recives data from every route the connecting container is subscribed to.
// Data is a JSON string.
func WSConnect(href string) (chan []byte, error) {

	out := make(chan []byte)

	storeURL, _ := GetStoreURLFromDsHref(href)

	token, err := requestToken(storeURL+"/ws", "GET")

	if err != nil {
		return nil, err
	}

	storeWsURL := s.Replace(storeURL, "https://", "wss://", 1) + "/ws"

	fmt.Println("WSConnect storeWsURL = ", storeWsURL)

	wsConfig, _ := websocket.NewConfig(storeWsURL, "http://localhost:8080")

	wsHeaders := http.Header{
		"Origin":    {"http://localhost:8080"},
		"x-api-key": {token},
	}

	//dial the ws
	wsConfig.TlsConfig = getDataboxTslConfig()
	wsConfig.Header = wsHeaders
	ws, err1 := websocket.DialConfig(wsConfig)
	if err1 != nil {
		fmt.Println("WSConnect DialConfig error:: ", err1)
		return nil, err1
	}

	go func() {
		buffer := make([]byte, 512)
		for {
			bytesRead, readErr := ws.Read(buffer)
			if readErr == nil && bytesRead > 0 {
				msg := buffer[:bytesRead]
				out <- msg
			}
			if readErr != nil {
				fmt.Println("WS Read error", readErr)
			}
		}

	}()

	return out, nil
}

//WSSubscribe Subscribes the caller to write notifications for a given route.
func WSSubscribe(href string, storeType string) (string, error) {
	dataSourceID, _ := GetDsIdFromDsHref(href)
	storeURL, _ := GetStoreURLFromDsHref(href)

	resp, err := makeStoreRequest(storeURL+"/sub/"+dataSourceID+"/"+storeType, "GET")

	return resp, err
}

//WSUnsubscribe Unsubscribes the caller to write notifications for a given route
func WSUnsubscribe(href string, storeType string) (string, error) {

	dataSourceID, _ := GetDsIdFromDsHref(href)
	storeURL, _ := GetStoreURLFromDsHref(href)

	resp, err := makeStoreRequest(storeURL+"/unsub/"+dataSourceID+"/"+storeType, "GET")

	return resp, err

}
