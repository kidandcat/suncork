package main

import (
	"net/http"
	"strconv"
)

type Adapter func(http.ResponseWriter, *http.Request, Data, *Session)

func mw(next Adapter) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		data := Data{}
		var s *Session
		if id, err := getCookie(r, CookieName); err == nil {
			idn, err := strconv.Atoi(id.Value)
			if err == nil {
				s = sess[idn]
			} else {
				s = session_start(&w)
			}
		} else {
			s = session_start(&w)
		}
		next(w, r, data, s)
	})
}
