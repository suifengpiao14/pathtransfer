package pathtransfer_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/suifengpiao14/pathtransfer"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
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
		tu := pathtransfer.TransferUnit{
			Path: pathtransfer.Path("func.vocabulary.SetLimit.input"),
			Type: "int",
		}
		arg, err := tu.FuncParameter()
		require.NoError(t, err)
		fmt.Println(arg)
	})
	t.Run("object arg", func(t *testing.T) {
		tu := pathtransfer.TransferUnit{
			Path: pathtransfer.Path("func.vocabulary.SetLimit.input.pagination.index"),
			Type: "int",
		}
		arg, err := tu.FuncParameter()
		require.NoError(t, err)
		fmt.Println(arg)
	})

}

func TestGetTransferFuncname(t *testing.T) {

	transfers := make(pathtransfer.Transfers, 0)

	transfers = append(transfers, pathtransfer.Transfer{
		Src: pathtransfer.TransferUnit{Path: "user.name", Type: "string"},
		Dst: pathtransfer.TransferUnit{Path: "data.userName"},
	})

	transfers = append(transfers, pathtransfer.Transfer{
		Src: pathtransfer.TransferUnit{Path: "func.vocabulary.SetLimit.input.index", Type: "int"},
		Dst: pathtransfer.TransferUnit{Path: "data.pagination.index"},
	})

	transfers = append(transfers, pathtransfer.Transfer{
		Src: pathtransfer.TransferUnit{Path: "func.vocabulary.SetLimit.input.size", Type: "int"},
		Dst: pathtransfer.TransferUnit{Path: "data.pagination.size"},
	})

	transfers = append(transfers, pathtransfer.Transfer{
		Src: pathtransfer.TransferUnit{Path: "func.vocabulary.SetLimit.output.offset", Type: "int"},
		Dst: pathtransfer.TransferUnit{Path: "data.limit.offset"},
	})

	transfers = append(transfers, pathtransfer.Transfer{
		Src: pathtransfer.TransferUnit{Path: "func.vocabulary.SetLimit.output.size", Type: "int"},
		Dst: pathtransfer.TransferUnit{Path: "data.limit.size"},
	})

	var data = `{"data":{"pagination":{"index":1,"size":20}},"user":{"name":"testName"}}`

	out, err := pathtransfer.CallTransferFunc(transfers, []byte(data), func(funcname string, input []byte) (out []byte, err error) {
		index := gjson.GetBytes(input, "index").Int()
		size := gjson.GetBytes(input, "size").Int()
		offset := index * size
		out, err = sjson.SetBytes(out, "offset", offset)
		if err != nil {
			return nil, err
		}
		out, err = sjson.SetBytes(out, "size", size)
		if err != nil {
			return nil, err
		}
		return out, nil
	})
	require.NoError(t, err)
	offset := gjson.Get(string(out), "data.limit.offset").Int()
	require.Equal(t, int64(20), offset)
}

func TestFil(t *testing.T) {
	allVocabularies := `
	func.SetLimit.input.index@int:Dictionary.pagination.index
func.SetLimit.input.size@int:Dictionary.pagination.size
func.SetLimit.output.offset@int:Dictionary.limit.offset
func.SetLimit.output.size@int:Dictionary.limit.size
	`
	limitVocabularies := `
	categorizePagination.input.offset@int:Dictionary.limit.offset
categorizePagination.input.path:Dictionary.education.categorize.path
categorizePagination.input.size@int:Dictionary.limit.size
	`
	allTransfers := pathtransfer.Parse(allVocabularies)
	limitTransfers := pathtransfer.Parse(limitVocabularies)
	funcTransfers := pathtransfer.FilterFuncTransfers(allTransfers, limitTransfers)
	require.Equal(t, 4, len(funcTransfers))
}
