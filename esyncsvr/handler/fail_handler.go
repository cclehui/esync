package handler

import (
	"context"
	"errors"
	"fmt"

	"github.com/cclehui/esync/esyncsvr/esyncsvc"
	"github.com/cclehui/esync/esyncutil"
)

// 会失败的handler
type FailHandler struct {
	FailNum int
}

func (handler *FailHandler) GetHandlerID() string {
	return "esync_fail_220209"
}

func (handler *FailHandler) Do(ctx context.Context, params *esyncsvc.HandlerParams) error {
	handlerInfo, _ := params.EventDefaultDao.GetHandlerInfo()

	if handler.FailNum > 0 && handlerInfo.FailCount >= handler.FailNum {
		esyncutil.GetLogger().Infof(ctx, "FailHandler 第:%d次执行, 成功", handlerInfo.FailCount+1)
		return nil
	}

	return errors.New(fmt.Sprintf("FailHandler 第:%d次执行, 失败", handlerInfo.FailCount+1))
}
