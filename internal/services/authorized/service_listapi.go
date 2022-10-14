package authorized

import (
	"github.com/imlida/go-gin-api/internal/pkg/core"
	"github.com/imlida/go-gin-api/internal/repository/mysql"
	"github.com/imlida/go-gin-api/internal/repository/mysql/authorized_api"
)

type SearchAPIData struct {
	BusinessKey string `json:"business_key"` // 调用方key
}

func (s *service) ListAPI(ctx core.Context, searchAPIData *SearchAPIData) (listData []*authorized_api.AuthorizedApi, err error) {

	qb := authorized_api.NewQueryBuilder()
	qb = qb.WhereIsDeleted(mysql.EqualPredicate, -1)

	if searchAPIData.BusinessKey != "" {
		qb.WhereBusinessKey(mysql.EqualPredicate, searchAPIData.BusinessKey)
	}

	listData, err = qb.
		OrderById(false).
		QueryAll(s.db.GetDb("default").WithContext(ctx.RequestContext()))
	if err != nil {
		return nil, err
	}

	return
}
