package errors

import (
	"encoding/json"
	"net/http"

	"github.com/robertokbr/blinkchat/src/domain/logger"
)

type API struct {
	Code    int    `json:"code"`
	Err     string `json:"error"`
	Message string `json:"message"`
}

func (e *API) Error() string {
	return e.Message
}

func NewAPI(code int, err string, message string) *API {
	return &API{
		Code:    code,
		Err:     err,
		Message: message,
	}
}

func NotFound(message string) *API {
	return NewAPI(http.StatusNotFound, "Not Found", message)
}

func BadRequest(message string) *API {
	return NewAPI(http.StatusBadRequest, "Bad Request", message)
}

func InternalServerError(message string) *API {
	return NewAPI(http.StatusInternalServerError, "Internal Server Error", message)
}

func Unauthorized(message string) *API {
	return NewAPI(http.StatusUnauthorized, "Unauthorized", message)
}

func Conflict(message string) *API {
	return NewAPI(http.StatusConflict, "Conflict", message)
}

func HandleError(w http.ResponseWriter, r *http.Request, err error) {
	if err != nil {
		apiError, ok := err.(*API)

		if !ok {
			apiError = InternalServerError(err.Error())
			logger.Errorf("%s: %s", r.URL.Path, err.Error())
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(apiError.Code)
		json.NewEncoder(w).Encode(apiError)
	}
}
