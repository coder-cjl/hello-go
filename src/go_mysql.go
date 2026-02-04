package main

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

type Person struct {
	UserId   int    `db:"user_id"`
	Username string `db:"username"`
	Sex      string `db:"sex"`
	Email    string `db:"email"`
}

func (Person) TableName() string {
	return "person"
}

type Place struct {
	Country string `db:"country"`
	City    string `db:"city"`
	TelCode int    `db:"telcode"`
}

var Db *sqlx.DB

func init() {
	database, err := sqlx.Open("mysql", "root:123456@~!@tcp(localhost:3306)/dev")
	if err != nil {
		Log.Error("open mysql failed")
		return
	}
	Db = database
}

func InsertPerson(p Person) error {
	resp, err := Db.NamedExec(`INSERT INTO person (username, sex, email) VALUES (:username, :sex, :email)`, &p)
	if err != nil {
		return err
	}

	id, err := resp.LastInsertId()
	if err != nil {
		return err
	}
	Log.Info("inserted person id:", id)

	return err
}

// func InsertPlace(pl Place) error {
// 	_, err := Db.NamedExec(`INSERT INTO place (country, city, telcode) VALUES (:country, :city, :telcode)`, &pl)
// 	return err
// }

type MySQL struct{}

// insert
func ts1() {
	p := Person{
		Username: "Mike",
		Sex:      "M",
		Email:    "mike@qq.com",
	}
	err := InsertPerson(p)
	if err != nil {
		Log.Error("InsertPerson failed:", err)
	} else {
		Log.Info("InsertPerson succeeded")
	}
}

// select
func ts2() {
	var person []Person
	err := Db.Select(&person, "SELECT user_id, username, sex, email FROM person WHERE user_id > ?", 0)
	if err != nil {
		Log.Error("Select person failed:", err)
		return
	}
	for _, p := range person {
		Log.Info("Person:", p)
	}
}

// update
func ts3() {
	resp, err := Db.Exec("UPDATE person set username=? where user_id=?", "孙悟空", 1001)
	if err != nil {
		Log.Error("Update person failed:", err)
		return
	}
	row, err := resp.RowsAffected()
	if err != nil {
		Log.Error("Get RowsAffected failed:", err)
		return
	}
	Log.Info("Updated rows:", row)
}

// delete
func ts4() {
	resp, err := Db.Exec("DELETE FROM person where user_id=?", 1002)
	if err != nil {
		Log.Error("Delete person failed:", err)
		return
	}
	row, err := resp.RowsAffected()
	if err != nil {
		Log.Error("Get RowsAffected failed:", err)
		return
	}
	Log.Info("Deleted rows:", row)
}

// 事务
func ts5() {
	tx, err := Db.Beginx()
	if err != nil {
		Log.Error("Begin transaction failed:", err)
		return
	}

	resp, err := tx.Exec("insert into person(username, sex, email) values (?, ?, ?)", "stu001", "M", "stu001@qq.com")
	if err != nil {
		tx.Rollback()
		Log.Error("Insert person failed:", err)
		return
	}
	id, err := resp.LastInsertId()
	if err != nil {
		tx.Rollback()
		Log.Error("Get LastInsertId failed:", err)
		return
	}
	Log.Info("Inserted person id in transaction:", id)

	err = tx.Commit()
	if err != nil {
		Log.Error("Commit transaction failed:", err)
		return
	}
	Log.Info("Transaction committed successfully")
}

// func ts6() {
// 	p := Person{
// 		Username: "Lucy",
// 		Sex:      "F",
// 		Email:    "lucy@qq.com",
// 	}
// 	result := Db
// }

func (m MySQL) Test() {
	ts5()
}

// pl := Place{
