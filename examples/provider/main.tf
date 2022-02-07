terraform {
  required_providers {
    tosca = {
      source = "puccini/tosca"
    }
  }
}

provider tosca {
  username = "education"
  password = "test123"
  host     = "http://localhost:19090"
}

resource tosca_clout hello {
  service_template = "./hello.yaml"
}

output clout {
  value = tosca_clout.hello
}
