// Code generated by ogen, DO NOT EDIT.

package api

import (
	"net/http"
	"net/url"

	"github.com/go-faster/errors"

	"github.com/ogen-go/ogen/conv"
	"github.com/ogen-go/ogen/middleware"
	"github.com/ogen-go/ogen/ogenerrors"
	"github.com/ogen-go/ogen/uri"
	"github.com/ogen-go/ogen/validate"
)

// GetLogsForMatchParams is parameters of GetLogsForMatch operation.
type GetLogsForMatchParams struct {
	// ETF2L match ID.
	MatchID int
}

func unpackGetLogsForMatchParams(packed middleware.Parameters) (params GetLogsForMatchParams) {
	{
		key := middleware.ParameterKey{
			Name: "match_id",
			In:   "path",
		}
		params.MatchID = packed[key].(int)
	}
	return params
}

func decodeGetLogsForMatchParams(args [1]string, argsEscaped bool, r *http.Request) (params GetLogsForMatchParams, _ error) {
	// Decode path: match_id.
	if err := func() error {
		param := args[0]
		if argsEscaped {
			unescaped, err := url.PathUnescape(args[0])
			if err != nil {
				return errors.Wrap(err, "unescape path")
			}
			param = unescaped
		}
		if len(param) > 0 {
			d := uri.NewPathDecoder(uri.PathDecoderConfig{
				Param:   "match_id",
				Value:   param,
				Style:   uri.PathStyleSimple,
				Explode: false,
			})

			if err := func() error {
				val, err := d.DecodeValue()
				if err != nil {
					return err
				}

				c, err := conv.ToInt(val)
				if err != nil {
					return err
				}

				params.MatchID = c
				return nil
			}(); err != nil {
				return err
			}
		} else {
			return validate.ErrFieldRequired
		}
		return nil
	}(); err != nil {
		return params, &ogenerrors.DecodeParamError{
			Name: "match_id",
			In:   "path",
			Err:  err,
		}
	}
	return params, nil
}

// GetMatchForLogParams is parameters of GetMatchForLog operation.
type GetMatchForLogParams struct {
	LogID int
}

func unpackGetMatchForLogParams(packed middleware.Parameters) (params GetMatchForLogParams) {
	{
		key := middleware.ParameterKey{
			Name: "log_id",
			In:   "path",
		}
		params.LogID = packed[key].(int)
	}
	return params
}

func decodeGetMatchForLogParams(args [1]string, argsEscaped bool, r *http.Request) (params GetMatchForLogParams, _ error) {
	// Decode path: log_id.
	if err := func() error {
		param := args[0]
		if argsEscaped {
			unescaped, err := url.PathUnescape(args[0])
			if err != nil {
				return errors.Wrap(err, "unescape path")
			}
			param = unescaped
		}
		if len(param) > 0 {
			d := uri.NewPathDecoder(uri.PathDecoderConfig{
				Param:   "log_id",
				Value:   param,
				Style:   uri.PathStyleSimple,
				Explode: false,
			})

			if err := func() error {
				val, err := d.DecodeValue()
				if err != nil {
					return err
				}

				c, err := conv.ToInt(val)
				if err != nil {
					return err
				}

				params.LogID = c
				return nil
			}(); err != nil {
				return err
			}
		} else {
			return validate.ErrFieldRequired
		}
		return nil
	}(); err != nil {
		return params, &ogenerrors.DecodeParamError{
			Name: "log_id",
			In:   "path",
			Err:  err,
		}
	}
	return params, nil
}

// GetPlayersParams is parameters of GetPlayers operation.
type GetPlayersParams struct {
	ID []int
}

func unpackGetPlayersParams(packed middleware.Parameters) (params GetPlayersParams) {
	{
		key := middleware.ParameterKey{
			Name: "id",
			In:   "query",
		}
		params.ID = packed[key].([]int)
	}
	return params
}

func decodeGetPlayersParams(args [0]string, argsEscaped bool, r *http.Request) (params GetPlayersParams, _ error) {
	q := uri.NewQueryDecoder(r.URL.Query())
	// Decode query: id.
	if err := func() error {
		cfg := uri.QueryParameterDecodingConfig{
			Name:    "id",
			Style:   uri.QueryStyleForm,
			Explode: false,
		}

		if err := q.HasParam(cfg); err == nil {
			if err := q.DecodeParam(cfg, func(d uri.Decoder) error {
				return d.DecodeArray(func(d uri.Decoder) error {
					var paramsDotIDVal int
					if err := func() error {
						val, err := d.DecodeValue()
						if err != nil {
							return err
						}

						c, err := conv.ToInt(val)
						if err != nil {
							return err
						}

						paramsDotIDVal = c
						return nil
					}(); err != nil {
						return err
					}
					params.ID = append(params.ID, paramsDotIDVal)
					return nil
				})
			}); err != nil {
				return err
			}
			if err := func() error {
				if params.ID == nil {
					return errors.New("nil is invalid value")
				}
				return nil
			}(); err != nil {
				return err
			}
		} else {
			return validate.ErrFieldRequired
		}
		return nil
	}(); err != nil {
		return params, &ogenerrors.DecodeParamError{
			Name: "id",
			In:   "query",
			Err:  err,
		}
	}
	return params, nil
}

// GetTeamParams is parameters of GetTeam operation.
type GetTeamParams struct {
	ID int
}

func unpackGetTeamParams(packed middleware.Parameters) (params GetTeamParams) {
	{
		key := middleware.ParameterKey{
			Name: "id",
			In:   "path",
		}
		params.ID = packed[key].(int)
	}
	return params
}

func decodeGetTeamParams(args [1]string, argsEscaped bool, r *http.Request) (params GetTeamParams, _ error) {
	// Decode path: id.
	if err := func() error {
		param := args[0]
		if argsEscaped {
			unescaped, err := url.PathUnescape(args[0])
			if err != nil {
				return errors.Wrap(err, "unescape path")
			}
			param = unescaped
		}
		if len(param) > 0 {
			d := uri.NewPathDecoder(uri.PathDecoderConfig{
				Param:   "id",
				Value:   param,
				Style:   uri.PathStyleSimple,
				Explode: false,
			})

			if err := func() error {
				val, err := d.DecodeValue()
				if err != nil {
					return err
				}

				c, err := conv.ToInt(val)
				if err != nil {
					return err
				}

				params.ID = c
				return nil
			}(); err != nil {
				return err
			}
		} else {
			return validate.ErrFieldRequired
		}
		return nil
	}(); err != nil {
		return params, &ogenerrors.DecodeParamError{
			Name: "id",
			In:   "path",
			Err:  err,
		}
	}
	return params, nil
}
