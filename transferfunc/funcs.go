package transferfunc

import (
	"fmt"
)

//Limit 设置SQL分页
func Limit(index int, size int) (offset int, limit int) {
	limit = size
	offset = index * size
	return offset, limit
}

//LikePrefix 设置like 前置 %
func LikePrefix(value string) (newValue string) {
	if value == "" {
		return ""
	}
	newValue = fmt.Sprintf("%%%s", value)
	return
}

//LikeSuffix 设置like 后置 %
func LikeSuffix(value string) (newValue string) {
	if value == "" {
		return ""
	}
	newValue = fmt.Sprintf("%s%%", value)
	return
}

//Like 设置like 前置后置 %
func Like(value string) (newValue string) {
	if value == "" {
		return ""
	}
	newValue = fmt.Sprintf("%%%s%%", value)
	return
}
