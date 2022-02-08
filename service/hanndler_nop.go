package service

import (
	"context"

	"github.com/cclehui/esync/esyncutil"
)

// 没有操作的handler
type HandlerNop struct{}

func (handler *HandlerNop) GetHandlerID() string {
	return "esync_nop_220208"
}

func (handler *HandlerNop) Do(ctx context.Context, params *HandlerParams) error {
	esyncutil.GetLogger().Infof(ctx, "HandlerNop handler 执行......")
	return nil
}
