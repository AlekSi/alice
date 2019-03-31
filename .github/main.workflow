workflow "Push" {
  on = "push"
  resolves = ["Push / Test with Go 1.12", "Push / Test with Go master"]
}

workflow "PR" {
  on = "pull_request"
  resolves = ["PR / Test with Go 1.12", "PR / Test with Go master"]
}

action "Push / Test with Go 1.12" {
  uses = "./.github/tests-go-1.12"
  secrets = ["CODECOV_TOKEN"]
}

action "Push / Test with Go master" {
  uses = "./.github/tests-go-master"
  secrets = ["CODECOV_TOKEN"]
}

action "PR / Test with Go 1.12" {
  uses = "./.github/tests-go-1.12"
  secrets = ["CODECOV_TOKEN"]
}

action "PR / Test with Go master" {
  uses = "./.github/tests-go-master"
  secrets = ["CODECOV_TOKEN"]
}
