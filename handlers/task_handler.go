package handlers

import (
    "database/sql"
    "net/http"
    "strconv"

    "github.com/gin-gonic/gin"
    "github.com/VishalMauriya/task-management/db"
    "github.com/VishalMauriya/task-management/models"
)

// Create a new task
func CreateTask(c *gin.Context) {
    var task models.Task
    if err := c.ShouldBindJSON(&task); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
        return
    }

    // Save the task to the database
    result, err := db.DB.Exec("INSERT INTO tasks (title, description, due_date, status) VALUES (?, ?, ?, ?)",
        task.Title, task.Description, task.DueDate, task.Status)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create task"})
        return
    }

    // Get the last inserted ID and assign it to the task
    lastInsertID, err := result.LastInsertId()
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create task"})
        return
    }
    
    // Convert int64 to int
    task.ID = int(lastInsertID)

    c.JSON(http.StatusCreated, task)
}

// Retrieve a task
func GetTask(c *gin.Context) {
    taskID, err := strconv.Atoi(c.Param("id"))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid task ID"})
        return
    }

    // Get the task from the database
    var task models.Task
    err = db.DB.QueryRow("SELECT * FROM tasks WHERE id=?", taskID).
        Scan(&task.ID, &task.Title, &task.Description, &task.DueDate, &task.Status)
    if err != nil {
        if err == sql.ErrNoRows {
            c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
            return
        }
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get task"})
        return
    }

    c.JSON(http.StatusOK, task)
}

// Update a task
func UpdateTask(c *gin.Context) {
    taskID, err := strconv.Atoi(c.Param("id"))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid task ID"})
        return
    }

    // Check if the task exists in the database before attempting update
    var existingTaskID int
    err = db.DB.QueryRow("SELECT id FROM tasks WHERE id=?", taskID).Scan(&existingTaskID)
    if err != nil {
        if err == sql.ErrNoRows {
            c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
            return
        }
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check task existence"})
        return
    }

    var updatedTask models.Task
    if err := c.ShouldBindJSON(&updatedTask); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
        return
    }

    // Update the task in the database
    _, err = db.DB.Exec("UPDATE tasks SET title=?, description=?, due_date=?, status=? WHERE id=?",
        updatedTask.Title, updatedTask.Description, updatedTask.DueDate, updatedTask.Status, taskID)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update task"})
        return
    }

    // Get the updated task from the database
    var task models.Task
    err = db.DB.QueryRow("SELECT * FROM tasks WHERE id=?", taskID).
        Scan(&task.ID, &task.Title, &task.Description, &task.DueDate, &task.Status)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve updated task details"})
        return
    }

    c.JSON(http.StatusOK, task)

}

// Delete a task
func DeleteTask(c *gin.Context) {
    taskID, err := strconv.Atoi(c.Param("id"))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid task ID"})
        return
    }

    // Check if the task exists in the database before attempting deletion
    var existingTaskID int
    err = db.DB.QueryRow("SELECT id FROM tasks WHERE id=?", taskID).Scan(&existingTaskID)
    if err != nil {
        if err == sql.ErrNoRows {
            c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
            return
        }
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check task existence"})
        return
    }

    // Delete the task from the database
    _, err = db.DB.Exec("DELETE FROM tasks WHERE id=?", taskID)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete task"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "Task deleted successfully"})
}

// List all tasks
func ListTasks(c *gin.Context) {
    rows, err := db.DB.Query("SELECT * FROM tasks")
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get tasks"})
        return
    }
    defer rows.Close()

    var tasks []models.Task
    for rows.Next() {
        var task models.Task
        err := rows.Scan(&task.ID, &task.Title, &task.Description, &task.DueDate, &task.Status)
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get tasks"})
            return
        }
        tasks = append(tasks, task)
    }

    if tasks == nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Tasks not found"})
        return
    }

    c.JSON(http.StatusOK, tasks)
}
