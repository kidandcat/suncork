package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"strings"
)

var letterRunes = []rune("1234567890ABCDEFGHIJKLMNOPQRSTUVWXYZ")

func randStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func randInt(min, max int) int {
	return rand.Intn(max-min) + min
}

func IDGeneratorPayment() string {
	i := randStringRunes(10)
	if paymentFindByID(i).ID != "" {
		return IDGeneratorPayment()
	}
	return i
}

func IDGeneratorOrder() string {
	i := randStringRunes(10)
	if orderFindByID(i).ID != "" {
		return IDGeneratorOrder()
	}
	return i
}

func convAIToAS(input interface{}) []string {
	s := make([]string, len(input.([]interface{})))
	for i, v := range input.([]interface{}) {
		s[i] = fmt.Sprint(v)
	}
	return s
}

func getShippingOptionByID(id string) (ShippingOption, bool) {
	conf := getConfig()
	for _, v := range conf.ShippingOptions {
		if v.ID == id {
			return v, true
		}
	}
	return ShippingOption{}, false
}

func (a ShippingAddress) toString() string {
	out, e := json.Marshal(a)
	if err(e) {
		return ""
	}
	return string(out)
}

func getShippingAddressByID(s string) (ShippingAddress, bool) {
	var sa ShippingAddress
	e := json.Unmarshal([]byte(s), &sa)
	if err(e) {
		return sa, false
	}
	return sa, true
}

func optionsFromJSON(op string) []ProductOption {
	sa := []ProductOption{}
	jopts := strings.Split(op, "::")
	for _, v := range jopts {
		var o ProductOption
		e := json.Unmarshal([]byte(v), &o)
		if e == nil {
			sa = append(sa, o)
		} else {
			err(e)
		}
	}
	return sa
}

func optionsToJSON(op []ProductOption) string {
	res := ""
	for _, v := range op {
		out, _ := json.Marshal(v)
		if res == "" {
			res += string(out)
		} else {
			res += "::" + string(out)
		}
	}
	return res
}
