// pages/device/devices.js
const {updateDevice} = require("../../api/index");
Page({
    data: {
        devices: [],
        listenId: -1,
    },
    onLoad(options) {
    },
    onShow() {
        const app = getApp();
        const that = this;
        this.data.listenId = app.watch("devices", (v) => {
            that.setData({
                devices: v,
            });
        });
        app.getDevcieList();
    },
    onHide() {
        if (this.data.listenId != -1) {
            getApp().unwatch(this.data.listenId);
        }
        this.data.listenId = -1;
    },
    onUnload() {
        if (this.data.listenId != -1) {
            getApp().unwatch(this.data.listenId);
        }
        this.data.listenId = -1;
    },
    onItemClicked(event) {
        const deviceId = event.currentTarget.dataset["deviceId"];
        wx.navigateTo({
            url: "/pages/devices/deviceInfo?deviceId=" + deviceId,
        });
    },
    switchName(e) {
        let app = getApp();
        let deviceId = e.currentTarget.dataset.deviceId;
        let req = {
            deviceId: deviceId,
        };
        wx.showModal({
            title: "请输入新名称",
            editable: true,
            placeholderText: app.globalData.deviceMap[deviceId].name,
            success: (res) => {
                if (res.confirm) {
                    if (res.content === "") {
                        wx.showToast({
                            title: "名称必填",
                            icon: "error",
                            mask: true,
                        });
                        return;
                    }
                    req.name = res.content;
                    updateDevice(req).then(() => app.getDevcieList());
                }
            },
        });
    },
    async onPullDownRefresh() {
        let app = getApp();
        await app.getDevcieList();
        wx.stopPullDownRefresh();
    },
});
