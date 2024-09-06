type actor {
  field email {
    type = "string"
    label = "Email"
  }
  field password {
    type = "string"
    label = "Password"
  }
  field block_reason {
    type = "string"
    label = "Block reason"
    optional = true
  }
}


