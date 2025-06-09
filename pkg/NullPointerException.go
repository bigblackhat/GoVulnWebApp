package pkg

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

type User struct {
	name string
}

// 模拟数据库搜索
func GetUserById(id int) interface{} {
	var user User
	if id <= 0 {
		return nil
	} else {
		user.name = "Admin"
		return user
	}
}

// 从数据库获取数据时，如果获取到的是没有初始化过的User，但没有处理报错，就会引发空指针异常
func NPE1(w http.ResponseWriter, req *http.Request) {
	id, _ := strconv.Atoi(req.FormValue("id"))
	user := GetUserById(id)
	u := user.(User)
	fmt.Fprintf(w, "Hello %s!", u.name)
}

func NPE2(w http.ResponseWriter, req *http.Request) {
	var data map[string]interface{}
	json.Unmarshal([]byte(req.FormValue("content")), &data)

	// 漏洞代码如下，直接转int，且不捕获error
	roleID := data["role_id"].(int)
	// 正确代码如下，先转string，再转int
	//roleID, _ := strconv.Atoi(data["role"].(string))
	fmt.Fprintf(w, "Hello %d!", roleID)
	//fmt.Fprintf(w, "Hello %s!", "test")
}
