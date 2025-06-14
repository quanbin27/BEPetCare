package main

import (
	"context"
	"time"

	"github.com/quanbin27/commons/genproto/users"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// UserStore defines the interface for data storage operations
type UserStore interface {
	GetUserByID(ctx context.Context, id int32) (*User, error)
	GetUserByEmail(ctx context.Context, email string) (*User, error)
	CreateUser(ctx context.Context, user *User) (int32, error)
	UpdateInfo(ctx context.Context, userID int32, updatedData map[string]interface{}) error
	UpdatePassword(ctx context.Context, userID int32, password string) error
	GetNameByID(ctx context.Context, id int32) (string, error)
	GetUsersByIDs(ctx context.Context, userIDs []int32) ([]User, error)
	CreateRole(ctx context.Context, userId int32, roleId int32) error
	UpdateRole(ctx context.Context, userId int32, roleId int32) error
	GetRole(ctx context.Context, userId int32) (int32, error)
	GetAllCustomers(ctx context.Context) ([]User, error)
	GetAllUsers(ctx context.Context) ([]UserWithRole, error)
	UpdateUser(ctx context.Context, input UserWithRole) error
	GetCustomersPaginated(ctx context.Context, page int32, pageSize int32) ([]User, int64, error)
	GetCustomersByName(ctx context.Context, nameFilter string) ([]User, error)
	GetBranchByEmployeeID(ctx context.Context, userID int32) (int32, error)
}

// UserService defines the interface for business logic operations with internal types
type UserService interface {
	Register(ctx context.Context, email, password, name string) (string, error)
	Login(ctx context.Context, email, password string, rememberMe bool) (string, string, error)                                             // Trả về status và token
	ChangeInfo(ctx context.Context, userID int32, email, name, address, phoneNumber string) (string, string, string, string, string, error) // Trả về status, email, name
	ChangePassword(ctx context.Context, userID int32, oldPassword, newPassword string) (string, error)                                      // Trả về status
	GetUserInfo(ctx context.Context, id int32) (*User, error)
	GetUserInfoByEmail(ctx context.Context, email string) (*User, error)
	VerifyEmail(ctx context.Context, token string) (int32, error)
	ForgotPassword(ctx context.Context, email string) error
	ResetPassword(ctx context.Context, userID int32, token, newPassword string) error
	GetAllCustomers(ctx context.Context) ([]User, error)
	GetAllUsers(ctx context.Context) ([]UserWithRole, error)
	EditUser(ctx context.Context, input UserWithRole) error
	GetCustomersPaginated(ctx context.Context, page int32, pageSize int32) ([]User, int64, error)
	GetCustomersByName(ctx context.Context, nameFilter string) ([]User, error)
	GetBranchByEmployeeID(ctx context.Context, employeeID int32) (int32, error)
	CreateUser(ctx context.Context, email, name, phoneNumber string) (int32, error)
}

// User represents a user in the internal system
type User struct {
	ID          int32  `gorm:"primaryKey"`
	Email       string `gorm:"unique;not null"`
	Password    string `gorm:"not null"`
	Name        string `gorm:"not null"`
	PhoneNumber string `gorm:"size:11"`
	Address     string
	Roles       []Role    `gorm:"many2many:user_roles;"` // Quan hệ nhiều-nhiều với Role
	BranchID    *int32    `gorm:"index"`                 // ID của chi nhánh hiện tại
	CreatedAt   time.Time `gorm:"autoCreateTime"`
}
type PendingUser struct {
	Email    string
	Password string
	Name     string
	Token    string
	Expires  time.Time
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
type UserWithRole struct {
	ID          int32  `json:"id"`
	Email       string `json:"email"`
	Name        string `json:"name"`
	PhoneNumber string `json:"phone_number"`
	Address     string `json:"address"`
	RoleID      int32  `json:"role_id"`
	BranchID    *int32 `json:"branch_id,omitempty"`
}

// Helper functions to convert between internal User and protobuf User
func toProtoUser(u *User) *users.User {
	return &users.User{
		ID:          u.ID,
		Email:       u.Email,
		Name:        u.Name,
		PhoneNumber: u.PhoneNumber,
		Address:     u.Address,
		CreatedAt:   timestamppb.New(u.CreatedAt),
	}
}

func fromProtoRegisterRequest(req *users.RegisterRequest) *User {
	return &User{
		Email:    req.Email,
		Password: req.Password, // Lưu ý: Service sẽ mã hóa password
		Name:     req.Name,
	}
}
