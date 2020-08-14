package main

import (
	"fmt"

	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"
	scas "github.com/qiangmzsx/string-adapter/v2"
)

var modelStr = `
[request_definition]
r = sub, obj, act

[policy_definition]
p = sub, obj, act

[role_definition]
g = _, _
g2 = _, _

[policy_effect]
e = some(where (p.eft == allow))

[matchers]
m = g(r.sub, p.sub) && g2(r.obj, p.obj) && r.act == p.act
`

var policyStr = `
p, alice, data1, read
p, bob, data2, write
p, data_group_admin, data_group, write

g, alice, data_group_admin
g2, data1, data_group
g2, data2, data_group
`

func main() {
	// 設定 Model 模型
	// 參考 https://casbin.org/docs/en/model-storage
	m, _ := model.NewModelFromString(modelStr)

	// 設定 Policy 具體規則
	// 可以參考 https://casbin.org/docs/en/adapters
	p := scas.NewAdapter(policyStr)

	// 建立 Enforcer 需要輸入 Model 和 Policy
	// casbin 提供很多的 adapter 讓開發者用自己適合的方式填充資料 (file, db, redis, cloud...)
	// 範例選擇了最方便的直接讀文字 XDD
	enforcer, _ := casbin.NewEnforcer(m, p)

	// 實際測試時間:
	// 在 model [request_definition] 設定了三個變數輸入
	sub, obj, act := "alice", "data2", "write"
	result, _ := enforcer.Enforce(sub, obj, act)
	fmt.Printf("enforce1: %t\n", result) // true

	sub, obj, act = "alice", "data2", "read"
	result, _ = enforcer.Enforce(sub, obj, act)
	fmt.Printf("enforce2: %t\n", result) // false (因為 data_group_admin 對 data_group 沒有 read)
}
