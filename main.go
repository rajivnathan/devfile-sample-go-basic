package main

import (
	"fmt"
	"net/http"
	"os"
)

var port = os.Getenv("PORT")

func main() {
	if port == "" {
		port = "8080"
	}
	http.HandleFunc("/", HelloServer)
	http.ListenAndServe(fmt.Sprintf("0.0.0.0:%s", port), nil)
}

func HelloServer(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path[1:]
	if path != "" {
		resp, err := http.Get("http://pagerduty.com/")
		if err != nil {
			fmt.Fprintf(w, "There was an error pinging pagerduty!")
		}
		fmt.Fprintf(w, "Resposne Status: %s", resp.Status)
	} else {
		fmt.Fprint(w, "Hello World!")
	}
}
