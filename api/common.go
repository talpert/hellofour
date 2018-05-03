package api

import (
	"net/http"

	"encoding/json"
	"fmt"
	"io"

	"github.com/InVisionApp/rye"
	"github.com/sirupsen/logrus"
)

func decodeJSONInput(body io.ReadCloser, target interface{}, l *logrus.Entry) *rye.Response {
	if err := json.NewDecoder(body).Decode(target); err != nil {
		l.Error(err)
		return &rye.Response{
			Err:        fmt.Errorf("Unable to decode request JSON: %v", err.Error()),
			StatusCode: http.StatusInternalServerError,
		}
	}

	return nil
}

func respondAsJSON(rw http.ResponseWriter, code int, body interface{}, l *logrus.Entry) *rye.Response {
	respJSON, err := json.Marshal(body)
	if err != nil {
		l.Errorf("failed to marshal response: %v", body)
		return &rye.Response{
			Err:        fmt.Errorf("Unable to generate response JSON: %v", err),
			StatusCode: http.StatusInternalServerError,
		}
	}

	rye.WriteJSONResponse(rw, code, respJSON)

	return nil
}
