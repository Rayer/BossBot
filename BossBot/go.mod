module BossBotLib

require github.com/pkg/errors v0.8.0

replace Utilities v0.0.0 => ../Utilities

require (
	Utilities v0.0.0
	github.com/sirupsen/logrus v1.2.0
)
