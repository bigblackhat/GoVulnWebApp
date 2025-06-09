package main

import (
	"GolandWeb/pkg"
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"html/template"
	"net/http"
)

func indexHandler(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("tpl/index.html")
	t.Execute(w, nil)
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

func WebServerRun() {
	fmt.Println("Go-Web, Listening on port 8080")
	fmt.Println("       _==/          i     i          \\==_\n     /XX/            |\\___/|            \\XX\\\n   /XXXX\\            |XXXXX|            /XXXX\\\n  |XXXXXX\\_         _XXXXXXX_         _/XXXXXX|\n XXXXXXXXXXXxxxxxxxXXXXXXXXXXXxxxxxxxXXXXXXXXXXX\n|XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX|\nXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX\n|XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX|\n XXXXXX/^^^^\"\\XXXXXXXXXXXXXXXXXXXXX/^^^^^\\XXXXXX\n  |XXX|       \\XXX/^^\\XXXXX/^^\\XXX/       |XXX|\n    \\XX\\       \\X/    \\XXX/    \\X/       /XX/\n       \"\\       \"      \\X/      \"      /\"\n\n------------------------------------------------\n")

	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/home", homeHandler)
	http.HandleFunc("/QueryUserById", pkg.QueryUserById)
	http.HandleFunc("/QueryMultiUser", pkg.QueryMultiUser)
	http.HandleFunc("/InsertUser", pkg.InsertUser)
	http.HandleFunc("/UpdateUser", pkg.UpdateUser)
	http.HandleFunc("/DeleteUser", pkg.DeleteUser)

	// 以下是靶场路由
	// 跨站脚本攻击
	http.HandleFunc("/XSS/reflectXSS1", pkg.ReflectXSS1)
	http.HandleFunc("/XSS/reflectXSS2", pkg.ReflectXSS2)
	http.HandleFunc("/XSS/reflectXSS3", pkg.ReflectXSS3)
	http.HandleFunc("/XSS/reflectXSS4", pkg.ReflectXSS4)
	http.HandleFunc("/XSS/NoXSS1", pkg.NoXSS1)

	// SQL注入
	http.HandleFunc("/SQLi/SQLi1", pkg.SQLi1)

	// 空指针异常
	http.HandleFunc("/NullPointerException/NPE1", pkg.NPE1)
	http.HandleFunc("/NullPointerException/NPE2", pkg.NPE2)

	// 命令注入
	http.HandleFunc("/CommandInjection/cmdi1", pkg.Cmdi1)

	// 路径穿越
	http.HandleFunc("/DirTraversal/FileRead", pkg.FileRead)
	http.HandleFunc("/DirTraversal/FileWrite", pkg.FileWrite)
	http.HandleFunc("/DirTraversal/FileDelete", pkg.FileDelete)
	http.HandleFunc("/DirTraversal/FileUpload", pkg.FileUpload)

	// 启动HTTP服务器，监听8080端口
	http.ListenAndServe(":8080", nil)
}

func main() {
	WebServerRun()
}
