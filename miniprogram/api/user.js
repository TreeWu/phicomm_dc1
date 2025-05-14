const {postMark, getMark} = require("../utils/request");

export function systemInfo() {
    return getMark("/v1/wechat/system_info");
}

export function checkHost() {
    return getMark("/v1/check_host");
}

export function login() {
    return new Promise((resolve, reject) => {
        wx.login({
            success: (res) => {
                getMark("/v1/wechat/miniapp/jscode2session", {
                    code: res.code,
                })
                    .then((r) => {
                        wx.setStorageSync("userinfo", r);
                        resolve();
                    })
                    .catch((e) => reject(e));
            },
            fail: (err) => reject(err),
        });
    });
}

export function updateUser(data) {
    return postMark("/v1/wechat/miniapp/update_user", data);
}

export function userInfo() {
    return getMark("/v1/wechat/miniapp/user_info");
}
