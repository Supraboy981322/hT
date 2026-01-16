package main

import("slices")

func init() {
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
				 case "port": port = ":"+next()
				 default: erorF("invalid arg ("+og_arg+")", nil)
				}
			} else {
				for _, a := range arg {
					switch a {
					 case 'p': port = ":"+next()
					 default: erorF("invalid arg ("+og_arg+")", nil)
					}
				}
			}
		}
	}
}
