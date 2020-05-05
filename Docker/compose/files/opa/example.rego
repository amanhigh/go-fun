package httpapi.authz

# bob is alice's manager, and betty is charlie's.
subordinates = {"alice": [], "charlie": [], "bob": ["alice"], "betty": ["charlie"]}

# HTTP API request
import input

default allow = false
default is_manager = false

# Allow users to get their own salaries.
allow {
  some username
  input.method == "GET"
  input.path = ["finance", "salary", username]
  input.user == username
}

# Allow managers to get their subordinates' salaries.
allow {
  some username
  input.method == "GET"
  input.path = ["finance", "salary", username]
  subordinates[input.user][_] == username
}

# Check who is Manager
is_manager {
    count(subordinates[input.user]) > 0
}