package main 



import (
	"net/http"
)

func getUser(req *http.Request) user {

	var u user
	c, err := req.Cookie("session")
	if err != nil {
		return u
	}

	if un, ok := sessionsDatabase[c.Value]; ok {
		u = userDatabase[un]
	}
	return u
}


func isLoggedIn(req *http.Request) bool {
	c, err := req.Cookie("session")
	if err != nil {
		return false
	}

	uName := sessionsDatabase[c.Value]
	_, ok := userDatabase[uName]

	return ok
}



