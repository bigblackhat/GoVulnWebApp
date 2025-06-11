package pkg

import (
	"fmt"
	"net/http"
	"os/exec"
)

func Cmdi1(w http.ResponseWriter, r *http.Request) {
	url := r.FormValue("url") // 获取用户输入的URL参数

	// 🚫 危险操作：直接拼接用户输入到系统命令
	cmd := exec.Command("sh", "-c", "curl "+url) // 未过滤的命令注入点

	// ✅️ 安全写法：将wget作为命令，url作为参数，类似预编译的效果，攻击者无法注入，因为无论攻击者输入什么，都会被当作url的一部分进行解析
	//cmd := exec.Command("curl", url)

	output, err := cmd.CombinedOutput()
	if err != nil {
		http.Error(w, "下载失败: "+err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "下载成功！输出：\n%s", output)
}

func NoCmdi(w http.ResponseWriter, r *http.Request) {
	url := r.FormValue("url")
	cmd := exec.Command("curl", url)
	output, err := cmd.CombinedOutput()
	if err != nil {
		http.Error(w, "下载失败: "+err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(w, "下载成功！输出：\n%s", output)
}
