package repo

import (
	"context"

	"encore.app/gen/pgdb"
	"encore.app/internal/pkg/errs"
)

func (r *CaptainRepo) GetByPhone(ctx context.Context, phone string) (*pgdb.Captain, error) {
	captain, err := r.queries.GetCaptainByPhone(ctx, phone)
	if err != nil {
		if err.Error() == "no rows in result set" || err.Error() == "sql: no rows in result set" {
			return nil, ErrNotFound
		}
		return nil, errs.Internal(err, "failed to query captain by phone")
	}
	return &captain, nil
}

func (r *CaptainRepo) Create(ctx context.Context, name, phone string) (*pgdb.Captain, error) {
	err := r.queries.CreateCaptain(ctx, pgdb.CreateCaptainParams{
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
