package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-sql-driver/mysql"
)

var db *sql.DB

type Todo struct {
	Id         int    `json:"id"`
	UserId     int    `json:"user_id"`
	Title      string `json:"title"`
	CurrStatus string `json:"curr_status"`
}

type User struct {
	Id       int    `json:"id"`
	Username string `json:"username"`
}

func main() {
	// Capture connection properties.
	cfg := mysql.Config{
		User:   os.Getenv("DBUSER"),
		Passwd: os.Getenv("DBPASS"),
		Net:    "tcp",
		Addr:   "127.0.0.1:3306",
		DBName: "todoList",
	}
	var err error
	db, err = sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		log.Fatal(err)
	}
	pingErr := db.Ping()
	if pingErr != nil {
		log.Fatal(pingErr)
	}
	fmt.Println("Connected!")

	router := gin.Default()
	router.GET("/userapp/:userid", getuserList)
	router.POST("/users/:userid", createUserTodo)
	router.DELETE("/users/:userid", deleteUserTodo)

	router.Run("localhost:8080")
}

func getuserList(c *gin.Context) {
	user_id, err := strconv.Atoi(c.Param("userid"))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	rows, err := db.Query("SELECT id, user_id, title, curr_status FROM todos WHERE user_id = ?", user_id)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var todos []Todo
	for rows.Next() {
		var todo Todo
		if err := rows.Scan(&todo.Id, &todo.UserId, &todo.Title, &todo.CurrStatus); err != nil {
			log.Fatal(err)
		}
		todos = append(todos, todo)
	}

	c.JSON(http.StatusOK, todos)
}

func createUserTodo(c *gin.Context) {
	user_id, err := strconv.Atoi(c.Param("userid"))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	var todo Todo
	if err := c.BindJSON(&todo); err != nil {
		fmt.Errorf("Bind error %v", err)
	}
	todo.UserId = user_id

	result, err := db.Exec("INSERT INTO todos (user_id, title, curr_status) VALUES (?, ?, ?)",
		todo.UserId, todo.Title, todo.CurrStatus)
	if err != nil {
		log.Fatal(err)
	}
	lastInsertID, err := result.LastInsertId()
	if err != nil {
		log.Fatal(err)
	}
	todo.Id = int(lastInsertID)

	c.JSON(http.StatusCreated, todo)
}

func deleteUserTodo(c *gin.Context) {
	user_id, err := strconv.Atoi(c.Param("userid"))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}
	title := c.Query("title")
	result, err := db.Exec("DELETE FROM todos WHERE user_id = ? AND title = ?", user_id, title)
	if err != nil {
		log.Fatal(err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Fatal(err)
	}

	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Todo not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Todo deleted successfully"})
}
