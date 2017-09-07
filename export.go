package libDatabox

import "os"
import "fmt"

var exportServiceURL = os.Getenv("DATABOX_EXPORT_SERVICE_ENDPOINT")

func ExportLongpoll(destination string, payload string) (string, error) {

	//TODO patload must be an escaped json string detect it it is not and error or escape it!!

	var jsonStr = `{\"id\":\"\",\"uri\":\"" + destination + "\",\"data\":" + payload + "}`

	fmt.Println("Sending ", jsonStr)

	res, err := makeStoreRequestPOST(exportServiceURL+"/lp/export", jsonStr)

	return res, err
}
