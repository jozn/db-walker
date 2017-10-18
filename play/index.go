package x

import "github.com/jmoiron/sqlx"

// GetChatByChatKey Generated from index 'PRIMARY' -- retrieves a row from 'os.chat' as a Chat.
func GetChatByChatKey(db sqlx.DB, chatKey string) (*Chat, error) {
	var err error

	const sqlstr = `SELECT * ` +
		`FROM os.chat ` +
		`WHERE ChatKey = ?`

	XOLog(sqlstr, chatKey)
	c := Chat{
		_exists: true,
	}

	err = db.Get(&c, sqlstr, chatKey)
	if err != nil {
		XOLogErr(err)
		return nil, err
	}

	OnChat_LoadOne(&c)

	return &c, nil
}

// GetCommentsById Generated from index 'PRIMARY' -- retrieves a row from 'os.comments' as a Comments.
func GetCommentsById(db sqlx.DB, id int) (*Comments, error) {
	var err error

	const sqlstr = `SELECT * ` +
		`FROM os.comments ` +
		`WHERE Id = ?`

	XOLog(sqlstr, id)
	c := Comments{
		_exists: true,
	}

	err = db.Get(&c, sqlstr, id)
	if err != nil {
		XOLogErr(err)
		return nil, err
	}

	OnComments_LoadOne(&c)

	return &c, nil
}

// GetCommentsByIdAndUserIdAndPostIdAndTextAndCreatedTime Generated from index 'AllPostUser' -- retrieves a row from 'os.comments' as a Comments.
func GetCommentsByIdAndUserIdAndPostIdAndTextAndCreatedTime(db sqlx.DB, userId int, postId int, createdTime int, id int) (*Comments, error) {
	var err error

	const sqlstr = `SELECT * ` +
		`FROM os.comments ` +
		`WHERE UserId = ? AND PostId = ? AND CreatedTime = ? AND Id = ?`

	XOLog(sqlstr, userId, postId, createdTime, id)
	c := Comments{
		_exists: true,
	}

	err = db.Get(&c, sqlstr, userId, postId, createdTime, id)
	if err != nil {
		XOLogErr(err)
		return nil, err
	}

	OnComments_LoadOne(&c)

	return &c, nil
}

// GetCommentsByPostId Generated from index 'PostId' -- retrieves a row from 'os.comments' as a Comments.
func GetCommentsByPostId(db sqlx.DB, postId int) ([]*Comments, error) {
	var err error

	const sqlstr = `SELECT * ` +
		`FROM os.comments ` +
		`WHERE PostId = ?`

	XOLog(sqlstr, postId)
	res := []*Comments{}
	err = db.Select(&res, sqlstr, postId)
	if err != nil {
		XOLogErr(err)
		return res, err
	}
	OnComments_LoadMany(res)

	return res, nil
}
