package services

import (
	"goblocks/app/services/blocks"

	"go.uber.org/fx"
)

var Module = fx.Module(
	"services",
	fx.Provide(
		blocks.NewBlockManager,
	),
)
