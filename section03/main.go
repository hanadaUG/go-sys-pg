package main

import (
	"archive/zip"
	"crypto/rand"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

//「バイト列 b を書き込み、書き込んだバイト数 n と、エラーが起きた場合はそのエラー error を返す」という振る舞いを
// インタフェース という型で定義している
// プログラムで外部からデータを読み込むための機能もGo言語のインタフェースとして抽象化されています。
//type Read interface{
//Read(b []byte) (n int, err error)
//}
func main() {
	//q1()
	//q2()
	//q3()
	//q4()
}

//
func sample01() {
	// io.Readerインターフェースを満たす構造体
	r := os.Stdin
	// 1024バイトのバッファをmakeで作る
	buffer := make([]byte, 1024)
	// sizeは実際に読み込んだバイト数、errはエラー
	size, err := r.Read(buffer)
	fmt.Print(size, err)

	// バッファの管理をしつつ、何度も Read メソッドを読んでデータを最後まで読み込むなど、
	// 読み込み処理を書くたびに同じようなコードを書かなければなりません。
}

// 読み込みの補助関数1
func sample02() {
	reader, _ := os.Open("test.txt")
	buffer, _ := ioutil.ReadAll(reader) // 終端記号に当たるまですべてのデータを読み込んで返す
	fmt.Print(buffer)
	fmt.Print("\n", string(buffer))
}

// 読み込みの補助関数2
func sample03() {
	reader, _ := os.Open("test.txt")

	buffer := make([]byte, 4)
	io.ReadFull(reader, buffer) // 4バイト読み込めないとエラー
	fmt.Print(buffer)
	fmt.Print("\n", string(buffer))
}

// コピーの補助関数
// section02/main.goで触れたのでスキップ
// バッファの確保、読み込み、書き込むという手順をio.Copyがよしなに処理してくれる

// 標準入力
func sample04() {
	for {
		buffer := make([]byte, 5)
		size, err := os.Stdin.Read(buffer)
		if err == io.EOF {
			fmt.Println("EOF")
			break
		}
		fmt.Printf("size=%d input='%s'\n", size, string(buffer))
	}
}

// Q1. ファイルのコピー
func q1() {
	oldFile, err := os.Open("old.txt")
	defer oldFile.Close()
	if err != nil {
		panic(err)
	}
	oldFile.Write([]byte("old\n"))

	newFile, err := os.Create("new.txt")
	defer newFile.Close()
	if err != nil {
		panic(err)
	}
	io.Copy(newFile, oldFile)
}

// Q2. テスト用の適当なサイズのファイルを作成
func q2() {
	buffer := make([]byte, 1024)
	rand.Reader.Read(buffer)
	file, err := os.Create("rand.txt")
	defer file.Close()
	if err != nil {
		panic(err)
	}
	file.Write(buffer)

	// io.CopyN(file, rand.Reader, 1024)
}

// Q3. zipファイルの書き込み
func q3() {
	file, _ := os.Create("test.zip")
	defer file.Close()
	zipWriter := zip.NewWriter(file) // zip書き込み用の構造体、io.Writerではない
	defer zipWriter.Close()

	writer, _ := zipWriter.Create("new.txt") // Createでio.Writerを返す、zipに書き込むファイル名を指定しただけ
	//newFile, _ := os.Open("new.txt")			// 書き込む内容は別途読み込む必要がある
	//defer newFile.Close()

	//io.Copy(writer,newFile)

	io.Copy(writer, strings.NewReader("test"))
}

// Q4. zipファイルをウェブサーバからダウンロード
func q4() {
	http.HandleFunc("/", handler)
	http.ListenAndServe(":8080", nil)
}

func handler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/zip")
	w.Header().Set("Content-Disposition", "attachment; filename=sample.zip")
	zipWriter := zip.NewWriter(w) // zip書き込み用の構造体、書き込み先をhttp.ResponseWriterにしている
	defer zipWriter.Close()

	writer, _ := zipWriter.Create("aaa.txt")
	io.Copy(writer, strings.NewReader("zip"))
}
