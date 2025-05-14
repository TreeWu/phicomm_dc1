//app.js
const {userInfo, systemInfo, devcieList} = require("api/index");
App({
    data: {
        timer: null,
        flushInterval: 0,
        hasload: false,
    },
    onShow: function () {
        this.run();
    },
    run: async function () {
        try {
            if (this.data.timer != null) {
                clearInterval(this.data.timer);
            }
            let res = await systemInfo();
            this.globalData.systemInfo = res;
            this.data.flushInterval = res.flushInterval;
            if (res.flushInterval != 0) {
                this.data.timer = setInterval(
                    this.getDevcieList,
                    this.data.flushInterval
                );
            }
            if (!this.data.hasload) {
                let user = await userInfo();
                const userinfo = wx.getStorageSync("userinfo");
                userinfo.avatar = user.avatar;
                userinfo.nickname = user.nickname;
                wx.setStorageSync("userinfo", userinfo);
                this.data.hasload = true;
            }
        } catch (err) {
            console.log(err)
            setTimeout(() => {
                wx.showToast({
                    title: err.data.message,
                    icon: "error",
                });
            }, 50);
        }
    },
    stop: function () {
        if (this.data.timer != null) {
            clearInterval(this.data.timer);
            this.data.timer = null;
        }
    },
    onHide: function () {
        this.stop();
    },
    getDevcieList() {
        const self = this;
        devcieList()
            .then((data) => {
                const devcies = data.devices;
                let deviceMap = {};
                devcies.forEach((item) => {
                    deviceMap[item.deviceId] = item;
                });
                self.globalData.deviceMap = deviceMap;
                self.globalData.devices = devcies;
            })
            .catch((err) => {
                clearInterval(this.data.timer);
                this.data.timer = null;
            });
    },
    unwatch: function (wathchId) {
        for (let key in this.globalData.listenKey) {
            if (this.globalData.listenKey[key][wathchId]) {
                delete this.globalData.listenKey[key][wathchId];
            }
        }
    },
    watch: function (key, method) {
        var obj = this.globalData;
        if (obj.listenKey[key]) {
            let listenId = obj.listenId++;
            obj.listenKey[key][listenId] = method;
            return listenId;
        }
        let listenId = obj.listenId++;
        obj.listenKey[key] = {};
        obj.listenKey[key][listenId] = method;

        Object.defineProperty(obj, key, {
            configurable: true,
            enumerable: true,
            set: function (value) {
                if (this.listenKey[key]) {
                    for (let item of Object.values(this.listenKey[key])) {
                        item(value);
                    }
                }
                this["_" + key] = value;
            },
            get: function () {
                return this["_" + key];
            },
        });
        return listenId;
    },
    globalData: {
        deviceMap: {},
        devices: [],
        listenKey: {},
        listenId: 1,
        systemInfo: {},
    },
});
