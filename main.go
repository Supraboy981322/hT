package main

import(
	"os"
	"errors"
	"net/http"
)

var (
	args = os.Args[1:]
	port = ":3784"
	log Log
)

func main() {
	http.HandleFunc("/", hanConn)
	log.Info("listening port{%s}", port[1:])
	log.Fatal("%v", http.ListenAndServe(port, nil))
}

func hanConn(w http.ResponseWriter, r *http.Request) {
	req_page := r.URL.Path[1:]
	resp := 200 
	f_B, e := os.ReadFile(req_page)
	if e != nil {
		if errors.Is(e, os.ErrNotExist) {
			http.Error(w, "not found", 404) ; resp = 404
		} else {
			log.Err("page{%s} : %v", req_page, e)
			http.Error(w, "server err", 500)  ; resp = 500
		}
	}

	log.Req("page{%q} ; resp{%d}", req_page, resp)
	w.Write(f_B)
}
