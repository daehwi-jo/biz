function getContentNumber(pageNo, contentNum) {
    if (pageNo * 1 > 1) {
        return (pageNo - 1) * (contentNum || 10)
    } else {
        return 0
    }
}

function getPagination(totalData, currentPage, dataPerPage, fname) {

    const pageCount = 5;
    const totalPage = Math.ceil(totalData / dataPerPage);    // 총 페이지 수

    const pageGroup = Math.ceil(currentPage / pageCount);    // 페이지 그룹

    let last = pageGroup * pageCount;    // 화면에 보여질 마지막 페이지 번호
    if (last > totalPage)
        last = totalPage;
    let first = last - (pageCount - 1);    // 화면에 보여질 첫번째 페이지 번호
    const next = last + 1;
    const prev = first - 1;
    let pHtml = "<ul class='pagination justify-content-center mb-0' style='cursor: pointer'>";
    pHtml += "<li class='page-item'><a class='page-link first' onclick=\"" + fname + "(" + 1 + ");\"> <i class='fa fa-angle-double-left'></i></a></li>";
    if (first > 5) {
        pHtml += "<li class='page-item'><a class ='page-link prev' onclick=\"" + fname + "(" + prev + ");\" ><i class='fa fa-angle-left' aria-hidden='true'></i></a> </li>";
    }
    for (let j = first; j <= last; j++) {
        if (currentPage === (j)) {
            pHtml += "<li class='page-item active'><a class='page-link'  onclick=\"" + fname + "(" + j + ");\">" + j + "</a></li>";
        } else if (j > 0) {
            pHtml += "<li class='page-item'><a class='page-link'  onclick=\"" + fname + "(" + j + ");\">" + j + "</a></li>";
        }
    }

    if (totalPage > next) {
        pHtml += "<li class='page-item'> <a class='page-link last' onclick=\"" + fname + "(" + next + ");\"> <i class='fa fa-angle-right' aria-hidden='true'></i> </a></li>";
    }

    if (next > 5 && next < totalPage) {
        pHtml += "<li class='page-item'><a class='page-link next' onclick=\"" + fname + "(" + totalPage + ");\" aria-label='Next'> <i class='fa fa-angle-double-right'></i></a></li>";
    }
    pHtml += "</ul>"

    document.getElementById("pageNav").innerHTML = pHtml
}

function getPaginationName(totalData, currentPage, dataPerPage, fname, name) {

    const pageCount = 5;
    const totalPage = Math.ceil(totalData / dataPerPage);    // 총 페이지 수

    const pageGroup = Math.ceil(currentPage / pageCount);    // 페이지 그룹

    let last = pageGroup * pageCount;    // 화면에 보여질 마지막 페이지 번호
    if (last > totalPage)
        last = totalPage;
    let first = last - (pageCount - 1);    // 화면에 보여질 첫번째 페이지 번호
    const next = last + 1;
    const prev = first - 1;
    let pHtml = "<ul class='pagination justify-content-center mb-0' style='cursor: pointer'>";
    pHtml += "<li class='page-item'><a class='page-link first' onclick=\"" + fname + "(" + 1 + ");\"> <i class='fa fa-angle-double-left'></i></a></li>";
    if (first > 5) {
        pHtml += "<li class='page-item'><a class ='page-link prev' onclick=\"" + fname + "(" + prev + ");\" ><i class='fa fa-angle-left' aria-hidden='true'></i></a> </li>";
    }
    for (let j = first; j <= last; j++) {
        if (currentPage === (j)) {
            pHtml += "<li class='page-item active'><a class='page-link'  onclick=\"" + fname + "(" + j + ");\">" + j + "</a></li>";
        } else if (j > 0) {
            pHtml += "<li class='page-item'><a class='page-link'  onclick=\"" + fname + "(" + j + ");\">" + j + "</a></li>";
        }
    }

    if (totalPage > next) {
        pHtml += "<li class='page-item'> <a class='page-link last' onclick=\"" + fname + "(" + next + ");\"> <i class='fa fa-angle-right' aria-hidden='true'></i> </a></li>";
    }

    if (next > 5 && next < totalPage) {
        pHtml += "<li class='page-item'><a class='page-link next' onclick=\"" + fname + "(" + totalPage + ");\" aria-label='Next'> <i class='fa fa-angle-double-right'></i></a></li>";
    }
    pHtml += "</ul>"

    document.getElementById(name).innerHTML = pHtml
}

