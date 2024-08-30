package repositories

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"saranasistemsolusindo.com/gusen-admin/internal/models"
)

// UserRepository interface
type UserRepository interface {
	GetUserByLoginID(loginID string) (*models.User, error)
	FetchUsers(ctx context.Context, offset, size int, keyword string) ([]models.UserResponse, error)
	Create(ctx context.Context, user *models.User) error
	Update(ctx context.Context, userID string, user *models.User) error
	GetTotalUsers(ctx context.Context) (int, error)
	ChangePin(ctx context.Context, userID, pin string) error
	ChangePassword(ctx context.Context, userID, password string) error
	DeactiveUser(ctx context.Context, userID string) error
}

// UserRepositoryImpl struct
type UserRepositoryImpl struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) UserRepository {
	return &UserRepositoryImpl{db: db}
}

// GetUserByLoginID fetches a user by their login ID
func (repo *UserRepositoryImpl) GetUserByLoginID(loginID string) (*models.User, error) {
	user := &models.User{}

	query := `SELECT LOGIN_ID, FULL_NAME, PASSWD, PASSWD_EXPDATE, USER_STATUS, ORDERRESTRICTIONS, PIN, PIN_EXPDATE, MASTER_CLIENT_CD, LAST_LOGIN, CREATE_BY, CREATE_DT, UPDATE_BY, UPDATE_DT, PHOTO_ID, IS_PROTLAM, EMAIL, CLIENT_BIRTH_DT, CITY, TELEPON FROM TLOTSUSER WHERE LOGIN_ID = :1`
	err := repo.db.QueryRow(query, loginID).Scan(
		&user.LoginID, &user.FullName, &user.Password, &user.PasswordExpDate, &user.UserStatus, &user.OrderRestrictions, &user.PIN, &user.PINExpDate, &user.MasterClientCD, &user.LastLogin, &user.CreateBy, &user.CreateDT, &user.UpdateBy, &user.UpdateDT, &user.PhotoID, &user.IsProtlAm, &user.Email, &user.ClientBirthDT, &user.City, &user.Telepon,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	return user, nil
}

func (r *UserRepositoryImpl) FetchUsers(ctx context.Context, offset, size int, keyword string) ([]models.UserResponse, error) {
	var users []models.UserResponse

	// Prepare the keyword for SQL LIKE clause
	keyword = "%" + keyword + "%"

	// Query to fetch users with pagination and keyword filtering
	query := `
        SELECT LOGIN_ID, FULL_NAME, USER_STATUS, ORDERRESTRICTIONS, CREATE_DT, UPDATE_DT
        FROM (
            SELECT LOGIN_ID, FULL_NAME, USER_STATUS, ORDERRESTRICTIONS, CREATE_DT, UPDATE_DT,
                   ROW_NUMBER() OVER (ORDER BY CREATE_DT) AS rnum
            FROM TLOTSUSER
            WHERE FULL_NAME LIKE :1 OR EMAIL LIKE :1
        )
        WHERE rnum BETWEEN :2 AND :3
    `
	// Calculate the row numbers for pagination
	startRow := offset + 1
	endRow := offset + size

	// Print the query and parameters
	rows, err := r.db.QueryContext(ctx, query, keyword, keyword, startRow, endRow)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var user models.UserResponse
		if err := rows.Scan(
			&user.LoginID, &user.FullName, &user.UserStatus, &user.OrderRestrictions, &user.CreateDT, &user.UpdateDT,
		); err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}

func (r *UserRepositoryImpl) Create(ctx context.Context, user *models.User) error {
	query := `INSERT INTO TLOTSUSER (LOGIN_ID, FULL_NAME, PASSWD, PASSWD_EXPDATE, USER_STATUS, DESC_STATUS, ORDERRESTRICTIONS, PIN, PIN_EXPDATE, MASTER_CLIENT_CD, LAST_LOGIN, CREATE_BY, CREATE_DT, UPDATE_BY, UPDATE_DT, PHOTO_ID, IS_PROTLAM, EMAIL, CLIENT_BIRTH_DT, CITY, TELEPON) VALUES (:1, :2, :3, :4, :5, :6, :7, :8, :9, :10, :11, :12, :13, :14, :15, :16, :17, :18, :19, :20, :21)`
	_, err := r.db.ExecContext(ctx, query, user.LoginID, user.FullName, user.Password, user.PasswordExpDate, user.UserStatus, user.DescStatus, user.OrderRestrictions, user.PIN, user.PINExpDate, user.MasterClientCD, user.LastLogin, user.CreateBy, user.CreateDT, user.UpdateBy, user.UpdateDT, user.PhotoID, user.IsProtlAm, user.Email, user.ClientBirthDT, user.City, user.Telepon)
	return err
}

func (r *UserRepositoryImpl) Update(ctx context.Context, userID string, user *models.User) error {
	query := `
        UPDATE TLOTSUSER
        SET login_id = ?, full_name = ?, email = ?, updated_at = NOW()
        WHERE id = ?
    `

	_, err := r.db.ExecContext(ctx, query, user.LoginID, user.FullName, user.Email, userID)
	if err != nil {
		return err
	}

	return nil
}

func (r *UserRepositoryImpl) GetTotalUsers(ctx context.Context) (int, error) {
	var total int
	query := `SELECT COUNT(*) FROM TLOTSUSER`
	err := r.db.QueryRowContext(ctx, query).Scan(&total)
	if err != nil {
		return 0, err
	}
	return total, nil
}

func (r *UserRepositoryImpl) DeactiveUser(ctx context.Context, userID string) error {
	query := `UPDATE TLOTSUSER SET USER_STATUS = 'S' WHERE LOGIN_ID = :1`

	// Execute the query
	q, err := r.db.ExecContext(ctx, query, userID)
	fmt.Println(err)
	if err != nil {
		return err
	}
	fmt.Println(q)

	return nil
}

func (r *UserRepositoryImpl) ChangePassword(ctx context.Context, userID, password string) error {
	query := `UPDATE TLOTSUSER SET PASSWD = :1 WHERE LOGIN_ID = :2`

	// Execute the query
	_, err := r.db.ExecContext(ctx, query, password, userID)
	if err != nil {
		return err
	}

	return nil
}

func (r *UserRepositoryImpl) ChangePin(ctx context.Context, userID, pin string) error {
	query := `UPDATE TLOTSUSER SET PIN = :1 WHERE LOGIN_ID = :2`

	// Execute the query
	_, err := r.db.ExecContext(ctx, query, pin, userID)
	if err != nil {
		return err
	}

	return nil
}
