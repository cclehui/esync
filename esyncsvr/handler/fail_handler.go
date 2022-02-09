package handler

import (
	"context"

	"github.com/cclehui/esync/esyncsvr/service"
	"github.com/cclehui/esync/esyncutil"
)

// 会失败的handler
type FailHandler struct {
	FailNum int
}

func (handler *FailHandler) GetHandlerID() string {
	return "esync_fail_220209"
}

func (handler *FailHandler) Do(ctx context.Context, params *service.HandlerParams) error {
	esyncutil.GetLogger().Infof(ctx, "FailHandler 执行......")
	return nil
}
