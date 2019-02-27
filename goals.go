package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

// Uses local MYSQL database.
var db *sql.DB

// Goal struct used to import/export data to MySQL database has 4 columns.
// ID is primary key & auto increments in MySQL database
type Goal struct {
	ID         int
	Goal       string
	Typeofgoal string
	Notes      string
}

//function that opens MySQL DB connection on local server ===================================
func Dbconn() (db *sql.DB) {
	db, err := sql.Open("mysql", "root:dennisjohn@/goals")
	if err != nil {
		fmt.Println("problem getting goals index from DB: ", err)
	}
	// checkError(err)
	err = db.Ping()
	if err != nil {
		fmt.Println("problem pinging DB: ", err)
	}
	return db
}

//Index function that lists all goals on main index page=============================================
func Index(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	db := Dbconn()
	goals, err := db.Query("SELECT * FROM goals.main") //Query into DB goals, table name "main"
	if err != nil {
		fmt.Println("problem with index function & retrieving goals index from DB: ", err)
	}

	var AllGoals []Goal
	var tempvariable Goal

	//scans data from "goals" Query above and saves into a temporary structs called tempvariable
	//tempvariable is then appended into AllGoals, the main struct that contains all goals
	for goals.Next() {
		err = goals.Scan(&tempvariable.ID, &tempvariable.Goal, &tempvariable.Typeofgoal, &tempvariable.Notes)
		AllGoals = append(AllGoals, tempvariable)
	}

	//I used this to test to print to terminal & to make sure DB was updating.
	fmt.Println(AllGoals)

	//parse & return index.html file into variable t
	t, err := template.ParseFiles("index.html")
	if err != nil {
		fmt.Println("template parsing error: ", err)
	}

	//t is executed with "AllGoals" being sent to the index.html page
	err = t.Execute(w, AllGoals)
	if err != nil {
		fmt.Println("template executing error: ", err)
	}
	defer db.Close() //not sure why I needed this, but everyone seemed to do it
}

//Displays new form to create a new goal  ====================================
func New(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	t, err := template.ParseFiles("new.html")
	if err != nil {
		fmt.Println("template parsing error: ", err)
	}
	err = t.Execute(w, nil)
	if err != nil {
		fmt.Println("template executing error: ", err)
	}

}

//show route - When button is clicked, shows more details about each goal ===========================
func Show(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	db := Dbconn()
	vars := mux.Vars(r) //gets parameters from URL and saves into vars variable
	ID := vars["id"]    //ID variable identifies goals that needs to be show. Sent later to Query
	fmt.Fprintf(w, "Goals ID is : %v\n", ID)

	goals, err := db.Query("SELECT * FROM goals.main WHERE uid=?", ID) //Query into MySQL Database

	var tempvariable Goal
	//scans data from "goals" Query above and saves into a temporary structs called tempvariable
	//tempvariable will later be sent to the show form to
	for goals.Next() {
		err = goals.Scan(&tempvariable.ID, &tempvariable.Goal, &tempvariable.Typeofgoal, &tempvariable.Notes)
	}
	fmt.Println(tempvariable) //test to print to terminal to see if the above code works

	//parse show.html form and save into variable t
	t, err := template.ParseFiles("show.html")
	if err != nil {
		log.Print("template parsing error: ", err)
	}

	//tempvariable containing data is sent to t, which is parsed "show.html" form
	err = t.Execute(w, tempvariable)
	if err != nil {
		fmt.Println("template executing error: ", err)
	}
	defer db.Close()
}

//Create route - creates a new goal into database. Received data from "new" as POST ==========================
func Create(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	db := Dbconn()
	//form data from "new.html" is parsed and saved into Newgoal, Newtypeofgoal, Newnotes
	Newgoal := r.FormValue("Goal")
	Newtypeofgoal := r.FormValue("Typeofgoal")
	Newnotes := r.FormValue("Notes")
	fmt.Println(Newgoal, Newtypeofgoal, Newnotes) //test to see if above works
	//inserts New goal fields into MySQL database using Insert function
	result, err := db.Exec("INSERT INTO goals.main (goal,typeofgoal,notes) VALUES(?, ?, ?)", Newgoal, Newtypeofgoal, Newnotes)
	if err != nil {
		fmt.Println("Insertion problem:  ", err)
	}

	fmt.Println(result)                           //test to see if above Query works
	fmt.Println(Newgoal, Newtypeofgoal, Newnotes) //test to see if form data from "new" parses.
	defer db.Close()
	http.Redirect(w, r, "/goals", http.StatusSeeOther) //redirect to Index
}

//Shows Edit form with populated fields ===================================
func Edit(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	fmt.Fprintf(w, "This is the edit page \n")
	db := Dbconn()
	vars := mux.Vars(r) //gets ID and saves into vars variable. Will be used for Query below
	ID := vars["id"]    //ID variable is from parameter

	goals, err := db.Query("SELECT * FROM goals.main WHERE uid=?", ID)
	if err != nil {
		fmt.Println("Problem with Show Query: ", err)
	}
	fmt.Println(goals) //test to see if above Query works

	var tempvariable Goal
	//gets data and saves into a temporary struct so it can be sent to edit.html
	for goals.Next() {
		goals.Scan(&tempvariable.ID, &tempvariable.Goal, &tempvariable.Typeofgoal, &tempvariable.Notes)
	}
	fmt.Println(tempvariable) //test to see is tempvariable stores correct data

	//parse & return edit.html file
	t, err := template.ParseFiles("edit.html")
	if err != nil {
		fmt.Println("template parsing error: ", err)
	}

	//temp variable is sent to edit.html form
	err = t.Execute(w, tempvariable)
	if err != nil {
		fmt.Println("template executing error: ", err)
	}
	defer db.Close()

}

//Update route - receives data from edit form
func Update(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	db := Dbconn()
	vars := mux.Vars(r) //vars will store parameters, specifically ID
	ID := vars["id"]
	//Updated form data is parsed & stored into "updated" variables
	Updatedgoal := r.FormValue("Goal")
	Updatedtypeofgoal := r.FormValue("Typeofgoal")
	Updatednotes := r.FormValue("Notes")
	fmt.Println(Updatedgoal, Updatedtypeofgoal, Updatednotes) //test to see above data works
	//Update query
	updatedrow, err := db.Exec("UPDATE goals.main SET goal=?, typeofgoal=?, notes=? WHERE uid=?", Updatedgoal, Updatedtypeofgoal, Updatednotes, ID)
	if err != nil {
		fmt.Println("problem updating goals index from DB: ", err)
	}
	fmt.Println(updatedrow)
	defer db.Close()
	http.Redirect(w, r, "/goals", http.StatusSeeOther) //redirect to Index

	// fmt.Fprintf(w, "Update Goals ID is : %v\n", ID)
}

//delete route that deleted a goal =================================
func Delete(w http.ResponseWriter, r *http.Request) {
	db := Dbconn()
	vars := mux.Vars(r) //vars will store parameters from URL
	ID := vars["id"]

	//Delete Query for database
	_, err := db.Exec("DELETE FROM goals.main WHERE uid=?", ID)
	if err != nil {
		fmt.Println("problem with deleting goals index from DB: ", err)
	}
	defer db.Close()
	http.Redirect(w, r, "/goals", http.StatusSeeOther) //redirect back to goals page
}

//Main handles all routes===================================
func main() {
	r := mux.NewRouter()
	r.HandleFunc("/goals", Index).Methods("GET")
	r.HandleFunc("/goals", Create).Methods("POST")
	r.HandleFunc("/goals/new", New)
	r.HandleFunc("/goals/{id}", Show).Methods("GET")
	r.HandleFunc("/goals/{id}", Update).Methods("POST")
	r.HandleFunc("/goals/{id}/edit", Edit)
	r.HandleFunc("/goals/{id}/delete", Delete)
	http.ListenAndServe(":80", r)
}
