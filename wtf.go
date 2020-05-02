package wtf

import (
	"fmt"
	"reflect"
)

type HookFunction func(*Failure)

var userDefinedDefaults = map[string]*Failure{}
var userDefinedHooks []HookFunction

func AddDefaultCaseFailure(inCase error, message string, code int) {
	rt := reflect.TypeOf(inCase).String()
	userDefinedDefaults[rt] = New(message, code)
}

func AddUnknownErrorHookFailure(f HookFunction) {
	userDefinedHooks = append(userDefinedHooks, f)
}

func callHooks(err *Failure) {
	for _, v := range userDefinedHooks {
		v(err)
	}
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
	Code         int         `json:"code"`
	Message      string      `json:"message"`
	ExtraMessage *string     `json:"details"`
	origin       interface{} `json:"-"`
}

func (e Failure) MessageString() string {
	var extra string
	if nil != e.ExtraMessage {
		extra = fmt.Sprintf("!extra[ %s ]", e.ExtraMessage)
	}
	return fmt.Sprintf("%s %s", e.Message, extra)
}

func (e Failure) Error() string {
	return e.MessageString()
}

func (e *Failure) WithMessage(message string) *Failure {
	e.ExtraMessage = &message
	return e
}

func (e *Failure) WitCode(code int) *Failure {
	e.Code = code
	return e
}

func (e *Failure) setOrigin(origin interface{}) *Failure {
	e.origin = origin
	return e
}

func (e *Failure) GetOrigin() interface{} {
	return e.origin
}

func (e *Failure) Hook() *Failure {
	callHooks(e)
	return e
}

func Wrap(err interface{}) *Failure {
	switch err.(type) {
	case Failure:
		return err.(*Failure)
	case error:
		userDef, ok := userDefinedDefaults[reflect.TypeOf(err).String()]
		if true == ok {
			return userDef.setOrigin(err)
		}
		return New(err.(error).Error(), mainConfig.DefaultErrorCode).setOrigin(err)
	case string:
		return New(err.(string), mainConfig.DefaultErrorCode).setOrigin(err)
	default:
		return New(mainConfig.WrapUnknownMessage, mainConfig.DefaultErrorCode).setOrigin(err)
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
