package actions

import (
	"database/sql"
	"errors"
	"htmx-practice/db"
	"htmx-practice/types"
	"log"

	"golang.org/x/crypto/bcrypt"
)

// Login performs a login action by querying the database for the user and comparing the password
// If the user is found and the password matches, the user is returned
// If the user is not found or the password does not match, an error is returned
func Login(email, passwd string) (*types.People, error) {
	conn, err := db.Init()
	if err != nil {
		log.Printf("Cannot create the db connection: %s\n", err)
		return nil, errors.New("service is down")
	}

	defer conn.Close()

	// define variables for the query
	var person types.People

	err = conn.QueryRow(
		`SELECT passwd, fname, lname, pid, rid FROM People WHERE email=?`,
		email,
	).Scan(
		&person.Passwd,
		&person.Fname,
		&person.Lname,
		&person.Pid,
		&person.Rid,
	)

	switch {
	case err == sql.ErrNoRows:
		log.Printf("No user with email: %s", email)
		return &types.People{}, errors.New("incorrect email or password")
	case err != nil:
		log.Printf("Error querying the database: %s", err)
		return &types.People{}, errors.New("service is down")
	}

	err = bcrypt.CompareHashAndPassword([]byte(person.Passwd), []byte(passwd))
	if err != nil {
		log.Printf("Password does not match: %s", err)
		return &types.People{}, errors.New("incorrect email or password")
	}

	log.Printf("User logged in: %s", email)
	person.Passwd = ""
	return &person, nil
}
