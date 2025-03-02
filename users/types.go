package main

import (
	"context"
	"github.com/quanbin27/commons/genproto/users"
	"time"
)

type UserStore interface {
	GetUserByID(ctx context.Context, id int32) (*User, error)
	GetUserByEmail(ctx context.Context, email string) (*User, error)
	CreateUser(ctx context.Context, user *User) error
	UpdateInfo(ctx context.Context, userID int32, updatedData map[string]interface{}) error
	UpdatePassword(ctx context.Context, userID int32, password string) error
	GetNameByID(ctx context.Context, id int32) (string, error)
	GetUsersByIDs(ctx context.Context, userIDs []int32) ([]User, error)
}
type UserService interface {
	Register(ctx context.Context, user *users.RegisterRequest) (*users.RegisterResponse, error)
	Login(ctx context.Context, login *users.LoginRequest) (*users.LoginResponse, error)
	ChangeInfo(ctx context.Context, update *users.ChangeInfoRequest) (*users.ChangeInfoResponse, error)
	ChangePassword(ctx context.Context, update *users.ChangePasswordRequest) (*users.ChangePasswordResponse, error)
	GetUserInfo(ctx context.Context, id *users.GetUserInfoRequest) (*users.User, error)
	GetUserInfoByEmail(ctx context.Context, email *users.GetUserInfoByEmailRequest) (*users.User, error)
}

type User struct {
	ID        int32     `gorm:"primaryKey"`
	Email     string    `gorm:"unique;not null"`
	Password  string    `gorm:"not null"`
	Name      string    `gorm:"not null"`
	Roles     []Role    `gorm:"many2many:user_roles;"` // Quan hệ nhiều-nhiều với Role
	BranchID  *int32    `gorm:"index"`                 // ID của chi nhánh hiện tại
	CreatedAt time.Time `gorm:"autoCreateTime"`
}

// Role - Bảng quyền (Admin, Employee, Customer)
type Role struct {
	ID    int32  `gorm:"primaryKey"`
	Name  string `gorm:"unique;not null"`
	Users []User `gorm:"many2many:user_roles;"`
}

// UserRole - Bảng trung gian giữa User và Role
type UserRole struct {
	UserID int32 `gorm:"primaryKey"`
	RoleID int32 `gorm:"primaryKey"`
}

// EmployeeBranch - Quản lý nhân viên thuộc chi nhánh nào
type EmployeeBranch struct {
	UserID   int32 `gorm:"primaryKey"`
	BranchID int32 `gorm:"primaryKey"`
}
