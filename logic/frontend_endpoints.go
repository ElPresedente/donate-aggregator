package logic

import (
	"context"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

func NotifyDBChange(ctx context.Context, data []string) {
	runtime.EventsEmit(ctx, "db_updated", data)
}

func SomeFrontendTask() {
	//do smth
}

func TestEndpoint(arg string) {

}
