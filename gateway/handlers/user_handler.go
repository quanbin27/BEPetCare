package handlers

import (
	"github.com/quanbin27/commons/auth"
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
func (h *UserHandler) RegisterRoutes(e *echo.Group) {
	e.POST("/users/register", h.Register)
	e.GET("/users/verify", h.VerifyEmail)
	e.POST("/users/login", h.LoginUser)
	e.PUT("/users/change-info", h.ChangeInfo, auth.WithJWTAuth())
	e.PUT("/users/change-password", h.ChangePassword, auth.WithJWTAuth())
	e.GET("/users/info/:id", h.GetUserInfo, auth.RoleMiddleware("1", "2"))
	e.GET("/users/info/me", h.GetMyInfo, auth.WithJWTAuth())
	e.GET("/users/info-by-email", h.GetUserInfoByEmail)
	e.POST("/users/forgot-password", h.ForgotPassword)
	e.POST("/users/reset-password", h.ResetPassword)
	e.GET("/helloworld", helloWorld)
}
func helloWorld(c echo.Context) error {
	return c.String(http.StatusOK, "Hello World")
}
func (h *UserHandler) Register(c echo.Context) error {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
		Name     string `json:"name"`
	}
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
	}

	if req.Email == "" || req.Password == "" || req.Name == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Email, password, and name are required"})
	}

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
		"message": resp.Status, // "Verification email sent"
	})
}

// VerifyEmail xử lý xác minh email
func (h *UserHandler) VerifyEmail(c echo.Context) error {
	token := c.QueryParam("token")
	if token == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Token is required"})
	}

	ctx := c.Request().Context()
	resp, err := h.client.VerifyEmail(ctx, &pb.VerifyEmailRequest{Token: token})
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

	return c.JSON(http.StatusOK, map[string]int32{
		"id": resp.Id,
	})
}

// LoginUser xử lý yêu cầu đăng nhập người dùng
func (h *UserHandler) LoginUser(c echo.Context) error {
	var req struct {
		Email      string `json:"email"`
		Password   string `json:"password"`
		RememberMe bool   `json:"rememberMe"`
	}
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
	}

	// Lấy context từ Echo request
	ctx := c.Request().Context()

	resp, err := h.client.Login(ctx, &pb.LoginRequest{
		Email:      req.Email,
		Password:   req.Password,
		RememberMe: req.RememberMe,
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
	id, err := auth.GetUserIDFromContext(c)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	var req struct {
		Email       string `json:"email"`
		Name        string `json:"name"`
		PhoneNumber string `json:"phoneNumber"`
		Address     string `json:"address"`
	}
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
	}

	// Lấy context từ Echo request
	ctx := c.Request().Context()

	resp, err := h.client.ChangeInfo(ctx, &pb.ChangeInfoRequest{
		Id:          id,
		Email:       req.Email,
		Name:        req.Name,
		PhoneNumber: req.PhoneNumber,
		Address:     req.Address,
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
		"status":      resp.Status,
		"email":       resp.Email,
		"name":        resp.Name,
		"address":     resp.Address,
		"phoneNumber": resp.PhoneNumber,
	})
}

// ChangePassword xử lý yêu cầu thay đổi mật khẩu
func (h *UserHandler) ChangePassword(c echo.Context) error {
	id, err := auth.GetUserIDFromContext(c)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	var req struct {
		OldPassword string `json:"oldPassword"`
		NewPassword string `json:"newPassword"`
	}
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
	}
	ctx := c.Request().Context()
	resp, err := h.client.ChangePassword(ctx, &pb.ChangePasswordRequest{
		Id:          id,
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
func (h *UserHandler) GetMyInfo(c echo.Context) error {

	id, err := auth.GetUserIDFromContext(c)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	// Lấy context từ Echo request
	ctx := c.Request().Context()

	resp, err := h.client.GetUserInfo(ctx, &pb.GetUserInfoRequest{ID: id})
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
	response := UserResponse{
		UserID:      resp.ID,
		Name:        resp.Name,
		Email:       resp.Email,
		PhoneNumber: resp.PhoneNumber,
		Address:     resp.Address,
	}

	return c.JSON(http.StatusOK, response)
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
func (h *UserHandler) ForgotPassword(c echo.Context) error {
	var req struct {
		Email string `json:"email"`
	}
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
	}

	ctx := c.Request().Context()
	_, err := h.client.ForgotPassword(ctx, &pb.ForgotPasswordRequest{Email: req.Email})
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

	return c.JSON(http.StatusOK, map[string]string{"message": "Password reset email sent"})
}
func (h *UserHandler) ResetPassword(c echo.Context) error {
	var req struct {
		UserID      string `json:"userId"`
		Token       string `json:"token"`
		NewPassword string `json:"newPassword"`
	}
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
	}
	id, err := strconv.ParseInt(req.UserID, 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid ID format, must be an integer"})
	}
	if id == 0 || req.Token == "" || req.NewPassword == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Missing required fields"})
	}

	ctx := c.Request().Context()
	_, err = h.client.ResetPassword(ctx, &pb.ResetPasswordRequest{
		UserID:      int32(id),
		Token:       req.Token,
		NewPassword: req.NewPassword,
	})
	if err != nil {
		if grpcErr, ok := status.FromError(err); ok {
			switch grpcErr.Code() {
			case codes.InvalidArgument:
				return c.JSON(http.StatusBadRequest, map[string]string{"error": grpcErr.Message()})
			case codes.Unauthenticated:
				return c.JSON(http.StatusUnauthorized, map[string]string{"error": grpcErr.Message()})
			case codes.Internal:
				return c.JSON(http.StatusInternalServerError, map[string]string{"error": grpcErr.Message()})
			}
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Password reset successfully"})
}
