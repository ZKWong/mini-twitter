package main

import (
    "fmt"
    "log"
    "net/http"
		"io/ioutil"
)

func handler(w http.ResponseWriter, r *http.Request) {
		var sHtml string
		if r.URL.Path[1:] == "vue2516.js" {
			b, err := ioutil.ReadFile("vue2516.js")
			if err != nil{
				fmt.Print(err)
			}
			sHtml := string(b)
			fmt.Fprintf(w,"%s",sHtml)
			return
		}
		if r.URL.Path[1:] == "login.css" {
			b, err := ioutil.ReadFile("login.css")
			if err != nil{
				fmt.Print(err)
			}
			sHtml := string(b)
			w.Header().Add("Content-Type","text/css")
			fmt.Fprintf(w,"%s",sHtml)
			return
		}
		if r.URL.Path[1:] == "checklogin" {
			var username string
			var password string
			r.ParseForm()
			for key, value:= range r.Form{
				if key == "username" {
					username = fmt.Sprintf("%s",value)
				}
				if key == "password" {
					password = fmt.Sprintf("%s",value)
				}
			}
			sHtml = "login failed " + username +" "+  password
			if username=="[test]"{
				if password=="[1234]" {
					sHtml = "login successful"
				}
			}
			fmt.Fprintf(w,"%s",sHtml)
			return
		}
		sHtml = "<!DOCTYPE html>\n" +
						"<html>\n" +
						"<head>\n"+
						"<title>Mini-Twitter</title>\n"+
						"<script src='vue2516.js'></script>\n"+
						"<link rel ='stylesheet' type = 'text/css' href = 'login.css'/>\n"+
						"</head>\n"+
						"<div id='container'>\n"  +
						"<div id='main'>\n"  +
						"<h1>Mini Twitter</h1>\n"+
						" Username:&#160; <input v-model = 'username' >\n"+
						" Password:&#160; <input v-model = 'password' type ='password'>\n"+
						"<button v-on:click='login()' id = 'login'>Login</button>\n"+
						"</div>\n" +
						"</div>\n" +
						"</html>\n" +
						"<script>\n" +
						"	var main = new Vue({\n" +
						"   el: '#main',\n" + 
						"   data: {\n" +
						"			username: '',\n" +
						"     password: ''\n" +
						"   },\n" +
						"   methods: {\n" +
						"     login: function(){\n" +
						"       checklogin(this.username,this.password);" +
						"     }\n" +
						"   }\n" +
						" })\n" +
						" function checklogin(usrname,pwd) {\n" +
						"   var xhr = new XMLHttpRequest();\n" +	
						"// Register the embedded handler function\n" +	
						"	  xhr.onreadystatechange = function () {\n" +	
						"		if (xhr.readyState == 4 && xhr.status == 200) {\n" +	
						"			var loginresult = xhr.responseText;\n" +	
						"     alert(loginresult);" +	
						"		}\n" +	
						"	}\n" +	
						"	xhr.open('GET', 'checklogin?username=' + usrname + '&password=' + pwd);\n" +	
						"	xhr.send(null);\n" +	
						" }\n" + 
						"</script>" 
						
	  fmt.Fprintf(w,"%s",sHtml)
}

func main() {
    http.HandleFunc("/", handler)
    log.Fatal(http.ListenAndServe(":8080", nil))
}
