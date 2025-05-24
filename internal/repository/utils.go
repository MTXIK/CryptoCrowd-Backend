package repository

import (
	"crypto/rand"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/base64"
	"fmt"
	"strings"
)

// generateSalt генерирует случайный соль заданного размера
func generateSalt(size int) ([]byte, error) {
	salt := make([]byte, size)
	_, err := rand.Read(salt)
	if err != nil {
		return nil, err
	}

	return salt, nil
}

// hashPasswordSHA256 хеширует пароль с использованием SHA-256 и соли
func hashPasswordSHA256(password string) (string, error) {
	salt, err := generateSalt(16)
	if err != nil {
		return "", err
	}

	hash := sha256.New()
	hash.Write(salt)
	hash.Write([]byte(password))
	hashedPassword := hash.Sum(nil)

	return fmt.Sprintf("%s:%s",
		base64.StdEncoding.EncodeToString(salt),
		base64.StdEncoding.EncodeToString(hashedPassword)), nil
}

// checkPassword проверяет, соответствует ли введенный пароль хешу
func checkPassword(storedHash, password string) bool {
	parts := strings.Split(storedHash, ":")
	if len(parts) != 2 {
		return false
	}

	salt, err := base64.StdEncoding.DecodeString(parts[0])
	if err != nil {
		return false
	}

	storedPasswordHash, err := base64.StdEncoding.DecodeString(parts[1])
	if err != nil {
		return false
	}

	hash := sha256.New()
	hash.Write(salt)
	hash.Write([]byte(password))
	hashedPassword := hash.Sum(nil)

	return subtle.ConstantTimeCompare(hashedPassword, storedPasswordHash) == 1
}
