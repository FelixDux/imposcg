var render = function (template, node) {
    node.innerHTML = template;
};

var template = '<h1>Watch This Space ...</h1>';
render(template, document.getElementById('main'));