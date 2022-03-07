/**
 * @Time: 2022/3/7 14:36
 * @Author: yt.yin
 */

package jwt

import (
	"github.com/gin-gonic/gin"
	"github.com/goworkeryyt/go-core/jwt"
	"net/http"
	"strings"
)

// JWTAuth JWT 认证中间件
func JWTAuth() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		path := ctx.Request.URL.Path
		if strings.Contains(path, "swagger") {
			ctx.Next()
			return
		}
		if strings.Contains(path, "login") || strings.Contains(path, "health") || strings.Contains(path, "captcha") {
			ctx.Next()
			return
		}
		token := ctx.Request.Header.Get("ACCESS_TOKEN")
		if token == "" {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"code":    -1,
				"message": "请求未携带token,无访问权限！",
			})
			ctx.Abort()
			return
		}
		j := jwt.NewJWT()
		// 解析token包含的信息
		claims, err := j.ResolveToken(token)
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"code":    -1,
				"message": err.Error(),
			})
			ctx.Abort()
			return
		}
		ctx.Set("claims", claims)
		ctx.Next()
	}
}
