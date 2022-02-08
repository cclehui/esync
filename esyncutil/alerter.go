package esyncutil

import "context"

// cclehui_test
type Alerter interface {
	Errorf(ctx context.Context, format string, args ...interface{})
	Warnf(ctx context.Context, format string, args ...interface{})
	Infof(ctx context.Context, format string, args ...interface{})
}
