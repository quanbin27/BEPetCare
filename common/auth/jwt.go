package auth

import (
	"fmt"
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"github.com/quanbin27/commons/config"
	"github.com/quanbin27/commons/utils"
	"net/http"
	"strconv"
	"time"
)

func CreateJWT(secret []byte, userID int32, seconds int64, roleId int32) (string, error) {
	expiration := time.Second * time.Duration(seconds)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":   strconv.Itoa(int(userID)),
		"expiredAt": time.Now().Add(expiration).Unix(),
		"role":      roleId,
	})
	tokenString, err := token.SignedString(secret)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}
func RoleMiddleware(allowedRoles ...string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Lấy token từ header hoặc query
			tokenString := utils.GetTokenFromRequest(c)
			if tokenString == "" {
				return c.JSON(http.StatusUnauthorized, map[string]string{
					"error": "unauthorized",
				})
			}

			// Xác thực token
			token, err := ValidateJWT(tokenString)
			if err != nil {
				return c.JSON(http.StatusUnauthorized, map[string]string{"error": err.Error()})
			}
			if !token.Valid {
				return c.JSON(http.StatusUnauthorized, map[string]string{"error": "invalid token"})
			}

			// Lấy claims từ token
			claims := token.Claims.(jwt.MapClaims)
			role, ok := claims["role"].(string)
			if !ok {
				return c.JSON(http.StatusUnauthorized, map[string]string{"error": "role not found in token"})
			}

			// Kiểm tra role có trong danh sách allowedRoles không
			for _, allowedRole := range allowedRoles {
				if role == allowedRole {
					userID, err := strconv.Atoi(claims["user_id"].(string))
					if err != nil {
						return c.JSON(http.StatusUnauthorized, map[string]string{"error": "invalid user_id"})
					}
					c.Set("user_id", int32(userID))
					c.Set("role", role)
					return next(c) // Role hợp lệ, tiếp tục
				}
			}

			return c.JSON(http.StatusForbidden, map[string]string{
				"error": "insufficient permissions",
			})
		}
	}
}
func WithJWTAuth() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Lấy token từ header hoặc query
			tokenString := utils.GetTokenFromRequest(c)
			if tokenString == "" {
				return c.JSON(http.StatusUnauthorized, map[string]string{
					"error": "unauthorized",
				})
			}

			// Xác thực token
			token, err := ValidateJWT(tokenString)
			if err != nil {
				return c.JSON(http.StatusUnauthorized, map[string]string{"error": err.Error()})
			}
			if !token.Valid {
				return c.JSON(http.StatusUnauthorized, map[string]string{"error": "invalid token"})
			}

			// Lấy claims từ token
			claims := token.Claims.(jwt.MapClaims)
			userID, err := strconv.Atoi(claims["user_id"].(string))
			if err != nil {
				return c.JSON(http.StatusUnauthorized, map[string]string{"error": err.Error()})
			}
			c.Set("user_id", userID)
			// Gọi hàm tiếp theo
			return next(c) // Tiếp tục xử lý yêu cầu
		}
	}
}
func GetUserIDFromContext(c echo.Context) (int32, error) {

	userID, ok := c.Get("user_id").(int32) // Giả sử types.User là kiểu dữ liệu của bạn
	if !ok || userID == 0 {
		return 0, echo.NewHTTPError(http.StatusUnauthorized, "User not found in context")
	}

	// Trả về userID
	return userID, nil // Giả sử user.ID là trường chứa userID
}

func ValidateJWT(tokenString string) (*jwt.Token, error) {
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(config.Envs.JWTSecret), nil
	})
}
