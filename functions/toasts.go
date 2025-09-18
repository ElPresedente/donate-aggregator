package functions

import (
	"context"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

func toastRun(ctx context.Context, message string, msgType string) {
	if message == "" {
		message = "!"
	}

	toastData := map[string]interface{}{
		"message": message,
		"type":    msgType,
	}
	runtime.EventsEmit(ctx, "toastExec", toastData)
}

func ToastSuccessRun(ctx context.Context, message string) {
	msgType := "success"
	toastRun(ctx, message, msgType)
}

func ToastErrorRun(ctx context.Context, message string) {
	msgType := "error"
	toastRun(ctx, message, msgType)
}

func ToastInfoRun(ctx context.Context, message string) {
	msgType := "info"
	toastRun(ctx, message, msgType)
}
