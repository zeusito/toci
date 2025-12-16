package logger

import (
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func MustConfigure() {
	zerolog.TimeFieldFormat = time.RFC3339
	log.Logger = log.With().Caller().Logger()
}
