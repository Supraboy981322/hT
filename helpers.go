package main

import (
	"os"
	"fmt"
	"time"
	"strconv"
	"net/http"
)

type Log struct {}

func strip_esc(msg string) string {
	var res string
	var esc bool
	for _, c := range msg {
		if esc {
			if c == ';' || c == '[' { continue }
			_, e := strconv.Atoi(string(c))
			if e != nil { esc = false }
		} else if c == '\033' {
			esc = true
		} else { res += string(c) }
	}
	return res
}

func (l Log) generic(pre string, msg string, a ...any) {
	msg = "\033[0m\033[37m[\033[0m"+pre+"\033[0m\033[37m]:\033[0m "+msg
	t := time.Now().Format("2006/01/02 15:04:05")
	fmt.Printf(msg+" \033[0;1mtime\033[0;38;2;130;139;184m{\033[0m"+t+"\033[0;38;2;130;139;184m}\033[0m"+"\n", a...)
}

func (l Log) Info(msg string, a ...any) {
	pre := "\033[38;2;130;170;255mINFO\033[0m"
	l.generic(pre, msg, a...)
}

func (l Log) Fatal(msg string, a ...any) {
	pre := "\033[1;38;2;255;117;127mFATAL\033[0m"
	l.generic(pre, msg, a...)
	os.Exit(1)
}

func (l Log) Err(msg string, a ...any) {
	pre := "\033[38;2;255;117;127mERR\033[0m"
	l.generic(pre, msg, a...)
}

func (l Log) Req(msg string, a ...any) {
	pre := "\033[38;2;255;199;119mreq\033[0m"
  l.generic(pre, msg, a...)
}

func (l Log) ReqParams(page string, resp int) {
	format :="page\033[0;38;2;130;139;184m{\033[0m"+
			"%q\033[0;38;2;130;139;184m} ; \033[0m"+
			"resp\033[0;38;2;130;139;184m{\033[0m%v}\033[0m"
	if len(page) > 0 {
		if page[0] != '/' { page = "/"+page }
	}
	l.Req(format, page, resp)
}

func (l Log) HttpErr(
	w http.ResponseWriter,
	msg string,
	page string,
	e error,
	resp int,
) {
	l.Err("page\033[0;38;2;130;139;184m{\033[0m"+
			"%q\033[0;38;2;130;139;184m} ; \033[0m"+
			"\033[38;2;255;117;127merr"+
			"\033[0;38;2;130;139;184m{"+
			"\033[0m%v}\033[0m", page, e)
	http.Error(w, "server err", resp)
}

func eror(msg string, e error) {
	msg = "\033[1;31mERR:\033[0m "+msg
	if e != nil { msg = fmt.Sprintf("%s: %v"+msg, e) }
	fmt.Fprintln(os.Stderr, msg)
}

func erorF(msg string, e error) {
	eror(msg, e) ; os.Exit(1)
}
