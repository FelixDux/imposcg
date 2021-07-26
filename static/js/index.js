import {rendererForNode} from './render.js';
import{getAPIInfo} from './api-data.js';
import {FullPathBuilder, Path, PathsHolder} from './components.js';

function processAPIInfo(data) {
    const renderer = rendererForNode("main");

    const pathBuilder = new FullPathBuilder(data);

    const paths = new PathsHolder(data);

    renderer(`<p>${pathBuilder.fullPath("/some/path/or/other")}</p>${paths.html()}`);
}


getAPIInfo(processAPIInfo);
