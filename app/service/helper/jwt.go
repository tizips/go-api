package helper

import (
	"fmt"
	"github.com/golang-module/carbon/v2"
	"saas/app/helper/crypt"
	"saas/app/helper/str"
	"saas/kernel/app"
)

func JwtToken(id any) string {
	now := carbon.Now()
	return crypt.Md5(fmt.Sprintf("%s%v%d%s", app.Cfg.Server.Name, id, now.Timestamp(), str.Random(8)))
}
