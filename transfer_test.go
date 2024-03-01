package transfer_test

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/suifengpiao14/lineschema"
	"github.com/tidwall/gjson"
)

type user struct {
	Name   string `json:"name"`
	UserId int    `json:"userId"`
}

func TestToGoTypeTransfer(t *testing.T) {

	t.Run("nomal", func(t *testing.T) {
		lschemaRaw := `
	version=http://json-schema.org/draft-07/schema#,id=out,direction=out
	fullname=code,format=int,description=业务状态码,comment=业务状态码,example=0
	fullname=message,description=业务提示,comment=业务提示,example=ok
	fullname=items,type=array,description=数组,comment=数组,example=-
	fullname=items[].id,format=int,description=主键,comment=主键,example=0
	fullname=items[].title,description=广告标题,comment=广告标题,example=新年豪礼
	fullname=items[].windowIds[],format=int,description=窗口Id集合,comment=窗口Id集合,example=[1,23,4]
	fullname=items[].windowIds1[],type=int,format=int,description=窗口Id集合,comment=窗口Id集合,example=[1,23,4]
	fullname=items[].windowIds2,type=array,format=string,description=窗口Id集合,comment=窗口Id集合,example=[1,23,4]
	fullname=pagination,type=object,description=对象,comment=对象
	fullname=pagination.index,format=int,description=页索引,0开始,comment=页索引,0开始,example=0
	fullname=pagination.size,format=int,description=每页数量,comment=每页数量,example=10
	fullname=pagination.total,format=int,description=总数,comment=总数,example=60
	`
		lschema, err := lineschema.ParseLineschema(lschemaRaw)
		require.NoError(t, err)
		pathMap := lschema.TransferToFormat().String()
		expected := `{code:code.@tonum,message:message.@tostring,items:{id:items.#.id.@tonum,title:items.#.title.@tostring,windowIds:items.#.windowIds.#.@tonum,windowIds1:items.#.windowIds1.#.@tonum,windowIds2:items.#.windowIds2.#.@tostring}|@groupPlus:0,pagination:{index:pagination.index.@tonum,size:pagination.size.@tonum,total:pagination.total.@tonum}}`
		assert.Equal(t, expected, pathMap)

		input := `{"code":200,"message":"ok","items":[{"id":1,"title":"test1","windowIds":[1,2,3],"windowIds1":[1,2,3],"windowIds2":[1,2,3]},{"id":2,"title":"test2","windowIds":[4,5,6],"windowIds1":[4,5,6]},"windowIds2":[4,5,6]}],"pagination":{"index":0,"size":10,"total":100}}`
		_ = input
		// out := gjson.Get(input, pathMap).String()
		// fmt.Println(out)
	})

	t.Run("struct", func(t *testing.T) {
		lineSchema := lineschema.ToGoTypeTransfer(new(user)).String()
		expected := `{name:@this.name.@tostring,userId:@this.userId.@tonum}`
		assert.Equal(t, expected, lineSchema)
	})

	t.Run("slice[struct]", func(t *testing.T) {
		users := make([]user, 0)
		lineSchema := lineschema.ToGoTypeTransfer(users).String()
		expected := `{name:@this.#.name.@tostring,userId:@this.#.userId.@tonum}|@groupPlus:0`
		assert.Equal(t, expected, lineSchema)
	})
	t.Run("array[struct]", func(t *testing.T) {
		users := [2]user{}
		lineSchema := lineschema.ToGoTypeTransfer(users).String()
		expected := `{name:@this.#.name.@tostring,userId:@this.#.userId.@tonum}|@groupPlus:0`
		assert.Equal(t, expected, lineSchema)
	})

	t.Run("array[int]", func(t *testing.T) {
		ids := [2]string{}
		lineSchema := lineschema.ToGoTypeTransfer(ids).String()
		expected := `@this.#.@tostring`
		assert.Equal(t, expected, lineSchema)
	})

	t.Run("int", func(t *testing.T) {
		id := 2
		lineSchema := lineschema.ToGoTypeTransfer(id).String()
		expected := `@this.@tonum`
		assert.Equal(t, expected, lineSchema)
	})

	t.Run("complex", func(t *testing.T) {
		packschema := `version=http://json-schema.org/draft-07/schema#,id=out
fullname=code,format=int,required,title=业务状态码,default=0,comment=业务状态码,example=0
fullname=message,required,title=业务提示,default=ok,comment=业务提示,example=ok
fullname=items[].id,format=int,required,title=主键,comment=主键,example=1
fullname=items,type=array,title=-,comment=-
fullname=items[].name,required,title=项目标识,comment=项目标识,example=advertise
fullname=items[].title,required,title=名称,comment=名称
fullname=items[].config,required,title=项目curd配置内容,comment=项目curd配置内容
fullname=items[].createdAt,format=datetime,required,title=创建时间,comment=创建时间,example=2023-01-1200:00:00
fullname=items[].updatedAt,format=datetime,required,title=修改时间,comment=修改时间,example=2023-01-3000:00:00
fullname=pagination.index,format=int,required,title=页索引,0开始,default=0,comment=页索引,0开始,example=0
fullname=pagination.size,format=int,required,title=每页数量,default=10,comment=每页数量,example=10
fullname=pagination.total,format=int,required,title=总数,comment=总数,example=60`
		lschema, err := lineschema.ParseLineschema(packschema)
		require.NoError(t, err)
		gjsonPath := lschema.TransferToFormat().Reverse().String()
		expected := `{code:code.@tostring,message:message.@tostring,items:{id:items.#.id.@tostring,name:items.#.name.@tostring,title:items.#.title.@tostring,config:items.#.config.@tostring,createdAt:items.#.createdAt.@tostring,updatedAt:items.#.updatedAt.@tostring}|@groupPlus:0,pagination:{index:pagination.index.@tostring,size:pagination.size.@tostring,total:pagination.total.@tostring}}`
		assert.Equal(t, expected, gjsonPath)

		data := `{"code":0,"message":"","items":[{"id":2,"name":"advertise1","title":"广aa告","config":"{\"navs\":[1]}","createdAt":"","updatedAt":""}],"pagination":{"index":0,"size":10,"total":1}}`
		_ = data
		// out := gjson.Get(data, gjsonPath).String()
		// fmt.Println(out)
	})

	t.Run("complex2", func(t *testing.T) {
		packschema := `version=http://json-schema.org/draft-07/schema#,id=out
fullname=code,format=int,required,title=业务状态码,default=0,comment=业务状态码,example=0
fullname=message,required,title=业务提示,default=ok,comment=业务提示,example=ok
fullname=navs[].id,format=int,required,title=主键,comment=主键
fullname=navs[].name,required,title=名称,comment=名称
fullname=navs[].title,required,title=标题,comment=标题
fullname=navs[].route,required,title=路由,comment=路由
fullname=navs[].sort,format=int,required,title=排序,comment=排序`
		lschema, err := lineschema.ParseLineschema(packschema)
		require.NoError(t, err)
		gjsonPath := lschema.TransferToFormat().Reverse().String()

		expected := `{code:code.@tostring,message:message.@tostring,navs:{id:navs.#.id.@tostring,name:navs.#.name.@tostring,title:navs.#.title.@tostring,route:navs.#.route.@tostring,sort:navs.#.sort.@tostring}|@groupPlus:0}`
		assert.Equal(t, expected, gjsonPath)

		data := `{"code":0,"message":"","navs":[{"id":1,"name":"creative","title":"广告创意","route":"creativeList","sort":99},{"id":2,"name":"plan","title":"广告计划","route":"planList","sort":98},{"id":3,"name":"window","title":"橱窗","route":"windowList","sort":97}]}`
		_ = data
		// out := gjson.Get(data, gjsonPath).String()
		// fmt.Println(out)

	})
	t.Run("object with no children", func(t *testing.T) {
		packschema := `version=http://json-schema.org/draft-07/schema#,id=out
fullname=code,format=int,required,title=业务状态码,default=0,comment=业务状态码,example=0
fullname=message,required,title=业务提示,default=ok,comment=业务提示,example=ok
fullname=uiSchema,type=object,required,title=uiSchema对象,comment=uiSchema对象`
		lschema, err := lineschema.ParseLineschema(packschema)
		require.NoError(t, err)
		gjsonPath := lschema.TransferToFormat().Reverse().String()
		expected := `{code:code.@tostring,message:message.@tostring,uiSchema:uiSchema}`
		assert.Equal(t, expected, gjsonPath)

		data := `{"code":0,"message":"","uiSchema":""}`
		_ = data
		// out := gjson.Get(data, gjsonPath).String()
		// fmt.Println(out)
	})
	t.Run("deep array", func(t *testing.T) {
		packschema := `version=http://json-schema.org/draft-07/schema#,id=out
		fullname=services[].name,required,title=项目标识,comment=项目标识,example=advertise
fullname=services[].servers[].name,required,title=服务标识,comment=服务标识,example=dev
fullname=services[].servers[].title,required,title=服务名称,comment=服务名称,example=dev
`
		lschema, err := lineschema.ParseLineschema(packschema)
		require.NoError(t, err)
		gjsonPath := lschema.TransferToFormat().Reverse().String()
		expected := `{services:{name:services.#.name.@tostring,servers:{name:services.#.servers.#.name.@tostring,title:services.#.servers.#.title.@tostring}|@groupPlus:1}|@groupPlus:0}`
		assert.Equal(t, expected, gjsonPath)

		data := `{"code":0,"message":"","services":[{"id":1,"name":"advertise","title":"广告服务","createdAt":"2023-11-25 22:32:16","updatedAt":"2023-11-25 22:32:16","servers":[{"name":"dev","title":"广告服务开发环境"},{"name":"dev2","title":"广告服务开发环境"}]}],"pagination":{"index":0,"size":10,"total":1}}`
		out := gjson.Get(data, gjsonPath).String()
		excepted := `{"services":[{"name":"advertise","servers":[{"name":"dev","title":"广告服务开发环境"},{"name":"dev2","title":"广告服务开发环境"}]}]}`
		assert.Equal(t, excepted, out)
	})

	t.Run("deep array2 ", func(t *testing.T) {
		packschema := `version=http://json-schema.org/draft-07/schema#,id=out
		fullname=services[].name,required,title=项目标识,comment=项目标识,example=advertise
fullname=services[].serverIds[],required,format=int,title=服务标识,comment=服务标识,example=dev
`
		lschema, err := lineschema.ParseLineschema(packschema)
		require.NoError(t, err)
		gjsonPath := lschema.TransferToFormat().Reverse().String()
		expected := `{services:{name:services.#.name.@tostring,serverIds:services.#.serverIds.#.@tostring}|@groupPlus:0}`
		assert.Equal(t, expected, gjsonPath)

		data := `{"code":0,"message":"","services":[{"name":"advertise","serverIds":[1,2,3]}],"pagination":{"index":0,"size":10,"total":1}}`
		_ = data
		// out := gjson.Get(data, gjsonPath).String()
		// fmt.Println(out)
	})

	t.Run("complexe 3", func(t *testing.T) {
		packschema := `version=http://json-schema.org/draft-07/schema#,id=out
fullname=code,format=int,required,title=业务状态码,default=0,comment=业务状态码,example=0
fullname=message,required,title=业务提示,default=ok,comment=业务提示,example=ok
fullname=services[].id,format=int,required,title=主键,comment=主键,example=1
fullname=services[].name,required,title=项目标识,comment=项目标识,example=advertise
fullname=services[].title,required,title=名称,comment=名称
fullname=services[].document,required,title=是,default=产品文档地址,comment=是
fullname=services[].createdAt,format=datetime,required,title=创建时间,comment=创建时间,example=2023-01-1200:00:00
fullname=services[].updatedAt,format=datetime,required,title=修改时间,comment=修改时间,example=2023-01-3000:00:00
fullname=services[].servers[].name,required,title=服务标识,comment=服务标识,example=dev
fullname=services[].servers[].title,required,title=服务名称,comment=服务名称,example=dev
fullname=services[].navs[].name,required,title=名称,comment=名称
fullname=services[].navs[].title,required,title=标题,comment=标题
fullname=services[].navs[].route,required,title=路由,comment=路由
fullname=services[].navs[].sort,format=int,required,title=排序,comment=排序
fullname=pagination.index,format=int,required,title=页索引,0开始,default=0,comment=页索引,0开始,example=0
fullname=pagination.size,format=int,required,title=每页数量,default=10,comment=每页数量,example=10
fullname=pagination.total,format=int,required,title=总数,comment=总数,example=60`
		lschema, err := lineschema.ParseLineschema(packschema)
		require.NoError(t, err)
		path := lschema.TransferToFormat().Reverse().String()
		expected := `{code:code.@tostring,message:message.@tostring,services:{id:services.#.id.@tostring,name:services.#.name.@tostring,title:services.#.title.@tostring,document:services.#.document.@tostring,createdAt:services.#.createdAt.@tostring,updatedAt:services.#.updatedAt.@tostring,servers:{name:services.#.servers.#.name.@tostring,title:services.#.servers.#.title.@tostring}|@groupPlus:1,navs:{name:services.#.navs.#.name.@tostring,title:services.#.navs.#.title.@tostring,route:services.#.navs.#.route.@tostring,sort:services.#.navs.#.sort.@tostring}|@groupPlus:1}|@groupPlus:0,pagination:{index:pagination.index.@tostring,size:pagination.size.@tostring,total:pagination.total.@tostring}}`
		assert.Equal(t, expected, path)

		data := `{"code":0,"message":"","services":[{"id":6,"name":"advertise1","title":"广告服务","document":"","createdAt":"2023-12-02 23:01:04","updatedAt":"2023-12-02 23:01:04","servers":[],"navs":[]},{"id":1,"name":"advertise","title":"广告服务","document":"","createdAt":"2023-11-25 22:32:16","updatedAt":"2023-11-25 22:32:16","servers":[{"name":"dev","title":"开发环境"},{"name":"prod","title":"开发环境"}],"navs":[{"name":"creative","title":"广告创意","route":"/advertise/creativeList","sort":99},{"name":"plan","title":"广告计划","route":"/advertise/planList","sort":98},{"name":"window","title":"橱窗","route":"/advertise/windowList","sort":97},{"name":"crativeList","title":"广告服务","route":"/creativeList","sort":4}]}],"pagination":{"index":0,"size":10,"total":2}}`
		_ = data
		/* newData := gjson.Get(data, path).String()
		fmt.Println(newData) */
	})

	t.Run("complex_array_object", func(t *testing.T) {
		unpackSchema := `version=http://json-schema.org/draft-07/schema#,id=out
fullname=service.name,required,title=服务名称,comment=服务名称
fullname=service.title,required,title=服务标题,comment=服务标题
fullname=service.document,required,title=服务文档地址,comment=服务文档地址
fullname=servers,type=array,required,title=服务,comment=服务
fullname=servers[].name,required,title=服务名称,comment=服务名称
fullname=servers[].title,required,title=服务标题,comment=服务标题
fullname=servers[].url,required,title=服务地址,comment=服务地址
fullname=servers[].proxy,required,title=服务代理地址,comment=服务代理地址
fullname=servers[].env,required,title=环境变量,comment=环境变量
fullname=navigates,type=array,required,title=导航,comment=导航
fullname=navigates[].name,required,title=导航名称,comment=导航名称
fullname=navigates[].title,required,title=导航标题,comment=导航标题
fullname=navigates[].route,required,title=导航路由,comment=导航路由
fullname=navigates[].sort,format=int,required,title=排序,comment=排序
fullname=dataSchemas,type=array,required,title=页面元素,comment=页面元素
fullname=dataSchemas[].name,required,title=元素名称,comment=元素名称
fullname=dataSchemas[].serviceName,required,title=服务名称,comment=服务名称
fullname=dataSchemas[].navRote,required,title=前端页面路由,comment=前端页面路由
fullname=dataSchemas[].parentNavRoute,required,title=上一级路由,comment=上一级路由
fullname=dataSchemas[].scene,required,title=场景,description=场景(list,create,edit,delete),comment=场景(list,create,edit,delete)
fullname=dataSchemas[].description,required,title=元素描述,comment=元素描述
fullname=dataSchemas[].request[].type,required,title=字段类型,comment=字段类型
fullname=dataSchemas[].request[].title,required,title=字段标签,comment=字段标签
fullname=dataSchemas[].request[].fullname,required,title=字段全称,comment=字段全称
fullname=dataSchemas[].request[].name,required,title=字段名称,comment=字段名称
fullname=dataSchemas[].request[].primaryKey,format=bool,required,title=是否为主键,comment=是否为主键
fullname=dataSchemas[].request[].required,format=bool,required,title=是否必填,comment=是否必填
fullname=dataSchemas[].request[].scene,required,title=使用场景,comment=使用场景
fullname=dataSchemas[].response[].type,required,title=字段类型,comment=字段类型
fullname=dataSchemas[].response[].title,required,title=字段标签,comment=字段标签
fullname=dataSchemas[].response[].fullname,required,title=字段全称,comment=字段全称
fullname=dataSchemas[].response[].name,required,title=字段名称,comment=字段名称
fullname=dataSchemas[].response[].primaryKey,format=bool,required,title=是否为主键,comment=是否为主键
fullname=dataSchemas[].response[].required,format=bool,required,title=是否必填,comment=是否必填
fullname=dataSchemas[].response[].scene,required,title=使用场景,comment=使用场景
fullname=dataSchemas[].action.url,required,title=请求地址,comment=请求地址
fullname=dataSchemas[].action.method,required,title=请求方法,comment=请求方法
fullname=code,format=int,required,title=业务状态码,default=0,comment=业务状态码,example=0
fullname=message,required,title=业务提示,default=ok,comment=业务提示,example=ok`
		data := `{"service":{"name":"advertise","title":"广告服务","document":"http://document.com/ap"},"servers":[{"name":"dev","title":"开发环境","url":"http://ad.micor.cn","proxy":"http://127.0.0.1:8083","env":""},{"name":"test","title":"测试环境","url":"http://ad.micor.cn","proxy":"","env":""},{"name":"prod","title":"正式环境","url":"http://ad.micor.cn","proxy":"","env":""}],"navigates":[{"name":"window","title":"橱窗","route":"/windowList","sort":"97"},{"name":"creative","title":"广告创意","route":"/creativeList","sort":"99"},{"name":"plan","title":"广告计划","route":"/planList","sort":"98"}],"dataSchemas":[{"name":"","serviceName":"advertise","navRote":"/creativeList","parentNavRoute":"","scene":"list","description":"创意列表","request":[{"type":"string","title":"广告计划Id","fullname":"planId","name":"planId","primaryKey":"false","required":"true","scene":""},{"type":"string","title":"名称","fullname":"name","name":"name","primaryKey":"false","required":"true","scene":""},{"type":"string","title":"页索引","fullname":"index","name":"index","primaryKey":"false","required":"true","scene":"page"},{"type":"string","title":"每页数量","fullname":"size","name":"size","primaryKey":"false","required":"true","scene":"page"}],"response":[{"type":"string","title":"业务状态码","fullname":"code","name":"code","primaryKey":"false","required":"false","scene":"businessStatus"},{"type":"string","title":"业务提示","fullname":"message","name":"message","primaryKey":"false","required":"false","scene":"businessStatus"},{"type":"string","title":"主键","fullname":"items[].id","name":"id","primaryKey":"false","required":"false","scene":"identify"},{"type":"string","title":"广告计划Id","fullname":"items[].planId","name":"planId","primaryKey":"false","required":"false","scene":""},{"type":"string","title":"名称","fullname":"items[].name","name":"name","primaryKey":"false","required":"false","scene":""},{"type":"string","title":"广告内容","fullname":"items[].content","name":"content","primaryKey":"false","required":"false","scene":""},{"type":"string","title":"创建时间","fullname":"items[].createdAt","name":"createdAt","primaryKey":"false","required":"false","scene":""},{"type":"string","title":"修改时间","fullname":"items[].updatedAt","name":"updatedAt","primaryKey":"false","required":"false","scene":""},{"type":"string","title":"页索引","fullname":"pagination.index","name":"index","primaryKey":"false","required":"false","scene":"page"},{"type":"string","title":"每页数量","fullname":"pagination.size","name":"size","primaryKey":"false","required":"false","scene":"page"},{"type":"string","title":"总数","fullname":"pagination.total","name":"total","primaryKey":"false","required":"false","scene":"page"}],"action":{"url":"/admin/v1/creative/list","method":"POST"}},{"name":"新增","serviceName":"advertise","navRote":"","parentNavRoute":"/creativeList","scene":"crate","description":"新增创意","request":[{"type":"string","title":"广告计划Id","fullname":"planId","name":"planId","primaryKey":"false","required":"true","scene":""},{"type":"string","title":"名称","fullname":"name","name":"name","primaryKey":"false","required":"true","scene":""},{"type":"string","title":"广告内容","fullname":"content","name":"content","primaryKey":"false","required":"true","scene":""}],"response":[{"type":"string","title":"业务状态码","fullname":"code","name":"code","primaryKey":"false","required":"false","scene":"businessStatus"},{"type":"string","title":"业务提示","fullname":"message","name":"message","primaryKey":"false","required":"false","scene":"businessStatus"}],"action":{"url":"/admin/v1/creative/add","method":"POST"}},{"name":"修改","serviceName":"advertise","navRote":"","parentNavRoute":"/creativeList","scene":"edit","description":"更新创意","request":[{"type":"string","title":"主键","fullname":"id","name":"id","primaryKey":"false","required":"true","scene":"identify"},{"type":"string","title":"名称","fullname":"name","name":"name","primaryKey":"false","required":"true","scene":""},{"type":"string","title":"广告内容","fullname":"content","name":"content","primaryKey":"false","required":"true","scene":""}],"response":[{"type":"string","title":"业务状态码","fullname":"code","name":"code","primaryKey":"false","required":"false","scene":"businessStatus"},{"type":"string","title":"业务提示","fullname":"message","name":"message","primaryKey":"false","required":"false","scene":"businessStatus"}],"action":{"url":"/admin/v1/creative/update","method":"POST"}},{"name":"删除","serviceName":"advertise","navRote":"","parentNavRoute":"/creativeList","scene":"delete","description":"删除创意","request":[{"type":"string","title":"主键","fullname":"id","name":"id","primaryKey":"false","required":"true","scene":""},{"type":"string","title":"广告计划Id","fullname":"planId","name":"planId","primaryKey":"false","required":"true","scene":""}],"response":[{"type":"string","title":"业务状态码","fullname":"code","name":"code","primaryKey":"false","required":"false","scene":"businessStatus"},{"type":"string","title":"业务提示","fullname":"message","name":"message","primaryKey":"false","required":"false","scene":"businessStatus"}],"action":{"url":"/admin/v1/creative/del","method":"POST"}}],"code":"0","message":"ok"}`
		_ = data
		lschema, err := lineschema.ParseLineschema(unpackSchema)
		require.NoError(t, err)
		path := lschema.TransferToFormat().Reverse().String()
		exceptedPath := `{service:{name:service.name.@tostring,title:service.title.@tostring,document:service.document.@tostring},servers:{name:servers.#.name.@tostring,title:servers.#.title.@tostring,url:servers.#.url.@tostring,proxy:servers.#.proxy.@tostring,env:servers.#.env.@tostring}|@groupPlus:0,navigates:{name:navigates.#.name.@tostring,title:navigates.#.title.@tostring,route:navigates.#.route.@tostring,sort:navigates.#.sort.@tostring}|@groupPlus:0,dataSchemas:{name:dataSchemas.#.name.@tostring,serviceName:dataSchemas.#.serviceName.@tostring,navRote:dataSchemas.#.navRote.@tostring,parentNavRoute:dataSchemas.#.parentNavRoute.@tostring,scene:dataSchemas.#.scene.@tostring,description:dataSchemas.#.description.@tostring,request:{type:dataSchemas.#.request.#.type.@tostring,title:dataSchemas.#.request.#.title.@tostring,fullname:dataSchemas.#.request.#.fullname.@tostring,name:dataSchemas.#.request.#.name.@tostring,primaryKey:dataSchemas.#.request.#.primaryKey.@tostring,required:dataSchemas.#.request.#.required.@tostring,scene:dataSchemas.#.request.#.scene.@tostring}|@groupPlus:1,response:{type:dataSchemas.#.response.#.type.@tostring,title:dataSchemas.#.response.#.title.@tostring,fullname:dataSchemas.#.response.#.fullname.@tostring,name:dataSchemas.#.response.#.name.@tostring,primaryKey:dataSchemas.#.response.#.primaryKey.@tostring,required:dataSchemas.#.response.#.required.@tostring,scene:dataSchemas.#.response.#.scene.@tostring}|@groupPlus:1,action:{url:dataSchemas.#.action.url.@tostring,method:dataSchemas.#.action.method.@tostring}|@groupPlus:0}|@groupPlus:0,code:code.@tostring,message:message.@tostring}`
		assert.Equal(t, exceptedPath, path)
		/* newData := gjson.Get(data, path).String()
		fmt.Println(newData) */

	})

}

func TestDeepArrWithSimplArr(t *testing.T) {
	gjsonPath := `{services:{name:services.#.name.@tostring,serverIds:services.#.serverIds.#.@tostring}|@group}`
	data := `{"code":0,"message":"","services":[{"name":"advertise","serverIds":[1,2,3]}],"pagination":{"index":0,"size":10,"total":1}}`
	out := gjson.Get(data, gjsonPath).String()
	fmt.Println(out)

}

func TestDeepArray(t *testing.T) {
	jsonStr := `{"code":0,"message":"","services":[{"id":1,"name":"advertise","title":"广告服务","createdAt":"2023-11-25 22:32:16","updatedAt":"2023-11-25 22:32:16","servers":[{"name":"dev","title":"广告服务开发环境"},{"name":"dev2","title":"广告服务开发环境"}]}],"pagination":{"index":0,"size":10,"total":1}}`
	path := `{services:{name:services.#.name.@tostring,servers:[{name:services.#.servers.#.name.@tostring|@flatten,title:services.#.servers.#.title.@tostring|@flatten}|@group}]|@group}`
	newJson := gjson.Get(jsonStr, path).String()

	fmt.Println(newJson)
}
func TestArray(t *testing.T) {
	jsonStr := `[{"name":"test1"}]`
	path := `#.name`
	newJson := gjson.Get(jsonStr, path).String()
	fmt.Println(newJson)
}

func TestStructArrayPath(t *testing.T) {
	jsonStr := `[{"name":"张三","userId":"1"},{"name":"李四","userId":"2"}]`
	path := `[{name:@this.#.name.@tostring,userId:@this.#.userId.@tonum}|@group]`
	newJson := gjson.Get(jsonStr, path).String()
	fmt.Println(newJson)
}

func TestSimpleArrayPath(t *testing.T) {
	jsonStr := `[1,2,3]`
	path := `@this.#.@tostring`
	newJson := gjson.Get(jsonStr, path).String()
	fmt.Println(newJson)
}
func TestValuePath(t *testing.T) {
	jsonStr := `"1"`
	path := `@this.@tonum`
	newJson := gjson.Get(jsonStr, path).String()
	fmt.Println(newJson)
}

type UerNoJsonTag struct {
	Name      string
	ID        int
	CreatedAt string
	Update_at string
}

func TestJsonUmarsh(t *testing.T) {
	u := UerNoJsonTag{}
	data := `{"name":"张三","id":2,"createdAt":"2023-11-24 16:10:00","Update_at":"2023-11-24 16:10:00"}`
	json.Unmarshal([]byte(data), &u)
	b, _ := json.Marshal(u)
	s := string(b)
	fmt.Println(s)
}

func TestTransfer1(t *testing.T) {

	jsonData := `{"code":0,"services":[{"id":6,"servers":[]},{"id":1,"servers":[{"name":"dev","title":"开发环境"},{"name":"prod","title":"开发环境"}]}]}`

	result := gjson.Get(jsonData, "**")
	paths := result.Array()

	for _, path := range paths {
		fmt.Println("Path:", path.String())
	}

}

func TestTransfer2(t *testing.T) {

	jsonData := `{"service":{"name":"advertise","title":"广告服务","document":"http://document.com/ap"},"servers":[{"name":"test","title":"测试环境","url":"http://ad.micor.cn","proxy":"","env":""},{"name":"dev","title":"开发环境","url":"http://ad.micor.cn","proxy":"http://127.0.0.1:8083","env":""},{"name":"prod","title":"正式环境","url":"http://ad.micor.cn","proxy":"","env":""}],"navigates":[{"name":"plan","title":"广告计划","route":"/planList","sort":"98"},{"name":"window","title":"橱窗","route":"/windowList","sort":"97"},{"name":"creative","title":"广告创意","route":"/creativeList","sort":"99"}],"dataSchemas":[{"name":"","serviceName":"advertise","navRote":"/creativeList","parentNavRoute":"","scene":"list","description":"创意列表","request":[{"type":"string","title":"广告计划Id","fullname":"planId","name":"planId","primaryKey":"false","required":"true","scene":""},{"type":"string","title":"名称","fullname":"name","name":"name","primaryKey":"false","required":"true","scene":""},{"type":"string","title":"页索引","fullname":"index","name":"index","primaryKey":"false","required":"true","scene":"page"},{"type":"string","title":"每页数量","fullname":"size","name":"size","primaryKey":"false","required":"true","scene":"page"}],"response":[{"type":"string","title":"业务状态码","fullname":"code","name":"code","primaryKey":"false","required":"false","scene":"businessStatus"},{"type":"string","title":"业务提示","fullname":"message","name":"message","primaryKey":"false","required":"false","scene":"businessStatus"},{"type":"string","title":"主键","fullname":"items[].id","name":"id","primaryKey":"false","required":"false","scene":"identify"},{"type":"string","title":"广告计划Id","fullname":"items[].planId","name":"planId","primaryKey":"false","required":"false","scene":""},{"type":"string","title":"名称","fullname":"items[].name","name":"name","primaryKey":"false","required":"false","scene":""},{"type":"string","title":"广告内容","fullname":"items[].content","name":"content","primaryKey":"false","required":"false","scene":""},{"type":"string","title":"创建时间","fullname":"items[].createdAt","name":"createdAt","primaryKey":"false","required":"false","scene":""},{"type":"string","title":"修改时间","fullname":"items[].updatedAt","name":"updatedAt","primaryKey":"false","required":"false","scene":""},{"type":"string","title":"页索引","fullname":"pagination.index","name":"index","primaryKey":"false","required":"false","scene":"page"},{"type":"string","title":"每页数量","fullname":"pagination.size","name":"size","primaryKey":"false","required":"false","scene":"page"},{"type":"string","title":"总数","fullname":"pagination.total","name":"total","primaryKey":"false","required":"false","scene":"page"}],"action":{"url":"/admin/v1/creative/list","method":"POST"}},{"name":"新增","serviceName":"advertise","navRote":"","parentNavRoute":"/creativeList","scene":"crate","description":"新增创意","request":[{"type":"string","title":"广告计划Id","fullname":"planId","name":"planId","primaryKey":"false","required":"true","scene":""},{"type":"string","title":"名称","fullname":"name","name":"name","primaryKey":"false","required":"true","scene":""},{"type":"string","title":"广告内容","fullname":"content","name":"content","primaryKey":"false","required":"true","scene":""}],"response":[{"type":"string","title":"业务状态码","fullname":"code","name":"code","primaryKey":"false","required":"false","scene":"businessStatus"},{"type":"string","title":"业务提示","fullname":"message","name":"message","primaryKey":"false","required":"false","scene":"businessStatus"}],"action":{"url":"/admin/v1/creative/add","method":"POST"}},{"name":"修改","serviceName":"advertise","navRote":"","parentNavRoute":"/creativeList","scene":"edit","description":"更新创意","request":[{"type":"string","title":"主键","fullname":"id","name":"id","primaryKey":"false","required":"true","scene":"identify"},{"type":"string","title":"名称","fullname":"name","name":"name","primaryKey":"false","required":"true","scene":""},{"type":"string","title":"广告内容","fullname":"content","name":"content","primaryKey":"false","required":"true","scene":""}],"response":[{"type":"string","title":"业务状态码","fullname":"code","name":"code","primaryKey":"false","required":"false","scene":"businessStatus"},{"type":"string","title":"业务提示","fullname":"message","name":"message","primaryKey":"false","required":"false","scene":"businessStatus"}],"action":{"url":"/admin/v1/creative/update","method":"POST"}},{"name":"删除","serviceName":"advertise","navRote":"","parentNavRoute":"/creativeList","scene":"delete","description":"删除创意","request":[{"type":"string","title":"主键","fullname":"id","name":"id","primaryKey":"false","required":"true","scene":""},{"type":"string","title":"广告计划Id","fullname":"planId","name":"planId","primaryKey":"false","required":"true","scene":""}],"response":[{"type":"string","title":"业务状态码","fullname":"code","name":"code","primaryKey":"false","required":"false","scene":"businessStatus"},{"type":"string","title":"业务提示","fullname":"message","name":"message","primaryKey":"false","required":"false","scene":"businessStatus"}],"action":{"url":"/admin/v1/creative/del","method":"POST"}}],"code":"0","message":"ok"}`
	path := `{service:{name:service.name.@tostring,title:service.title.@tostring,document:service.document.@tostring},servers:{name:servers.#.name.@tostring,title:servers.#.title.@tostring,url:servers.#.url.@tostring,proxy:servers.#.proxy.@tostring,env:servers.#.env.@tostring}|@groupPlus:0,navigates:{route:navigates.#.route.@tostring,sort:navigates.#.sort.@tonum,name:navigates.#.name.@tostring,title:navigates.#.title.@tostring}|@groupPlus:0,dataSchemas:{parentNavRoute:dataSchemas.#.parentNavRoute.@tostring,scene:dataSchemas.#.scene.@tostring,request:{scene:dataSchemas.#.request.#.scene.@tostring,type:dataSchemas.#.request.#.type.@tostring,title:dataSchemas.#.request.#.title.@tostring,fullname:dataSchemas.#.request.#.fullname.@tostring,name:dataSchemas.#.request.#.name.@tostring,primaryKey:dataSchemas.#.request.#.primaryKey.@tobool,required:dataSchemas.#.request.#.required.@tobool}|@groupPlus:1,response:{name:dataSchemas.#.response.#.name.@tostring,primaryKey:dataSchemas.#.response.#.primaryKey.@tobool,required:dataSchemas.#.response.#.required.@tobool,scene:dataSchemas.#.response.#.scene.@tostring,type:dataSchemas.#.response.#.type.@tostring,title:dataSchemas.#.response.#.title.@tostring,fullname:dataSchemas.#.response.#.fullname.@tostring}|@groupPlus:1,action:{url:dataSchemas.#.action.url.@tostring,method:dataSchemas.#.action.method.@tostring}|@groupPlus:0,serviceName:dataSchemas.#.serviceName.@tostring,navRote:dataSchemas.#.navRote.@tostring,name:dataSchemas.#.name.@tostring,description:dataSchemas.#.description.@tostring}|@groupPlus:0,code:code.@tonum,message:message.@tostring}`
	//path := `{action:{url:dataSchemas.#.action.url.@tostring,method:dataSchemas.#.action.method.@tostring}}`
	newData := gjson.Get(jsonData, path).String()
	fmt.Println(newData)

}

func TestTransferJson(t *testing.T) {
	jsonStr := `[{"DatabaseConfig":{"databaseName":"ad","tablePrefix":"","columnPrefix":"","deletedAtColumn":"deleted_at","logLevel":"","version":"","extaConfigs":null},"TableName":"creative","PrimaryKey":"id","DeleteColumn":"deleted_at","Columns":[{"Prefix":"","CamelName":"Id","ColumnName":"id","Name":"id","Type":"int","Comment":"主键","Tag":"","Nullable":false,"Enums":[],"AutoIncrement":true,"DefaultValue":"","OnCreate":false,"OnUpdate":false,"OnDelete":false},{"Prefix":"","CamelName":"PlanId","ColumnName":"plan_id","Name":"plan_id","Type":"string","Comment":"广告计划Id","Tag":"","Nullable":false,"Enums":[],"AutoIncrement":false,"DefaultValue":"","OnCreate":false,"OnUpdate":false,"OnDelete":false},{"Prefix":"","CamelName":"Name","ColumnName":"name","Name":"name","Type":"string","Comment":"名称","Tag":"","Nullable":false,"Enums":[],"AutoIncrement":false,"DefaultValue":"","OnCreate":false,"OnUpdate":false,"OnDelete":false},{"Prefix":"","CamelName":"Content","ColumnName":"content","Name":"content","Type":"string","Comment":"广告内容","Tag":"","Nullable":true,"Enums":[],"AutoIncrement":false,"DefaultValue":"","OnCreate":false,"OnUpdate":false,"OnDelete":false},{"Prefix":"","CamelName":"CreatedAt","ColumnName":"created_at","Name":"created_at","Type":"string","Comment":"创建时间","Tag":"","Nullable":true,"Enums":[],"AutoIncrement":false,"DefaultValue":"current_timestamp()","OnCreate":true,"OnUpdate":false,"OnDelete":false},{"Prefix":"","CamelName":"UpdatedAt","ColumnName":"updated_at","Name":"updated_at","Type":"string","Comment":"修改时间","Tag":"","Nullable":true,"Enums":[],"AutoIncrement":false,"DefaultValue":"current_timestamp()","OnCreate":false,"OnUpdate":true,"OnDelete":false},{"Prefix":"","CamelName":"DeletedAt","ColumnName":"deleted_at","Name":"deleted_at","Type":"string","Comment":"删除时间","Tag":"","Nullable":true,"Enums":[],"AutoIncrement":false,"DefaultValue":"NULL","OnCreate":false,"OnUpdate":false,"OnDelete":true}],"EnumsConst":[],"Comment":"广告物料","TableDef":null},{"DatabaseConfig":{"databaseName":"ad","tablePrefix":"","columnPrefix":"","deletedAtColumn":"deleted_at","logLevel":"","version":"","extaConfigs":null},"TableName":"plan","PrimaryKey":"id","DeleteColumn":"deleted_at","Columns":[{"Prefix":"","CamelName":"Id","ColumnName":"id","Name":"id","Type":"int","Comment":"主键","Tag":"","Nullable":false,"Enums":[],"AutoIncrement":true,"DefaultValue":"","OnCreate":false,"OnUpdate":false,"OnDelete":false},{"Prefix":"","CamelName":"AdvertiserId","ColumnName":"advertiser_id","Name":"advertiser_id","Type":"string","Comment":"广告主","Tag":"","Nullable":false,"Enums":[],"AutoIncrement":false,"DefaultValue":"","OnCreate":false,"OnUpdate":false,"OnDelete":false},{"Prefix":"","CamelName":"Name","ColumnName":"name","Name":"name","Type":"string","Comment":"名称","Tag":"","Nullable":false,"Enums":[],"AutoIncrement":false,"DefaultValue":"","OnCreate":false,"OnUpdate":false,"OnDelete":false},{"Prefix":"","CamelName":"Position","ColumnName":"position","Name":"position","Type":"string","Comment":"位置编码","Tag":"","Nullable":false,"Enums":[],"AutoIncrement":false,"DefaultValue":"","OnCreate":false,"OnUpdate":false,"OnDelete":false},{"Prefix":"","CamelName":"BeginAt","ColumnName":"begin_at","Name":"begin_at","Type":"string","Comment":"投放开始时间","Tag":"","Nullable":true,"Enums":[],"AutoIncrement":false,"DefaultValue":"NULL","OnCreate":false,"OnUpdate":false,"OnDelete":false},{"Prefix":"","CamelName":"EndAt","ColumnName":"end_at","Name":"end_at","Type":"string","Comment":"投放结束时间","Tag":"","Nullable":true,"Enums":[],"AutoIncrement":false,"DefaultValue":"NULL","OnCreate":false,"OnUpdate":false,"OnDelete":false},{"Prefix":"","CamelName":"Did","ColumnName":"did","Name":"did","Type":"int","Comment":"出价","Tag":"","Nullable":true,"Enums":[],"AutoIncrement":false,"DefaultValue":"0","OnCreate":false,"OnUpdate":false,"OnDelete":false},{"Prefix":"","CamelName":"LandingPage","ColumnName":"landing_page","Name":"landing_page","Type":"string","Comment":"落地页","Tag":"","Nullable":false,"Enums":[],"AutoIncrement":false,"DefaultValue":"","OnCreate":false,"OnUpdate":false,"OnDelete":false},{"Prefix":"","CamelName":"CreatedAt","ColumnName":"created_at","Name":"created_at","Type":"string","Comment":"创建时间","Tag":"","Nullable":true,"Enums":[],"AutoIncrement":false,"DefaultValue":"current_timestamp()","OnCreate":true,"OnUpdate":false,"OnDelete":false},{"Prefix":"","CamelName":"UpdatedAt","ColumnName":"updated_at","Name":"updated_at","Type":"string","Comment":"修改时间","Tag":"","Nullable":true,"Enums":[],"AutoIncrement":false,"DefaultValue":"current_timestamp()","OnCreate":false,"OnUpdate":true,"OnDelete":false},{"Prefix":"","CamelName":"DeletedAt","ColumnName":"deleted_at","Name":"deleted_at","Type":"string","Comment":"删除时间","Tag":"","Nullable":true,"Enums":[],"AutoIncrement":false,"DefaultValue":"NULL","OnCreate":false,"OnUpdate":false,"OnDelete":true}],"EnumsConst":[],"Comment":"广告计划","TableDef":null},{"DatabaseConfig":{"databaseName":"ad","tablePrefix":"","columnPrefix":"","deletedAtColumn":"deleted_at","logLevel":"","version":"","extaConfigs":null},"TableName":"window","PrimaryKey":"id","DeleteColumn":"deleted_at","Columns":[{"Prefix":"","CamelName":"Id","ColumnName":"id","Name":"id","Type":"int","Comment":"主键","Tag":"","Nullable":false,"Enums":[],"AutoIncrement":true,"DefaultValue":"","OnCreate":false,"OnUpdate":false,"OnDelete":false},{"Prefix":"","CamelName":"MediaId","ColumnName":"media_id","Name":"media_id","Type":"string","Comment":"媒体Id","Tag":"","Nullable":false,"Enums":[],"AutoIncrement":false,"DefaultValue":"","OnCreate":false,"OnUpdate":false,"OnDelete":false},{"Prefix":"","CamelName":"Position","ColumnName":"position","Name":"position","Type":"string","Comment":"位置编码","Tag":"","Nullable":false,"Enums":[],"AutoIncrement":false,"DefaultValue":"","OnCreate":false,"OnUpdate":false,"OnDelete":false},{"Prefix":"","CamelName":"Name","ColumnName":"name","Name":"name","Type":"string","Comment":"位置名称","Tag":"","Nullable":false,"Enums":[],"AutoIncrement":false,"DefaultValue":"","OnCreate":false,"OnUpdate":false,"OnDelete":false},{"Prefix":"","CamelName":"Remark","ColumnName":"remark","Name":"remark","Type":"string","Comment":"广告位描述(建议记录位置、app名称等)","Tag":"","Nullable":false,"Enums":[],"AutoIncrement":false,"DefaultValue":"","OnCreate":false,"OnUpdate":false,"OnDelete":false},{"Prefix":"","CamelName":"Scheme","ColumnName":"scheme","Name":"scheme","Type":"string","Comment":"广告素材的格式规范","Tag":"","Nullable":true,"Enums":[],"AutoIncrement":false,"DefaultValue":"","OnCreate":false,"OnUpdate":false,"OnDelete":false},{"Prefix":"","CamelName":"CreatedAt","ColumnName":"created_at","Name":"created_at","Type":"string","Comment":"创建时间","Tag":"","Nullable":true,"Enums":[],"AutoIncrement":false,"DefaultValue":"current_timestamp()","OnCreate":true,"OnUpdate":false,"OnDelete":false},{"Prefix":"","CamelName":"UpdatedAt","ColumnName":"updated_at","Name":"updated_at","Type":"string","Comment":"修改时间","Tag":"","Nullable":true,"Enums":[],"AutoIncrement":false,"DefaultValue":"current_timestamp()","OnCreate":false,"OnUpdate":true,"OnDelete":false},{"Prefix":"","CamelName":"DeletedAt","ColumnName":"deleted_at","Name":"deleted_at","Type":"string","Comment":"删除时间","Tag":"","Nullable":true,"Enums":[],"AutoIncrement":false,"DefaultValue":"NULL","OnCreate":false,"OnUpdate":false,"OnDelete":true}],"EnumsConst":[],"Comment":"广告位表","TableDef":null}]`
	//jsonStr := `[{"Enums":[]}]`
	newJson, err := lineschema.TransferJson(jsonStr, func(transfer lineschema.Transfers) (newTransfer lineschema.Transfers) {
		return transfer
	})
	require.NoError(t, err)
	require.JSONEq(t, jsonStr, newJson)
}

func TestTransferJsonChange(t *testing.T) {
	jsonStr := `[{"DatabaseConfig":{"databaseName":"ad","tablePrefix":"","columnPrefix":"","deletedAtColumn":"deleted_at","logLevel":"","version":"","extaConfigs":null},"TableName":"creative","PrimaryKey":"id","DeleteColumn":"deleted_at","Columns":[{"Prefix":"","CamelName":"Id","ColumnName":"id","Name":"id","Type":"int","Comment":"主键","Tag":"","Nullable":false,"Enums":[],"AutoIncrement":true,"DefaultValue":"","OnCreate":false,"OnUpdate":false,"OnDelete":false},{"Prefix":"","CamelName":"PlanId","ColumnName":"plan_id","Name":"plan_id","Type":"string","Comment":"广告计划Id","Tag":"","Nullable":false,"Enums":[],"AutoIncrement":false,"DefaultValue":"","OnCreate":false,"OnUpdate":false,"OnDelete":false},{"Prefix":"","CamelName":"Name","ColumnName":"name","Name":"name","Type":"string","Comment":"名称","Tag":"","Nullable":false,"Enums":[],"AutoIncrement":false,"DefaultValue":"","OnCreate":false,"OnUpdate":false,"OnDelete":false},{"Prefix":"","CamelName":"Content","ColumnName":"content","Name":"content","Type":"string","Comment":"广告内容","Tag":"","Nullable":true,"Enums":[],"AutoIncrement":false,"DefaultValue":"","OnCreate":false,"OnUpdate":false,"OnDelete":false},{"Prefix":"","CamelName":"CreatedAt","ColumnName":"created_at","Name":"created_at","Type":"string","Comment":"创建时间","Tag":"","Nullable":true,"Enums":[],"AutoIncrement":false,"DefaultValue":"current_timestamp()","OnCreate":true,"OnUpdate":false,"OnDelete":false},{"Prefix":"","CamelName":"UpdatedAt","ColumnName":"updated_at","Name":"updated_at","Type":"string","Comment":"修改时间","Tag":"","Nullable":true,"Enums":[],"AutoIncrement":false,"DefaultValue":"current_timestamp()","OnCreate":false,"OnUpdate":true,"OnDelete":false},{"Prefix":"","CamelName":"DeletedAt","ColumnName":"deleted_at","Name":"deleted_at","Type":"string","Comment":"删除时间","Tag":"","Nullable":true,"Enums":[],"AutoIncrement":false,"DefaultValue":"NULL","OnCreate":false,"OnUpdate":false,"OnDelete":true}],"EnumsConst":[],"Comment":"广告物料","TableDef":null},{"DatabaseConfig":{"databaseName":"ad","tablePrefix":"","columnPrefix":"","deletedAtColumn":"deleted_at","logLevel":"","version":"","extaConfigs":null},"TableName":"plan","PrimaryKey":"id","DeleteColumn":"deleted_at","Columns":[{"Prefix":"","CamelName":"Id","ColumnName":"id","Name":"id","Type":"int","Comment":"主键","Tag":"","Nullable":false,"Enums":[],"AutoIncrement":true,"DefaultValue":"","OnCreate":false,"OnUpdate":false,"OnDelete":false},{"Prefix":"","CamelName":"AdvertiserId","ColumnName":"advertiser_id","Name":"advertiser_id","Type":"string","Comment":"广告主","Tag":"","Nullable":false,"Enums":[],"AutoIncrement":false,"DefaultValue":"","OnCreate":false,"OnUpdate":false,"OnDelete":false},{"Prefix":"","CamelName":"Name","ColumnName":"name","Name":"name","Type":"string","Comment":"名称","Tag":"","Nullable":false,"Enums":[],"AutoIncrement":false,"DefaultValue":"","OnCreate":false,"OnUpdate":false,"OnDelete":false},{"Prefix":"","CamelName":"Position","ColumnName":"position","Name":"position","Type":"string","Comment":"位置编码","Tag":"","Nullable":false,"Enums":[],"AutoIncrement":false,"DefaultValue":"","OnCreate":false,"OnUpdate":false,"OnDelete":false},{"Prefix":"","CamelName":"BeginAt","ColumnName":"begin_at","Name":"begin_at","Type":"string","Comment":"投放开始时间","Tag":"","Nullable":true,"Enums":[],"AutoIncrement":false,"DefaultValue":"NULL","OnCreate":false,"OnUpdate":false,"OnDelete":false},{"Prefix":"","CamelName":"EndAt","ColumnName":"end_at","Name":"end_at","Type":"string","Comment":"投放结束时间","Tag":"","Nullable":true,"Enums":[],"AutoIncrement":false,"DefaultValue":"NULL","OnCreate":false,"OnUpdate":false,"OnDelete":false},{"Prefix":"","CamelName":"Did","ColumnName":"did","Name":"did","Type":"int","Comment":"出价","Tag":"","Nullable":true,"Enums":[],"AutoIncrement":false,"DefaultValue":"0","OnCreate":false,"OnUpdate":false,"OnDelete":false},{"Prefix":"","CamelName":"LandingPage","ColumnName":"landing_page","Name":"landing_page","Type":"string","Comment":"落地页","Tag":"","Nullable":false,"Enums":[],"AutoIncrement":false,"DefaultValue":"","OnCreate":false,"OnUpdate":false,"OnDelete":false},{"Prefix":"","CamelName":"CreatedAt","ColumnName":"created_at","Name":"created_at","Type":"string","Comment":"创建时间","Tag":"","Nullable":true,"Enums":[],"AutoIncrement":false,"DefaultValue":"current_timestamp()","OnCreate":true,"OnUpdate":false,"OnDelete":false},{"Prefix":"","CamelName":"UpdatedAt","ColumnName":"updated_at","Name":"updated_at","Type":"string","Comment":"修改时间","Tag":"","Nullable":true,"Enums":[],"AutoIncrement":false,"DefaultValue":"current_timestamp()","OnCreate":false,"OnUpdate":true,"OnDelete":false},{"Prefix":"","CamelName":"DeletedAt","ColumnName":"deleted_at","Name":"deleted_at","Type":"string","Comment":"删除时间","Tag":"","Nullable":true,"Enums":[],"AutoIncrement":false,"DefaultValue":"NULL","OnCreate":false,"OnUpdate":false,"OnDelete":true}],"EnumsConst":[],"Comment":"广告计划","TableDef":null},{"DatabaseConfig":{"databaseName":"ad","tablePrefix":"","columnPrefix":"","deletedAtColumn":"deleted_at","logLevel":"","version":"","extaConfigs":null},"TableName":"window","PrimaryKey":"id","DeleteColumn":"deleted_at","Columns":[{"Prefix":"","CamelName":"Id","ColumnName":"id","Name":"id","Type":"int","Comment":"主键","Tag":"","Nullable":false,"Enums":[],"AutoIncrement":true,"DefaultValue":"","OnCreate":false,"OnUpdate":false,"OnDelete":false},{"Prefix":"","CamelName":"MediaId","ColumnName":"media_id","Name":"media_id","Type":"string","Comment":"媒体Id","Tag":"","Nullable":false,"Enums":[],"AutoIncrement":false,"DefaultValue":"","OnCreate":false,"OnUpdate":false,"OnDelete":false},{"Prefix":"","CamelName":"Position","ColumnName":"position","Name":"position","Type":"string","Comment":"位置编码","Tag":"","Nullable":false,"Enums":[],"AutoIncrement":false,"DefaultValue":"","OnCreate":false,"OnUpdate":false,"OnDelete":false},{"Prefix":"","CamelName":"Name","ColumnName":"name","Name":"name","Type":"string","Comment":"位置名称","Tag":"","Nullable":false,"Enums":[],"AutoIncrement":false,"DefaultValue":"","OnCreate":false,"OnUpdate":false,"OnDelete":false},{"Prefix":"","CamelName":"Remark","ColumnName":"remark","Name":"remark","Type":"string","Comment":"广告位描述(建议记录位置、app名称等)","Tag":"","Nullable":false,"Enums":[],"AutoIncrement":false,"DefaultValue":"","OnCreate":false,"OnUpdate":false,"OnDelete":false},{"Prefix":"","CamelName":"Scheme","ColumnName":"scheme","Name":"scheme","Type":"string","Comment":"广告素材的格式规范","Tag":"","Nullable":true,"Enums":[],"AutoIncrement":false,"DefaultValue":"","OnCreate":false,"OnUpdate":false,"OnDelete":false},{"Prefix":"","CamelName":"CreatedAt","ColumnName":"created_at","Name":"created_at","Type":"string","Comment":"创建时间","Tag":"","Nullable":true,"Enums":[],"AutoIncrement":false,"DefaultValue":"current_timestamp()","OnCreate":true,"OnUpdate":false,"OnDelete":false},{"Prefix":"","CamelName":"UpdatedAt","ColumnName":"updated_at","Name":"updated_at","Type":"string","Comment":"修改时间","Tag":"","Nullable":true,"Enums":[],"AutoIncrement":false,"DefaultValue":"current_timestamp()","OnCreate":false,"OnUpdate":true,"OnDelete":false},{"Prefix":"","CamelName":"DeletedAt","ColumnName":"deleted_at","Name":"deleted_at","Type":"string","Comment":"删除时间","Tag":"","Nullable":true,"Enums":[],"AutoIncrement":false,"DefaultValue":"NULL","OnCreate":false,"OnUpdate":false,"OnDelete":true}],"EnumsConst":[],"Comment":"广告位表","TableDef":null}]`
	newJson, err := lineschema.TransferJson(jsonStr, func(transfer lineschema.Transfers) (newTransfer lineschema.Transfers) {
		transfer = transfer.ModifyDstPath(lineschema.PathModifyFnSmallCameCase)
		//transfer = transfer.ModifySrcPath(lineschema.PathModifyFnString)
		return transfer
	})
	require.NoError(t, err)
	fmt.Println(newJson)
}
