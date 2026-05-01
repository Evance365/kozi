package handlers

import (
    "database/sql"
    "net/http"
    "sort"

    "github.com/gin-gonic/gin"
    "github.com/Evance365/kozi/internal/models"
)

func GetMatches(db *sql.DB) gin.HandlerFunc {
    return func(c *gin.Context) {
        studentID := c.Query("student_id")
        level := c.Query("level")

        if studentID == "" {
            c.JSON(http.StatusBadRequest, gin.H{"error": "student_id is required"})
            return
        }

        rows, err := db.Query(
            "SELECT points FROM results WHERE student_id = $1",
            studentID,
        )
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "could not fetch results"})
            return
        }
        defer rows.Close()

        var allPoints []int
        for rows.Next() {
            var p int
            if err := rows.Scan(&p); err != nil {
                c.JSON(http.StatusInternalServerError, gin.H{"error": "could not read points"})
                return
            }
            allPoints = append(allPoints, p)
        }

        if len(allPoints) == 0 {
            c.JSON(http.StatusNotFound, gin.H{"error": "no results found for this student"})
            return
        }

        sort.Sort(sort.Reverse(sort.IntSlice(allPoints)))

        top := allPoints
        if len(top) > 7 {
            top = top[:7]
        }

        total := 0
        for _, p := range top {
            total += p
        }

        query := `
            SELECT
                c.name,
                c.level,
                i.name,
                COALESCE(i.location, ''),
                cu.cutoff_points,
                cu.year
            FROM cutoffs cu
            JOIN courses c ON c.id = cu.course_id
            JOIN institutions i ON i.id = cu.institution_id
            WHERE cu.cutoff_points <= $1
        `

        args := []interface{}{total}

        if level != "" {
            query += " AND c.level = $2"
            args = append(args, level)
        }

        query += " ORDER BY cu.cutoff_points DESC"

        matchRows, err := db.Query(query, args...)
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "could not fetch matches"})
            return
        }
        defer matchRows.Close()

        var matches []models.Match
        for matchRows.Next() {
            var m models.Match
            m.StudentPoints = total
            if err := matchRows.Scan(
                &m.CourseName,
                &m.CourseLevel,
                &m.InstitutionName,
                &m.Location,
                &m.CutoffPoints,
                &m.Year,
            ); err != nil {
                c.JSON(http.StatusInternalServerError, gin.H{"error": "could not read match"})
                return
            }
            matches = append(matches, m)
        }

        if len(matches) == 0 {
            c.JSON(http.StatusOK, gin.H{
                "student_points": total,
                "matches":        []models.Match{},
            })
            return
        }

        c.JSON(http.StatusOK, gin.H{
            "student_points": total,
            "matches":        matches,
        })
    }
}
