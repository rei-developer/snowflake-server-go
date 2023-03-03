package main

import (
	"fmt"
	"net/http"
)

func main() {
	srv := newServer(":8080")
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		fmt.Printf("failed to start server: %s\n", err)
	}
}
