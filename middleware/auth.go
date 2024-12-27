package middleware

import (
	"a21hc3NpZ25tZW50/model"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

func Auth() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		cookie, err := ctx.Cookie("session_token")
		if err != nil {
			if ctx.GetHeader("Content-Type") == "application/json" {
				ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			} else {
				ctx.Redirect(http.StatusSeeOther, "/client/login")
			}
			ctx.Abort()
			return
		}

		token, err := jwt.ParseWithClaims(cookie, &model.Claims{}, func(token *jwt.Token) (interface{}, error) {
			return model.JwtKey, nil
		})

		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid token"})
			ctx.Abort()
			return
		}

		if claims, ok := token.Claims.(*model.Claims); ok && token.Valid {
			ctx.Set("email", claims.Email)
		} else {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			ctx.Abort()
			return
		}
		ctx.Next()
	}
}
