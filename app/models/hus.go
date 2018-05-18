package models

import (
	"github.com/revel/revel"
	"encoding/json"
	"math/rand"
	"mime/multipart"
	"os"
	"io/ioutil"
)

type Hus struct {
	*revel.Controller
}

func (hus *Hus) Log(data interface{}, overwrite bool) error {
	jsonData, err := json.MarshalIndent(data, "", "\t")
	if overwrite {
		err = ioutil.WriteFile("/tmp/debug", jsonData, 0664)
	} else {
		f, err := os.OpenFile("/tmp/debug", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0664)
		if err != nil {
			return err
		}
		defer f.Close()
		_, err = f.Write(jsonData)
		_, err = f.WriteString("\n\n")
	}

	if err != nil {
		return err
	}
	return nil
}

func (hus *Hus) BaseUrl() string {
	baseUrl, _ := revel.Config.String("app.baseurl")

	return baseUrl
}

func (hus *Hus) RandomString(n int) string {
	var letterRunes = []rune("abcdefghijklmnopqrstuvwxyz0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ")

	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func (hus *Hus) DecodeObjSession(encode_str string, obj interface{}) {
	json.Unmarshal([]byte(encode_str), obj)
}

func (hus *Hus) DeleteDirFile(path string) bool {
	if path == "" {
		return false
	}

	path = revel.BasePath + path
	if err := os.RemoveAll(path); err != nil {
		return false
	}
	return true
}

func (hus *Hus) ValidateFiles(files []*multipart.FileHeader, types []string) bool {
	for _, file := range files {
		check := false
		for _, fileType := range types {
			if file.Header.Get("Content-Type") == fileType {
				check = true
			}
		}
		if !check {
			return false
		}
	}

	return true;
}
