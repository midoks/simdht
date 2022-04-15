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

var SysStatus runtime.MemStats

func init() {
	go func() {
		runtime.ReadMemStats(&SysStatus)
		time.Sleep(time.Microsecond * 100)
	}()
}

func GcInfo(ctx *fasthttp.RequestCtx) {

	fmt.Fprintf(ctx, "GC 下次内存回收量\t:%s\r\n", tools.SizeFormat(float64(SysStatus.NextGC)))
	fmt.Fprintf(ctx, "GC 距离上次时间\t:%0.2fs\r\n", float64(time.Now().UnixNano()-int64(SysStatus.LastGC))/1000/1000/1000)
	fmt.Fprintf(ctx, "GC 暂停时间总量\t:%0.2fs\r\n", float64(SysStatus.PauseTotalNs)/1000/1000/1000)
	fmt.Fprintf(ctx, "GC 上次暂停时间\t:%.3f\r\n", float64(SysStatus.PauseNs[(SysStatus.NumGC+255)%256])/1000/1000/1000)
	fmt.Fprintf(ctx, "GC 执行次数	:%d\r\n", SysStatus.NumGC)

	fmt.Fprintf(ctx, "Goroutine Num:%d\r\n", runtime.NumGoroutine())
	fmt.Fprintf(ctx, "GC Status\r\n")
}
