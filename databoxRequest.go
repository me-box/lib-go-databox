package libDatabox

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"time"
)

func NewDataboxHTTPsAPI() *http.Client {
	//pass the default databox paths
	return NewDataboxHTTPsAPIWithPaths(DefaultHTTPSRootCertPath)
}

func NewDataboxHTTPsAPIWithPaths(cmRootCaPath string) *http.Client {

	//setup the https root cert
	CM_HTTPS_CA_ROOT_CERT, err := ioutil.ReadFile(cmRootCaPath)

	var tr *http.Transport
	if err != nil {
		fmt.Println("Warning:: failed to read root certificate certs will not be checked.")
		tr = &http.Transport{
			Dial: (&net.Dialer{
				Timeout: 5 * time.Second,
			}).Dial,
			TLSHandshakeTimeout: 5 * time.Second,
			DisableCompression:  true,
			MaxIdleConns: 100,
			MaxIdleConnsPerHost:10,
			MaxConnsPerHost:10,
			IdleConnTimeout:5,
		}

	} else {
		roots := x509.NewCertPool()
		ok := roots.AppendCertsFromPEM([]byte(CM_HTTPS_CA_ROOT_CERT))
		if !ok {
			fmt.Println("Warning:: failed to parse root certificate")
		}

		databoxTlsConfig := &tls.Config{RootCAs: roots}
		tr = &http.Transport{
			TLSClientConfig: databoxTlsConfig,
			Dial: (&net.Dialer{
				Timeout: 5 * time.Second,
			}).Dial,
			TLSHandshakeTimeout: 5 * time.Second,
			DisableCompression:  true,
		}
	}

	databoxClient := &http.Client{
		Transport: tr,
		Timeout:   time.Second * 10,
	}

	return databoxClient
}
