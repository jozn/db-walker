package x

import (
	"strconv"
	"time"
)

//Chat Events

func OnChat_AfterInsert(row *Chat) {
	RowCache.Set("Chat:"+row.ChatKey, row, time.Hour*0)
}

func OnChat_AfterUpdate(row *Chat) {
	RowCache.Set("Chat:"+row.ChatKey, row, time.Hour*0)
}

func OnChat_AfterDelete(row *Chat) {
	RowCache.Delete("Chat:" + row.ChatKey)
}

func OnChat_LoadOne(row *Chat) {
	RowCache.Set("Chat:"+row.ChatKey, row, time.Hour*0)
}

func OnChat_LoadMany(rows []*Chat) {
	for _, row := range rows {
		RowCache.Set("Chat:"+row.ChatKey, row, time.Hour*0)
	}
}

//Comments Events

func OnComments_AfterInsert(row *Comments) {
	RowCache.Set("Comments:"+strconv.Itoa(row.Id), row, time.Hour*0)
}

func OnComments_AfterUpdate(row *Comments) {
	RowCache.Set("Comments:"+strconv.Itoa(row.Id), row, time.Hour*0)
}

func OnComments_AfterDelete(row *Comments) {
	RowCache.Delete("Comments:" + strconv.Itoa(row.Id))
}

func OnComments_LoadOne(row *Comments) {
	RowCache.Set("Comments:"+strconv.Itoa(row.Id), row, time.Hour*0)
}

func OnComments_LoadMany(rows []*Comments) {
	for _, row := range rows {
		RowCache.Set("Comments:"+strconv.Itoa(row.Id), row, time.Hour*0)
	}
}
