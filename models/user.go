package models

import (
	"time"

	"gorm.io/gorm"
)

type UserLevel string

const (
	UserLevelAdmin   UserLevel = "ADMIN"
	UserLevelSupport UserLevel = "SUPPORT"
)

type UserRoleAction string

const (
	UserRoleAll    UserRoleAction = "ALL"
	UserRoleInsert UserRoleAction = "INSERT"
	UserRoleUpdate UserRoleAction = "UPDATE"
	UserRoleQuery  UserRoleAction = "QUERY"
)

type UserStatus string

const (
	UserStatusActive  UserStatus = "ACTIVE"
	UserStatusDisable UserStatus = "DISABLE"
	UserStatusSuspend UserStatus = "SUSPEND"
)

// User represents the user model
type UserTable struct {
	UserID     uint           `gorm:"primaryKey;column:user_id" json:"user_id"`
	Username   string         `gorm:"column:username" json:"username"`
	Telephone  string         `gorm:"column:telephone" json:"telephone"`
	Email      string         `gorm:"column:email" json:"email"`
	Password   string         `gorm:"column:password" json:"password"`
	Level      UserLevel      `gorm:"column:level;type:enum('ADMIN','SUPPORT');default:SUPPORT" json:"level,omitempty"`
	RoleAction UserRoleAction `gorm:"column:role_action;type:enum('ALL','INSERT','UPDATE','QUERY');default:QUERY" json:"role_action,omitempty"`
	Status     UserStatus     `gorm:"column:status;type:enum('ACTIVE','DISABLE','SUSPEND');default:ACTIVE" json:"status,omitempty"`
	OpID       uint           `gorm:"column:op_id" json:"op_id"`
	CreatedAt  time.Time      `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt  time.Time      `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"column:deleted_at" json:"deleted_at"`
}

type AuthUser struct {
	UserId     uint           `json:"user_id"`
	Telephone  string         `json:"telephone"`
	Email      string         `json:"email"`
	Username   string         `json:"username"`
	RoleAction UserRoleAction `json:"role_action"`
	Password   string         `json:"-"`
	OpId       *uint          `json:"op_id"`
	Level      UserLevel      `json:"level"`
	Status     UserStatus     `json:"-"`
}

type User struct {
	UserId     uint           `json:"user_id"`
	Telephone  string         `json:"telephone"`
	Email      string         `json:"email"`
	Username   string         `json:"username"`
	OpId       *uint          `json:"op_id"`
	Level      UserLevel      `json:"level"`
	RoleAction UserRoleAction `json:"role_action"`
	Status     UserStatus     `json:"status"`
	CreatedAt  string         `json:"created_at"`
	UpdatedAt  *string        `json:"updated_at"`
	DeletedAt  *string        `json:"deleted_at"`
}

type SaveUser struct {
	Telephone  string          `json:"telephone"`
	Email      string          `json:"email"`
	Username   string          `json:"username"`
	Level      *UserLevel      `gorm:"column:level" json:"level,omitempty"`
	RoleAction *UserRoleAction `gorm:"column:role_action" json:"role_action,omitempty"`
	Status     *UserStatus     `gorm:"column:status" json:"status,omitempty"`
	Password   *string         `json:"password,omitempty"`
}

type DeleteUser struct {
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

type SaveUserPassword struct {
	Status   string `json:"status"`
	Password string `json:"password"`
}

func (User) TableName() string {
	return "users"
}

func (AuthUser) TableName() string {
	return "users"
}

func (SaveUser) TableName() string {
	return "users"
}

func (SaveUserPassword) TableName() string {
	return "users"
}

func (DeleteUser) TableName() string {
	return "users"
}

func (UserTable) TableName() string {
	return "users"
}
