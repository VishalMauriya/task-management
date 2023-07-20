package main

import (
    "github.com/VishalMauriya/task-management/db"
    "github.com/VishalMauriya/task-management/handlers"

    "github.com/gin-gonic/gin"
)

func main() {
    // Initialized the database
    db.InitDB("./tasks.db")

    defer db.DB.Close()

    // Created a new Gogin router
    router := gin.Default()

    // API endpoints
    router.POST("/tasks", handlers.CreateTask)
    router.GET("/tasks/:id", handlers.GetTask)
    router.PUT("/tasks/:id", handlers.UpdateTask)
    router.DELETE("/tasks/:id", handlers.DeleteTask)
    router.GET("/tasks", handlers.ListTasks)
    
    // Starting the server
    router.Run("localhost:8080")
}
