package x

import (
	"ms/sun/base"
	"strconv"
)

func (c _StoreImpl) GetChatByChatKey(ChatKey string) (*Chat, bool) {
	o, ok := RowCache.Get("Chat:" + ChatKey)
	if ok {
		if obj, ok := o.(*Chat); ok {
			return obj, true
		}
	}
	obj2, err := ChatByChatKey(base.DB, ChatKey)
	if err == nil {
		return obj2, true
	}
	XOLogErr(err)
	return nil, false
}

func (c _StoreImpl) PreLoadChatByChatKeys(ids []string) {
	not_cached := make([]string, 0, len(ids))

	for _, id := range ids {
		_, ok := RowCache.Get("Chat:" + id)
		if !ok {
			not_cached = append(not_cached, id)
		}
	}

	if len(not_cached) > 0 {
		NewChat_Selector().ChatKey_In(not_cached).GetRows(base.DB)
	}
}

// yes 222 string

func (c _StoreImpl) GetCommentsById(Id int) (*Comments, bool) {
	o, ok := RowCache.Get("Comments:" + strconv.Itoa(Id))
	if ok {
		if obj, ok := o.(*Comments); ok {
			return obj, true
		}
	}
	obj2, err := CommentsById(base.DB, Id)
	if err == nil {
		return obj2, true
	}
	XOLogErr(err)
	return nil, false
}

func (c _StoreImpl) PreLoadCommentsByIds(ids []int) {
	not_cached := make([]int, 0, len(ids))

	for _, id := range ids {
		_, ok := RowCache.Get("Comments:" + strconv.Itoa(id))
		if !ok {
			not_cached = append(not_cached, id)
		}
	}

	if len(not_cached) > 0 {
		NewComments_Selector().Id_In(not_cached).GetRows(base.DB)
	}
}

// yes 222 int
