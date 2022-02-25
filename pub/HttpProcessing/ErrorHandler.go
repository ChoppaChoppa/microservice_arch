package HttpProcessing

import (
	"encoding/json"
	"fmt"
	"go.uber.org/zap"
	"net/http"
)

type CustomError struct {
	Message string `json:"message"`
}

func HttpError(w http.ResponseWriter, err error, msgForLogger string, msgForResponse string, code int) {
	w.Header().Set("Content-Type", "application/json")
	fmt.Println("1")
	ce := CustomError{
		Message: msgForResponse,
	}

	fmt.Println("2")
	res, errGetJson := GetJsonByte(ce)
	if errGetJson != nil {
		zap.S().Errorw("marshal", "error", errGetJson)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("error"))
		return
	}

	fmt.Println("3")
	fmt.Println(msgForLogger + ": " + err.Error())
	w.WriteHeader(code)
	w.Write(res)
}

func GetJsonByte(v interface{}) ([]byte, error) {
	usersJson, errJson := json.Marshal(v)
	if errJson != nil {
		return nil, errJson
	}

	return usersJson, nil
}