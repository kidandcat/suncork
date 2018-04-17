package main

import (
	"net/http"
	"strconv"
	"time"
)

const CookieName = "__SUNCORK__ID__"

var sess map[int]*Session

type Option struct {
	Name  string
	Value string
}
type ProductWithOption struct {
	Product Product
	Options []Option
}

type Session struct {
	lang  string
	cart  []ProductWithOption
	admin bool
}

func initSessionDB() {
	sess = map[int]*Session{}
}

func mwSession(next Adapter) Adapter {
	return Adapter(func(w http.ResponseWriter, r *http.Request, data Data, s *Session) {
		print("Session middleware start")
		if id, err := getCookie(r, CookieName); err == nil {
			idn, err := strconv.Atoi(id.Value)
			_, ok := sess[idn]
			if err == nil && ok {
				s = sess[idn]
			} else {
				s = session_start(&w)
			}
		} else {
			s = session_start(&w)
		}
		print("Session middleware end")
		next(w, r, data, s)
	})
}

func session_start(w *http.ResponseWriter) *Session {
	print("session_start start")
	index := randInt(1, 9223372036854775807)
	_, ok := sess[index]
	for ok {
		index := randInt(1, 9223372036854775807)
		_, ok = sess[index]
	}
	setCookie(w, http.Cookie{
		Name:     CookieName,
		Value:    strconv.Itoa(index),
		Path:     "/",
		HttpOnly: true,
		Expires:  time.Now().AddDate(1, 0, 0),
	})
	sess[index] = &Session{
		cart: []ProductWithOption{},
		lang: "es",
	}
	print("session_start end")
	return sess[index]
}

func setCookie(w *http.ResponseWriter, cookie http.Cookie) {
	http.SetCookie(*w, &cookie)
}

func getCookie(r *http.Request, key string) (*http.Cookie, error) {
	cookie, e := r.Cookie(key)
	if err(e) {
		return nil, e
	}
	return cookie, nil
}
