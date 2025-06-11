package main

import (
	"context"
	"errors"
	"fmt"
	"gorm.io/gorm"
)

type Store struct {
	db *gorm.DB
}

func NewStore(db *gorm.DB) *Store {
	return &Store{db}
}

func (s *Store) GetUserByEmail(ctx context.Context, email string) (*User, error) {
	var user User
	result := s.db.WithContext(ctx).Where("email = ?", email).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}
func (s *Store) CreateRole(ctx context.Context, userId int32, roleId int32) error {
	err := s.db.WithContext(ctx).Create(&UserRole{UserID: userId, RoleID: roleId}).Error
	return err
}
func (s *Store) UpdateRole(ctx context.Context, userId int32, roleId int32) error {
	err := s.db.WithContext(ctx).Save(&UserRole{UserID: userId, RoleID: roleId}).Error
	return err
}
func (s *Store) GetNameByID(ctx context.Context, id int32) (string, error) {
	var user User
	result := s.db.WithContext(ctx).Where("id = ?", id).First(&user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return "", fmt.Errorf("no user found with ID %d", id)
		}
		return "", result.Error
	}
	return user.Name, nil
}

func (s *Store) GetUsersByIDs(ctx context.Context, userIDs []int32) ([]User, error) {
	var users []User
	err := s.db.WithContext(ctx).Where("id IN ?", userIDs).Find(&users).Error
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (s *Store) GetUserByID(ctx context.Context, id int32) (*User, error) {
	var user User
	result := s.db.WithContext(ctx).Unscoped().Where("id = ?", id).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}
func (s *Store) GetRole(ctx context.Context, id int32) (int32, error) {
	var userRole UserRole
	if err := s.db.WithContext(ctx).Where("user_id = ?", id).First(&userRole).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return 0, nil // Không tìm thấy role nhưng không phải lỗi
		}
		return 0, err // Lỗi khác khi truy vấn database
	}
	return userRole.RoleID, nil
}

func (s *Store) CreateUser(ctx context.Context, user *User) (int32, error) {
	result := s.db.WithContext(ctx).Create(&user)
	return user.ID, result.Error
}

func (s *Store) UpdateInfo(ctx context.Context, userID int32, updatedData map[string]interface{}) error {
	allowedFields := map[string]bool{
		"name":         true,
		"email":        true,
		"address":      true,
		"phone_number": true,
	}
	for key := range updatedData {
		if !allowedFields[key] {
			delete(updatedData, key) // Xóa các trường không hợp lệ
		}
	}

	result := s.db.WithContext(ctx).Model(&User{}).Where("id = ?", userID).Updates(updatedData)
	return result.Error
}

func (s *Store) UpdatePassword(ctx context.Context, userID int32, password string) error {
	return s.db.WithContext(ctx).Model(&User{}).Where("id = ?", userID).Update("password", password).Error
}
func (s *Store) GetAllCustomers(ctx context.Context) ([]User, error) {
	var users []User
	err := s.db.WithContext(ctx).
		Joins("JOIN user_roles ON users.id = user_roles.user_id").
		Where("user_roles.role_id = ?", 1).
		Find(&users).Error
	if err != nil {
		return nil, err
	}
	return users, nil
}
func (s *Store) GetAllUsers(ctx context.Context) ([]UserWithRole, error) {
	var users []User
	err := s.db.WithContext(ctx).
		Preload("Roles").
		Find(&users).Error
	if err != nil {
		return nil, err
	}

	var result []UserWithRole
	for _, user := range users {
		var maxRoleID int32
		for _, role := range user.Roles {
			if role.ID > maxRoleID {
				maxRoleID = role.ID
			}
		}

		result = append(result, UserWithRole{
			ID:          user.ID,
			Email:       user.Email,
			Name:        user.Name,
			PhoneNumber: user.PhoneNumber,
			Address:     user.Address,
			RoleID:      maxRoleID,
			BranchID:    user.BranchID,
		})
	}

	return result, nil
}
func (s *Store) UpdateUser(ctx context.Context, input UserWithRole) error {
	// 1. Tìm user theo ID
	var user User
	if err := s.db.WithContext(ctx).Preload("Roles").First(&user, input.ID).Error; err != nil {
		return fmt.Errorf("user not found: %w", err)
	}

	// 2. Cập nhật thông tin cơ bản
	user.Email = input.Email
	user.Name = input.Name
	user.PhoneNumber = input.PhoneNumber
	user.Address = input.Address
	user.BranchID = input.BranchID

	// 3. Ghi lại thông tin user
	if err := s.db.WithContext(ctx).Save(&user).Error; err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	// 4. Cập nhật vai trò (Roles): thay toàn bộ roles bằng role mới
	if err := s.db.WithContext(ctx).Model(&user).Association("Roles").Clear(); err != nil {
		return fmt.Errorf("failed to clear roles: %w", err)
	}

	if input.RoleID > 0 {
		var role Role
		if err := s.db.WithContext(ctx).First(&role, input.RoleID).Error; err != nil {
			return fmt.Errorf("invalid role ID %d: %w", input.RoleID, err)
		}
		if err := s.db.WithContext(ctx).Model(&user).Association("Roles").Append(&role); err != nil {
			return fmt.Errorf("failed to assign role: %w", err)
		}
	}

	return nil
}

// Get customers with pagination
func (s *Store) GetCustomersPaginated(ctx context.Context, page int32, pageSize int32) ([]User, int64, error) {
	var users []User
	var total int64

	query := s.db.WithContext(ctx).
		Joins("JOIN user_roles ON users.id = user_roles.user_id").
		Where("user_roles.role_id = ?", 1)

	// Count total records for pagination
	if err := query.Model(&User{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Apply pagination
	offset := (page - 1) * pageSize
	err := query.
		Limit(int(pageSize)).
		Offset(int(offset)).
		Find(&users).Error
	if err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

// Get customers filtered by name
func (s *Store) GetCustomersByName(ctx context.Context, nameFilter string) ([]User, error) {
	var users []User
	err := s.db.WithContext(ctx).
		Joins("JOIN user_roles ON users.id = user_roles.user_id").
		Where("user_roles.role_id = ?", 1).
		Where("users.name LIKE ?", "%"+nameFilter+"%").
		Find(&users).Error
	if err != nil {
		return nil, err
	}
	return users, nil
}
func (s *Store) GetBranchByEmployeeID(ctx context.Context, userID int32) (int32, error) {
	var eb User
	err := s.db.WithContext(ctx).Where("id = ?", userID).First(&eb).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return 0, fmt.Errorf("no branch found for user ID %d", userID)
		}
		return 0, err
	}
	if eb.BranchID == nil {
		return -1, fmt.Errorf("no branch associated with user ID %d", userID)
	}
	branchID := *eb.BranchID
	return branchID, nil
}
