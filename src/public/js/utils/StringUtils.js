function numberWithCommas(str) {
    if (str === "%!s(<nil>)") {
        return 0
    }
    return str.toString().replace(/\B(?=(\d{3})+(?!\d))/g, ',');
}

function textNumComma(obj){
    obj.value = numberWithCommas(String((obj.value).replace(/[^0-9]/g, "")).replaceAll(",",""))
}

function textNumUnComma(obj){
    return obj.value.replaceAll(",","")
}

function getHpNo(hpNo) {
    const slice = (number,num1,num2,num3) =>
    { return hpNo.slice(0, num1) + '-' + hpNo.slice(num1, num1 + num2) + '-' + hpNo.slice(num1 + num2, num1 + num2 + num3) }

    hpNo = hpNo.replaceAll("-", "")

    const cases = {
        9:slice(hpNo, 2, 3, 4),
        10:(hpNo.slice(0, 2) === "02")
            ?slice(hpNo, 2, 4, 4)
            :slice(hpNo, 3, 3, 4),
        11: slice(hpNo, 3, 4, 4)
    }

    return cases[hpNo.length] || hpNo
}

function textNumHpNum(obj){
    obj.value = getHpNo(String((obj.value).replace(/[^0-9]/g, "")).replaceAll("-",""))
}

function Format_comma2(str) {
    str = String(str);
    var minus = str.substring(0, 1);

    str = str.replace(/[^\d]+/g, '');
    str = str.replace(/(\d)(?=(?:\d{3})+(?!\d))/g, '$1,');

    //음수일 경우
    if (minus == "-") str = "-" + str;

    return str;
}

function maxLengthCheck(object) { //number type input maxLength
    if (object.value.length > object.maxLength) {
        object.value = object.value.slice(0, object.maxLength);
    }
}

function checkUrlForm(strUrl) {
    let expUrl = /^http(s)?\:\/\//i;

    return expUrl.test(strUrl);
}

