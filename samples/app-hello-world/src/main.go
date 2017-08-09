package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	databox "github.com/me-box/lib-go-databox"
)

var dataSourceTest = databox.JsonUnmarshal(os.Getenv("DATASOURCE_test"))
var storeURL = databox.GetStoreURLFromDsHref(dataSourceTest["href"].(string))

func main() {

	fmt.Printf(storeURL + "\n")

	databox.WaitForStoreStatus(storeURL)

	//start the https server for the app UI
	http.HandleFunc("/ui", func(w http.ResponseWriter, r *http.Request) {

		data, err := databox.StoreJSONGetlatest(dataSourceTest["href"].(string))

		if err != nil {
			fmt.Fprintf(w, "<html><body><h1>hello world! from a databox app</h1><p>error:: "+err.Error()+"</p></body></html>\n")
		}

		fmt.Fprintf(w, "<html><body><h1>hello world! from a databox app</h1><p>Latest from the driver"+data+"</p></body></html>\n")
	})
	//The https server is setup to offer the configuration UI for your app
	//you can use any framework you like to display the interface and parse
	//user input.

	log.Fatal(http.ListenAndServeTLS(":8080", databox.GetHttpsCredentials(), databox.GetHttpsCredentials(), nil))
}
