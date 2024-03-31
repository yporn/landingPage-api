package utils

import (
	"database/sql"
	"fmt"
	"strconv"

	"github.com/gofiber/fiber/v2"
)


func LogActivity(db *sql.DB, userID, action, details string) error {
 
    stmt, err := db.Prepare("INSERT INTO activity_logs (user_id, action, details) VALUES ($1, $2, $3)")
    if err != nil {
        return err
    }
    defer stmt.Close()

    // Execute the SQL statement with the provided parameters
    _, err = stmt.Exec(userID, action, details)
    if err != nil {
        return err
    }

    return nil
}

func GetUserIDFromContext(c *fiber.Ctx) int {
    userIDStr := c.Locals("userId")
    if userIDStr == nil {
        return 0
    }

    userID, err := strconv.Atoi(userIDStr.(string))
    if err != nil {
        // Handle error when conversion fails
        fmt.Println("Error converting userID to int:", err)
        return 0
    }

    return userID
}