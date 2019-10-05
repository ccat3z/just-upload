package main

import (
	"fmt"
	"html/template"
	"io"
	"net/http"
	"os"
)

const page = `
<!DOCTYPE html>
<html>
  <head>
    <title>Upload</title>
  </head>
  <body>
    <form action="/" method="post" enctype="multipart/form-data">
     <div>
       <input type="file" id="file" name="file">
	 </div>
	 <button>Submit</button>
    </form>
    <span style="color: red;">{{.}}</span>
  </body>
</html>
`

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", upload)
	http.ListenAndServe(":7088", mux)
}

func upload(w http.ResponseWriter, r *http.Request) {
	var msg string

	defer func() {
		t, err := template.New("index").Parse(page)
		if err != nil {
			panic(err)
		}

		err = t.Execute(w, msg)
		if err != nil {
			panic(err)
		}
	}()

	if r.Method == "POST" {
		r.ParseMultipartForm(32 << 10)

		file, handler, err := r.FormFile("file")
		if err != nil {
			msg = err.Error()
			return
		}
		defer file.Close()

		if _, err := os.Stat(handler.Filename); !os.IsNotExist(err) {
			msg = fmt.Sprintf("file %v exist", handler.Filename)
			return
		}

		f, err := os.OpenFile(handler.Filename, os.O_WRONLY|os.O_CREATE, 0666)
		if err != nil {
			msg = err.Error()
			return
		}
		defer f.Close()

		io.Copy(f, file)
		msg = fmt.Sprintf("get %v", handler.Filename)
	}
}
