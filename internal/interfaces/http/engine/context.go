package engine

import (
	"context"

	"github.com/google/uuid"
	"github.com/pawannn/taskflow-pawan-kalyan/backend/internal/pkg/requestContext"
)

func (e *HttpEngine) SetContext(ctx context.Context, reqContext *requestContext.ReqContext) context.Context {
	var request requestContext.ReqContext

	if reqContext == nil {
		reqID := uuid.New().String()
		request = requestContext.ReqContext{
			ReqID:     reqID,
			UserID:    nil,
			UserEmail: nil,
		}
	} else {
		request = *reqContext
	}

	c := context.WithValue(ctx, requestContext.RequestKey, request)
	return c
}

func (e *HttpEngine) ParseContext(ctx context.Context) *requestContext.ReqContext {
	val := ctx.Value(requestContext.RequestKey)
	if val == nil {
		return &requestContext.ReqContext{}
	}

	meta, ok := val.(requestContext.ReqContext)
	if !ok {
		return &requestContext.ReqContext{}
	}

	return &meta
}
