package api

import (
	"context"
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"github.com/gnasnik/titan-container-api/core/dao"
	"github.com/gnasnik/titan-container-api/core/errors"
	"github.com/gnasnik/titan-container-api/core/generated/model"
	"github.com/gnasnik/titan-container-api/core/logger"
	"github.com/gnasnik/titan-container-api/utils"
	"github.com/mssola/user_agent"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"time"
)

const (
	loginStatusFailure = iota
	loginStatusSuccess
)

type login struct {
	Username string `form:"username" json:"username" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
}

type loginResponse struct {
	Token  string `json:"token"`
	Expire string `json:"expire"`
}

var identityKey = "id"

func jwtGinMiddleware(secretKey string) (*jwt.GinJWTMiddleware, error) {
	return jwt.New(&jwt.GinJWTMiddleware{
		Realm:             "User",
		Key:               []byte(secretKey),
		Timeout:           time.Hour,
		MaxRefresh:        24 * time.Hour,
		IdentityKey:       identityKey,
		SendAuthorization: true,
		PayloadFunc: func(data interface{}) jwt.MapClaims {
			if v, ok := data.(*model.User); ok {
				return jwt.MapClaims{
					identityKey: v.Username,
				}
			}
			return jwt.MapClaims{}
		},
		IdentityHandler: func(c *gin.Context) interface{} {
			claims := jwt.ExtractClaims(c)
			return &model.User{
				Username: claims[identityKey].(string),
			}
		},
		LoginResponse: func(c *gin.Context, code int, token string, expire time.Time) {
			c.JSON(http.StatusOK, gin.H{
				"code": 0,
				"data": loginResponse{
					Token:  token,
					Expire: expire.Format(time.RFC3339),
				},
			})
		},
		LogoutResponse: func(c *gin.Context, code int) {
			c.JSON(http.StatusOK, gin.H{
				"code": 0,
			})
		},
		Authenticator: func(c *gin.Context) (interface{}, error) {
			var loginParams login
			if err := c.ShouldBind(&loginParams); err != nil {
				return "", jwt.ErrMissingLoginValues
			}
			userID := loginParams.Username
			password := loginParams.Password
			userAgent := c.Request.Header.Get("User-Agent")
			ua := user_agent.New(userAgent)
			os := ua.OS()
			explorer, _ := ua.Browser()
			clientIP := utils.GetClientIP(c.Request)
			location := utils.GetLocationByIP(clientIP)

			loginLog := &model.LoginLog{
				Ipaddr:        clientIP,
				Browser:       explorer,
				Os:            os,
				LoginLocation: location,
			}

			user, err := loginByPassword(c.Request.Context(), userID, password)
			if err != nil {
				loginLog.Status = loginStatusFailure
				loginLog.Msg = err.Error()
				logger.AddLoginLog(loginLog)
				return nil, err
			}

			loginLog.LoginUsername = userID
			loginLog.Status = loginStatusSuccess
			loginLog.Msg = "success"
			logger.AddLoginLog(loginLog)

			return user, nil
		},
		//Authorizator: func(data interface{}, c *gin.Context) bool {
		//	if v, ok := data.(*model2.User); ok && v.Username == "admin" {
		//		return true
		//	}
		//
		//	return false
		//},
		Unauthorized: func(c *gin.Context, code int, message string) {
			c.JSON(http.StatusOK, gin.H{
				"code":    401,
				"msg":     message,
				"success": false,
			})
		},
		// TokenLookup is a string in the form of "<source>:<name>" that is used
		// to extract token from the request.
		// Optional. Default value "header:Authorization".
		// Possible values:
		// - "header:<name>"
		// - "query:<name>"
		// - "cookie:<name>"
		// - "param:<name>"
		//TokenLookup: "header: Authorization, query: token, cookie: jwt",
		TokenLookup: "header: JwtAuthorization",
		// TokenLookup: "query:token",
		// TokenLookup: "cookie:token",

		// TokenHeadName is a string in the header. Default value is "Bearer"
		TokenHeadName: "Bearer",

		// TimeFunc provides the current time. You can override it to use another time value. This is useful for testing or if your server uses a different time zone than your tokens.
		TimeFunc: time.Now,

		RefreshResponse: func(c *gin.Context, code int, token string, t time.Time) {
			c.Next()
		},
	})
}

func loginByPassword(ctx context.Context, username, password string) (interface{}, error) {
	user, err := dao.GetUserByUsername(ctx, username)
	if err != nil {
		log.Errorf("get user by username: %v", err)
		return nil, errors.ErrUserNotFound
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PassHash), []byte(password)); err != nil {
		log.Errorf("can't compare hash %s ans password %s: %v", user.PassHash, password, err)
		return nil, errors.ErrInvalidPassword
	}

	return &model.User{Uuid: user.Uuid, Username: user.Username, Role: user.Role}, nil
}
