package main

import (
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"strconv"
	"strings"

	"go.opentelemetry.io/contrib/bridges/otelslog"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
)

// 注意，Trace/Span体系,Logger体系以及Metrics体系共享相同遥测对象Name（默认是程序包名）
// !!! 有助于实现跨信号的聚合
const name = "go.opentelemetry.io/otel/example/dice"

var (
	tracer = otel.Tracer(name) //!!! 追踪器是用来创建和管理跨度（Span）的核心组件
	/**
	   !!! Meter（仪表））是用来创建和管理度量项（Metrics）的核心组件，
	   !!! Meter（仪表）通过创建不同的量具（Instruments）来完成不同度量项的测量数据的度量（mesuarement）。
	   !!! Meter 由程序中MeterProvider 创建。 OpenTelemetry 提供了一个默认的全局 MeterProvider实现。
	   !!! 通过调用全局 MeterProvider 实现的 Meter 方法并传入一个名称来创建一个 Meter 实例。语句为：
	   !!! otel.Meter(name)
	**/
	meter = otel.Meter(name) //!!! Meter与trace，log使用相同的名字有助于实现跨信号的聚合
	//!!!logger的名字通常是包名，用于定位log的来源
	logger = otelslog.NewLogger(name)
	//!!! rollCnt 是一个 Int64Counter 类型的度量器，用于记录“骰子投掷”活动的次数
	rollCnt metric.Int64Counter
)

func init() {
	var err error
	rollCnt, err = meter.Int64Counter("dice.rolls",
		metric.WithDescription("The number of rolls by roll value"),
		metric.WithUnit("{roll}"))
	if err != nil {
		panic(err)
	}
}

// rolldice 掷骰子 游戏的http 处理器
func rolldice(w http.ResponseWriter, r *http.Request) {
	/**
	Start会创建一个Span,和一个包含新创建的Span的context.Context。如果 `ctx` 中
	提供的 context.Context 已包含了一个Span，则新创建的Span将成为该Span的子Span，
	否则它将是根Span。可以通过提供 `WithNewRoot()` 作为SpanOption 来覆盖此行为，
	这样，即使 `ctx` 中已有了Span，新创建的Span也会成为根Span。创建Span时，建议使用
	`WithAttributes()`SpanOption提供所有已知的 Span 属性，因为采样器只能访问创
	建 Span 时提供的属性。任何被创建的 Span 也必须被结束（End），这是使用者的责任。如
	果Span 未结束，此 API 的实现可能会泄漏内存或其他资源。
	//!!! contex是建立Trace/Span体系与Logger体系以及Metrics体系关联关系的桥梁
	**/
	ctx, span := tracer.Start(r.Context(), "roll")
	defer span.End() //如果不结束span，则可能会泄漏内存或其他资源。

	roll := 1 + rand.Intn(6) // 随机生成1-6的数字，模拟掷骰子

	var msg string
	if player := r.PathValue("player"); player != "" {
		kv := strings.Split(player, "=")
		if len(kv) != 2 {
			msg = "未知玩家正在投掷骰"
		} else {
			playerName := kv[1]
			msg = fmt.Sprintf("玩家%s正在投掷骰子", playerName)
		}

	} else {
		msg = "匿名玩家正在投掷骰子"
	}
	//!!! 通过使用相同的来建立Trace/Span体系与Logger体系的关系，即：Logger体系可记录下每条日志的trace_id/span_id
	logger.InfoContext(ctx, msg, "result", roll)

	rollValueAttr := attribute.Int("骰子.点数", roll)
	span.SetAttributes(rollValueAttr)
	/**
	Trace 的 Span 与 Metrics 的 Instruments（如 Int64Counter）并没有“父子关系”或绑定关系。
	两者通过 context.Context 进行关联：用携带当前Span 的 ctx 去记录指标时，
	SDK 会将该 Span 的上下文用于度量数据的关联，常见形式是生成带 trace_id/span_id
	的采样数据（Exemplars）。在后端，会据此把某个数据点和具体Span 关联起来。
	具体方式是Instruments（如 Int64Counter）在记录采样数据时，比如 counter.Add(ctx, value, ...)
	时，传入当前活跃Span的ctx即可。如果该Span被采样且启用了Exemplars（数据采样），
	SDK会把该数据点的 exemplar（采样） 标记上 trace_id/span_id。同时，要通过 Attributes
	对齐采样的数据点: 在 Span 和 Metrics 上使用一致的业务属性（如 user_id, order_id）便于
	查询侧的关联。最后还要共享 Resource: Trace 和 Metrics 共用 Resource（服务名、版本等），
	助力实现跨信号的聚合。
	**/
	rollCnt.Add(ctx, 1, metric.WithAttributes(rollValueAttr)) //!!! 使用与span相同的contex

	resp := "投掷的骰子点数为:" + strconv.Itoa(roll) + "\n"
	if _, err := io.WriteString(w, resp); err != nil {
		logger.ErrorContext(ctx, err.Error())
	}
}
