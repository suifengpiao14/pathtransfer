package pathtransfer

import (
	"bytes"
	_ "embed"
	"fmt"
	"strings"
	"text/template"

	"github.com/pkg/errors"
	_ "github.com/suifengpiao14/gjsonmodifier"
	"github.com/tidwall/gjson"
)

const (
	Transfer_Top_Namespace_Func = "func."
)

type CallFunc struct {
	Package  string         `json:"package"`
	Input    FuncParameters `json:"input"`
	Output   FuncParameters `json:"output"`
	FuncName string         `json:"funcName"`
}

type CallFuncs []CallFunc

//go:embed callfunc.tpl
var CallfuncTpl string

func (cfs CallFuncs) Script(language string) (script string, err error) {
	t, err := template.New("").Parse(CallfuncTpl)
	if err != nil {
		return "", err
	}
	var w bytes.Buffer
	tmpl := t.Lookup(language)
	if tmpl == nil {
		err = errors.Errorf("unsport script language:%s", language)
		return "", err
	}
	err = tmpl.Execute(&w, cfs)
	if err != nil {
		return "", err
	}
	return w.String(), nil
}

func (cfs CallFuncs) FirstPackage() (packageName string) {
	for _, cf := range cfs {
		packageName = cf.Package
	}
	return packageName
}

type FuncParameter struct {
	Direction string `json:"direction"` // 标记入参，出参
	Package   string `json:"package"`
	FuncName  string `json:"funcName"`
	Name      string `json:"name"`
	Path      string `json:"path"`
	Type      string `json:"type"`
}

func (fp FuncParameter) String() (s string) {
	s = fp.Path
	if !strings.EqualFold(fp.Type, "object") && !strings.EqualFold(fp.Type, "array") {
		s = fmt.Sprintf("%s@%s", s, fp.Type)
	}
	return s
}

// TypeConvertFunc 类型转换函数
func (fp FuncParameter) TypeConvertFunc() (fnName string) {
	m := map[string]string{
		"int": "ToInt", //使用 cast.XXX 方法
	}
	fnName, ok := m[strings.ToLower(fp.Type)]
	if !ok {
		fnName = "String"
	}
	return fnName
}

func (fp FuncParameter) IsIn() bool {
	return isIn(fp.Path)
}

func (fp FuncParameter) IsOut() bool {
	return isOut(fp.Path)
}

type FuncParameters []FuncParameter

func (fps *FuncParameters) AddReplace(funcParameters ...FuncParameter) {

	for _, fp := range funcParameters {
		exists := false
		for i, fp2 := range *fps {
			if strings.EqualFold(fp2.String(), fp.String()) {
				(*fps)[i] = fp
				exists = true
				break
			}
		}
		if !exists {
			*fps = append(*fps, fp)
		}
	}
}

func (fps FuncParameters) Names() (names string) {
	arr := make([]string, 0)
	for _, fp := range fps {
		arr = append(arr, fp.Name)
	}
	names = strings.Join(arr, ",")
	return names
}

func (fps FuncParameters) First() (fp *FuncParameter, ok bool) {
	for _, fp := range fps {
		return &fp, true
	}
	return nil, false
}

// SplitInOut 分割出入/出参数转换关系
func (fps FuncParameters) SplitInOut() (in FuncParameters, out FuncParameters) {
	in, out = make(FuncParameters, 0), make(FuncParameters, 0)
	for _, t := range fps {
		if t.IsIn() {
			in.AddReplace(t)
		} else if t.IsOut() {
			out.AddReplace(t)
		}
	}
	return in, out
}

func (fps FuncParameters) GroupByFuncName() (groupd map[string]FuncParameters) {
	groupd = map[string]FuncParameters{}
	for _, fp := range fps {
		funcName := fmt.Sprintf("%s.%s", fp.Package, fp.FuncName)
		funcName = strings.TrimPrefix(funcName, ".")
		if _, ok := groupd[funcName]; !ok {
			groupd[funcName] = make(FuncParameters, 0)
		}
		funcPs := groupd[funcName]
		funcPs.AddReplace(fp)
		groupd[funcName] = funcPs
	}
	return groupd
}

func (fps FuncParameters) CallFuncs() (callFuncs CallFuncs) {
	callFuncs = make(CallFuncs, 0)
	gourpd := fps.GroupByFuncName()
	for _, funcParams := range gourpd {
		first, ok := funcParams.First()
		if !ok {
			continue
		}
		callFunc := CallFunc{
			Package:  first.Package,
			FuncName: first.FuncName,
		}
		callFunc.Input, callFunc.Output = funcParams.SplitInOut()
		callFuncs = append(callFuncs, callFunc)
	}
	return
}

var (
	ERROR_TRANSFER_PATH_NAMESPACE_NOT_FUNC = errors.New("not func namespace path")
	ERROR_TRANSFER_PATH_DIRECTION_MISSING  = errors.New("missing direction")
)

// ExplainFuncPath 解析函数格式路径
func ExplainFuncPath(funcPath string) (funcParameter *FuncParameter, err error) {
	funcParameter = &FuncParameter{}
	if !strings.HasPrefix(funcPath, Transfer_Top_Namespace_Func) {
		err = errors.WithMessagef(ERROR_TRANSFER_PATH_NAMESPACE_NOT_FUNC, "func path require prefix:%s,got:%s", Transfer_Top_Namespace_Func, funcPath)
		return nil, err
	}
	if strings.Contains(funcPath, Transfer_Direction_input) {
		funcParameter.Direction = Transfer_Direction_input
	} else if strings.Contains(funcPath, Transfer_Direction_output) {
		funcParameter.Direction = Transfer_Direction_output
	}
	if funcParameter.Direction == "" {
		err = errors.WithMessagef(ERROR_TRANSFER_PATH_DIRECTION_MISSING, "func path format required %s[packageName.]funcName[%s|%s]argName ,got:%s",
			Transfer_Top_Namespace_Func,
			Transfer_Direction_input,
			Transfer_Direction_output,
			funcPath,
		)
		return nil, err
	}

	funcPath = strings.TrimPrefix(funcPath, Transfer_Top_Namespace_Func)
	arr := strings.SplitN(funcPath, funcParameter.Direction, 2)
	funcName, arg := arr[0], arr[1]

	funcParameter.FuncName, funcParameter.Name = funcName, arg
	lastDot := strings.LastIndex(funcParameter.FuncName, ".")
	if lastDot > -1 {
		funcParameter.Package, funcParameter.FuncName = funcParameter.FuncName[:lastDot], funcParameter.FuncName[lastDot+1:]
	}
	atIndex := typeAtIndex(funcParameter.Name)
	if atIndex > -1 {
		funcParameter.Name, funcParameter.Type = funcParameter.Name[:atIndex], funcParameter.Name[atIndex+1:]
	}
	firstDot := strings.Index(funcParameter.Name, ".")
	if firstDot > -1 {
		funcParameter.Name = funcParameter.Name[:firstDot] // 保留第一层
		funcParameter.Type = "object"
	}
	if strings.HasSuffix(funcParameter.Name, "#") {
		funcParameter.Name = strings.TrimSuffix(funcParameter.Name, "#") // 删除结尾的#
		funcParameter.Type = "array"
	}
	funcParameter.Path = fmt.Sprintf("%s%s%s%s", Transfer_Top_Namespace_Func, funcName, funcParameter.Direction, funcParameter.Name) // 剔除name后面的部分，对于对象和原始的funcpath有区别
	return funcParameter, nil
}

var (
	ERROR_TRANSFER_FUNC_NAME_NOT_FOUND = errors.New("not found transfer func name")
)

// GetTransferFuncname 根据输入数据,以及目标key路径,从transfers中选者合适的函数,返回函数名
func GetTransferFuncname(transfers Transfers, data string, dstKeys []string) (funcName string, err error) {
	transfers = transfers.GetByNamespace(Transfer_Top_Namespace_Func)
	dstTransfers := transfers.FilterByDst(dstKeys...)
	funcNames := dstTransfers.GetSrcNamespace(Transfer_Direction_output)
	if len(funcNames) == 0 { // 没有函数名,说明本次无需转换
		return "", nil
	}
	for _, funcName := range funcNames {
		inputNamespace := JoinPath(funcName, Transfer_Direction_input)
		inputTransfers := transfers.GetByNamespace(inputNamespace)
		allInputKeyExist := true
		for _, t := range inputTransfers {
			if !gjson.Get(data, t.Dst.Path).Exists() {
				allInputKeyExist = false
				break
			}
		}
		if allInputKeyExist {
			funcName = strings.TrimPrefix(funcName, Transfer_Top_Namespace_Func)
			return funcName, nil
		}
	}
	//函数名不为空,但是找不到转换函数，则报错
	err = errors.WithMessagef(ERROR_TRANSFER_FUNC_NAME_NOT_FOUND, "funcNames:%s,dstKeys:%s,input:%s", strings.Join(funcNames, ","), strings.Join(dstKeys, ","), data)
	return "", err
}
