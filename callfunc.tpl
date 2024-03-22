{{- define "go" -}}
{{$package:=.FirstPackage}}
{{if $package}}
package {{.FirstPackage}}
{{end}}
import (
    "github.com/spf13/cast"
    "github.com/suifengpiao14/goscript/yaegi"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
)
{{range $callFunc:=.}}
func Call{{$callFunc.FuncName}}(input string) (outputDTO *yaegi.OutputDTO){
        {{range $arg:=$callFunc.Input -}}
            {{$arg.Name}} := cast.{{$arg.TypeConvertFunc}}(gjson.Get(input, "{{$arg.Path.NameWithoutSpace}}").String())
        {{ end}}
        var out string
		var err error
		outputDTO = &yaegi.OutputDTO{}
        {// 避免局部变量冲突
            {{$callFunc.Output.Names}}:={{$callFunc.FuncName}}({{$callFunc.Input.Names}})
            {{range $arg:=$callFunc.Output -}}
                out, err = sjson.Set(out, "{{$arg.Path.NameWithoutSpace}}", {{$arg.Name}})
                if err != nil {
                   outputDTO.Err = err
				    return outputDTO
                }
            {{ end}}
        }
        outputDTO.Data = out
		return outputDTO
	}
{{end}}
{{- end -}}

