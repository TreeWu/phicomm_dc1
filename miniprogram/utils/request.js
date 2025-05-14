const {baseUrl} = require("../utils/config");
const request = (url, method, data, header = {}) => {
    let h = {
        "Content-Type": "application/json", // 设置请求的 header
        ...header,
    };
    const userInfo = wx.getStorageSync("userinfo");
    if (userInfo && userInfo.token) {
        const token = userInfo.token.accessToken;
        if (token) {
            h["Authorization"] = "Bearer " + token;
        }
    }
    return new Promise((resolve, reject) => {
        wx.request({
            url: baseUrl + url, // 拼接完整的url
            method: method,
            data: data,
            header: h,
            success(res) {
                if (res.statusCode >= 200 && res.statusCode < 300) {
                    // 请求成功的处理
                    resolve(res.data);
                } else {
                    if (res.statusCode == 401) {
                        wx.showToast({
                            title: "请先登录",
                            icon: "error",
                        });
                        wx.clearStorageSync("userinfo");
                        setTimeout(() => {
                            wx.switchTab({
                                url: "/pages/user/user",
                            });
                        }, 2000);
                    }
                    reject(res);
                }
            },
            fail: reject,
        });
    });
};

// 导出封装的request方法
module.exports = {
    get: (url, data, header) => request(url, "GET", data, header),
    post: (url, data, header) => request(url, "POST", data, header),
    postMark: async (url, data, header) => {
        try {
            wx.showLoading({mask: true, title: "请求中"});
            return await request(url, "POST", data, header);
        } finally {
            wx.hideLoading({noConflict: false});
        }
    },
    getMark: async (url, data, header) => {
        wx.showLoading({mask: true, title: "请求中"});
        try {
            return await request(url, "GET", data, header);
        } finally {
            wx.hideLoading();
        }
    },
};
