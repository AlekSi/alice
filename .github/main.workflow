workflow "Main" {
  on = "push"
  resolves = ["Run tests with Go 1.12", "Run tests with Go master"]
}

action "Run tests with Go 1.12" {
  uses = "./.github/tests-go-1.12"
}

action "Run tests with Go master" {
  uses = "./.github/tests-go-master"
}
