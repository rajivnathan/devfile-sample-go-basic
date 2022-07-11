package main

import (
	"bytes"
	"context"
	"crypto/tls"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"time"
)

const namespace = "rajivtest"
const svcName = "devfile-sample-python-basic-git"
const endpoint = "rajiv"

var port = flag.Int("p", 8080, "server port")
var url = fmt.Sprintf("http://%s.%s.svc.cluster.local:8080/%s", svcName, namespace, endpoint)

var httpClient = newHTTPClient()

func main() {
	fmt.Printf("url: %s", url)
	flag.Parse()
	http.HandleFunc("/", HelloServer)
	http.HandleFunc("/rajiv", HelloRajiv)
	http.ListenAndServe(fmt.Sprintf("0.0.0.0:%d", *port), nil)
}

func HelloServer(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, %s!", r.URL.Path[1:])
}

func HelloRajiv(w http.ResponseWriter, r *http.Request) {
	result := doRequest(httpClient)
	fmt.Println(result)
	fmt.Fprint(w, result)
}

func doRequest(httpClient *http.Client) string {

	// url := fmt.Sprintf("http://%s.che-db-cleaner/%s", namespace, userID)

	// create request
	req, err := http.NewRequestWithContext(context.TODO(), http.MethodGet, url, nil)
	if err != nil {
		return err.Error()
	}

	// if queryParams != nil {
	// 	req.URL.RawQuery = queryParams.Encode()
	// }

	// do the request
	res, err := httpClient.Do(req)
	if err != nil {
		return err.Error()
	}

	defer closeResponse(res)
	resBody, readError := readBody(res.Body)
	if readError != nil {
		fmt.Println("error while reading body of the response")
		return err.Error()
	}
	return fmt.Sprintf("Response status: '%s' Body: '%s'", res.Status, resBody)
}

// newHTTPClient returns a new HTTP client with some timeout and TLS values configured
func newHTTPClient() *http.Client {
	var netTransport = &http.Transport{
		Dial: (&net.Dialer{
			Timeout: 5 * time.Second,
		}).Dial,
		TLSHandshakeTimeout: 5 * time.Second,
		TLSClientConfig:     &tls.Config{InsecureSkipVerify: true}, // nolint:gosec
	}
	var httpClient = &http.Client{
		Timeout:   time.Second * 10,
		Transport: netTransport,
	}
	return httpClient
}

// readBody reads body from a ReadCloser and returns it as a string
func readBody(body io.ReadCloser) (string, error) {
	buf := new(bytes.Buffer)
	_, err := buf.ReadFrom(body)
	return buf.String(), err
}

// closeResponse reads the body and close the response. To be used to prevent file descriptor leaks.
func closeResponse(res *http.Response) {
	if res != nil {
		io.Copy(ioutil.Discard, res.Body) //nolint: errcheck
		defer res.Body.Close()
	}
}
