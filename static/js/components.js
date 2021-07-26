import{extractFromAPIInfo, kvObjectToPairs} from './api-data.js';

class FullPathBuilder {
    constructor(apiData) {

        const setter = path => this.basePath = JSON.stringify(path).replace(/\"/g, "");

        extractFromAPIInfo(apiData, 'basePath', setter);
    }

    fullPath(path) {
        return `${this.basePath}${path}`;
    }
}

class Parameter {
    constructor(apiData) {
        this.attributes = [];

        this.name = apiData['name'];
        this.description = apiData['description'];

        Object.keys(apiData).forEach( key => {
            if (!['name', 'description', 'in'].includes(key)) {
                this.attributes.push({name: key, value: apiData[key]});
            }
        })
    }

    html() {
        const attributeList = this.attributes.reduce(
            (acc, attribute) => `${acc}<li>${attribute.name}: ${attribute.value}`, ''
        );

        return `${this.name}: ${this.description} <ul>${attributeList}</ul>`;
    }
}

class Path {
    constructor(path, apiData) {

        this.path = path;

        const processPostData = data => {
            this.summary = data['summary'];
            this.description = data['description'];
            this.parameters = data['parameters'].map((e, i) => new Parameter(e));
        };

        extractFromAPIInfo(apiData, 'post', processPostData);
    }

    navHtml() {
        return `
        <div class="tooltip">
        <a>
        ${this.summary}
        <span class="tooltiptext">${this.description}</span>
        </a>
        </div>`;
    }

    debug() {

        const parameterList = this.parameters.reduce(
            (acc, parameter) => `${acc}<li>${parameter.html()}`,
            ""
        );

        return `${this.path}: ${this.description} <ul>${parameterList}</ul>`;
    }
}

class PathsHolder {
    constructor(apiData) {
        this.paths = [];

        const setter = paths => {const pairs = kvObjectToPairs(paths); pairs.forEach(pair => {
            this.paths.push(new Path(pair[0], pair[1]));
        });};

        extractFromAPIInfo(apiData, 'paths', setter);
    }

    html() {
        return `<div className="topnav">
        ${this.paths.reduce((prev, curr) => prev.concat(curr.navHtml()), "")}
        </div>`;
    }
}

class Header {
    constructor(apiData) {
        const setter = info => {
            this.title = info.title;
            this.version = info.version;
            this.description = info.description;
        }

        extractFromAPIInfo(apiData, 'info', setter);
    }

    html() {
        return `<header class="imposc-header">
                    <h1>${this.title}</h1>
                    <small>Version ${this.version}</small><br>
                    <small>${this.description}</small>
                </header>`;
    }
}

export {FullPathBuilder, Header, PathsHolder};
