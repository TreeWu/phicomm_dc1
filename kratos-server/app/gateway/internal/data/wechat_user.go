package data

import (
	"context"
	"gorm.io/gorm"
)

type WechatUser struct {
	ID         int64  `gorm:"column:id;primary_key" json:"id"`
	AppId      string `gorm:"column:app_id;type:varchar(128)" json:"app_id"`
	Openid     string `gorm:"column:openid;type:varchar(256)" json:"openid"`
	Unionid    string `gorm:"column:unionid;type:varchar(256)" json:"unionid"`
	SessionKey string `gorm:"column:session_key" json:"session_key"`
	Avatar     string `gorm:"column:avatar" json:"avatar"`
	Email      string `gorm:"column:email;type:varchar(64)" json:"email"`
	Nickname   string `gorm:"column:nickname;type:varchar(64)" json:"nickname"`
	CreateTime int64  `gorm:"column:create_time" json:"create_time"`
	UserCode   string `gorm:"column:user_code;type:varchar(64);index:idx_user_code" json:"user_code"`
}

// TableName sets the insert table name for this struct type
func (s *WechatUser) TableName() string {
	return "wechat_user"
}

type WechatUserRepo struct {
	db *gorm.DB
}

func NewWechatUserRepo(data *Data) *WechatUserRepo {
	return &WechatUserRepo{
		db: data.db,
	}
}

func (w *WechatUserRepo) Create(ctx context.Context, user WechatUser) error {
	return w.db.Create(&user).Error
}
func (w *WechatUserRepo) Update(ctx context.Context, user WechatUser) error {
	return w.db.Debug().Model(&user).Where("user_code = ? ", user.UserCode).Updates(user).Error
}
func (w *WechatUserRepo) Find(ctx context.Context, user WechatUser) (*WechatUser, error) {
	var u WechatUser
	if err := w.db.Model(&WechatUser{}).Where(user).First(&u).Error; err != nil {
		return nil, err
	}
	return &u, nil
}
