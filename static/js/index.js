import {rendererForNode} from './render.js';
import{getAPIInfo} from './api-data.js';
import {FullPathBuilder, Header, PathsHolder} from './components.js';
import {addEventListeners} from './listeners.js'

function processAPIInfo(data) {
    const renderer = rendererForNode("main");

    const pathBuilder = new FullPathBuilder(data);

    const header = new Header(data);

    const paths = new PathsHolder(data);

    renderer(`${header.html()}${paths.html()}`);
}


getAPIInfo(processAPIInfo);

// Wait until the document is ready
document.addEventListener("DOMContentLoaded", function() { 
    addEventListeners();
});
