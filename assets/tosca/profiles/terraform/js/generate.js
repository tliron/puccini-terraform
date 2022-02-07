
const traversal = require('tosca.lib.traversal');
const tosca = require('tosca.lib.utils');

traversal.coerce();

let tf = '';

for (let id in clout.vertexes) {
	let vertex = clout.vertexes[id];

    let terraform = getTerraform(vertex);
    if (terraform.resource) {
        //tf += JSON.stringify(terraform.resource)
        tf += 'resource ' + terraform.resource + ' ' + vertex.properties.name + ' {\n'
        for (let name in vertex.properties.properties) {
            let property = vertex.properties.properties[name];
            tf += '  ' + name + ' = \"' + property.replace(/"/g, '\\"') + '\"\n';
        }
        tf += '}\n';
    }
    /*if (isTerraform(vertex)) {
        tf += JSON.stringify(vertex.properties.properties)
    }*/
}

function getTerraform(vertex) {
    let terraform = {};
    if (tosca.isTosca(vertex, 'NodeTemplate')) {
        for (let name in vertex.properties.types) {
            let type = vertex.properties.types[name];
            if (type.metadata) {
                terraform.resource = type.metadata['terraform.resource'];
            }
        }
    }
    return terraform;
}

puccini.write(tf);