package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	//databox "github.com/me-box/lib-go-databox"
	databox "github.com/toshbrown/lib-go-databox"
)

func main() {

	//Set up the needed databox components to communicate with other parts of the databox
	var dataSourceTest, dataSourceStoreURL, _ = databox.HypercatToDataSourceMetadata(os.Getenv("DATASOURCE_test"))
	arbiterClient, _ := databox.NewArbiterClient(databox.DefaultArbiterKeyPath, databox.DefaultStorePublicKeyPath, os.Getenv("DATABOX_ARBITER_ENDPOINT"))
	coreStoreClient := databox.NewCoreStoreClient(arbiterClient, databox.DefaultStorePublicKeyPath, dataSourceStoreURL, false)

	//start the https server for the app UI
	http.HandleFunc("/ui", func(w http.ResponseWriter, r *http.Request) {

		//Read data from the datasource
		data, err := coreStoreClient.TSBlobJSON.Latest(dataSourceTest.DataSourceID)

		if err != nil {
			fmt.Fprintf(w, "<html><body><h1>hello world! from a databox app</h1><p>error:: "+err.Error()+"</p></body></html>\n")
		}

		fmt.Fprintf(w, "<html><body><h1>hello world! from a databox app</h1><p>Latest from the driver"+string(data)+"</p></body></html>\n")
	})
	//The https server is setup to offer the configuration UI for your app
	//you can use any framework you like to display the interface and parse
	//user input.

	log.Fatal(http.ListenAndServeTLS(":8080", databox.DefaultHTTPSCertPath, databox.DefaultHTTPSCertPath, nil))
}
