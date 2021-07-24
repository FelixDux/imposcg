var render = function (nodeId, content) {
    document.getElementById(nodeId).innerHTML = `${content}`;
};

function message2JSON(message) {
    return {"message": message};
}

function getAPIInfo(callback) {
    let theJson = message2JSON("Watch this space ...");

    fetch("/swagger/doc.json")
    .then(response => response.json())
    .then(data => callback(data))
    .catch(error => callback(message2JSON(`${error}`)));
}

function processAPIInfo(data) {
    render("main", JSON.stringify(data));
}


getAPIInfo(processAPIInfo);
