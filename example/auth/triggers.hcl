trigger create {
  event_type = "create"

  logic {
    condition "event.password" {
      equals=["event.password_confirmation"]
      on_failure {
        exception = "Password and password confirmation do not match"
      }
    }
  }
  logic {
    set "actor.email" {
      copy = "event.email"
    }
    set "actor.password_hash" {
      copy = "event.password"
      strconv = ["password_hash"]
    }
  }

  transition = "await_moderation"
}
