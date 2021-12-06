package my

func CreatePost(user *User, msg string, pId int) {
	cmt := Comment{
		UserId:  int(user.Model.ID),
		PostId:  pId,
		Message: msg,
	}

	DB.Create(&cmt)
	return
}

func FindPostById(id string, pst *Post) {
	DB.Where("id = ?", id).First(pst)
	return
}

func SelectJoinedTable(cmts *[]CommentJoin, tableName string, selectClause string, joinClause string, whereClause string, pid string, order string) {
	DB.Table(tableName).Select(selectClause).Joins(joinClause).Where(whereClause, pid).Order("created_at " + order).Find(cmts)
	return
}
