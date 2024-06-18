package common

import (
	"encoding/json"
	"io"
	"net/http"
)

func UnmarshalResponse[T any](r *http.Response) (T, error) {
	var response T

	body, err := io.ReadAll(r.Body)

	if err != nil {
		return response, err
	}

	err = json.Unmarshal(body, &response)

	return response, err
}