package repo

import (
	"context"

	"encore.app/gen/pgdb"
	"encore.app/internal/pkg/errs"
)

func (r *UserRepo) GetByPhone(ctx context.Context, phone string) (*pgdb.User, error) {
	captain, err := r.queries.GetUserByPhone(ctx, phone)

	if err != nil {
		if err.Error() == "no rows in result set" || err.Error() == "sql: no rows in result set" {
			return nil, ErrNotFound
		}
		return nil, errs.Internal(err, "failed to query captain by phone")
	}

	return &captain, nil
}

func (r *UserRepo) Create(ctx context.Context, name, phone string) (*pgdb.User, error) {

	err := r.queries.CreateUser(ctx, pgdb.CreateUserParams{
		Name:  name,
		Phone: phone,
	})

	if err != nil {
		return nil, errs.Internal(err, "failed to create captain")
	}

	captain, err := r.GetByPhone(ctx, phone)
	if err != nil {
		return nil, errs.Internal(err, "failed to fetch created captain")
	}

	return captain, nil
}
