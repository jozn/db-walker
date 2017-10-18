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

// TriggerLog - PRIMARY

//field//field//field

///// Generated from index 'Id'.
func (c _StoreImpl) TriggerLog_ById(Id int) (*TriggerLog, bool) {
	o, ok := RowCacheIndex.Get("TriggerLog_Id:" + fmt.Sprintf("%v", Id))
	if ok {
		if obj, ok := o.(*TriggerLog); ok {
			return obj, true
		}
	}

	row, err := NewTriggerLog_Selector().Id_Eq(Id).GetRow(base.DB)
	if err == nil {
		RowCacheIndex.Set("TriggerLog_Id:"+fmt.Sprintf("%v", row.Id), row, 0)
		return row, true
	}

	XOLogErr(err)
	return nil, false
}

func (c _StoreImpl) PreLoadTriggerLog_ByIds(Ids []int) {
	not_cached := make([]int, 0, len(Ids))

	for _, id := range Ids {
		_, ok := RowCacheIndex.Get("TriggerLog_Id:" + fmt.Sprintf("%v", id))
		if !ok {
			not_cached = append(not_cached, id)
		}
	}

	if len(not_cached) > 0 {
		rows, err := NewTriggerLog_Selector().Id_In(not_cached).GetRows(base.DB)
		if err == nil {
			for _, row := range rows {
				RowCacheIndex.Set("TriggerLog_Id:"+fmt.Sprintf("%v", row.Id), row, 0)
			}
		}
	}
}
