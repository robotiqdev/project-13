package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/robotiqdev/project-13/internal/router"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	addr := fmt.Sprintf(":%s", port)

	r := router.New()
	if err := http.ListenAndServe(addr, r); err != nil {
		fmt.Fprintf(os.Stderr, "server error: %v\n", err)
		os.Exit(1)
	}
}
