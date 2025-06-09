package pkg

import (
	"html/template"
	"io"
	"net/http"
	"os"
)

func FileRead(w http.ResponseWriter, r *http.Request) {
	file := r.FormValue("path")
	// 漏洞利用示例: /read?file=../../etc/passwd
	data, err := os.ReadFile(file) // 直接使用用户输入路径
	if err != nil {
		http.Error(w, "File not found", http.StatusNotFound)
		return
	}
	w.Write(data)
}

func FileWrite(w http.ResponseWriter, r *http.Request) {
	file := "./" + r.FormValue("path")
	content := r.FormValue("content")

	// 漏洞利用示例: /write?name=../../malicious.sh&content=rm -rf /
	if err := os.WriteFile(file, []byte(content), 0644); err != nil {
		http.Error(w, "Write failed", http.StatusInternalServerError)
		return
	}
	w.Write([]byte("File created"))
}

func FileDelete(w http.ResponseWriter, r *http.Request) {
	file := "./" + r.FormValue("path")
	// 漏洞利用示例: /delete?file=../../critical-system-file
	if err := os.Remove(file); err != nil {
		http.Error(w, "Delete failed", http.StatusInternalServerError)
		return
	}
	w.Write([]byte("File deleted"))
}

func FileUpload(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		t, _ := template.ParseFiles("tpl/DirTrav/upload.html")
		t.Execute(w, nil)
	} else {
		r.ParseMultipartForm(10 << 20) // 10MB max memory

		file, header, err := r.FormFile("file")
		if err != nil {
			http.Error(w, "Invalid upload", http.StatusBadRequest)
			return
		}
		defer file.Close()

		// 漏洞点1：未验证文件类型
		// 漏洞点2：直接使用原始文件名
		dst, err := os.Create(header.Filename) // 危险：原始文件名
		if err != nil {
			http.Error(w, "Save failed", http.StatusInternalServerError)
			return
		}
		defer dst.Close()

		io.Copy(dst, file) // 漏洞点3：未限制写入路径
		w.Write([]byte("Upload success"))
	}
}
