package script

import (
	"strings"

	"github.com/pkg/errors"
	"github.com/suifengpiao14/pathtransfer/script/yaegi"
)

const (
	SCRIPT_LANGUAGE_GO = yaegi.SCRIPT_LANGUAGE_GO
)

type ScriptI interface {
	Language() string
	Compile() (err error)
	Run(script string) (out string, err error)
	CallFuncScript(funcName string, input string) (callFuncScript string) //最终调用函数代码
}

func NewScriptEngine(language string) (scriptI ScriptI, err error) {
	m := map[string]func() ScriptI{
		SCRIPT_LANGUAGE_GO: func() ScriptI {
			return yaegi.NewScriptGo()
		},
	}
	fn, ok := m[strings.ToLower(language)]
	if !ok {
		err = errors.Errorf("not found script engine by language:%s", language)
		return nil, err
	}
	return fn(), nil
}
