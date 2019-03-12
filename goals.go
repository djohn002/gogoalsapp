package main

import (
	"goalsbasiclocalapp/server"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

//Main calls the MakeRouter function & creates server===================================
func main() {
	s, err := server.NewServer()
	if err != nil {
		log.Fatal(err)
	}
	defer s.Close()

	http.ListenAndServe(":80", s.MakeRouter())
}
