package data

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"gorm.io/gorm"
)

type InviteConfirmStatus = int32

var (
	InviteConfirmStatusPadding  InviteConfirmStatus = 1 //  等待确认
	InviteConfirmStatusAccepted InviteConfirmStatus = 2 //  已接受
	InviteConfirmStatusRejected InviteConfirmStatus = 3 //  已拒绝
)

type UserDevice struct {
	Id            int64               `json:"id" gorm:"column:id;primary_key"`
	UserId        int64               `json:"user_id" gorm:"column:user_id;uniqueIndex:uniq_device"`
	DeviceId      string              `json:"device_id" gorm:"column:device_id;uniqueIndex:uniq_device;type:varchar(64)"`
	IsShare       bool                `json:"is_share" gorm:"column:is_share"`
	ShareFrom     int64               `json:"share_from" gorm:"column:share_from"`
	InviteConfirm InviteConfirmStatus `json:"invite_confirm" gorm:"column:invite_confirm"`
}

func (d *UserDevice) TableName() string {
	return "user_device"
}

type UserDeviceDao struct {
	log *log.Helper
	db  *gorm.DB
}

func NewUserDeviceDao(db *gorm.DB, logger log.Logger) *UserDeviceDao {
	return &UserDeviceDao{
		log: log.NewHelper(log.With(logger, "module", "data/userdevice")),
		db:  db,
	}
}

func (dao *UserDeviceDao) FindUserDevices(ctx context.Context, userid int64) ([]UserDevice, error) {
	var deviceIds []UserDevice
	err := dao.db.WithContext(ctx).Model(&UserDevice{}).Where("user_id = ?", userid).Pluck("device_id", &deviceIds).Error
	return deviceIds, err
}

func (dao *UserDeviceDao) InsertUserDevice(ctx context.Context, userDevice UserDevice) error {
	return dao.db.WithContext(ctx).Create(&userDevice).Error
}

func (dao *UserDeviceDao) UnbindByDevice(ctx context.Context, deviceId string) error {
	return dao.db.WithContext(ctx).Delete(&UserDevice{}, "device_id = ?", deviceId).Error
}

func (dao *UserDeviceDao) FindDeviceOwner(ctx context.Context, deviceId string) (*UserDevice, error) {
	var ud UserDevice
	if err := dao.db.WithContext(ctx).Model(&UserDevice{}).Where("device_id = ? and is_share = 0", deviceId).First(&ud).Error; err != nil {
		return nil, err
	}
	return &ud, nil
}

func (dao *UserDeviceDao) GetDevice(ctx context.Context, device UserDevice) (*UserDevice, error) {
	var ud UserDevice
	if err := dao.db.WithContext(ctx).Model(&UserDevice{}).Where(device).First(&ud).Error; err != nil {
		return nil, err
	}
	return &ud, nil
}

func (dao *UserDeviceDao) Update(ctx context.Context, device UserDevice) (*UserDevice, error) {
	if err := dao.db.WithContext(ctx).Model(&UserDevice{}).Where("id = ? ", device.Id).Updates(&device).Error; err != nil {
		return nil, err
	}
	return &device, nil
}

func (dao *UserDeviceDao) Delete(ctx context.Context, device UserDevice) error {
	return dao.db.WithContext(ctx).Model(&UserDevice{}).Delete(&device).Error
}

func (dao *UserDeviceDao) GetDevices(ctx context.Context, device UserDevice) ([]UserDevice, error) {
	var ud []UserDevice
	if err := dao.db.WithContext(ctx).Model(&UserDevice{}).Where(device).Find(&ud).Error; err != nil {
		return nil, err
	}
	return ud, nil
}
