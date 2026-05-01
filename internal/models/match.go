package models

type Match struct {
    CourseName      string `json:"course_name"`
    CourseLevel     string `json:"course_level"`
    InstitutionName string `json:"institution_name"`
    Location        string `json:"location"`
    CutoffPoints    int    `json:"cutoff_points"`
    StudentPoints   int    `json:"student_points"`
    Year            int    `json:"year"`
}
