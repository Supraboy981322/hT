package main

import(
	"os"
	"fmt"
	"mime"
	"slices"
	"errors"
	"strings"
	"strconv"
)

func init() {
	mime.AddExtensionType(".js", "application/javascript")
	mime.AddExtensionType(".mjs", "application/javascript")
	read_config()
	parse_args()
}

func read_config() {
	if slices.Contains(args, "--no-config") { return }
	fi_B, e := os.ReadFile("hT.conf")
	if errors.Is(e, os.ErrNotExist) { return }
	log.Info("using config") ; if e != nil {
		if errors.Is(e, os.ErrNotExist) { return }
		erorF_raw(e)
	}

	fi_str := string(fi_B)
	lines := strings.Split(fi_str, "\n")

	for i, l := range lines {
		li_N := "'"+strconv.Itoa(i+1)+"'"
		l = strings.TrimSpace(l)

		//skip empty lines and comments
		if len(l) == 0 { continue }
		if l[0] == '#' { continue }

		fields := strings.Split(l, ":") 
		if len(fields) > 2 || len(fields) < 2 {
			erorF("not a key-value pair (line "+li_N+")", nil)
		}

		k := strings.TrimSpace(fields[0])
		if len(k) == 0 {
			erorF("missing key (line "+li_N+")", nil)
		}
		v := strings.TrimSpace(fields[1])
		if len(v) == 0 {
			erorF("missing value for key '"+k+"' (line "+li_N+")", nil)
		}

		esc_v := fmt.Sprintf("%q", v)
		esc_k := fmt.Sprintf("%q", k)
		switch k {
		 case "dir":
			if s, e := os.Stat(v); e != nil {
				if errors.Is(e, os.ErrNotExist) {
					erorF("configured dir (line "+li_N+" of 'hT.conf') doesn't exist", nil)
				} else { erorF_raw(e) } 
			} else if !s.IsDir() {
				erorF("configured dir (line "+li_N+" of hT.conf) isn't a dir", nil)
			} else { dir = v }
		 case "port":
			if _, e := strconv.Atoi(v); e != nil {
				erorF("invalid value for "+esc_k+"; "+esc_v+" is not a number", nil)
			} else { port = ":"+v }
		 default:
			erorF("invalid config key: "+esc_k+" (line "+li_N+")", nil)
		}
	}
}

func parse_args() {
	if len(args) > 0 {
		var tak []int
		for i, arg := range args {
			if len(arg) < 1 || slices.Contains(tak, i) { continue }
			og_arg := arg ; arg = arg[1:]
			next := func() string {
				if i+1 >= len(args) {
					erorF("used '"+og_arg+"' but recieved no value", nil)
					return ""
				} else { tak = append(tak, i+1) ; return args[i+1] }
				panic(nil)
			}
			if arg[0] == '-' && len(arg) > 1 {
				switch arg[1:] {
				 case "no-config": continue //already handled when parsed config
				 case "dir", "directory": next()
				 case "port": port = ":"+next()
				 default: erorF("invalid arg ("+og_arg+")", nil)
				}
			} else {
				for _, a := range arg {
					switch a {
					 case 'd': dir = next()
					 case 'p': port = ":"+next()
					 default: erorF("invalid arg ("+og_arg+")", nil)
					}
				}
			}
		}
	}
}
