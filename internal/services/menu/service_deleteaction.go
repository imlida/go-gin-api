package menu

import (
	"github.com/imlida/go-gin-api/internal/pkg/core"
	"github.com/imlida/go-gin-api/internal/repository/mysql"
	"github.com/imlida/go-gin-api/internal/repository/mysql/menu_action"

	"gorm.io/gorm"
)

func (s *service) DeleteAction(ctx core.Context, id int32) (err error) {
	// 先查询 id 是否存在
	_, err = menu_action.NewQueryBuilder().
		WhereIsDeleted(mysql.EqualPredicate, -1).
		WhereId(mysql.EqualPredicate, id).
		First(s.db.GetDb("default").WithContext(ctx.RequestContext()))

	if err == gorm.ErrRecordNotFound {
		return nil
	}

	data := map[string]interface{}{
		"is_deleted":   1,
		"updated_user": ctx.SessionUserInfo().UserName,
	}

	qb := menu_action.NewQueryBuilder()
	qb.WhereId(mysql.EqualPredicate, id)
	err = qb.Updates(s.db.GetDb("default").WithContext(ctx.RequestContext()), data)
	if err != nil {
		return err
	}

	return
}
