package server

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

//Basic unit test that tests database connection.  Calls DBconn().
//I tested this test by changing DB name to ensure that it failed.
func TestDbconn(t *testing.T) {
	_, err := Dbconn()
	if err != nil {
		t.Errorf("Failed to connect to DB: %s", err)
	}
}

func TestIndexHandler(t *testing.T) {
	s, err := NewServer()
	if err != nil {
		t.Fatal(err)
	}
	defer s.Close()

	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/goals", nil)
	s.Index(w, r) //calls Index function & send w & r

	if w.Code != http.StatusOK {
		t.Error(w.Code, string(w.Body.String()))
	}
	fmt.Print(w.Code) //prints status code
}
