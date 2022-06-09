class AlertModal {

    constructor(title, content) {
        if (this.addModal()) {
            this.setModalContent(title, content)
        }
    }

    addModal() { //해당객체 생성시 자동으로 알람 모달 추가됨
        let modalHtmlString = "<div class=\"col-md-4\">"
        modalHtmlString += "<div class=\"modal fade\" id=\"AlertModal\" tabindex=\"-2\" role=\"dialog\" aria-hidden=\"true\">"
        modalHtmlString += "<div class=\"modal-dialog modal-dialog-centered\" role=\"document\">"
        modalHtmlString += "<div class=\"modal-content\">"
        modalHtmlString += "<div class=\"modal-header\" id='modalHead'><h5 class=\"modal-title\" id='AlertModalHead'></h5></div>"
        modalHtmlString += "<div class=\"modal-body\" id='modalBody'><p class=\"mb-0\" id='AlertModalContent'></p></div>"
        modalHtmlString += "<div class=\"modal-footer\" id='modalFoot'>"
        modalHtmlString += "<button type=\"button\" " +
            "class=\"btn bg-gradient-primary\" data-bs-dismiss=\"modal\" id=\"AlertModalCommit\" onclick=''>확인" +
            "</button></div>"
        modalHtmlString += "</div></div></div></div>"

        document.body.innerHTML += modalHtmlString

        return true
    }

    setModalContent(title, content) {
        document.getElementById("AlertModalHead").textContent = title
        document.getElementById("AlertModalContent").textContent = content
    }
}