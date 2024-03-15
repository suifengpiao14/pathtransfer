package yaegi_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/suifengpiao14/pathtransfer/script/yaegi"
)



func TestScriptGo(t *testing.T) {
	code := `
		package vocabulary
		func SetLimit(index int,size int)(offset int,limit int){
			limit=size
			offset=index*size
			return offset,limit
		}
	`
	callCode := `
package vocabulary
import (
	"errors"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
	"github.com/spf13/cast"
)

func CallSetLimit(input string) (out string, err error) {
        index := cast.ToInt(gjson.Get(input, "func.vocabulary.SetLimit.input.index").String())
        size :=  cast.ToInt(gjson.Get(input, "func.vocabulary.SetLimit.input.size").String())
        
        {// 避免局部变量冲突
            offset,size:=SetLimit(index,size)
            if err !=nil{
                return "",err
            }
            out, err = sjson.Set(out, "func.vocabulary.SetLimit.output.offset", offset)
                if err != nil {
                    return "", err
                }
            out, err = sjson.Set(out, "func.vocabulary.SetLimit.output.size", size)
                if err != nil {
                    return "", err
                }
            
        }
		return out, errors.New("hhah")
	}

`
	engine := yaegi.NewScriptGo()
	engine.WriteCode(code, callCode)
	input := `{"func":{"vocabulary":{"SetLimit":{"input":{"index":1,"size":20}}}}}`
	runCode := fmt.Sprintf("vocabulary.CallSetLimit(`%s`)", input)
	out, err := engine.Run(runCode)
	require.NoError(t, err)
	fmt.Println(out)

}
