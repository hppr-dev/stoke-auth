package usr

import (
	"github.com/rs/zerolog"
)


var logger zerolog.Logger

func SetLogger(rootLogger zerolog.Logger) {
	logger = rootLogger.With().Str("package", "stoke.internal.usr").Logger()
}
