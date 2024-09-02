package handlers

import (
	"context"
	"database/sql"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"saranasistemsolusindo.com/gusen-admin/internal/handlers/responses"
	"saranasistemsolusindo.com/gusen-admin/internal/usecases"
)

type LogHandler struct {
	userUseCase usecases.UserUseCase
}

func NewLogHandler(db *sql.DB) (*LogHandler, error) {
	userUseCase, err := usecases.NewUserUseCase(db)
	if err != nil {
		return nil, err
	}
	return &LogHandler{userUseCase: *userUseCase}, nil
}

func (h *LogHandler) GetLogHistory(c echo.Context) error {
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
