package main

import (
	"errors"
	"html/template"
	"log"
	"net/http"
	"os"
	"regexp"
)

// 每次渲染页面调用ParseFiles并不高效 因此创建全局变量
// template.Must 是一个方便的包装器 当传入非 nil 的 error 值时会引发 panic 否则会原样返回 *Template
var templates = template.Must(template.ParseFiles("edit.html", "view.html"))

// 验证地址栏输入避免文件泄漏
// 正则用括号包裹是为了创建捕获组 便于后面提取
var validPath = regexp.MustCompile("^/(edit|save|view)/([a-zA-Z0-9]+)$")

type Page struct {
	Title string
	Body  []byte // io 库使用的格式
}

// 持久存储页面数据
func (p *Page) save() error {
	filename := p.Title + ".txt"
	// 第三个参数是文件权限 0600 表示该文件应仅以当前用户的读写权限创建
	return os.WriteFile(filename, p.Body, 0600)
}

// 加载页面
func loadPage(title string) (*Page, error) {
	filename := title + ".txt"
	body, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return &Page{Title: title, Body: body}, nil
}

// http.Request 是结构体 而结构体是值类型 因此需要指针引用
// http.ResponseWriter 是interface接口 接口是引用类型 不需要再用指针
// 值类型：基础数据类型、结构体、固定长度的数组
// 自带指针的容器：interface接口、切片、map、channel、func
func viewHandler(w http.ResponseWriter, r *http.Request) {
	// 旧
	// // 得到的path比如是/view/test len("/view/")就是去掉/view/部分
	// title := r.URL.Path[len("/view/"):]

	// 新
	title, err := getTitle(w, r)
	if err != nil {
		return
	}

	p, err := loadPage(title)
	if err != nil {
		// 请求的页面不存在 应该重定向到编辑页面
		http.Redirect(w, r, "/edit/"+title, http.StatusFound)
		return
	}

	// 旧
	// fmt.Fprintf(w, "<h1>%s</h1><div>%s</div>", p.Title, p.Body)

	// 新
	renderTemplate(w, "view", p)
}
func editHandler(w http.ResponseWriter, r *http.Request) {
	// 旧
	// title := r.URL.Path[len("/edit/"):]

	// 新
	title, err := getTitle(w, r)
	if err != nil {
		return
	}

	p, err := loadPage(title)
	if err != nil {
		p = &Page{Title: title} // loadPage 返回的是*Page类型 因此这里赋值修改也必须是指针
	}

	// 旧 这样硬编码的 HTML 很难看 所以使用 html/template 包
	// fmt.Fprintf(w, "<h1>Editing %s</h1>"+
	// 	"<form action=\"/save/%s\" method=\"POST\">"+
	// 	"<textarea name=\"body\">%s</textarea><br>"+
	// 	"<input type=\"submit\" value=\"Save\">"+
	// 	"</form>",
	// 	p.Title, p.Title, p.Body)

	// 新
	renderTemplate(w, "edit", p)
}
func renderTemplate(w http.ResponseWriter, tmpl string, p *Page) {
	// 旧
	// t, err := template.ParseFiles(tmpl + ".html") // 读取html文件内容并返回*template.Template
	// if err != nil {
	// 	http.Error(w, err.Error(), http.StatusInternalServerError)
	// 	return
	// }
	// // Execute 执行模板 将生成的 HTML 写入 http.ResponseWriter
	// err = t.Execute(w, p)

	// 新
	err := templates.ExecuteTemplate(w, tmpl+".html", p)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
func saveHandler(w http.ResponseWriter, r *http.Request) {
	// 旧
	// title := r.URL.Path[len("/save/"):]

	// 新
	title, err := getTitle(w, r)
	if err != nil {
		return
	}

	// 从表单数据中取name="body"的值
	body := r.FormValue("body")
	p := &Page{Title: title, Body: []byte(body)}
	// 注意！这里不能用 := 是因为前面已经声明过同名的
	// 而其他方法中可以用 p, err := 是因为含有未声明变量
	err = p.save() // 写入文件
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/view/"+title, http.StatusFound) // 重定向到 /view/ 页面
}
func getTitle(w http.ResponseWriter, r *http.Request) (string, error) {
	// 匹配字符串 返回所有捕获组 以字符串切片的形式返回
	m := validPath.FindStringSubmatch(r.URL.Path)
	if m == nil {
		http.NotFound(w, r)
		return "", errors.New("无效路径")
	}
	// m[0] 是整个匹配的字符串
	return m[2], nil
}

func main() {
	// p1 := &Page{Title: "TestPage", Body: []byte("This is a sample Page.")}
	// p1.save()
	// p2, _ := loadPage("TestPage")
	// fmt.Println(string(p2.Body))

	http.HandleFunc("/view/", viewHandler) // 写成/view 会导致/view/abc无法匹配 只能匹配到/view
	http.HandleFunc("/edit/", editHandler)
	http.HandleFunc("/save/", saveHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
