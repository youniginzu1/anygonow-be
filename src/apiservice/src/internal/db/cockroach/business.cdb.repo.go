package cockroach

import (
	"context"
	"fmt"
	"time"

	"github.com/aqaurius6666/apiservice/src/internal/db/business"
	"github.com/aqaurius6666/apiservice/src/internal/db/feedback"
	"github.com/aqaurius6666/apiservice/src/internal/lib"
	"github.com/aqaurius6666/apiservice/src/internal/var/c"
	"github.com/aqaurius6666/go-utils/database"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel"
	"golang.org/x/xerrors"
	"gorm.io/gorm"
)

func applySearchBusiness(db *gorm.DB, search *business.Search) *gorm.DB {
	if search.ID != uuid.Nil {
		db = db.Where(business.Business{
			BaseModel: database.BaseModel{
				ID: search.ID,
			},
		})
	}

	if search.Status != nil {
		db = db.Where(business.Business{
			Status: search.Status,
		})
	}

	if search.Skip != 0 {
		db = db.Offset(search.Skip)
	}

	if search.Limit != 0 {
		db = db.Limit(search.Limit)
	}
	if search.Query == c.SORT_QUERY_DEFAULT && search.OrderBy != "" && search.OrderType != "" && search.CategoryId != uuid.Nil {
		db = db.Order(`"businesses"."start_date" DESC`)
		db = db.Order(fmt.Sprintf("%s %s", search.OrderBy, search.OrderType))
	}
	if search.Query == c.SORT_QUERY_REVIEW {
		db = db.Order(`"rating"."review" DESC NULLS LAST, "rating"."rate" DESC, "order"."request" DESC`)
	}
	if search.Query == c.SORT_QUERY_REQUEST {
		db = db.Order(`"order"."request" DESC NULLS LAST, "rating"."rate" DESC, "rating"."review" DESC`)
	}
	if search.InvitationCode != nil {
		db = db.Where(business.Business{
			InvitationCode: search.InvitationCode,
		})
	}
	if search.Mail != nil {
		db = db.Where(`"businesses"."mail" like ?`, *search.Mail+"%")
	}
	if search.Phone != nil {
		db = db.Where(`"businesses"."phone" like ?`, *search.Phone+"%")
	}
	if search.ValidateMail != nil {
		db = db.Where(`"businesses"."mail" = ?`, *search.ValidateMail)
	}
	return db
}

func getSubQueryOrder(db *gorm.DB, search *business.Search) *gorm.DB {
	db = db.Table(`"orders"`)
	db = db.Select(`"orders"."business_id",count(business_id) as "request"`)
	db = db.Group(`"orders"."business_id"`)
	db = db.Where(`"orders"."status" = ?`, c.ORDER_STATUS_COMPLETED)
	return db
}

func getSubQueryByCategory(db *gorm.DB, search *business.Search) *gorm.DB {
	var field = make([]string, 0)
	field = append(field, `DISTINCT ON("businesses"."id") "businesses".*`)
	if search.CategoryId == uuid.Nil {
		return db.Model(business.Business{})
	}
	db = db.Table("businesses")
	db = db.Joins(`left join "services" on "services"."business_id" = "businesses"."id" and "services"."status" = 0`)
	db = db.Where(`"services"."category_id" = ?`, search.CategoryId)
	if search.Zipcode != nil {
		db = db.Where(`cast(? as varchar) = any("businesses"."zipcodes" :: varchar[])`, *search.Zipcode)
	}
	if search.Query == c.SORT_QUERY_DEFAULT {
		now := time.Now().UnixMilli()
		db = db.Joins(`left join "advertise_orders" on "advertise_orders"."service_id" = any("businesses"."services" :: uuid[]) and "advertise_orders"."start_date" < ? and "advertise_orders"."end_date" > ?`, now, now)
		field = append(field, `"advertise_orders"."start_date"`)
	}
	db = db.Select(field)
	return db
}
func getSubQueryRating(db *gorm.DB, search *business.Search) *gorm.DB {
	db = db.Table("feedbacks")
	db = db.Select(`sum(rate) / cast(count(rate) as float8) as "rate", "business_id", count(*) as "review"`)
	db = db.Group(`"feedbacks"."business_id"`)
	return db
}

func (u *ServerCDBRepo) SelectBusiness(ctx context.Context, search *business.Search) (*business.Business, error) {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(u.SelectBusiness))
	defer span.End()
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	r := business.Business{}
	if err := applySearchBusiness(u.Db.Joins("Contact"), search).WithContext(ctx).Select(search.Fields).First(&r).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			err = business.ErrNotFound
		}
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err, ctx)
		return nil, err
	}
	return &r, nil
}

func (u *ServerCDBRepo) InsertBusiness(ctx context.Context, value *business.Business) (*business.Business, error) {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(u.InsertBusiness))
	defer span.End()
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	if err := u.Db.WithContext(ctx).Create(value).Error; err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err, ctx)
		return nil, err
	}
	return value, nil
}

func (u *ServerCDBRepo) ListBusinesss(ctx context.Context, search *business.Search) ([]*business.Business, error) {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(u.ListBusinesss))
	defer span.End()
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	r := make([]*business.Business, 0)

	subQuery := u.Db.Group("business_id").Select("sum(rate) / cast(count(rate) as float8) as rate, business_id, count(*) as review").Model(feedback.Feedback{})
	if err := applySearchBusiness(u.Db.Joins(`left join "contacts" on "contacts"."id" = "businesses"."contact_id"`), search).WithContext(ctx).
		Joins(`left join (?) as rating on "rating"."business_id" = "businesses"."id"`, subQuery).
		Joins(`left join "services" on "services"."id" = any("businesses"."services" :: uuid[])`).
		Joins(`left join "categories" on "services"."category_id" = "categories"."id"`).
		Group(`"businesses"."id", "businesses"."name", "logo_url", "banner_url", "contact_id", "website", "description", "businesses"."services", "mail", "businesses"."status", "contacts"."zipcode", "phone", "rate", "review", "businesses"."ref_status"`).
		Select(search.Fields).Find(&r).Error; err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err, ctx)
		return nil, err
	}
	return r, nil
}

func (u *ServerCDBRepo) ListBusinesssWithRating(ctx context.Context, search *business.Search) ([]*business.Business, error) {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(u.ListBusinesssWithRating))
	defer span.End()
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	r := make([]*business.Business, 0)
	db := u.Db
	if err := applySearchBusiness(db, search).Table(`(?) as "businesses"`, getSubQueryByCategory(db, search)).
		WithContext(ctx).
		Select(`"businesses".*, "contacts".*, "rating".*, "order"."request"`).
		Joins(`left join (?) as "rating" on "rating"."business_id" = "businesses"."id"`, getSubQueryRating(db, search)).
		Joins(`left join "contacts" on "businesses"."contact_id" = "contacts"."id"`).
		Joins(`left join (?) as "order" on "order"."business_id" = "businesses"."id"`, getSubQueryOrder(db, search)).
		Find(&r).Error; err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err, ctx)
		return nil, err
	}
	return r, nil
}

func (u *ServerCDBRepo) DeleteBusiness(ctx context.Context, search *business.Search) error {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(u.DeleteBusiness))
	defer span.End()
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	if err := applySearchBusiness(u.Db, search).WithContext(ctx).Delete(&business.Business{}).Error; err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err, ctx)
		return err
	}
	return nil
}

func (u *ServerCDBRepo) UpdateBusiness(ctx context.Context, search *business.Search, value *business.Business) error {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(u.UpdateBusiness))
	defer span.End()
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	if err := applySearchBusiness(u.Db, search).WithContext(ctx).Updates(value).Error; err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err, ctx)
		return err
	}
	return nil
}

func (u *ServerCDBRepo) TotalBusiness(ctx context.Context, search *business.Search) (*int64, error) {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(u.TotalBusiness))
	defer span.End()
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	var r int64
	if err := applySearchBusiness(getSubQueryByCategory(u.Db, search), search).WithContext(ctx).Model(business.Business{}).Count(&r).Error; err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err, ctx)
		return nil, err
	}
	return &r, nil
}

func (u *ServerCDBRepo) GetMapIdName(ctx context.Context, search *business.Search) (map[string]*string, error) {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(u.GetMapIdName))
	defer span.End()
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	res := make([]*business.Business, 0)
	if err := u.Db.WithContext(ctx).
		Raw(`select id, name from ((?) union (?))`,
			u.Db.Raw(`select id, name from "businesses"`),
			u.Db.Raw(`select id, trim(concat(first_name, ' ', last_name)) as name from "users"`),
		).
		Where(` = any(? :: uuid[])`, search.BothId).Scan(&res).Error; err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err, ctx)
		return nil, err
	}
	v := make(map[string]*string)
	for _, r := range res {
		v[r.ID.String()] = r.Name
	}
	return v, nil
}

func (u *ServerCDBRepo) GetTotalZipcodes(ctx context.Context, search *business.Search) (*int64, error) {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(u.GetTotalZipcodes))
	defer span.End()
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	r := business.Business{}
	if err := u.Db.WithContext(ctx).Raw(`select count(*) as "count_zipcodes" from (select distinct unnest(zipcodes) from businesses) as "aa"`).Model(business.Business{}).First(&r).Error; err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err, ctx)
		return nil, err
	}
	return r.CountZipcodes, nil
}

func (u *ServerCDBRepo) ListBusinessOptimize(ctx context.Context, search *business.Search) ([]*business.Business, error) {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(u.ListBusinessOptimize))
	defer span.End()
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	r := make([]*business.Business, 0)

	db := u.Db
	if err := applySearchBusiness(db, search).WithContext(ctx).
		Select(search.Fields).
		Joins(`left join (select "id", "zipcode" from "contacts") as "contacts" on "contacts"."id" = "businesses"."contact_id"`).
		Joins(`left join (?) as "rating" on "rating"."business_id" = "businesses"."id"`, getSubQueryRating(db, search)).
		Joins(`left join (?) as "aggregation" on "aggregation"."id" = "businesses"."id"`, getSubQueryAgg(db)).Find(&r).Error; err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err, ctx)
		return nil, err
	}

	return r, nil
}

func getSubQueryAgg(db *gorm.DB) *gorm.DB {
	db = db.Select([]string{
		`services."business_id" as "id"`,
		`array_agg("c"."id")        as "service_id"`,
		`array_agg("c"."name")      as "service_name"`,
	})
	db = db.Table("services")
	db = db.Joins(`left join "categories" c on "services"."category_id" = "c"."id" group by "services"."business_id"`)
	return db
}
