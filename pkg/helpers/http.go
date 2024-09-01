package helpers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
)

var cors = "*"

var unknownErrorResponse Response = Response{
	Status:  http.StatusInternalServerError,
	Success: false,
	Error:   []string{fmt.Errorf("unknown error (HTTP %d)", http.StatusInternalServerError).Error()},
	Objects: nil,
}
var unknownErrorResponseJSON []byte

func init() {
	corsOverride := strings.TrimSpace(os.Getenv("DJANGOLANG_CORS"))
	if corsOverride == "" {
		log.Printf("DJANGOLANG_CORS empty or unset; defaulted to '*'")
	}

	b, err := json.Marshal(unknownErrorResponse)
	if err != nil {
		panic(err)
	}

	unknownErrorResponseJSON = b
}

type Response struct {
	Status  int      `json:"status"`
	Success bool     `json:"success"`
	Error   []string `json:"error,omitempty"`
	Objects any      `json:"objects,omitempty"`
}

type TypedResponse[T any] struct {
	Status  int      `json:"status"`
	Success bool     `json:"success"`
	Error   []string `json:"error,omitempty"`
	Objects []*T     `json:"objects,omitempty"`
}

func GetResponse(status int, err error, objects any, prettyFormats ...bool) (int, Response, []byte, error) {
	prettyFormat := len(prettyFormats) > 1 && prettyFormats[0]

	if status >= http.StatusBadRequest && err == nil {
		err = fmt.Errorf("unspecified error (HTTP %d)", status)
	}

	errorMessage := []string{}
	if err != nil {
		errorMessage = strings.Split(err.Error(), "; ")
	}

	response := Response{
		Status:  status,
		Success: err == nil,
		Error:   errorMessage,
		Objects: objects,
	}

	var b []byte

	if status != http.StatusNoContent {
		if prettyFormat {
			b, err = json.MarshalIndent(response, "", "    ")
		} else {
			b, err = json.Marshal(response)
		}
	}

	if err != nil {
		response = Response{
			Status:  http.StatusInternalServerError,
			Success: false,
			Error: []string{
				fmt.Sprintf("failed to marshal to JSON; status: %d", status),
				fmt.Sprintf("err: %v", err),
				fmt.Sprintf("objects: %#+v", objects),
			},
			Objects: objects,
		}

		if prettyFormat {
			b, err = json.MarshalIndent(response, "", "    ")
			if err != nil {
				return http.StatusInternalServerError, unknownErrorResponse, unknownErrorResponseJSON, err
			}
		} else {
			b, err = json.Marshal(response)
			if err != nil {
				return http.StatusInternalServerError, unknownErrorResponse, unknownErrorResponseJSON, err
			}
		}

		return http.StatusInternalServerError, response, b, err
	}

	return status, response, b, nil
}

func WriteResponse(w http.ResponseWriter, status int, b []byte) {
	w.Header().Add("Access-Control-Allow-Origin", cors)
	w.Header().Add("Content-Type", "application/json")

	w.WriteHeader(status)
	_, err := w.Write(b)
	if err != nil {
		log.Printf("warning: failed WriteResponse for status: %d: b: %s: %s", status, string(b), err.Error())
		return
	}
}

func HandleErrorResponse(w http.ResponseWriter, status int, err error) {
	status, _, b, _ := GetResponse(status, err, nil, true)
	WriteResponse(w, status, b)
}

func HandleObjectsResponse(w http.ResponseWriter, status int, objects any) []byte {
	status, _, b, _ := GetResponse(status, nil, objects)
	WriteResponse(w, status, b)
	return b
}
