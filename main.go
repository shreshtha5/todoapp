package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/go-sql-driver/mysql"
)

var db *sql.DB

type app struct {
	Id          int    `json:"id"`
	Title       string `json:"title"`
	Curr_status string `json:"curr_status"`
}

func main() {
	// Capture connection properties.
	fmt.Println("reached here")
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
	router.GET("/todoapp", getList)
	router.PUT("/todoapp/:id", updateList)
	router.POST("/todoapp", putList)
	router.DELETE("/todoapp/:id", deleteitem)
	router.Run("localhost:8080")
}

func getList(c *gin.Context) {
	var apps []app
	rows, err := db.Query("select * from todoapp")

	fmt.Printf("type %T", rows)

	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		var alb app
		if err := rows.Scan(&alb.Id, &alb.Title, &alb.Curr_status); err != nil {
			log.Fatal(err)
		}
		apps = append(apps, alb)

	}
	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}
	c.JSON(http.StatusOK, apps)
}

func putList(c *gin.Context) {
	var item app
	if err := c.BindJSON(&item); err != nil {
		fmt.Errorf("Bind error %v", err)
	}
	result, err := db.Exec("INSERT INTO todoapp (id,title, curr_status) VALUES (?, ?, ?)", item.Id, item.Title, item.Curr_status)
	if err != nil {
		fmt.Errorf("putList: %v", err)
	}
	id, err := result.LastInsertId()

	if err != nil {
		fmt.Errorf("putList: %v", err)
	}
	fmt.Printf("ID of added album: %v\n", id)
	c.IndentedJSON(http.StatusCreated, item)
}

func updateList(c *gin.Context) {
	//var alb app
	/*var item app
	if err := c.BindJSON(&item); err != nil {
		fmt.Errorf("Bind error %v", err)
	}*/
	id := c.Param("id")
	fmt.Println("id from parameters ts %v", id)

	result, err := db.Exec("UPDATE todoapp SET curr_status = 'Done' WHERE id == ?", id)
	i, err := result.RowsAffected()
	if err != nil {
		fmt.Errorf("putList: %v", err)
	}
	fmt.Printf("ID of updated row is : %v\n", i)
	/*row, err := db.Query("select * from todoapp where id == ?", item.Id)
	if err != nil {
		fmt.Errorf("updatelist %v", err)
	}

	if err := row.Scan(&alb.Id, &alb.Title, &alb.Curr_status); err != nil {
		if err == sql.ErrNoRows {
			fmt.Errorf("albumsById no such album")
		}
		fmt.Errorf("albumsById")
	}

	c.IndentedJSON(http.StatusCreated, alb)*/
}

func deleteitem(c *gin.Context) {
	id := c.Param("id")
	fmt.Println("id from parameters ts %v", id)

	row, err := db.Exec("DELETE FROM todoapp WHERE id = ?", id)
	if err != nil {
		log.Fatal(err)
	}
	rowsAffected, err := row.RowsAffected()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Record deleted successfully %v", rowsAffected)
}
