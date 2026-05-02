package models

type Match struct {
    CourseName      string  `json:"course_name"`
    CourseLevel     string  `json:"course_level"`
    InstitutionName string  `json:"institution_name"`
    Location        string  `json:"location"`
    CutoffPoints    float64 `json:"cutoff_points"`
    StudentPoints   float64 `json:"student_points"`
    Year            int     `json:"year"`
}
