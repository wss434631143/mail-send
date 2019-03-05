package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"gopkg.in/gomail.v2"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
)

type AlarmEmailBody struct {
	From     string `json:"from"`
	PassWord string `json:"password"`
	To       string `json:"tos"`
	Subject  string `json:"subject"`
	Body     string `json:"content"`
	Host     string `json:"host"`
	Port     int    `json:"port"`
}

type Error struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

var (
	EmailUser     = os.Getenv("EMAILUSER")
	EmailPassword = os.Getenv("EMAILPASSWORD")
	EmailHost     = os.Getenv("EMAILHOST")
	EmailPort     = os.Getenv("EMAILPORT")
)

func _build(code int, defval string, custom ...string) Error {
	msg := defval
	if len(custom) > 0 {
		msg = custom[0]
	}
	return Error{
		Code: code,
		Msg:  msg,
	}
}

// 400
func BadRequestError(msg ...string) Error {
	return _build(http.StatusBadRequest, "bad request", msg...)
}

func String(r *http.Request, key string, defVal string) string {
	if val, ok := r.URL.Query()[key]; ok {
		if val[0] == "" {
			return defVal
		}
		return strings.TrimSpace(val[0])
	}

	if r.Form == nil {
		err := r.ParseForm()
		if err != nil {
			panic(BadRequestError())
		}
	}

	val := r.Form.Get(key)
	if val == "" {
		return defVal
	}

	return strings.TrimSpace(val)
}

func MustString(r *http.Request, key string, displayName ...string) string {
	val := String(r, key, "")
	if val == "" {
		name := key
		if len(displayName) > 0 {
			name = displayName[0]
		}
		panic(BadRequestError(fmt.Sprintf("%s is necessary", name)))
	}
	return val
}

//发送告警邮件
func AlarmEMailSend(context AlarmEmailBody) {
	m := gomail.NewMessage()
	m.SetHeader("From", context.From)
	m.SetHeader("To", strings.Split(context.To, ",")...)
	m.SetHeader("Subject", context.Subject)
	m.SetBody("text/plain", context.Body)
	d := gomail.NewDialer(context.Host, context.Port, context.From, context.PassWord)
	if err := d.DialAndSend(m); err != nil {
		log.Print(err)
	}
}

//接受发送告警邮件请求
func AlarmEMailHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		tos := MustString(r, "tos")
		subject := MustString(r, "subject")
		content := MustString(r, "content")
                log.Printf("tos: %s", tos)
		body, _ := ioutil.ReadAll(r.Body)
		log.Printf("accept from %s", r.RemoteAddr)
		port, _ := strconv.Atoi(EmailPort)
		context := AlarmEmailBody{To: tos, Subject: subject, Body: content, From: EmailUser, Host: EmailHost, Port: port, PassWord: EmailPassword}
		json.Unmarshal(body, &context)
		go AlarmEMailSend(context)
		w.WriteHeader(http.StatusAccepted)
	}
}

func main() {
	EnvArray := []string{EmailPassword, EmailHost, EmailUser, EmailPort}
	for _, v := range EnvArray {
		if len(v) == 0 {
			log.Fatal("环境变量配置缺失!")
		}
	}
	host := flag.String("host", ":8011", "the host port listen")
	flag.Parse()
	log.Printf("开始监听，%s\n", *host)
	http.HandleFunc("/api/v1/msg/alarm/email/", AlarmEMailHandler)
	http.ListenAndServe(*host, nil)
}
