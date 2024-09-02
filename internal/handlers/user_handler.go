package handlers

import (
	"context"
	"crypto/sha1"
	"database/sql"
	"encoding/hex"
	"net/http"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
	"saranasistemsolusindo.com/gusen-admin/internal/handlers/requests"
	"saranasistemsolusindo.com/gusen-admin/internal/handlers/responses"
	"saranasistemsolusindo.com/gusen-admin/internal/usecases"
)

// UserHandler struct
type UserHandler struct {
	userUseCase usecases.UserUseCase
}

// NewUserHandler creates a new UserHandler
func NewUserHandler(db *sql.DB) (*UserHandler, error) {
	userUseCase, err := usecases.NewUserUseCase(db)
	if err != nil {
		return nil, err
	}
	return &UserHandler{userUseCase: *userUseCase}, nil
}

type UserResponse struct {
	LoginID           string `json:"LoginID"`
	FullName          string `json:"FullName"`
	UserStatus        string `json:"UserStatus"`
	OrderRestrictions string `json:"OrderRestrictions"`
	CreateDT          string `json:"CreateDT"`
	UpdateDT          string `json:"UpdateDT"`
}

func (h *UserHandler) LoginAdmin(c echo.Context) error {
	var loginReq requests.LoginRequest
	if err := c.Bind(&loginReq); err != nil {
		return c.JSON(http.StatusBadRequest, responses.BaseResponse{
			StatusCode: http.StatusBadRequest,
			Message:    "Invalid request payload",
		})
	}

	id := strings.ToUpper(loginReq.LoginID)

	token, err := h.userUseCase.LoginAdmin(id, loginReq.Password)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, responses.BaseResponse{
			StatusCode: http.StatusUnauthorized,
			Message:    err.Error(),
		})
	}

	return c.JSON(http.StatusOK, responses.BaseResponse{
		StatusCode: http.StatusOK,
		Message:    "Login successful",
		Data:       map[string]string{"token": token},
	})
}

func (h *UserHandler) GetUserPaginated(c echo.Context) error {
	page, err := strconv.Atoi(c.QueryParam("page"))
	if err != nil || page < 1 {
		page = 1
	}

	size, err := strconv.Atoi(c.QueryParam("size"))
	if err != nil || size < 1 {
		size = 10
	}

	keyword := c.QueryParam("keyword")
	if keyword == "" {
		keyword = ""
	}

	offset := (page - 1) * size

	users, err := h.userUseCase.FetchUsers(context.Background(), offset, size, keyword)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, responses.BaseResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to fetch users",
		})
	}

	total, err := h.userUseCase.GetTotalUsers(context.Background())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, responses.BaseResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to fetch total user",
		})
	}

	response := map[string]interface{}{
		"users": users,
		"page":  page,
		"size":  size,
		"total": total,
	}

	return c.JSON(http.StatusOK, responses.BaseResponse{
		StatusCode: http.StatusOK,
		Message:    "Fetch User Successful",
		Data:       response,
	})
}

func hashString(s string) string {
	h := sha1.New()
	h.Write([]byte(s))
	return hex.EncodeToString(h.Sum(nil))
}

func (h *UserHandler) GetUserByLoginId(c echo.Context) error {
	loginID := c.Param("login_id")

	user, err := h.userUseCase.GetUserByLoginId(c.Request().Context(), loginID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to fetch user"})
	}

	return c.JSON(http.StatusOK, responses.BaseResponse{
		StatusCode: http.StatusOK,
		Message:    "Fetch User Successful",
		Data:       user,
	})
}

// GetLogHistory handles fetching log history
func (h *UserHandler) GetLogHistory(c echo.Context) error {
	page, err := strconv.Atoi(c.QueryParam("page"))
	if err != nil || page < 1 {
		page = 1
	}

	size, err := strconv.Atoi(c.QueryParam("size"))
	if err != nil || size < 1 {
		size = 10
	}

	keyword := c.QueryParam("keyword")
	if keyword == "" {
		keyword = ""
	}

	offset := (page - 1) * size

	userLogin, total, err := h.userUseCase.GetUserLoginPaginated(context.Background(), offset, size, keyword)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, responses.BaseResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to fetch log history",
		})
	}

	response := map[string]interface{}{
		"Log":   userLogin,
		"page":  page,
		"size":  size,
		"total": total,
	}

	return c.JSON(http.StatusOK, responses.BaseResponse{
		StatusCode: http.StatusOK,
		Message:    "Fetch User Successful",
		Data:       response,
	})
}

func (u *UserHandler) GetClientByLoginID(c echo.Context) error {
	loginID := c.Param("login_id")
	client, err := u.userUseCase.GetClientListByLoginID(c.Request().Context(), loginID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to fetch user"})
	}

	return c.JSON(http.StatusOK, responses.BaseResponse{
		StatusCode: http.StatusOK,
		Message:    "Fetch User Successful",
		Data:       client,
	})
}

func (u *UserHandler) GetAvailableClients(c echo.Context) error {
	clients, err := u.userUseCase.GetClientNotInUser(c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to fetch clients"})
	}

	return c.JSON(http.StatusOK, responses.BaseResponse{
		StatusCode: http.StatusOK,
		Message:    "Fetch Clients Successful",
		Data:       clients,
	})
}

func (u *UserHandler) GetClientByClientID(c echo.Context) error {
	client := c.QueryParam("client_id")
	if client == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Client ID is required"})
	}

	clients, err := u.userUseCase.GetClientDetailByClientCD(c.Request().Context(), client)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to fetch clients"})
	}
	return c.JSON(http.StatusOK, responses.BaseResponse{
		StatusCode: http.StatusOK,
		Message:    "Fetch Clients Successful",
		Data:       clients,
	})
}

func (u *UserHandler) UpdateClientByUserLogin(c echo.Context) error {
	return c.JSON(http.StatusOK, "GetClientByClientName")
}

func (u *UserHandler) DeactiveUser(c echo.Context) error {
	loginID := c.Param("login_id")
	err := u.userUseCase.DeactiveUser(c.Request().Context(), loginID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to deactivate user"})
	}
	return c.JSON(http.StatusOK, map[string]string{"message": "User deactivated successfully"})
}

func (u *UserHandler) ResetPin(c echo.Context) error {
	var req requests.ResetPin

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
	}

	if len(req.Pin) != 6 {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Password must be longer than 6 characters"})
	}

	err := u.userUseCase.ResetPin(c.Request().Context(), req.LoginID, req.Pin)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to reset pin"})
	}
	return c.JSON(http.StatusOK, map[string]string{"message": "Pin reset successfully"})
}

func (u *UserHandler) ResetPassword(c echo.Context) error {
	var req requests.ResetPassword

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
	}

	err := u.userUseCase.ResetPass(c.Request().Context(), req.LoginID, req.Password)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to reset password"})
	}
	return c.JSON(http.StatusOK, map[string]string{"message": "Password reset successfully"})
}

func (u *UserHandler) GetCities(c echo.Context) error {
	cities, err := u.userUseCase.GetCities(c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to fetch cities"})
	}
	return c.JSON(http.StatusOK, responses.BaseResponse{
		StatusCode: http.StatusOK,
		Message:    "Fetch Cities Successful",
		Data:       cities,
	})
}
