package pkg

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"net/http"
	"text/template"
)

type user struct {
	Id         int
	Name       string
	Passwd     string
	Level      int
	Mail       string
	Age        int
	Create_at  []uint8
	Updated_at []uint8
}

var db *sql.DB

func InitDB() (db *sql.DB, err error) {
	db_config := "root:12345678@tcp(127.0.0.1:3306)/vulnapp"
	db, err = sql.Open("mysql", db_config)
	if err != nil {
		//panic(err.Error())
		return nil, err
	}

	// 添加连接池配置
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)

	err = db.Ping()
	if err != nil {
		panic(err.Error())
		return nil, err
	}
	return db, nil
}

func QueryUserById(w http.ResponseWriter, req *http.Request) {
	var err error
	db, err = InitDB()

	sqlStr := "select * from vulnapp.user where id=?"
	var u user
	err = db.QueryRow(sqlStr, 1).Scan(&u.Id, &u.Name, &u.Mail, &u.Age, &u.Passwd, &u.Create_at, &u.Updated_at)
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}
	t, _ := template.ParseFiles("./tpl/SQLi/QueryUserById.gtpl")
	err = t.Execute(w, u)
	if err != nil {
		w.Write([]byte(err.Error()))
		panic(err)
	}
}

func QueryMultiUser(w http.ResponseWriter, req *http.Request) {
	sqlStr := "select * from user where id > 0"

	var err error
	db, err = InitDB()

	rows, err := db.Query(sqlStr)
	if err != nil {
		w.Write([]byte(err.Error()))
	}
	defer rows.Close()

	for rows.Next() {
		var u user
		err = rows.Scan(&u.Id, &u.Name, &u.Mail, &u.Age, &u.Passwd, &u.Create_at, &u.Updated_at)
		if err != nil {
			w.Write([]byte(err.Error()))
		}
		fmt.Fprintf(w, "id: %d, name: %s, mail: %s, age: %s, passwd: %s, create_at: %s, update_at: %s\n", u.Id, u.Name, u.Mail, u.Age, u.Passwd, u.Create_at, u.Updated_at)
	}
}

func InsertUser(w http.ResponseWriter, req *http.Request) {
	var err error
	db, err = InitDB()
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}

	sqlStr := "insert into user(name,mail,age,passwd) values (?,?,?,?)"
	ret, er := db.Exec(sqlStr, "test", "test@mail.com", 100, "123456")
	if er != nil {
		w.Write([]byte(er.Error()))
		return
	}
	fmt.Fprintf(w, "insert success,id is %d", ret)
}

func UpdateUser(w http.ResponseWriter, req *http.Request) {
	var err error
	db, err = InitDB()
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}
	sqlStr := "update user set name=?,mail=?,age=?,passwd=? where id=?"
	ret, er := db.Exec(sqlStr, "test", "test@mail.com", 100, "123456", 1)
	if er != nil {
		w.Write([]byte(er.Error()))
		return
	}
	fmt.Fprintf(w, "update success,id is %d", ret)
}

func DeleteUser(w http.ResponseWriter, req *http.Request) {
	var err error
	db, err = InitDB()
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}
	sqlStr := "delete from user where id > ?"
	ret, er := db.Exec(sqlStr, 5)
	if er != nil {
		w.Write([]byte(er.Error()))
		return
	}
	fmt.Fprintf(w, "delete success,id is %d", ret)
}

func SQLi1(w http.ResponseWriter, req *http.Request) {
	var err error
	db, err = InitDB()
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}
	id := req.FormValue("id")
	sqlStr := fmt.Sprintf("select name from user where id = '%s'", id)
	var name string
	err = db.QueryRow(sqlStr).Scan(&name)
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}
	fmt.Fprintf(w, "name is %s", name)
}
