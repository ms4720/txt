package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"reflect"
	"strconv"
	"strings"
)

var funcs = map[string]interface{}{
	"slice": func(what interface{}, slice ...interface{}) (interface{}, error) {
		v := reflect.ValueOf(what)
		switch v.Kind() {
		case reflect.Slice, reflect.Array, reflect.String:
		default:
			return nil, fmt.Errorf("can't index item of type %s", v.Type())
		}
		if ln := len(slice); ln == 0 || ln > 2 {
			return nil, fmt.Errorf("slice takes 2 or 3 parameters, given %d", ln)
		}

		ln := v.Len()
		start, err := index(slice[0])
		if err != nil {
			return nil, err
		}

		usestop := false
		stop := ln
		if len(slice) == 2 {
			usestop = true
			stop, err = index(slice[1])
			if err != nil {
				return nil, err
			}
		}

		if start == stop {
			return v.Slice(0, 0).Interface(), nil
		}

		if usestop && start < 0 && stop < 0 {
			if start > stop {
				return nil, fmt.Errorf("invalid indicies (%d, %d)", start, stop)
			}
			//flip so correct when we un-negative them
			start, stop = stop, start
		}

		if start < 0 {
			start = ln + start
		}
		if usestop && stop < 0 {
			stop = ln + stop
		}

		if start > stop {
			if usestop {
				return nil, fmt.Errorf("invalid indicies (%d, %d)", start, stop)
			} else {
				return v.Slice(0, 0).Interface(), nil
			}
		}

		//clamp
		if usestop && stop > ln {
			stop = ln
		}

		return v.Slice(start, stop).Interface(), nil
	},
	"nl": func(s string) string {
		if len(s) == 0 || s[len(s)-1] == '\n' {
			return s
		}
		return s + "\n"
	},
	"parseCSV": func(header, input string) (interface{}, error) {
		hdr := splitHeader(header)
		return CSV(hdr, rdr(input))
	},
	"parseJSON": func(input string) (interface{}, error) {
		return JSON(rdr(input))
	},
	"parseLine": func(RS, LP, header, input string) (interface{}, error) {
		if RS == "" {
			RS = *RecordSeparator
		}
		if LP == "" {
			LP = *LinePattern
		}
		hdr := splitHeader(header)
		return SubmatchSplit(hdr, RS, LP, rdr(input))
	},
	"parse": func(RS, FS, header, input string) (interface{}, error) {
		if RS == "" {
			RS = *RecordSeparator
		}
		if FS == "" {
			FS = *FieldSeparator
		}
		hdr := splitHeader(header)
		return Split(hdr, RS, FS, rdr(input))
	},
	"quoteCSV": func(s string) string {
		hasQuote := strings.Index(s, `"`) > 0
		hasComma := strings.Index(s, ",") > 0
		if hasComma && !hasQuote {
			return `"` + s + `"`
		}
		if hasQuote {
			return `"` + strings.Replace(s, `"`, `""`, -1) + `"`
		}
		return s
	},

	"toJSON": func(v interface{}) (string, error) {
		bs, err := json.Marshal(v)
		return string(bs), err
	},

	"read": func(f string) (string, error) {
		bs, err := ioutil.ReadFile(f)
		return string(bs), err
	},

	"equalFold": strings.EqualFold,
	"fields":    strings.Fields,
	"join": func(sep string, a []string) string {
		return strings.Join(a, sep)
	},
	"lower":      strings.ToLower,
	"upper":      strings.ToUpper,
	"title":      strings.ToTitle,
	"trimCutset": swapArgs(strings.Trim),
	"trimLeft":   swapArgs(strings.TrimLeft),
	"trimRight":  swapArgs(strings.TrimRight),
	"trimPrefix": swapArgs(strings.TrimPrefix),
	"trimSuffix": swapArgs(strings.TrimSuffix),
	"trim":       strings.TrimSpace,

	"quoteGo":      strconv.Quote,
	"quoteGoASCII": strconv.QuoteToASCII,

	"match": func(pattern, src string) (bool, error) {
		r, err := cmpl(pattern)
		if err != nil {
			return false, err
		}
		return r.MatchString(src), nil
	},
	"find": func(pattern, src string) ([]string, error) {
		r, err := cmpl(pattern)
		if err != nil {
			return nil, err
		}
		return r.FindAllString(src, -1), nil
	},
	"replace": func(pattern, template, src string) (string, error) {
		r, err := cmpl(pattern)
		if err != nil {
			return "", err
		}
		return r.ReplaceAllString(src, template), nil
	},
	"split": func(pattern, src string) ([]string, error) {
		r, err := cmpl(pattern)
		if err != nil {
			return nil, err
		}
		return r.Split(src, -1), nil
	},

	"env": os.Getenv,

	"exec": func(name string, args ...string) string {
		return run(exec.Command(name, args...))
	},
	"pipe": func(name string, args ...string) (string, error) {
		if len(args) == 0 {
			return "", errors.New("pipe requires an input as the last argument")
		}
		last := len(args) - 1
		input := args[last]
		args = args[:last]
		cmd := exec.Command(name, args...)
		cmd.Stdin = strings.NewReader(input)
		return run(cmd), nil
	},
}
