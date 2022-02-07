terraform {
  required_providers {
    docker = {
      source = "kreuzwerker/docker"
    }
  }
}

provider docker {
  host = "unix:///run/user/1000/podman/podman.sock"
}

resource docker_image fedora {
  name = "fedora:latest"
}

resource docker_container hello {
  image = docker_image.fedora.latest
  name = "hello"
}
