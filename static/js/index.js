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

class FullPathBuilder {
    constructor(apiData) {
        this.basePath = "";

        const setter = path => this.basePath = path.replace(/\"/g, "");

        extractFromAPIInfo(apiData, 'basePath', setter);
    }

    fullPath(path) {
        return `${this.basePath}${path}`;
    }
}

class Parameter {
    constructor(apiData) {
    }
}

class Path {
    constructor(path, apiData) {
        this.path = path;
        this.description = "";

        // this.parameters = [];

        function processPostData(data) {
            this.description = data['description'];
        }

        extractFromAPIInfo(apiData, 'post', processPostData);
    }

    html() {
        return `${this.path}: ${this.description}`;
    }
}

class PathsHolder {
    constructor(apiData) {
        this.paths = [];

        const setter = paths => {data=JSON.parse(paths); data.forEach(k, v => {
            this.paths.push(new Path(k, v));
        });};

        extractFromAPIInfo(apiData, 'paths', setter);
    }

    // html() {
    //     return `<ul>${this.paths.reduce(prev, curr => prev.concat("<li>", curr))}</ul>`;
    // }
}

function processAPIInfo(data) {
    const renderer = rendererForNode("main");

    const pathBuilder = new FullPathBuilder(data);

    const paths = new PathsHolder(data);

    // renderer(`<p>${pathBuilder.fullPath("/some/path/or/other")}</p>${paths.html()}`);
}


getAPIInfo(processAPIInfo);
