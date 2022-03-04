package HttpProcessing

import (
	"encoding/json"
	"go.uber.org/zap"
	"net/http"
)

type CustomError struct {
	Message string
}

func ErrorHandler(w http.ResponseWriter, err error, msgForLogger string, msgForResponse string, code int) {
	w.Header().Set("Content-Type", "application/json")
	ce := CustomError{
		Message: msgForResponse,
	}

	res, errGetJson := json.Marshal(ce)
	if errGetJson != nil {
		zap.S().Errorw("marshal", "error", errGetJson)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("error"))
		return
	}

	w.WriteHeader(code)
	w.Write(res)
}
