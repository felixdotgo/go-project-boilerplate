package httpapi

import (
	"context"
	"time"

	"github.com/0x46656C6978/go-project-boilerplate/cmd/api/config"
	"github.com/0x46656C6978/go-project-boilerplate/cmd/api/entity"
	"github.com/0x46656C6978/go-project-boilerplate/cmd/api/service"
	"github.com/0x46656C6978/go-project-boilerplate/pkg/conv"
	v1 "github.com/0x46656C6978/go-project-boilerplate/rpc/api/auth/v1"
	"github.com/golang-jwt/jwt/v5"
)

type AuthHttpApi struct {
	v1.UnimplementedAuthServiceServer
	s   service.UserServiceInterface
	cfg *config.Config
}

func NewAuthServiceServer(cfg *config.Config, s service.UserServiceInterface) v1.AuthServiceServer {
	return &AuthHttpApi{
		s:   s,
		cfg: cfg,
	}
}

func (u *AuthHttpApi) Login(ctx context.Context, req *v1.Auth_LoginRequest) (*v1.Auth_LoginResponse, error) {
	user, err := u.s.FindByEmail(ctx, req.GetData().GetEmail())
	if err != nil {
		return nil, err
	}
	err = u.s.VerifyCredentials(ctx, user, req.Data.GetEmail(), req.GetData().GetPassword())
	if err != nil {
		return nil, err
	}

	signedStr, err := u.generateJWTToken(user)
	if err != nil {
		return nil, err
	}
    return &v1.Auth_LoginResponse{
		Data: &v1.Auth_LoginResponseData{
			Token: signedStr,
		},
	}, nil
}

func (u *AuthHttpApi) Regiter(ctx context.Context, req *v1.Auth_RegisterRequest) (*v1.Auth_RegisterResponse, error) {
	user, err := u.s.FindByEmail(ctx, req.GetData().GetEmail())
	if err != nil {
		return nil, err
	}
	if user != nil {
		return nil, err
	}

	user = &entity.User{}
	user.Email = req.GetData().GetEmail()
	err = user.SetPassword(req.GetData().GetPassword())
	if err != nil {
		return nil, err
	}

	err = u.s.Create(ctx, user)
	if err != nil {
		return nil, err
	}

    return &v1.Auth_RegisterResponse{
	}, nil
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
