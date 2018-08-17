package libDatabox

import (
	"encoding/json"
	"os"
	"strconv"
	s "strings"
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
var Arbiter *ArbiterClient

//a unique ID per test run so data will not collide
var dsID string

func Setup() {

	var err error
	hostname, _ := os.Hostname()

	Arbiter, err := NewArbiterClient("", "", "tcp://127.0.0.1:4444")
	if err != nil {
		panic("Cant connect to Zest server. Did you start one? " + err.Error())
	}

	StoreClient = NewCoreStoreClient(Arbiter, "", "tcp://127.0.0.1:5555", false)
	if err != nil {
		panic("Cant connect to Zest server. Did you start one? " + err.Error())
	}

	dsID = "test" + strconv.Itoa(int(time.Now().UnixNano()/int64(time.Millisecond)))

	Arbiter.RegesterDataboxComponent(hostname, "secret", DataboxTypeApp)

	Arbiter.GrantContainerPermissions(ContainerPermissions{
		Name: hostname,
		Route: Route{
			Target: "127.0.0.1",
			Path:   "/*",
			Method: "POST",
		},
	})

	Arbiter.GrantContainerPermissions(ContainerPermissions{
		Name: hostname,
		Route: Route{
			Target: "127.0.0.1",
			Path:   "/*",
			Method: "GET",
		},
	})
}

func Teardown() {
	//todo
}

func TestRegisterDatasource(t *testing.T) {

	dsmd := DataSourceMetadata{
		DataSourceID:   dsID,
		Vendor:         "testing",
		ContentType:    "application/json",
		StoreType:      "ts",
		Description:    "A test DS",
		DataSourceType: "test",
	}

	err := StoreClient.RegisterDatasource(dsmd)
	if err != nil {
		t.Errorf("RegisterDatasource failed expected err to be nil got %s", err.Error())
	}

	rootCat, getErr := StoreClient.GetStoreDataSourceCatalogue("tcp://127.0.0.1:5555")
	if getErr != nil {
		t.Errorf("GetDatasourceCatalogue failed expected err to be nil got %s", getErr.Error())
	}
	catByteArray, _ := json.Marshal(rootCat)

	dsmdByteArray, _ := StoreClient.dataSourceMetadataToHypercat(dsmd, "tcp://127.0.0.1:5555")

	cont := s.Contains(string(catByteArray), string(dsmdByteArray))
	if cont != true {
		t.Errorf("GetDatasourceCatalogue Error '%s' does not contain  %s", string(catByteArray), string(dsmdByteArray))
	}
}
