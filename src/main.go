package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"golang.org/x/crypto/bcrypt"
	"github.com/gorilla/mux"

	uuid "github.com/satori/go.uuid"
)


var tpl *template.Template


type user struct {
	Unemail string
	FirstName string
	SecondName string
	Password []byte
	UserRecipies reciepeDBnew
}

var userDatabase = map[string]user{}
var sessionsDatabase = map[string]string{}


func init() {
	tpl = template.Must(template.ParseGlob("templates/*"))
}



func main() {
	r := mux.NewRouter()
	r.HandleFunc("/", indexFunc)
	r.HandleFunc("/dashboard", dashboardFunc)
	r.HandleFunc("/signup", signUpFunc)
	r.HandleFunc("/login", loginFunc)
	r.HandleFunc("/logout", logoutFunc)
	r.HandleFunc("/reciepe/creation", createRecipe)
	r.HandleFunc("/reciepe/{category}/{name}", viewReciepe)
	r.HandleFunc("/reciepe/delete/{category}/{name}", deleteReciepe)
	log.Fatal(http.ListenAndServe(":8080", r))
}

func indexFunc(w http.ResponseWriter, req *http.Request) {
	err := tpl.ExecuteTemplate(w, "homePage.gohtml", nil)
	if err != nil {
		log.Fatalln(err)
	}
}

func dashboardFunc(w http.ResponseWriter, req *http.Request) {
	u := getUser(req)
	fmt.Println(u, sessionsDatabase, isLoggedIn(req))
	if !isLoggedIn(req) {
		http.Redirect(w, req, "/", http.StatusSeeOther)
		return
	}

	err := tpl.ExecuteTemplate(w, "dashboard.gohtml", u)
	if err != nil {
		log.Fatalln(err)
	}
}

func signUpFunc(w http.ResponseWriter, req *http.Request) {
	// check if the user is already logged in
	if isLoggedIn(req) {
		http.Redirect(w, req, "/dashboard", http.StatusSeeOther)
		return
	}

	if req.Method == http.MethodPost {

		uName := req.FormValue("useremail")
		fName := req.FormValue("firstName")
		sName := req.FormValue("secondName")
		password := req.FormValue("password")


		sID, _ := uuid.NewV4()

		// Create a cookie that will contain the session ID
		c := &http.Cookie{
			Name: "session",
			Value: sID.String(),
		}

		if _, ok := userDatabase[uName]; ok{
			http.Error(w, "Username already taken", http.StatusForbidden)
			return
		}

		// set that cookie to the browser
		http.SetCookie(w, c)
		fmt.Println(c.Value, sessionsDatabase, uName)
		// save session
		sessionsDatabase[c.Value] = uName
		bs, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)


		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		} 

		// save the user info
		u := user{uName, fName, sName, bs, recipesDB}
		userDatabase[uName] = u

		// redirects to the dashboard page
		http.Redirect(w, req, "/dashboard", http.StatusSeeOther)
		return
	}
	
	err := tpl.ExecuteTemplate(w, "signUp.gohtml", nil)
	if err != nil {
		log.Fatalln(err)
	}
}

func loginFunc(w http.ResponseWriter, req *http.Request) {
	if isLoggedIn(req) {
		http.Redirect(w, req, "/dashboard", http.StatusSeeOther)
		return
	}

	if req.Method == http.MethodPost {
		un := req.FormValue("email")
		p := req.FormValue("password")

		if _, ok := userDatabase[un]; !ok{
			http.Error(w, "Username and/or password do not match", http.StatusForbidden)
			return
		}

		u := userDatabase[un]
		err := bcrypt.CompareHashAndPassword(u.Password, []byte(p))
		if err != nil {
			http.Error(w, "Username and/or password do not match", http.StatusForbidden)
			return
		}

		sID, _ := uuid.NewV4()
		c := &http.Cookie{
			Name: "session",
			Value: sID.String(),
		}

		http.SetCookie(w, c)
		sessionsDatabase[c.Value] = un
		http.Redirect(w, req, "/dashboard", http.StatusSeeOther)
		return
	}
	err := tpl.ExecuteTemplate(w, "login.gohtml", nil)
	if err != nil {
		log.Fatalln(err)
	}
}

func logoutFunc(w http.ResponseWriter, req *http.Request) {
	if !isLoggedIn(req) {
		http.Redirect(w, req, "/", http.StatusSeeOther)
		return
	}
	c, _ := req.Cookie("session")
	// delete the session
	delete(sessionsDatabase, c.Value)
	// remove the cookie
	c = &http.Cookie{
		Name:   "session",
		Value:  "",
		MaxAge: -1,
	}
	http.SetCookie(w, c)
	http.Redirect(w, req, "/login", http.StatusSeeOther)
}


