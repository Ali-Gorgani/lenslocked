package main

import (
	"fmt"

	"github.com/Ali-Gorgani/lenslocked/models"
)

func main() {
	cfg := models.DefultPostgresConfig()
	db, err := models.Open(cfg)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	us := models.UserService{
		DB: db,
	}
	user, err := us.Create("test1@test.com", "test1")
	if err != nil {
		panic(err)
	}
	fmt.Println(user)

	// // Create a table
	// _, err = db.Exec(`
	// 	CREATE TABLE IF NOT EXISTS users (
	// 		id SERIAL PRIMARY KEY,
	// 		name TEXT,
	// 		email TEXT NOT NULL UNIQUE
	// 	);

	// 	CREATE TABLE IF NOT EXISTS orders (
	// 		id SERIAL PRIMARY KEY,
	// 		user_id INT NOT NULL,
	// 		amount INT,
	// 		description TEXT
	// 	);
	// `)
	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Println("Tables created!")

	// // Insert some data
	// type User struct {
	// 	Name  string
	// 	Email string
	// }

	// user := User{
	// 	Name:  "New User 1",
	// 	Email: "NewUser1@example.com",
	// }

	// // Select some data
	// row := db.QueryRow(`
	// 	INSERT INTO users (name, email)
	// 	VALUES
	// 	($1, $2) RETURNING id;`, user.Name, user.Email)
	// var id int
	// err = row.Scan(&id)
	// if err != nil {
	// 	panic(err)
	// }

	// fmt.Println("Inserted user with id", id)

	// userID := 1
	// for i := range(5) {
	// 	amount := (i+1) * 100
	// 	desc := fmt.Sprintf("Fake order #%d", i+1)
	// 	_, err := db.Exec(`
	// 		INSERT INTO orders(user_id, amount, description)
	// 		VALUES($1, $2, $3);`, userID, amount, desc)
	// 	if err != nil {
	// 		panic(err)
	// 	}
	// }
	// fmt.Println("Created fake orders.")

	// type Order struct {
	// 	ID int
	// 	UserID int
	// 	Amount int
	// 	Description string
	// }

	// var orders []Order
	// rows, err := db.Query(`
	// SELECT id, amount, description
	// FROM orders
	// WHERE user_id=$1;`, userID)
	// if err != nil {
	// 	panic(err)
	// }
	// defer rows.Close()

	// for rows.Next() {
	// 	var order Order
	// 	order.UserID = userID
	// 	err := rows.Scan(&order.ID, &order.Amount, &order.Description)
	// 	if err != nil {
	// 		panic(err)
	// 	}
	// 	orders = append(orders, order)
	// }
	// if rows.Err() != nil {
	// 	panic(rows.Err())
	// }
	// fmt.Println("Orders: ", orders)
}
