package main

import (
	"bufio"
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
)

//「バイト列 p を書き込み、書き込んだバイト数 n と、エラーが起きた場合はそのエラー error を返す」という振る舞いを
// インタフェース という型で定義している
// いろいろなものを「ファイル」のように扱うために、システムコールではファイルディスクリプタ、Goではインタフェースを使って抽象化している
//type Writer interface{
//Write(p []byte) (n int, err error)
//}
func main() {
	fmt.Println("Hello World!")
}

// ファイル出力
func fileWrite() {
	file, err := os.Create("test.txt")
	if err != nil {
		panic(err)
	}
	file.Write([]byte("os.File example\n"))
	file.Close()
}

// 標準出力
func stdWrite() {
	os.Stdout.Write([]byte("os.Stdout example\n"))
}

// バッファ
func bufferWrite() {
	var buffer bytes.Buffer
	buffer.Write([]byte("bytes.Buffer example\n")) // バッファ書にき込み = バッファに蓄積、あとで纏めて結果を受け取る
	fmt.Println(buffer.String())
}

// コネクション
func connectionWrite() {
	// Connはio.Readerとio.Writerの両方のinterfaceを満たしている
	conn, err := net.Dial("tcp", "ascii.jp:80")
	if err != nil {
		panic(err)
	}
	conn.Write([]byte("GET / HTTP/1.0\r\nHost: ascii.jp\r\n\r\n"))
	// io.Copyはio.Readerからio.Writerにデータを渡すとき利用する
	// https://ascii.jp/elem/000/001/252/1252961/
	// 第3回   低レベルアクセスへの入り口（2）：io.Reader前編で説明される
	io.Copy(os.Stdout, conn) // 第一引数：コピー先、第二引数：コピー元
}

// indexにアクセスした場合の振る舞いを定義している
func handler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("http.ResponseWriter sample"))
}

// indexにアクセスした場合の振る舞いを登録しサーバーを起動する
func server() {
	http.HandleFunc("/", handler)
	http.ListenAndServe(":8080", nil)
}

// フィルタ
func multiWrite() {
	file, err := os.Create("multiwriter.txt")
	if err != nil {
		panic(err)
	}
	// file.Write([]byte("os.File example\n"))
	// os.Stdout.Write([]byte("os.Stdout example\n"))
	writer := io.MultiWriter(file, os.Stdout)          // 書き出し先を指定する
	io.WriteString(writer, "io.MultiWriter example\n") // ファイル、標準出力の両方に同時に書き出す
}

// 圧縮
func gzipWrite() {
	file, err := os.Create("test.txt.gz")
	if err != nil {
		panic(err)

	}
	writer := gzip.NewWriter(file)
	writer.Header.Name = "test.txt"
	writer.Write([]byte("gzip.Writer example\n"))
	writer.Close()
}

// 出力結果を一時的にためておいて、ある程度の分量ごとにまとめて書き出す
// 参考：https://www.mas9612.net/posts/golang-bufio-writer/
func bufioWrite() {
	buffer := bufio.NewWriterSize(os.Stdout, 8)
	buffer.WriteString("123456")
	//buffer.Flush()
	buffer.WriteString("abc\n")
	buffer.Flush()
}

// json
func jsonWrite() {
	encoder := json.NewEncoder(os.Stdout) // Encoderに渡すio.Writerを変更することで出力先をサーバーやブラウザに変更できる
	encoder.SetIndent("", "    ")
	encoder.Encode(map[string]string{
		"example": "encoding/json",
		"hello":   "world",
	})
}

// net/httpパッケージのRequest構造体 = HTTPリクエストを取り扱う構造体
func requestWrite() {
	request, err := http.NewRequest("GET", "http://ascii.jp", nil)
	if err != nil {
		panic(err)
	}
	request.Header.Set("X-TEST", "ヘッダーも追加できます")
	request.Write(os.Stdout)
}
