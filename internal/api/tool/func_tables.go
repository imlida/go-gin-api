/*
 * @Author: lida lidaemail@qq.com
 * @LastEditors: lida lidaemail@qq.com
 */
package tool

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/imlida/go-gin-api/internal/code"
	"github.com/imlida/go-gin-api/internal/pkg/core"
)

type tablesRequest struct {
	DbName string `form:"db_name"` // 数据库名称
}

type tablesResponse struct {
	List []tableData `json:"list"` // 数据表列表
}

type tableData struct {
	Name    string `json:"table_name"`    // 数据表名称
	Comment string `json:"table_comment"` // 数据表备注
}

// Tables 查询 Table
// @Summary 查询 Table
// @Description 查询 Table
// @Tags API.tool
// @Accept application/x-www-form-urlencoded
// @Produce json
// @Param db_name formData string true "数据库名称"
// @Success 200 {object} tablesResponse
// @Failure 400 {object} code.Failure
// @Router /api/tool/data/tables [post]
// @Security LoginToken
func (h *handler) Tables() core.HandlerFunc {
	return func(c core.Context) {
		req := new(tablesRequest)
		res := new(tablesResponse)
		if err := c.ShouldBindForm(req); err != nil {
			c.AbortWithError(core.Error(
				http.StatusBadRequest,
				code.ParamBindError,
				code.Text(code.ParamBindError)).WithError(err),
			)
			return
		}

		//分割字符串
		dbName := strings.Split(req.DbName, "|")

		sqlTables := fmt.Sprintf("SELECT `table_name`,`table_comment` FROM `information_schema`.`tables` WHERE `table_schema`= '%s'", strings.TrimSpace(dbName[1]))

		// TODO 后期支持查询多个数据库
		rows, err := h.db.GetDb(strings.TrimSpace(dbName[0])).Raw(sqlTables).Rows()
		if err != nil {
			c.AbortWithError(core.Error(
				http.StatusBadRequest,
				code.MySQLExecError,
				code.Text(code.MySQLExecError)).WithError(err),
			)
			return
		}

		defer rows.Close()

		for rows.Next() {
			var info tableData
			err = rows.Scan(&info.Name, &info.Comment)
			if err != nil {
				fmt.Printf("execute query tables action error,had ignored, detail is [%v]\n", err.Error())
				continue
			}

			res.List = append(res.List, info)
		}

		c.Payload(res)
	}
}
