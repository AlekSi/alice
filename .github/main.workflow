workflow "Push" {
  on = "push"
  resolves = ["Push: test with Go 1.12", "Push: test with Go master"]
}

workflow "PR" {
  on = "pull_request"
  resolves = ["PR: test with Go 1.12", "PR: test with Go master"]
}

action "Push: test with Go 1.12" {
  uses = "./.github/tests-go-1.12"
  secrets = ["CODECOV_TOKEN"]
}

action "Push: test with Go master" {
  uses = "./.github/tests-go-master"
  secrets = ["CODECOV_TOKEN"]
}

action "PR: test with Go 1.12" {
  uses = "./.github/tests-go-1.12"
  secrets = ["CODECOV_TOKEN"]
}

action "PR: test with Go master" {
  uses = "./.github/tests-go-master"
  secrets = ["CODECOV_TOKEN"]
}
