package pathtransfer_test

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/suifengpiao14/pathtransfer"
	"github.com/tidwall/gjson"
)

type user struct {
	Name   string `json:"name"`
	UserId int    `json:"userId"`
}

func TestToGoTypeTransfer(t *testing.T) {

	t.Run("struct", func(t *testing.T) {

		lineSchema := pathtransfer.ToGoTypeTransfer(new(user)).GjsonPath()
		expected := `{name:@this.name.@tostring,userId:@this.userId.@tonum}`
		assert.Equal(t, expected, lineSchema)
	})

	t.Run("slice[struct]", func(t *testing.T) {
		users := make([]user, 0)
		lineSchema := pathtransfer.ToGoTypeTransfer(users).GjsonPath()
		expected := `{name:@this.#.name.@tostring,userId:@this.#.userId.@tonum}|@groupPlus:0`
		assert.Equal(t, expected, lineSchema)
	})
	t.Run("array[struct]", func(t *testing.T) {
		users := [2]user{}
		lineSchema := pathtransfer.ToGoTypeTransfer(users).GjsonPath()
		expected := `{name:@this.#.name.@tostring,userId:@this.#.userId.@tonum}|@groupPlus:0`
		assert.Equal(t, expected, lineSchema)
	})

	t.Run("array[int]", func(t *testing.T) {
		ids := [2]string{}
		lineSchema := pathtransfer.ToGoTypeTransfer(ids).GjsonPath()
		expected := `@this.#.@tostring`
		assert.Equal(t, expected, lineSchema)
	})

	t.Run("int", func(t *testing.T) {
		id := 2
		lineSchema := pathtransfer.ToGoTypeTransfer(id).GjsonPath()
		expected := `@this.@tonum`
		assert.Equal(t, expected, lineSchema)
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
	newJson, err := pathtransfer.RebuildJson(jsonStr, func(transfer pathtransfer.Transfers) (newTransfer pathtransfer.Transfers) {
		return transfer
	})
	require.NoError(t, err)
	require.JSONEq(t, jsonStr, newJson)
}

func TestTransferJsonChange(t *testing.T) {
	jsonStr := `[{"DatabaseConfig":{"databaseName":"ad","tablePrefix":"","columnPrefix":"","deletedAtColumn":"deleted_at","logLevel":"","version":"","extaConfigs":null},"TableName":"creative","PrimaryKey":"id","DeleteColumn":"deleted_at","Columns":[{"Prefix":"","CamelName":"Id","ColumnName":"id","Name":"id","Type":"int","Comment":"主键","Tag":"","Nullable":false,"Enums":[],"AutoIncrement":true,"DefaultValue":"","OnCreate":false,"OnUpdate":false,"OnDelete":false},{"Prefix":"","CamelName":"PlanId","ColumnName":"plan_id","Name":"plan_id","Type":"string","Comment":"广告计划Id","Tag":"","Nullable":false,"Enums":[],"AutoIncrement":false,"DefaultValue":"","OnCreate":false,"OnUpdate":false,"OnDelete":false},{"Prefix":"","CamelName":"Name","ColumnName":"name","Name":"name","Type":"string","Comment":"名称","Tag":"","Nullable":false,"Enums":[],"AutoIncrement":false,"DefaultValue":"","OnCreate":false,"OnUpdate":false,"OnDelete":false},{"Prefix":"","CamelName":"Content","ColumnName":"content","Name":"content","Type":"string","Comment":"广告内容","Tag":"","Nullable":true,"Enums":[],"AutoIncrement":false,"DefaultValue":"","OnCreate":false,"OnUpdate":false,"OnDelete":false},{"Prefix":"","CamelName":"CreatedAt","ColumnName":"created_at","Name":"created_at","Type":"string","Comment":"创建时间","Tag":"","Nullable":true,"Enums":[],"AutoIncrement":false,"DefaultValue":"current_timestamp()","OnCreate":true,"OnUpdate":false,"OnDelete":false},{"Prefix":"","CamelName":"UpdatedAt","ColumnName":"updated_at","Name":"updated_at","Type":"string","Comment":"修改时间","Tag":"","Nullable":true,"Enums":[],"AutoIncrement":false,"DefaultValue":"current_timestamp()","OnCreate":false,"OnUpdate":true,"OnDelete":false},{"Prefix":"","CamelName":"DeletedAt","ColumnName":"deleted_at","Name":"deleted_at","Type":"string","Comment":"删除时间","Tag":"","Nullable":true,"Enums":[],"AutoIncrement":false,"DefaultValue":"NULL","OnCreate":false,"OnUpdate":false,"OnDelete":true}],"EnumsConst":[],"Comment":"广告物料","TableDef":null},{"DatabaseConfig":{"databaseName":"ad","tablePrefix":"","columnPrefix":"","deletedAtColumn":"deleted_at","logLevel":"","version":"","extaConfigs":null},"TableName":"plan","PrimaryKey":"id","DeleteColumn":"deleted_at","Columns":[{"Prefix":"","CamelName":"Id","ColumnName":"id","Name":"id","Type":"int","Comment":"主键","Tag":"","Nullable":false,"Enums":[],"AutoIncrement":true,"DefaultValue":"","OnCreate":false,"OnUpdate":false,"OnDelete":false},{"Prefix":"","CamelName":"AdvertiserId","ColumnName":"advertiser_id","Name":"advertiser_id","Type":"string","Comment":"广告主","Tag":"","Nullable":false,"Enums":[],"AutoIncrement":false,"DefaultValue":"","OnCreate":false,"OnUpdate":false,"OnDelete":false},{"Prefix":"","CamelName":"Name","ColumnName":"name","Name":"name","Type":"string","Comment":"名称","Tag":"","Nullable":false,"Enums":[],"AutoIncrement":false,"DefaultValue":"","OnCreate":false,"OnUpdate":false,"OnDelete":false},{"Prefix":"","CamelName":"Position","ColumnName":"position","Name":"position","Type":"string","Comment":"位置编码","Tag":"","Nullable":false,"Enums":[],"AutoIncrement":false,"DefaultValue":"","OnCreate":false,"OnUpdate":false,"OnDelete":false},{"Prefix":"","CamelName":"BeginAt","ColumnName":"begin_at","Name":"begin_at","Type":"string","Comment":"投放开始时间","Tag":"","Nullable":true,"Enums":[],"AutoIncrement":false,"DefaultValue":"NULL","OnCreate":false,"OnUpdate":false,"OnDelete":false},{"Prefix":"","CamelName":"EndAt","ColumnName":"end_at","Name":"end_at","Type":"string","Comment":"投放结束时间","Tag":"","Nullable":true,"Enums":[],"AutoIncrement":false,"DefaultValue":"NULL","OnCreate":false,"OnUpdate":false,"OnDelete":false},{"Prefix":"","CamelName":"Did","ColumnName":"did","Name":"did","Type":"int","Comment":"出价","Tag":"","Nullable":true,"Enums":[],"AutoIncrement":false,"DefaultValue":"0","OnCreate":false,"OnUpdate":false,"OnDelete":false},{"Prefix":"","CamelName":"LandingPage","ColumnName":"landing_page","Name":"landing_page","Type":"string","Comment":"落地页","Tag":"","Nullable":false,"Enums":[],"AutoIncrement":false,"DefaultValue":"","OnCreate":false,"OnUpdate":false,"OnDelete":false},{"Prefix":"","CamelName":"CreatedAt","ColumnName":"created_at","Name":"created_at","Type":"string","Comment":"创建时间","Tag":"","Nullable":true,"Enums":[],"AutoIncrement":false,"DefaultValue":"current_timestamp()","OnCreate":true,"OnUpdate":false,"OnDelete":false},{"Prefix":"","CamelName":"UpdatedAt","ColumnName":"updated_at","Name":"updated_at","Type":"string","Comment":"修改时间","Tag":"","Nullable":true,"Enums":[],"AutoIncrement":false,"DefaultValue":"current_timestamp()","OnCreate":false,"OnUpdate":true,"OnDelete":false},{"Prefix":"","CamelName":"DeletedAt","ColumnName":"deleted_at","Name":"deleted_at","Type":"string","Comment":"删除时间","Tag":"","Nullable":true,"Enums":[],"AutoIncrement":false,"DefaultValue":"NULL","OnCreate":false,"OnUpdate":false,"OnDelete":true}],"EnumsConst":[],"Comment":"广告计划","TableDef":null},{"DatabaseConfig":{"databaseName":"ad","tablePrefix":"","columnPrefix":"","deletedAtColumn":"deleted_at","logLevel":"","version":"","extaConfigs":null},"TableName":"window","PrimaryKey":"id","DeleteColumn":"deleted_at","Columns":[{"Prefix":"","CamelName":"Id","ColumnName":"id","Name":"id","Type":"int","Comment":"主键","Tag":"","Nullable":false,"Enums":[],"AutoIncrement":true,"DefaultValue":"","OnCreate":false,"OnUpdate":false,"OnDelete":false},{"Prefix":"","CamelName":"MediaId","ColumnName":"media_id","Name":"media_id","Type":"string","Comment":"媒体Id","Tag":"","Nullable":false,"Enums":[],"AutoIncrement":false,"DefaultValue":"","OnCreate":false,"OnUpdate":false,"OnDelete":false},{"Prefix":"","CamelName":"Position","ColumnName":"position","Name":"position","Type":"string","Comment":"位置编码","Tag":"","Nullable":false,"Enums":[],"AutoIncrement":false,"DefaultValue":"","OnCreate":false,"OnUpdate":false,"OnDelete":false},{"Prefix":"","CamelName":"Name","ColumnName":"name","Name":"name","Type":"string","Comment":"位置名称","Tag":"","Nullable":false,"Enums":[],"AutoIncrement":false,"DefaultValue":"","OnCreate":false,"OnUpdate":false,"OnDelete":false},{"Prefix":"","CamelName":"Remark","ColumnName":"remark","Name":"remark","Type":"string","Comment":"广告位描述(建议记录位置、app名称等)","Tag":"","Nullable":false,"Enums":[],"AutoIncrement":false,"DefaultValue":"","OnCreate":false,"OnUpdate":false,"OnDelete":false},{"Prefix":"","CamelName":"Scheme","ColumnName":"scheme","Name":"scheme","Type":"string","Comment":"广告素材的格式规范","Tag":"","Nullable":true,"Enums":[],"AutoIncrement":false,"DefaultValue":"","OnCreate":false,"OnUpdate":false,"OnDelete":false},{"Prefix":"","CamelName":"CreatedAt","ColumnName":"created_at","Name":"created_at","Type":"string","Comment":"创建时间","Tag":"","Nullable":true,"Enums":[],"AutoIncrement":false,"DefaultValue":"current_timestamp()","OnCreate":true,"OnUpdate":false,"OnDelete":false},{"Prefix":"","CamelName":"UpdatedAt","ColumnName":"updated_at","Name":"updated_at","Type":"string","Comment":"修改时间","Tag":"","Nullable":true,"Enums":[],"AutoIncrement":false,"DefaultValue":"current_timestamp()","OnCreate":false,"OnUpdate":true,"OnDelete":false},{"Prefix":"","CamelName":"DeletedAt","ColumnName":"deleted_at","Name":"deleted_at","Type":"string","Comment":"删除时间","Tag":"","Nullable":true,"Enums":[],"AutoIncrement":false,"DefaultValue":"NULL","OnCreate":false,"OnUpdate":false,"OnDelete":true}],"EnumsConst":[],"Comment":"广告位表","TableDef":null}]`
	newJson, err := pathtransfer.RebuildJson(jsonStr, func(transfer pathtransfer.Transfers) (newTransfer pathtransfer.Transfers) {
		transfer = transfer.ModifyDstPath(pathtransfer.PathModifyFnSmallCameCase)
		return transfer
	})
	require.NoError(t, err)
	fmt.Println(newJson)
}

func TestPathBaseName(t *testing.T) {
	path := "abc.input.pagination.index"
	baseName := pathtransfer.Path(path).NameWithoutSpace()
	require.Equal(t, "pagination.index", baseName)

}
