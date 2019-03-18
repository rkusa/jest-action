# Github Action for Jest (with Annotations)

This Action runs your [Jest](https://github.com/facebook/jest) test suite and adds annotations to the Github check the action is run in.

Annotaiton Example: https://github.com/rkusa/jest-action-example/pull/1/files

![Annotation Example](screenshot.png)

## Usage

```hcl
workflow "Tests" {
  on = "push"
  resolves = ["Jest"]
}

action "Dependencies" {
  uses = "actions/npm@master"
  args = "install"
}

action "Jest" {
  uses = "docker://rkusa/jest-action:latest"
  secrets = ["GITHUB_TOKEN"]
  args = ""
  needs = ["Dependencies"]
}
```

### Secrets

* `GITHUB_TOKEN` - **Required**. Required to add annotations to the check that is executing the Github action.

### Environment variables

* `JEST_CMD` - **Optional**. The path the Jest command - defaults to `./node_modules/.bin/jest`.

#### Example

To run Jest, either use the published docker image ...

```hcl
action "Jest" {
  uses = "docker://rkusa/jest-action:latest"
  secrets = ["GITHUB_TOKEN"]
  args = ""
}
```

... or the Github repo:

```hcl
action "Jest" {
  uses = "rkusa/jest-action@master"
  secrets = ["GITHUB_TOKEN"]
  args = ""
}
```

## License

The Dockerfile and associated scripts and documentation in this project are released under the [MIT License](LICENSE).

Container images built with this project include third party materials. View license information for [Node.js](https://github.com/nodejs/node/blob/master/LICENSE), [Node.js Docker project](https://github.com/nodejs/docker-node/blob/master/LICENSE), [Jest](https://github.com/facebook/jest/blob/master/LICENSE), [Go](https://golang.org/LICENSE), [google/go-github](https://github.com/google/go-github/blob/master/LICENSE) or [ldez/ghactions](https://github.com/ldez/ghactions/blob/master/LICENSE). As with all Docker images, these likely also contain other software which may be under other licenses. It is the image user's responsibility to ensure that any use of this image complies with any relevant licenses for all software contained within.