package main

import(
	"os"
	"time"
	"bytes"
	"errors"
	"strings"
	"net/http"
	"path/filepath"
)

var (
	embed_replacements = make(map[string]string)
	placeholder_replacements = make(map[string]string)
	page_overrides = make(map[string]string)
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
	if page_overrides[req_page] != "" { req_page = page_overrides[req_page] }
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

	{
		f := strings.Split(req_page, ".")
		if len(f) > 1 {
			f_B = replace_placeholders(f_B)
			f_B = populate_embeds(f_B)
		}
	}

	log.ReqParams(req_page, resp)
	http.ServeContent(w, r, req_page, time.Now(), bytes.NewReader(f_B))
}

func replace_placeholders(og []byte) []byte {
	res := string(og)
	
	for p, r := range placeholder_replacements {
		res = strings.ReplaceAll(res, p, r)
	}

	return []byte(res)
}

func populate_embeds(og []byte) []byte {
	res := string(og)
	for p, fi_N := range embed_replacements {
		if !strings.Contains(res, p) { continue }
		fi_N = filepath.Join(dir, fi_N)

		fi_b, e := os.ReadFile(fi_N)
		if e != nil { log.Err("%v", e) }
		if len(fi_b) < 1 { continue }

		fi_str := string(fi_b)
		if fi_str[len(fi_str)-1] == '\n' {
			fi_str = fi_str[:len(fi_str)-1]
		}
		res = strings.ReplaceAll(res, p, fi_str)
	}
	return []byte(res)
}
