package pathtransfer

import (
	"bytes"
	"encoding/json"
	"fmt"
	"reflect"
	"sort"
	"strings"

	"github.com/pkg/errors"
	"github.com/spf13/cast"
	"github.com/suifengpiao14/funcs"
	"github.com/suifengpiao14/gjsonmodifier"
	"github.com/tidwall/gjson"
)

type Path string

func (path Path) String() string {
	return string(path)
}
func (path Path) EqualFold(path2 Path) bool {
	return strings.EqualFold(path.String(), path2.String())
}

func (path Path) HasNamespace(namespace string) bool {
	return strings.HasPrefix(path.String(), namespace)
}

func (path Path) IsIn() bool {
	return isIn(path)
}

func (path Path) IsOut() bool {
	return isOut(path)
}

func (path Path) TrimNamespace(namespace string) (newPath Path) {
	newPath = Path(strings.TrimPrefix(strings.TrimPrefix(path.String(), namespace), "."))
	return newPath
}

func (path Path) SplitByIO() (namespace string, localName string) {
	delim := ""
	if path.IsIn() {
		delim = Transfer_Direction_input
	} else if path.IsOut() {
		delim = Transfer_Direction_output
	}
	if delim == "" {
		return "", path.String()
	}
	arr := strings.SplitN(path.String(), delim, 2)
	namespace, localName = strings.Trim(arr[0], "."), strings.Trim(arr[1], ".")
	return namespace, localName
}

// TrimIONamespace 剔除输入/输出命名空间后的名称
func (path Path) TrimIONamespace() (localName string) {
	_, localName = path.SplitByIO()

	return localName
}

type TransferUnit struct {
	Path Path   `json:"path"`
	Type string `json:"type"`
}

func (tu TransferUnit) String() string {
	var w bytes.Buffer
	w.WriteString(tu.Path.String())
	if tu.Type != "" {
		w.WriteString(fmt.Sprintf("@%s", tu.Type))
	}
	return w.String()
}

// FuncParameter 解析函数格式路径
func (tu TransferUnit) FuncParameter() (funcParameter *FuncParameter, err error) {
	funcPath := tu.Path
	funcParameter = &FuncParameter{
		Type: tu.Type,
	}
	if !funcPath.HasNamespace(Transfer_Top_Namespace_Func) {
		err = errors.WithMessagef(ERROR_TRANSFER_PATH_NAMESPACE_NOT_FUNC, "func path require prefix:%s,got:%s", Transfer_Top_Namespace_Func, funcPath)
		return nil, err
	}
	if funcPath.IsIn() {
		funcParameter.Direction = Transfer_Direction_input
	} else if funcPath.IsOut() {
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

	funcPath = funcPath.TrimNamespace(Transfer_Top_Namespace_Func)
	funcName, arg := funcPath.SplitByIO()

	funcParameter.FuncName, funcParameter.Name = funcName, arg
	lastDot := strings.LastIndex(funcParameter.FuncName, ".")
	if lastDot > -1 {
		funcParameter.Package, funcParameter.FuncName = funcParameter.FuncName[:lastDot], funcParameter.FuncName[lastDot+1:]
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
	funcParameter.Path = JoinPath(Transfer_Top_Namespace_Func, funcName, funcParameter.Direction, funcParameter.Name) // 剔除name后面的部分，对于对象和原始的funcpath有区别
	return funcParameter, nil
}

const (
	TransferUnit_Type_Int    = "int"
	TransferUnit_Type_String = "string"
)

type Transfer struct {
	Src TransferUnit `json:"src"`
	Dst TransferUnit `json:"dst"`
}

func (t Transfer) String() (s string) {
	var w bytes.Buffer
	w.WriteString(t.Src.String())
	w.WriteString(":")
	w.WriteString(t.Dst.String())
	return w.String()
}
func (t Transfer) IsIn() bool {
	return isIn(t.Src.Path)
}
func (t Transfer) IsOut() bool {
	return isOut(t.Src.Path)
}

// 外界不可以直接初始化,
type Transfers []Transfer

func (a Transfers) Len() int           { return len(a) }
func (a Transfers) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a Transfers) Less(i, j int) bool { return a[i].Src.Path < a[j].Src.Path }

func NewTransfer() (transfer Transfers) {
	return Transfers{}
}

const (
	Transfer_Direction_input  = ".input"  //函数入参
	Transfer_Direction_output = ".output" //函数出参
)

type InOUt interface {
	IsIn() bool
	IsOut() bool
}

func isIn(path Path) bool {
	return strings.Contains(path.String(), Transfer_Direction_input)
}

func isOut(path Path) bool {
	return strings.Contains(path.String(), Transfer_Direction_output)
}

// SplitInOut 分割出入/出参数转换关系
func (transfers Transfers) SplitInOut() (in Transfers, out Transfers) {
	in, out = make(Transfers, 0), make(Transfers, 0)
	for _, t := range transfers {
		if t.IsIn() {
			in.AddReplace(t)
		} else if t.IsOut() {
			out.AddReplace(t)
		}
	}
	return in, out
}

// GetAllDst 获取所有的dst path(筛选函数场景有使用,将目标transfers的dst提取出来，看在func transfers 内是否存在,从而确定转换函数)
func (transfers Transfers) GetAllDst() (dsts []Path) {
	dsts = make([]Path, 0)
	m := map[Path]struct{}{}
	for _, t := range transfers {
		path := t.Dst.Path
		if _, ok := m[path]; ok {
			continue
		}
		dsts = append(dsts, path)
	}
	return dsts
}

// 新增，存在替换
func (transfer *Transfers) AddReplace(transferItems ...Transfer) {
	for _, transferItem := range transferItems {
		exists := false
		for i, item := range *transfer {
			if strings.EqualFold(item.String(), transferItem.String()) {
				(*transfer)[i] = transferItem
				exists = true
				break
			}
		}
		if !exists {
			*transfer = append(*transfer, transferItem)
		}
	}
}

func (transfer Transfers) Reverse() (reversedTransfer Transfers) {
	reversedTransfer = Transfers{}
	for _, item := range transfer {
		refersedItem := Transfer{
			Src: item.Dst,
			Dst: item.Src,
		}
		reversedTransfer = append(reversedTransfer, refersedItem)
	}
	return reversedTransfer
}

// appendTypeToPath 在来源路径上增加上目标类型转换函数
func (t Transfers) appendTypeToPath() (newT Transfers) {
	newT = make(Transfers, 0)
	for _, transfer := range t {
		if transfer.Src.Path == "" { // 路径为空,使用当前数据(如 userTotal.output  去除命名空间后为空,实际数据库返回也是一个整形,没有key)
			transfer.Src.Path = Path("@this")
		}
		transferFunc, ok := DefaultTransferTypes.GetByType(transfer.Dst.Type)
		if ok {
			transfer.Src.Path = Path(fmt.Sprintf("%s%s", transfer.Src.Path.String(), transferFunc.ConvertFn)) //存在映射函数,则修改,否则保持原样
		}
		newT = append(newT, transfer)
	}

	return newT

}

func (t Transfers) Sort() {
	sort.Sort(t)
}
func (t Transfers) Marshal() (tjson string, err error) {
	b, err := json.Marshal(t)
	if err != nil {
		return "", err
	}
	tjson = string(b)
	return tjson, nil

}

// FilterBySrc 通过srcpath 过滤
func (ts Transfers) FilterBySrc(srcPaths ...Path) (subTransfers Transfers) {
	subTransfers = make(Transfers, 0)
	for _, srcPath := range srcPaths {
		for _, t := range ts {
			if t.Src.Path.EqualFold(srcPath) {
				subTransfers.AddReplace(t)
				break
			}
		}
	}

	return subTransfers
}

// FilterByDst 通过srcpath 过滤(已知目标词典path,找转换函数时会用到)
func (ts Transfers) FilterByDst(dstPaths ...Path) (subTransfers Transfers) {
	subTransfers = make(Transfers, 0)
	for _, dstPath := range dstPaths {
		for _, t := range ts {
			if t.Dst.Path.EqualFold(dstPath) {
				subTransfers.AddReplace(t)
				break
			}
		}
	}

	return subTransfers
}

// GetSrcNamespace 获取所有命名空间 delim 一般为.input.|.output.
func (ts Transfers) GetSrcNamespace(delim string) (namespaces []string) {
	namespaces = make([]string, 0)
	m := make(map[string]struct{})
	for _, t := range ts {
		if index := strings.Index(t.Src.Path.String(), delim); index > -1 {
			namespace := t.Src.Path[:index]
			m[namespace.String()] = struct{}{}
		}
	}

	for namespace := range m {
		namespaces = append(namespaces, namespace)
	}

	return namespaces

}

type transfersKeys []string

func (tks *transfersKeys) AppendIgnore(key string) { // 存在忽略
	for _, existsKey := range *tks {
		if existsKey == key {
			return
		}

	}
	*tks = append(*tks, key)
}

type transfersModel struct {
	keys transfersKeys
	m    map[string]any
}

func (ts Transfers) GetByNamespace(namespace string) (subTransfer Transfers) {
	namespace = fmt.Sprintf("%s.", strings.TrimRight(namespace, ".")) // 确保以.结尾
	subTransfer = make(Transfers, 0)
	for _, t := range ts {
		if t.Src.Path.HasNamespace(namespace) {
			subTransfer.AddReplace(t)
		}
	}
	return subTransfer
}

func JoinPath(paths ...string) (newPath Path) {
	arr := make([]string, 0)
	for _, path := range paths {
		if path == "" {
			continue
		}
		arr = append(arr, strings.Trim(path, "."))
	}
	newPath = Path(strings.Join(arr, "."))
	return newPath
}

func (t Transfers) GjsonPath() (gjsonPath string) {
	newT := t.appendTypeToPath()
	m := &transfersModel{
		keys: make([]string, 0),
		m:    make(map[string]any),
	}
	if len(newT) == 0 {
		return ""
	}
	if len(newT) == 1 && newT[0].Dst.Path == "" { // 后续代码默认为对象，在开头增加 . 如只有一个，则不可默认，源字符串输出即可
		return newT[0].Src.Path.String()
	}
	for _, item := range newT {
		dst := item.Dst
		dstPath := strings.ReplaceAll(dst.Path.String(), "@this.", "") // 目标地址 @this. 删除
		dstPath = strings.TrimPrefix(dstPath, ".")
		if !strings.HasPrefix(dstPath, "#") {
			dstPath = fmt.Sprintf(".%s", dstPath) // 非数组，统一标准化前缀
		}

		arr := strings.Split(dstPath, ".")
		l := len(arr)
		ref := m
		for i, key := range arr {
			if l == i+1 { // 处理最后一个
				if (*ref).m[key] == nil {
					(*ref).keys.AppendIgnore(key)
					(*ref).m[key] = item.Src.Path // 第一次默认设置为字符串类型, 如果已经存在,不再修改成字符串(//当类型为 object,array 的在后面,之前有子元素时,忽略)
				}

				continue
			}
			var ok bool
			if _, ok = (*ref).m[key]; !ok {
				(*ref).keys.AppendIgnore(key)
				(*ref).m[key] = &transfersModel{
					keys: make([]string, 0),
					m:    make(map[string]any),
				}
			}
			if ok {
				_, ok = (*ref).m[key].(*transfersModel) //检验类型( //当类型为 object,array 的在前面先设置时 (fullname=items, type=array )其类型不为map)
			}
			if !ok {
				(*ref).keys.AppendIgnore(key)
				(*ref).m[key] = &transfersModel{
					keys: make([]string, 0),
					m:    make(map[string]any),
				}
			}
			ref = (*ref).m[key].(*transfersModel) // 本次递进一定成功
		}

	}
	w, _ := t.recursionWrite(m, false, 0)
	gjsonPath = w.String()

	return gjsonPath
}

// 生成路径
func (t Transfers) recursionWrite(m *transfersModel, parentIsArray bool, depth int) (w bytes.Buffer, childrenIsArray bool) {
	writeComma := false
	for _, k := range m.keys {
		v := (*m).m[k]
		if writeComma {
			w.WriteString(",")
		}
		writeComma = true
		ref, ok := v.(*transfersModel)
		if !ok {
			switch k {
			case "#":
				childrenIsArray = true
				w.WriteString(cast.ToString(v))
			case "":
				w.WriteString(cast.ToString(v))

			default:
				w.WriteString(fmt.Sprintf("%s:%s", k, cast.ToString(v)))
			}
			continue
		}
		var subw bytes.Buffer
		currentIsArray := k == "#"
		if currentIsArray {
			depth++
		}
		subw, subChildrenIsArray := t.recursionWrite(ref, currentIsArray, depth) //isWrapBraces 必须使用外出定义,才能返回true到上一个函数
		subwKey := subw.String()
		if !subChildrenIsArray { //不会被{}包裹,则使用{} 将子内容包裹，表示对象整体(@group 执行后会自动生成{},此处要排除这种情况)
			subwKey = fmt.Sprintf("{%s}", subwKey)
			if parentIsArray {
				subwKey = fmt.Sprintf("%s|@groupPlus:%d", subwKey, depth-1) // 上一级也为数组时，需要包裹到[]中
			}
		}
		var subStr string
		switch k {
		case "#":
			childrenIsArray = true
			subStr = fmt.Sprintf("%s|@groupPlus:%d", subwKey, depth-1)
		case "":
			subStr = subwKey
		default:
			subStr = fmt.Sprintf("%s:%s", k, subwKey)
		}
		w.WriteString(subStr)
	}
	return w, childrenIsArray
}

// PathModifyFn 路径修改函数
type PathModifyFn func(path Path) (newPath Path)

// PathModifyFnSmallCameCase 将路径改成小驼峰格式
func PathModifyFnSmallCameCase(path Path) (newPath Path) {
	arr := strings.Split(path.String(), ".")
	l := len(arr)
	newArr := make([]string, l)
	for i := 0; i < l; i++ {
		if arr[i] == "#" {
			newArr[i] = arr[i]
		} else {
			newArr[i] = funcs.CamelCase(arr[i], false, false)
		}
	}
	newPath = Path(strings.Join(newArr, "."))
	return
}

// PathModifyFnSnakeCase 将路径转为下划线格式
func PathModifyFnSnakeCase(path Path) (newPath Path) {
	arr := strings.Split(path.String(), ".")
	l := len(arr)
	newArr := make([]string, l)
	for i := 0; i < l; i++ {
		newArr[i] = funcs.SnakeCase(arr[i])
	}
	newPath = Path(strings.Join(newArr, "."))
	return
}

// PathModifyFnLower 将路径转为小写格式
func PathModifyFnLower(path Path) (newPath Path) {
	return Path(strings.ToLower(path.String()))
}

// PathModifyFnString 路径后面增加@tostring
func PathModifyFnString(path Path) (newPath Path) {
	return Path(fmt.Sprintf("%s.@tostring", path.String()))
}

// PathModifyFnTrimPrefixFn 生成剔除前缀修改函数
func PathModifyFnTrimPrefixFn(prefix Path) (pathModifyFn PathModifyFn) {
	return func(path Path) (newPath Path) {
		return path.TrimNamespace(prefix.String())
	}
}

// ModifyPath 修改转换路径
func (t Transfers) ModifyDstPath(dstPathModifyFns ...PathModifyFn) (nt Transfers) {
	nt = make(Transfers, 0)
	for _, l := range t {
		src := l.Src
		dst := l.Dst
		for _, fn := range dstPathModifyFns {
			if fn != nil {
				dst.Path = fn(dst.Path)
			}

		}
		item := Transfer{
			Src: src,
			Dst: dst,
		}
		nt.AddReplace(item)
	}
	return nt
}
func (t Transfers) ModifySrcPath(srcPathModifyFns ...PathModifyFn) (nt Transfers) {
	nt = make(Transfers, 0)
	for _, l := range t {
		src := l.Src
		dst := l.Dst
		for _, fn := range srcPathModifyFns {
			if fn != nil {
				src.Path = fn(src.Path)
			}
		}
		item := Transfer{
			Src: src,
			Dst: dst,
		}
		nt.AddReplace(item)
	}
	return nt
}

// GetCallFnScript 从transfer中获取调用函数的动态脚本
func (ts Transfers) GetCallFnScript(language string) (callScript string, err error) {
	funcParameters := make(FuncParameters, 0)
	for _, t := range ts {
		funcParameter, err := t.Src.FuncParameter()
		if errors.Is(err, ERROR_TRANSFER_PATH_NAMESPACE_NOT_FUNC) || errors.Is(err, ERROR_TRANSFER_PATH_DIRECTION_MISSING) {
			err = nil
			continue
		}
		if err != nil {
			return "", err
		}
		funcParameters.AddReplace(*funcParameter)

	}

	callFuncs := funcParameters.CallFuncs()
	callScript, err = callFuncs.Script(language)
	if err != nil {
		return "", err
	}
	return callScript, nil
}

func (ts Transfers) String() (s string) {
	var w bytes.Buffer
	for _, t := range ts {
		w.WriteString(t.String())
		w.WriteString("\n")
	}
	return w.String()
}

type TransferType struct {
	Type      string `json:"type"`      // 对应类型
	ConvertFn string `json:"convertFn"` // 转换函数名称
}
type TransferTypes []TransferType

func (ts TransferTypes) GetByType(typ string) (t *TransferType, ok bool) {
	for _, transfer := range ts {
		if strings.EqualFold(transfer.Type, typ) {
			return &transfer, true
		}
	}
	return nil, false
}

// DefaultTransferTypes schema format 转类型
var DefaultTransferTypes = TransferTypes{
	{Type: "int", ConvertFn: ".@tonum"},
	{Type: "integer", ConvertFn: ".@tonum"},
	{Type: "number", ConvertFn: ".@tonum"},
	{Type: "float", ConvertFn: ".@tonum"},
	{Type: "bool", ConvertFn: ".@tobool"},
	{Type: "boolean", ConvertFn: ".@tobool"},
	{Type: "string", ConvertFn: ".@tostring"},
}

// ToGoTypeTransfer 根据go结构体json tag以及类型生成转换
func ToGoTypeTransfer(dst any) (lineschemaTransfer Transfers) {
	if dst == nil {
		return nil
	}
	rv := reflect.Indirect(reflect.ValueOf(dst))
	rt := rv.Type()
	return toGoTypeTransfer(rt, "@this")
}

func toGoTypeTransfer(rt reflect.Type, prefix Path) (lineschemaTransfer Transfers) {
	switch rt.Kind() {
	case reflect.Array, reflect.Slice:
		lineschemaTransfer = toGoTypeTransfer(rt.Elem(), Path(fmt.Sprintf("%s.#", prefix)))
	case reflect.Struct:
		lineschemaTransfer = str2StructTransfer(rt, prefix)
	case reflect.Int64, reflect.Float64, reflect.Int:
		lineschemaTransfer = str2SimpleTypeTransfer("number", prefix)
	case reflect.Bool:
		lineschemaTransfer = str2SimpleTypeTransfer("bool", prefix)
	case reflect.String:
		lineschemaTransfer = str2SimpleTypeTransfer("string", prefix)
	}

	for i := range lineschemaTransfer {
		t := &lineschemaTransfer[i]
		// 删除前缀 @this
		t.Dst.Path = Path(strings.TrimPrefix(t.Dst.Path.String(), "@this"))
	}

	return lineschemaTransfer
}

func str2SimpleTypeTransfer(typ string, path Path) (lineschemaTransfer Transfers) {
	if path == "" {
		path = "@this"
	}
	return Transfers{
		Transfer{
			Dst: TransferUnit{
				Path: path,
				Type: typ,
			},
			Src: TransferUnit{
				Path: path,
				Type: "string",
			},
		},
	}
}

func str2StructTransfer(rt reflect.Type, pPrefix Path) (transfers Transfers) {
	if rt.Kind() != reflect.Struct {
		return nil
	}
	prefix := pPrefix.String()
	if prefix != "" {
		prefix = strings.TrimRight(prefix, ".")
		prefix = fmt.Sprintf("%s.", prefix)
	}
	transfers = make(Transfers, 0)
	for i := 0; i < rt.NumField(); i++ {
		field := rt.Field(i)
		typ := field.Type.String()
		tag := field.Tag.Get("json")
		if tag == "-" {
			continue // Skip fields without json tag or with "-" tag
		}

		isString := strings.Contains(tag, ",string")
		if isString {
			typ = "string"
		}
		commIndex := strings.Index(tag, ",")
		if commIndex > -1 {
			tag = tag[:commIndex] // 取,前的内容
		}

		fieldType := field.Type
		filedTK := field.Type.Kind()
		switch filedTK {
		case reflect.Slice, reflect.Array, reflect.Struct:
			subPrefix := fmt.Sprintf("%s%s", prefix, tag)
			subTransfer := str2StructTransfer(fieldType, Path(subPrefix))
			transfers.AddReplace(subTransfer...)
			continue // 复合类型，只收集子值
		}
		if tag == "" {
			tag = field.Name // 根据json.Umarsh/Marsh 发现未写json tag时，默认使用列名称，此处兼容保持一致
		}
		path := Path(fmt.Sprintf("%s%s", prefix, tag))
		linschemaT := Transfer{
			Dst: TransferUnit{
				Path: path,
				Type: typ,
			},
			Src: TransferUnit{
				Path: path,
				Type: "string",
			},
		}
		transfers = append(transfers, linschemaT)
	}

	return transfers
}

// RebuildJson 重建json数据 修改json数据的key, 比如下划线修改为小驼峰
func RebuildJson(s string, modifyTransferFn func(transfer Transfers) (newTransfer Transfers)) (newS string, err error) {
	paths := gjsonmodifier.GetAllPath(s)
	transfers := make(Transfers, 0)
	for _, pathStr := range paths {
		path := Path(pathStr)
		transfer := Transfer{
			Src: TransferUnit{
				Path: path,
			},
			Dst: TransferUnit{
				Path: path,
			},
		}
		transfers = append(transfers, transfer)
	}
	if modifyTransferFn != nil {
		transfers = modifyTransferFn(transfers)
	}
	gjsonPath := transfers.String()
	newS = gjson.Get(s, gjsonPath).String()
	return newS, nil
}
