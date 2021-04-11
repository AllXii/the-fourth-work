package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"sync"
	"time"
)

type User struct {
	username string
	password string
	article []string
}


var (
	//储存用户信息
	Users []User
	locker sync.Mutex
	//记录当前登录的用户序号
	pos = 0
	//记录对应用户的文章数量
	kase = [100] int {0}
)

func login(c *gin.Context){
	locker.Lock()
	u := c.Request.FormValue("username")
	p := c.Request.FormValue("password")
	locker.Unlock()
	exist := false
	flag := false
	for i, item := range Users{
		if u == item.username && p == item.password{
			pos = i
			flag = true
			time.Sleep( 300 * time.Millisecond)
			c.Redirect(http.StatusMovedPermanently, "/user/home")
			break
		}else if u == item.username && p != item.password{
			exist = true
			flag = true
			break
		}
	}
	if !flag{
		c.HTML(http.StatusOK, "登录页面.tmpl", gin.H{
			"flag": 1,
		})
	}
	if exist{
		c.HTML(http.StatusOK, "登录页面.tmpl", gin.H{
			"flag": 2,
		})
	}

}

func register(c *gin.Context){
	locker.Lock()
	u := c.Request.FormValue("username")
	locker.Unlock()
	if check(u){
		c.HTML(http.StatusOK, "注册页面.tmpl", gin.H{
			"flag": 1,
		})
	}else {
		locker.Lock()
		//注册信息的存储
		Users = append(Users, User{
			username: u,
			password: c.Request.FormValue("password"),
		})
		locker.Unlock()
		time.Sleep( 300 * time.Millisecond)
		c.Redirect(http.StatusMovedPermanently, "/user/login")
	}
}

//@check
//@判断用户名是否重复
func check(username string) bool{
	for _, item := range Users{
		if item.username == username{
			return true
		}
	}
	return false
}

func publish(c *gin.Context){
	Users[pos].article = append(Users[pos].article, c.Request.FormValue("article"))
	c.HTML(http.StatusOK, "home.tmpl", gin.H{
		"flag": kase[pos],
		"username": Users[pos].username,
		"article": Users[pos].article[kase[pos]],
		"datetime": time.Now().Format("2006-1-2 15:04:05"),
	})
	kase[pos] += 1
}
func main(){
	r := gin.Default()
	r.LoadHTMLGlob("第四次作业/*")
	user := r.Group("/user")
	{
		user.GET("/home", func(c *gin.Context) {
			c.HTML(200, "home.tmpl", gin.H{
				"username": Users[pos].username,
			})
		})
		user.POST("/home", publish)


		//传递进去的flag所对应的变量，是判断弹窗的打印内容
		user.GET("/login", func(c *gin.Context) {
			c.HTML(200, "登录页面.tmpl",gin.H{
				"flag": 0,
			})
		})
		user.POST("/login", login)

		user.GET("register", func(c *gin.Context) {
			c.HTML(200, "注册页面.tmpl", gin.H{
				"flag": 0,
			})
		})
		user.POST("register", register)
	}
	//http://127.0.0.1:8080/user/login
	r.Run(":8080")
}