package models

type Result struct {
    Subject string `json:"subject"`
    Grade   string `json:"grade"`
    Points  int    `json:"points"`
}

type StudentResults struct {
    StudentID string   `json:"student_id"`
    Results   []Result `json:"results"`
}
