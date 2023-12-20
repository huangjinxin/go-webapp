// 导入包
package main

// 导入包
import (
	"fmt"
	"html/template"
	"net/http"
	"os"
	"strings"
)

// 定义路由
// homeHandler 处理首页请求
func homeHandler(w http.ResponseWriter, r *http.Request) {
	data, err := os.ReadFile("data.txt")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// 分割字符串为行
	lines := strings.Split(string(data), "\n")

	tmplData := struct {
		Lines []string // 使用字符串的slice
	}{
		Lines: lines,
	}

	tmpl, _ := template.ParseFiles("html/template.html")
	tmpl.Execute(w, tmplData)
}

// 定义路由
func submitHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	// 解析multipart表单
	// 设置最大内存
	err := r.ParseMultipartForm(1024)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	name := r.FormValue("name")

	fmt.Printf("接收到的姓名: %s\n", name) // 这里将打印到服务器的控制台

	// 将数据写入本地文件
	file, err := os.OpenFile("data.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer file.Close()

	if _, err := file.WriteString(name + "\n"); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// 确认信息
	fmt.Fprintf(w, "Received: %s", name)
}

// 主函数
func main() {
	// 主函数
	http.HandleFunc("/", homeHandler)         // 配置路由处理器，当访问根路径时调用homeHandler函数
	http.HandleFunc("/submit", submitHandler) // 配置路由处理器，当访问/submit路径时调用submitHandler函数
	http.ListenAndServe(":8088", nil)         // 启动HTTP服务器并监听指定端口8088，接受所有请求
}
