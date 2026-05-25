package handlers

import (
	"fmt"
	"go-fiber-api/models"
	"go-fiber-api/utils"
	"log"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// type UserHandler struct {
// 	db *gorm.DB
// }

// func NewUserHandler(db *gorm.DB) *UserHandler {
// 	return &UserHandler{db: db}
// }

func (h *ApiHandler) GetUsers(c *fiber.Ctx) error {
	page := c.Query("page", "1")
	text := c.Query("text")
	deleted := c.Query("get_deleted") == "true"

	var pageInt int
	fmt.Sscanf(page, "%d", &pageInt)
	if pageInt < 1 {
		pageInt = 1
	}

	pageSize := 10
	offset := (pageInt - 1) * pageSize

	var users []models.User
	var total int64
	var col string = `user_id, username, telephone, email, op_id, level, role_action, status, 
	DATE_FORMAT(created_at, '%d/%m/%Y %r') as created_at, 
	DATE_FORMAT(updated_at, '%d/%m/%Y %r') as updated_at,
	CASE
		WHEN deleted_at IS NULL THEN NULL
		ELSE DATE_FORMAT(deleted_at, '%d/%m/%Y %r')
	END as deleted_at`

	query := h.db.Model(&models.User{}).Select(col)

	if deleted {
		query.Unscoped()
	}
	if text != "" {
		text = "%" + text + "%"
		query = query.Where("username LIKE ? OR telephone LIKE ?", text, text)
	}

	query.Count(&total)

	query.Offset(offset).Limit(pageSize).Find(&users)

	return c.JSON(fiber.Map{
		"data":        users,
		"page":        pageInt,
		"total":       total,
		"total_pages": (total + int64(pageSize) - 1) / int64(pageSize),
	})
}

func (h *ApiHandler) GetUsersRaw(c *fiber.Ctx) error {
	page := c.Query("page", "1")
	text := c.Query("text")
	deleted := c.Query("get_deleted") == "true"

	var pageInt int
	fmt.Sscanf(page, "%d", &pageInt)

	if pageInt < 1 {
		pageInt = 1
	}

	pageSize := 10
	offset := (pageInt - 1) * pageSize

	var users []models.User
	var total int64

	baseQuery := `
		FROM users
		WHERE 1=1
	`

	var args []interface{}

	// exclude soft deleted
	if !deleted {
		baseQuery += " AND deleted_at IS NULL"
	}

	// search
	if text != "" {
		baseQuery += `
			AND (
				username LIKE ?
				OR telephone LIKE ?
			)
		`

		search := "%" + text + "%"
		args = append(args, search, search)
	}

	// count query
	countQuery := "SELECT COUNT(*) " + baseQuery

	err := h.db.Raw(countQuery, args...).Scan(&total).Error
	if err != nil {
		fmt.Printf("Failed to count users: %v\n", err)

		return c.JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to get users",
		})
	}

	// data query
	dataQuery := `SELECT
		user_id,
		username,
		telephone,
		email,
		op_id,
		level,
		role_action,
		status,
		DATE_FORMAT(created_at, '%d/%m/%Y %r') as created_at, 
		DATE_FORMAT(updated_at, '%d/%m/%Y %r') as updated_at,
		DATE_FORMAT(deleted_at, '%d/%m/%Y %r') as deleted_at
	` + baseQuery + `
		LIMIT ?
		OFFSET ?`

	args = append(args, pageSize, offset)

	err = h.db.Raw(dataQuery, args...).Scan(&users).Error
	if err != nil {
		fmt.Printf("Failed to get users: %v\n", err)

		return c.JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to get users",
		})
	}

	return c.JSON(fiber.Map{
		"data":        users,
		"page":        pageInt,
		"total":       total,
		"total_pages": (total + int64(pageSize) - 1) / int64(pageSize),
	})
}

func (h *ApiHandler) GetUser(c *fiber.Ctx) error {
	id := c.Params("id")
	deleted := c.Query("get_deleted") == "true"
	var col string = `user_id, username, telephone, email, op_id, level, role_action, status, 
	DATE_FORMAT(created_at, '%d/%m/%Y %r') as created_at, 
	DATE_FORMAT(updated_at, '%d/%m/%Y %r') as updated_at,
	CASE
		WHEN deleted_at IS NULL THEN NULL
		ELSE DATE_FORMAT(deleted_at, '%d/%m/%Y %r')
	END as deleted_at`
	var user models.User
	query := h.db.Select(col).First(&user, id)
	if deleted {
		query.Unscoped()
	}
	err := query.Error
	if err != nil {
		return c.JSON(fiber.Map{"status": "error", "message": "User not found", "data": user})
	}
	return c.JSON(fiber.Map{"status": "success", "data": user})
}

func (h *ApiHandler) CreateUser(c *fiber.Ctx) error {
	var user models.SaveUser
	if err := c.BodyParser(&user); err != nil {
		log.Printf("Failed to read request body: %v\n", err)
		return c.JSON(fiber.Map{"status": "error", "message": "Invalid request body"})
	}

	if user.Password != nil || *user.Password != "" {
		hash, err := utils.HashPassword(*user.Password)
		if err != nil {
			log.Printf("Failed to hash password: %v\n", err)
			return c.JSON(fiber.Map{"status": "error", "message": "Failed to create user"})
		}
		user.Password = &hash
	}

	if user.Level == nil || *user.Level == "" {
		lvl := models.UserLevelSupport
		user.Level = &lvl
	}

	if user.RoleAction == nil || *user.RoleAction == "" {
		ra := models.UserRoleQuery
		user.RoleAction = &ra
	}

	if user.Status == nil || *user.Status == "" {
		status := models.UserStatusActive
		user.Status = &status
	}

	err := h.db.Create(&user).Error
	if err != nil {
		log.Printf("Failed to create user: %v\n", err)
		return c.JSON(fiber.Map{"status": "error", "message": "Failed to create user"})
	}

	return c.JSON(fiber.Map{"status": "success", "message": "User created successfully"})
}

func (h *ApiHandler) UpdateUser(c *fiber.Ctx) error {
	id := c.Params("id")
	var user models.UserTable
	if err := h.db.First(&user, id).Error; err != nil {
		log.Printf("User not found: %v\n", err)
		return c.JSON(fiber.Map{"status": "error", "message": "User not found"})
	}
	var input models.SaveUser
	if err := c.BodyParser(&input); err != nil {
		log.Printf("Failed to read request body: %v\n", err)
		return c.JSON(fiber.Map{"status": "error", "message": "Invalid request body"})
	}

	if input.Password != nil && *input.Password != "" {
		valid := utils.CheckPassword(*input.Password, user.Password)
		if !valid {
			hash, err := utils.HashPassword(*input.Password)
			if err != nil {
				log.Printf("Failed to hash password: %v\n", err)
				c.JSON(fiber.Map{"status": "error", "message": "Failed to update user"})
			}
			input.Password = &hash
		}
	} else {
		input.Password = &user.Password
	}

	if input.Level == nil || *input.Level == "" {
		input.Level = &user.Level
	}

	if input.RoleAction == nil || *input.RoleAction == "" {
		input.RoleAction = &user.RoleAction
	}

	if input.Status == nil || *input.Status == "" {
		input.Status = &user.Status
	}

	log.Printf("Update data: %+v\n", input)

	err := h.db.Model(&user).Where("user_id = ?", id).Update("updated_at", gorm.Expr("NOW()")).Updates(input).Error
	if err != nil {
		log.Printf("Failed to update user: %v\n", err)
		return c.JSON(fiber.Map{"status": "error", "message": "Failed to update user"})
	}
	return c.JSON(fiber.Map{"status": "success", "message": "User updated successfully"})
}

func (h *ApiHandler) SoftDeleteUser(c *fiber.Ctx) error {
	id := c.Params("id")
	var user models.DeleteUser
	// err := h.db.Model(&user).Where("user_id = ?", id).Updates(user).Error
	err := h.db.Where("user_id = ?", id).Delete(&user).Error
	if err != nil {
		log.Printf("Failed to doft delete user: %v\n", err)
		return c.JSON(fiber.Map{"status": "error", "message": "User not found"})
	}
	return c.JSON(fiber.Map{"status": "success", "message": "User deleted"})
}

func (h *ApiHandler) HardDeleteUser(c *fiber.Ctx) error {
	id := c.Params("id")
	var user models.UserTable
	err := h.db.Unscoped().Where("user_id = ?", id).Delete(&user, id).Error
	if err != nil {
		log.Printf("Failed to hard delete user: %v\n", err)
		return c.JSON(fiber.Map{"status": "error", "message": "User not found"})

	}
	return c.JSON(fiber.Map{"status": "success", "message": "User deleted"})
}
