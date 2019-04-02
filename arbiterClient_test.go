package libDatabox

import (
	"testing"
)

func TestRequestToken(t *testing.T) {

	token, err := Arbiter.RequestToken("tcp://127.0.0.1:5555/ts/test", "POST", ``)

	if err != nil {
		t.Errorf("Call to RequestToken failed with error %s", err.Error())
	}

	if len(token) == 0 {
		t.Errorf("Token to short")
	}

}

func TestRequestDeligatedToken(t *testing.T) {

	token, err := Arbiter.RequestDeligatedToken("core-logger", "tcp://127.0.0.1:5555/ts/test", "POST", ``)

	if err != nil {
		t.Errorf("Call to RequestToken failed with error %s", err.Error())
	}

	if len(token) == 0 {
		t.Errorf("Token to short")
	}

}
