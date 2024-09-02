package usecases

import (
	"context"
	"crypto/sha1"
	"database/sql"
	"encoding/hex"
	"errors"

	"saranasistemsolusindo.com/gusen-admin/internal/constants"
	"saranasistemsolusindo.com/gusen-admin/internal/models"
	"saranasistemsolusindo.com/gusen-admin/internal/repositories"
	"saranasistemsolusindo.com/gusen-admin/internal/utils/jwt"
)

type UserUseCase struct {
	userRepo      repositories.UserRepository
	userLoginRepo repositories.UserLogRepository
	clientRepo    repositories.ClientRepository
	cityRepo      repositories.CityRepository
}

// NewUserUseCase creates a new UserUseCase
func NewUserUseCase(db *sql.DB) (*UserUseCase, error) {
	userRepo := repositories.NewUserRepository(db)
	userLoginRepo := repositories.NewUserLogRepository(db)
	clientRepo := repositories.NewClientRepository(db)
	cityRepo := repositories.NewCityRepository(db)
	return &UserUseCase{userRepo: userRepo, userLoginRepo: userLoginRepo, clientRepo: clientRepo, cityRepo: cityRepo}, nil
}

// LoginAdmin handles the login for admin users
func (uc *UserUseCase) LoginAdmin(loginID, password string) (string, error) {
	// Fetch user by loginID
	user, err := uc.userRepo.GetUserByLoginID(loginID)
	if err != nil {
		return "", errors.New("user not found")
	}
	// Check if the password matches
	if err := comparePassword(user.Password, password); err != nil {
		return "", errors.New("invalid password")
	}

	if user.OrderRestrictions != constants.IS_ADMIN {
		return "", errors.New("user not Admin")
	}

	token, err := jwt.GenerateJWT(user.LoginID, user.OrderRestrictions)
	if err != nil {
		// Handle error
		return "", err
	}

	return token, nil
}

func (u *UserUseCase) FetchUsers(ctx context.Context, offset, size int, keyword string) ([]models.UserResponse, error) {
	return u.userRepo.FetchUsers(ctx, offset, size, keyword)
}

func (uc *UserUseCase) CreateUser(ctx context.Context, user *models.User) error {
	return uc.userRepo.Create(ctx, user)
}

func (u *UserUseCase) UpdateUser(ctx context.Context, userID string, user *models.User) error {
	// Call the repository method to update the user in the database
	return u.userRepo.Update(ctx, userID, user)
}

func (u *UserUseCase) GetUserByLoginId(ctx context.Context, loginID string) (*models.UserResponseID, error) {
	user, err := u.userRepo.GetUserByLoginID(loginID)
	if err != nil {
		return nil, err
	}

	// Map the user data to the new UserResponse struct
	userResponse := &models.UserResponseID{
		LoginID:           user.LoginID,
		FullName:          user.FullName,
		UserStatus:        user.UserStatus,
		OrderRestrictions: user.OrderRestrictions,
		Password:          &user.Password,
		PasswordExpDate:   &user.PasswordExpDate,
		Pin:               &user.PIN,
		PinExpDate:        &user.PINExpDate,
		MasterClientCD:    &user.MasterClientCD.String,
		LastLogin:         &user.LastLogin.Time,
		CreateBy:          user.CreateBy,
		CreateDT:          user.CreateDT,
		UpdateBy:          user.UpdateBy,
		UpdateDT:          user.UpdateDT,
		IsProtlAm:         &user.IsProtlAm.String,
		Email:             &user.Email.String,
		ClientBirthDT:     &user.ClientBirthDT.Time,
		City:              &user.City.Int16,
		Telepon:           &user.Telepon.String,
	}

	return userResponse, nil
}

func (u *UserUseCase) GetTotalUsers(ctx context.Context) (int, error) {
	return u.userRepo.GetTotalUsers(ctx)
}

func (u *UserUseCase) GetUserLoginPaginated(ctx context.Context, offset, size int, keyword string) ([]models.UserLogin, int, error) {
	user, err := u.userLoginRepo.GetLogHistoryPaginated(ctx, offset, size, keyword)
	if err != nil {
		return nil, 0, err
	}
	total, err := u.userLoginRepo.GetTotalUserLogin(ctx)
	if err != nil {
		return nil, 0, err
	}
	return user, total, nil
}

func comparePassword(hashedPassword, plainPassword string) error {
	// Create a new SHA-1 hash
	hasher := sha1.New()

	// Write the plain password to the hasher
	hasher.Write([]byte(plainPassword))

	// Get the SHA-1 hash in bytes
	sha1Hash := hasher.Sum(nil)

	// Convert the SHA-1 hash to a hexadecimal string
	sha1HashString := hex.EncodeToString(sha1Hash)

	// Compare the hashed password with the stored hash
	if hashedPassword != sha1HashString {
		return errors.New("invalid password")
	}

	return nil
}

func (u *UserUseCase) GetClientListByLoginID(ctx context.Context, loginID string) ([]*models.UserClient, error) {
	return u.clientRepo.GetListClientByLoginID(loginID)
}

func (u *UserUseCase) GetClientNotInUser(ctx context.Context) ([]*models.ClientDetail, error) {
	return u.clientRepo.GetClientNotInUser()
}
func (u *UserUseCase) GetClientDetailByClientCD(ctx context.Context, client string) ([]*models.ClientDetail, error) {
	return u.clientRepo.GetClientByClientID(client)
}

func (u *UserUseCase) DeactiveUser(ctx context.Context, login string) error {
	return u.userRepo.DeactiveUser(ctx, login)
}

func (u *UserUseCase) ResetPass(ctx context.Context, login string, password string) error {
	hasher := sha1.New()
	hasher.Write([]byte(password))
	hashedPassword := hasher.Sum(nil)
	hashedPasswordHex := hex.EncodeToString(hashedPassword)
	return u.userRepo.ChangePin(ctx, login, hashedPasswordHex)
}

func (u *UserUseCase) ResetPin(ctx context.Context, login string, pin string) error {
	hasher := sha1.New()
	hasher.Write([]byte(pin))
	hashedPin := hasher.Sum(nil)
	hashedPinHex := hex.EncodeToString(hashedPin)
	return u.userRepo.ChangePin(ctx, login, hashedPinHex)
}

func (u *UserUseCase) GetCities(ctx context.Context) ([]*models.City, error) {
	return u.cityRepo.GetAllCities(ctx)
}
