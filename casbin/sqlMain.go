package main

import (
	"fmt"

	"github.com/casbin/casbin/v2"
	xormadapter "github.com/casbin/xorm-adapter/v2"
	_ "github.com/lib/pq"
)

func main() {
	// 初始化 xorm adpater，在這裡與 DB 連線
	// 沒有指定 db 的話會幫你建立一個 casbin (加入 dbname=abc 可以指定使用 DB abc)
	// 進去會幫你檢查有沒有 casbin_rule 的資料表，沒有的話也會幫你加進去
	driverName := "postgres"
	pgUser, pgPassWd := "root", "root"
	dataSource := fmt.Sprintf("user=%s password=%s host=127.0.0.1 port=5432 sslmode=disable", pgUser, pgPassWd)
	a, _ := xormadapter.NewAdapter(driverName, dataSource) // Your driver and data source.

	// 直接讀檔案最快啦
	enforcer, _ := casbin.NewEnforcer("model_rbac.conf", a)

	// 第一次跑的話會需要 Policy 資料
	// 直接寫個 AddPolicy 會自動幫你同步進去
	// enforcer.EnableAutoSave(false) // 關閉自動同步
	{
		// 如果是 enforcer.AddPolicy 的話會幫你填第一個 "p"
		enforcer.AddNamedPolicy("p", "alice", "data1", "read")
		enforcer.AddNamedPolicy("p", "bob", "data2", "write")
		enforcer.AddNamedPolicy("p", "data_group_admin", "data_group", "write")
		// enforcer.AddGroupingPolicies 的話會幫你填第一個 "g"
		enforcer.AddNamedGroupingPolicy("g", "alice", "data_group_admin")
		enforcer.AddNamedGroupingPolicy("g2", "data1", "data_group")
		enforcer.AddNamedGroupingPolicy("g2", "data2", "data_group")
	}

	// Load the policy from DB.
	enforcer.LoadPolicy()

	// 實際測試時間:
	// 在 model [request_definition] 設定了三個變數輸入
	sub, obj, act := "alice", "data2", "write"
	result, _ := enforcer.Enforce(sub, obj, act)
	fmt.Printf("enforce1: %t\n", result) // true

	sub, obj, act = "alice", "data2", "read"
	result, _ = enforcer.Enforce(sub, obj, act)
	fmt.Printf("enforce2: %t\n", result) // false (因為 data_group_admin 對 data_group 沒有 read)
}
