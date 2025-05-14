//"2023-04-23"
export function getCurrentDate() {
    const now = new Date();
    const year = now.getFullYear();
    const month = String(now.getMonth() + 1).padStart(2, '0'); // 月份从0开始，所以要+1
    const day = String(now.getDate()).padStart(2, '0');

    return `${year}-${month}-${day}`;
}

//"14:30:00"
export function getCurrentTime() {
    const now = new Date();
    const second = String(now.getSeconds()).padStart(2, '0'); // 补零，如 "09"
    const hours = String(now.getHours()).padStart(2, '0'); // 补零，如 "09"
    const minutes = String(now.getMinutes()).padStart(2, '0'); // 补零，如 "05"
    return `${hours}:${minutes}:${second}`;
}

//"[12,23,4]"
export function getCurrentTimeArray() {
    const now = new Date();
    return [
        now.getHours(), now.getMinutes(), now.getSeconds()
    ]
}