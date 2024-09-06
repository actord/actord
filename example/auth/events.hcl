event create {
  field email {
    type = "string"
    label = "Email"
  }

  field password {
    type = "string"
    label = "Password"
  }

  field password_confirmation {
    type = "string"
    label = "Password Confirmation"
  }
}

event approve {
}

event block {
  field reason {
    type = "string"
    label = "Reason"
  }
}

event unblock {
}

event login {
  field password {
    type = "string"
    label = "Password"
  }
}
