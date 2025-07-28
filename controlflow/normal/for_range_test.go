package normal

import (
	"iter"
	"testing"
)

/*
*
迭代器是一个函数，它将序列中的连续元素传递给一个"回调函数(callback)"，该函数通常命名为 yield。
该函数在序列结束时停止，或者当 yield 返回 false 时停止，表示提前终止迭代。
yield是生产的意思，该函数的返回值表示是否让迭代器生产元素，该函数的输入就是迭代器当前遍历到的值。
Go语言的iter标准包中定义了 [Seq] 和 [Seq2]（发音与“seek”相同，即“sequence”的第一个音节）
作为迭代器的简写形式，用于传递每个序列元素的 1 或 2 个值 给yield回调函数：
!!! 从上面的Push迭代器函数的例子可以看出，Seq 和 Seq2 就是一种Push迭代器。。

	type (
	    Seq[V any]     func(yield func(V) bool)
	    Seq2[K, V any] func(yield func(K, V) bool)
	)

Seq2 表示成对值的序列，通常为键值对或索引值对。
yield 回调函数返回 true 表示迭代器应继续处理序列中的下一个元素，返回 false 表示应停止。
例如，[maps.Keys] 返回一个迭代器，该迭代器生成Map m 的键序列，实现如下：

	func Keys[Map ~map[K]V, K comparable, V any](m Map) iter.Seq[K] {
		return func(yield func(K) bool) {
			for k := range m {
				if !yield(k) {
					return
				}
			}
		}
	}

更多示例可参见[The Go Blog: Range Over Function Types]。

!!! 迭代器函数通常由[range 循环]调用，例如：

	func PrintAll[V any](seq iter.Seq[V]) {
	    for v := range seq {
	        fmt.Println(v)
	    }
	}

# !!!命名规范

迭代器函数和方法的命名基于正在遍历的序列：

	// All 返回一个遍历 s 中所有元素的迭代器。
	func (s *Set[V]) All() iter.Seq[V]

集合类型的迭代器方法通常命名为 All，因为它遍历集合中所有值的序列。
对于包含多个可能序列的类型，迭代器的名称可以指示正在提供的序列：

	    // Cities 返回该国主要城市的迭代器。
	    func (c *Country) Cities() iter.Seq[*City]
		// Languages 返回该国官方语言的迭代器。
	    func (c *Country) Languages() iter.Seq[string]

如果迭代器需要额外配置，构造函数
可以接受额外的配置参数：

	    // Scan 返回一个遍历键值对的迭代器，其中 min ≤ key ≤ max。
	    func (m *Map[K, V]) Scan(min, max K) iter.Seq2[K, V]

		// Split 返回一个迭代器，遍历 s 的子字符串（可能为空），
	    // 这些子字符串由 sep 分隔。
	    func Split(s, sep string) iter.Seq[string]

当存在多种可能的迭代顺序时，方法名称可能指示该顺序：

	    // All 返回一个从头到尾遍历列表的迭代器。
	    func (l *List[V]) All() iter.Seq[V]

	    // Backward 返回一个从尾到头遍历列表的迭代器。
		func (l *List[V]) Backward() iter.Seq[V]

	    // Preorder 返回一个遍历语法树中所有节点（包括指定根节点）的迭代器，
	    // 采用深度优先先序遍历，即先访问父节点再访问其子节点。
	    func Preorder(root Node) iter.Seq[Node]

# !!!仅能单次使用的迭代器

大多数迭代器都提供了遍历整个序列的能力：
当调用迭代器时，它会执行必要的初始化操作以启动序列，然后依次调用序列中的元素，
最后在返回前进行清理。再次调用迭代器将再次遍历该序列。
某些迭代器打破了这一惯例，仅允许遍历序列一次。这些“单次使用迭代器”通常从
无法回卷重新开始的数据流中读取值。在提前停止后再次调用迭代器可能会继续流，
但在序列完成后再次调用它将不会返回任何值。返回单次使用迭代器的函数或方法的文档注释应说明此事实：

	// Lines 返回从 r 读取的行序列的迭代器。
	// 它返回一个单次使用迭代器。
	func (r *Reader) Lines() iter.Seq[string]

# !!!提取式迭代的场景
接受或返回迭代器的函数和方法应使用标准的 [Seq] 或 [Seq2] 类型，以确保
与范围循环和其他迭代器适配器的兼容性。标准迭代器可以被视为“推式迭代器”，它们将值推送到 yield 函数。

!!! 有时范围循环并非消费序列值的最自然方式。在此情况下，iter包提供了一个[Pull]方法， 将标准推送迭代器
!!! 转换为“拉取迭代器”，可用于逐个拉取序列中的值。[Pull] 启动一个迭代器并返回一组函数
!!! ——next拉取函数和——stop函数——分别用于从迭代器中获取下一个值,并停止迭代器。例如：

// Pairs 返回一个迭代器，该迭代器遍历 seq 中连续的值对。

	unc Pairs[V any](seq iter.Seq[V]) iter.Seq2[V, V] {
			return func(yield func(V, V) bool) {
				next, stop := iter.Pull(seq)
				defer stop() //!!!  确保迭代器在使用后停止，避免内存泄漏。
				for {
					v1, ok1 := next()
					if !ok1 {
						return
					}
					v2, ok2 := next()
					// If ok2 is false, v2 should be the
					// zero value; yield one last pair.
					if !yield(v1, v2) {
						return
					}
					if !ok2 {
						return
					}
				}
			}
		}

如果客户端未完成序列的消费，必须调用 stop，这允许迭代器函数完成并返回。如示例所示，
确保这一点的传统方法是使用 defer。
!!! # 标准库使用

标准库中的几个包提供了基于迭代器的 API，其中最值得注意的是 [maps] 和 [slices] 包。
例如，[maps.Keys] 返回一个遍历映射键的迭代器，而 [slices.Sorted] 将迭代器的值收集到一个切片中，
对它们进行排序，并返回该切片，因此要遍历映射的排序键：

	for _, key := range slices.Sorted(maps.Keys(m)) {
	    ...
	}

!!! # 改变
迭代器仅提供序列的值，而不提供任何直接修改序列的方式。如果迭代器希望在迭代过程中提供修改序列的机制，
通常的做法是定义一个带有额外操作的位置类型，然后提供一个基于位置的迭代器。例如，树的实现可能提供：
// Positions 返回一个遍历序列中位置的迭代器。

	    func (t *Tree[V]) Positions() iter.Seq[*Pos]
	    // Pos 表示序列中的一个位置。
	    // 它仅在传递给 yield 调用时有效。
	    type Pos[V any] struct { ... }
		// Pos 返回光标处的值。
	    func (p *Pos[V]) Value() V
	    // Delete 删除迭代过程中此处的值。
	    func (p *Pos[V]) Delete()
	    // Set 修改光标处的值 v。
	    func (p *Pos[V]) Set(v V)

然后，客户端可以使用以下代码从树中删除无聊的值：

	for p := range t.Positions() {
	    if boring(p.Value()) {
	        p.Delete()
	    }
	}

[Go 博客：遍历函数类型]：https://go.dev/blog/range-functions
[range 循环]：https://go.dev/ref/spec#For_range

*
*/
type Set[E comparable] struct {
	//!!!struct{}是一个以结构体为底层类型的字面类型，该类型没有成员，
	// !!! 表示"空结构体类型，该类型的值不占用内存空间。
	//!!! 常用来作为map的Value以表示K的集合。或者作为channel的元素类型用于发送不占内存的数据以解除阻塞。
	m map[E]struct{}
}

func NewSet[E comparable]() *Set[E] {
	return &Set[E]{
		m: make(map[E]struct{}),
	}
}
func (s *Set[E]) Add(e E) {
	s.m[e] = struct{}{}
}
func (s *Set[E]) Remove(e E) {
	delete(s.m, e) //删除元素,delete内置函数的逻辑是：如果元素为nil或不存在，则不做任何操作。
}
func (s *Set[E]) Contains(e E) bool {
	_, ok := s.m[e]
	return ok
}
func Union1[E comparable](s1, s2 *Set[E]) *Set[E] {
	result := NewSet[E]()
	for e := range s1.m { //!!! 如果希望使用for range 遍历s1而不是 s1.m。那该怎么办？
		result.Add(e)
	}
	for e := range s2.m {
		result.Add(e)
	}
	return result
}

// !!! Push函数是一个迭代器。
// !!! 它提供一种主动推出集合内部元素给来自外部的生产函数（yield）对给定的元素的生产加工。

func (s *Set[E]) Push(yield func(E) bool) {
	for e := range s.m {
		if yield(e) == false { //将集合元素推给外部传递来的处理函数，故称push
			return
		}
	}
}
func Test_push(t *testing.T) {
	s := NewSet[int]()
	s.Add(1)
	s.Add(2)
	s.Add(3)
	s.Add(4)
	s.Add(5)
	s.Add(6)
	s.Add(7)
	s.Add(8)
	s.Add(9)
	s.Add(10)
	s.Push(func(e int) bool {
		println(e)
		if e >= 5 {
			return false
		}
		return true //返回true表示继续处理下一个元素
	})
}

// !!! Pull 通过返回“数据拉取器函数”和“停止数据拉取函数”这两个函数，来使外部可以控制对齐其内部元素的拉取。
// !!!  func() (E, bool)是元素拉取函数，该函数可以拉取返回元素，即，每调用一次就拉取一次元素。
// !!! func()是停止拉取函数，该函数可以停止拉取元素。
// !!! 这个实现的问题在于如果使用方再没有拉取所有元素，也不调用返回的停止拉取函数，就会造成阻塞中的goroutine不会释放，而造成内存泄露
// !!! 在形式上，推式迭代器Push与拉式迭代器Pull的区别在于，推式迭代器把回调函数作为Push函数的输入参数，
// !!! 而拉式迭代器则把控制数据拉取的回调函数作为Pull函数的返回结果。想象一下：Push的回调函数就像一个有两端管道，可以接收
// !!! 数据元素，进行处理，返回是否继续推送的信号。而Pull所返回的回调函数就像是只有一端的管道，只能从一端拉取数据元素。

func (s *Set[E]) Pull() (next func() (E, bool), stop func()) {
	dataCh := make(chan E)            //!!! 使用缓冲通道来存储集合中的元素
	stopFlagCh := make(chan struct{}) //!!! 使用一个停止标志通道来控制拉取操作
	go func() {
		defer func() {
			close(dataCh)
			println("dataCh closed")

		}()

		for e := range s.m {
			select {
			case dataCh <- e: //!!! 将集合中的元素发送到dataCh通道
			case <-stopFlagCh:
				return //!!! 如果接收到停止信号，则退出循环
			}
		}
	}()
	next = func() (E, bool) { //该闭包捕获了dataCh
		e, ok := <-dataCh //!!! 从dataCh通道中拉取元素
		return e, ok      //!!! 返回拉取的元素和是否成功拉取的标志
	}
	stop = func() {
		close(stopFlagCh) //!!! 关闭停止标志通道，通知拉取操作停止

	}
	return next, stop //!!! 返回拉取函数和停止函数
}

func Test_pull(t *testing.T) {
	s := NewSet[int]()
	s.Add(1)
	s.Add(2)
	s.Add(3)
	s.Add(4)
	s.Add(5)
	s.Add(6)
	s.Add(7)
	s.Add(8)
	s.Add(9)
	s.Add(10)

	next, stop := s.Pull() //!!! 获取拉取函数和停止函数
	defer stop()           //!!! 确保stop函数可以被调用，避免next迭代器未全部读取数据，中途中断时pull，函数内部gorouine的阻塞所导致的内存泄露。
	for {
		e, ok := next() //!!! 调用拉取函数获取下一个元素

		if !ok { //!!! 如果没有更多元素了，退出循环
			break
		}
		println(e)  //!!! 打印拉取到的元素
		if e == 5 { //!!! 如果拉取到的元素大于等于5，则停止拉取
			break
		}
	}
	println("program over")
}

// !!! All函数返回一个迭代器，该迭代器可以提供给for range 语句来遍历Set集合中的所有元素。
func (s *Set[E]) All() iter.Seq[E] {
	return s.Push
}

// !!! 展示了如何根据Set集合所提供的迭代器来函数，在for range 循环语句中遍历Set集合中的所有元素。
func TestIteratorUsedByForRangeLoop(t *testing.T) {
	s := NewSet[int]()
	s.Add(1)
	s.Add(2)
	s.Add(3)
	s.Add(4)
	s.Add(5)
	s.Add(6)
	s.Add(7)
	s.Add(8)
	s.Add(9)
	s.Add(10)
	for e := range s.All() { //!!! 使用for range 循环来遍历Set集合中的所有元素
		println(e) //!!! 打印每个元素
	}
}
