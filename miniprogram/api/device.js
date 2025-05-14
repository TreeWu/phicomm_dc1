const {
    postMark,
    get
} = require("../utils/request")

export function updateDevice(data) {
    return postMark("/v1/device", data)
}

export function bindDevice(data) {
    return postMark("/v1/device/binding", data)
}

export function devcieList() {
    return get("/v1/device")
}

export function getShareList() {
    return get("/v1/share/list")
}

export function shareConfirm(data) {
    return postMark("/v1/share/confirm", data)
}

export function shareRevoke(data) {
    return postMark("/v1/share/revoke", data)
}

export function shareInvite(data) {
    return postMark("/v1/share/invite", data)
}
