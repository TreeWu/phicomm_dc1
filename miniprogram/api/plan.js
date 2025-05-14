const {
    postMark,
    getMark,
    get
} = require("../utils/request")

export function listPlan() {
    return getMark("/v1/plan/list")
}

export function switchPlan(data) {
    return postMark("/v1/plan/switch", data)
}

export function getPlan(id) {
    return getMark("/v1/plan/get?id=" + id)
}

export function savePlan(data) {
    return postMark("/v1/plan/create", data)
}

export function updatePlan(data) {
    return postMark("/v1/plan/update", data)
}

export function delPlan(data) {
    return postMark("/v1/plan/delete", data)
}

export function execPlan(data) {
    return postMark("/v1/plan/exec", data)
}