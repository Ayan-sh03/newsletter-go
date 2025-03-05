package utils

import (
	"encoding/json"
	"net/http"
)

func ParseRequestBody(r *http.Request, model interface{}) error {

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&model); err != nil {
		return err
	}
	return nil
}
