package biz_internal

import "errors"

func (b *BizStruct) communityIdCheck(communityId int) error {
	if communityId <= 0 {
		return errors.New("incorrect community id")
	}
	return nil
}
