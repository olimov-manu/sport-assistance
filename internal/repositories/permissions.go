package repositories

import (
	"context"
	"sport-assistance/pkg/myerrors"
)

func (r *Repository) GetPermissionsByRoleId(ctx context.Context, roleId uint64) ([]string, error) {
	const q = `
		SELECT p.name
		FROM role_permissions rp
		JOIN permissions p ON p.id = rp.permission_id
		WHERE rp.role_id = $1
		ORDER BY p.name
	`

	if err := ctx.Err(); err != nil {
		return nil, myerrors.NewRepositoryErr("контекст отменён перед выполнением запроса: ", err)
	}

	rows, err := r.postgres.Query(ctx, q, roleId)
	if err != nil {
		return nil, myerrors.NewRepositoryErr("не удалось получить permissions по role_id: ", err)
	}
	defer rows.Close()

	permissions := make([]string, 0)
	for rows.Next() {
		var permission string
		if err := rows.Scan(&permission); err != nil {
			return nil, myerrors.NewRepositoryErr("не удалось считать permission: ", err)
		}
		permissions = append(permissions, permission)
	}

	if err := rows.Err(); err != nil {
		return nil, myerrors.NewRepositoryErr("ошибка итерации по permissions: ", err)
	}

	return permissions, nil
}
