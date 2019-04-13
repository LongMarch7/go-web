package main

import (
    "fmt"
    "github.com/casbin/casbin"
    "github.com/LongMarch7/go-web/auth"
    _ "github.com/go-sql-driver/mysql"
)

func main() {
    a := auth.NewAdapter("mysql", "root:123456@tcp(127.0.0.1:13306)/", "test") // Your driver and data source.
    //e := casbin.NewEnforcer("config/rbac_model.conf", a)
    text :=
        `
[request_definition]
r = sub, obj, act

[policy_definition]
p = sub, obj, act

[role_definition]
g = _, _

[policy_effect]
e = some(where (p.eft == allow))

[matchers]
m = g(r.sub, p.sub) && r.obj == p.obj && r.act == p.act
`
    m := casbin.NewModel(text)
    e := casbin.NewEnforcer(m, a)
    e.LoadPolicy()

    police := []string{"alice", "data1", "read"}
    e.AddPolicy(police)
    e.AddRoleForUser("admin","alice")
    // Check the permission.
    if e.Enforce("admin", "data1", "read") {
        fmt.Println("allow")
    }else{
        fmt.Println("deny")
    }

    // Modify the policy.
    // e.AddPolicy(...)
    // e.RemovePolicy(...)

    // Save the policy back to DB.
    e.SavePolicy()
}