package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"html/template"
	"net/http"
)

// 静态文件
// html页面上用到的样式、css、js、图片

func Hello(ctx *gin.Context) {
	ctx.String(http.StatusOK, "hello world")
}

type UserInfo struct {
	// 通常情况下要写的全一点
	Username string `form:"username" json:"username"`
	Password string `form:"password" json:"password"`
}

func main() {
	r := gin.Default()

	// 第一个参数是模板文件里引入的前缀
	// 第二个是具体地址
	r.Static("/xxx", "./statics")

	// gin 框架中给模板添加自定义函数
	r.SetFuncMap(template.FuncMap{
		"safe": func(s string) template.HTML {
			return template.HTML(s)
		},
	})

	// 模板解析
	//r.LoadHTMLFiles("templates/index.tmpl")

	r.LoadHTMLGlob("templates/**/*") // 加载templates下的所有文件夹下的所有文件
	r.GET("/posts/index", func(context *gin.Context) {
		// 模板渲染 ，额偶群殴
		context.HTML(http.StatusOK, "posts/index.tmpl", gin.H{
			"title": "posts/gin框架模板渲染",
		})
	})

	// 重定向
	r.GET("/ready/index", func(c *gin.Context) {
		//c.JSON(http.StatusOK, gin.H{
		//	"status": "ok",
		//})

		// 重定向
		c.Redirect(http.StatusMovedPermanently, "https://www.sogou.com")
	})

	// 路由重定向
	r.GET("/a", func(c *gin.Context) {
		// 跳转到b对应的路由函数
		c.Request.URL.Path = "/b" // 把请求的URI地址修改
		r.HandleContext(c)        // 继续后续的处理
	})

	r.GET("/b", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "bbbb",
		})
	})

	r.GET("/upload", func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, "upload/index.tmpl", nil)
	})

	// 简单上传文件的操作
	r.POST("/upload_img", func(ctx *gin.Context) {
		f, err := ctx.FormFile("f1")
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
		} else {
			// 将读取到的文件保存在本地(服务端本地)
			dst := fmt.Sprintf("./%s", f.Filename)
			err := ctx.SaveUploadedFile(f, dst)
			if err != nil {
				ctx.JSON(http.StatusInternalServerError, gin.H{
					"error": err.Error(),
				})
			}
			ctx.JSON(http.StatusOK, gin.H{
				"message": "ok",
			})
		}
	})

	r.GET("/", Hello)

	r.GET("/users/index", func(context *gin.Context) {
		// 模板渲染 没有在模板里define起名字就以文件名为准
		context.HTML(http.StatusOK, "users/index.tmpl", gin.H{
			"title": "users/ <a href='www.baidu.com'>百度</a>",
		})
	})

	r.GET("/json", func(context *gin.Context) {
		// 方法1：使用map
		//data := map[string]interface{}{
		//	"name":    "无解",
		//	"message": "hello world",
		//	"age":     18,
		//}
		// 方法2
		data := gin.H{
			"name": "无解",
			"age":  19,
		}
		context.JSON(http.StatusOK, data)
	})

	// 方法3 使用结构体
	type msg struct {
		// 小写不行，模板里会取不到，小写是私有的
		Name    string `json:"name"` // 灵活使用结构体的tag来做定制化操作
		Message string
		Age     int
	}

	r.GET("/another_json", func(context *gin.Context) {
		data := msg{
			Name:    "无解",
			Message: "你好啊",
			Age:     19,
		}
		// json的序列化
		// 默认的go语言的json模块是通过反射去进行序列化
		// 结构体变量小写，是私有的，就读取不到
		context.JSON(http.StatusOK, data)
	})

	r.GET("/web", func(c *gin.Context) {
		// 遇事不决写注释
		// 获取浏览器里请求携带的querystring 参数

		// 1. 第一种
		name := c.Query("query")
		age := c.Query("age") // 多个键值对
		// 2. 第二种
		//name := c.DefaultQuery("query", "无解") // 没有传值就默认为无解

		// 3. 第三种
		//name, ok := c.GetQuery("query") // 取不到就返回false
		//if !ok {
		//	// 取不到
		//	name = "somebody"
		//}
		c.JSON(http.StatusOK, gin.H{
			"name": name,
			"age":  age,
		})
	})

	r.GET("/login", func(c *gin.Context) {
		c.HTML(http.StatusOK, "login/login.tmpl", nil)
	})

	r.POST("/login", func(c *gin.Context) {
		//username := c.PostForm("username")
		//password := c.PostForm("password") // 取到就返回值
		//username := c.DefaultPostForm("username", "dqwdwdqw")
		//password := c.DefaultPostForm("password", "****") // 取不到就返回一个默认值

		username, ok := c.GetPostForm("username")
		if !ok {
			username = "sb"
		}
		password, _ := c.GetPostForm("password")
		c.HTML(http.StatusOK, "index/index.tmpl", gin.H{
			"Name":     username,
			"Password": password,
		})
	})

	// 这里注意理由的匹配 最好再加个前缀
	r.GET("/user/:name/:age", func(c *gin.Context) {
		// 获取路径参数 querystring
		name := c.Param("name") // 返回的是字符串
		age := c.Param("age")

		c.JSON(http.StatusOK, gin.H{
			"name": name,
			"age":  age,
		})
	})

	// gin 使用参数绑定学习
	r.GET("/user", func(c *gin.Context) {
		//username := c.Query("username")
		//password := c.Query("password")
		//u := UserInfo{
		//	username: username,
		//	password: password,
		//}
		//fmt.Printf("%v\n", u)
		var user UserInfo // 声明一个UserInfo类型的变量user
		// 能根据你请求的参数自动匹配类型获取对应的数据
		err := c.ShouldBind(&user) // 通过反射找到结构体的对应的字段
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
		} else {
			fmt.Printf("%v\n", user)
			c.JSON(http.StatusOK, gin.H{
				"message": "OK",
			})
		}
	})

	// 运行服务器 监控3000端口
	r.Run(":3000")
}
