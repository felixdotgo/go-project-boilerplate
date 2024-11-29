package httpapi

import (
	"errors"
	"github.com/0x46656C6978/go-project-boilerplate/pkg/conv"
	"net/http"
	"time"

	"github.com/0x46656C6978/go-project-boilerplate/internal/config"
	"github.com/0x46656C6978/go-project-boilerplate/internal/entity"
	"github.com/0x46656C6978/go-project-boilerplate/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/sessions"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/google"
	"gorm.io/gorm"
)

type AuthHttpApi struct {
	s   service.UserServiceInterface
	cfg *config.Config
}

func NewAuthHttpApi(cfg *config.Config, s service.UserServiceInterface) *AuthHttpApi {
	return &AuthHttpApi{
		s:   s,
		cfg: cfg,
	}
}

func (u *AuthHttpApi) RegisterApiEndpoints(engine *gin.Engine) {
	engine.POST("/auth/login", u.Login)
	engine.POST("/auth/register", u.Register)
	engine.GET("/oauth", u.OAuthLogin)
	engine.GET("/oauth/callback", u.OAuthLoginCallback)
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (u *AuthHttpApi) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid request"})
		return
	}
	user, err := u.s.FindByEmail(c.Request.Context(), req.Email)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "user not found"})
		return
	}
	err = u.s.VerifyCredentials(c, user, req.Email, req.Password)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid credentials"})
		return
	}

	signedStr, err := u.generateJWTToken(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "internal error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "success",
		"data": gin.H{
			"token": signedStr,
		},
	})
}

type RegisterRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (u *AuthHttpApi) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid request"})
		return
	}
	user, err := u.s.FindByEmail(c.Request.Context(), req.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "internal error"})
		return
	}
	if user != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "user already exists"})
		return
	}

	user = &entity.User{}
	user.Email = req.Email
	err = user.SetPassword(req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "internal error"})
		return
	}

	err = u.s.Create(c.Request.Context(), user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "internal error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "success"})
}

func (u *AuthHttpApi) OAuthLogin(c *gin.Context) {
	store := sessions.NewCookieStore([]byte(u.cfg.JWT.Secret))
	store.MaxAge(conv.ToInt(u.cfg.JWT.Expire))
	store.Options.Path = "/"
	store.Options.HttpOnly = true // HttpOnly should always be enabled
	store.Options.Secure = false

	gothic.Store = store

	goth.UseProviders(
		google.New(u.cfg.OAuth.Google.ClientID, u.cfg.OAuth.Google.ClientSecret, u.cfg.OAuth.Google.RedirectURL),
	)
	gothic.BeginAuthHandler(c.Writer, c.Request)
}

func (u *AuthHttpApi) OAuthLoginCallback(c *gin.Context) {
	// you can dump the user variable to see what information is available
	user, err := gothic.CompleteUserAuth(c.Writer, c.Request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "internal error"})
		return
	}

	dbUser, err := u.s.FindByEmail(c.Request.Context(), user.Email)
	if dbUser == nil || errors.Is(err, gorm.ErrRecordNotFound) {
		dbUser = &entity.User{}
		dbUser.Email = user.Email
		dbUser.Password = ""
		err = u.s.Create(c.Request.Context(), dbUser)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "internal error"})
			return
		}
	}

	signedStr, err := u.generateJWTToken(dbUser)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "internal error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "success",
		"data": gin.H{
			"token": signedStr,
		},
	})
}

func (u *AuthHttpApi) generateJWTToken(user *entity.User) (string, error) {
	now := time.Now()
	exp := time.Duration(conv.ToInt64(u.cfg.JWT.Expire))
	claims := jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(now.Add(exp)),
		IssuedAt:  jwt.NewNumericDate(now),
		Issuer:    u.cfg.JWT.Issuer,
		Subject:   user.Email,
	}

	tok := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedStr, err := tok.SignedString([]byte(u.cfg.JWT.Secret))
	if err != nil {
		return "", err
	}
	return signedStr, nil
}
