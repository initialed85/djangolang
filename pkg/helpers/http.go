package helpers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

var unknownErrorResponse Response = Response{
	Status:  http.StatusInternalServerError,
	Success: false,
	Error:   fmt.Errorf("unknown error (HTTP %d)", http.StatusInternalServerError).Error(),
	Objects: nil,
}
var unknownErrorResponseJSON []byte

func init() {
	b, err := json.Marshal(unknownErrorResponse)
	if err != nil {
		panic(err)
	}

	unknownErrorResponseJSON = b
}

type Response struct {
	Status  int    `json:"status"`
	Success bool   `json:"success"`
	Error   string `json:"error,omitempty"`
	Objects any    `json:"objects,omitempty"`
}

func GetResponse(
	status int,
	err error,
	objects any,
) (int, Response, []byte, error) {
	if status >= http.StatusBadRequest && err == nil {
		err = fmt.Errorf("unspecified error (HTTP %d)", status)
	}

	errorMessage := ""
	if err != nil {
		errorMessage = err.Error()
	}

	response := Response{
		Status:  status,
		Success: err == nil,
		Error:   errorMessage,
		Objects: objects,
	}

	b, err := json.Marshal(response)
	if err != nil {
		response = Response{
			Status:  http.StatusInternalServerError,
			Success: false,
			Error: fmt.Errorf(
				"failed to marshal status: %d, err: %#+v, objects: %#+v to JSON: %s",
				status,
				errorMessage,
				objects,
				err,
			).Error(),
			Objects: objects,
		}

		b, err = json.Marshal(response)
		if err != nil {
			return http.StatusInternalServerError, unknownErrorResponse, unknownErrorResponseJSON, err
		}

		return http.StatusInternalServerError, response, b, err
	}

	return status, response, b, nil
}

func WriteResponse(w http.ResponseWriter, status int, b []byte) {
	w.WriteHeader(status)
	_, err := w.Write(b)
	if err != nil {
		log.Printf("warning: failed WriteResponse for status: %d: b: %s: %s", status, string(b), err.Error())
		return
	}
}

func HandleErrorResponse(w http.ResponseWriter, status int, err error) {
	status, _, b, _ := GetResponse(status, err, nil)
	WriteResponse(w, status, b)
}

func HandleObjectsResponse(w http.ResponseWriter, status int, objects any) {
	status, _, b, _ := GetResponse(status, nil, objects)
	WriteResponse(w, status, b)
}
