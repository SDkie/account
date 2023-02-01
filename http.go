package account

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

func encodeRequest(input interface{}) (io.Reader, error) {
	data := struct {
		Data interface{} `json:"data"`
	}{
		Data: input,
	}

	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(data)
	if err != nil {
		log.Printf("json encoding request failed with err: %v", err)
		return nil, err
	}
	return &buf, nil
}

type errorResponse struct {
	ErrorMessage string `json:"error_message"`
}

func decodeResponse(resp *http.Response, result interface{}) error {
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		var errResponse errorResponse
		err := json.NewDecoder(resp.Body).Decode(&errResponse)
		if err != nil {
			log.Printf("errorResponse decoding failed with err: %s\n", err)
			return err
		}

		err = fmt.Errorf("http request failed with err: %s and httpStatus:%s", errResponse.ErrorMessage, resp.Status)
		log.Println(err)
		return err
	}

	err := json.NewDecoder(resp.Body).Decode(result)
	if err != nil {
		log.Printf("decoding response body failed with err: %s\n", err)
		return err
	}

	return nil
}
