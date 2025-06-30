package middleware

import (
	"gongChang/config"
	"gongChang/models"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	UserID string      `json:"user_id"`
	Role   models.UserRole `json:"role"`
	jwt.RegisteredClaims
}

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			log.Printf("Authorization header is missing")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
			c.Abort()
			return
		}

		// 检查Bearer token格式
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			log.Printf("Invalid authorization header format: %s", authHeader)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization header format"})
			c.Abort()
			return
		}

		tokenString := parts[1]
		cfg, err := config.LoadConfig()
		if err != nil {
			log.Printf("Failed to load configuration: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load configuration"})
			c.Abort()
			return
		}

		// 检查JWT密钥是否已配置
		if cfg.JWT.Secret == "${JWT_SECRET}" {
			log.Printf("JWT secret key not configured")
			c.JSON(http.StatusInternalServerError, gin.H{"error": "JWT secret key not configured"})
			c.Abort()
			return
		}

		// 解析和验证token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return []byte(cfg.JWT.Secret), nil
		})

		if err != nil {
			log.Printf("Token validation failed: %v", err)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		if !token.Valid {
			log.Printf("Token is invalid")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token is invalid"})
			c.Abort()
			return
		}

		// 将用户ID和角色添加到上下文中
		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			if userID, ok := claims["user_id"].(string); ok {
				c.Set("user_id", userID)
			}
			if role, ok := claims["role"].(string); ok {
				c.Set("user_role", role)
			}
		}

		c.Next()
	}
}

func GenerateToken(userID string, role models.UserRole, secret string) (string, error) {
	claims := Claims{
		UserID: userID,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

func ValidateToken(tokenString string) (*Claims, error) {
	cfg, err := config.LoadConfig()
	if err != nil {
		return nil, err
	}

	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(cfg.JWT.Secret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, jwt.ErrSignatureInvalid
}

// FactoryRoleMiddleware 工厂角色验证中间件
func FactoryRoleMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 先执行认证中间件
		AuthMiddleware()(c)
		
		// 如果认证失败，直接返回
		if c.IsAborted() {
			return
		}
		
		// 检查用户角色
		userRole := c.GetString("user_role")
		if userRole != string(models.RoleFactory) {
			c.JSON(http.StatusForbidden, gin.H{"error": "仅工厂角色可以访问此功能"})
			c.Abort()
			return
		}
		
		c.Next()
	}
} 