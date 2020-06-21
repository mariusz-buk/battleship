package sessions

import (
	"net/http"

	"github.com/gorilla/sessions"
)

var (
	key   = []byte("MSRvg*A_w\\h86Qs6_*wLJDD4;dBN_dXc")
	store = sessions.NewFilesystemStore("tmp", key)
)

func CheckSession(writer http.ResponseWriter, request *http.Request) *sessions.Session {
	session, _ := store.Get(request, "session")

	if session.IsNew {
		session.Values["step"] = 1
		err := session.Save(request, writer)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
		}
	}

	return session
}

func ClearSession(writer http.ResponseWriter, request *http.Request) {
	session, _ := store.Get(request, "session")

	session.Values = nil
	err := session.Save(request, writer)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
	}
}
