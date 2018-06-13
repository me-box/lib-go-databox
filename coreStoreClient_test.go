package libDatabox

import (
	"net/http"
	"os"
	"strconv"
	"testing"
	"time"
)

func TestMain(m *testing.M) {
	Setup()
	retCode := m.Run()
	Teardown()
	os.Exit(retCode)
}

var StoreClient *CoreStoreClient

//a unique ID per test run so data will not collide
var dsID string

func Setup() {

	dr := &http.Client{}
	ac := NewArbiterClient("", dr, "https://arbiter:8080/")

	var err error
	StoreClient = NewCoreStoreClient(dr, &ac, "", "tcp://127.0.0.1:5555", false)
	if err != nil {
		panic("Cant connect to Zest server. Did you start one? " + err.Error())
	}

	dsID = "test" + strconv.Itoa(int(time.Now().UnixNano()/int64(time.Millisecond)))
}

func Teardown() {
	//todo
}
