import{extractFromAPIInfo, kvObjectToPairs, getParameterSymbols} from './api-data.js';
import {rendererForNode} from './render.js';

const refreshForm =  rendererForNode('form');

class ParameterSymbols {
    constructor() {
        this.symbols = new Map();;

        const setter = data => {
            data.Symbols.forEach( (e, _) => {
                this.symbols.set(e.Parameter, e.Symbol);
            });
        };

        getParameterSymbols(setter);
    }

    lookup(parameter) {
        if (this.symbols.has(parameter)) {
            return this.symbols.get(parameter);
        }
        else {
            return parameter;
        }
    }
}

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
    constructor(apiData, symbols) {
        this.attributes = [];

        this.name = symbols.lookup(apiData['name']);
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
    constructor(path, apiData, symbols) {

        this.path = path;

        const processPostData = data => {
            this.summary = data['summary'];
            this.id = this.summary.replace(/\s/g, "");
            this.description = data['description'];
            this.parameters = data['parameters'].map((e, i) => new Parameter(e, symbols));
        };

        extractFromAPIInfo(apiData, 'post', processPostData);
    }

    addListener() {
        const formContent = this.formHtml();
        document.getElementById(this.id).addEventListener("click", function () {refreshForm(formContent);});
    }

    navHtml() {
        return `
        <div class="tooltip">
        <a id='${this.id}'">
        ${this.summary}
        <span class="tooltiptext">${this.description}</span>
        </a>
        </div>`;
    }

    formHtml() {

        const parameterList = this.parameters.reduce(
            (acc, parameter) => `${acc} ${parameter.html()}`,
            ""
        );

        return `<p>${this.description}</p>${parameterList}`;
    }
}

class PathsHolder {
    constructor(apiData, symbols) {
        this.paths = [];

        const setter = paths => {const pairs = kvObjectToPairs(paths); pairs.forEach(pair => {
            if ('post' in pair[1]) {
                this.paths.push(new Path(pair[0], pair[1], symbols))
            };
        });};

        extractFromAPIInfo(apiData, 'paths', setter);
    }

    html() {
        return `<div className="topnav">
        ${this.paths.reduce((prev, curr) => prev.concat(curr.navHtml()), "")}
        </div>
        <div id="form"></div>`;
    }

    addListeners() {
        this.paths.forEach(path => path.addListener());
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

export {FullPathBuilder, Header, PathsHolder, ParameterSymbols};
