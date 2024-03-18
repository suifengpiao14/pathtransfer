#函数选择需求
定义func[.packageName].funcName.(input/output).argName@type:valueKey 为转换函数协议,如:
```
user.name@string:data.userName// 这个不代表函数入参,因为没有func. 开头,也没有.input/output. 区别是输入还是输出
func.vocabulary.SetLimit.input.index@int:data.pagination.index //代表SetLimit函数的入参index 取m['data.pagination.index'] 的值并转为整数
func.vocabulary.SetLimit.input.size@int:data.pagination.size //代表SetLimit函数的入参size 取m['data.pagination.size'] 的值并转为整数
func.vocabulary.SetLimit.output.offset@int:data.limit.offset //代表SetLimit函数的出参offset 设置为m['data.limit.offset'] 的值
func.vocabulary.SetLimit.output.size@int:data.limit.size //代表SetLimit函数的出参size 设置为m['data.limit.size'] 的值
	`
```

函数协议解析后获得结构体 Transfers
```
type Transfers []Transfer
type Transfer struct {
	Src TransferUnit `json:"src"`
	Dst TransferUnit `json:"dst"`
}
type TransferUnit struct {
	Path string `json:"path"`
	Type string `json:"type"`
}

```


案例go代码:
```
package vocabulary
var data=map[string]any{
    "user.name":"testName",
    "data.pagination.index":1,
    "data.pagination.size":20
}
func SetLimit(index int,size int)(limitOffset int,limitSize int){
    limitOffset=index*size
    limitSize=size
    return 
}
offst,size:=SetLimit(1,20)
data["data.limit.offset"]=offst
data["data.limit.size"]=size
```

假设已知道函数转换协议:
```
user.name@string:data.userName// 这个不代表函数入参,因为没有func. 开头,也没有.input/output. 区别是输入还是输出
func.vocabulary.SetLimit.input.index@int:data.pagination.index //代表SetLimit函数的入参index 取m['data.pagination.index'] 的值并转为整数
func.vocabulary.SetLimit.input.size@int:data.pagination.size //代表SetLimit函数的入参size 取m['data.pagination.size'] 的值并转为整数
func.vocabulary.SetLimit.output.offset@int:data.limit.offset //代表SetLimit函数的出参offset 设置为m['data.limit.offset'] 的值
func.vocabulary.SetLimit.output.size@int:data.limit.size //代表SetLimit函数的出参size 设置为m['data.limit.size'] 的值
func.vocabulary.TrimName.input.name@string:data.userName //代表TrimName函数的入参name取m["data.userName"]的值
func.vocabulary.TrimName.output.name@string:data.userName //代表TrimName函数的出参name设置为取m["data.userName"]的值
	`
```
数据
```
var data=map[string]any{
    "user.name":"testName",
    "data.pagination.index":1,
    "data.pagination.size":20
}
```
需要在data 中增加数据 data.limit.offset和data.limit.size 的值
请实现函数todo 部分,使得TestGetTransferFuncname 测试通过
```
func GetTransferFuncname(transfers Transfers,data map[string]any,dstKeys []string)(funcName string){
    //todo 实现选函数逻辑
    return 
}

func TestGetTransferFuncname(t *testing.T){


transfers:=make(Transfers,0)

transfers=append(transfers,Transfer{
    Src:TransferUnit{Path:"func.vocabulary.TrimName.input.name",Type:"string",},
    Dst:TransferUnit{Path:"data.userName",}
})

transfers=append(transfers,Transfer{
    Src:TransferUnit{Path:"func.vocabulary.TrimName.output.name",Type:"string",},
    Dst:TransferUnit{Path:"data.userName",}
})

transfers=append(transfers,Transfer{
    Src:TransferUnit{Path:"user.name",Type:"string",},
    Dst:TransferUnit{Path:"data.userName",}
})

transfers=append(transfers,Transfer{
    Src:TransferUnit{Path:"func.vocabulary.SetLimit.input.index",Type:"int",},
    Dst:TransferUnit{Path:"data.pagination.index",}
})

transfers=append(transfers,Transfer{
    Src:TransferUnit{Path:"func.vocabulary.SetLimit.input.size",Type:"int",},
    Dst:TransferUnit{Path:"data.pagination.size",}
})

transfers=append(transfers,Transfer{
    Src:TransferUnit{Path:"func.vocabulary.SetLimit.output.index",Type:"int",},
    Dst:TransferUnit{Path:"data.limit.index",}
})

transfers=append(transfers,Transfer{
    Src:TransferUnit{Path:"func.vocabulary.SetLimit.output.size",Type:"int",},
    Dst:TransferUnit{Path:"data.limit.size",}
})



var data=map[string]any{
    "user.name":"testName",
    "data.pagination.index":1,
    "data.pagination.size":20
}
dstKey:=[]string{
    "data.limit.offset",
    "data.limit.size",
}
funcName:=GetTransferFuncname(transfers,data)
require.Eq("SetLimit",funcName)

}



```

