package cronexpr

import (
	"github.com/gorhill/cronexpr"
	log "github.com/sirupsen/logrus"
)

func test() {
	cron := cronexpr.MustParse("0 0 0 * * * *")
	log.Info(cron)
}
