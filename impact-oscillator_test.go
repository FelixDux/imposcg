package main

import (
    "fmt"
    "net/http"
    "net/http/httptest"
    "testing"
)
func TestSwaggerRoute(t *testing.T) {

    ts := httptest.NewServer(setupServer())
    // Shut down the server and block until all requests have gone through
    defer ts.Close()

	resp, err := http.Get(fmt.Sprintf("%s/swagger/index.html", ts.URL))

    if err != nil {
        t.Fatalf("Expected no error, got %v", err)
    }

    if resp.StatusCode != 200 {
        t.Fatalf("Expected status code 200, got %v", resp.StatusCode)
    }

    val, ok := resp.Header["Content-Type"]

    // Assert that the "content-type" header is actually set
    if !ok {
        t.Fatalf("Expected Content-Type header to be set")
    }

    // Assert that it was set as expected
    if val[0] != "text/html; charset=utf-8" {
        t.Fatalf("Expected \"text/html; charset=utf-8\", got %s", val[0])
    }
}
