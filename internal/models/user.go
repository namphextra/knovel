package models

type Role string

const (
    RoleEmployee Role = "employee"
    RoleEmployer Role = "employer"
)

type User struct {
    ID       uint   `json:"id" gorm:"primaryKey"`
    Name     string `json:"name"`
    Email    string `json:"email" gorm:"unique"`
    Password string `json:"-"` // "-" ensures password is never sent in JSON responses
    Role     Role   `json:"role"`
} 