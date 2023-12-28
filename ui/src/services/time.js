


const formatTime = (dateTimeString) => {

    // 将时间字符串转换为 Date 对象
    const dateTime = new Date(dateTimeString);
    // 获取年、月、日、小时、分钟和秒
    const year = dateTime.getFullYear();
    const month = ("0" + (dateTime.getMonth() + 1)).slice(-2);
    const day = ("0" + dateTime.getDate()).slice(-2);
    const hours = ("0" + dateTime.getHours()).slice(-2);
    const minutes = ("0" + dateTime.getMinutes()).slice(-2);
    const seconds = ("0" + dateTime.getSeconds()).slice(-2);

    // 格式化为中国地区更易读的格式
    const formattedDateTime = `${year}-${month}-${day} ${hours}:${minutes}:${seconds}`;
    return formattedDateTime;
}



export {
    formatTime
}