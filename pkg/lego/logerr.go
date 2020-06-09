package lego

import (
	"emperror.dev/emperror"
	"github.com/vseinstrumentiru/lego/internal/lego"
	"logur.dev/logur"
)

type LogErr = lego.LogErr

type ErrorMatcher = lego.ErrorMatcher

type ContextExtractor = lego.ContextExtractor

func NewLogErr(logger logur.LoggerFacade, handler emperror.ErrorHandlerFacade) LogErr {
	return lego.NewLogErr(logger, handler)
}

// NewContextAwareLogger returns a new Logger instance that can extract information from a contexttool.
func NewContextAwareLogErr(logger logur.LoggerFacade, handler emperror.ErrorHandlerFacade, extractor ContextExtractor) LogErr {
	return lego.NewContextAwareLogErr(logger, handler, extractor)
}

type LogErrF = lego.LogErrF

func NewLogErrF(logErr LogErr) LogErrF {
	return lego.NewLogErrF(logErr)
}
