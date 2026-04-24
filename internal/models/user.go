package models

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/xrcuo/xrcuo-lib/db"
	"golang.org/x/crypto/bcrypt"
)

// User 用户模型
type User struct {
	ID        int64     `json:"id"`
	Username  string    `json:"username"`
	Password  string    `json:"-"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// CreateUsersTable 创建用户表
func CreateUsersTable() error {
	sqlStr := `
	CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		username TEXT NOT NULL UNIQUE,
		password TEXT NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);
	`
	_, err := db.GetDB().Exec(sqlStr)
	return err
}

// UserCreate 创建用户
func UserCreate(username, password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("密码加密失败: %v", err)
	}

	_, err = db.GetDB().Exec(
		"INSERT INTO users (username, password) VALUES (?, ?)",
		username, string(hashedPassword),
	)
	if err != nil {
		return fmt.Errorf("创建用户失败: %v", err)
	}
	return nil
}

// UserGetByUsername 根据用户名获取用户
func UserGetByUsername(username string) (*User, error) {
	row := db.GetDB().QueryRow(
		"SELECT id, username, password, created_at, updated_at FROM users WHERE username = ?",
		username,
	)
	user := &User{}
	var createdAt, updatedAt sql.NullTime
	err := row.Scan(&user.ID, &user.Username, &user.Password, &createdAt, &updatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("用户不存在")
		}
		return nil, err
	}
	if createdAt.Valid {
		user.CreatedAt = createdAt.Time
	}
	if updatedAt.Valid {
		user.UpdatedAt = updatedAt.Time
	}
	return user, nil
}

// UserVerifyPassword 验证密码
func UserVerifyPassword(username, password string) (*User, error) {
	user, err := UserGetByUsername(username)
	if err != nil {
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return nil, fmt.Errorf("密码错误")
	}

	return user, nil
}

// UserUpdatePassword 修改密码
func UserUpdatePassword(username, oldPassword, newPassword string) error {
	_, err := UserVerifyPassword(username, oldPassword)
	if err != nil {
		return err
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("密码加密失败: %v", err)
	}

	_, err = db.GetDB().Exec(
		"UPDATE users SET password = ?, updated_at = CURRENT_TIMESTAMP WHERE username = ?",
		string(hashedPassword), username,
	)
	if err != nil {
		return fmt.Errorf("更新密码失败: %v", err)
	}
	return nil
}

// UserExists 检查用户是否存在
func UserExists() (bool, error) {
	var count int
	err := db.GetDB().QueryRow("SELECT COUNT(*) FROM users").Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
