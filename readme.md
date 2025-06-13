# Golang Web靶场

![GoVulnWebApp](https://socialify.git.ci/bigblackhat/GoVulnWebApp/image?description=1&font=KoHo&forks=1&issues=1&owner=1&pattern=Signal&pulls=1&stargazers=1&theme=Auto)

---

## 项目部署

下载项目后执行如下命令部署：

```shell
git clone https://github.com/bigblackhat/GoVulnWebApp.git GoVulnWebApp
cd GoVulnWebApp
go build -o server
```

**将gwva.sql导入Mysql**，然后启动项目：

```shell
./server
```

访问地址：http://127.0.0.1:8080 进入靶场。

也可以指定启动端口：`./server -port=8081`。

### 演示

项目启动：

```shell
leo:~/GoVulnWebApp (main *%) $ ./server
Go-Web, Listening on port: 8080

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

Server is running...
Visit: http://localhost:8080
```

浏览器访问：http://localhost:8080

![](img/login.png)

默认管理员账号密码：admin/123456

开始玩转靶场～

![](img/index.png)

#### 漏洞清单

* 暴力破解-Brute Force
* 跨站脚本攻击-XSS
* 跨站请求伪造-CSRF
* SQL注入
* 命令注入
* 不安全的文件读取
* 不安全的文件上传
* 不安全的文件删除
* 不安全的文件写入
* 越权漏洞
* JWT-Key硬编码
* 不安全的URL重定向
* 服务端请求伪造-SSRF
* More。。。（详见writeup）

## Todo

- [ ] 完善writeup
- [ ] 后台用户管理功能，区分权限级别，仅管理员可访问，此处留洞（越权，csrf）
- [ ] 存储型XSS
- [ ] 二次SQL注入
- [ ] 越权，未授权等
- [ ] 应用Gin框架
- [ ] Docker


- [x] 空指针异常
- [x] 反射型XSS
- [x] SQL注入
- [x] 命令注入
- [x] 任意文件读取、写、删除、上传
- [x] 服务端请求伪造（SSRF）
- [x] 服务端模板注入（SSTI）
- [x] url跳转
- [x] 权限绕过

## Version Log

### v1.6

* 增加了writeup（持续完善....）

### v1.5

* 优化了权限机制，全面接入jwt
* 代码解耦

### v1.4

* 优化鉴权
* 优化数据库与SQLi-demo
* 增加注册功能
* 增加越权漏洞
* 增加初始化数据库文件

### v1.3

* 优化服务器启动相关逻辑（模板出错时提示）
* 增加NoCmdi漏洞demo
* 增加权限机制（登入、登出、接口鉴权等）
* 增加url跳转、权限绕过等漏洞demo

### v1.2

* 更新index.html
* 增加SSTI、SSRF
* 优化服务器启动相关逻辑（端口占用时提示）

### v1.1

* 增加路径穿越类漏洞（任意文件操作）
* 增加index.html界面

### v1.0

* 基础漏洞类型：XSS、SQL注入、空指针异常、命令注入

## Thanks

感谢 @[大白哥](https://github.com/1derian) 提供的诸多idea，在他的帮助下，该项目得以丰富更多的漏洞类型