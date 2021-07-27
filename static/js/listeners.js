import{getAPIInfo} from './api-data.js';
import {PathsHolder} from './components.js';

function listenersFromAPI(data, symbols) {

    const paths = new PathsHolder(data, symbols);

    paths.addListeners();
}

function addEventListeners(symbols) {
    const fetcher = data => listenersFromAPI(data, symbols);
    getAPIInfo(fetcher);
}

export {addEventListeners};
