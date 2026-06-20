package service

import "portfolio-backend/internal/core/ports"

type HashManager struct {
	hashManager ports.HashManager
}

func NewHashManager(hashManager ports.HashManager) *HashManager {
	return &HashManager{hashManager: hashManager}
}

func (c *HashManager) Hash(text string) (string, error) {
	return c.hashManager.Hash(text)
}

func (c *HashManager) Compare(hashedText, textToCompare string) error {
	return c.hashManager.Compare(hashedText, textToCompare)
}
