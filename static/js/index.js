import {rendererForNode} from './render.js';
import{getAPIInfo} from './api-data.js';
import {FullPathBuilder, Header, PathsHolder} from './components.js';

function processAPIInfo(data) {
    const renderer = rendererForNode("main");

    const pathBuilder = new FullPathBuilder(data);

    const header = new Header(data);

    const paths = new PathsHolder(data);

    renderer(`<p>${header.html()}</p>${paths.html()}`);
}


getAPIInfo(processAPIInfo);
