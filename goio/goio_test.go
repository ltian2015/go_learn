package goio

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"io"
	"math/rand/v2"
	"os"
	"strings"
	"testing"
)

func TestBasicIo(t *testing.T) {
	var s = "Hello 中国 is a string"

	var r io.Reader = strings.NewReader(s)
	//从源流拷贝所有字节到目标流中，这里是把r中的所有字节写到标准输出设备中。
	n, e := io.Copy(os.Stdout, r)
	println("拷贝了", n, "个字节")
	CheckError(e)
	// Hello it is a string
}

func CheckError(e error) {
	if e != nil {
		fmt.Println(e)
	}
}
func TestSectionReader(t *testing.T) {
	reader := strings.NewReader("a0a1a2a3a4a5a6a7a8a9\n")
	secReader := io.NewSectionReader(reader, 0, 10)
	//!!! 当SeeK方法的whence参数值为io.SeekEnd时，有效的offset参数值应为负数。
	//!!! 同样，当SeeK方法的whence参数值为io.SeekStart时，有效的offset参数值应为整数。
	//!!! 总之，当Seek方法设定的参照位置与相对于参照位置的偏移量使得实际读取位置超过了合理位置，
	//!!! 那么，read方法就会返回 (n=0,err=EOF)
	secReader.Seek(-5, io.SeekEnd)
	buf := make([]byte, 20)
	n, err := secReader.Read(buf)
	CheckError(err)
	println("读取了", n, "个字节为：", string(buf))
	//!!!  会返回 (n=0,err=EOF)
	secReader.Seek(-5, io.SeekStart)
	n, err = secReader.Read(buf)
	CheckError(err)
	println("读取了", n, "个字节为：", string(buf))

}
func TestTeeReader(t *testing.T) {

}

func TestBufferOp(t *testing.T) {
	buf := make([]byte, 0, 4096)
	println("len=", len(buf), " cap=", cap(buf))
	for i := 0; i < 4097; i++ {
		buf = append(buf, 1)
	}
	println("len=", len(buf), " cap=", cap(buf))
	buf = append(buf, 1)[:len(buf)]
	println("len=", len(buf), " cap=", cap(buf))
	io.MultiReader()
	io.Pipe()
	io.MultiWriter()
}

func generateDataFile(size, dataRange int, fileName string) (int, error) {
	const BUFFER_SIZE int = 256
	f, err := os.Create(fileName)
	if err != nil {
		return 0, err
	}
	defer f.Close()
	buffer := make([]byte, 0)
	bufferCounts := size / BUFFER_SIZE
	leftBufferSize := size - BUFFER_SIZE*bufferCounts
	var fileSize int = 0
	for i := 0; i < bufferCounts; i++ {
		for j := 0; j < BUFFER_SIZE; j++ {
			ri := rand.IntN(dataRange)
			buffer = binary.AppendVarint(buffer, int64(ri))
		}
		nw, err := f.Write(buffer)
		buffer = buffer[:0] //清空buffer，保留内存
		if err != nil {
			return fileSize, err
		}
		fileSize += nw
	}
	if leftBufferSize > 0 {
		var leftBuffer []byte = make([]byte, 0)
		for i := 0; i < leftBufferSize; i++ {
			ri := rand.IntN(dataRange)
			leftBuffer = binary.AppendVarint(leftBuffer, int64(ri))
		}
		nw, err := f.Write(leftBuffer)
		if err != nil {
			return fileSize, err
		}
		fileSize += nw
	}
	return fileSize, nil
}
func readIntDataFile(fileName string) ([]int, error) {
	f, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	buf := make([]int, 0)
	r := bufio.NewReader(f)
	for x, err := binary.ReadVarint(r); err != io.EOF || err == nil; x, err = binary.ReadVarint(r) {
		buf = append(buf, int(x))
	}
	return buf, nil

}

func TestWriteThenReadFile(t *testing.T) {
	const NUMBERS = 1024 * 1024 * 100
	const RANGE = 1024 * 1024
	n, err := generateDataFile(NUMBERS, RANGE, "test.data")
	if err == nil {
		fmt.Printf("file size is  %d bytes,contains %d int numbers \n", n, NUMBERS)
	}
	intNumbs, err := readIntDataFile("test.data")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("read %d int numbers from file \n", len(intNumbs))
}
