package main

import (
	"fmt"
)

// 実行
// $ go run ./section01/main.go

// デバッガー
// Run > Debug...

// エラー
// could not launch process: debugserver or lldb-server not found: install XCode's command line tools or lldb-server
// 下記コマンドを実行して解決する
// $ xcode-select --install
func main() {
	// デバッガーを使って"Hello World!" プログラムの、
	// さらに下のレイヤのシステムコールを「見る」
	fmt.Println("Hello World!")
	// print.go内でFprintlnをラップしたPrintlnを呼び出す
	// さまざまな型を受け取るために、interface{} の可変長引数になっている
	// Fprintlnの第一引数にos.Stdoutを固定で渡している
	/* ---------------------------------------------------------------------------------------------------------------*/
	// Println formats using the default formats for its operands and writes to standard output.
	// Spaces are always added between operands and a newline is appended.
	// It returns the number of bytes written and any write error encountered.
	// func Println(a ...interface{}) (n int, err error) {
	// return Fprintln(os.Stdout, a...)
	// }

	// w = os.Stdout(標準出力)
	/* ---------------------------------------------------------------------------------------------------------------*/
	// Fprintln formats using the default formats for its operands and writes to w.
	// Spaces are always added between operands and a newline is appended.
	// It returns the number of bytes written and any write error encountered.
	// func Fprintln(w io.Writer, a ...interface{}) (n int, err error) {
	// p := newPrinter()		// 文字列をフォーマット文字列に従って整形するプリンターを作成
	// p.doPrintln(a)			// 文字列を生成
	// n, err = w.Write(p.buf)	//
	// p.free()					// プリンターの解放
	// return					// 定義した返り値を返す = n, err
	// }

	// File.writeを呼び出す
	/* ---------------------------------------------------------------------------------------------------------------*/
	// Write writes len(b) bytes to the File.
	// It returns the number of bytes written and an error, if any.
	// Write returns a non-nil error when n != len(b).
	//func (f *File) Write(b []byte) (n int, err error) {
	//	if err := f.checkValid("write"); err != nil {
	//		return 0, err
	//	}
	//	n, e := f.write(b)
	//	if n < 0 {
	//		n = 0
	//	}
	//	if n != len(b) {
	//		err = io.ErrShortWrite
	//	}
	//
	//	epipecheck(f, e)
	//
	//	if e != nil {
	//		err = f.wrapErr("write", e)
	//	}
	//
	//	return n, err
	//}

	// file_unix.goなのでUnix系OS固有のコードに飛んでいる
	// ファイルディスクリプタのWriteを呼び出す
	/* ---------------------------------------------------------------------------------------------------------------*/
	// write writes len(b) bytes to the File.
	// It returns the number of bytes written and an error, if any.
	//func (f *File) write(b []byte) (n int, err error) {
	//	n, err = f.pfd.Write(b)
	//	runtime.KeepAlive(f)
	//	return n, err
	//}

	// syscall.Write = システムコール
	/* ---------------------------------------------------------------------------------------------------------------*/
	// Write implements io.Writer.
	//func (fd *FD) Write(p []byte) (int, error) {
	//	if err := fd.writeLock(); err != nil {
	//		return 0, err
	//	}
	//	defer fd.writeUnlock()
	//	if err := fd.pd.prepareWrite(fd.isFile); err != nil {
	//		return 0, err
	//	}
	//	var nn int
	//	for {
	//		max := len(p)
	//		if fd.IsStream && max-nn > maxRW {
	//			max = nn + maxRW
	//		}
	//		n, err := syscall.Write(fd.Sysfd, p[nn:max])
	//		if n > 0 {
	//			nn += n
	//		}
	//		if nn == len(p) {
	//			return nn, err
	//		}
	//		if err == syscall.EAGAIN && fd.pd.pollable() {
	//			if err = fd.pd.waitWrite(fd.isFile); err == nil {
	//				continue
	//			}
	//		}
	//		if err != nil {
	//			return nn, err
	//		}
	//		if n == 0 {
	//			return nn, io.ErrUnexpectedEOF
	//		}
	//	}
	//}
}
