async function RequestGETApi(url) { //SELECT
    return fetch(url)
        .then(async function (response) {
            return response.json();
        })
}

async function RequestPOSTApi(url, object, header) { //INSERT
    return fetch(url, {
        method: "POST",
        headers : header || {'Content-Type': 'application/json;', 'Accept': '*/*'},
        body: JSON.stringify(object),
    }).then(function (response) {
        return response.json();
    })
}

async function RequestPUTApi(url, object, header) { //UPDATE
    return fetch(url, {
        method: "PUT",
        headers: header || {'Content-Type': 'application/json;', 'Accept': '*/*'},
        body: JSON.stringify(object),
    }).then(function (response) {
        return response.json();
    })
}

async function RequestDELETEApi(url) { //DELETE
    return fetch(url, {
        method: "DELETE",
    }).then(function (response) {
        return response.json();
    })
}
