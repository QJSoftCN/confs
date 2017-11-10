package confs

import (
	"io/ioutil"
	"log"
	"strings"
	"encoding/json"
	"path"
	"path/filepath"
)

var all map[string](map[string]interface{})

const (
	default_conf_root = "conf"
	json_file         = "json"

	conf_locale = "locale"
	top_conf    = "top"
)

var confRoot string

func SetConfDir(confDir string) {
	confRoot = filepath.Clean(confDir)
	//re load confs
	loadConfs()
}

func loadConfs() {
	if len(confRoot) == 0 {
		confRoot = default_conf_root
	}
	all = make(map[string](map[string]interface{}))

	files, err := ioutil.ReadDir(confRoot)
	if err != nil {
		log.Fatal(err)
		return
	}

	for index, file := range files {
		log.Println("read No.", index+1, " confige file name:", file.Name())
		switch strings.ToLower(path.Ext(file.Name())) {
		case json_file:
			readConfByJson(file.Name())
		default:
			readConfByJson(file.Name())
		}
	}
}

func init() {
	loadConfs()
}

func readConfByJson(fileName string) {
	bs, err := ioutil.ReadFile(confRoot + "/" + fileName)
	if err != nil {
		log.Println(err)
	}

	var fmap map[string]interface{}
	err = json.Unmarshal(bs, &fmap)
	if err != nil {
		log.Println(err)
	}

	lastDotIndex := strings.LastIndex(fileName, ".")
	group := fileName[0:lastDotIndex]

	all[group] = fmap
	log.Println("put ", group, " confs")
}

func Locale() string {
	return GetString(top_conf, conf_locale)
}

func get(group string, keys ...string) (interface{}, bool) {
	var t map[string]interface{}
	t = all[group]
	for _, key := range keys {
		v, ok := (t[key]).(map[string]interface{})
		if ok {
			t = v
		} else {
			val, ok := t[key]
			return val, ok
		}
	}
	return t, true
}

func GetString(group string, keys ...string) string {
	val, ok := get(group, keys...)
	if ok {
		s, _ := val.(string)
		return s
	} else {
		return ""
	}
}

func GetInt(group string, keys ...string) int {
	val, ok := get(group, keys...)
	if ok {
		s, _ := val.(float64)
		return int(s)
	} else {
		return 0
	}
}

func GetFloat(group string, keys ...string) (float64, bool) {
	val, ok := get(group, keys...)
	if ok {
		s, suc := val.(float64)
		return s, suc
	} else {
		return 0, ok
	}
}
