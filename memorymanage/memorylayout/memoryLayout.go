package memorylayout

import (
	"fmt"
	"reflect"
)

/**
  在GO语言中，对每种数据类型进行内存分配时，考虑到了向CPU寄存器加载内存数据时的性能，
  内存块的分配以参考CPU寄存器的容量大小作为基准。
  对于64位计算机而言，CPU中有多个8位，16位，32位和64位四种容量大小的寄存器，
  这四种容量大小的寄存器可以分别加载1，2，4，8个字节的数据。
  因而，GO语言内置的各种基本类型的大小也就是1，2，4，8个字节。
  对于基本类型而言，其内存分配时需要对齐寄存器的字节数量也就分别是1，2，4，8。
  这样就可以“一次性”地将基本类型完整地从内存中加载到相应的寄存器中，提升CPU利用率。
  因此，GO语言把基本类型的内存分配与加载单位定义为该类型的对齐量（align），也就是对齐寄存器的字节数量。

  像struct这样的多个基本类型复合而成的复杂类型的内存分配与布局，也是希望可以“对齐”CPU的寄存器的字节数量，
  从而来能够“成块地”加载数据，同时又保证基本类型一定不会“被分两次加载到不同的寄存器中”，
  从而能够确保内存数据在CPU寄存器中的完整性，同时又能提升了CPU利用效率。
  再思考如何实现这一目标之前，我们需要牢记，无论复杂类型数据结构内部嵌套了多少层复杂类型数据结构，
  分解到最后仍然是由多个基本类型组成，所以，对于GO编译器而言，
  复杂结构所包含的所有基本类型所需要占用的内存数量，即，需求量，以及基本类型的“对齐量”都是可以知道。
  为了实现这一目标，GO把struct中所有字段类型的对齐量的最大值作为该struct类型的对齐量，
  按此对齐量分配和加载内存块，我们将其称为“内存对齐块”或“对齐块”。按此“对齐量”块划分“对齐块”
  并“成块地”加载数据到对应大小的CPU寄存器中，可以实现期望的目标。
  补充一点，因为在64为计算机中基本类型最大的寄存器对齐量就是8个字节，所以struct类型的对齐量最大的可能值就是8。

  由于struct可能由多个对齐量不同的基本类型字段组成，如何保证最大程度地节省内存呢？
  GO所采取的主要策略是：
  1.按照字段在struct中声明的顺序来将其放置到“对齐块”中。
  2.尽可能地将字段数据合理（考虑完整性与效率）地装到一个对齐块，并同时保证加载到CPU的寄存器中的字段数据不会破坏基本类型的完整性，
  也就是需要考虑字段类型的对齐量。这势必导致对对齐块分割。分割方法是，块内所有字段类型的最大对齐量作为块内对齐量，
  用”块内对齐量“将对齐块分割为多个”小对齐块“。如果其中某个”小对齐块“仍然装有多个字段，
  那么，该”小对齐块“也必须遵循相同的分割方法再次进行块内的分割。

 struc对齐块的划分，以及块内分割布局的大体算法是：
 1. 求出所有字段的对齐量，其中最大值作为结构体的“结构对齐量”
 2. 按“结构对齐量”初始化一个“结构对齐块”。并作为“当前块”，把“当前块”的大小，
    也就是 “结构对齐量” 设置为”剩余容量“。
 3. 把第一个字段作为“当前字段”，并进行当前字段在当前块的填充处理。

    3.1计算“当前字段”的 ”字段需求对齐块数“。
	    用当前字段实际所需的字节数除以”字段对齐量“即为”字段需求对齐块数“，如不能整除，向上圆整。

    3.2 计算”当前字段“在”当前块“中的”起始填充位置“、“可用字段对齐块数”
	    3.2.1 用“剩余容量” 除以 “字段对齐量”，商即为 当前字段的 “可用字段对齐块数”。
	    3.2.2 “块对齐量” 减去  “可用填充块数” 与”字段对齐量”的乘积，即为“起始填充位置”。
	3.3 计算“当前字段”所需要的“新增结构对齐块数”。
	    如果“可用字段对齐块数”小于”需求字段对齐块数“，那么“新增结构对齐块数”为0个。
        否则“新增结构对齐块数”=取整(（”需求字段对齐块数“ -“可用字段对齐块数”) /(“结构对齐量”/“字段对齐量”)） +1

    3.4 更新“当前块”及其“剩余容量”
        3.4.1 如果“新增结构对齐块数”为0，那么“当前块”不变。 当前块的剩余容量=“结构对齐量” -“起始填充位置”— （”字段需求对齐块数“ *”字段对齐量”）
		3.4.2 如果“新增结构对齐块数”大于0，那么“当前块”即为最后的“结构对齐块”。
		  当前块的剩余容量=“新增结构对齐块数” * “结构对齐量” +“可用字段对齐块数” * ”字段对齐量” - 当前字段实际所需的字节数
4.取下一个字段作为当前字段，回到第三步进行处理。

  这种内存对齐机制，优化了CPU对truct内存数据的加载性能，也消除了CPU错误加载基本类型数据的可能。
  但是，在这种机制下，虽然优化了CPU 的计算性能，但是如果开发者对struct内部字段变量声明的顺序不合理，
  也会导致某个字段本应与其他字段在同一个对齐块中，但却单独占据了一个对齐块，导致了内存浪费，不得不加以注意。
  下面的函数以友好的方式打印了各种合理与不合理的struct内部字段你变量声明的顺序。
  总之，按照“字段对齐量”从小到大的顺序声明struct的字段变量就不会出现内存的浪费。尤其是内存大小为0的空结构体类型的字段，
  必须放在所有字段的最前面声明。

**/
//使用反射机制，展示GOlang中，复合类型（主要是struct）内存布局信息。
//为了整体学习，把所有与MemoryLayoutTest1相关的没有必须要共享的类型都定义在该函数内部。
func MemoryLayoutTest1() {
	//占用1个字节的结构体,结构体的内存对齐量为1，
	//因为结构体的内存对齐量为1，所以按照1倍数布局该结构体的内存块，在运行期间，加载结构体数据到CPU寄存器时，
	//按照0,1，2，3，4，5...诸如此类1的整数倍的数值，于相对结构体起始位置的偏移量处内存中，
	//取出结构体中对齐量大小的数据块加载到CPU长度相应的寄存器中（8位寄存器）。
	type StructOfSize1 struct {
		aByte byte //1个字节,byte类型按1个字节对齐
	}
	//占用2个字节的结构体，结构体的内存对齐量为1,因为所有字段对齐量中，最大字段对齐量为1
	type StructOfSize2 struct {
		aByte1 byte //1个字节,地址偏移量为0
		aByte2 byte //1个字节,地址偏移量为1
	}
	//占用3个字节的结构体，结构体的内存对齐量为1,因为所有字段对齐量中，最大字段对齐量为1
	type StructOfSize3 struct {
		aByte1 byte //1个字节,地址偏移量为0
		aByte2 byte //1个字节,地址偏移量为1
		aByte3 byte //1个字节,地址偏移量为2
	}
	// 理论上只需要3个字节，结构体自身对齐量为2，按照对齐机制，实际需要（need）4个字节,
	//因为结构体的内存对齐量为2，所以编译器按照2的整数数布局该结构体的内存对齐块，运行期间，加载结构体数据到CPU寄存器时，
	//按照0,2，4，6，8，10...诸如此类2的整数倍的数值，于相对结构体起始位置的偏移量处的内存中，
	//取出结构体中对齐量大小数据块加载到CPU长度相应的寄存器中（16位寄存器）。
	type StructOfSize3n4 struct {
		aByte1  byte  //1个字节,字段对齐量为1，地址偏移量为0，在第一个对齐块中
		twoByte int16 //2个字节,字段对齐量为2，地址偏移量为2（最大），在第二个对齐块中
	}
	//字段声明顺序不合理导致的内存浪费，理论上共需要4个字节，该结构体自身对齐量为2，
	//按照对齐机制，可以用4个字节，但确实际占用了6个字节，浪费了内存！！！
	//因为结构体的内存对齐量为2，所以按照2的整数数布局该结构体的内存对齐块，运行期间，加载结构体数据到CPU寄存器时，
	//按照0,2，4，6，8，10...诸如此类2的整数倍的数值，于相对结构体起始位置的偏移量处的内存中，
	//取出结构体中对齐量大小的数据块加载到CPU长度相应的寄存器中（16位寄存器）。
	//
	type StructOfSize4n6 struct {
		aByte1  byte  //1个字节，字段对齐量为1，偏移量为0，独占一个内存布局块
		twoByte int16 //2个字节,字段对齐量为2，，偏移量为2，独占一个内存布局块。
		aByte2  byte  //1个字节，字段对齐量为1，偏移量为4，独占一个内存布局块。
	}
	//与上一个类型相比，字段声明合理的结构体。
	//理论上4个字节，实际也占用4个字节。该结构体自身对齐量为2，
	//因为结构体的内存对齐量为2，所以按照2的整数数布局该结构体的内存块(每次分配或加载2个字节的“对齐式内存块”)，运行期间，加载结构体数据到CPU寄存器时，
	//按照0,2，4，6，8，10...诸如此类2的整数倍的数值，于相对结构体起始位置的偏移量处的内存中，
	//取出结构体中对齐量大小的数据块加载到CPU长度相应的寄存器中（16位寄存器）。
	type StructOfSize4 struct {
		aByte1  byte  //1个字节，偏移量为0，,在第一个对齐式内存块中
		aByte2  byte  //1个字节，偏移量为1，与aByte1被布局到同一个内存对齐块中，在读取时与aByte1一起被加载到CPU中。
		twoByte int16 //2个字节,int16类型按2个字节对齐，偏移量为2，在第二个对齐式内存块中
	}
	//理论上只需要5个字节。该结构体自身对齐量为2，根据对齐机制，实际占用6个字节。
	//因为结构体的内存对齐量为2，所以按照2的整数数布局该结构体的内存对齐块(每次分配或加载2个字节的“对齐式内存块”)，运行期间，加载结构体数据到CPU寄存器时，
	//按照0,2，4，6，8，10...诸如此类2的整数倍的数值，于相对结构体起始位置的偏移量处的内存中，
	//取出结构体中对齐量大小的数据块加载到CPU长度相应的寄存器中（16位寄存器）。
	type StructOfSize5n6 struct {
		aByte1   byte  //1个字节，字段对齐量为1，地址偏移量为0,在第一个对齐式内存块中
		aByte2   byte  //1个字节，字段对齐量为1，地址偏移量为1,在第一个对齐式内存块中，与aByte1在同一个块中。
		aByte3   byte  //1个字节，字段对齐量为1，地址偏移量为2,在第二个对齐式内存块中
		TwoBytes int16 //2个字节,字段对齐量为2，地址偏移量为4,在第三个对齐式内存块中
	}
	//包含了空结构体对象（即0大小数据）的结构类型数据，理论上占用8个字节，实际上也占用了8个字节
	//因为所有字段对齐量中，最大字段对齐量为8，所以该结构体自身对齐量为8.
	//所以按照8的整数数布局该结构体的内存对齐块，运行期间，加载结构体数据到CPU寄存器时，
	//按照0,8，16，24，32，40...诸如此类8的整数倍的数值，于相对结构体起始位置的偏移量处的内存中，
	//取出结构体中对齐量大小的数据块加载到CPU长度相应的寄存器中（64位寄存器）。
	type StructOfSize8 struct {
		empty   struct{} //0个字节,0字节类型的对齐量为1，地址偏移量为0，在第一个对齐式内存块中。
		anInt64 int64    //8个字节，int64类型的对齐量为8，地址偏移量为0，在第一个对齐式内存块中。
	}
	//因为所有字段对齐量中，最大字段对齐量为8，所以该结构体自身对齐量为8.
	//典型的不合理的内部布局类型，主要没处理好所包含了空结构体对象（即0大小数据）的结构类型数据的声明顺序。
	//理论上只有8个字节，按照内存对齐机制，应该占用8个字节，但实际上占用了16个字节，浪费了内存！
	//按照8的整数数布局该结构体的内存对齐块，运行期间，加载结构体数据到CPU寄存器时，
	//按照0,8，16，24，32，40...诸如此类8的整数倍的数值，于相对结构体起始位置的偏移量处的内存中，
	//取出结构体中对齐量大小的数据块加载到CPU长度相应的寄存器中（64位寄存器）。
	type StructOfSize8n16 struct {
		anInt64 int64    //8个字节，int64类型的对齐量为8，偏移量为0,在第一个对齐式内存块中。
		empty   struct{} //0个字节,0字节类型的对齐量为1，偏移量为8,在第一个对齐式内存块中。因为该字段不能与anInt32字段共用同一个偏移地址。所以，必须分配一个新的8字节的对齐式内存块。
	}
	//占据了16个字节。结构体自身对齐量为8个字节。
	type StructOfSize16 struct {
		aByte      byte  //1个字节，偏移量为0,在第一个对齐式内存块中。
		twoBytes   int16 //2个字节，偏移量为2,在第一个对齐式内存块中。
		fourBytes  int32 //4个字节, 偏移量为4,在第一个对齐式内存块中。
		eightBytes int64 //8个字节，偏移量为8,在第二个对齐式内存块中。
	}
	//理论上只需要有13个字节，按照内存对齐机制，应需要16个字节，但实际上却使用了24个字节，浪费了内存！！
	type StructOfSize16n24 struct {
		myBool  bool    // 1 byte，偏移量为0，在第一个对齐内存块中
		myFloat float64 // 8 bytes，偏移量为8，在第二个对齐内存块中
		myInt   int32   // 4 bytes，偏移量为16，在第三个对齐内存块中
	}
	//类型自身的对量为8，因此每8个字节分配一个对齐块。
	type StructOfSize15n32 struct {
		a bool  //1字节，a,b可以放在一个对齐块中，但由于两个字段对齐量是1和4，那么只能把这个对齐块划分为4个字节的两个子块。前4个字节放a，后4个人字节放b
		b int32 //4字节，
		c int8  //1字节，由于a，b
		d int64 //8字节，
		e byte  //1字节
	}
	type StructOfSize15n16 struct {
		a bool
		c int8
		e byte
		b int32
		d int //在64位系统中为64位，相当于int64
	}
	type StructOfSize32 struct {
		aByte   byte   //1个字节，偏移量为0，在第一个对齐块中。
		aShort  int16  //2个字节，偏移量为2，在第一个对齐块中。
		anInt32 int32  //4个字节，偏移量为4，在第一个对齐块中。
		aSlice  []byte //24个字节，偏移量为8，在第二个对齐块开始的三个对齐块中。
	}
	type StructOfSize10n12 struct {
		sixBytes [6]byte //6个字节，对齐量为1
		fourByte int32   //4个字节，对齐量为4
	}
	type StructOfSize16n16 struct {
		sixBytes [6]byte //6个字节，对齐量为1,虽然该字段6个字节大小，但对齐量为1，可以安全地跨越多个连续的对齐块。
		twoBytes [6]byte //2个字节，对齐量为2，虽然该字段6个字节大小，但对齐量为1，可以安全地跨越多个连续的对齐块。
		fourByte int32   //4个字节，对齐量为4
	}
	type StructOfSize4n4 struct {
		twoBytes int16
		aByte1   byte
		aByte2   byte
		aByte3   byte
	}
	PrintMemoryLayout(StructOfSize1{}, "StructOfSize1")
	PrintMemoryLayout(StructOfSize2{}, "StructOfSize2")
	PrintMemoryLayout(StructOfSize3{}, "StructOfSize3")
	PrintMemoryLayout(StructOfSize3n4{}, "StructOfSize3n4")
	PrintMemoryLayout(StructOfSize4n6{}, "StructOfSize4n6")
	PrintMemoryLayout(StructOfSize4{}, "StructOfSize4")
	PrintMemoryLayout(StructOfSize5n6{}, "StructOfSize5n6")
	PrintMemoryLayout(StructOfSize8n16{}, "StructOfSize8n16")
	PrintMemoryLayout(StructOfSize8{}, "StructOfSize8")
	PrintMemoryLayout(StructOfSize16{}, "StructOfSize16")
	PrintMemoryLayout(StructOfSize16n24{}, "StructOfSize16n24")
	PrintMemoryLayout(StructOfSize32{}, "StructOfSize32")
	PrintMemoryLayout(StructOfSize15n32{}, "StructOfSize15n32")
	PrintMemoryLayout(StructOfSize15n16{}, "StructOfSize15n16")
	PrintMemoryLayout(StructOfSize10n12{}, "StructOfSize10n12")
	PrintMemoryLayout(StructOfSize16n16{}, "StructOfSize16n16")
	PrintMemoryLayout(StructOfSize4n4{}, "StructOfSize4n4")
}

//定义了一个函数字面量并赋值给一个函数变量，该函数主要用于打印一个给定的（结构体）变量的内存布局。
func PrintMemoryLayout(varStruct interface{}, typeName string) {
	typ := reflect.TypeOf(varStruct)                         //给定结构体的反射类型
	varTypeSize := int(typ.Size())                           //给定结构体变量类型所占内存字节数
	varTypeAlign := typ.Align()                              //给定结构体变量类型的内存对齐块的字节数
	alignBlockCountOfVarStruct := varTypeSize / varTypeAlign //给定结构体变量类型所占用的对齐块数量。
	fmt.Printf("结构体 %s 对齐量为 %d字节 ,占用 %d 个内存对齐块,占用内存总计 %d字节 ,各字段内存布局情况如下：\n", typeName, varTypeAlign, alignBlockCountOfVarStruct, varTypeSize)
	n := typ.NumField()
	for i := 0; i < n; i++ {
		field := typ.Field(i)
		fieldSize := int(field.Type.Size()) //字段字节数
		alignBlockIndex := (int(field.Offset) / varTypeAlign) + 1
		var alignBlockCountOfField int //字段所占用的对齐块数量
		if fieldSize <= varTypeAlign { //字段的大小小于外部结构体类型的对齐量，只能（自己或与其他字段一起）占用1个对齐块。
			alignBlockCountOfField = 1
		} else if (fieldSize % varTypeAlign) == 0 { //字段的大小大于外部结构体类型的对齐量，且是整数倍
			alignBlockCountOfField = (fieldSize / varTypeAlign)
		} else {
			alignBlockCountOfField = (fieldSize / varTypeAlign) + 1
		}
		FieldAlignBolckUsedInfo := fmt.Sprintf("位于第 %d 个对齐块开始的 %d 个对齐块中", alignBlockIndex, alignBlockCountOfField)
		fmt.Printf("field  %s 其类型对齐量为 %d字节  需用内存 %d 字节,%s 相对对结构体偏移量为 %d;\n", field.Name, field.Type.FieldAlign(), fieldSize, FieldAlignBolckUsedInfo, field.Offset)
	}
	fmt.Println("------------------------------------------------------------------------------")
}
