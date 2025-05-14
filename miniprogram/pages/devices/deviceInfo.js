const {updateDevice} = require("../../api/index");
const {postMark} = require("../../utils/request");
Page({
    data: {
        deviceId: "",
        device: {},
        command: {},
        listenId: -1,
    },
    onLoad(options) {
        console.log(options)
        const deviceId = options.deviceId;
        if (!deviceId) {
            wx.showModal({
                content: "设备不存在",
                showCancel: false,
                complete: (res) => {
                    wx.navigateBack();
                },
            });
        } else {
            this.setData({
                deviceId: deviceId,
            });
        }
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
    onShow() {
        const app = getApp();
        this.data.listenId = app.watch("deviceMap", this.flushDevice.bind(this));
        app.getDevcieList();
    },
    flushDevice(devices) {
        let that = this;
        let device = devices[that.data.deviceId];
        const keys = Object.keys(device.command);
        const command = {};
        keys.forEach((key) => {
            command[key] = device.command[key] === 1;
        });
        that.setData({
            device: device,
            command: command,
        });
        wx.setNavigationBarTitle({
            title: device.name,
        });
    },
    commandChange(e) {
        if (!this.data.device.isOnline) {
            wx.showToast({
                title: "设备离线",
                icon: "error",
                mask: true,
            });
            let command = this.data.command;
            command[e.currentTarget.dataset.switchkey] = !e.detail.value;
            this.setData({
                command: command,
            });
            return;
        }
        let command = this.data.command;
        command[e.currentTarget.dataset.switchkey] = e.detail.value;
        if (e.detail.value && e.currentTarget.dataset.switchkey != "switchMain") {
            command["switchMain"] = true;
        }
        this.setData({
            command: command,
        });
        let req = {
            device_type: "dc1",
            device_id: this.data.deviceId,
            dc1: {
                switchMain: this.data.command.switchMain ? 1 : 0,
                switch1: this.data.command.switch1 ? 1 : 0,
                switch2: this.data.command.switch2 ? 1 : 0,
                switch3: this.data.command.switch3 ? 1 : 0,
            },
        };
        postMark("/v1/command", req);
    },
    recoverChange(e) {
        const that = this;
        if (e.detail.value) {
            wx.showModal({
                title: "确认",
                content: "排插重新联网后会恢复成最后的开关状态，确定开启吗？",
                success(res) {
                    if (res.confirm) {
                        that.data.device.recover = true;
                        updateDevice({
                            deviceId: that.data.deviceId,
                            recover: true,
                        }).finally(() => {
                            getApp().getDevcieList();
                        });
                    } else {
                        that.setData({"device.recover": false});
                    }
                },
            });
        } else {
            that.data.device.recover = false;
            updateDevice({
                deviceId: that.data.deviceId,
                recover: false,
            });
        }
    },
    editSwitchName(e) {
        let switchName = e.currentTarget.dataset.switchname;
        let that = this;
        let req = {
            deviceId: that.data.deviceId,
        };
        wx.showModal({
            title: "请输入新名称",
            editable: true,
            placeholderText: this.data.device[switchName],
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
                    req[switchName] = res.content;
                    updateDevice(req).finally(() => {
                        getApp().getDevcieList();
                    });
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
