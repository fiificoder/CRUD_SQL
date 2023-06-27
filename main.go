package main

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

type Person struct {
	Id            int
	First_name    string
	Last_name     string
	Gender        string
	Date_of_birth string
	Email         string
}

var Db *sql.DB

func init() {
	var err error
	Db, err = sql.Open("postgres", "user=username dbname=databasename password=dbpassword sslmode=disable")
	if err != nil {
		panic(err)
	}
}

// getting all persons (a lmit of your choice)
func Persons(limit int) (persons []Person, err error) {
	rows, err := Db.Query("SELECT first_name, last_name, gender, date_of_birth, email FROM person LIMIT $1", limit)
	if err != nil {
		return
	}
	defer rows.Close()

	for rows.Next() {
		people := Person{}
		err = rows.Scan(&people.First_name, &people.Last_name, &people.Gender, &people.Date_of_birth, &people.Email)
		if err != nil {
			return
		}
		persons = append(persons, people)
	}
	return
}

// creating a psot
func (person *Person) Create() (err error) {
	statement := "INSERT INTO person (first_name, last_name, gender, date_of_birth, email) VALUES ($1, $2, $3, $4, $5) RETURNING id"
	stmt, err := Db.Prepare(statement)
	if err != nil {
		return
	}
	defer stmt.Close()

	err = stmt.QueryRow(person.First_name, person.Last_name, person.Gender, person.Date_of_birth, person.Email).Scan(&person.Id)
	return
}

// retrieving a person by id
func GetPerson(id int) (person Person, err error) {
	person = Person{}
	err = Db.QueryRow("SELECT id, first_name, last_name, gender, date_of_birth, email FROM person WHERE id = $1", id).Scan(&person.Id, &person.First_name, &person.Last_name, &person.Gender, &person.Date_of_birth, person.Email)
	return
}

// Update a person after retrieving or viewing it
func (person *Person) Update() (err error) {
	_, err = Db.Exec("UPDATE person SET first_name = $2, last_name = $3, gender = $4, date_of_birth = $5, email = $6 WHERE id = $1", person.Id, person.First_name, person.Last_name, person.Gender, person.Date_of_birth, person.Email)
	return
}

// delete a post after retrieving or viewing it
func (person *Person) Delete() (err error) {
	_, err = Db.Exec("DELETE FROM person WHERE id = $1", person.Id)
	return
}

func main() {
	// setting a post for our CREATE function
	person := Person{First_name: "Oprah", Last_name: "Winfrey", Gender: "Female", Date_of_birth: "2002-10-14", Email: "mavisbeacon43@gmail.com"}
	fmt.Println(person)

	person.Create()
	fmt.Println(person)

	// setting up our GetPerson function
	readPerson, _ := GetPerson(10)
	fmt.Println(readPerson)

	// setting up our Update function
	// first we use the GetPerson function to view before we update
	viewPerson, _ := GetPerson(16)
	fmt.Println(viewPerson)
	viewPerson.First_name = "Jack"
	viewPerson.Last_name = "Grealish"
	viewPerson.Gender = "Male"
	viewPerson.Date_of_birth = "1998-04-19"
	viewPerson.Email = "jackgrealish10.1@gmail.com"
	viewPerson.Update()

	// getting our preffered limit of persons
	persons, _ := Persons(10)
	fmt.Println(persons)

	// didn't call delete cause I really have nothing to delete :)
}
