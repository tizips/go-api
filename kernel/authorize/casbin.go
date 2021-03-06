package authorize

import (
	"fmt"
	"github.com/casbin/casbin/v2"
	adapter "github.com/casbin/gorm-adapter/v3"
	"github.com/gookit/color"
	"saas/kernel/config"
	"saas/kernel/data"
)

const ROOT = 888

var Casbin *casbin.Enforcer

func InitCasbin() {

	a, err := adapter.NewAdapterByDBUseTableName(data.Database, config.Values.Database.Prefix, "sys_casbin")
	if err != nil {
		color.Errorf("Casbin new adapter error: %v", err)
		return
	}

	Casbin, err = casbin.NewEnforcer(config.Application.Path+"/conf/casbin.conf", a)
	if err != nil {
		color.Errorf("Casbin new enforcer error: %v", err)
		return
	}

}

func NameByAdmin(id any) string {
	return fmt.Sprintf("admin:%v", id)
}

func NameByRole(id any) string {
	return fmt.Sprintf("role:%v", id)
}

func Root(id any) bool {
	exist, _ := Casbin.HasRoleForUser(NameByAdmin(id), NameByRole(ROOT))
	return exist
}
