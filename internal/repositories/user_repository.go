package repositories

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"reflect"
	"time"

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

	keyword = "%" + keyword + "%"
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
	startRow := offset + 1
	endRow := offset + size
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
	query := `INSERT INTO TLOTSUSER (LOGIN_ID, FULL_NAME, PASSWD, ORDERRESTRICTIONS, PIN, CREATE_BY, CREATE_DT, UPDATE_BY, UPDATE_DT) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`
	_, err := r.db.ExecContext(ctx, query, user.LoginID, user.FullName, user.Password, user.OrderRestrictions, user.PIN, "system", time.Now(), "system", time.Now())
	return err
}

func (r *UserRepositoryImpl) Update(ctx context.Context, userID string, user *models.User) error {
	query := `UPDATE TLOTSUSER SET FULL_NAME = ?, TELEPON = ?, EMAIL = ?, CITY = ?, USER_STATUS = ? WHERE LOGIN_ID = ?`
	params := []interface{}{
		user.FullName,
		user.Telepon,
		user.Email,
		user.City,
		user.UserStatus,
		userID,
	}

	fmt.Println("Query:", query)
	for i, param := range params {
		fmt.Printf("Parameter %d: %v (Type: %v)\n", i+1, param, reflect.TypeOf(param))
	}

	result, err := r.db.ExecContext(ctx, query, params...)
	if err != nil {
		fmt.Println("Error:", err)
		return err
	}

	fmt.Println("Result:", result)
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

	_, err := r.db.ExecContext(ctx, query, userID)
	if err != nil {
		return err
	}

	return nil
}

func (r *UserRepositoryImpl) ChangePassword(ctx context.Context, userID, password string) error {
	query := `UPDATE TLOTSUSER SET PASSWD = :1 WHERE LOGIN_ID = :2`

	_, err := r.db.ExecContext(ctx, query, password, userID)
	if err != nil {
		return err
	}

	return nil
}

func (r *UserRepositoryImpl) ChangePin(ctx context.Context, userID, pin string) error {
	query := `UPDATE TLOTSUSER SET PIN = :1 WHERE LOGIN_ID = :2`

	_, err := r.db.ExecContext(ctx, query, pin, userID)
	if err != nil {
		return err
	}

	return nil
}
