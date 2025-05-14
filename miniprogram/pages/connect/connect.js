// pages/connect/connect.js
const AB2String = (arrayBuffer) => {
    const unit8Arr = new Uint8Array(arrayBuffer);
    const encodedString = String.fromCharCode.apply(null, unit8Arr);
    const decodedString = decodeURIComponent(escape(encodedString)); // 没有这一步中文会乱码
    return decodedString;
};
const {bindDevice, checkHost} = require("../../api/index");
Page({
    /**
     * 页面的初始数据
     */
    data: {
        respTimeout: 10000, // 等待配网超时时间
        send_data: {},
        ssid: "",
        password: "",
        address: "192.168.4.1",
        port: "7550",
        canUseUdp: false,
        udpInstance: null,
        message: "",
        host: "",
        resp: {},
    },
    onLoad() {
        const canIUse = wx.canIUse("createUDPSocket");
        if (!canIUse) {
            wx.showModal({
                title: "微信版本过低，暂不支持本功能",
            });
            this.setData({
                canUseUdp: canIUse,
            });
        }
    },
    onShow() {
        let data = getApp().globalData;
        this.setData({host: data.systemInfo.host});
    },
    handleInput(e) {
        let value = e.detail.value;
        let item = e.target.dataset.item;
        var obj = this.data;
        obj[item] = value;
        this.setData(obj);
    },
    bindDevice() {
        let that = this;
        if (this.data.resp.status == 200) {
            bindDevice({
                mac: this.data.resp.result.mac,
                uuid: this.data.resp.uuid,
            })
                .then((res) => {
                    wx.switchTab({
                        url: "/pages/devices/devices",
                    });
                })
                .catch((e) => {
                    wx.showModal({
                        showCancel: false,
                        content: e.data.message,
                    });
                });
        } else {
            wx.showModal({
                title: "绑定失败",
                content: "请先对排插进行联网再绑定",
                showCancel: false,
            });
        }
    },
    connect() {
        const that = this;
        that.close();
        wx.showModal({
            title: "注意",
            content: "请先连接[PHI_PLUG1]开头的WIFI",
            confirmText: "已连接",
            complete: (res) => {
                if (res.confirm) {
                    that.setData({
                        resp: {},
                    });
                    let c = setTimeout(() => {
                        if (!that.data.resp.status) {
                            wx.hideLoading();
                            that.close();
                            wx.showModal({
                                title: "联网失败",
                                content: "连接超时",
                                showCancel: false,
                            });
                        }
                    }, that.data.respTimeout);
                    wx.showLoading({
                        title: "联网中...",
                        mask: true,
                    });
                    let udpInstance = wx.createUDPSocket();
                    udpInstance.onListening(() => {
                        udpInstance.setTTL(10);
                    });
                    udpInstance.onMessage((res) => {
                        clearTimeout(c);
                        wx.hideLoading();
                        //{"header":"phi-plug-0001","uuid":"00010","status":200,"msg":"set wifi success","result":{"mac":"a4:7b:9d:06:a0:e6"}}
                        let message = AB2String(res.message);
                        let resp = JSON.parse(message);
                        that.setData({
                            resp: resp,
                        });
                        if (resp.status == 200) {
                            setTimeout(() => {
                                wx.showToast({
                                    title: "联网成功，尽快完成绑定",
                                });
                            }, 50);
                        }
                    });
                    udpInstance.onError((errMsg) => {
                        console.log(errMsg);
                        wx.hideLoading();
                        wx.showModal({
                            title: "联网失败",
                            content: errMsg,
                            showCancel: false,
                        });
                    });
                    udpInstance.onClose(() => {
                        console.log("onClose");
                        wx.hideLoading();
                    });
                    let port = udpInstance.bind();
                    let data = {
                        header: "phi-plug-0001",
                        uuid: "" + port,
                        action: "wifi=",
                        auth: "" + port,
                        params: {
                            ssid: that.data.ssid,
                            password: that.data.password,
                        },
                    };
                    that.setData({
                        send_data: data,
                        udpInstance: udpInstance,
                    });
                    udpInstance.send({
                        address: that.data.address,
                        port: that.data.port,
                        message: JSON.stringify(data) + "\n",
                    });
                }
            },
        });
    },
    close() {
        if (this.data.udpInstance != null) {
            this.data.udpInstance.close();
            this.data.udpInstance = null;
        }
    },
    async checkHost() {
        try {
            let res = await checkHost();
            wx.showModal({
                content:
                    getApp().globalData.systemInfo.host == res.host
                        ? "劫持成功"
                        : "劫持失败",
            });
        } catch (e) {
            wx.showModal({
                content: "检查失败",
            });
        }
    },
});
