package engine

import (
	"context"

	"github.com/google/uuid"
)

type httpContext string

var RequestKey httpContext = "reqKey"

type ReqContext struct {
	ReqID     string  `json:"req_id"`
	UserID    *string `json:"user_id"`
	UserEmail *string `json:"user_email"`
}

func (e *HttpEngine) SetContext(ctx context.Context, reqContext *ReqContext) context.Context {
	var request ReqContext

	if reqContext == nil {
		reqID := uuid.New().String()
		request = ReqContext{
			ReqID:     reqID,
			UserID:    nil,
			UserEmail: nil,
		}
	} else {
		request = *reqContext
	}

	c := context.WithValue(ctx, RequestKey, request)
	return c
}

func (e *HttpEngine) ParseContext(ctx context.Context) *ReqContext {
	val := ctx.Value(RequestKey)
	if val == nil {
		return &ReqContext{}
	}

	meta, ok := val.(ReqContext)
	if !ok {
		return &ReqContext{}
	}

	return &meta
}
