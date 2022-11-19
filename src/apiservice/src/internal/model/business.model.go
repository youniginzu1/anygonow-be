package model

import (
	"context"
	"math/rand"
	"time"

	"github.com/aqaurius6666/apiservice/src/internal/db/business"
	"github.com/aqaurius6666/apiservice/src/internal/db/contact"
	"github.com/aqaurius6666/apiservice/src/internal/lib"
	"github.com/aqaurius6666/apiservice/src/internal/var/c"
	"github.com/aqaurius6666/apiservice/src/pb"
	"github.com/aqaurius6666/go-utils/database"
	"github.com/aqaurius6666/go-utils/utils"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel"
	"golang.org/x/xerrors"
)

var (
	_                             BusinessModel = (*ServerModel)(nil)
	RANDOM_INVITATION_CODE_LENGTH int           = 6
)

type BusinessModel interface {
	GetBusinessById(ctx context.Context, id interface{}) (*business.Business, error)
	CreateBusiness(ctx context.Context, id interface{}, mail, refCode, phone string) (*business.Business, error)
	ConvertBusinessToProto(*business.Business) *pb.Business
	ConvertBusinessToProtos([]*business.Business) []*pb.Business
	UpdateBusiness(ctx context.Context, id interface{}, b *business.Business) error
	GetNearBussiness(ctx context.Context, zipcode *string) ([]*business.Business, error)
	ListBusinesses(context.Context, *business.Search) ([]*business.Business, error)
	ListBusinessesWithRating(context.Context, *business.Search) ([]*business.Business, error)
	TotalBusinesses(context.Context, *business.Search) (*int64, error)
	GetBusinessByReferralCode(ctx context.Context, referralCode *string) (*business.Business, error)
	GetTotalZipcodes(context.Context, *business.Search) (*int64, error)
	GetBusiness(ctx context.Context, search *business.Search) (*business.Business, error)
}

func (s *ServerModel) ConvertBusinessToProtos(u []*business.Business) []*pb.Business {
	arr := make([]*pb.Business, 0)
	for _, uu := range u {
		arr = append(arr, s.ConvertBusinessToProto(uu))
	}
	return arr
}
func (s *ServerModel) ConvertBusinessToProto(u *business.Business) *pb.Business {
	upb := new(pb.Business)
	if u.ID != uuid.Nil {
		upb.Id = u.ID.String()
	}
	if u.Name != nil {
		upb.Name = *u.Name
	}
	if u.Contact != nil && u.Contact.Zipcode != nil {
		upb.Zipcode = *u.Contact.Zipcode
	}
	if u.Status != nil {
		upb.Status = c.ACCOUNT_STATUS(*u.Status)
	}
	if u.RefStatus != nil {
		upb.RefStatus = c.STATUS_VERIFY_REFERRAL_CODE(*u.RefStatus)
	}
	if u.Mail != nil {
		upb.Mail = *u.Mail
	}
	if u.BannerUrl != nil {
		upb.BannerImage = *u.BannerUrl
	}
	if u.ContactId != uuid.Nil {
		upb.ContactId = u.ContactId.String()
	}
	if u.Phone != nil {
		upb.Phone = *u.Phone
	}
	if u.LogoUrl != nil {
		upb.LogoImage = *u.LogoUrl
	}
	if u.Website != nil {
		upb.Website = *u.Website
	}
	if u.Description != nil {
		upb.Descriptions = *u.Description
	}
	if u.Zipcodes != nil {
		upb.Zipcodes = u.Zipcodes
	}
	if u.Services != nil {
		tmp := make([]string, 0)
		for _, s := range u.Services {
			tmp = append(tmp, s.String())
		}
		upb.Services = tmp
	}
	if u.Zipcode != nil {
		upb.Zipcode = *u.Zipcode
	}
	if u.ServiceName != nil && u.ServiceId != nil {
		for i := range u.ServiceName {
			upb.ServiceInfo = append(upb.ServiceInfo, &pb.ServiceGroup{
				ServiceName: u.ServiceName[i],
				ServiceId:   u.ServiceId[i],
			})
		}
	}
	if u.StartDate != nil {
		upb.StartDate = *u.StartDate
	}
	return upb
}

func (s *ServerModel) ListBusinessesWithRating(ctx context.Context, search *business.Search) ([]*business.Business, error) {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(s.ListBusinessesWithRating))
	defer span.End()

	b, err := s.Repo.ListBusinesssWithRating(ctx, search)
	if err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err, ctx)
		return nil, err
	}

	return b, nil
}

func (s *ServerModel) ListBusinesses(ctx context.Context, search *business.Search) ([]*business.Business, error) {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(s.ListBusinesses))
	defer span.End()

	b, err := s.Repo.ListBusinessOptimize(ctx, search)
	if err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err, ctx)
		return nil, err
	}

	return b, nil
}

func (s *ServerModel) TotalBusinesses(ctx context.Context, search *business.Search) (*int64, error) {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(s.TotalBusinesses))
	defer span.End()

	c, err := s.Repo.TotalBusiness(ctx, search)

	if err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err, ctx)
		return nil, err
	}

	return c, nil
}

func (s *ServerModel) GetBusinessByReferralCode(ctx context.Context, referralCode *string) (*business.Business, error) {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(s.GetBusinessByReferralCode))
	defer span.End()

	if referralCode == nil {
		return nil, nil
	}

	bus, err := s.Repo.SelectBusiness(ctx, &business.Search{
		Business: business.Business{
			InvitationCode: referralCode,
		},
		DefaultSearchModel: database.DefaultSearchModel{
			Fields: []string{`"businesses"."id"`, `"businesses"."invitation_code"`, `"businesses"."free_contact"`},
		},
	})

	if err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err, ctx)
		return nil, err
	}

	return bus, nil
}

func (s *ServerModel) GetNearBussiness(ctx context.Context, zipcode *string) ([]*business.Business, error) {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(s.GetNearBussiness))
	defer span.End()

	b, err := s.Repo.ListBusinesss(ctx, &business.Search{
		Zipcode: zipcode,
		DefaultSearchModel: database.DefaultSearchModel{
			Limit: 6,
			Fields: []string{`"businesses"."id"`,
				`"businesses"."name"`, `"businesses"."logo_url"`,
				`"businesses"."banner_url"`,
				`"businesses"."contact_id"`,
				`"businesses"."website"`,
				`"businesses"."description"`,
				`"businesses"."services"`,
				"rate",
				"review",
			},
		},
	})
	if err != nil {
		return nil, xerrors.Errorf("%w", err)
	}
	return b, nil
}
func (s *ServerModel) UpdateBusiness(ctx context.Context, id interface{}, u *business.Business) error {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(s.UpdateBusiness))
	defer span.End()

	uid, err := lib.ToUUID(id)
	if err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err, ctx)
		return err
	}
	err = s.Repo.UpdateBusiness(ctx, &business.Search{
		Business: business.Business{
			BaseModel: database.BaseModel{ID: uid},
		},
	}, u)
	if err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err, ctx)
		return err
	}
	return nil
}

func (s *ServerModel) CreateBusiness(ctx context.Context, id interface{}, mail, phone, refCode string) (*business.Business, error) {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(s.CreateBusiness))
	defer span.End()

	uid, err := lib.ToUUID(id)
	if err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err, ctx)
		return nil, err
	}

	invitationCode := CreateInvitationCodeForBusiness(RANDOM_INVITATION_CODE_LENGTH)
	rCode := utils.SafeStrPtr(refCode)

	refStatus := c.STATUS_VERIFY_REFERRAL_CODE_VERIFYING
	if rCode == nil {
		refStatus = c.STATUS_VERIFY_REFERRAL_CODE_NONE
	} else {
		_, err := s.GetBusinessByReferralCode(ctx, rCode)
		if err != nil {
			err = xerrors.Errorf("%w", err)
			lib.RecordError(span, err, ctx)
			return nil, err
		}
	}

	b, err := s.Repo.InsertBusiness(ctx, &business.Business{
		BaseModel:      database.BaseModel{ID: uid},
		ContactId:      uid,
		Mail:           &mail,
		Phone:          &phone,
		InvitationCode: utils.SafeStrPtr(invitationCode),
		RefCode:        rCode,
		RefStatus:      utils.Int32Ptr(int32(refStatus)),
	})
	if err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err, ctx)
		return nil, err
	}
	_, err = s.Repo.InsertContact(ctx, &contact.Contact{
		BaseModel: database.BaseModel{
			ID: uid,
		},
	})
	if err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err, ctx)
		return nil, err
	}
	return b, nil
}

func (s *ServerModel) GetTotalZipcodes(ctx context.Context, search *business.Search) (*int64, error) {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(s.GetTotalZipcodes))
	defer span.End()

	total, err := s.Repo.GetTotalZipcodes(ctx, search)
	if err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err, ctx)
		return nil, err
	}
	return total, nil
}

func CreateInvitationCodeForBusiness(n int) string {
	rand.Seed(time.Now().UnixNano())
	letterRunes := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func (s *ServerModel) GetBusinessById(ctx context.Context, id interface{}) (*business.Business, error) {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(s.GetBusinessById))
	defer span.End()

	uid, err := lib.ToUUID(id)
	if err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err, ctx)
		return nil, err
	}

	b, err := s.Repo.SelectBusiness(ctx, &business.Search{
		Business: business.Business{
			BaseModel: database.BaseModel{ID: uid},
		},
		DefaultSearchModel: database.DefaultSearchModel{
			Fields: []string{
				`"businesses"."id"`,
				"name",
				"phone",
				"logo_url",
				"banner_url",
				"website",
				"mail",
				"contact_id",
				"description",
				"services",
				"zipcodes",
				"status",
				"ref_code",
				"invitation_code",
				"free_contact",
				"ref_status",
			},
		},
	})
	if err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err, ctx)
		return nil, err
	}
	return b, nil
}

func (s *ServerModel) GetBusiness(ctx context.Context, search *business.Search) (*business.Business, error) {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(s.GetBusiness))
	defer span.End()

	u, err := s.Repo.SelectBusiness(ctx, search)
	if err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err, ctx)
		return nil, err
	}
	return u, nil
}
