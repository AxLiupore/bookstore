package serializer

import (
	"bookstore/conf"
	"bookstore/model"
)

type User struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Status   string `json:"status"`
	Avatar   string `json:"avatar"`
	CreateAt int64  `json:"create_at"`
}

func BuildUser(user *model.User) *User {
	return &User{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
		Status:   user.Status,
		Avatar:   conf.Host + conf.HttpPort + conf.AvatarPath + user.Avatar,
		CreateAt: user.CreatedAt.Unix(),
	}
}
