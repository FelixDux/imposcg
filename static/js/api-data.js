function getAPIInfo(callback) {
    fetch("/swagger/doc.json")
    .then(response => response.json())
    .then(data => callback(data))
    .catch(error => console.log(`${error}`));
}

function extractFromAPIInfo(data, key, callback) {
    if (key in data) {
        callback(data[key]);
    }
    else
    {
        console.log(`Could not find key '${key}' in ${JSON.stringify(data)}`);
    }
}

function kvObjectToPairs(obj) {
    const keys = Object.keys(obj);
    const values = Object.values(obj);

    return keys.map( (element, index) => [element, values[index]] );
}

export {getAPIInfo, extractFromAPIInfo, kvObjectToPairs};