package repositories

import (
	"database/sql"

	"saranasistemsolusindo.com/gusen-admin/internal/models"
)

type ClientRepository interface {
	GetListClientByLoginID(LoginID string) ([]*models.UserClient, error)
	getClientDetail(client *models.UserClient) ([]*models.ClientDetail, error)
	GetClientNotInUser() ([]*models.ClientDetail, error)
	GetClientByClientID(clientID string) ([]*models.ClientDetail, error)
}

// UserRepositoryImpl struct
type ClientRepositoryImpl struct {
	db *sql.DB
}

// NewUserRepository creates a new instance of UserRepositoryImpl
func NewClientRepository(db *sql.DB) ClientRepository {
	return &ClientRepositoryImpl{db: db}
}

func (repo *ClientRepositoryImpl) GetListClientByLoginID(LoginID string) ([]*models.UserClient, error) {
	var client models.UserClient
	query := `SELECT LOGIN_ID, CLIENT_CD, CREATE_DT, CREATE_BY FROM TLOTSUSER_CLIENT WHERE LOGIN_ID = :1`
	err := repo.db.QueryRow(query, LoginID).Scan(
		&client.LoginID, &client.ClientCD, &client.CreateDT, &client.CreateBy,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return []*models.UserClient{}, nil
		}
		return nil, err
	}
	return []*models.UserClient{&client}, nil
}

func (repo *ClientRepositoryImpl) getClientDetail(client *models.UserClient) ([]*models.ClientDetail, error) {
	query := `SELECT CLIENT_CD, CLIENT_NAME FROM TLM_MV_MST_CLIENT WHERE CLIENT_CD = :1`

	rows, err := repo.db.Query(query, client.ClientCD)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var clientDetails []*models.ClientDetail
	for rows.Next() {
		var clientDetail models.ClientDetail
		if err := rows.Scan(&clientDetail.ClientCD, &clientDetail.ClientName); err != nil {
			return nil, err
		}
		clientDetails = append(clientDetails, &clientDetail)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return clientDetails, nil
}

// Get CLient from TLM_MV_MST_CLIENT who not in TLOTSUSER_CLIENT
func (repo *ClientRepositoryImpl) GetClientNotInUser() ([]*models.ClientDetail, error) {
	query := `SELECT CLIENT_CD, CLIENT_NAME FROM TLM_MV_MST_CLIENT WHERE CLIENT_CD NOT IN (SELECT CLIENT_CD FROM TLOTSUSER_CLIENT)`

	rows, err := repo.db.Query(query)
	if err != nil {
		if err == sql.ErrNoRows {
			return []*models.ClientDetail{}, nil
		}
	}
	defer rows.Close()

	var clientDetails []*models.ClientDetail
	for rows.Next() {
		var clientDetail models.ClientDetail
		if err := rows.Scan(&clientDetail.ClientCD, &clientDetail.ClientName); err != nil {
			return nil, err
		}
		clientDetails = append(clientDetails, &clientDetail)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return clientDetails, nil
}

func (repo *ClientRepositoryImpl) GetClientByClientID(clientID string) ([]*models.ClientDetail, error) {
	query := `SELECT CLIENT_CD, CLIENT_NAME FROM TLM_MV_MST_CLIENT WHERE CLIENT_CD LIKE :clientID AND CLIENT_CD NOT IN (SELECT CLIENT_CD FROM TLOTSUSER_CLIENT)`

	// Format clientID for the LIKE operator
	clientID = "%" + clientID + "%"

	rows, err := repo.db.Query(query, sql.Named("clientID", clientID))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var clientDetails []*models.ClientDetail
	for rows.Next() {
		var clientDetail models.ClientDetail
		if err := rows.Scan(&clientDetail.ClientCD, &clientDetail.ClientName); err != nil {
			return nil, err
		}
		clientDetails = append(clientDetails, &clientDetail)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return clientDetails, nil
}
