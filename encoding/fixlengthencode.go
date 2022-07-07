package encoding

import (
	"encoding/binary"
	"fmt"
)

//把固定长度的int32整数编码为二进制的字节数组\切片([]byte)
func EncodeInt32() {
	v := uint32(500)
	fmt.Printf("数字500的二进制为：%b\n", v)
	fmt.Printf("数字500的16进制为：%X\n", v)
	//大端序按照人类的阅读与书写习惯，无论是十进制还是16进制，都是先写入高位数字，再写入低位数字。
	//比如十进制数字的500的16进制表示位01 F4,按照大端序，则先写入高位数字01，再写入低位数字F4.
	//在内存中，先写入的内存相对地址较小，后写入的内存相对地址较高。所以，如果以两字节存储，应为01 F4，
	//用4个字节存储就是00 00 01 F4 （高位在前）
	bigEndianBuf := make([]byte, 4)
	//小端序按照CPU的阅读与书写习惯，无论是十进制还是16进制，都是先写入低位数字，再写入高位数字。
	// 比如十进制数字的500的16进制为01F4,按照小端序，则先写入低位数字F4，再写入高位数字01.
	///在内存中，先写入的内存相对地址较小，后写入的内存相对地址较高。所以，占用两个字节 就是F4 01,
	//如果占用4个人字节就是 F4 01 00 00,(越高位越向后，高位在后）

	littleEndianBuf := make([]byte, 4)

	binary.BigEndian.PutUint32(bigEndianBuf, v)
	binary.LittleEndian.PutUint32(littleEndianBuf, v)

	fmt.Printf("数字500的大端序编码的16进制为：%X\n", bigEndianBuf)
	fmt.Printf("数字500的小端序编码的16进制为：%X\n", littleEndianBuf)
	//由于littleEndianBuf2所初始化的2个字节内存小于int32的4个字节的长度，所以会抛出panic。
	littleEndianBuf2 := make([]byte, 2)
	binary.LittleEndian.PutUint32(littleEndianBuf2, v)
}

//这个函数把字节数组解码为int32数字。
func DecodeInt32() {
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
