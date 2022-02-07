package eventrp

import (
	"context"

	"git2.qingtingfm.com/podcaster/papi-go/define"
	"git2.qingtingfm.com/podcaster/papi-go/repository/eventrp/common"
	"git2.qingtingfm.com/podcaster/papi-go/repository/eventrp/handler/program/openapiprogram"
	"git2.qingtingfm.com/podcaster/papi-go/repository/eventrp/handler/program/papiprogram"
)

type HandlerBase interface {
	GetHandlerID() string // 建议日期加三位数字 比如 210916002
	Do(ctx context.Context, params *common.HandlerParams) error
}

func GetHandlerList(handlerParams *common.HandlerParams) []HandlerBase {
	result := make([]HandlerBase, 0)

	// 不同事件类型不同处理方法
	switch handlerParams.EventData.EventType {
	case define.ETypeOpenApiUploadProgramCheck:
		result = append(result,
			&openapiprogram.TransStatusCallbackHandler{},
			&openapiprogram.TransFailHandler{})
	case define.ETypePapiProgramCreate:
		result = append(result, &papiprogram.TriggerTransWaitEventHandler{})
	case define.ETypePapiProgramTransWait:
		result = append(result, &papiprogram.TransFailHandler{})
	case define.ETypePapiProgramBatchSort:
		result = append(result, &papiprogram.BatchSortHandler{})
	default:
		return result
	}

	return result
}
