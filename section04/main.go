package main

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"encoding/csv"
	"fmt"
	"hash/crc32"
	"io"
	"io/ioutil"
	"os"
	"strings"
)

func main() {
	sample13()
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

// 改行／単語で区切る
var source1 = `1行目
2行目
3行目`

func sample07() {
	reader := bufio.NewReader(strings.NewReader(source1))
	for {
		line, err := reader.ReadString('\n') // 任意の文字で分割する
		fmt.Printf("%#v\n", line)
		if err == io.EOF {
			break
		}
	}
}

var source2 = `1行目 2行目 3行目`

// sample07をScannerで書き換えたコード
// 競技プログラミングで標準入力を処理する時に多用する
func sample08() {
	scanner := bufio.NewScanner(strings.NewReader(source2))
	// 分割処理を単語区切りに設定
	scanner.Split(bufio.ScanWords)
	for scanner.Scan() {
		fmt.Printf("%#v\n", scanner.Text())
	}
}

// データ型を指定して解析
var source3 = "123 1.234 1.0e4 test"

func sample09() {
	reader := strings.NewReader(source3)
	var i int
	var f, g float64
	var s string
	fmt.Fscan(reader, &i, &f, &g, &s) // スペース区切りであることを前提としている
	fmt.Printf("i=%#v f=%#v g=%#v s=%#v\n", i, f, g, s)
}

// その他の形式の決まったフォーマットの文字列の解析
var csvSource = `13101,"100  ","1000003","ﾄｳｷｮｳﾄ","ﾁﾖﾀﾞｸ","ﾋﾄﾂﾊﾞｼ(1ﾁｮｳﾒ)","東京都","千代田区","一ツ橋（１丁目）",1,0,1,0,0,0
13101,"101  ","1010003","ﾄｳｷｮｳﾄ","ﾁﾖﾀﾞｸ","ﾋﾄﾂﾊﾞｼ(2ﾁｮｳﾒ)","東京都","千代田区","一ツ橋（２丁目）",1,0,1,0,0,0
13101,"100  ","1000012","ﾄｳｷｮｳﾄ","ﾁﾖﾀﾞｸ","ﾋﾋﾞﾔｺｳｴﾝ","東京都","千代田区","日比谷公園",0,0,0,0,0,0
13101,"102  ","1020093","ﾄｳｷｮｳﾄ","ﾁﾖﾀﾞｸ","ﾋﾗｶﾜﾁｮｳ","東京都","千代田区","平河町",0,0,1,0,0,0
13101,"102  ","1020071","ﾄｳｷｮｳﾄ","ﾁﾖﾀﾞｸ","ﾌｼﾞﾐ","東京都","千代田区","富士見",0,0,1,0,0,0
`

func sample10() {
	reader := strings.NewReader(csvSource)
	csvReader := csv.NewReader(reader)
	for {
		line, err := csvReader.Read()
		if err == io.EOF {
			break
		}
		fmt.Println(line[2], line[6:9])
	}
}

// ストリームを自由に操るio.Reader／io.Writer
// 出力先を一つにまとめる
func sample11() {
	header := bytes.NewBufferString("----- HEADER -----\n")
	content := bytes.NewBufferString("Example of io.MultiReader\n")
	footer := bytes.NewBufferString("----- FOOTER -----\n")

	reader := io.MultiReader(header, content, footer)
	// すべてのreaderをつなげた出力が表示
	io.Copy(os.Stdout, reader)
}

// 出力先を分岐させる
func sample12() {
	var buffer bytes.Buffer
	reader := bytes.NewBufferString("Example of io.TeeReader\n")
	teeReader := io.TeeReader(reader, &buffer)
	// データを読み捨てる
	_, _ = ioutil.ReadAll(teeReader)

	// けどバッファに残ってる
	fmt.Println(buffer.String())
}

// io.Pipe
// https://golang.org/pkg/io/#example_Pipe
// 参考
// https://medium.com/eureka-engineering/file-uploads-in-go-with-io-pipe-75519dfa647b
// https://christina04.hatenablog.com/entry/2017/01/06/190000
// ただし「このstructをio.Readerとして渡したいなぁ」と思った時に、bytes.Bufferにエンコードすると一時的にそのデータを保持することになります。
// 小さいデータならまだ良いですが、大きいデータを扱う場合はやはり無駄にメモリを消費してしまいます。
// そこでio.Pipeを使うと特に内部バッファを保つ必要もなく、io.Readerとして渡すことが可能になります。
// https://qiita.com/m0a/items/bba395b2fc9cd160e441
// (1) bufferを準備(bytes.Buffer)
// (2) bufferに対して動画等データなどの書き込み処理
// (3) bufferをPOST処理に渡す。
// 動画のサイズが大きいとbufferが実機のメモリを超えてしまい上手く動かなくなったようです
// 変更後
// (1) のbufferの代わりにio.Pipeのwriter側を準備
// (2)の処理をgo routineとして並列実行
// (3)のPOST処理にio.PipeのReader側を渡す
// つまりPOST処理でデータが必要なタイミングではじめてfileからデータを読み込む動作になるわけです。
//
// ファイルを読み込みながらネットワークへの送信を行う挙動となるため ファイルサイズ分のバッファは不要となります
// 遅延評価というやつです
func sample13() {
	fmt.Println("---------- main E")
	r, w := io.Pipe()

	fmt.Println("---------- call goroutine E")
	go func() {
		fmt.Println("---------- goroutine E")
		fmt.Fprint(w, "some text to be read\n")
		w.Close()
		fmt.Println("---------- goroutine X")
	}()
	fmt.Println("---------- call goroutine X")

	buf := new(bytes.Buffer)
	fmt.Println("---------- Read E")
	buf.ReadFrom(r)
	fmt.Println("---------- Read X ")

	fmt.Print(buf.String())
	fmt.Println("---------- main X")
}
