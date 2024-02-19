package main

import (
	"fmt"
	"net/http"
	"os"
	"time"
)

func main() {
	kvHandler := &KVServer{
		kvStore: NewKeyValueStore(),
	}

	server := &http.Server{
		Addr:         fmt.Sprintf(":%s", os.Getenv("PORT")),
		Handler:      kvHandler,
		ReadTimeout:  time.Second * 5,
		WriteTimeout: time.Second * 5,
	}

	err := server.ListenAndServe()
	if err != nil {
		panic(err)
	}
}
