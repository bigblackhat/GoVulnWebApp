package main

import (
	"GolandWeb/lib"
	"GolandWeb/pkg"
	"database/sql"
	"flag"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/golang-jwt/jwt"
	"net"
	"net/http"
	"os"
	"strconv"
)

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

	http.HandleFunc("/index", lib.AuthMiddleware(lib.IndexHandler))
	http.HandleFunc("/", lib.AuthMiddleware(lib.IndexHandler))
	http.HandleFunc("/login", lib.LoginHandler)
	http.HandleFunc("/register", lib.RegisterHandler)
	http.HandleFunc("/logout", lib.LogoutHandler)

	// 测试handler
	http.HandleFunc("/home", lib.HomeHandler)
	http.HandleFunc("/QueryUserById", pkg.QueryUserById)
	http.HandleFunc("/QueryMultiUser", pkg.QueryMultiUser)
	http.HandleFunc("/InsertUser", pkg.InsertUser)
	http.HandleFunc("/UpdateUser", pkg.UpdateUser)
	http.HandleFunc("/DeleteUser", pkg.DeleteUser)

	// 以下是靶场路由
	// 跨站脚本攻击
	http.HandleFunc("/XSS/reflectXSS1", lib.AuthMiddleware(pkg.ReflectXSS1))
	http.HandleFunc("/XSS/reflectXSS2", lib.AuthMiddleware(pkg.ReflectXSS2))
	http.HandleFunc("/XSS/reflectXSS3", lib.AuthMiddleware(pkg.ReflectXSS3))
	http.HandleFunc("/XSS/reflectXSS4", lib.AuthMiddleware(pkg.ReflectXSS4))
	http.HandleFunc("/XSS/NoXSS1", lib.AuthMiddleware(pkg.NoXSS1))

	// SQL注入
	http.HandleFunc("/SQLi/SQLi1", lib.AuthMiddleware(pkg.SQLi1))

	// 空指针异常
	http.HandleFunc("/NullPointerException/NPE1", lib.AuthMiddleware(pkg.NPE1))
	http.HandleFunc("/NullPointerException/NPE2", lib.AuthMiddleware(pkg.NPE2))

	// 命令注入
	http.HandleFunc("/CommandInjection/cmdi1", lib.AuthMiddleware(pkg.Cmdi1))
	http.HandleFunc("/CommandInjection/nocmd", lib.AuthMiddleware(pkg.NoCmdi))

	// 路径穿越
	http.HandleFunc("/DirTraversal/FileRead", lib.AuthMiddleware(pkg.FileRead))
	http.HandleFunc("/DirTraversal/FileWrite", lib.AuthMiddleware(pkg.FileWrite))
	http.HandleFunc("/DirTraversal/FileDelete", lib.AuthMiddleware(pkg.FileDelete))
	http.HandleFunc("/DirTraversal/FileUpload", lib.AuthMiddleware(pkg.FileUpload))

	// SSRF
	http.HandleFunc("/SSRF/SSRF1", lib.AuthMiddleware(pkg.SSRF1))
	http.HandleFunc("/SSRF/NoSSRF", lib.AuthMiddleware(pkg.NoSSRF))

	// SSTI
	http.HandleFunc("/SSTI/SSTI1", lib.AuthMiddleware(pkg.SSTI1))
	http.HandleFunc("/SSTI/SSTI2", lib.AuthMiddleware(pkg.SSTI2))

	// URL跳转
	http.HandleFunc("/Redirect/Redirect1", lib.AuthMiddleware(pkg.Redirect1))
	http.HandleFunc("/Redirect/Redirect2", lib.AuthMiddleware(pkg.Redirect2))
	http.HandleFunc("/Redirect/Redirect3", lib.AuthMiddleware(pkg.Redirect3))

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
	//test.Test1()
}
