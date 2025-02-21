package main

import (
	"gorm.io/gorm"
)

type store struct {
	db *gorm.DB
}

func NewStore(db *gorm.DB) *store {
	return &store{db: db}
}
func (s *store) Create() error {
	return nil
}
