package handlers

import (
	"github.com/labstack/echo/v4"
	"github.com/quanbin27/commons/auth"
	pb "github.com/quanbin27/commons/genproto/users"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"net/http"
	"strconv"
)

type UserHandler struct {
	client pb.UserServiceClient
}

func NewUserHandler(client pb.UserServiceClient) *UserHandler {
	return &UserHandler{client: client}
}

// RegisterRoutes registers routes for the Users service with the "/users" prefix
func (h *UserHandler) RegisterRoutes(e *echo.Group) {
	e.POST("/auth/register", h.Register)
	e.POST("/auth/verify", h.VerifyEmail)
	e.POST("/auth/login", h.LoginUser)
	e.POST("/auth/forgot-password", h.ForgotPassword)
	e.POST("/auth/reset-password", h.ResetPassword)

	// User Management Routes
	e.GET("/users/me", h.GetMyInfo, auth.WithJWTAuth())
	e.PUT("/users/me", h.ChangeInfo, auth.WithJWTAuth())
	e.PUT("/users/me/password", h.ChangePassword, auth.WithJWTAuth())

	// User Information Retrieval
	e.GET("/users/:id", h.GetUserInfo, auth.RoleMiddleware("1", "2"))
	e.GET("/users/email", h.GetUserInfoByEmail, auth.RoleMiddleware("1", "2"))

	e.GET("/hello-world", helloWorld)
	e.Static("/swagger", "docs")
}

// helloWorld returns a simple "Hello World" message
// @Summary Test endpoint
// @Description Returns a "Hello World" string to verify the API is working
// @Tags General
// @Produce plain
// @Success 200 {string} string "Hello World"
// @Router /hello-world [get]
func helloWorld(c echo.Context) error {
	return c.String(http.StatusOK, "Hello World")
}

// Register handles user registration
// @Summary Register a new user
// @Description Registers a new user with email, password, and name. Sends a verification email.
// @Tags Users
// @Accept json
// @Produce json
// @Param request body object{email=string,password=string,name=string} true "User registration details"
// @Success 200 {object} object{message=string} "Verification email sent"
// @Failure 400 {object} object{error=string} "Invalid request or missing fields"
// @Failure 409 {object} object{error=string} "User already exists"
// @Failure 500 {object} object{error=string} "Internal server error"
// @Router /auth/register [post]
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
		"message": resp.Status,
	})
}

// VerifyEmail handles email verification
// @Summary Verify user email
// @Description Verifies a user's email using a token sent in the verification email
// @Tags Users
// @Produce json
// @Param token query string true "Verification token"
// @Success 200 {object} object{id=int32} "User ID"
// @Failure 400 {object} object{error=string} "Token is required or invalid"
// @Failure 500 {object} object{error=string} "Internal server error"
// @Router /auth/verify [get]
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

// LoginUser handles user login
// @Summary User login
// @Description Logs in a user with email and password, optionally with "remember me" functionality
// @Tags Users
// @Accept json
// @Produce json
// @Param request body object{email=string,password=string,rememberMe=boolean} true "Login credentials"
// @Success 200 {object} object{status=string,token=string} "Login successful with JWT token"
// @Failure 400 {object} object{error=string} "Invalid request"
// @Failure 401 {object} object{error=string} "Invalid credentials"
// @Failure 500 {object} object{error=string} "Internal server error"
// @Router /auth/login [post]
func (h *UserHandler) LoginUser(c echo.Context) error {
	var req struct {
		Email      string `json:"email"`
		Password   string `json:"password"`
		RememberMe bool   `json:"rememberMe"`
	}
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
	}

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

// ChangeInfo handles updating user information
// @Summary Update user info
// @Description Updates the authenticated user's information (email, name, phone, address)
// @Tags Users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body object{email=string,name=string,phoneNumber=string,address=string} true "Updated user info"
// @Success 200 {object} object{status=string,email=string,name=string,address=string,phoneNumber=string} "Info updated"
// @Failure 400 {object} object{error=string} "Invalid request"
// @Failure 500 {object} object{error=string} "Internal server error"
// @Router /user/me [put]
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

// ChangePassword handles password change
// @Summary Change user password
// @Description Changes the authenticated user's password
// @Tags Users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body object{oldPassword=string,newPassword=string} true "Password change details"
// @Success 200 {object} object{status=string} "Password changed successfully"
// @Failure 400 {object} object{error=string} "Invalid request or wrong old password"
// @Failure 500 {object} object{error=string} "Internal server error"
// @Router /users/me/password [put]
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

// GetUserInfo retrieves user information by ID
// @Summary Get user info by ID
// @Description Retrieves user information for a specific user ID (requires role 1 or 2)
// @Tags Users
// @Produce json
// @Security BearerAuth
// @Param id path int true "User ID"
// @Success 200 {object} object{id=int32,name=string,email=string,phoneNumber=string,address=string} "User info"
// @Failure 400 {object} object{error=string} "Invalid ID format"
// @Failure 404 {object} object{error=string} "User not found"
// @Failure 500 {object} object{error=string} "Internal server error"
// @Router /users/{id} [get]
func (h *UserHandler) GetUserInfo(c echo.Context) error {
	idStr := c.Param("id")
	if idStr == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "ID is required"})
	}

	id, err := strconv.ParseInt(idStr, 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid ID format, must be an integer"})
	}

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

// GetMyInfo retrieves the authenticated user's information
// @Summary Get authenticated user's info
// @Description Retrieves the information of the currently authenticated user
// @Tags Users
// @Produce json
// @Security BearerAuth
// @Success 200 {object} object{id=int32,name=string,email=string,phoneNumber=string,address=string} "User info"
// @Failure 400 {object} object{error=string} "Invalid request"
// @Failure 404 {object} object{error=string} "User not found"
// @Failure 500 {object} object{error=string} "Internal server error"
// @Router /users/me [get]
func (h *UserHandler) GetMyInfo(c echo.Context) error {
	id, err := auth.GetUserIDFromContext(c)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

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

// GetUserInfoByEmail retrieves user information by email
// @Summary Get user info by email
// @Description Retrieves user information using their email address
// @Tags Users
// @Produce json
// @Param email query string true "User email"
// @Success 200 {object} object{id=int32,name=string,email=string,phoneNumber=string,address=string} "User info"
// @Failure 400 {object} object{error=string} "Email is required"
// @Failure 404 {object} object{error=string} "User not found"
// @Failure 500 {object} object{error=string} "Internal server error"
// @Router /users/email [get]
func (h *UserHandler) GetUserInfoByEmail(c echo.Context) error {
	email := c.QueryParam("email")
	if email == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Email is required"})
	}

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

// ForgotPassword handles forgot password requests
// @Summary Forgot password
// @Description Sends a password reset email to the user
// @Tags Users
// @Accept json
// @Produce json
// @Param request body object{email=string} true "User email"
// @Success 200 {object} object{message=string} "Password reset email sent"
// @Failure 400 {object} object{error=string} "Invalid request"
// @Failure 404 {object} object{error=string} "User not found"
// @Failure 500 {object} object{error=string} "Internal server error"
// @Router /auth/forgot-password [post]
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

// ResetPassword handles password reset
// @Summary Reset password
// @Description Resets the user's password using a token from the forgot password email
// @Tags Users
// @Accept json
// @Produce json
// @Param request body object{userId=string,token=string,newPassword=string} true "Password reset details"
// @Success 200 {object} object{message=string} "Password reset successfully"
// @Failure 400 {object} object{error=string} "Invalid request or missing fields"
// @Failure 401 {object} object{error=string} "Invalid token"
// @Failure 500 {object} object{error=string} "Internal server error"
// @Router /auth/reset-password [post]
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
