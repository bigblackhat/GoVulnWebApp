package lib

import (
	"GolandWeb/pkg"
	"html/template"
	"net/http"
	"strconv"
	"time"
)

var secretKey string = "your-256-bit-very-secure-secret-key"

func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		c, err := r.Cookie("token")
		if err != nil {
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}

		value := c.Value
		username, expiration, err := ParseJWT(value, secretKey)

		if username == "" || time.Now().After(expiration) {
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}

		next(w, r)
	}
}

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("tpl/index.html")
	//t, err := template.ParseFiles("test/tmp.html")
	if err != nil {
		w.Write([]byte(err.Error()))
		return // 确保返回
	}
	t.Execute(w, nil)
}

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		t, err := template.ParseFiles("tpl/register.html")
		if err != nil {
			w.Write([]byte(err.Error()))
			return
		}
		t.Execute(w, nil)
	} else {
		username := r.FormValue("username")
		password := r.FormValue("password")
		var level string
		level = r.FormValue(level)
		if level == "" {
			level = "2"
		}
		levelInt, _ := strconv.Atoi(level)

		email := r.FormValue("email")

		if username == "" || password == "" {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("账号或密码不能为空"))
			return
		}

		var err error
		db, err := pkg.InitDB()
		sql := "select name from manager where name =?"
		var name string
		_ = db.QueryRow(sql, username).Scan(&name) // 抛弃error
		if name != "" {
			w.Write([]byte("账号已存在"))
			return
		}

		sql = "insert into manager(name,passwd,level,mail) values(?,?,?,?)"
		_, err = db.Exec(sql, username, password, levelInt, email)
		if err != nil {
			w.Write([]byte(err.Error()))
			return
		}
		//w.Write([]byte("注册成功"))
		http.Redirect(w, r, "/login", http.StatusFound)
	}
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		t, err := template.ParseFiles("tpl/login.html")
		if err != nil {
			w.Write([]byte(err.Error()))
			return
		}
		t.Execute(w, nil)
	} else {
		username := r.FormValue("username")
		password := r.FormValue("password")
		if username == "" || password == "" {
			w.Write([]byte("username or password is empty"))
			return
		}
		db, err := pkg.InitDB()
		sql := "select name from manager where name =? and passwd=?"
		var name string
		err = db.QueryRow(sql, username, password).Scan(&name)
		if err != nil {
			w.Write([]byte(err.Error()))
			return
		}
		if name != "" {
			expirationDuration := 24 * time.Hour
			jwt, err := GenerateJWT(name, secretKey, expirationDuration)
			if err != nil {
				w.Write([]byte(err.Error()))
			}
			// 创建 Cookie 对象
			cookie := &http.Cookie{
				Name:   "token",
				Value:  jwt,
				MaxAge: 3600, // 有效期（秒）
			}
			// 添加到响应头
			http.SetCookie(w, cookie)
			http.Redirect(w, r, "/", http.StatusFound)
			w.WriteHeader(http.StatusOK)
			return
		}
	}
}

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	delCookie := &http.Cookie{Name: "token", MaxAge: -1}
	http.SetCookie(w, delCookie)
	files, err := template.ParseFiles("tpl/logout.html")
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}
	files.Execute(w, nil)
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	t, _ := template.New("home").Parse("<h1>Hello, It's Home, {{.Name}}</h1>")
	t.Execute(w, map[string]string{"Name": "Leo's Go Web"})
}
