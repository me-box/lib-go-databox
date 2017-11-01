package libDatabox

import (
	"errors"
	"io/ioutil"
	"strings"

	zest "github.com/toshbrown/goZestClient"
)

type KeyValueClient struct {
	zestC     zest.ZestClient
	zEndpoint string
	dEndpoint string
}

func NewKeyValueClient(reqEndpoint string, enableLogging bool) (KeyValueClient, error) {

	serverKey, err := ioutil.ReadFile("/run/secrets/ZMQ_PUBLIC_KEY")
	if err != nil {
		return KeyValueClient{}, err
	}

	kvc := KeyValueClient{}
	kvc.zEndpoint = reqEndpoint
	kvc.dEndpoint = strings.Replace(reqEndpoint, ":5555", ":5556", 1)
	kvc.zestC = zest.New(kvc.zEndpoint, kvc.dEndpoint, string(serverKey), enableLogging)

	return kvc, nil
}

func (kvc KeyValueClient) Write(path string, payload string) error {

	token, err := requestToken(kvc.zEndpoint+path, "POST")
	if err != nil {
		return err
	}

	err = kvc.zestC.Post(token, path, payload)
	if err != nil {
		return errors.New("Error posting data: " + err.Error())
	}

	return nil

}

func (kvc KeyValueClient) Read(path string) (string, error) {

	token, err := requestToken(kvc.zEndpoint+path, "GET")
	if err != nil {
		return "", err
	}

	resp, getErr := kvc.zestC.Get(token, path)
	if getErr != nil {
		return "", errors.New("Error posting data: " + err.Error())
	}

	return resp, nil

}
