package wtf

import (
	_ "errors"
	"fmt"
)

type HookFunction func(error)

var userDefinedDefaults = map[error]*Failure{}
var userDefinedHooks []HookFunction

func AddDefaultCaseFailure(inCase error, message string, code int) {
	userDefinedDefaults[inCase] = New(message, code)
}

func AddUnknownErrorHookFailure(f HookFunction) {
	userDefinedHooks = append(userDefinedHooks, f)
}

func init() {
	userDefinedHooks = make([]HookFunction, 0)
}

type Config struct {
	DefaultErrorCode   int
	WrapUnknownMessage string
}

var defaultConfig = Config{
	DefaultErrorCode:   100,
	WrapUnknownMessage: "Unknown Failure",
}

var mainConfig = defaultConfig

func Configure(config Config) {
	mainConfig = config
}

type Failure struct {
	Code         int    `json:"code"`
	Message      string `json:"message"`
	ExtraMessage string `json:"-"`
}

func (e Failure) MessageString() string {
	return fmt.Sprintf("%s !extra[%s]", e.Message, e.ExtraMessage)
}

func (e Failure) Error() string {
	return e.MessageString()
}

func (e *Failure) WithMessage(message string) *Failure {
	e.ExtraMessage = message
	return e
}

func (e *Failure) WitCode(code int) *Failure {
	e.Code = code
	return e
}

func Wrap(err interface{}) *Failure {

	switch err.(type) {
	case Failure:
		return err.(*Failure)
	case error:

		userDef, ok := userDefinedDefaults[err.(error)]
		if true == ok {
			return userDef
		}

		return New(err.(error).Error(), mainConfig.DefaultErrorCode)
	case string:
		return New(err.(string), mainConfig.DefaultErrorCode)
	default:
		return New(mainConfig.WrapUnknownMessage, mainConfig.DefaultErrorCode)
	}
}

func (e Failure) Panic() {
	panic(e)
}

func New(message string, code int) *Failure {
	return &Failure{Code: code, Message: message}
}

func Define(message string, code int) *Failure {
	return New(message, code)
}

func IsFailure(toTest error) bool {
	switch toTest.(type) {
	case Failure:
		return true
	}
	return false
}

func Default() *Failure {
	return New(mainConfig.WrapUnknownMessage, mainConfig.DefaultErrorCode)
}
