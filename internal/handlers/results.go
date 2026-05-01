package handlers

import (
    "database/sql"
    "net/http"

    "github.com/gin-gonic/gin"
    "github.com/Evance365/kozi/internal/models"
)

func PostResults(db *sql.DB) gin.HandlerFunc {
    return func(c *gin.Context) {
        var input models.StudentResults

        if err := c.ShouldBindJSON(&input); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
            return
        }

        if input.StudentID == "" {
            c.JSON(http.StatusBadRequest, gin.H{"error": "student_id is required"})
            return
        }

        if len(input.Results) == 0 {
            c.JSON(http.StatusBadRequest, gin.H{"error": "results cannot be empty"})
            return
        }

        tx, err := db.Begin()
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "could not start transaction"})
            return
        }

        _, err = tx.Exec(
            "DELETE FROM results WHERE student_id = $1",
            input.StudentID,
        )
        if err != nil {
            tx.Rollback()
            c.JSON(http.StatusInternalServerError, gin.H{"error": "could not clear old results"})
            return
        }

        for _, r := range input.Results {
            if r.Points < 0 || r.Points > 12 {
                tx.Rollback()
                c.JSON(http.StatusBadRequest, gin.H{"error": "points must be between 0 and 12"})
                return
            }

            _, err := tx.Exec(
                "INSERT INTO results (student_id, subject, grade, points) VALUES ($1, $2, $3, $4)",
                input.StudentID, r.Subject, r.Grade, r.Points,
            )
            if err != nil {
                tx.Rollback()
                c.JSON(http.StatusInternalServerError, gin.H{"error": "could not save result"})
                return
            }
        }

        if err := tx.Commit(); err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "could not commit results"})
            return
        }

        c.JSON(http.StatusCreated, gin.H{"message": "results saved"})
    }
}
