package service

import (
	"github.com/tdmdh/lornian-backend/services/auth-service/internal/types"
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(hash), nil
}

func ComparePasswords(hashed string, plain []byte) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashed), plain)
	return err == nil
}

func ResetPassword(hashedPassword, newPassword string) (string, error) {
	if err := ValidatePasswordStrength(newPassword); err != nil {
		return "", err
	}

	if ComparePasswords(hashedPassword, []byte(newPassword)) {
		return "", types.ErrSamePassword
	}

	hashed, err := HashPassword(newPassword)
	if err != nil {
		return "", err
	}

	return hashed, nil
}

func ValidatePasswordStrength(password string) error {
	if len(password) < 8 {
		return types.ErrPasswordTooWeak
	}

	hasUpper := false
	hasLower := false
	hasDigit := false
	hasSpecial := false

	for _, char := range password {
		switch {
		case char >= 'A' && char <= 'Z':
			hasUpper = true
		case char >= 'a' && char <= 'z':
			hasLower = true
		case char >= '0' && char <= '9':
			hasDigit = true
		case char >= 32 && char <= 126: 
			if !((char >= 'A' && char <= 'Z') || (char >= 'a' && char <= 'z') || (char >= '0' && char <= '9')) {
				hasSpecial = true
			}
		}
	}

	requiredTypes := 0
	if hasUpper {
		requiredTypes++
	}
	if hasLower {
		requiredTypes++
	}
	if hasDigit {
		requiredTypes++
	}
	if hasSpecial {
		requiredTypes++
	}

	if requiredTypes < 3 {
		return types.ErrPasswordTooWeak
	}

	return nil
}

func CompareCurrentPassword(hashedPassword, currentPassword string) error {
	if !ComparePasswords(hashedPassword, []byte(currentPassword)) {
		return types.ErrIncorrectCurrentPassword
	}
	return nil
}

func ChangePassword(hashedPassword, currentPassword, newPassword string) (string, error) {
	if err := CompareCurrentPassword(hashedPassword, currentPassword); err != nil {
		return "", err
	}

	newHashedPassword, err := ResetPassword(hashedPassword, newPassword)
	if err != nil {
		return "", err
	}

	return newHashedPassword, nil
}
