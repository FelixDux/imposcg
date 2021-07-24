function render(nodeId, content) {
    document.getElementById(nodeId).innerHTML = `${content}`;
};

function rendererForNode(nodeId) {
    return content => render(nodeId, content);
}

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

// class NavBarInfo {
//     constructor(apiData) {
//         this
//     }
// }

function extractFromAPIInfo(data, key, callback) {
    if (key in data) {
        info = data[key];
        callback(JSON.stringify(info));
    }
    else
    {
        callback(`Could not find key '${key}' in JSON data`);
    }
}

class PathBuilder {
    constructor(apiData) {
        this.basePath = "";

        const setter = path => this.basePath = path.replace(/\"/g, "");

        extractFromAPIInfo(apiData, 'basePath', setter);
    }

    fullPath(path) {
        return `${this.basePath}${path}`;
    }
}

function processAPIInfo(data) {
    const renderer = rendererForNode("main");

    const pathBuilder = new PathBuilder(data);

    renderer(pathBuilder.fullPath("/some/path/or/other"))
}


getAPIInfo(processAPIInfo);
