package tplfunc

import (
	"text/template"
	"time"
)

var Date = []Func{Now(), FormatDate()}

func Now() Func {
	return func(funcMap template.FuncMap) {
		funcMap["now"] = func() int64 {
			return time.Now().Unix()
		}
	}
}

func FormatDate() Func {
	return func(funcMap template.FuncMap) {
		funcMap["formatDate"] = func(args ...any) string {
			switch len(args) {
			case 0:
				panic("formatDate() requires at least 1 argument")
			case 1:
				switch v := args[0].(type) {
				case int, int64:
					return time.Unix(v.(int64), 0).Format("20060102150405")
				default:
					panic("formatDate() requires an int or int64 as its first argument")
				}
			case 2:
				switch v := args[0].(type) {
				case int, int64:
					return time.Unix(v.(int64), 0).Format(args[1].(string))
				default:
					panic("formatDate() requires an int or int64 as its first argument")
				}
			default:
				panic("formatDate() requires at most 2 arguments")
			}
		}
	}
}
