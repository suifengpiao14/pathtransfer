package pathtransfer

import (
	"encoding/json"
	"strings"
)

/**
line stransfer example:
api.getUser.input.id@int:db.user.Fuser_id@int
api.getUser.input.name:db.user.Fname
**/

func Parse(s string) (ts Transfers) {
	s = strings.TrimSpace(s)
	rows := strings.Split(s, "\n")
	ts = make(Transfers, 0)
	for _, row := range rows {
		row = strings.TrimSpace(row)
		if row == "" {
			continue
		}
		var src, dst, srcType, dstType string
		src, dst = row, row
		colonIndex := strings.Index(src, ":")
		if colonIndex > -1 {
			dst = src[colonIndex+1:]
			src = src[:colonIndex]
		}
		srcAtIndex := strings.LastIndex(src, "@")
		if srcAtIndex > -1 {
			srcType = src[srcAtIndex+1:]
			src = src[:srcAtIndex]
		}
		dstAtIndex := strings.LastIndex(dst, "@")
		if dstAtIndex > -1 {
			dstType = dst[dstAtIndex+1:]
			dst = dst[:dstAtIndex]
		}
		t := Transfer{
			Src: TransferUnit{
				Path: src,
				Type: srcType,
			},
			Dst: TransferUnit{
				Path: dst,
				Type: dstType,
			},
		}
		ts = append(ts, t)
	}
	return ts
}

type TransferJson string

//Vocabulary 转换为Vocabulary对象
func (tJson TransferJson) Vocabularies() (vocabularies Transfers, err error) {
	if tJson == "" {
		return nil, nil
	}
	vocabularies = make(Transfers, 0)
	err = json.Unmarshal([]byte(tJson), &vocabularies)
	if err != nil {
		return nil, err
	}
	return vocabularies, nil
}
