package main

import (
	"fmt"
	"net/http"
)

func main() {
	fmt.Println("Starting hello-world server...")

	http.Handle("/", helloHandler())
	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic(err)
	}

}
