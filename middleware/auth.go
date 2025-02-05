package middleware

import (
	"github.com/gin-gonic/gin"
	"go-gin-demo/api"
	"net/http"
)

const (
	TokenName   = "Authorization"
	TokenPREFIX = "Bearer: "
)

func tokenErr(c *gin.Context) {
	api.Fail(c, api.ResponseJson{
		Status: http.StatusUnauthorized,
		Msg:    "非法的token",
	})
}

//func Auth() func(c *gin.Context) {
//	return func(c *gin.Context) {
//		token := c.GetHeader(TokenName)
//		if token == "" || !strings.HasPrefix(token, TokenPREFIX) {
//			tokenErr(c)
//			return
//		}
//
//		token = token[len(TokenPREFIX):]
//		iJwtCustomClaims, err := utils.ParseToken(token)
//		currentUserID := iJwtCustomClaims.ID
//		if err != nil || currentUserID == 0 {
//			tokenErr(c)
//			return
//		}
//		// token 和 redis 不一致
//		stUserId := strconv.Itoa(int(currentUserID))
//		stRedisUserIdKey := strings.Replace(constans.LOGIN_USER_TOKEN_REDIS_KEY, "{ID}", stUserId, -1)
//		confTokenExpire := viper.GetDuration("jwt.tokenExpire") * time.Minute
//		stRedisToken, err := global.RedisClient.Get(stRedisUserIdKey)
//
//		if err != nil || token != stRedisToken {
//			tokenErr(c)
//			return
//		}
//		// token 失效，直接返回
//		nTokenExpireDuration, err := global.RedisClient.GetExpireTTl(stRedisUserIdKey)
//		if err != nil || nTokenExpireDuration <= 0 {
//			tokenErr(c)
//			return
//		}
//		// token的续期
//		if nTokenExpireDuration.Seconds() < confTokenExpire.Seconds() {
//			stNewToken, err := user_service.GenTokenAndSetToken2Redis(currentUserID, iJwtCustomClaims.Name)
//			if err != nil {
//				tokenErr(c)
//				return
//			}
//			c.Header("token", stNewToken)
//		}
//
//		c.Set(constans.LOGIN_USER, user_model.LoginUser{
//			ID:   currentUserID,
//			Name: iJwtCustomClaims.Name,
//		})
//		c.Next()
//
//	}
//}
