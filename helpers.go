package main

import ("fmt";"os";og_log "log";"strconv")

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
	og_log.Printf(msg, a...)
}

func (l Log) Info(msg string, a ...any) {
	l.generic("\033[34mINFO\033[0m", msg, a...)
}

func (l Log) Fatal(msg string, a ...any) {
	l.generic("\033[1;31mFATAL\033[0m", msg, a...)
	os.Exit(1)
}

func (l Log) Err(msg string, a ...any) {
	l.generic("\033[31mERR\033[0m", msg, a...)
}

func (l Log) Req(msg string, a ...any) {
  l.generic("\033[33mreq\033[0m", msg, a...)
}

func eror(msg string, e error) {
	msg = "ERR: "+msg
	if e != nil { msg = fmt.Sprintf("%s: %v"+msg, e) }
	fmt.Fprintln(os.Stderr, msg)
}

func erorF(msg string, e error) {
	eror(msg, e) ; os.Exit(1)
}
