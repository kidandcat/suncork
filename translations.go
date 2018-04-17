package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

func getTrans(lang string) map[string]interface{} {
	print("getTrans", lang, "opening locales/"+lang+".json")
	file, _ := os.Open("locales/" + lang + ".json")
	print("json loaded -> decodifying json")
	decoder := json.NewDecoder(file)
	var f interface{}
	decoder.Decode(&f)
	if _, ok := f.(map[string]interface{}); !ok {
		print("failed to decode json -> returning empty translations")
		return getTrans("es")
	}
	print("json decoded -> returning")
	return f.(map[string]interface{})
}

func setTrans(lang string, data map[string]string) {
	trans, e := json.Marshal(data)
	e = ioutil.WriteFile("locales/"+lang+".json", trans, 0644)
	if err(e) {
		print("Error writting translations to file", e.Error())
	}
}
