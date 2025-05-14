// pages/device-share/device-share.js

const {
    getShareList,
    shareConfirm,
    shareRevoke,
    shareInvite,
} = require("../../api/index");

Page({
    data: {
        statusMap: {
            pending: "等待确认",
            accepted: "已接受",
            rejected: "已拒绝",
            revoked: "已撤回",
        },
        activeTab: "shared", // 当前激活的标签页
        sharedList: [
            // 我分享的设备列表
        ],
        receivedList: [
            // 分享给我的设备列表
        ],
    },

    // 切换标签页
    switchTab: function (e) {
        const tab = e.currentTarget.dataset.tab;
        this.setData({
            activeTab: tab,
        });
    },
    // 接受分享
    handleAccept: async function (e) {
        const id = e.currentTarget.dataset.id;
        await shareConfirm({shareId: id, confirm: true});
        this.loadSharedList();
    },
    // 拒绝分享
    handleReject: async function (e) {
        const id = e.currentTarget.dataset.id;
        await shareConfirm({shareId: id, confirm: false});
        this.loadSharedList();
    },
    // 撤回分享
    handleRevoke: async function (e) {
        const id = e.currentTarget.dataset.id;
        await shareRevoke({shareId: id});
        this.loadSharedList();
    },
    // 退出分享
    handleExit: async function (e) {
        const id = e.currentTarget.dataset.id;
        await shareConfirm({shareId: id, confirm: false});
        this.loadSharedList();
    },
    onShow: function () {
        this.loadSharedList();
    },

    // 加载我分享的设备列表
    loadSharedList: async function () {
        try {
            let res = await getShareList();
            this.setData({
                sharedList: res.ownerShares,
                receivedList: res.fromOtherShares,
            });
        } catch (err) {
            this.setData({
                sharedList: [],
                receivedList: [],
            });
        }
    },
    handleReturn() {
        wx.switchTab({
            url: "/pages/devices/devices",
        });
    },
    handelCreate() {
        const that = this;
        const devices = getApp().globalData.devices;
        const deviceNames = devices.map((item) => item.name);
        wx.showActionSheet({
            alertText: "设备选择",
            itemList: deviceNames,
            success(res) {
                let device = devices[res.tapIndex];
                wx.showModal({
                    editable: true,
                    placeholderText: "输入对方的用户ID",
                    title: `是否分享设备【${device.name}】`,
                    complete: (res) => {
                        if (res.confirm) {
                            shareInvite({userCode: res.content, deviceId: device.deviceId})
                                .then(() => {
                                    that.loadSharedList();
                                })
                                .catch((err) => {
                                    setTimeout(() => {
                                        wx.showToast({
                                            title: err.data.message,
                                            icon: "error",
                                        });
                                    }, 50);
                                });
                        }
                    },
                });
            },
        });
    },
    async onPullDownRefresh() {
        await this.loadSharedList()
        wx.stopPullDownRefresh()
    }
});
