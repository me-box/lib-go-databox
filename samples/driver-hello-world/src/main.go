package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	databox "github.com/toshbrown/lib-go-databox"
)

func main() {

	//Set up the needed databox components to communicate with other parts of the databox
	var storeEndPoint = os.Getenv("DATABOX_ZMQ_ENDPOINT")
	arbiterClient, _ := databox.NewArbiterClient(databox.DefaultArbiterKeyPath, databox.DefaultStorePublicKeyPath, os.Getenv("DATABOX_ARBITER_ENDPOINT"))
	coreStoreClient := databox.NewCoreStoreClient(arbiterClient, databox.DefaultStorePublicKeyPath, storeEndPoint, false)

	//Register your data source so apps can find it
	testDataSource := databox.DataSourceMetadata{
		Description:    "Hello world test data",
		ContentType:    "application/json",
		Vendor:         "databox",
		DataSourceType: "test",
		DataSourceID:   "test",
		StoreType:      "timseries",
		IsActuator:     false,
	}
	err := coreStoreClient.RegisterDatasource(testDataSource)
	if err != nil {
		panic("Can't register data source with store: " + err.Error())
	}

	//write in some data to the timeseries blob store
	go func() {

		for {
			var data = map[string]string{"data": "Hello World " + time.Now().Format(time.RFC850) + " !"}
			res, _ := json.Marshal(data)
			coreStoreClient.TSBlobJSON.Write(testDataSource.DataSourceID, res)
			time.Sleep(1000 * time.Millisecond)
		}

	}()

	//start the https server for the driver UI
	http.HandleFunc("/ui", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "<html><body><h1>hello world! from a databox driver</h1></body></html>\n")
	})
	//The https server is setup to offer the configuration UI for your driver
	//you can use any framework you like to display the interface and parse
	//user input.
	log.Fatal(http.ListenAndServeTLS(":8080", databox.DefaultHTTPSCertPath, databox.DefaultHTTPSCertPath, nil))

}
