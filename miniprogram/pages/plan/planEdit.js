const {getPlan, savePlan, updatePlan, delPlan} = require("../../api/index");
const {getCurrentDate, getCurrentTimeArray} = require("../../utils/date");
const PLAN_TYPE_AUTO = "PLAN_TYPE_AUTO";
const PLAN_TYPE_MANUAL = "PLAN_TYPE_MANUAL";
Page({
    data: {
        plan: {
            id: 0,
            name: "",
            planType: PLAN_TYPE_AUTO,
            cron: "",
            enabled: true,
            devices: [],
        },
        typeOptions: [
            {name: "自动执行", value: PLAN_TYPE_AUTO},
            {name: "手动执行", value: PLAN_TYPE_MANUAL},
        ],
        timeArray: [
            Array.from({length: 24}, (_, i) => i.toString().padStart(2, "0")), // 时 0-23],
            Array.from({length: 60}, (_, i) => i.toString().padStart(2, "0")), // 分 0-59
            Array.from({length: 60}, (_, i) => i.toString().padStart(2, "0")), // 秒0-59
        ],
        timeIndex: getCurrentTimeArray(),
        typeIndex: 0,
        repeatOptions: [
            {name: "执行一次", value: "once"},
            {name: "每天", value: "daily"},
            {name: "每周", value: "weekly"},
            {name: "自定义", value: "custom"},
        ],
        repeatIndex: 0, // 默认每天
        weekDays: [
            {name: "一", value: 1, selected: false},
            {name: "二", value: 2, selected: false},
            {name: "三", value: 3, selected: false},
            {name: "四", value: 4, selected: false},
            {name: "五", value: 5, selected: false},
            {name: "六", value: 6, selected: false},
            {name: "日", value: 0, selected: false},
        ],
        specificDate: getCurrentDate(),
        cronExpression: "",
        deviceMap: null,
        devices: [],
        allDevices: [],
    },

    onLoad(options) {
        if (options.id) {
            // 编辑现有计划，加载数据
            this.loadPlanData(options.id);
        } else {
            // 新建计划，生成默认CRON
            this.generateCron();
        }
        let app = getApp();
        this.setData({
            deviceMap: app.globalData.deviceMap,
            allDevices: app.globalData.devices,
        });
    },

    loadPlanData(id) {
        const that = this;
        getPlan(id)
            .then((res) => {
                var typeIndex = res.plan.planType == PLAN_TYPE_AUTO ? 0 : 1;
                that.setData({
                    plan: res.plan,
                    typeIndex: typeIndex,
                });
                that.parsePlan(res.plan);
            })
            .catch((err) => {
                wx.showToast({
                    title: err.data.message,
                    icon: "error",
                    complete: () => {
                        setTimeout(() => {
                            wx.navigateBack();
                        }, 1500);
                    },
                });
            });
    },

    typeChange(e) {
        this.setData({
            typeIndex: e.detail.value,
            "plan.planType": this.data.typeOptions[e.detail.value].value,
        });
        this.generateCron();
    },
    timeChange(e) {
        this.setData({
            execTime: e.detail.value,
        });
        this.generateCron();
    },
    repeatChange(e) {
        this.setData({
            repeatIndex: e.detail.value,
        });
        this.generateCron();
    },
    toggleWeekDay(e) {
        const index = e.currentTarget.dataset.index;
        const weekDays = this.data.weekDays.map((day, i) => {
            if (i === index) {
                return {...day, selected: !day.selected};
            }
            return day;
        });
        this.setData({
            weekDays,
        });
        this.generateCron();
    },
    dateChange(e) {
        this.setData({
            specificDate: e.detail.value,
        });
        this.generateCron();
    },
    parsePlan(plan) {
        let devices = [];
        for (let i = 0; i < plan.devices.length; i++) {
            let device = plan.devices[i];
            let d = {
                deviceId: device.deviceId,
                switchMain: false,
                switchMainSel: false,
                switch1Sel: false,
                switch1: false,
                switch2Sel: false,
                switch2: false,
                switch3Sel: false,
                switch3: false,
            };
            if (device.switchMain != null) {
                d.switchMainSel = true;
                d.switchMain = device.switchMain == 1;
            }
            if (device.switch1 != null) {
                d.switch1Sel = true;
                d.switch1 = device.switch1 == 1;
            }
            if (device.switch2 != null) {
                d.switch2Sel = true;
                d.switch2 = device.switch2 == 1;
            }
            if (device.switch3 != null) {
                d.switch3Sel = true;
                d.switch3 = device.switch3 == 1;
            }
            devices.push(d);
        }
        this.setData({
            devices: devices,
        });

        if (plan.planType != PLAN_TYPE_AUTO) {
            return;
        }
        this.setData({
            cronExpression: plan.cron,
        });
        var cron = plan.cron;
        // 解析现有的CRON表达式来初始化表单
        const parts = cron.split(" ");
        if (parts.length < 6) {
            wx.showToast({
                title: "计划解析失败",
                icon: "error",
            });
            return;
        }
        let isNumber = function (val) {
            return !isNaN(val) && !isNaN(parseInt(val));
        };
        const second = parts[0];
        const minute = parts[1];
        const hour = parts[2];
        const dayOfMonth = parts[3];
        const month = parts[4];
        const dayOfWeek = parts[5];
        // 设置时间
        if (isNumber(second) && isNumber(minute) && isNumber(hour)) {
            this.setData({
                timeIndex: [parseInt(hour), parseInt(minute), parseInt(second)],
            });
        }
        // 判断重复模式
        if (
            isNumber(second) &&
            isNumber(minute) &&
            isNumber(hour) &&
            isNumber(dayOfMonth) &&
            isNumber(month) &&
            dayOfWeek == "*"
        ) {
            // 指定日期
            const now = new Date();
            const year = now.getFullYear();
            this.setData({
                repeatIndex: 0,
                specificDate: `${year}-${month}-${dayOfMonth}`,
            });
        } else if (
            isNumber(second) &&
            isNumber(minute) &&
            isNumber(hour) &&
            !isNumber(dayOfMonth) &&
            !isNumber(month) &&
            dayOfWeek != "*"
        ) {
            // 每周
            const selectedDays = dayOfWeek.split(",").map((d) => parseInt(d));
            const weekDays = this.data.weekDays.map((day) => {
                return {...day, selected: selectedDays.includes(day.value)};
            });
            this.setData({
                repeatIndex: 2,
                weekDays,
            });
        } else if (
            isNumber(second) &&
            isNumber(minute) &&
            isNumber(hour) &&
            dayOfMonth == "*" &&
            month == "*" &&
            (dayOfWeek == "*" || dayOfWeek == "?")
        ) {
            // 每天
            this.setData({
                repeatIndex: 1,
            });
        } else {
            this.setData({
                repeatIndex: 3,
            });
        }
    },
    generateCron() {
        let cron = "";
        let second = this.data.timeIndex[2];
        let minute = this.data.timeIndex[1];
        let hour = this.data.timeIndex[0];
        switch (parseInt(this.data.repeatIndex)) {
            case 0: // 执行一次
                const [year, month, day] = this.data.specificDate.split("-");
                cron = `${second} ${minute} ${hour} ${day} ${month} ?`;
                break;
            case 1: // 每天
                cron = `${second} ${minute} ${hour} * * ?`;
                break;
            case 2: // 每周
                const selectedDays = this.data.weekDays
                    .filter((day) => day.selected)
                    .map((day) => day.value)
                    .join(",");
                cron = `${second} ${minute} ${hour} * * ${selectedDays || "*"}`;
                break;
            case 3: // 自定义
                cron = this.data.cronExpression;
        }
        this.setData({
            cronExpression: cron,
        });
    },
    savePlan(e) {
        const formData = e.detail.value;
        const plan = {
            ...this.data.plan,
            name: formData.name,
            planType: this.data.typeOptions[this.data.typeIndex].value,
            cron: formData.cronExpression,
        };
        if (plan.name == "") {
            wx.showToast({
                title: "计划名称不能为空",
                icon: "error",
            });
            return;
        }
        if (this.data.devices.length == 0) {
            wx.showToast({
                title: "设备为空",
                icon: "error",
            });
            return;
        }
        plan.devices = this.data.devices.map((item) => {
            let device = {
                deviceId: item.deviceId,
            };
            if (item.switchMainSel) {
                device.switchMain = item.switchMain ? 1 : 0;
            }
            if (item.switch1Sel) {
                device.switch1 = item.switch1 ? 1 : 0;
            }
            if (item.switch2Sel) {
                device.switch2 = item.switch2 ? 1 : 0;
            }
            if (item.switch3Sel) {
                device.switch3 = item.switch3 ? 1 : 0;
            }
            return device;
        });

        let f = updatePlan;
        if (this.data.plan.id == 0) {
            f = savePlan;
        }
        f({plan: plan})
            .then((res) => {
                setTimeout(() => {
                    wx.showToast({
                        title: "保存成功",
                        icon: "success",
                        duration: 1500,
                        complete: () => {
                            setTimeout(() => {
                                wx.navigateBack();
                            }, 1500);
                        },
                    });
                }, 50);
            })
            .catch((err) => {
                setTimeout(() => {
                    wx.showToast({
                        title: err.data.message,
                        icon: "error",
                    });
                }, 50);
            });
    },
    bindMultiPickerChange: function (e) {
        this.setData({
            timeIndex: e.detail.value,
        });
        this.generateCron();
    },
    addDevcie() {
        const that = this;
        wx.showActionSheet({
            itemList: this.data.allDevices.map((item) => item.name),
            success(res) {
                let d = that.data.allDevices[res.tapIndex];
                if (
                    that.data.devices.filter((item) => item.deviceId == d.deviceId)
                        .length != 0
                ) {
                    wx.showToast({
                        title: "设备已存在",
                        icon: "error",
                    });
                    return;
                }
                that.data.devices.push({
                    deviceId: d.deviceId,
                    switchMain: false,
                    switchMainSel: false,
                    switch1Sel: false,
                    switch1: false,
                    switch2Sel: false,
                    switch2: false,
                    switch3Sel: false,
                    switch3: false,
                });
                console.log(that.data.devices);
                that.setData({
                    devices: that.data.devices,
                });
            },
            fail(res) {
                console.log(res);
            },
        });
    },
    switchSelChange(e) {
        let devices = this.data.devices.map((item) => {
            if (item.deviceId == e.currentTarget.dataset.deviceId) {
                item.switchMainSel = false;
                item.switch1Sel = false;
                item.switch2Sel = false;
                item.switch3Sel = false;
                e.detail.value.forEach((element) => {
                    item[element] = true;
                });
            }
            return item;
        });
        this.setData({
            devices: devices,
        });
    },
    switchChange(e) {
        let devices = this.data.devices.map((item) => {
            if (item.deviceId == e.currentTarget.dataset.deviceId) {
                item[e.currentTarget.dataset.switchKey] = e.detail.value;
            }
            return item;
        });
        this.setData({
            devices: devices,
        });
    },
    delDevice(e) {
        let devices = this.data.devices.filter((item) => {
            return item.deviceId != e.currentTarget.dataset.deviceId;
        });
        this.setData({
            devices: devices,
        });
    },
    async delPlan(e) {
        try {
            await delPlan({id: this.data.plan.id});
            wx.navigateBack();
        } catch (err) {
            setTimeout(() => {
                wx.showToast({
                    title: err.data.message,
                    icon: "error",
                    complete: () => {
                        setTimeout(() => {
                            wx.navigateBack();
                        }, 1500);
                    },
                });
            }, 50);
        }
    },
});
