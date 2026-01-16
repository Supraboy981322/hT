package main

import(
	"os"
	"time"
	"bytes"
	"errors"
	"net/http"
	"path/filepath"
)

var (
	args = os.Args[1:]
	port = ":3784"
	dir string
	log Log
)

func main() {
	http.HandleFunc("/", hanConn)
	log.Info("listening... port{%s}", port[1:])
	log.Fatal("%v", http.ListenAndServe(port, nil))
}

func hanConn(w http.ResponseWriter, r *http.Request) {
	req_page := r.URL.Path[1:]
	if req_page == "" { req_page = "index.html" }
	resp := 200
	fi_P := filepath.Join(dir, req_page)
	f_B, e := os.ReadFile(fi_P)
	if e != nil {
		if errors.Is(e, os.ErrNotExist) {
			http.Error(w, "not found", 404) ; return
		} else {
			log.HttpErr(w, "server err", req_page, e, 500) ; return
		}
	}

	log.ReqParams(req_page, resp)
	http.ServeContent(w, r, req_page, time.Now(), bytes.NewReader(f_B))
}
