package menu

import (
	"github.com/imlida/go-gin-api/internal/pkg/core"
	"github.com/imlida/go-gin-api/internal/repository/mysql"
	"github.com/imlida/go-gin-api/internal/repository/mysql/menu"
)

func (s *service) UpdateUsed(ctx core.Context, id int32, used int32) (err error) {
	data := map[string]interface{}{
		"is_used":      used,
		"updated_user": ctx.SessionUserInfo().UserName,
	}

	qb := menu.NewQueryBuilder()
	qb.WhereId(mysql.EqualPredicate, id)
	err = qb.Updates(s.db.GetDb("default").WithContext(ctx.RequestContext()), data)
	if err != nil {
		return err
	}

	return
}
