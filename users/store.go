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
		"name":        true,
		"email":       true,
		"address":     true,
		"phoneNumber": true,
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
