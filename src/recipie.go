package main

import (
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

// func init() {
// 	tpl = template.Must(template.ParseGlob("templates/*"))
// }

type recipe struct {
	Name string
	Category string
	Description string
	Steps []string
}



type reciepeDBnew struct {
	Breakfast map[string]recipe
	Lunch map[string]recipe
	Dinner map[string]recipe
	Other map[string]recipe
}

var recipesDB = reciepeDBnew{
	Breakfast: map[string]recipe{},
	Lunch:     map[string]recipe{},
	Dinner:    map[string]recipe{},
	Other:     map[string]recipe{},
}


func createRecipe(w http.ResponseWriter, req *http.Request) {
	u := getUser(req)
	if !isLoggedIn(req) {
		http.Redirect(w, req, "/", http.StatusSeeOther)
		return
	}

	if req.Method == http.MethodPost {
		ctg := req.FormValue("category")
		rpm := req.FormValue("recipeName")
		dpt := req.FormValue("description")
		itn := req.FormValue("instructions")
		var cr recipe

		if ctg == "Breakfast" {
			if _, ok := recipesDB.Breakfast[rpm]; ok {
				http.Error(w, "Recipie Already Exists", http.StatusForbidden)
				return
			}
			newItn := strings.Split(itn, "\n")
			cr = recipe{rpm, ctg, dpt, newItn}
			u.UserRecipies.Breakfast[rpm] = cr
			http.Redirect(w, req, "/reciepe/Breakfast/"+rpm, http.StatusSeeOther)
			return
		} else if ctg == "Lunch" {
			if _, ok := recipesDB.Lunch[rpm]; ok {
				http.Error(w, "Recipie Already Exists", http.StatusForbidden)
				return
			}
			newItn := strings.Split(itn, "\n")
			cr = recipe{rpm, ctg, dpt, newItn}
			u.UserRecipies.Lunch[rpm] = cr
			http.Redirect(w, req, "/reciepe/Lunch/"+rpm, http.StatusSeeOther)
			return
		} else if ctg == "Dinner" {
			if _, ok := recipesDB.Dinner[rpm]; ok {
				http.Error(w, "Recipie Already Exists", http.StatusForbidden)
				return
			}
			newItn := strings.Split(itn, "\n")
			cr = recipe{rpm, ctg, dpt, newItn}
			u.UserRecipies.Dinner[rpm] = cr
			http.Redirect(w, req, "/reciepe/Dinner/"+rpm, http.StatusSeeOther)
			return
		} else {
			var cr recipe
			if _, ok := recipesDB.Other[rpm]; ok {
				http.Error(w, "Recipie Already Exists", http.StatusForbidden)
				return
			}
			newItn := strings.Split(itn, "\n")
			cr = recipe{rpm, ctg, dpt, newItn}
			u.UserRecipies.Other[rpm] = cr
			http.Redirect(w, req, "/reciepe/Other/"+rpm, http.StatusSeeOther)
			return
		}

	}

	err := tpl.ExecuteTemplate(w, "createrecipe.gohtml", nil)
	if err != nil {
		log.Fatalln(err)
		return
	}
}

func viewReciepe(w http.ResponseWriter, req *http.Request) {
	u := getUser(req)
	if !isLoggedIn(req) {
		http.Redirect(w, req, "/", http.StatusSeeOther)
		return
	}
	vars := mux.Vars(req)
	cName := vars["category"]
	rName := vars["name"]
	var cr recipe
	if cName == "Breakfast" {
		cr = u.UserRecipies.Breakfast[rName]
	} else if cName == "Lunch" {
		cr = u.UserRecipies.Lunch[rName]
	} else if cName == "Dinner" {
		cr = u.UserRecipies.Dinner[rName]
	} else {
		cr = u.UserRecipies.Other[rName]
	}

	err := tpl.ExecuteTemplate(w, "viewReciepe.gohtml", cr)
	if err != nil {
		log.Fatalln(err)
	}

}

func deleteReciepe(w http.ResponseWriter, req *http.Request) {
	u := getUser(req)
	if !isLoggedIn(req) {
		http.Redirect(w, req, "/", http.StatusSeeOther)
		return
	}
	vars := mux.Vars(req)
	cName := vars["category"]
	rName := vars["name"]

	if cName == "Breakfast" {
		delete(u.UserRecipies.Breakfast, rName)
	} else if cName == "Lunch" {
		delete(u.UserRecipies.Lunch, rName)
	} else if cName == "Dinner" {
		delete(u.UserRecipies.Dinner, rName)
	} else {
		delete(u.UserRecipies.Other, rName)
	}

	http.Redirect(w, req, "/dashboard", http.StatusSeeOther)
	return
}
