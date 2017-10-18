package x

import (
	"fmt"
	"ms/sun/base"
)

// Chat - PRIMARY

// Comments - PRIMARY

// Comments - AllPostUser

//field//field//field

///// Generated from index 'PostId'.
func (c _StoreImpl) Comments_ByPostId(PostId int) (*Comments, bool) {
	o, ok := RowCacheIndex.Get("Comments_PostId:" + fmt.Sprintf("%v", PostId))
	if ok {
		if obj, ok := o.(*Comments); ok {
			return obj, true
		}
	}

	row, err := NewComments_Selector().PostId_Eq(PostId).GetRow(base.DB)
	if err == nil {
		RowCacheIndex.Set("Comments_PostId:"+fmt.Sprintf("%v", row.PostId), row, 0)
		return row, true
	}

	XOLogErr(err)
	return nil, false
}

func (c _StoreImpl) PreLoadComments_ByPostIds(PostIds []int) {
	not_cached := make([]int, 0, len(PostIds))

	for _, id := range PostIds {
		_, ok := RowCacheIndex.Get("Comments_PostId:" + fmt.Sprintf("%v", id))
		if !ok {
			not_cached = append(not_cached, id)
		}
	}

	if len(not_cached) > 0 {
		rows, err := NewComments_Selector().PostId_In(not_cached).GetRows(base.DB)
		if err == nil {
			for _, row := range rows {
				RowCacheIndex.Set("Comments_PostId:"+fmt.Sprintf("%v", row.PostId), row, 0)
			}
		}
	}
}
