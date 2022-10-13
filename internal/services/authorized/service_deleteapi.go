package authorized

import (
	"github.com/imlida/go-gin-api/configs"
	"github.com/imlida/go-gin-api/internal/pkg/core"
	"github.com/imlida/go-gin-api/internal/repository/mysql"
	"github.com/imlida/go-gin-api/internal/repository/mysql/authorized_api"
	"github.com/imlida/go-gin-api/internal/repository/redis"

	"gorm.io/gorm"
)

func (s *service) DeleteAPI(ctx core.Context, id int32) (err error) {
	// 先查询 id 是否存在
	authorizedApiInfo, err := authorized_api.NewQueryBuilder().
		WhereIsDeleted(mysql.EqualPredicate, -1).
		WhereId(mysql.EqualPredicate, id).
		First(s.db.GetDb("Read").WithContext(ctx.RequestContext()))

	if err == gorm.ErrRecordNotFound {
		return nil
	}

	data := map[string]interface{}{
		"is_deleted":   1,
		"updated_user": ctx.SessionUserInfo().UserName,
	}

	qb := authorized_api.NewQueryBuilder()
	qb.WhereId(mysql.EqualPredicate, id)
	err = qb.Updates(s.db.GetDb("Write").WithContext(ctx.RequestContext()), data)
	if err != nil {
		return err
	}

	s.cache.Del(configs.RedisKeyPrefixSignature+authorizedApiInfo.BusinessKey, redis.WithTrace(ctx.Trace()))
	return
}
