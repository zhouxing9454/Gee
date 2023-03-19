package gee

import (
	"fmt"
	"log"
	"net/http"
	"runtime"
	"strings"
)

func Recovery() HandlerFunc {
	return func(context *Context) {
		defer func() {
			if err := recover(); err != nil {
				message := fmt.Sprintf("%s", err)
				log.Printf("%s\n\n", trace(message))
				context.Fail(http.StatusInternalServerError, "Internal Server Error")
			}
		}()
		context.Next()
	}
}

// print stack trace for debug
func trace(message string) string {
	var pcs [32]uintptr
	n := runtime.Callers(3, pcs[:]) // skip first 3 caller

	var str strings.Builder
	str.WriteString(message + "\nTraceback:")
	for _, pc := range pcs[:n] {
		fn := runtime.FuncForPC(pc)
		file, line := fn.FileLine(pc)
		str.WriteString(fmt.Sprintf("\n\t%s:%d", file, line))
	}
	return str.String()
}

//这段代码实现了一个简单的函数 trace，该函数用于输出当前调用栈的追踪信息，以及包含在该信息中的指定消息。
//具体来说，函数的输入是一个字符串 message，表示在输出中需要包含的消息。
//函数的输出是一个字符串，其中包含了消息以及当前调用栈的追踪信息。
//函数首先定义了一个长度为32的 uintptr 数组 pcs，用于存储当前调用栈的程序计数器（PC）地址。
//然后使用 runtime.Callers 函数获取当前调用栈中的程序计数器地址，从第3个调用者开始获取（因为前两个调用者是 trace 函数本身以及包含它的函数）。
//接下来，函数使用 strings.Builder 创建一个字符串构建器，将输入的 message 添加到该构建器中，并添加一个 Traceback: 的标签作为追踪信息的起始。
//然后，函数遍历刚刚获取的程序计数器地址，使用 runtime.FuncForPC 函数获取与该地址相关联的函数，
//进而获取该函数所在的文件和行号，最后将这些信息添加到字符串构建器中。
//函数最终返回字符串构建器中的字符串形式的追踪信息。
