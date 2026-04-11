package engine

import (
	"context"

	"github.com/pawannn/taskflow-pawan-kalyan/backend/internal/pkg/requestContext"
	"github.com/pawannn/taskflow-pawan-kalyan/backend/internal/utils"
)

func (e *HttpEngine) SetContext(ctx context.Context, reqContext *requestContext.ReqContext) context.Context {
	if reqContext == nil {
		reqContext = &requestContext.ReqContext{
			ReqID: utils.GenerateUUID(),
		}
	}

	return context.WithValue(ctx, requestContext.RequestKey, reqContext)
}

func (e *HttpEngine) ParseContext(ctx context.Context) *requestContext.ReqContext {
	val := ctx.Value(requestContext.RequestKey)
	if val == nil {
		return &requestContext.ReqContext{}
	}

	meta, ok := val.(*requestContext.ReqContext)
	if !ok {
		return &requestContext.ReqContext{}
	}

	return meta
}
