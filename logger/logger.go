package logger

import (
	env "raptor/config"
	"fmt"
	"github.com/rs/zerolog/log"
)

var Logger = log.With().Logger()

func Trace(msg interface{}) {
	if env.GetString("log_lvl") == "debug" {
		Logger.Trace().Interface("message", fmt.Sprintf("%v", msg)).Msg("")
	}
}

func Debug(msg interface{}) {
	if env.GetString("log_lvl") == "debug" {
		Logger.Debug().Interface("message", fmt.Sprintf("%v", msg)).Msg("")
	}
}

func Info(msg interface{}) {
	if env.GetString("log_lvl") == "debug" || env.GetString("log_lvl") == "info" {
		Logger.Info().Interface("message", fmt.Sprintf("%v", msg)).Msg("")
	}
}

func Warn(msg interface{}) {
	if env.GetString("log_lvl") == "debug" || env.GetString("log_lvl") == "warn" {
		Logger.Warn().Interface("message", fmt.Sprintf("%v", msg)).Msg("")
	}
}

func Error(msg interface{}) {
	if env.GetString("log_lvl") == "debug" || env.GetString("log_lvl") == "error" {
		Logger.Error().Interface("message", fmt.Sprintf("%v", msg)).Msg("")

	}
}
