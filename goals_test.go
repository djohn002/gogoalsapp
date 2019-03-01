package main

import (
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
