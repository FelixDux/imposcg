function getAPIInfo(callback) {
    fetch("/swagger/doc.json")
    .then(response => response.json())
    .then(data => callback(data))
    .catch(error => callback(message2JSONObject(`${error}`)));
}

function extractFromAPIInfo(data, key, callback) {
    if (key in data) {
        callback(data[key]);
    }
    else
    {
        callback(message2JSONObject(`Could not find key '${key}' in ${JSON.stringify(data)}`));
    }
}

function kvObjectToPairs(obj) {
    const keys = Object.keys(obj);
    const values = Object.values(obj);

    return keys.map( (element, index) => [element, values[index]] );
}

function message2JSONObject(message) {
    return {"message": message};
}

export {getAPIInfo, extractFromAPIInfo, kvObjectToPairs};