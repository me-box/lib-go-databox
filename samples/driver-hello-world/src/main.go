package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	databox "github.com/me-box/lib-go-databox"
)

var dataStoreHref = os.Getenv("DATABOX_STORE_ENDPOINT")

func main() {

	//Wait for your store to be created. Then you can read and write values into it.
	databox.WaitForStoreStatus(dataStoreHref)

	//Register your data source so apps can find it
	_, err := databox.RegisterDatasource(dataStoreHref, databox.StoreMetadata{
		Description:    "Hello world test data",
		ContentType:    "application/json",
		Vendor:         "databox",
		DataSourceType: "test",
		DataSourceID:   "test",
		StoreType:      "ts",
		IsActuator:     false,
	})

	if err != nil {
		panic(err)
	}

	var dataSourceHref = dataStoreHref + "/test"

	//write in some data
	go func() {

		for {
			var data = map[string]string{"data": "Hello World " + time.Now().Format(time.RFC850) + " !"}
			res, _ := json.Marshal(data)
			databox.StoreJSONWriteTS(dataSourceHref, string(res[:]))
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
	log.Fatal(http.ListenAndServeTLS(":8080", databox.GetHttpsCredentials(), databox.GetHttpsCredentials(), nil))

}
