var render = function (nodeId, content) {
    document.getElementById(nodeId).innerHTML = `${content}`;
};

function message2JSON(message) {
    return {"message": message};
}

class DocInfo {
    constructor(data) {
        this.json = data;
    }

    put(nodeId) {
        render(nodeId, JSON.stringify(this.json));
    }
}

const info = fetch("/swagger/doc.json")
.then(response => response.json())
.then(data => new DocInfo(data))
.catch(error => new DocInfo(message2JSON(`${error}`)))
;

info.put("main");

