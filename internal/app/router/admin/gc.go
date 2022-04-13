package admin

import (
	"fmt"
	"runtime"
	"time"

	"github.com/midoks/simdht/internal/tools"
	"github.com/valyala/fasthttp"
)

// initTime is the time when the application was initialized.
var initTime = time.Now()

var sysStatus runtime.MemStats

func updateSystemStatus() {

	m := new(runtime.MemStats)
	runtime.ReadMemStats(m)
	sysStatus = *m
}

func GcInfo(ctx *fasthttp.RequestCtx) {

	updateSystemStatus()

	fmt.Fprintf(ctx, "GC 下次内存回收量\t:%s\r\n", tools.SizeFormat(float64(sysStatus.NextGC)))
	fmt.Fprintf(ctx, "GC 距离上次时间\t:%0.2fs\r\n", float64(time.Now().UnixNano()-int64(sysStatus.LastGC))/1000/1000/1000)
	fmt.Fprintf(ctx, "GC 暂停时间总量\t:%0.2fs\r\n", float64(sysStatus.PauseTotalNs)/1000/1000/1000)
	fmt.Fprintf(ctx, "GC 上次暂停时间\t:%.3f\r\n", float64(sysStatus.PauseNs[(sysStatus.NumGC+255)%256])/1000/1000/1000)
	fmt.Fprintf(ctx, "GC 执行次数	 :%d\r\n", sysStatus.NumGC)

	fmt.Fprintf(ctx, "Goroutine Num:%d\r\n", runtime.NumGoroutine())

	// fmt.Fprintf(ctx, "sysStatus, %s!\n", sysStatus)
	fmt.Fprintf(ctx, "GC Status\r\n")
}
