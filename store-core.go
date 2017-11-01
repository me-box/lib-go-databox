package libDatabox

import (
	"errors"

	zest "github.com/toshbrown/goZestClient"
)

type KeyValueClient struct {
	zestC     zest.ZestClient
	zEndpoint string
}

func NewKeyValueClient(ReqEndpoint string, DealerEndpoint string, ServerKey string, enableLogging bool) KeyValueClient {

	kvc := KeyValueClient{}
	kvc.zEndpoint = ReqEndpoint
	kvc.zestC = zest.New(ReqEndpoint, DealerEndpoint, ServerKey, enableLogging)

	return kvc
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
