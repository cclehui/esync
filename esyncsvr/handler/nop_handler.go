package handler

import (
	"context"

	"github.com/cclehui/esync/esyncsvr/esyncsvc"
	"github.com/cclehui/esync/esyncutil"
)

// 没有操作的handler
type NopHandler struct{}

func (handler *NopHandler) GetHandlerID() string {
	return "esync_nop_220208"
}

func (handler *NopHandler) Do(ctx context.Context, params *esyncsvc.HandlerParams) error {
	esyncutil.GetLogger().Infof(ctx, "HandlerNop handler 执行......")
	return nil
}
