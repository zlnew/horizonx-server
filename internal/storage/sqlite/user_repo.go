package sqlite

import (
	"context"
	"database/sql"
	"fmt"

	"horizonx-server/internal/domain"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) domain.UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) GetUsers(ctx context.Context, opts domain.ListOptions) ([]*domain.User, int64, error) {
	baseQuery := "FROM users"
	whereQuery := ""
	args := []any{}

	if opts.Search != "" {
		whereQuery = " WHERE email LIKE ?"
		args = append(args, "%"+opts.Search+"%")
	}

	var total int64

	if opts.IsPaginate {
		countQuery := "SELECT COUNT(*) " + baseQuery + whereQuery
		if err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&total); err != nil {
			return nil, 0, fmt.Errorf("failed to count users: %w", err)
		}
	}

	selectQuery := "SELECT id, email, password FROM users" + whereQuery

	if opts.IsPaginate {
		offset := (opts.Page - 1) * opts.Limit
		selectQuery += " LIMIT ? OFFSET ?"
		args = append(args, opts.Limit, offset)
	} else {
		selectQuery += " LIMIT 1000"
	}

	rows, err := r.db.QueryContext(ctx, selectQuery, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to query users: %w", err)
	}
	defer rows.Close()

	var users []*domain.User
	for rows.Next() {
		var user domain.User
		if err := rows.Scan(&user.ID, &user.Email, &user.Password); err != nil {
			return nil, 0, err
		}
		users = append(users, &user)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

func (r *UserRepository) GetUserByID(ctx context.Context, ID int64) (*domain.User, error) {
	query := `SELECT id, email, password FROM users WHERE id = ?`

	row := r.db.QueryRowContext(ctx, query, ID)

	var user domain.User
	if err := row.Scan(&user.ID, &user.Email, &user.Password); err != nil {
		if err == sql.ErrNoRows {
			return nil, domain.ErrUserNotFound
		}
		return nil, err
	}

	return &user, nil
}

func (r *UserRepository) GetUserByEmail(ctx context.Context, email string) (*domain.User, error) {
	query := `SELECT id, email, password FROM users WHERE email = ?`

	row := r.db.QueryRowContext(ctx, query, email)

	var user domain.User
	if err := row.Scan(&user.ID, &user.Email, &user.Password); err != nil {
		if err == sql.ErrNoRows {
			return nil, domain.ErrUserNotFound
		}
		return nil, err
	}

	return &user, nil
}

func (r *UserRepository) CreateUser(ctx context.Context, user *domain.User) error {
	query := `INSERT INTO users (email, password) VALUES (?, ?)`

	result, err := r.db.ExecContext(ctx, query, user.Email, user.Password)
	if err != nil {
		return fmt.Errorf("failed to insert user: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get last insert id: %w", err)
	}

	user.ID = id

	return nil
}

func (r *UserRepository) UpdateUser(ctx context.Context, user *domain.User, userID int64) error {
	query := `UPDATE users SET email = ?, password = ? WHERE id = ?`

	result, err := r.db.ExecContext(ctx, query, user.Email, user.Password, userID)
	if err != nil {
		return fmt.Errorf("failed to execute update query: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to retrieve affected rows: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("user with ID %d not found", userID)
	}

	return nil
}

func (r *UserRepository) DeleteUser(ctx context.Context, userID int64) error {
	query := `DELETE FROM users WHERE id = ?`

	result, err := r.db.ExecContext(ctx, query, userID)
	if err != nil {
		return fmt.Errorf("failed to execute delete query: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to retrieve affected rows: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("user with ID %d not found", userID)
	}

	return nil
}
