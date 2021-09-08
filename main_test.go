/*
 * @Author: wxvirus
 * @Date: 2021-08-11 23:26:53
 * @LastEditTime: 2021-08-25 00:05:22
 * @LastEditors: Please set LastEditors
 * @Description: In User Settings Edit
 * @FilePath: /gin_demo/main_test.go
 */
package main

import (
	"fmt"
	"html/template"
	"net/http"
)

type UserInfo struct {
	Name   string
	Gender string
	Age    int
}

func sayHello(w http.ResponseWriter, r *http.Request) {
	// 2. 解析模板
	t, err := template.ParseFiles("./hello.tmpl")
	if err != nil {
		fmt.Println(t)
		fmt.Printf("Parse template failed, err:%v", err)
		return
	}
	// 3. 渲染模板
	user := UserInfo{
		Name:   "小王子",
		Gender: "男", // go语言大小写有特殊含义
		Age:    18,
	}
	// map没必要让key进行首字母大写
	m1 := map[string]interface{}{
		"name":   "小王子",
		"gender": "男",
		"age":    19,
	}
	hobbyList := []string{
		"篮球",
		"足球",
		"双色球",
	}
	//err = t.Execute(w, user)
	err = t.Execute(w, map[string]interface{}{
		"user":  user,
		"m1":    m1,
		"hobby": hobbyList,
	})
	if err != nil {
		fmt.Println("render template failed", err)
		return
	}
}

// 模板嵌套案例 开始
func f1(w http.ResponseWriter, r *http.Request) {
	// 定义一个函数
	// 返回两个值得情况下第二个必须是error，要么就只有一个返回值
	k := func(name string) (string, error) {
		return name + "年轻又帅气!", nil
	}

	// 名字一定要与模板对应上
	t := template.New("f.tmpl")
	// 一定要在解析模板文件之前告诉模板引擎现在多了一个自定义函数
	t.Funcs(template.FuncMap{
		"kua99": k,
	})

	_, err := t.ParseFiles("./f.tmpl")
	if err != nil {
		fmt.Printf("parse template failed, err:%v", err)
		return
	}
	name := "无解"
	// 渲染模板
	t.Execute(w, name)
}

func demo1(w http.ResponseWriter, r *http.Request) {

	t, err := template.ParseFiles("./t.tmpl", "./ul.tmpl")
	if err != nil {
		fmt.Printf("parse template fialed , err:%v\n", err)
		return
	}
	name := "无解"
	t.Execute(w, name)
}

// 模板嵌套案例 结束

// 模板继承案例 开始
func index(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("./base.tmpl", "./index.tmpl")
	if err != nil {
		fmt.Printf("parse template fialed , err:%v\n", err)
		return
	}
	msg := "这是index页面"
	t.ExecuteTemplate(w, "index.tmpl", msg)
}

func home(w http.ResponseWriter, r *http.Request) {
	// 先加载根模板
	t, err := template.ParseFiles("./base.tmpl", "./home.tmpl")
	if err != nil {
		fmt.Printf("parse template fialed , err:%v\n", err)
		return
	}
	msg := "这是home页面"
	t.ExecuteTemplate(w, "home.tmpl", msg)
}

// 模板继承案例 结束

func xss(w http.ResponseWriter, r *http.Request) {
	// 解析模板之前定义一个自定义的函数 safe
	t, err := template.New("xss.tmpl").Funcs(template.FuncMap{
		"safe": func(s string) template.HTML {
			return template.HTML(s)
		},
	}).ParseFiles("./xss.tmpl")
	if err != nil {
		fmt.Printf("parse template failed, err: %v\n", err)
		return
	}
	// 渲染模板
	// 以字符串的形式进行解析
	str1 := "<script>alert(111);</script>"
	str2 := "<a href='www.baidu.com'>百度</a>"
	t.Execute(w, map[string]string{
		"str1": str1,
		"str2": str2,
	})
}

func main() {
	http.HandleFunc("/", f1)
	http.HandleFunc("/tmplDemo", demo1)
	http.HandleFunc("/index", index)
	http.HandleFunc("/home", home)
	http.HandleFunc("/xss", xss)
	err := http.ListenAndServe(":3000", nil)
	if err != nil {
		fmt.Printf("HTTP server start failed, err:%v", err)
		return
	}
}
