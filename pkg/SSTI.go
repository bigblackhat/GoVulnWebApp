package pkg

import (
	"fmt"
	ht "html/template"
	"net/http"
	"os/exec"
	"strings"
)

/*
场景1，模板为用户完全可控

- 输出任意内容：
	temp={{%22Hello%22}}
- xss效果：
	temp={{define%20"T1"}}<script>alert(1)</script>{{end}}%20{{template%20"T1"}}
- 获取环境变量：
	第一步，代码中通过parse.Execute将变量传给模板中渲染，以下面demo举例，将http.Request传给了模板
	第二步，设计payload，获取http.Request的值：
		temp={{.URL}}
*/

func SSTI1(w http.ResponseWriter, r *http.Request) {
	temp := r.URL.Query().Get("temp")
	if temp == "" {
		http.Error(w, "temp parameter is empty", http.StatusBadRequest)
		return
	}
	parse, err := ht.New("ssti1").Parse(temp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	parse.Execute(w, r)
}

/*
场景2，模板内容部分为用户可控

- 输出任意内容：
	name={{printf%20"1"}}
但无法如场景1那样实现XSS效果，因为那个XSS payload会导致模板被重复定义，引发报错
	name=</h1>{{define "T2"}}<script>alert(1)</script>{{end}} {{template "T2"}}<!--

- 获取环境变量：
	name={{.Name}},{{.Password}}      // 输出指定数据内容
	name={{.}}                        // 输出全部数据内容

- 调用函数，执行命令：
	name={{.Secret%20"whoami"}}
*/

type manage struct {
	Id       int
	Name     string
	Password string
}

func (m manage) Secret(test string) string {
	output, _ := exec.Command(test).CombinedOutput()
	return string(output)
}
func (m manage) Test1() string {
	s := "111"
	return "test success" + s
}
func (m manage) Test2(id string, tt string) string {
	s := "aa" + id + tt
	return "test2 success" + s
}

func SSTI2(w http.ResponseWriter, r *http.Request) {
	mng := &manage{1, "admin", "123456"}
	arg := strings.Join(r.URL.Query()["name"], "")
	tpl1 := fmt.Sprintf(`<h1>Hi, ` + arg + `</h1> Your name is ` + arg + `!`)
	html, err := ht.New("login").Parse(tpl1)
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}
	ht.Must(html, err)
	html.Execute(w, mng)
}
