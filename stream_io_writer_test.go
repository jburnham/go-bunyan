package bunyan

import (
	"bufio"
	"log"
	"os"
	"testing"
)

func TestIOWriterStream(t *testing.T) {
	//str1 := new(bytes.Buffer)
	//str2 := new(bytes.Buffer)
	var str1 = os.Stdout
	var str2 = os.Stdout
	buf1 := bufio.NewWriter(str1)
	buf2 := bufio.NewWriter(str2)
	s1 := NewIOWriterStream(buf1, Info, nil)
	s2 := NewIOWriterStream(buf2, Info, nil)

	l1 := NewLogger("Test", []StreamInterface{s1})
	l2 := NewLogger("Test 2", []StreamInterface{s2})

	l1.Fatal("Test")
	l2.Fatal("Test 2")
	buf1.Flush()
	log.Println(s1)
	//written := str1.Bytes()
	//log.Println(written)
	var str3 string
	lens, _ := buf1.WriteString(str3)
	log.Println(lens)
	//log.Println(str1.Len())
	//rbuf1 := bufio.NewReader(&str1)
	//rbuf1.Seek(0, 0)
	//log.Println(*buf2)
	//line, _, _ := rbuf1.ReadLine()
	//log.Print(line)
}
