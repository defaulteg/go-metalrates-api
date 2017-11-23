package listener

import (
	"net/http"
	"log"
)

// Start listening for incoming requests
func InvokeApiServiceListener() {
	go func() {
		// Make new router
		router := NewRouter()

		//Set static resource folder
		router.PathPrefix("/").Handler(http.FileServer(http.Dir("./static/")))

		// Start listening for responses
		log.Fatal(http.ListenAndServe(":8080", router))
	}()
}




