package main

import (
	"fmt"

	"github.com/casbin/casbin/v2"
)

func main() {
	// 直接讀檔案最快啦
	enforcer, _ := casbin.NewEnforcer("model_rbac.conf", "policy_rbac.csv")

	// 實際測試時間:
	// 在 model [request_definition] 設定了三個變數輸入
	sub, obj, act := "alice", "data2", "write"
	result, _ := enforcer.Enforce(sub, obj, act)
	fmt.Printf("enforce1: %t\n", result) // true

	sub, obj, act = "alice", "data2", "read"
	result, _ = enforcer.Enforce(sub, obj, act)
	fmt.Printf("enforce2: %t\n", result) // false (因為 data_group_admin 對 data_group 沒有 read)
}
