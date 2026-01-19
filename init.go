package main

import(
	"os"
	"fmt"
	"slices"
	"errors"
	"strings"
	"strconv"
)

func init() {
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
		 case "override":
			if e := parse_overrides(v); e != nil {
				erorF_raw(e)
			}
		 case "placeholder replacement":
			if e := populate_placeholder_replacements(v); e != nil {
				erorF_raw(e)
			}
		case "inline embed replacement":
			if e := init_populate_embeds(v); e != nil {
				erorF_raw(e)
			}
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
					erorF("used '"+og_arg+"' but received no value", nil)
					return ""

				} else { tak = append(tak, i+1) ; return args[i+1] }
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

func parse_overrides(v_R string) error {
	entries := strings.Split(v_R, ";")
	for _, e_R := range entries {
		e_R = strings.TrimSpace(e_R)
		if len(e_R) == 0 { continue } 
		fields := strings.Split(e_R, "=")
		if len(fields) < 2 { return errors.New("not a key-value pair") }
		k := strings.TrimSpace(fields[0])
		v := strings.TrimSpace(fields[1])
		page_overrides[k] = v
	}

	return nil
} 

func populate_placeholder_replacements(v_R string) error { 
	entries := strings.Split(v_R, ";")
	for _, e_R := range entries {
		e_R = strings.TrimSpace(e_R)
		if len(e_R) == 0 { continue } 
		fields := strings.Split(e_R, "=")
		if len(fields) < 2 { return errors.New("not a key-value pair") }
		k := strings.TrimSpace(fields[0])
		v := strings.TrimSpace(fields[1])
		placeholder_replacements[k] = v
	}

	return nil
}

func init_populate_embeds(v_R string) error {
	entries := strings.Split(v_R, ";")
	for _, e_R := range entries {
		e_R = strings.TrimSpace(e_R)
		if len(e_R) == 0 { continue }
		fields := strings.Split(e_R, "=")
		if len(fields) < 2 { return errors.New("Not a key-value pair") }
		k := strings.TrimSpace(fields[0])
		v := strings.TrimSpace(fields[1])
		p := struct {
			mem string
			loc rune
			esc bool
			res []string
		}{ res: []string{}, loc: '_' }
		for _, c := range v {
			if p.esc {
				switch c {
				 case 'n': p.mem += "\n"
				 case '\\': p.esc = false
				 default: p.mem += string(c)
				}
				p.esc = false ; continue
			}
			switch c {
				case '(': if p.loc == '_' {
					p.loc = 'i'
				} else { p.mem += string(c) }

			 case ')': if p.loc == 'i' {
					p.res = append(p.res, p.mem)
					p.mem = "" ; p.loc = '_'
				} else { p.mem += string(c) }

			 case '\\': p.esc = true

			 case '"':
				switch p.loc {
				 case '_':
					if p.mem != "" {
						p.res = append(p.res, p.mem)
						p.mem = ""
					}; p.loc = 'q'
				 case 'q':
					p.res = append(p.res, p.mem)
					p.mem = "" ; p.loc = '_'
				 default: p.mem += string(c)
				}

			 default:
				if p.loc == '_' { continue } else { p.mem += string(c) }
			}
		}

		for i, t := range p.res { fmt.Printf("%d{%s}\n", i, t) }

		if len(p.res) < 3 {
			if len(p.res) == 1 {
				p.res = append(p.res, p.res[0])
				p.res[0] = "" ; p.res = append(p.res, "")
			} else {
				for _, f := range p.res {
					if _, e := os.Stat(f); e != nil {
						p.res = []string{"", f, ""} ; break
					}
				} 
			}
		}; if len(p.res) < 3 { return errors.New("invalid value") }

		em := embed_stuff {
			filename: strings.TrimSpace(p.res[1]),
			pre: strings.TrimSpace(p.res[0]),
			after: strings.TrimSpace(p.res[2]),
		}
		fmt.Println(em.filename)
		embed_replacements[k] = em
	}
	return nil
}
