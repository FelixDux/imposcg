import{getAPIInfo} from './api-data.js';
import {PathsHolder} from './components.js';

function listenersFromAPI(data) {

    const paths = new PathsHolder(data);

    paths.addListeners();
}

function addEventListeners() {
    getAPIInfo(listenersFromAPI);
}

export {addEventListeners};
