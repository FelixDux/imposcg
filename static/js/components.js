import{extractFromAPIInfo, kvObjectToPairs} from './api-data.js';

class FullPathBuilder {
    constructor(apiData) {
        this.basePath = "";

        const setter = path => this.basePath = JSON.stringify(path).replace(/\"/g, "");

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

        const processPostData = data => this.description = data['description'];

        extractFromAPIInfo(apiData, 'post', processPostData);
    }

    html() {
        return `${this.path}: ${this.description}`;
    }
}

class PathsHolder {
    constructor(apiData) {
        this.paths = [];

        const setter = paths => {const pairs = kvObjectToPairs(paths); pairs.forEach(pair => {
            console.log(`${pair[0]}: ${pair[1]}`);
            this.paths.push(new Path(pair[0], pair[1]));
        });};

        extractFromAPIInfo(apiData, 'paths', setter);
    }

    html() {
        return `<ul>${this.paths.reduce((prev, curr) => prev.concat("<li>", curr.html()), "")}</ul>`;
    }
}

export {FullPathBuilder, Parameter, Path, PathsHolder};
