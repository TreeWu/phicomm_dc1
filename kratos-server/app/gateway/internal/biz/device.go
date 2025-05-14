package biz

import (
	"context"
	"errors"
	"github.com/duke-git/lancet/v2/condition"
	"github.com/duke-git/lancet/v2/slice"
	"github.com/go-kratos/kratos/v2/log"
	v1 "github.com/treewu/phicomm_dc1/api/gateway/v1"
	"github.com/treewu/phicomm_dc1/app/common/data"
	. "github.com/treewu/phicomm_dc1/app/common/data"
	"github.com/treewu/phicomm_dc1/app/gateway/internal/conf"
	data2 "github.com/treewu/phicomm_dc1/app/gateway/internal/data"
	"github.com/treewu/phicomm_dc1/pkg/server/dc1server"
	"gorm.io/gorm"
	"strings"
	"time"
)

type DeviceBiz struct {
	wechatRepo     *data2.WechatUserRepo
	deviceRepo     *data.DeviceDao
	log            *log.Helper
	userDeviceRepo *data.UserDeviceDao
	cf             *conf.Data_Wechat
	*CommonBiz
}

func NewDeviceBiz(
	logger log.Logger,
	repo *data.DeviceDao,
	wechatdao *data2.WechatUserRepo,
	userdevice *data.UserDeviceDao,
	biz *CommonBiz,
	cf *conf.Data,
) *DeviceBiz {
	return &DeviceBiz{
		userDeviceRepo: userdevice,
		wechatRepo:     wechatdao,
		deviceRepo:     repo,
		log:            log.NewHelper(log.With(logger, "module", "biz/device")),
		CommonBiz:      biz,
		cf:             cf.Wechat,
	}

}

func (d *DeviceBiz) ListDevice(ctx context.Context, req *v1.ListDeviceRequest) (*v1.ListDeviceReply, error) {
	user, err := d.GetUser(ctx)
	if err != nil {
		return nil, err
	}
	uds, err := d.userDeviceRepo.FindUserDevices(ctx, user.ID)
	if err != nil {
		return nil, err
	}
	var ids []string
	for _, ud := range uds {
		if !ud.IsShare || (ud.IsShare && ud.InviteConfirm == data.InviteConfirmStatusAccepted) {
			ids = append(ids, ud.DeviceId)
		}
	}
	var ds []*v1.Dc1Device
	if len(ids) != 0 {
		devices, err := d.deviceRepo.ListDevice(ctx, ids)
		if err != nil {
			return nil, err
		}
		for _, device := range devices {
			ds = append(ds, d.ToProto(device))
		}
	}
	return &v1.ListDeviceReply{Devices: ds}, nil
}

func (d *DeviceBiz) UpdateDevice(ctx context.Context, req *v1.UpdateDeviceRequest) (*v1.UpdateDeviceReply, error) {
	return &v1.UpdateDeviceReply{}, d.deviceRepo.UpdateDevice(ctx, req)
}

func (d *DeviceBiz) BindDevice(ctx context.Context, req *v1.DeviceConnectReq) (*v1.DeviceConnectReply, error) {
	user, err := d.GetUser(ctx)
	if err != nil {
		return nil, err
	}
	mac := strings.ToUpper(req.Mac)
	device, err := d.deviceRepo.FindByMac(mac)
	if err != nil {
		return nil, err
	}
	if time.Now().Add(-time.Minute*5).UnixMilli() > device.LastActivatedAt {
		return nil, ServerError("超时未绑定，请重新对排插联网")
	}
	owner, err := d.userDeviceRepo.FindDeviceOwner(ctx, device.DeviceId)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ServerError("服务异常")
	}

	if owner != nil && owner.UserId != user.ID {
		d.log.Warnf("设备已被其他用户绑定 ,deviceId %s", device.DeviceId)
		return nil, ResourceWithoutPermissionError("设备已被其他用户绑定")
	}
	if owner != nil && owner.UserId == user.ID {
		d.log.Warnf("设备已绑定 ,deviceId %s", device.DeviceId)
		return &v1.DeviceConnectReply{
			DeviceId: device.DeviceId,
		}, nil
	}

	err = d.userDeviceRepo.InsertUserDevice(ctx, data.UserDevice{
		UserId:        user.ID,
		DeviceId:      device.DeviceId,
		IsShare:       false,
		ShareFrom:     0,
		InviteConfirm: 0,
	})
	return &v1.DeviceConnectReply{
		DeviceId: device.DeviceId,
	}, err
}

func (d *DeviceBiz) ToProto(device data.Device) *v1.Dc1Device {
	v := &v1.Dc1Device{
		Id:              device.ID,
		DeviceId:        device.DeviceId,
		Name:            device.Name,
		DeviceType:      device.DeviceType,
		DetalKwh:        device.DetalKWh,
		Recover:         device.Recover,
		Switch1Name:     device.Switch1Name,
		Switch2Name:     device.Switch2Name,
		Switch3Name:     device.Switch3Name,
		LastActivatedAt: device.LastActivatedAt,
		IsOnline:        time.Now().Add(-d.cf.OnlineInterval.AsDuration()).UnixMilli() < device.LastMessageAt,
	}

	if device.Status != nil {
		command := &v1.Dc1Command{
			Switch_3:   dc1server.Int32Pointer(*device.Status / 1000),
			Switch_2:   dc1server.Int32Pointer(*device.Status % 1000 / 100),
			Switch_1:   dc1server.Int32Pointer(*device.Status % 1000 % 100 / 10),
			SwitchMain: dc1server.Int32Pointer(*device.Status % 1000 % 100 % 10),
		}
		v.Command = command
	}
	if device.I != nil {
		v.I = *device.I
	}
	if device.P != nil {
		v.P = *device.P
	}
	if device.V != nil {
		v.V = *device.V
	}

	return v
}

// CheckDeviceOwner 检查设备是否属于当前用户
func (d *DeviceBiz) CheckDeviceOwner(ctx context.Context, deviceId string) (*data.UserDevice, error) {
	user, err := d.GetUser(ctx)
	if err != nil {
		return nil, err
	}
	getDevice, err := d.userDeviceRepo.GetDevice(ctx, data.UserDevice{
		UserId:   user.ID,
		DeviceId: deviceId,
	})
	if err != nil {
		return nil, err
	}
	if getDevice.IsShare {
		return nil, ResourceWithoutPermissionError("非设备拥有者")
	}
	return getDevice, nil
}

// CheckDevicePermission 检测当前用户是否有设备操作权限
// 如果是邀请设备，应该检测是否被邀请
// 非邀请设备即是Owner
func (d *DeviceBiz) CheckDevicePermission(ctx context.Context, deviceId string) (*data.UserDevice, error) {
	user, err := d.GetUser(ctx)
	if err != nil {
		return nil, err
	}
	device, err := d.userDeviceRepo.GetDevice(ctx, data.UserDevice{
		DeviceId: deviceId,
		UserId:   user.ID,
	})
	if err != nil {
		return nil, err
	}
	if (device.IsShare && device.InviteConfirm == data.InviteConfirmStatusAccepted) || (!device.IsShare) {
		return device, nil
	}
	return nil, ResourceWithoutPermissionError("没有操作权限")
}

// ShareInvite 邀请设备控制
//
// - 检测是否拥有者
// - 检查是否重复邀请
// - 插入邀请记录
func (d *DeviceBiz) ShareInvite(ctx context.Context, req *v1.ShareInviteReq) (*v1.ShareInviteReply, error) {

	user, err := d.GetUser(ctx)
	if err != nil {
		return nil, err
	}
	getDevice, err := d.userDeviceRepo.GetDevice(ctx, data.UserDevice{
		UserId:   user.ID,
		DeviceId: req.DeviceId,
	})
	if err != nil {
		return nil, err
	}
	if getDevice.IsShare {
		return nil, ResourceWithoutPermissionError("非设备拥有者")
	}

	find, err := d.wechatRepo.Find(ctx, data2.WechatUser{UserCode: req.UserCode})
	if err != nil {
		return nil, RecordNotFoundError("用户不存在")
	}
	device, err := d.userDeviceRepo.GetDevice(ctx, data.UserDevice{UserId: find.ID, DeviceId: req.DeviceId})
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, SqlError(err.Error())
	}

	if device != nil {
		if device.InviteConfirm != data.InviteConfirmStatusAccepted {
			device.InviteConfirm = data.InviteConfirmStatusPadding
			_, err = d.userDeviceRepo.Update(ctx, *device)
		} else {
			return nil, ServerError("重复邀请")
		}
	} else {
		err = d.userDeviceRepo.InsertUserDevice(ctx, data.UserDevice{
			UserId:        find.ID,
			DeviceId:      req.DeviceId,
			IsShare:       true,
			ShareFrom:     user.ID,
			InviteConfirm: data.InviteConfirmStatusPadding,
		})
	}

	return &v1.ShareInviteReply{}, err
}

// ShareConfirm 确认分享
func (d *DeviceBiz) ShareConfirm(ctx context.Context, req *v1.ShareConfirmReq) (*v1.ShareConfirmReply, error) {
	user, err := d.GetUser(ctx)
	if err != nil {
		return nil, err
	}
	device, err := d.userDeviceRepo.GetDevice(ctx, data.UserDevice{
		Id: req.ShareId,
	})
	if err != nil {
		return nil, err
	}
	if device.UserId != user.ID {
		return nil, ResourceWithoutPermissionError("没有操作权限")
	}
	if req.Confirm {
		device.InviteConfirm = data.InviteConfirmStatusAccepted
	} else {
		device.InviteConfirm = data.InviteConfirmStatusRejected
	}
	_, err = d.userDeviceRepo.Update(ctx, *device)
	return &v1.ShareConfirmReply{}, err
}

func (d *DeviceBiz) ShareRevoke(ctx context.Context, req *v1.ShareRevokeReq) (*v1.ShareRevokeReply, error) {
	user, err := d.GetUser(ctx)
	if err != nil {
		return nil, err
	}
	device, err := d.userDeviceRepo.GetDevice(ctx, data.UserDevice{Id: req.ShareId})
	if err != nil {
		return nil, err
	}
	if device.ShareFrom != user.ID {
		return nil, ResourceWithoutPermissionError("没有操作权限")
	}
	return &v1.ShareRevokeReply{}, d.userDeviceRepo.Delete(ctx, data.UserDevice{Id: req.ShareId})
}

func (d *DeviceBiz) GetShareList(ctx context.Context, req *v1.GetShareListReq) (*v1.GetShareListReply, error) {
	var resp = v1.GetShareListReply{}
	user, err := d.GetUser(ctx)
	if err != nil {
		return nil, err
	}
	devices, err := d.userDeviceRepo.GetDevices(ctx, data.UserDevice{ShareFrom: user.ID})
	if err != nil {
		return nil, err
	}
	resp.OwnerShares = slice.Map(devices, func(_ int, device data.UserDevice) *v1.GetShareListReply_ShareInfo {
		id, _ := d.deviceRepo.FindByDeviceId(device.DeviceId)
		find, _ := d.wechatRepo.Find(ctx, data2.WechatUser{ID: device.UserId})

		return &v1.GetShareListReply_ShareInfo{
			Id:          device.Id,
			DeviceId:    device.DeviceId,
			DeviceName:  condition.Ternary(id != nil, id.Name, ""),
			UserName:    condition.Ternary(find != nil, find.Nickname, ""),
			ShareStatus: v1.ShareStatus(device.InviteConfirm),
		}
	})

	devices, err = d.userDeviceRepo.GetDevices(ctx, data.UserDevice{UserId: user.ID, IsShare: true})
	if err != nil {
		return nil, err
	}
	resp.FromOtherShares = slice.Map(devices, func(_ int, device data.UserDevice) *v1.GetShareListReply_ShareInfo {
		id, _ := d.deviceRepo.FindByDeviceId(device.DeviceId)
		find, _ := d.wechatRepo.Find(ctx, data2.WechatUser{ID: device.UserId})

		return &v1.GetShareListReply_ShareInfo{
			Id:          device.Id,
			DeviceId:    device.DeviceId,
			DeviceName:  condition.Ternary(id != nil, id.Name, ""),
			UserName:    condition.Ternary(find != nil, find.Nickname, ""),
			ShareStatus: v1.ShareStatus(device.InviteConfirm),
		}
	})

	return &resp, nil
}
