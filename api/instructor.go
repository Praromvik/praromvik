package api

type Instructor struct {
	InstructorId string `json:"instructorId"`
	Name         string `json:"name"`
	Phone        string `json:"phone"`
	Email        string `json:"email"`
}
