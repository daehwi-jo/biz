function excelDownload(apiUrl, Name, excelHeader) { //excelHeader = [목록에 들어갈 테이블 헤더]
    RequestGETApi(apiUrl).then(function (data) {
        if (data.resultCode === "00") {
            let blob = new Blob([getExcelHtml(Name, data.resultData.listData, excelHeader)], {type: "application/csv;charset=utf-8;"});
            let downloadLink = document.createElement('a');
            downloadLink.href = URL.createObjectURL(blob);
            downloadLink.download = Name + ".xls";
            document.body.appendChild(downloadLink);
            downloadLink.click();
            document.body.removeChild(downloadLink);
        }
    })
}

function getExcelHtml(Name, data, excelHeader) {

    let htmlString = '<html xmlns:x="urn:schemas-microsoft-com:office:excel">'
    htmlString += '<head><meta http-equiv="content-type" content="application/vnd.ms-excel; charset=UTF-8"><xml>'
    htmlString += '<x:ExcelWorkbook><x:ExcelWorksheets><x:ExcelWorksheet>'
    htmlString += '<x:Name>' + Name + '</x:Name>'
    htmlString += '<x:WorksheetOptions><x:Panes></x:Panes></x:WorksheetOptions>'
    htmlString += '</x:ExcelWorksheet></x:ExcelWorksheets></x:ExcelWorkbook>'
    htmlString += '</xml></head>'
    htmlString += '<body>' + getExcelTable(data, excelHeader).outerHTML + '</body>'
    htmlString += '</html>'
    return htmlString

}

function getExcelTable(data, excelHeader) {

    let ExcelTable = document.createElement('table')
    let ExcelHead = document.createElement('thead')
    let ExcelBody = document.createElement('tbody')

    //헤드
    let headRow = document.createElement('tr')
    for (let i = 0; i < excelHeader.length; i++) {
        let cell = document.createElement('th')
        //이곳에서 각 행의 넓이 지정
        //각행의 높이 는 각 행에서 높이를 지정해 줘야함
        cell.style.textAlign = 'center'
        cell.style.width = '100px'
        console.log(cell.outerHTML)
        cell.appendChild(document.createTextNode(excelHeader[i])) //이 부분에 데이터 삽입
        headRow.appendChild(cell)
    }
    ExcelHead.appendChild(headRow)

    //바디
    for (let i = 0; i < data.length; i++) {
        let bodyRow = document.createElement('tr')
        for (let j = 1; j < excelHeader.length + 1; j++) {
            //형식 지정에 따라 데이터 형식 편집 "mso-number-format:ㅁㅁㅁ;" ㅁㅁㅁ이곳에 엑셀에서 사용할 형식기입;
            let cell = (j===8||j===7)
                    ? "<td style=\"mso-number-format:'0';\">"+data[i]["excel" + j]+"</td>"
                    : "<td style=\"mso-number-format:'\@'\">"+data[i]["excel" + j]+"</td>"
            bodyRow.innerHTML += cell
        }
        ExcelBody.appendChild(bodyRow)
    }

    ExcelTable.appendChild(ExcelHead)
    ExcelTable.appendChild(ExcelBody)

    return ExcelTable
}