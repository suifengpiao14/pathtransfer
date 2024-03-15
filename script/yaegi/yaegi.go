package yaegi

import (
	"github.com/pkg/errors"
	_ "github.com/spf13/cast"
	_ "github.com/syyongx/php2go"
	_ "github.com/tidwall/gjson"
	_ "github.com/tidwall/sjson"
	"github.com/traefik/yaegi/interp"
	"github.com/traefik/yaegi/stdlib"
)

const (
	SCRIPT_LANGUAGE_GO = "go"
)

type ScriptGo struct {
	engine *interp.Interpreter
	code   []string
}

func (sgo ScriptGo) Language() string {
	return SCRIPT_LANGUAGE_GO
}

func (sgo *ScriptGo) WriteCode(codes ...string) {
	sgo.code = append(sgo.code, codes...)
}

func (sgo *ScriptGo) Compile() (err error) {
	engine := interp.New(interp.Options{})
	engine.Use(stdlib.Symbols)
	engine.Use(Symbols) //注册当前包结构体
	for _, code := range sgo.code {
		_, err = engine.Eval(code)
		if err != nil {
			err = errors.WithMessage(err, "init dynamic go script error")
			return err
		}
	}

	sgo.engine = engine
	return nil
}

func (sgo *ScriptGo) Run(script string) (out string, err error) {
	if sgo.engine == nil {
		err = sgo.Compile()
		if err != nil {
			return "", err
		}
	}
	rv, err := sgo.engine.Eval(script)
	if err != nil {
		return "", err
	}
	out = rv.String()
	return out, nil
}

func NewScriptGo() (sgo *ScriptGo) {
	return &ScriptGo{}
}

var Symbols = stdlib.Symbols

//go:generate go install github.com/traefik/yaegi/cmd/yaegi
//go:generate yaegi extract github.com/tidwall/gjson
//go:generate yaegi extract github.com/tidwall/sjson
//go:generate yaegi extract github.com/spf13/cast
//go:generate yaegi extract github.com/syyongx/php2go
//go:generate yaegi extract github.com/suifengpiao14/pathtransfer/script/yaegi/custom
