package handlers

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	pb "github.com/quanbin27/commons/genproto/users"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type UserHandler struct {
	client pb.UserServiceClient
}

func NewUserHandler(client pb.UserServiceClient) *UserHandler {
	return &UserHandler{client: client}
}

// RegisterRoutes đăng ký các route cho Users service với tiền tố "/users"
func (h *UserHandler) RegisterRoutes(e *echo.Echo) {
	e.POST("/users/register", h.RegisterUser)
	e.POST("/users/login", h.LoginUser)
	e.PUT("/users/change-info", h.ChangeInfo)
	e.PUT("/users/change-password", h.ChangePassword)
	e.GET("/users/info/:id", h.GetUserInfo)
	e.GET("/users/info-by-email", h.GetUserInfoByEmail)
}

// RegisterUser xử lý yêu cầu đăng ký người dùng
func (h *UserHandler) RegisterUser(c echo.Context) error {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
		Name     string `json:"name"`
	}
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
	}

	// Lấy context từ Echo request
	ctx := c.Request().Context()

	resp, err := h.client.Register(ctx, &pb.RegisterRequest{
		Email:    req.Email,
		Password: req.Password,
		Name:     req.Name,
	})
	if err != nil {
		if grpcErr, ok := status.FromError(err); ok {
			switch grpcErr.Code() {
			case codes.AlreadyExists:
				return c.JSON(http.StatusConflict, map[string]string{"error": grpcErr.Message()})
			case codes.Internal:
				return c.JSON(http.StatusInternalServerError, map[string]string{"error": grpcErr.Message()})
			}
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]string{
		"status": resp.Status,
	})
}

// LoginUser xử lý yêu cầu đăng nhập người dùng
func (h *UserHandler) LoginUser(c echo.Context) error {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
	}

	// Lấy context từ Echo request
	ctx := c.Request().Context()

	resp, err := h.client.Login(ctx, &pb.LoginRequest{
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		if grpcErr, ok := status.FromError(err); ok {
			switch grpcErr.Code() {
			case codes.Unauthenticated:
				return c.JSON(http.StatusUnauthorized, map[string]string{"error": grpcErr.Message()})
			case codes.Internal:
				return c.JSON(http.StatusInternalServerError, map[string]string{"error": grpcErr.Message()})
			}
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"status": resp.Status,
		"token":  resp.Token,
	})
}

// ChangeInfo xử lý yêu cầu thay đổi thông tin người dùng
func (h *UserHandler) ChangeInfo(c echo.Context) error {
	var req struct {
		ID    string `json:"id"`
		Email string `json:"email"`
		Name  string `json:"name"`
	}
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
	}

	// Chuyển đổi ID từ string sang int32
	id, err := strconv.ParseInt(req.ID, 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid ID format, must be an integer"})
	}

	// Lấy context từ Echo request
	ctx := c.Request().Context()

	resp, err := h.client.ChangeInfo(ctx, &pb.ChangeInfoRequest{
		Id:    int32(id),
		Email: req.Email,
		Name:  req.Name,
	})
	if err != nil {
		if grpcErr, ok := status.FromError(err); ok {
			switch grpcErr.Code() {
			case codes.InvalidArgument:
				return c.JSON(http.StatusBadRequest, map[string]string{"error": grpcErr.Message()})
			case codes.Internal:
				return c.JSON(http.StatusInternalServerError, map[string]string{"error": grpcErr.Message()})
			}
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"status": resp.Status,
		"email":  resp.Email,
		"name":   resp.Name,
	})
}

// ChangePassword xử lý yêu cầu thay đổi mật khẩu
func (h *UserHandler) ChangePassword(c echo.Context) error {
	var req struct {
		ID          string `json:"id"`
		OldPassword string `json:"old_password"`
		NewPassword string `json:"new_password"`
	}
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
	}

	// Chuyển đổi ID từ string sang int32
	id, err := strconv.ParseInt(req.ID, 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid ID format, must be an integer"})
	}

	// Lấy context từ Echo request
	ctx := c.Request().Context()

	resp, err := h.client.ChangePassword(ctx, &pb.ChangePasswordRequest{
		Id:          int32(id),
		OldPassword: req.OldPassword,
		NewPassword: req.NewPassword,
	})
	if err != nil {
		if grpcErr, ok := status.FromError(err); ok {
			switch grpcErr.Code() {
			case codes.InvalidArgument:
				return c.JSON(http.StatusBadRequest, map[string]string{"error": grpcErr.Message()})
			case codes.Internal:
				return c.JSON(http.StatusInternalServerError, map[string]string{"error": grpcErr.Message()})
			}
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]string{
		"status": resp.Status,
	})
}

// GetUserInfo xử lý yêu cầu lấy thông tin người dùng theo ID
func (h *UserHandler) GetUserInfo(c echo.Context) error {
	idStr := c.Param("id")
	if idStr == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "ID is required"})
	}

	// Chuyển đổi ID từ string sang int32
	id, err := strconv.ParseInt(idStr, 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid ID format, must be an integer"})
	}

	// Lấy context từ Echo request
	ctx := c.Request().Context()

	resp, err := h.client.GetUserInfo(ctx, &pb.GetUserInfoRequest{ID: int32(id)})
	if err != nil {
		if grpcErr, ok := status.FromError(err); ok {
			switch grpcErr.Code() {
			case codes.NotFound:
				return c.JSON(http.StatusNotFound, map[string]string{"error": grpcErr.Message()})
			case codes.Internal:
				return c.JSON(http.StatusInternalServerError, map[string]string{"error": grpcErr.Message()})
			}
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, resp)
}

// GetUserInfoByEmail xử lý yêu cầu lấy thông tin người dùng theo email
func (h *UserHandler) GetUserInfoByEmail(c echo.Context) error {
	email := c.QueryParam("email")
	if email == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Email is required"})
	}

	// Lấy context từ Echo request
	ctx := c.Request().Context()

	resp, err := h.client.GetUserInfoByEmail(ctx, &pb.GetUserInfoByEmailRequest{Email: email})
	if err != nil {
		if grpcErr, ok := status.FromError(err); ok {
			switch grpcErr.Code() {
			case codes.NotFound:
				return c.JSON(http.StatusNotFound, map[string]string{"error": grpcErr.Message()})
			case codes.Internal:
				return c.JSON(http.StatusInternalServerError, map[string]string{"error": grpcErr.Message()})
			}
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, resp)
}
