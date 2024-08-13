package dao

import (
	"context"
	"fmt"
	"github.com/Filecoin-Titan/titan/api/types"
	"github.com/gnasnik/titan-container-api/core/generated/model"
	"strings"
)

func AddNewProvider(ctx context.Context, provider *types.Provider) error {
	qry := `INSERT INTO providers (id, owner, area_id, remote_addr, ip, state, created_at) 
		        VALUES (:id, :owner, :area_id, :remote_addr, :ip, :state, :created_at) ON DUPLICATE KEY UPDATE  owner=:owner, remote_addr=:remote_addr, 
		            ip=:ip, state=:state`
	_, err := DB.NamedExecContext(ctx, qry, provider)

	return err
}

func AddProviderWithResource(ctx context.Context, provider []*model.ProviderWithResource) error {
	qry := `INSERT INTO provider_with_res (id, area_id, remote_addr, ip, state, cpu, gpu, memory, storage, created_at) 
		        VALUES (:id, :area_id, :remote_addr, :ip, :state, :cpu, :gpu, :memory, :storage,:created_at) ON DUPLICATE KEY UPDATE  area_id = values(area_id), remote_addr=values(remote_addr), 
		            ip=values(ip), state=values(state)`
	_, err := DB.NamedExecContext(ctx, qry, provider)

	return err
}

func GetProvidersWithResource(ctx context.Context, areaId string, option model.QueryOption) (int64, []*model.ProviderWithResource, error) {
	qry := `SELECT * from provider_with_res`
	var condition []string

	condition = append(condition, "state = 1")

	if option.ID != "" {
		condition = append(condition, fmt.Sprintf(`id = '%s'`, option.ID))
	}

	if option.UserID != "" {
		condition = append(condition, fmt.Sprintf(`owner = '%s'`, option.UserID))
	}

	if areaId != "" && areaId != "ALL" {
		condition = append(condition, fmt.Sprintf(`area_id = '%s'`, areaId))
	}

	countSql := "select count(1) from provider_with_res "

	if len(condition) > 0 {
		qry += ` WHERE `
		qry += strings.Join(condition, ` AND `)
		countSql += ` WHERE `
		countSql += strings.Join(condition, ` AND `)
	}
	var total int64

	err := DB.GetContext(ctx, &total, countSql)
	if err != nil {
		return 0, nil, err
	}

	if option.Page <= 0 {
		option.Page = 1
	}

	if option.Size <= 0 {
		option.Size = 10
	}

	offset := (option.Page - 1) * option.Size
	limit := option.Size
	qry += fmt.Sprintf(" ORDER BY id LIMIT %d OFFSET %d", limit, offset)

	var out []*model.ProviderWithResource
	err = DB.SelectContext(ctx, &out, qry)
	if err != nil {
		return 0, nil, err
	}
	return total, out, nil
}
