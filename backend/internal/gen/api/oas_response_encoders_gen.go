// Code generated by ogen, DO NOT EDIT.

package api

import (
	"net/http"

	"github.com/go-faster/errors"
	"github.com/go-faster/jx"

	ht "github.com/ogen-go/ogen/http"
)

func encodeGetLogsForMatchResponse(response GetLogsForMatchRes, w http.ResponseWriter) error {
	switch response := response.(type) {
	case *GetLogsForMatchOK:
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(200)

		e := new(jx.Encoder)
		response.Encode(e)
		if _, err := e.WriteTo(w); err != nil {
			return errors.Wrap(err, "write")
		}

		return nil

	case *ErrorStatusCode:
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		code := response.StatusCode
		if code == 0 {
			// Set default status code.
			code = http.StatusOK
		}
		w.WriteHeader(code)

		e := new(jx.Encoder)
		response.Response.Encode(e)
		if _, err := e.WriteTo(w); err != nil {
			return errors.Wrap(err, "write")
		}

		if code >= http.StatusInternalServerError {
			return errors.Wrapf(ht.ErrInternalServerErrorResponse, "code: %d, message: %s", code, http.StatusText(code))
		}
		return nil

	default:
		return errors.Errorf("unexpected response type: %T", response)
	}
}

func encodeGetMatchForLogResponse(response GetMatchForLogRes, w http.ResponseWriter) error {
	switch response := response.(type) {
	case *GetMatchForLogOK:
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(200)

		e := new(jx.Encoder)
		response.Encode(e)
		if _, err := e.WriteTo(w); err != nil {
			return errors.Wrap(err, "write")
		}

		return nil

	case *ErrorStatusCode:
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		code := response.StatusCode
		if code == 0 {
			// Set default status code.
			code = http.StatusOK
		}
		w.WriteHeader(code)

		e := new(jx.Encoder)
		response.Response.Encode(e)
		if _, err := e.WriteTo(w); err != nil {
			return errors.Wrap(err, "write")
		}

		if code >= http.StatusInternalServerError {
			return errors.Wrapf(ht.ErrInternalServerErrorResponse, "code: %d, message: %s", code, http.StatusText(code))
		}
		return nil

	default:
		return errors.Errorf("unexpected response type: %T", response)
	}
}

func encodeGetPlayersResponse(response *GetPlayersOK, w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(200)

	e := new(jx.Encoder)
	response.Encode(e)
	if _, err := e.WriteTo(w); err != nil {
		return errors.Wrap(err, "write")
	}

	return nil
}

func encodeGetTeamResponse(response *GetTeamOK, w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(200)

	e := new(jx.Encoder)
	response.Encode(e)
	if _, err := e.WriteTo(w); err != nil {
		return errors.Wrap(err, "write")
	}

	return nil
}

func encodeErrorResponse(response *ErrorStatusCode, w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	code := response.StatusCode
	if code == 0 {
		// Set default status code.
		code = http.StatusOK
	}
	w.WriteHeader(code)

	e := new(jx.Encoder)
	response.Response.Encode(e)
	if _, err := e.WriteTo(w); err != nil {
		return errors.Wrap(err, "write")
	}

	if code >= http.StatusInternalServerError {
		return errors.Wrapf(ht.ErrInternalServerErrorResponse, "code: %d, message: %s", code, http.StatusText(code))
	}
	return nil

}
