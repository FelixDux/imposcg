package frontend

import (
	"encoding/json"
    "io/ioutil"
    "os"
	"fmt"
)

func ReadJSON(jsonPath string) (*map[string]interface{}, error) {

    jsonFile, err := os.Open(jsonPath)
	
    if err != nil {
        return nil, err
    }
	
    defer jsonFile.Close()

    byteValue, _ := ioutil.ReadAll(jsonFile)

    var result map[string]interface{}
    json.Unmarshal([]byte(byteValue), &result)

	return &result, nil
}

type SwaggerReader struct {
	Prefix string
}

func NewSwaggerReader(prefix string) *SwaggerReader {
	return &SwaggerReader{Prefix: prefix}
}

func (reader SwaggerReader) SwaggerPath() string {
	return fmt.Sprintf("%sdocs/swagger.json", reader.Prefix)
}

func (reader SwaggerReader) ReadSwagger() (*map[string]interface{}, error) {
	return ReadJSON(reader.SwaggerPath())
}