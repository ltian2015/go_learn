package encoding

import (
	"encoding/binary"
	"fmt"
	"testing"
)

/**
!!!编码知识点1： 二进制字节的内存存储顺序：
小端序(little endian)，是指二进制数的低位数字(Least Significant Byte - LSB)在低地址的字节，
    高位数字(Most Significant Byte - MSB)在高地址字节。小端序复合CPU的阅读与书写习惯，
	即：先处理（写入或读取）低位数字，再处理（写入或读取）高位数字。
	比如十进制数字的500的16进制为 ：（高地址） 01 F4（低地址）,
	按照小端序，则先在低地址内存写入低位数字F4，再在高地址内存写入高位数字01。
大端序（big endian）则正好相反，按照人类的阅读与书写习惯，无论是十进制还是16进制，都是先处理（阅读或书写）高位
    数字，再处理（书写或阅读）低位数字。
	比如，十进制数字的500的16进制表示位01 F4,按照大端序，则先在低地址内存写入高位数字01，
      再在高地址内存写入低位数字4，应为：（高地址）F4 01（低地址），
!!!编码知识点2：有符号整数的表示。
    有符号整数的最高位是符号位，1表示为负数，为0表示整数。比如：
	一个字节的有符号整数范围是： -128 —— +127，也就是负2的7次方，到正的2的7次方-1.
	两个字节的有符号整数范围是： -32768 —— +32767，也就是负2的15次方，到正的2的15次方 -1.
    有符号整数采用2的补码表示，
	也就是：
	      最高位为0，其余位数为0的数，是最小的整数0.加1就是1，依次类推，一直加到表示符号位的最高位之外的所有
		           位为1时，就达到了正数的最大值，如果再加1，则最高位变为1，其余位变为0，也就是负数的最小值。
				   注意，有符号数的最大正数加1，就“溢出”为有符号数的最小值。+1的溢出得到的值相当于“按位取反”
		  最高位为1，其余位为0的数，是最小的负数，该数加1，也就是最低位为1的时候，表示次小的负数，
		  以此类推，当所有的位都为1的时候，表示最大的负数-1.
	      注意，0不是负数，此时，如果最大的负数-1再加1，则会导致字节溢出，变成最小的正数的0，+1的溢出得到的值相当于“按位取反”
          !!!因此，有符号数可以通过无符号数的补码进行计算。
		  !!!算法就是：求一个负数的二进制表示，可以把该负数对应的正的相反数按位取反，然后再加上二进制的1即可。
		             正数按位取反之后加1，与该正数减1后按位取反效果相等。
		  !!!        反向操作，将一个有符号的负数的二进制表示减1，然后再按位取反，即可得到相反数。

!!!编码知识点3：实现原子的无符号整数的减法操作。
  注意，原子的无符号整数类型没有减法操作，只有加法操作。所以要将其转换为与相反数的相加操作。
  !!! 就是让减数按位取反后加1，这样，加上这个结果，就相当于减去这个数。
  ^操作符做一元操作符时，是按位取反操作，即，^n 相当于^作为“异或”二元操作符时，m^n操作，
   m与n类型相同，所有位都是1的数。异或二元操作的结果是两个数，如果对应的位的值相异（不同），
   则结果对应位的值为1，相同为0.
   一个位上只有1，0，所以有三种情况：1—1；1-0；0—0,对应的按位计算的操作为：
	按位与 &  : 1-1得1，1-0得0，0-0得0
	按位或 |  : 1-1得1，1-0得1，0-0得0
	按位异或^ : 1-1得0，1-0得1，0-0得1
!!! 编码知识点4：所谓字符或者由字符在组成的字符串，其实也是unicode点码值，用rune类型表示。
!!! uft-8编码则是用变长字节来存储rune值。

**/

// 把固定长度的int32整数编码为二进制的字节数组\切片([]byte)
func TestMaxInt(t *testing.T) {
	//!!! 对于有符号整数，先取反后加1，就是对应的相反数，等效于  先减1后取反。
	//!!! 注意，最小的负数，对应的正数相反数不在取值范围之内。
	//整数场景验证
	var number1 int64 = 9
	var oppositeNumber1 int64 = ^number1 + 1
	var oppositeNumber1M2 int64 = ^(number1 - 1)
	println("number1=", number1, " oppositeNumber1=", oppositeNumber1)
	println("number1=", number1, " oppositeNumber1M2=", oppositeNumber1M2)
	//负数场景验证
	var number2 = -9
	var oppositeNumber2 = ^(number2 - 1)
	var oppositeNumber2M2 = ^number2 + 1
	println("number2=", number2, " oppositeNumber2=", oppositeNumber2)
	println("number2=", number2, " oppositeNumber2M2=", oppositeNumber2M2)

}

type SignedInteger interface {
	int | int8 | int16 | int32 | int64
}
type OverRangeError int64

func (ore OverRangeError) Error() string {

	return fmt.Sprintf("给定类型最小负数%d的相反数超界", int64(ore))
}

//  整数取相反数函数
func oppsiteNumber[SN SignedInteger](number SN) (SN, error) {
	var err error = nil
	oppsiteNumber := ^number + 1
	//!!! 补码表示的最小负数的按上述操作的相反数仍是自身，和0一样。
	if (number != 0) && (oppsiteNumber == number) {
		err = OverRangeError(number)
	}
	return oppsiteNumber, err
}

// 测试上面取相反数函数是否正确。
func TestOppsiteNumber(t *testing.T) {
	var aInt int64 = -0
	var aIntOpp, _ = oppsiteNumber(aInt)
	fmt.Printf("-(%d)=%d\n", aInt, aIntOpp)
	const maxInt64 int64 = 1<<63 - 1      //符号位为0，其他位为1，即：+9223372036854775807
	const minInt64 int64 = ^(1 << 63) + 1 //符号位为1，其他位为0，即：-9223372036854775808
	const minInt32 int32 = ^(1 << 31) + 1
	maxInt64Opp, _ := oppsiteNumber(maxInt64)
	fmt.Printf("-(%d)=%d\n", maxInt64, maxInt64Opp)
	minInt64Opp, err := oppsiteNumber(minInt64)
	if err != nil {
		println(err.Error())
	} else {
		fmt.Printf("-(%d)=%d\n", minInt64, minInt64Opp)
	}
	minInt32Opp, err := oppsiteNumber(minInt32)
	if err != nil {
		println(err.Error())
	} else {
		fmt.Printf("-(%d)=%d\n", minInt32, minInt32Opp)
	}
}

// 验证大端序和小端序整数编码在byte切片中。
func TestEncodeInt32(t *testing.T) {
	v := uint32(500)
	fmt.Printf("数字500的二进制为：%b\n", v)
	fmt.Printf("数字500的16进制为：%X\n", v)
	//存储大端序字节（高位字节在前，低位字节在后，即：高地址内存写入低位字节，低地址内存写入高位字节）
	bigEndianBuf := make([]byte, 4)
	//存储小段旭字节（低位字节在前，高位字节在后，即：低地址内存写入低位字节，高地址内存写入高位字节）
	littleEndianBuf := make([]byte, 4)

	binary.BigEndian.PutUint32(bigEndianBuf, v)
	binary.LittleEndian.PutUint32(littleEndianBuf, v)

	fmt.Printf("数字500的大端序编码的16进制为：%X\n", bigEndianBuf)
	fmt.Printf("数字500的小端序编码的16进制为：%X\n", littleEndianBuf)
	//由于littleEndianBuf2所初始化的2个字节内存小于int32的4个字节的长度，所以会抛出panic。
	littleEndianBuf2 := make([]byte, 2)
	bigEndianBuf2 := make([]byte, 2)
	binary.LittleEndian.PutUint16(littleEndianBuf2, uint16(v))
	binary.BigEndian.PutUint16(bigEndianBuf2, uint16(v))
	fmt.Printf("数字500的大端序编码的16进制为：%X\n", bigEndianBuf2)
	fmt.Printf("数字500的小端序编码的16进制为：%X\n", littleEndianBuf2)
}

// 这个函数把字节数组解码为int32数字。
func TestDecodeInt32(t *testing.T) {
	var buf1 = make([]byte, 4)
	buf1[0] = 0
	buf1[1] = 0
	buf1[2] = 0x01
	buf1[3] = 0xF4
	var v1 uint32 = binary.BigEndian.Uint32(buf1)
	fmt.Println("字节数组1的内容是：", buf1, "按照BigEndian编码为Unit32的值是：", v1)
	//	fmt.Println("字节数组变量v的值是：", v1)
	var buf2 = make([]byte, 4)
	buf2[0] = 0xF4
	buf2[1] = 0x01
	buf2[2] = 0
	buf2[3] = 0
	var v2 uint32 = binary.LittleEndian.Uint32(buf2)
	fmt.Println("字节数组2的内容是：", buf2, "按照LittleEndian编码为Unit32的值是：", v2)
}
