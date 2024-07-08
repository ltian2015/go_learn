package particular

import (
	"fmt"
	"strings"
)

type Tracer struct {
	//追踪的层级
	traceLevel int
	//定义文本缩进的占位符，缺省为是制表符"\t""
	traceIdentPlaceholder string
}

func NewTracer(traceIdentPlaceholder string) *Tracer {
	if traceIdentPlaceholder == "" {
		return &Tracer{traceIdentPlaceholder: "\t"}
	} else {
		return &Tracer{traceIdentPlaceholder: traceIdentPlaceholder}
	}
}

// 根据缩进级别，生成缩进占位符所组成的缩进字符串
func (t *Tracer) identLevel() string {
	return strings.Repeat(t.traceIdentPlaceholder, t.traceLevel-1)
}

// 打印追踪信息
func (t *Tracer) tracePrint(fs string) {
	fmt.Printf("%s%s\n", t.identLevel(), fs)
}

// 缩进级别加1
func (t *Tracer) incIdent() { t.traceLevel = t.traceLevel + 1 }

// 缩进级别减1
func (t *Tracer) decIdent() { t.traceLevel = t.traceLevel - 1 }

// 追踪启动
func (t *Tracer) Trace(msg string) string {
	t.incIdent() //增加缩进级别
	t.tracePrint("BEGIN " + msg)
	return msg
}

// 追踪终结
func (t *Tracer) Untrace(msg string) {
	t.tracePrint("END " + msg)
	t.decIdent() //减少缩进级别
}
