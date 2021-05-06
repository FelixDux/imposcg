package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
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

func TestIterationRoutes(t *testing.T) {

    ts := httptest.NewServer(setupServer())
    // Shut down the server and block until all requests have gone through
    defer ts.Close()

    IterationTests := [] struct {
        endpoint string
        contentType string
    }{
        {"data", "application/json; charset=utf-8"},
        {"image", "image/png"},
    }

	for _, data := range IterationTests {
            resp, err := http.PostForm(fmt.Sprintf("%s/api/iteration/%s", ts.URL, data.endpoint), 
            url.Values{"frequency": {"2.8"}, "offset": {"0"}, "r": {"0.8"}, "maxPeriods": {"100"}, "phi": {"0"}, "v": {"0"}, "numIterations": {"100"}})

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
        if val[0] != data.contentType {
            t.Fatalf("Expected \"%s\", got %s", data.contentType, val[0])
        }
    }

}