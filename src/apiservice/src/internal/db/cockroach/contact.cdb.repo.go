package cockroach

import (
	"context"

	"github.com/aqaurius6666/apiservice/src/internal/db/contact"
	"github.com/aqaurius6666/apiservice/src/internal/lib"
	"github.com/aqaurius6666/apiservice/src/internal/var/c"
	"github.com/aqaurius6666/go-utils/database"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel"
	"golang.org/x/xerrors"
	"gorm.io/gorm"
)

func applySearchContact(db *gorm.DB, search *contact.Search) *gorm.DB {
	if search.ID != uuid.Nil {
		db = db.Where(contact.Contact{
			BaseModel: database.BaseModel{
				ID: search.ID,
			},
		})
	}
	if len(search.Fields) > 0 {
		db = db.Select(search.Fields)
	}
	return db
}
func (u *ServerCDBRepo) DeleteContact(ctx context.Context, search *contact.Search) error {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(u.DeleteContact))
	defer span.End()
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	r := contact.Contact{}
	if err := applySearchContact(u.Db, search).WithContext(ctx).Delete(&r).Error; err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err, ctx)
		return err

	}
	return nil
}

func (u *ServerCDBRepo) SelectContact(ctx context.Context, search *contact.Search) (*contact.Contact, error) {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(u.SelectContact))
	defer span.End()
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	r := contact.Contact{}
	if err := applySearchContact(u.Db, search).WithContext(ctx).
		Joins(`left join "states" on "states"."id" = "contacts"."state_id"`).
		First(&r).
		Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			err = contact.ErrNotFound
		}
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err, ctx)
		return nil, err
	}
	return &r, nil
}

func (u *ServerCDBRepo) InsertContact(ctx context.Context, value *contact.Contact) (*contact.Contact, error) {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(u.InsertContact))
	defer span.End()
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	if err := u.Db.WithContext(ctx).Create(value).Error; err != nil {
		err = xerrors.Errorf("%w", contact.ErrInsertFail)
		lib.RecordError(span, err, ctx)
		return nil, err
	}
	return value, nil
}

func (u *ServerCDBRepo) ListContacts(ctx context.Context, search *contact.Search) ([]*contact.Contact, error) {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(u.ListContacts))
	defer span.End()
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	r := make([]*contact.Contact, 0)
	if err := applySearchContact(u.Db, search).WithContext(ctx).Find(&r).Error; err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err, ctx)
		return nil, err
	}
	return r, nil
}

func (u *ServerCDBRepo) UpdateContact(ctx context.Context, search *contact.Search, value *contact.Contact) error {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(u.UpdateContact))
	defer span.End()
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	if err := applySearchContact(u.Db, search).WithContext(ctx).Updates(value).Error; err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err, ctx)
		return err
	}
	return nil
}
