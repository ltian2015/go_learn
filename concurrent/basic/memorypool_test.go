/*
*
A Pool is a set of temporary objects that may be individually saved and retrieved.
Poos是一系列可以被单独存储和获取的临时对象。

Any item stored in the Pool may be removed automatically at any time without notification.
If the Pool holds the only reference when this happens, the item might be deallocated.
在Pool中存储的任何事物(item)都可以在任何时候被自动移除，而没有通知。当这种情况发生时，如果Pool持有了该事物（item）
的唯一引用，那么这个事物（item)就会被释放（deallocated）
A Pool is safe for use by multiple goroutines simultaneously.
在多个goroutines同时使用时，Pool是安全的。

Pool's purpose is to cache allocated but unused items for later reuse,
relieving pressure on the garbage collector.
Pool的目的是缓存已经分配内存但尚未使用的事物，以备后续重新使用，从而减轻垃圾回收器的压力。
That is, it makes it easy to build efficient, thread-safe free lists.
However, it is not suitable for all free lists.
也就是说，它很容构建高效的，线程的安全的自由列表（free list）。但是，它并不适用于所有的自由列表。

An appropriate use of a Pool is to manage a group of temporary items silently shared
among and potentially reused by concurrent independent clients of a package.
Pool provides a way to amortize allocation overhead across many clients.

Pool的合适用途是管理一组临时的事物（items），默默地共享这些事物，并可能被一个包中多个并发的独立客户端再次使用。
!!!Pool提供了一种跨多个客户端分摊内存开销的方法。

An example of good use of a Pool is in the fmt package, which maintains a dynamically-sized
store of temporary output buffers. The store scales under load (when many goroutines are
actively printing) and shrinks when quiescent.
使用Pool的一个很好的范例就是fmt包，fmt包维护了一个动态大小（danamically-sized）的临时的输出缓存存储。
这个存储的大小在有负荷（当很多goroutines都在进行打印）的时候扩张，在静默时收缩。

On the other hand, a free list maintained as part of a short-lived object is not a suitable
use for a Pool, since the overhead does not amortize well in that scenario.
It is more efficient to have such objects implement their own free list.
另一方面，作为短期存活对象的一部分而维护的自由列表不适合用于Pool，因为在那个场景中开销（overhead）不会被
很好地分摊（amortize）。

A Pool must not be copied after first use.
!!!Pool在第一次使用后就不允许被拷贝。

In the terminology of the Go memory model, a call to Put(x) “synchronizes before” a call
to Get returning that same value x. Similarly, a call to New returning x “synchronizes
before” a call to Get returning that same value x.
在Go内存模型术语中，Put(x)的调用“同步先于” 获取（Get）同一个x值的调用。同样，返回x值的新建（New）调用
“同步先于” 获取（Get）同一个x值的调用。
**/

package basic

import (
	"bytes"
	gzip "compress/gzip"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"strings"
	"sync"
	"testing"
)

// 存储1024个"how now brown cow"字符串的压缩结果
var gzcow bytes.Buffer
var buf = gzcow.Bytes()

// 主要是将字符串"how now brown cow"的1024个拷贝压缩，结果存放在包变量gzcow中。
func init() {
	var data strings.Builder
	//构造一个含1024个"how now brown cow"的字符串。
	for i := 0; i < 1024; i++ {
		data.WriteString("how now brown cow")
	}
	//压缩字符串
	gz := gzip.NewWriter(&gzcow)
	if _, err := gz.Write([]byte(data.String())); err != nil {
		log.Fatal(err)
	}
	if err := gz.Close(); err != nil {
		log.Fatal(err)
	}
	for j := 0; j < 100; j++ {
		var vi uint8 = 1
		buf = append(buf, vi)
		vi += 1
	}

	println("压缩前长度为 ：", len([]byte(data.String())), "压缩后长度为 ：", len(gzcow.Bytes()))
}

// 根据sync.Pool所提炼出来的接口定义
type pool interface {
	Get() interface{}
	Put(x interface{})
}

// nopool  这类 “pool” 只是每次简单地在堆上分配一个新的gzip Reader。
type nopool struct{}

func (*nopool) Get() interface{} { return new(gzip.Reader) }

func (*nopool) Put(x interface{}) {}

// 这是一个人为刻意设计的基准测试。
// gunziploop从pool中取出gzip.Reader(解压缩器)，使用该解压缩器完成对包变量gzcow中所存储的
// 对1024个“how now brown cow”压缩结果的解压缩后，再将解压缩器缓存到pool中。
// 以便后续再使用。
func gunziploop(b *testing.B, m pool) {
	for i := 0; i < b.N; i++ {
		//gzip.Reader是gzip解码器
		r := m.Get().(*gzip.Reader)
		//复位gzip.Reader，使之处于刚构建出来时的初始状态，同时将其底层reader替换为新的IO Reader
		//从而使gzip.Reader可以再次工作。
		r.Reset(bytes.NewReader(gzcow.Bytes()))
		//以gzip.Reader r为源，将读出的（解压缩）字节拷贝到目标的写入器ioutil.Discard中。
		if n, err := io.Copy(ioutil.Discard, r); err != nil {

			b.Fatal(err)

		} else if int(n) != 1024*len("how now brown cow") {

			b.Fatal("bad length")
		}
		//再将r缓存到pool中。因为pool中缓存的对象可能会被垃圾回收。
		m.Put(r)
	}
}

func BenchmarkGunzipNopool(b *testing.B) {
	gunziploop(b, new(nopool))
}

func BenchmarkGunzipPooled(b *testing.B) {
	var newGzipReaderFunc = func() any {
		return new(gzip.Reader)
	}
	var pool = &sync.Pool{New: newGzipReaderFunc}
	gunziploop(b, pool)
}

const maxalloc int = 9

// contrived benchmark: read either 10 or maxalloc-1
// bytes from the bigcow reader into an allocated buffer
// of that size. Hold 50 allocations at once to exercise
// the pooled allocators.
// 人为的基准测试：从此bigcow 读取器（reader）中读取10个或maxalloc-1个字节到一个
// 按照对应大小所分配buffer中。一次性持有50块分配的内存用于测试池化的内存分配。
func allocloop(b *testing.B, r io.ReadSeeker, m alloc) {
	for i := 0; i < b.N; i++ {
		var bufs [50][]byte
		for i := range bufs {
			n := 10
			if rand.Intn(10) == 0 { //取10以内的随机数作为开辟的内存大小，如果为0，则取maxalloc - 1
				n = maxalloc - 1
			}
			bufs[i] = m.Alloc(n)   //调用内存分配器，分配指定大小的内存（ bytes切片）
			if len(bufs[i]) != n { //处于谨慎，检查内存分配的结果
				b.Fatal("dishonest allocator")
			}
		}
		//从读取器中读取数据到buffer中
		for i := range bufs {
			r.Seek(0, 0)
			if _, err := io.ReadFull(r, bufs[i]); err != nil {
				b.Fatal(err)
			}
		}
		//释放分配的内存缓存
		for i := range bufs {
			m.Free(bufs[i])
		}
	}
}

// alloc定义了内存分配器的接口
type alloc interface {
	Alloc(n int) []byte
	Free([]byte)
}

// profile heap allocations with a simple wrapper for
// the alloc interface; frees are implicit with Go's GC.
// 使用alloc接口的简单的封装器来剖析（profile）对内存的分配；
// 隐含使用Go的垃圾回收器（GC）进行内存释放。
type heap struct{}

func (*heap) Alloc(n int) []byte {
	//私用make函数构建一个byte切片
	return make([]byte, n)
}

// 释放内存 ，实际上什么都没做，利用go的GC机制释放内存
func (*heap) Free([]byte) {}

var rcow = bytes.NewReader(buf)

func BenchmarkHeapAlloc(b *testing.B) {
	var data []byte = make([]byte, 10)
	println(rcow.Read(data))
	for _, vi := range buf {
		println(uint8(vi))
	}
	allocloop(b, rcow, &heap{})
}
