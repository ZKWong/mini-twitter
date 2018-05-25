package main

import (
    "fmt"
    "log"
    "net/http"
		"io/ioutil"
		"bufio"
		"os"
		"strings"
		"time"
		
)

func postmessage_page(r *http.Request, sHtml *string) {
	var username string
	var msg string

	r.ParseForm()
	for key, value:= range r.Form{
		if key == "username" {
			username = fmt.Sprintf("%s",value)
			username = username[1:len(username)-1]
		}
		if key == "msg" {
			msg = fmt.Sprintf("%s",value)
			msg = msg[1:len(msg)-1]
		}
	}

	file,err := os.OpenFile(username+"_pmsg", os.O_APPEND|os.O_RDWR|os.O_CREATE,0644)
	if err != nil {fmt.Print(err)}	
	
	timevalue := fmt.Sprintf("%s",time.Now())
	timestr   := strings.Split(timevalue,".")
	msg = "[" + timestr[0] + "] " + msg
	if _, err := file.Write([]byte(msg +"\n")); err != nil{
		fmt.Print(err)
	}
	if err := file.Close(); err != nil{
		fmt.Print(err)
	}

	*sHtml = "ok"
	return

}

func getfile(sFile string, sHtml *string){
	b, err := ioutil.ReadFile(sFile)
	if err != nil{
		fmt.Print(err)
	}
	*sHtml = string(b)
	return
}

func checklogin_page(r *http.Request, sHtml *string) {
	var username string
	var password string
	var querystr []string
	var key_value []string
	var usrname string
	var pwd string

	*sHtml = "0" 

	r.ParseForm()
	for key, value:= range r.Form{
		if key == "username" {
			username = fmt.Sprintf("%s",value)
			username = username[1:len(username)-1]
		}
		if key == "password" {
			password = fmt.Sprintf("%s",value)
			password = password[1:len(password)-1]
		}
	}

	file,err := os.Open("tbusers")
	if err != nil {
		fmt.Print(err)
	}	
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		querystr  = strings.Split(scanner.Text(),"&")
		key_value = strings.Split(querystr[0],"=")
		if key_value[0]=="username" {usrname=key_value[1]}
		key_value = strings.Split(querystr[1],"=")
		if key_value[0]=="password" {pwd=key_value[1]}
		if username==usrname{
			if password==pwd {
				*sHtml = "1"
				break
			}
		}
	}

	return
}

func new_signup_page(r *http.Request, sHtml *string) {  //Check if input is empty or not, prompt user to key in value
	var new_username string
	var new_password string
	var querystr []string
	var key_value []string
	var new_usrname string
	//var new_pwd string

	*sHtml = "Sign up failed " + new_username +" "+  new_password

	new_username = ""
	new_password = ""

	r.ParseForm()
	for key, value:= range r.Form{
		if key == "new_username" {
			new_username = fmt.Sprintf("%s",value)
			new_username = new_username[1:len(new_username)-1]
		}
		if key == "new_password" {
			new_password = fmt.Sprintf("%s",value)
			new_password = new_password[1:len(new_password)-1]
		}
	}
	if new_username =="" || new_password ==""{
		*sHtml = "Username or Password cannot be empty"
		return
	}

	file,err := os.OpenFile("tbusers", os.O_APPEND|os.O_RDWR|os.O_CREATE,0644)
	if err != nil {fmt.Print(err)}	
	defer file.Close()

	scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			querystr  = strings.Split(scanner.Text(),"&")
			key_value = strings.Split(querystr[0],"=")
			if key_value[0]=="username" {new_usrname=key_value[1]}
			if new_username==new_usrname{
				*sHtml = "Username already exist"
				return
			}
		}
	
	if _, err := file.Write([]byte("username=" + new_username + "&password=" + new_password +"\n")); err != nil{
		fmt.Print(err)
	}
	if err := file.Close(); err != nil{
		fmt.Print(err)
	}
	
	*sHtml = "Sign Up successful"
	return
}

func homepage(r *http.Request,sHtml *string) {
	var uname string
	var usrname string
	var following string
	var querystr []string
	var key_value []string

	uname = ""
	following = ""
	getfile("homepage.html",sHtml)
	r.ParseForm()
	for key, value:= range r.Form{
		if key == "username" {
			uname = fmt.Sprintf("%s",value)
			uname = uname[1:len(uname)-1]
			*sHtml = strings.Replace(*sHtml,"$username",uname,2)
			break
		}
	}
	if uname != "" {
		if _,err := os.Stat(uname + "_pmsg"); err == nil{
			var content_msg string
			getfile(uname + "_pmsg",&content_msg)
			content_msg = strings.Replace(content_msg,"\n","\\n",-1)
			*sHtml = strings.Replace(*sHtml,"$content",content_msg,1)
		}
		if _,err := os.Stat(uname + "_pmsg"); os.IsNotExist(err){
			*sHtml = strings.Replace(*sHtml,"$content","",1)
		}
	}
	file,err := os.Open("tbusers")
	if err != nil {
		fmt.Print(err)
	}	
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		querystr  = strings.Split(scanner.Text(),"&")
		key_value = strings.Split(querystr[0],"=")
		if key_value[0]=="username" {
			usrname=key_value[1]
			following = following + "<tr id = 'folname' v-on:click=\"follow('" + usrname + "')\"><td>" + usrname + "</td></tr>\n"
		}
	}
	if following != "" {
		following = "<table >\n" + "<tr ><th> Following </th></tr>\n" + following + "</table>\n" 
	}
	*sHtml = strings.Replace(*sHtml,"$following",following,1)
	return
}

func get_followmsg(r *http.Request,sHtml *string){
	var uname string
	
  uname = ""
  
	r.ParseForm()
	for key, value:= range r.Form{
		if key == "username" {
			uname = fmt.Sprintf("%s",value)
			uname = uname[1:len(uname)-1]
			break
		}
	}
	if uname != "" {
		if _,err := os.Stat(uname + "_pmsg"); err == nil{
			var content_msg string
			getfile(uname + "_pmsg",&content_msg)
			*sHtml = content_msg
		}
		if _,err := os.Stat(uname + "_pmsg"); os.IsNotExist(err){
			*sHtml = ""
		}
	}
}

func web_response(w http.ResponseWriter, r *http.Request) {
	var sHtml string
		
		switch r.URL.Path[1:] {
			case "vue2516.js": 
				getfile("vue2516.js",&sHtml)
			
			case "login.css" :
				w.Header().Add("Content-Type","text/css")
				getfile("login.css",&sHtml)
			
			case "checklogin":
				checklogin_page(r, &sHtml)

			case "new_signup":
				new_signup_page(r, &sHtml)

			case "homepage.html":
				homepage(r,&sHtml) //pass by reference
	
			case "homepage.css" :
				w.Header().Add("Content-Type","text/css")
				getfile("homepage.css",&sHtml)

			case "postmessage":
				postmessage_page(r,&sHtml)

			case "get_followmsg":
				get_followmsg(r, &sHtml)

			default:
				getfile("index.html",&sHtml)
		}
						
	  fmt.Fprintf(w,"%s",sHtml)
}


func main() {
    http.HandleFunc("/", web_response)
    log.Fatal(http.ListenAndServe(":8080", nil))
}
