process {
  state await_moderation {
    logic {
      label = "Moderator resolution"
      await {
        event approve {
          transition = "approved"
        }
        event block {
          transition = "blocked"
        }
        on_unknown {
          exception = "Unknown event"
        }
      }
    }
  }

  state approved {
    logic {
      label = "User approved"
      await {
        event block {
          transition = "blocked"
        }
        event login {
          handler = "auth"
        }
        on_unknown {
          exception = "Unknown event for approved user"
        }
      }
    }

  }

  state blocked {
    logic {
      set "actor.block_reason" {
        copy = "event.reason"
      }
    }
    logic {
      label = "User blocked"
      await {
        event unblock {
        }
        on_unknown {
          exception = "Unknown event for blocked user"
        }
      }
    }
    // user unblocked
    logic {
      set "actor.block_reason" {
        json = "null"
      }
      transition = "approved"
    }
  }

  handler auth {
    logic {
      verify_password "event.password" {
        hash = "actor.password_hash"
        set_to = "temp.is_valid"
      }
      reply {
        from = "temp"
      }
    }
  }

}
