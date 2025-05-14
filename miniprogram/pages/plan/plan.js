// pages/plan/plan.js
const {listPlan, switchPlan, execPlan} = require("../../api/index");
Page({
    /**
     * 页面的初始数据
     */
    data: {
        plans: [],
    },

    /**
     * 生命周期函数--监听页面显示
     */
    onShow() {
        this.getPlanList();
    },
    async getPlanList() {
        let data = await listPlan();
        this.setData({
            plans: data.plans,
        });
    },
    async chanagePlan(e) {
        const that = this;
        try {
            const planId = e.currentTarget.dataset.id;
            const isChecked = e.detail.value;
            var data = {id: planId, enable: isChecked};
            await switchPlan(data);
            setTimeout(() => {
                wx.showToast({
                    title: isChecked ? "计划已启用" : "计划已禁用",
                    icon: "success",
                    mask: true,
                });
            }, 50);
        } catch (error) {
            setTimeout(() => {
                wx.showToast({
                    title: error.data.message,
                    icon: "error",
                    mask: true,
                });
            }, 50);
        } finally {
            that.getPlanList();
        }
    },
    async navigateToEdit(e) {
        const planId = e.currentTarget.dataset.id;
        wx.navigateTo({
            url: `/pages/plan/planEdit?id=${planId}`,
        });
    },
    addNewPlan(e) {
        wx.navigateTo({
            url: `/pages/plan/planEdit`,
        });
    },
    async execPlan(e) {
        try {
            const plangId = e.currentTarget.dataset.id;
            let res = await execPlan({id: plangId});
            setTimeout(() => {
                wx.showToast({
                    title: res.message,
                    mask: true,
                });
            }, 50);
        } catch (err) {
            setTimeout(() => {
                wx.showToast({
                    title: err.data.message,
                    icon: "error",
                    mask: true,
                });
            }, 50);
        }
    },
    async onPullDownRefresh() {
        await this.getPlanList();
        wx.stopPullDownRefresh();
    },
});
