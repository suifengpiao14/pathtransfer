package pathtransfer

import (
	"encoding/json"
	"strings"
)

type TransferLine string

func (l TransferLine) Transfer() (ts Transfers) {
	return Parse(string(l))
}

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

//Unmarshal 转换为Transfers对象
func Unmarshal(tJson string) (vocabularies Transfers, err error) {
	vocabularies = make(Transfers, 0)
	err = json.Unmarshal([]byte(tJson), &vocabularies)
	if err != nil {
		return nil, err
	}
	return vocabularies, nil
}
func Marshal(tJson string) (vocabularies Transfers, err error) {
	vocabularies = make(Transfers, 0)
	err = json.Unmarshal([]byte(tJson), &vocabularies)
	if err != nil {
		return nil, err
	}
	return vocabularies, nil
}
