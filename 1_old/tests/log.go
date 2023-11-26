package tests

import "git.in.zhihu.com/antispam/datasupply/log"

var DefaultLogger = log.NewDefaultLog()

func init() {
	DefaultLogger.SetLevel(log.DebugLevel)
}
