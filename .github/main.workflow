workflow "Main" {
  on = "pull_request"
  resolves = ["Test with Go 1.12", "Test with Go master"]
}

action "Test with Go 1.12" {
  uses = "./.github/tests-go-1.12"
  secrets = ["CODECOV_TOKEN"]
}

action "Test with Go master" {
  uses = "./.github/tests-go-master"
  secrets = ["CODECOV_TOKEN"]
}
