# lib-go-databox

A Golang library for interfacing with Databox APIs.

see https://godoc.org/github.com/me-box/lib-go-databox for full documtatiosn


# Example

Set up the store and arbiter using setupTest.sh script

```go
package main

import (
	"fmt"
	libDatabox "github.com/me-box/lib-go-databox"
)

func main () {

    //Create a new client in testing mode outside databox
    const testArbiterEndpoint = "tcp://127.0.0.1:4444"
    const testStoreEndpoint = "tcp://127.0.0.1:5555"
    ac, _ := libDatabox.NewArbiterClient("./", "./", testArbiterEndpoint)
    storeClient := libDatabox.NewCoreStoreClient(ac, "./", DataboxStoreEndpoint, false)


    //write some data
    jsonData := `{"data":"This is a test"}`
	err := storeClient.TSBlobJSON.Write("testdata1", []byte(jsonData))
	if err != nil {
		libDatabox.Err("Error Write Datasource " + err.Error())
	}

    //Read some data
    jsonData, err := storeClient.TSBlobJSON.Latest("testdata1")
    if err != nil {
        libDatabox.Err("Error Write Datasource " + err.Error())
    }
    fmt.Println(jsonData)

}
```

More examples can be found in the [databox-quickstart guide](https://github.com/me-box/databox-quickstart)

## Development of databox was supported by the following funding
```
EP/N028260/1, Databox: Privacy-Aware Infrastructure for Managing Personal Data
EP/N028260/2, Databox: Privacy-Aware Infrastructure for Managing Personal Data
EP/N014243/1, Future Everyday Interaction with the Autonomous Internet of Things
EP/M001636/1, Privacy-by-Design: Building Accountability into the Internet of Things EP/M02315X/1, From Human Data to Personal Experience
```
