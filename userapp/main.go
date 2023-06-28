package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-sql-driver/mysql"
	"github.com/golang-jwt/jwt"
)

const (
	secretKey       = "your-secret-key"
	tokenExpiration = time.Hour * 24 // Token expiration time (e.g., 24 hours)
)

type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type Todo struct {
	ID         int    `json:"id"`
	UserID     int    `json:"user_id"`
	Title      string `json:"title"`
	CurrStatus string `json:"curr_status"`
}

var db *sql.DB

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
	router.POST("/register", register)
	router.POST("/login", login)

	//Authenticated routes
	authGroup := router.Group("/api")
	authGroup.Use(authMiddleware())
	{
		authGroup.GET("/todos", getTodos)
		authGroup.POST("/todos", createTodo)
		//authGroup.DELETE("/todos", deleteTodo)
	}

	// Run the server
	router.Run(":8080")
}
func register(c *gin.Context) {
	var user User

	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	// Check if the username is already taken
	existingUser, err := getUserByUsername(user.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch user from database"})
		return
	}

	if existingUser != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Username is already taken"})
		return
	}

	err = insertUser(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to register user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User registered successfully"})

}

func login(c *gin.Context) {
	var credentials struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := c.ShouldBindJSON(&credentials); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}
	user, err := getUserByUsername(credentials.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch user from database"})
		return
	}

	if user == nil || user.Password != credentials.Password {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}
	// Generate the JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		Subject:   fmt.Sprintf("%d", user.ID),
		ExpiresAt: time.Now().Add(tokenExpiration).Unix(),
	})

	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		log.Fatal(err)
	}

	c.JSON(http.StatusOK, gin.H{"token": tokenString})
}

func authMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization token not provided"})
			c.Abort()
			return
		}
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return []byte(secretKey), nil
		})

		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token at token parsing step"})
			fmt.Println(err)
			c.Abort()
			return
		}
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token at claims step"})
			c.Abort()
			return
		}

		userID, ok := claims["sub"].(string)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token at userID step"})
			c.Abort()
			return
		}
		user, err := getUserByID(userID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch user from database"})
			c.Abort()
			return
		}

		if user == nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token at database fetch step"})
			c.Abort()
			return
		}
		c.Set("user", user)

		// Continue to the next handler
		c.Next()
	}

}

func getTodos(c *gin.Context) {
	// Get the authenticated user from the context
	//User user
	userInterface, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "User not found in context"})
		return
	}

	user, ok := userInterface.(*User)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user object"})
		return
	}

	// Fetch todos specific to the user from the database
	todos, err := getTodosByUserID(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch todos from database"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Fetching todos for user", "user": user, "todos": todos})

}

func getUserByUsername(username string) (*User, error) {
	query := "SELECT id, username, userpass FROM newusers WHERE username = ?"
	row := db.QueryRow(query, username)

	var user User
	err := row.Scan(&user.ID, &user.Username, &user.Password)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &user, nil
}

func getUserByID(userID string) (*User, error) {
	query := "SELECT id, username, userpass FROM newusers WHERE id = ?"
	row := db.QueryRow(query, userID)

	var user User
	err := row.Scan(&user.ID, &user.Username, &user.Password)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &user, nil
}
func getTodosByUserID(userID int) ([]Todo, error) {
	query := "SELECT id, user_id, title, curr_status FROM newtodos WHERE user_id = ?"
	rows, err := db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var todos []Todo
	for rows.Next() {
		var todo Todo
		err := rows.Scan(&todo.ID, &todo.UserID, &todo.Title, &todo.CurrStatus)
		if err != nil {
			return nil, err
		}
		todos = append(todos, todo)
	}

	return todos, nil
}

func insertUser(user User) error {
	query := "INSERT INTO newusers (username, userpass) VALUES (?, ?)"
	_, err := db.Exec(query, user.Username, user.Password)
	if err != nil {
		return err
	}

	return nil
}

func createTodo(c *gin.Context) {
	userInterface, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "User not found in context"})
		return
	}

	user, ok := userInterface.(*User)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user object"})
		return
	}
	var todo Todo
	if err := c.ShouldBindJSON(&todo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	// Set the user ID for the new todo
	todo.UserID = user.ID
	fmt.Println(todo.UserID)

	// Insert the new todo into the database
	err := insertTodo(todo)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create todo"})
		fmt.Println(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Creating todo for user", "user": user})
}
func insertTodo(todo Todo) error {
	query := "INSERT INTO newtodos (user_id, title, curr_status) VALUES (?,?,?)"
	_, err := db.Exec(query, todo.UserID, todo.Title, todo.CurrStatus)
	fmt.Println(todo.UserID, todo.Title, todo.CurrStatus)
	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}
