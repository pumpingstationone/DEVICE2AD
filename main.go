package main

import (
	"fmt"
	"log"
	"net/http"
)

func checkAccess(device, tag string) bool {
	return true
}

// We should be getting a request from a device similar to:
// 		http://localhost:8080/authcheck?device=leblond&tag=12345
func authCheckServer(w http.ResponseWriter, r *http.Request) {
	// Make sure we have both parts of the request, otherwise do
	// nothing. This is *not* the same as returning a value to
	// the caller; we don't want to give the false impression
	// that there's a problem with access when there is really a
	// problem with the way we're being called.
	device, ok := r.URL.Query()["device"]
	if !ok || len(device[0]) < 1 {
		log.Println("Url Param 'device' is missing")
		return
	}

	tag, ok := r.URL.Query()["tag"]
	if !ok || len(tag[0]) < 1 {
		log.Println("Url Param 'tag' is missing")
		return
	}

	// Okay, if we're here, we have a valid request. Let's see
	// whether they have access or not...

	log.Printf("Got a request from %s to look up access for %s", device[0], tag[0])

	// Now we're going to check whether the tag is associated with
	// having valid access to the device; because we are assuming that
	// the caller was the device itself, we don't mess around with
	// fancy responses, simply a "1" if access is to be granted, and
	// "0" if otherwise
	if checkAccess(device[0], tag[0]) {
		fmt.Fprintf(w, "%d", 1)
	} else {
		fmt.Fprintf(w, "%d", 0)
	}
}

func main() {
	http.HandleFunc("/authcheck", authCheckServer)
	http.ListenAndServe(":8080", nil)
}
