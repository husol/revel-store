package models

import (
	b64 "encoding/base64"
	"path"
	"github.com/revel/revel"
	"bytes"
	"html/template"
	"path/filepath"
	"log"
)

type HusAjax struct {
	Messages string
	Html interface{}
}

func (husAjax *HusAjax) Fetch(file_path string, data interface{}) string {
	fileTmpl := path.Join(revel.ViewsPath, file_path)
	tmpl, err := template.New(filepath.Base(file_path)).Funcs(revel.TemplateFuncs).ParseFiles(fileTmpl)
	if err != nil {
		log.Fatal(err)
	}

	var w bytes.Buffer
	tmpl.Execute(&w, data)

	return w.String()
}

func (husAjax *HusAjax) SetHTML(id, content string) {
	html := make(map[string] string);
	html["id"] = id
	html["content"] = b64.StdEncoding.EncodeToString([]byte(content))

	husAjax.Html = html
}

func (husAjax *HusAjax) SetMessage(message string) {
	messages := husAjax.Messages + "<span>"+message+"</span><br/>"
	husAjax.Messages = messages
}

func (husAjax *HusAjax) OutData(result interface{}) interface{} {
	returnData := make(map[string]interface{})
	returnData["messages"] = husAjax.Messages
	returnData["html"] = husAjax.Html
	returnData["result"] = result
	return returnData
}
