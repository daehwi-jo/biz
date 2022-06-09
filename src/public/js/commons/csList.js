
function addCsList(type, KeyId, csUrl) {
    document.getElementById("csList").innerHTML = "<div class=\"col ms-3 p-3\"><h5>CS내역</h5>"
    + "<div class=\"row mt-2 me-3\"><div class=\"input-group mb-3\">"
    + "<textarea type=\"text\" class=\"form-control\" placeholder=\"CS 내용 입력\" id=\"csInputData\" style='resize: none;'></textarea>"
    + "<button class=\"btn btn-outline-primary mb-0\" type=\"submit\" id=\"button-addon2\" " +
                            "data-bs-toggle=\"modal\" data-bs-target=\"#CsAlertModal\" " +
                            "onClick='addCsContentWrite("+"\""+csUrl+"\""+","+"\""+type+"\""+","+"\""+KeyId+"\""+")'>내용 추가</button>"

    + "</div></div>"
    + "<div class=\"container-md-2 mt-2 me-3 addCsContainer\" id=\"CsList\"></div>"


    document.getElementById("CsAlert").innerHTML += "<div class=\"modal fade\" id=\"CsAlertModal\" tabindex=\"-1\" role=\"dialog\">" +
        "            <div class=\"modal-dialog modal-dialog-centered\" role=\"document\">" +
        "                <div class=\"modal-content\">" +
        "                    <div class=\"modal-header\">" +
        "                        <h5 class=\"modal-title\" id=\"AlertModalTitle\">내용 추가</h5>" +
        "                        <button type=\"button\" class=\"btn-close\" data-bs-dismiss=\"modal\" aria-label=\"Close\">" +
        "                            <span aria-hidden=\"true\">×</span>" +
        "                        </button>" +
        "                    </div>" +
        "                    <div class=\"modal-body\">" +
        "                        <p class=\"mb-0 ms-3\" id=\"modalContent\">CS 내역이 추가되었습니다.</p>" +
        "                    </div>" +
        "                    <div class=\"modal-footer\">" +
        "                        <button type=\"button\" class=\"btn bg-gradient-primary\" data-bs-dismiss=\"modal\"" +
        "                                id=\"AlertModalCommit\">확인" +
        "                        </button>" +
        "                    </div>" +
        "                </div>" +
        "            </div>" +
        "        </div>"

    return true
}

function addCsList2(type, KeyId, csUrl) {

    document.getElementById("CsAlert").innerHTML += "<div class=\"modal fade\" id=\"CsAlertModal\" tabindex=\"-1\" role=\"dialog\">" +
        "            <div class=\"modal-dialog modal-dialog-centered\" role=\"document\">" +
        "                <div class=\"modal-content\">" +
        "                    <div class=\"modal-header\">" +
        "                        <h5 class=\"modal-title\" id=\"AlertModalTitle\">내용 추가</h5>" +
        "                        <button type=\"button\" class=\"btn-close\" data-bs-dismiss=\"modal\" aria-label=\"Close\">" +
        "                            <span aria-hidden=\"true\">×</span>" +
        "                        </button>" +
        "                    </div>" +
        "                    <div class=\"modal-body\">" +
        "                        <p class=\"mb-0 ms-3\" id=\"modalContent\">CS 내역이 추가되었습니다.</p>" +
        "                    </div>" +
        "                    <div class=\"modal-footer\">" +
        "                        <button type=\"button\" class=\"btn bg-gradient-primary\" data-bs-dismiss=\"modal\"" +
        "                                id=\"AlertModalCommit\">확인" +
        "                        </button>" +
        "                    </div>" +
        "                </div>" +
        "            </div>" +
        "        </div>"

    return "<div class=\"col ms-3 p-3\"><h5>CS내역</h5>"
        + "<div class=\"row mt-2 me-3\"><div class=\"input-group mb-3\">"
        + "<textarea type=\"text\" class=\"form-control\" placeholder=\"CS 내용 입력\" id=\"csInputData\" style='resize: none;'></textarea>"
        + "<button class=\"btn btn-outline-primary mb-0\" type=\"submit\" id=\"button-addon2\" " +
        "data-bs-toggle=\"modal\" data-bs-target=\"#CsAlertModal\" " +
        "onClick='addCsContentWrite("+"\""+csUrl+"\""+","+"\""+type+"\""+","+"\""+KeyId+"\""+")'>내용 추가</button>"

        + "</div></div>"
        + "<div class=\"container-md-2 mt-2 me-3 addCsContainer\" id=\"CsList\"></div>"
}


function addCsContentWrite(csUrl, type, keyId) {
    RequestPOSTApi("/api/AddBizCsContent", {
        contents: document.getElementById("csInputData").value,
        keyId: keyId,
        type: type,
    }).then(function (data) {
        if (data.resultCode === "00") {
            RequestGETApi(csUrl + '?searchKeyId=' + keyId + '&searchType=' + type
            ).then(function (data) {
                setCsList(data.resultData)
            })
        }
    })
    document.getElementById("csInputData").value = ""
}

function setCsList(data) {

    let htmlString = ""

    if (data.length === 0) {
        htmlString += "<div class=\"row\">"
        htmlString += "     <div class=\"col\">"
        htmlString += "         <p class=\"text-center text-bold \" style='line-height: 100px;'>이력 없음</p>"
        htmlString += "     </div>"
        htmlString += "</div>"
        document.getElementById('CsList').innerHTML = htmlString
        return
    }

    for (const csData of data) {
        htmlString += "<div class=\"row\">"
        htmlString += "     <div class=\"col\">"
        htmlString += "         <p class=\"text-bold text-sm ms-2\">" + csData.contents.replaceAll(/\n/g, '<br/>') + "</p>"
        htmlString += "     </div>"
        htmlString += "     <div class=\"col addCsDate\">"
        htmlString += "         <p class=\"text-bold text-xs mt-1\">" + csData.regDate + "</p>"
        htmlString += "     </div>"
        htmlString += "</div>"
    }
    document.getElementById('CsList').innerHTML = htmlString

}