package x

// chat 'Chat'.
type Chat struct {
	ChatKey              string
	RoomKey              string
	RoomTypeEnumId       int
	UserId               int
	PeerUserId           int
	GroupId              int
	CreatedSe            int
	StartMessageIdFrom   int
	LastSeenMessageId    int
	UpdatedMs            int
	LastMessageId        int //just direct
	LastDeletedMessageId int
	LastSeqSeen          int
	LastSeqDelete        int
	CurrentSeq           int //just for peer to peer

	_exists, _deleted bool
}

/*
:= &Chat {
	ChatKey: "",
	RoomKey: "",
	RoomTypeEnumId: 0,
	UserId: 0,
	PeerUserId: 0,
	GroupId: 0,
	CreatedSe: 0,
	StartMessageIdFrom: 0,
	LastSeenMessageId: 0,
	UpdatedMs: 0,
	LastMessageId: 0,
	LastDeletedMessageId: 0,
	LastSeqSeen: 0,
	LastSeqDelete: 0,
	CurrentSeq: 0,
*/
// comments 'Comments'.
type Comments struct {
	Id          int
	UserId      int
	PostId      int
	Text        string
	CreatedTime int

	_exists, _deleted bool
}

/*
:= &Comments {
	Id: 0,
	UserId: 0,
	PostId: 0,
	Text: "",
	CreatedTime: 0,
*/
// trigger_log 'TriggerLog'.
type TriggerLog struct {
	Id         int
	TableName  string
	ChangeType string
	TargetId   int
	TargetStr  string

	_exists, _deleted bool
}

/*
:= &TriggerLog {
	Id: 0,
	TableName: "",
	ChangeType: "",
	TargetId: 0,
	TargetStr: "",
*/
