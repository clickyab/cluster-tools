package errhandler

import (
	"context"
	"net/http"

	"github.com/mssola/user_agent"
)

const errTemplate = `
	<!doctype html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport"
          content="width=device-width, user-scalable=no, initial-scale=1.0, maximum-scale=1.0, minimum-scale=1.0">
    <meta http-equiv="X-UA-Compatible" content="ie=edge">
    <title>Document</title>
</head>
<body>
	error occurred 404
</body>
</html>
`

func errCheck(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	ua := user_agent.New(r.UserAgent())
	name, _ := ua.Browser()
	if name == "Mozilla" || name == "Opera" || name == "Edge" || name == "Chrome" || name == "Chromium" || name == "Internet Explorer" {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(errTemplate))
		return
	}
	w.WriteHeader(http.StatusNotFound)
	w.Write([]byte("error occurred 404"))

}
