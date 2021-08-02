import{extractFromAPIInfo, kvObjectToPairs} from './api-data.js';
import {rendererForNode} from './render.js';

const refreshForm =  rendererForNode('form');

const imageRefresher = rendererForNode("imageTarget");

function refreshImage(blob, result) {
    imageRefresher(`<img src=${URL.createObjectURL(blob)} alt=${result} width="50%" align="center" />`);
}

function refreshImageWithError(error) {
    imageRefresher(`<h1>Error</h1>${error.toString()}`);
}

class ParameterSymbols {
    constructor(symbolsGetter) {
        this.symbols = new Map();

        const setter = (data) => {this.addSymbols(data)};

        symbolsGetter(setter);
    }

    addSymbols(data) {
        try {
            data.Symbols.forEach( (e, _) => {
                this.symbols.set(e.Parameter, e.Symbol);
            });
        }
        catch {
            console.log('Symbols lookup initialised with invalid data');
        }
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

        this.basePath = "";

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

        if ('name' in apiData && 'description' in apiData) {
            this.name = apiData.name;
            this.label = symbols.lookup(apiData.name);
            this.description = apiData.description;

            Object.keys(apiData).forEach( key => {
                if (!['name', 'description', 'in'].includes(key)) {
                    this.attributes.push({name: key, value: apiData[key]});
                }
            })
        }
        else {
            console.log(`Parameter initialised with invalid data: ${apiData}`);
        }
    }

    html() {
        function renderAttribute(name, value) {
            switch(name) {
                case "default":
                    return `value = ${value}`;
                case "minimum":
                    return `min = ${value}`;
                case "maximum":
                    return `max = ${value}`;
                case "required":
                    if (value) {
                        return " required";
                    }
                    else {
                        return " ";
                    }
                case "type":
                    if (value === "number") {
                        return `${name} = ${value} step=0.01 style="width: 6em"`;
                    }
                    else {
                        return `${name} = ${value} size=6`;
                    }
                default:
                    return `${name} = ${value}`;
            }
        }

        const attributeList = this.attributes.reduce(
            (acc, attribute) => `${acc} \n${renderAttribute(attribute.name,attribute.value)}`, ''
        );

        return `
        <tr class = "inputGroup">
          <td class = "inputGroup" >
            <div class="tooltip">${this.label}
            <span class="tooltiptext">${this.description}</span></div>
            </td>
          <td class = "inputGroup" >
            <input 
            id=${this.name}
            name=${this.name}
            ${attributeList}
          />
          </td>
        </tr>`;
    }
}

class Path {
    constructor(path, apiData, symbols) {

        this.path = path;

        this.state = {src: "", result: "", blob: {}};

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
        const action = this.path;
        document.getElementById(this.id).addEventListener("click", function () {refreshForm(formContent); document.getElementById('form').action = action;});
    }

    navHtml() {
        return `
        <div class="tooltip">
        <a id='${this.id}' class="topnav">
        ${this.summary}
        <span class="tooltiptext">${this.description}</span>
        </a>
        </div>`;
    }

    submitHtml() {
        if (this.parameters.length > 0) {
            return '<input type ="submit" value="Show" >'
        }
        else {
            return ""
        }
    }

    formHtml() {

        const parameterList = this.parameters.reduce(
            (acc, parameter) => `${acc} ${parameter.html()}`,
            ""
        );

        return `
        <div class="row">
        <div class="column left"><p /></div>
        <div class="column"><em>${this.description}</em>
      <table class="inputGroup"><tbody>${parameterList}</tbody></table>
      ${this.submitHtml()}</div>      
      <div class="column right" id="imageTarget"></div>
      </div>`;
    }
}

class PathsHolder {
    constructor(apiData, symbols) {
        this.paths = [];

        const pathBuilder = new FullPathBuilder(apiData);

        const setter = paths => {const pairs = kvObjectToPairs(paths); pairs.forEach(pair => {
            if ('post' in pair[1] && pair[0].endsWith("image/")) {
                this.paths.push(new Path(pathBuilder.fullPath(pair[0]), pair[1], symbols))
            };
        });};

        extractFromAPIInfo(apiData, 'paths', setter);
    }

    html() {
        return `<div><div class="topnav">
        ${this.paths.reduce((prev, curr) => prev.concat(curr.navHtml()), "")}
        </div>
        <form id="form"></form></div>`;
    }

    addListeners() {
        this.paths.forEach(path => path.addListener());

        let elem = document.getElementById('form');

        elem.addEventListener("submit", function (event) {
            fetch(elem.action,
              {
                method: 'POST',
                body: new FormData(elem)
              })
            .then(response => {
              if (response.status >= 400 || !response.ok) {
                  response.json().then((json) => {                    
                      if ('error' in json) {
                          throw json['error'];
                      }
                      else {
                          throw 'Unknown error';
                      }
                  });
              }
              else {
                return response.blob();
              }
            })
            .then(blob => {
                if (blob) {
                    refreshImage(blob, "Done");
                }
            })
            .catch((error) => {
              console.error('Fetch operation failed:', error);
              refreshImageWithError(error);
            });
            event.preventDefault();
            
        });
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

export {FullPathBuilder, Parameter, Header, PathsHolder, ParameterSymbols};
