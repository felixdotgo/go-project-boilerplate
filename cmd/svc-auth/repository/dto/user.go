package dto

import (
	"github.com/0x46656C6978/go-project-boilerplate/cmd/svc-auth/entity"
	"github.com/0x46656C6978/go-project-boilerplate/cmd/svc-auth/repository/model"
)

func FromUserEntity(ent *entity.User) *model.User {
	if ent == nil {
		return nil
	}
	return &model.User{
		ID:            ent.ID,
		Email:         ent.Email,
		Password:      ent.Password,
		Name:          ent.Name,
		FirstName:     ent.FirstName,
		LastName:      ent.LastName,
		AvatarURL:     ent.AvatarURL,
		EmailVerified: ent.EmailVerified,
		Locale:        ent.Locale,
		CreatedAt:     ent.CreatedAt,
		UpdatedAt:     ent.UpdatedAt,
	}
}

func ToUserEntity(mdl *model.User) *entity.User {
	if mdl == nil {
		return nil
	}
	return &entity.User{
		ID:            mdl.ID,
		Email:         mdl.Email,
		Password:      mdl.Password,
		Name:          mdl.Name,
		FirstName:     mdl.FirstName,
		LastName:      mdl.LastName,
		AvatarURL:     mdl.AvatarURL,
		EmailVerified: mdl.EmailVerified,
		Locale:        mdl.Locale,
		CreatedAt:     mdl.CreatedAt,
		UpdatedAt:     mdl.UpdatedAt,
	}
}

