package polls_biz

import (
	biz_internal "oprosdom.ru/core/internal/polls/biz/internal"
)

type BizInterface interface {

}

func NewBizFactory() BizInterface {
	return biz_internal.NewCallInternalBiz()
}
