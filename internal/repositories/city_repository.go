package repositories

import (
	"context"
	"database/sql"

	"saranasistemsolusindo.com/gusen-admin/internal/models"
)

type CityRepository interface {
	GetAllCities(ctx context.Context) ([]*models.City, error)
}

type CityRepositoryImpl struct {
	db *sql.DB
}

// NewUserRepository creates a new instance of UserRepositoryImpl
func NewCityRepository(db *sql.DB) CityRepository {
	return &CityRepositoryImpl{db: db}
}

func (repo *CityRepositoryImpl) GetAllCities(ctx context.Context) ([]*models.City, error) {
	var cities []*models.City

	query := `SELECT CITY_CD, CITY_NAME FROM TLD_MST_CITY`
	rows, err := repo.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var city models.City
		if err := rows.Scan(&city.CityCD, &city.CityName); err != nil {
			return nil, err
		}
		cities = append(cities, &city)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return cities, nil
}
