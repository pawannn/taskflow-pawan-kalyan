package engine

import (
	"context"

	requestcontext "github.com/pawannn/taskflow-pawan-kalyan/backend/internal/pkg/requestContext"
	utils "github.com/pawannn/taskflow-pawan-kalyan/backend/internal/utils"
)

// SetContext attaches request metadata (ReqContext) to the given context.
func (e *HttpEngine) SetContext(ctx context.Context, reqContext *requestcontext.ReqContext) context.Context {
	if reqContext == nil {
		reqContext = &requestcontext.ReqContext{
			ReqID: utils.GenerateUUID(),
		}
	}

	return context.WithValue(ctx, requestcontext.RequestKey, reqContext)
}

// ParseContext extracts request metadata (ReqContext) from the given context.
func (e *HttpEngine) ParseContext(ctx context.Context) *requestcontext.ReqContext {
	val := ctx.Value(requestcontext.RequestKey)
	if val == nil {
		return &requestcontext.ReqContext{}
	}

	meta, ok := val.(*requestcontext.ReqContext)
	if !ok {
		return &requestcontext.ReqContext{}
	}

	return meta
}
