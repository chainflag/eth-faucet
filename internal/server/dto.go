package server

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/chainflag/eth-faucet/internal/chain"
	"github.com/ethereum/go-ethereum/common"
)

type claimRequest struct {
	Address string `json:"address"`
}

type claimResponse struct {
	Message string `json:"msg"`
}

type infoResponse struct {
	Account         string `json:"account"`
	Network         string `json:"network"`
	Payout          string `json:"payout"`
	Symbol          string `json:"symbol"`
	HcaptchaSiteKey string `json:"hcaptcha_sitekey,omitempty"`
}

type malformedRequest struct {
	status  int
	message string
}

func (mr *malformedRequest) Error() string {
	return mr.message
}

const maxBodySize = 1024

func decodeJSONBody(r *http.Request, dst interface{}) error {
	if r.Body == nil {
		return &malformedRequest{status: http.StatusBadRequest, message: "Request body is empty"}
	}
	defer r.Body.Close()

	body, err := io.ReadAll(io.LimitReader(r.Body, maxBodySize+1))
	if err != nil {
		return &malformedRequest{status: http.StatusBadRequest, message: "Unable to read request body"}
	}
	if len(body) > maxBodySize {
		return &malformedRequest{status: http.StatusRequestEntityTooLarge, message: fmt.Sprintf("Request body must not be larger than %d bytes", maxBodySize)}
	}

	dec := json.NewDecoder(bytes.NewReader(body))
	dec.DisallowUnknownFields()

	if err := dec.Decode(dst); err != nil {
		var syntaxError *json.SyntaxError
		var unmarshalTypeError *json.UnmarshalTypeError

		switch {
		case errors.As(err, &syntaxError):
			msg := fmt.Sprintf("Request body contains badly-formed JSON (at position %d)", syntaxError.Offset)
			return &malformedRequest{status: http.StatusBadRequest, message: msg}
		case errors.Is(err, io.ErrUnexpectedEOF):
			msg := "Request body contains badly-formed JSON"
			return &malformedRequest{status: http.StatusBadRequest, message: msg}
		case errors.As(err, &unmarshalTypeError):
			msg := fmt.Sprintf("Request body contains an invalid value for the %q field (at position %d)", unmarshalTypeError.Field, unmarshalTypeError.Offset)
			return &malformedRequest{status: http.StatusBadRequest, message: msg}
		case strings.HasPrefix(err.Error(), "json: unknown field "):
			fieldName := strings.TrimPrefix(err.Error(), "json: unknown field ")
			msg := fmt.Sprintf("Request body contains unknown field %s", fieldName)
			return &malformedRequest{status: http.StatusBadRequest, message: msg}
		case errors.Is(err, io.EOF):
			msg := "Request body must not be empty"
			return &malformedRequest{status: http.StatusBadRequest, message: msg}
		default:
			return err
		}
	}

	return nil
}

func readAddress(r *http.Request) (string, error) {
	var claimReq claimRequest
	if err := decodeJSONBody(r, &claimReq); err != nil {
		return "", err
	}
	if !chain.IsValidAddress(claimReq.Address, false) {
		return "", &malformedRequest{status: http.StatusBadRequest, message: "invalid address"}
	}

	return common.HexToAddress(claimReq.Address).Hex(), nil
}

func renderJSON(w http.ResponseWriter, v interface{}, code int) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	return json.NewEncoder(w).Encode(v)
}
