package repository

import (
	"database/sql"
	"fmt"
	"strings"
	"user-api/internal/models"

	"github.com/jmoiron/sqlx"
)

// UserRepository интерфейс для работы с пользователями
type UserRepository interface {
	Create(user *models.CreateUserRequest) (*models.User, error)
	GetByID(id int) (*models.User, error)
	GetAll(page, pageSize int, filters map[string]interface{}) ([]models.User, int, error)
	Update(id int, user *models.UpdateUserRequest) (*models.User, error)
	Delete(id int) error
}

type userRepository struct {
	db *sqlx.DB
}

// NewUserRepository создает новый репозиторий пользователей
func NewUserRepository(db *sqlx.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(req *models.CreateUserRequest) (*models.User, error) {
	query := `
        INSERT INTO users (name, email, age)
        VALUES ($1, $2, $3)
        RETURNING id, name, email, age, created_at, updated_at
    `

	var user models.User
	err := r.db.QueryRowx(query, req.Name, req.Email, req.Age).StructScan(&user)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return &user, nil
}

func (r *userRepository) GetByID(id int) (*models.User, error) {
	query := `
        SELECT id, name, email, age, created_at, updated_at
        FROM users
        WHERE id = $1
    `

	var user models.User
	err := r.db.Get(&user, query, id)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("user not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return &user, nil
}

func (r *userRepository) GetAll(page, pageSize int, filters map[string]interface{}) ([]models.User, int, error) {
	// Построение WHERE clause для фильтрации
	var conditions []string
	var args []interface{}
	argCounter := 1

	if name, ok := filters["name"].(string); ok && name != "" {
		conditions = append(conditions, fmt.Sprintf("name ILIKE $%d", argCounter))
		args = append(args, "%"+name+"%")
		argCounter++
	}

	if email, ok := filters["email"].(string); ok && email != "" {
		conditions = append(conditions, fmt.Sprintf("email ILIKE $%d", argCounter))
		args = append(args, "%"+email+"%")
		argCounter++
	}

	if minAge, ok := filters["min_age"].(int); ok && minAge > 0 {
		conditions = append(conditions, fmt.Sprintf("age >= $%d", argCounter))
		args = append(args, minAge)
		argCounter++
	}

	if maxAge, ok := filters["max_age"].(int); ok && maxAge > 0 {
		conditions = append(conditions, fmt.Sprintf("age <= $%d", argCounter))
		args = append(args, maxAge)
		argCounter++
	}

	whereClause := ""
	if len(conditions) > 0 {
		whereClause = "WHERE " + strings.Join(conditions, " AND ")
	}

	// Подсчет общего количества
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM users %s", whereClause)
	var total int
	err := r.db.Get(&total, countQuery, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count users: %w", err)
	}

	// Получение данных с пагинацией
	offset := (page - 1) * pageSize
	query := fmt.Sprintf(`
        SELECT id, name, email, age, created_at, updated_at
        FROM users
        %s
        ORDER BY created_at DESC
        LIMIT $%d OFFSET $%d
    `, whereClause, argCounter, argCounter+1)

	args = append(args, pageSize, offset)

	var users []models.User
	err = r.db.Select(&users, query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get users: %w", err)
	}

	return users, total, nil
}

func (r *userRepository) Update(id int, req *models.UpdateUserRequest) (*models.User, error) {
	var updates []string
	var args []interface{}
	argCounter := 1

	if req.Name != "" {
		updates = append(updates, fmt.Sprintf("name = $%d", argCounter))
		args = append(args, req.Name)
		argCounter++
	}

	if req.Email != "" {
		updates = append(updates, fmt.Sprintf("email = $%d", argCounter))
		args = append(args, req.Email)
		argCounter++
	}

	if req.Age > 0 {
		updates = append(updates, fmt.Sprintf("age = $%d", argCounter))
		args = append(args, req.Age)
		argCounter++
	}

	if len(updates) == 0 {
		return r.GetByID(id)
	}

	updates = append(updates, "updated_at = CURRENT_TIMESTAMP")
	args = append(args, id)

	query := fmt.Sprintf(`
        UPDATE users
        SET %s
        WHERE id = $%d
        RETURNING id, name, email, age, created_at, updated_at
    `, strings.Join(updates, ", "), argCounter)

	var user models.User
	err := r.db.QueryRowx(query, args...).StructScan(&user)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("user not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to update user: %w", err)
	}

	return &user, nil
}

func (r *userRepository) Delete(id int) error {
	query := "DELETE FROM users WHERE id = $1"
	result, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("user not found")
	}

	return nil
}
