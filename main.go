package main

import (
	"GolandWeb/pkg"
	"database/sql"
	"flag"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/golang-jwt/jwt"
	_ "github.com/golang-jwt/jwt"
	"html/template"
	"net"
	"net/http"
	"os"
	"strconv"
	"time"
)

func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		c, err := r.Cookie("username")
		if err != nil {
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}

		value := c.Value
		if value == "" {
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}
		next(w, r)
	}
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("tpl/index.html")
	//t, err := template.ParseFiles("test/tmp.html")
	if err != nil {
		w.Write([]byte(err.Error()))
		return // 确保返回
	}
	t.Execute(w, nil)
}

func generateJWT(username string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user": username,
		"exp":  time.Now().Add(24 * time.Hour).Unix(),
	})
	return token.SignedString([]byte("your-secret-key"))
}

func registerHandler(w http.ResponseWriter, r *http.Request) {
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

func loginHandler(w http.ResponseWriter, r *http.Request) {
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
			// 创建 Cookie 对象
			cookie := &http.Cookie{
				Name:   "username",
				Value:  name,
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

func logoutHandler(w http.ResponseWriter, r *http.Request) {
	delCookie := &http.Cookie{Name: "username", MaxAge: -1}
	http.SetCookie(w, delCookie)
	files, err := template.ParseFiles("tpl/logout.html")
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}
	files.Execute(w, nil)
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	t, _ := template.New("home").Parse("<h1>Hello, It's Home, {{.Name}}</h1>")
	t.Execute(w, map[string]string{"Name": "Leo's Go Web"})
}

func mysqlTestHandler(w http.ResponseWriter, r *http.Request) {
	db_config := "root:12345678@tcp(127.0.0.1:3306)/mysql"
	db, err := sql.Open("mysql", db_config)
	if err != nil {
		panic(err.Error())
	}
	err = db.Ping()
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()
	w.Write([]byte("Mysql Connect Success~<br>" +
		"DB Config is: " + db_config + "<br>"))
}

var banner = `
       _==/          i     i          \==_
     /XX/            |\___/|            \XX\
   /XXXX\            |XXXXX|            /XXXX\
  |XXXXXX\_         _XXXXXXX_         _/XXXXXX|
 XXXXXXXXXXXxxxxxxxXXXXXXXXXXXxxxxxxxXXXXXXXXXXX
|XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX|
XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX
|XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX|
 XXXXXX/^^^^"\XXXXXXXXXXXXXXXXXXXXX/^^^^^\XXXXXX
  |XXX|       \XXX/^^\XXXXX/^^\XXX/       |XXX|
    \XX\       \X/    \XXX/    \X/       /XX/
       "\       "      \X/      "      /"
------------------------------------------------
`

// 检查端口是否被占用
func isPortOccupied(port int) bool {
	addr := ":" + strconv.Itoa(port)
	listener, err := net.Listen("tcp", addr)

	if err != nil {
		// 检查特定错误类型
		if opErr, ok := err.(*net.OpError); ok {
			if opErr.Op == "listen" {
				return true
			}
		}
		return false
	}

	// 端口可用，立即关闭
	listener.Close()
	return false
}

func WebServerRun(port int) {
	// 检查端口是否可用
	if isPortOccupied(port) {
		fmt.Printf("Error: Port %d is already in use. Please choose a different port.\n", port)
		os.Exit(1)
	}

	fmt.Println("Go-Web, Listening on port:", port)
	fmt.Println(banner)

	http.HandleFunc("/index", AuthMiddleware(indexHandler))
	http.HandleFunc("/", AuthMiddleware(indexHandler))
	http.HandleFunc("/login", loginHandler)
	http.HandleFunc("/register", registerHandler)
	http.HandleFunc("/logout", logoutHandler)

	// 测试handler
	http.HandleFunc("/home", homeHandler)
	http.HandleFunc("/QueryUserById", pkg.QueryUserById)
	http.HandleFunc("/QueryMultiUser", pkg.QueryMultiUser)
	http.HandleFunc("/InsertUser", pkg.InsertUser)
	http.HandleFunc("/UpdateUser", pkg.UpdateUser)
	http.HandleFunc("/DeleteUser", pkg.DeleteUser)

	// 以下是靶场路由
	// 跨站脚本攻击
	http.HandleFunc("/XSS/reflectXSS1", AuthMiddleware(pkg.ReflectXSS1))
	http.HandleFunc("/XSS/reflectXSS2", AuthMiddleware(pkg.ReflectXSS2))
	http.HandleFunc("/XSS/reflectXSS3", AuthMiddleware(pkg.ReflectXSS3))
	http.HandleFunc("/XSS/reflectXSS4", AuthMiddleware(pkg.ReflectXSS4))
	http.HandleFunc("/XSS/NoXSS1", AuthMiddleware(pkg.NoXSS1))

	// SQL注入
	http.HandleFunc("/SQLi/SQLi1", AuthMiddleware(pkg.SQLi1))

	// 空指针异常
	http.HandleFunc("/NullPointerException/NPE1", AuthMiddleware(pkg.NPE1))
	http.HandleFunc("/NullPointerException/NPE2", AuthMiddleware(pkg.NPE2))

	// 命令注入
	http.HandleFunc("/CommandInjection/cmdi1", AuthMiddleware(pkg.Cmdi1))
	http.HandleFunc("/CommandInjection/nocmd", AuthMiddleware(pkg.NoCmdi))

	// 路径穿越
	http.HandleFunc("/DirTraversal/FileRead", AuthMiddleware(pkg.FileRead))
	http.HandleFunc("/DirTraversal/FileWrite", AuthMiddleware(pkg.FileWrite))
	http.HandleFunc("/DirTraversal/FileDelete", AuthMiddleware(pkg.FileDelete))
	http.HandleFunc("/DirTraversal/FileUpload", AuthMiddleware(pkg.FileUpload))

	// SSRF
	http.HandleFunc("/SSRF/SSRF1", AuthMiddleware(pkg.SSRF1))
	http.HandleFunc("/SSRF/NoSSRF", AuthMiddleware(pkg.NoSSRF))

	// SSTI
	http.HandleFunc("/SSTI/SSTI1", AuthMiddleware(pkg.SSTI1))
	http.HandleFunc("/SSTI/SSTI2", AuthMiddleware(pkg.SSTI2))

	// URL跳转
	http.HandleFunc("/Redirect/Redirect1", AuthMiddleware(pkg.Redirect1))
	http.HandleFunc("/Redirect/Redirect2", AuthMiddleware(pkg.Redirect2))
	http.HandleFunc("/Redirect/Redirect3", AuthMiddleware(pkg.Redirect3))

	// 创建自定义HTTP服务器实例，直接使用已检测过的listener
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		fmt.Printf("Unexpected error after port check: %v\n", err)
		os.Exit(1)
	}
	// 启动服务器
	fmt.Println("Server is running...")
	fmt.Printf("Visit: http://localhost:%d\n", port)

	if err := http.Serve(listener, nil); err != nil {
		fmt.Printf("Server failed to start: %v\n", err)
		os.Exit(1)
	}
	// 启动HTTP服务器，监听8080端口
	//http.ListenAndServe(":8080", nil)
}

func main() {
	// 1. 解析命令行端口参数
	port := flag.Int("port", 8080, "HTTP server port")
	flag.Parse()

	// 2. 启动带端口检查的Web服务器
	WebServerRun(*port)
}
