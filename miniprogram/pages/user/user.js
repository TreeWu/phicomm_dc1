const {
    login,
    updateUser,
    userInfo
} = require("../../api/index")
Page({
    data: {
        userInfo: {},
        hasUserInfo: false,
    },
    onShow() {
        const userInfo = wx.getStorageSync('userinfo')
        if (userInfo) {
            this.setData({
                userInfo: userInfo,
                hasUserInfo: true,
            })
        } else {
            this.setData({
                userInfo: {},
                hasUserInfo: false
            })
        }
    },
    onChooseAvatar(e) {
        const {
            avatarUrl
        } = e.detail
        let userInfo = this.data.userInfo
        userInfo.avatar = avatarUrl
        this.setData({
            userInfo: userInfo
        })
    },
    login() {
        const that = this
        login().then(() => {
            const userInfo = wx.getStorageSync('userinfo')
            that.setData({
                userInfo: userInfo,
                hasUserInfo: true,
            })
        })
    },
    update() {
        const that = this
        updateUser({
            avatar: this.data.userInfo.avatar,
            nickname: this.data.userInfo.nickname
        }).catch(err => {
            wx.showToast({
                title: err,
            })
        }).finally(() => {
            that.freshUserInfo()
        })
    },
    handleInput(e) {
        const {
            value
        } = e.detail;
        let userInfo = this.data.userInfo
        userInfo.nickname = value
        this.setData({
            userInfo: userInfo,
        })
    },
    freshUserInfo() {
        const that = this
        userInfo().then(res => {
            const userInfo = wx.getStorageSync('userinfo')
            userInfo.avatar = res.avatar
            userInfo.nickname = res.nickname
            that.setData({
                userInfo: userInfo,
                hasUserInfo: true,
            })
            wx.setStorageSync('userinfo', userInfo)
        })
    },
    handelShare() {
        wx.navigateTo({
            url: "/pages/share/share"
        })
    }
})