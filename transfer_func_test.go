package pathtransfer_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/suifengpiao14/pathtransfer"
)

func TestCallFunc(t *testing.T) {
	transferLine := `
func.vocabulary.SetLimit.input.index@int:Dictionary.pagination.index
func.vocabulary.SetLimit.input.size@int:Dictionary.pagination.size
func.vocabulary.SetLimit.output.offset@int:Dictionary.limit.offset
func.vocabulary.SetLimit.output.size@int:Dictionary.limit.size
	`
	transfer := pathtransfer.Parse(transferLine)
	script, err := transfer.GetCallFnScript("go")
	require.NoError(t, err)
	fmt.Println(script)

}

func TestExplainFuncPath(t *testing.T) {
	t.Run("bas arg", func(t *testing.T) {
		funcPath := "func.vocabulary.SetLimit.input.index@int"
		arg, err := pathtransfer.ExplainFuncPath(funcPath)
		require.NoError(t, err)
		fmt.Println(arg)
	})
	t.Run("object arg", func(t *testing.T) {
		funcPath := "func.vocabulary.SetLimit.input.pagination.index@int"
		arg, err := pathtransfer.ExplainFuncPath(funcPath)
		require.NoError(t, err)
		fmt.Println(arg)
	})

}
