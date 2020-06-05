package gofun.authz

default allow = false

resource_mappings[mapping]{
	#Direct Match
	mapping:=data.role_mapping[input.resource][_]
}{
	#Climb Tree
    parent := data.resource_tree[input.resource][_]
    mapping:=data.role_mapping[parent][_]
}

allow {
    #For each Mapping
    m:=resource_mappings[_]
    match_mapping(m)
}

match_mapping(m) {
    #Match Principle
    match_principle(m)

    #Match Scope
    match_scope(m)

    #Match Role
    match_role(m.role)
}

#Negation Example
# !(a | b) = !a & !b
has_no_role {
	not allow
}

match_principle(m)
{
  #Direct Match Input principle
    m.principle==input.principle
}

match_principle(m)
{
  #Regex Match
    endswith(m.principle,".*")
    re_match(m.principle,input.principle)
}

match_role(roleName)
{
  #Direct Match Role
  role:=data.role[roleName]

  match_role_permissions(role)
}

match_role(roleName)
{
  #Match Linked Roles
  linkedRoleName:=data.role_group[roleName][_]
  role:=data.role[linkedRoleName]

  match_role_permissions(role)
}

match_role_permissions(role)
{
  #For Each Permission
  p:=role[_]

  #Match
  match_permission(p)
}

match_permission(p)
{
 #Direct Match Action and Resource Type
  p.action=input.action
  p.resourceType=input.resourceType
}

match_scope(m)
{
  #No Scope
  not m.scope
}

match_scope(m)
{
  #Direct Match
  m.scope==input.scope
}

match_scope(m)
{
    #Regex Match
    endswith(m.scope,".*")
    re_match(m.scope,input.scope)
}


#Test Cases
test_direct_match {
    allow with input as {
        "principle": "bob",
        "resource": "gujrat",
        "resourceType": "image",
        "action": "read"
    }

    allow with input as {
        "principle": "laxmi",
        "resource": "bengal",
        "resourceType": "image",
        "scope": "puja",
        "action": "read"
    }

    #Negative
    not allow with input as {
        "principle": "steeve",
        "resource": "gujrat",
        "resourceType": "image",
        "action": "read"
    }
    not allow with input as {
        "principle": "bob",
        "resource": "bombay",
        "resourceType": "image",
        "action": "read"
    }
}

test_parent_match {
    allow with input as {
        "principle": "bob",
        "resource": "surat",
        "resourceType": "image",
        "action": "read"
    }

    #Negative
    not allow with input as {
        "principle": "steeve",
        "resource": "gujrat",
        "resourceType": "image",
        "action": "read"
    }
    not allow with input as {
        "principle": "bob",
        "resource": "india",
        "resourceType": "image",
        "action": "read"
    }
}


test_regex_match {
    allow with input as {
      "principle": "alice",
      "resource": "gujrat",
      "resourceType": "image",
      "scope": "foo-1",
      "action": "read"
  }

  #Negative
  #User Mismatch
    not allow with input as {
      "principle": "mogambo",
      "resource": "gujrat",
      "resourceType": "image",
      "scope": "foo-1",
      "action": "read"
  }

   #Scope Mismatch
    not allow with input as {
      "principle": "alice",
      "resource": "gujrat",
      "resourceType": "image",
      "scope": "bar-1",
      "action": "read"
  }
}
