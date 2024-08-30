package repositories

import (
	"context"
	"database/sql"

	"saranasistemsolusindo.com/gusen-admin/internal/models"
)

type UserLogRepository interface {
	GetLogHistoryPaginated(ctx context.Context, offset, size int, keyword string) ([]models.UserLogin, error)
	GetTotalUserLogin(ctx context.Context) (int, error)
}

type UserLogRepositoryImpl struct {
	db *sql.DB
}

func NewUserLogRepository(db *sql.DB) UserLogRepository {
	return &UserLogRepositoryImpl{db: db}
}

func (r *UserLogRepositoryImpl) GetLogHistoryPaginated(ctx context.Context, offset, size int, keyword string) ([]models.UserLogin, error) {
	var userLog []models.UserLogin

	keyword = "%" + keyword + "%"
	query := `
    SELECT LOGINID, STATUS, ACTION_DATE, ACTION_TIME, CHANNEL_MEDIA, CHANNEL_DEVICE, IP_ADDRESS
    FROM (
        SELECT LOGINID, STATUS, ACTION_DATE, ACTION_TIME, CHANNEL_MEDIA, CHANNEL_DEVICE, IP_ADDRESS,
            ROW_NUMBER() OVER (ORDER BY ACTION_DATE DESC) AS r
        FROM TLD_USERLOGIN_LOG
        WHERE LOGINID LIKE :1
    )
    WHERE r > :2 AND r <= :3
`
	rows, err := r.db.QueryContext(ctx, query, keyword, offset, offset+size)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var user models.UserLogin
		err := rows.Scan(&user.LoginID, &user.Status, &user.ActionDate, &user.ActionTime, &user.ChannelMedia, &user.ChannelDevice, &user.IPAddress)
		if err != nil {
			return nil, err
		}
		userLog = append(userLog, user)
	}

	// Check for errors from iterating over rows
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return userLog, nil
}

func (r *UserLogRepositoryImpl) GetTotalUserLogin(ctx context.Context) (int, error) {
	query := `SELECT COUNT(*) FROM TLD_USERLOGIN_LOG`
	var total int
	err := r.db.QueryRowContext(ctx, query).Scan(&total)
	if err != nil {
		return 0, err
	}
	return total, nil
}
