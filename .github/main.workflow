workflow "Docker" {
  on = "push"
  resolves = ["Docker Push"]
}

action "Docker Build" {
  uses = "actions/docker/cli@master"
  args = "build -t jest-action ."
}

action "Docker Tag" {
  uses = "actions/docker/tag@master"
  needs = ["Docker Build"]
  args = "--no-sha jest-action rkusa/jest-action"
}

action "Docker Login" {
  uses = "actions/docker/login@master"
  needs = ["Docker Tag"]
  secrets = ["DOCKER_USERNAME", "DOCKER_PASSWORD"]
}

action "Docker Push" {
  uses = "actions/docker/cli@master"
  needs = ["Docker Login"]
  args = "push rkusa/jest-action"
}