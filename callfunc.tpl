{{- define "go" -}}
package {{.FirstPackage}}
import (
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
    "github.com/spf13/cast"
)
{{range $callFunc:=.}}
func Call{{$callFunc.FuncName}}(input string) (out string, err error) {
        {{range $arg:=$callFunc.Input -}}
            {{$arg.Name}} := cast.{{$arg.TypeConvertFunc}}(gjson.Get(input, "{{$arg.Path}}").String())
        {{ end}}
        {// 避免局部变量冲突
            {{$callFunc.Output.Names}}:={{$callFunc.FuncName}}({{$callFunc.Input.Names}})
            if err !=nil{
                return "",err
            }
            {{range $arg:=$callFunc.Output -}}
                out, err = sjson.Set(out, "{{$arg.Path}}", {{$arg.Name}})
                if err != nil {
                    return "", err
                }
            {{ end}}
        }
		return out, nil
	}
{{end}}
{{- end -}}

