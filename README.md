(work in progress)

TOSCA for Terraform
===================

[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)
[![Latest Release](https://img.shields.io/github/release/tliron/puccini-terraform.svg)](https://github.com/tliron/puccini-terraform/releases/latest)
[![Go Report Card](https://goreportcard.com/badge/github.com/tliron/puccini-terraform)](https://goreportcard.com/report/github.com/tliron/puccini-terraform)

Enable [TOSCA](https://www.oasis-open.org/committees/tosca/)
for [Terraform](https://www.terraform.io/)
using [Puccini](https://puccini.cloud/).

Two main features are supported:

TOSCA Preprocessor
------------------

Generate a Terraform module from TOSCA. Included is a Terraform profile. For example, this TOSCA:

```yaml
tosca_definitions_version: tosca_simple_yaml_1_3

imports:
- file: profiles/terraform/profile.yaml
  namespace_prefix: tf

topology_template:

  node_templates:

    hello:
      type: tf:LocalFile
      properties:
        filename: ./artifacts/hello.txt
        content: Hello World!
```

would become this Terraform module:

```hcl
resource local_file hello {
  filename = "./artifacts/hello.txt"
  content = "Hello World!"
}
```

TOSCA Provider
--------------

A Terraform resource that represents a TOSCA service template. Currently we simply compile it to a Clout,
but in the future we hope to allow integration with TOSCA orchestrators, such as
[Turandot](https://turandot.puccini.cloud/). Example:

```hcl
terraform {
  required_providers {
    tosca = {
      source = "puccini/tosca"
    }
  }
}

resource tosca_clout hello {
  service_template = "./hello.csar"
}

output clout {
  value = tosca_clout.hello
}
```
