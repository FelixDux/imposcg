import {rendererForNode} from './render.js';
import{getAPIInfo} from './api-data.js';
import {FullPathBuilder, Header, PathsHolder, ParameterSymbols} from './components.js';
import {addEventListeners} from './listeners.js'

function populate() {
    const symbols = new ParameterSymbols();

    function processAPIInfo(data) {
        const renderer = rendererForNode("main");

        const pathBuilder = new FullPathBuilder(data);

        const header = new Header(data);

        const paths = new PathsHolder(data, symbols);

        renderer(`${header.html()}${paths.html()}`);
    }

        getAPIInfo(processAPIInfo);

    // Wait until the document is ready
    document.addEventListener("DOMContentLoaded", function() { 
        addEventListeners(symbols);
    });
}

populate();

// TODO look at this for testing https://dev.to/thawkin3/how-to-unit-test-html-and-vanilla-javascript-without-a-ui-framework-4io
