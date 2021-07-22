package frontend

import (
	"testing"
)

func TestReadJSONError(t *testing.T) {
	_, gotErr := ReadJSON("badfile.json")

	if gotErr == nil {
		t.Errorf("Invalid JSON file should return an error")
	}
}

var reader = NewSwaggerReader("../")

func TestReadSwagger(t *testing.T) {
	swagData, _ := reader.ReadSwagger()

	if swagData == nil {
		t.Errorf("Valid swagger file should not return an error")
	}
}