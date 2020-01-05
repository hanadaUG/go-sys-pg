package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"hash/crc32"
	"io"
	"io/ioutil"
	"os"
	"strings"
)

func main() {
	sample05()
}

// 必要な部位を切り出すio.LimitReader／io.SectionReader
func sample01() {
	reader := os.Stdin
	// たくさんデータがあっても先頭の16バイトしか読み込めないようにする。
	lReader := io.LimitReader(reader, 4)
	buf, _ := ioutil.ReadAll(lReader)
	fmt.Print(string(buf))
}

func sample02() {
	reader := strings.NewReader("Example of io.SectionReader\n")
	sectionReader := io.NewSectionReader(reader, 14, 7)
	io.Copy(os.Stdout, sectionReader)
}

// エンディアン変換
// https://qiita.com/tobira-code/items/a03f39a02678d80bbd26
func sample03() {
	// 32ビットのビッグエンディアンのデータ（10000）
	data := []byte{0x0, 0x0, 0x27, 0x10}
	fmt.Printf("data           : %d\n", data)

	fmt.Printf("------------------------------------------------\n")

	var i int32
	// エンディアンの変換
	binary.Read(bytes.NewReader(data), binary.BigEndian, &i)
	fmt.Printf("i              : %d\n", i)
	fmt.Printf("i(BigEndian)   : %#016x\n", i)

	fmt.Printf("------------------------------------------------\n")

	binary.Read(bytes.NewReader(data), binary.LittleEndian, &i)
	fmt.Printf("i              : %d\n", i)
	fmt.Printf("i(LittleEndian): %#016x\n", i)

	//fmt.Printf("%#016x\n", binary.BigEndian.Uint32(data))
	//fmt.Printf("%#016x\n", binary.LittleEndian.Uint32(data))
}

func sample04() {
	// 32ビットのビッグエンディアンのデータ（AB）
	//data := []byte("AB")
	data := []byte{0x0, 0x0, 0x41, 0x42}
	fmt.Printf("data           : %s\n", data)

	fmt.Printf("------------------------------------------------\n")

	var i int32
	// エンディアンの変換
	binary.Read(bytes.NewReader(data), binary.BigEndian, &i)
	fmt.Printf("i(BigEndian)   : %#016x <-- AB\n", i)

	fmt.Printf("------------------------------------------------\n")

	binary.Read(bytes.NewReader(data), binary.LittleEndian, &i)
	fmt.Printf("i(LittleEndian): %#016x <-- BA\n", i)
}

func dumpChunk(chunk io.Reader) {
	var length int32
	binary.Read(chunk, binary.BigEndian, &length) // 長さ
	buffer := make([]byte, 4)
	chunk.Read(buffer) // 種類
	fmt.Printf("chunk '%v' (%d bytes)\n", string(buffer), length)
	if bytes.Equal(buffer, []byte("tEXt")) {
		rawText := make([]byte, length)
		chunk.Read(rawText)
		fmt.Println(string(rawText))
	}
}

func readChunks(file *os.File) []io.Reader {
	// チャンクを格納する配列
	var chunks []io.Reader

	// 最初の8バイトはシグニチャなので飛ばす
	file.Seek(8, 0)
	// whence
	// 0: ファイルの先頭からのoffset(先頭からスキップするバイト数)
	// 1: 今のSeek位置からのoffset(前回Seekした位置からスキップするバイト数）
	// 2: ファイルの末尾からのoffset(末尾から読む。この場合は負数にしないと読めません)

	var offset int64 = 8

	for {
		var length int32 // 長さ
		// lengthの型の分だけ読み込む、int32(4バイトの符号付き整数)なので4バイト読み込む
		err := binary.Read(file, binary.BigEndian, &length) // 長さ取得
		if err == io.EOF {
			break
		}
		// int64(length) = データ
		// 12 = 長さ(4バイト) + 種類(4バイト) + CRC(4バイト)
		chunks = append(chunks, io.NewSectionReader(file, offset, int64(length)+12))
		// 次のチャンクの先頭に移動
		// 現在位置は長さを読み終わった箇所なので
		// 種類(4バイト) + データ(lengthバイト) + CRC(4バイト)先に移動
		offset, _ = file.Seek(int64(length+8), 1)
	}
	return chunks
}

// PNGファイルを分析してみる
func sample05() {
	file, err := os.Open("Lenna.png")
	if err != nil {
		panic(err)
	}
	chunks := readChunks(file)
	for _, chunk := range chunks {
		dumpChunk(chunk)
	}
}

func textChunk(text string) io.Reader {
	byteData := []byte(text)
	var buffer bytes.Buffer
	binary.Write(&buffer, binary.BigEndian, int32(len(byteData)))
	buffer.WriteString("tEXt")
	buffer.Write(byteData)
	// CRCを計算して追加
	crc := crc32.NewIEEE()
	io.WriteString(crc, "tEXt")
	binary.Write(&buffer, binary.BigEndian, crc.Sum32())
	return &buffer
}

// PNG画像に秘密のテキストを入れてみる
func sample06() {
	file, err := os.Open("Lenna.png")
	if err != nil {
		panic(err)
	}
	defer file.Close()
	newFile, err := os.Create("Lenna2.png")
	if err != nil {
		panic(err)
	}
	defer newFile.Close()

	chunks := readChunks(file)
	// シグニチャ書き込み
	io.WriteString(newFile, "\x89PNG\r\n\x1a\n")
	// 先頭に必要なIHDRチャンクを書き込み
	io.Copy(newFile, chunks[0])
	// テキストチャンクを追加
	io.Copy(newFile, textChunk("ASCII PROGRAMMING++"))
	// 残りのチャンクを追加
	for _, chunk := range chunks[1:] {
		io.Copy(newFile, chunk)
	}
}
