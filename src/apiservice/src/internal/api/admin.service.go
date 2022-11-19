package api

import (
	"context"

	"github.com/aqaurius6666/apiservice/src/internal/db/advertise_package"
	"github.com/aqaurius6666/apiservice/src/internal/db/business"
	"github.com/aqaurius6666/apiservice/src/internal/db/category"
	"github.com/aqaurius6666/apiservice/src/internal/db/group"
	"github.com/aqaurius6666/apiservice/src/internal/db/user"
	"github.com/aqaurius6666/apiservice/src/internal/lib"
	"github.com/aqaurius6666/apiservice/src/internal/lib/validate"
	"github.com/aqaurius6666/apiservice/src/internal/model"
	"github.com/aqaurius6666/apiservice/src/internal/var/c"
	"github.com/aqaurius6666/apiservice/src/internal/var/e"
	"github.com/aqaurius6666/apiservice/src/pb"
	"github.com/aqaurius6666/go-utils/database"
	"github.com/aqaurius6666/go-utils/utils"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel"
	"golang.org/x/xerrors"
)

type AdminService struct {
	Model model.Server
}

func (s *AdminService) DeleteUser(ctx context.Context, req *pb.AdminUsersDeletePostRequest) (*pb.AdminUsersDeletePostResponse_Data, error) {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(s.DeleteUser))
	defer span.End()

	if f, ok := validate.RequiredFields(req, "Id", "XUserId"); !ok {
		err := xerrors.Errorf("%w", e.ErrMissingField(f))
		lib.RecordError(span, err, ctx)
		return nil, err
	}
	b, err := s.Model.GetUserById(ctx, req.Id)
	if err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err, ctx)
		return nil, err
	}
	if *b.Status == int32(c.ACCOUNT_STATUS_ACTIVE) {
		err = xerrors.Errorf("%w", e.ErrBodyInvalid)
		lib.RecordError(span, err, ctx)
		return nil, err
	}
	err = s.Model.DeleteUser(ctx, req.Id)
	if err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err, ctx)
		return nil, err
	}

	return &pb.AdminUsersDeletePostResponse_Data{}, nil
}

func (s *AdminService) DeleteBusiness(ctx context.Context, req *pb.AdminBusinessDeletePostRequest) (*pb.AdminBusinessDeletePostResponse_Data, error) {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(s.DeleteBusiness))
	defer span.End()

	if f, ok := validate.RequiredFields(req, "Id", "XUserId"); !ok {
		err := xerrors.Errorf("%w", e.ErrMissingField(f))
		lib.RecordError(span, err, ctx)
		return nil, err
	}
	b, err := s.Model.GetBusinessById(ctx, req.Id)
	if err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err, ctx)
		return nil, err
	}
	if *b.Status == int32(c.ACCOUNT_STATUS_ACTIVE) {
		err = xerrors.Errorf("%w", e.ErrBodyInvalid)
		lib.RecordError(span, err, ctx)
		return nil, err
	}
	err = s.Model.DeleteBusiness(ctx, req.Id)
	if err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err, ctx)
		return nil, err
	}

	return &pb.AdminBusinessDeletePostResponse_Data{}, nil
}

func (s *AdminService) BanBusiness(ctx context.Context, req *pb.AdminBusinessBanPostRequest) (*pb.AdminBusinessBanPostResponse_Data, error) {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(s.BanBusiness))
	defer span.End()

	if f, ok := validate.RequiredFields(req, "Id", "XUserId"); !ok {
		err := xerrors.Errorf("%w", e.ErrMissingField(f))
		lib.RecordError(span, err, ctx)
		return nil, err
	}
	b, err := s.Model.GetBusinessById(ctx, req.Id)
	if err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err, ctx)
		return nil, err
	}
	if *b.Status == int32(c.ACCOUNT_STATUS_INACTIVE) {
		err = xerrors.Errorf("%w", e.ErrBodyInvalid)
		lib.RecordError(span, err, ctx)
		return nil, err
	}
	_, _, err = s.Model.InactiveBusiness(ctx, req.Id)
	if err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err, ctx)

		return nil, err
	}
	err = s.Model.UpdateOrdersStatusByUser(ctx, req.Id, c.ORDER_STATUS_CANCELED)
	if err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err, ctx)
		return nil, err
	}

	return &pb.AdminBusinessBanPostResponse_Data{}, nil
}

func (s *AdminService) UnbanUser(ctx context.Context, req *pb.AdminUsersUnbanPostRequest) (*pb.AdminUsersUnbanPostResponse_Data, error) {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(s.UnbanUser))
	defer span.End()

	if f, ok := validate.RequiredFields(req, "Id", "XUserId"); !ok {
		return nil, e.ErrMissingField(f)
	}
	u, err := s.Model.GetUserById(ctx, req.Id)
	if err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err, ctx)
		return nil, err
	}
	if *u.Status == int32(c.ACCOUNT_STATUS_ACTIVE) {
		err = xerrors.Errorf("%w", e.ErrBodyInvalid)
		lib.RecordError(span, err, ctx)
		return nil, err
	}
	_, _, err = s.Model.ActiveUser(ctx, req.Id)
	if err != nil {
		return nil, xerrors.Errorf("%w", err)
	}

	return &pb.AdminUsersUnbanPostResponse_Data{}, nil
}

func (s *AdminService) UnbanBusiness(ctx context.Context, req *pb.AdminBusinessesUnbanPostRequest) (*pb.AdminBusinessesUnbanPostResponse_Data, error) {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(s.UnbanBusiness))
	defer span.End()

	if f, ok := validate.RequiredFields(req, "Id", "XUserId"); !ok {
		err := xerrors.Errorf("%w", e.ErrMissingField(f))
		lib.RecordError(span, err, ctx)
		return nil, err
	}
	b, err := s.Model.GetBusinessById(ctx, req.Id)
	if err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err, ctx)
		return nil, err
	}
	if *b.Status == int32(c.ACCOUNT_STATUS_ACTIVE) {
		err = xerrors.Errorf("%w", e.ErrBodyInvalid)
		lib.RecordError(span, err, ctx)
		return nil, err
	}
	_, _, err = s.Model.ActiveBusiness(ctx, req.Id)
	if err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err, ctx)

		return nil, err
	}

	return &pb.AdminBusinessesUnbanPostResponse_Data{}, nil
}

func (s *AdminService) BanUser(ctx context.Context, req *pb.AdminBanUserPostRequest) (*pb.AdminBanUserPostResponse_Data, error) {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(s.BanUser))
	defer span.End()

	if f, ok := validate.RequiredFields(req, "Id", "XUserId"); !ok {
		err := xerrors.Errorf("%w", e.ErrMissingField(f))
		lib.RecordError(span, err, ctx)
		return nil, err
	}
	u, err := s.Model.GetUserById(ctx, req.Id)
	if err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err, ctx)
		return nil, err
	}
	if *u.Status == int32(c.ACCOUNT_STATUS_INACTIVE) {
		err = xerrors.Errorf("%w", e.ErrBodyInvalid)
		lib.RecordError(span, err, ctx)
		return nil, err
	}
	_, _, err = s.Model.InactiveUser(ctx, req.Id)
	if err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err, ctx)
		return nil, err
	}
	err = s.Model.UpdateOrdersStatusByUser(ctx, req.Id, c.ORDER_STATUS_CANCELED)
	if err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err, ctx)
		return nil, err
	}

	return &pb.AdminBanUserPostResponse_Data{}, nil
}

func (s *AdminService) ListBusinesses(ctx context.Context, req *pb.AdminBusinessesGetRequest) (*pb.AdminBusinessesGetResponse_Data, error) {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(s.ListBusinesses))
	defer span.End()
	if f, ok := validate.RequiredFields(req, "XUserId"); !ok {
		err := xerrors.Errorf("%w", e.ErrMissingField(f))
		lib.RecordError(span, err, ctx)
		return nil, err
	}
	limit := lib.ParseInt32Val(req.Limit)
	offset := lib.ParseInt32Val(req.Offset)

	total, err := s.Model.TotalBusinesses(ctx, &business.Search{
		Business: business.Business{
			Phone: utils.SafeStrPtr(req.Phone),
			Mail:  utils.SafeStrPtr(req.Mail),
		},
	})
	if err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err, ctx)
		return nil, err
	}
	bus, err := s.Model.ListBusinesses(ctx, &business.Search{
		DefaultSearchModel: database.DefaultSearchModel{
			Skip:  int(offset),
			Limit: int(limit),
			Fields: []string{
				`"businesses"."id"`,
				`"businesses"."name"`,
				"logo_url",
				"banner_url",
				"contact_id",
				"website",
				"description",
				"mail",
				`"businesses"."status"`,
				`"businesses"."ref_status"`,
				`"contacts"."zipcode"`,
				`"phone"`,
				`"service_name"`,
				`"service_id"`,
			},
		},
		Business: business.Business{
			Phone: utils.SafeStrPtr(req.Phone),
			Mail:  utils.SafeStrPtr(req.Mail),
		},
	})
	if err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err, ctx)
		return nil, err
	}

	return &pb.AdminBusinessesGetResponse_Data{
		Result:     s.Model.ConvertBusinessToProtos(bus),
		Pagination: lib.Pagination(offset, limit, total),
	}, nil
}

func (s *AdminService) ListUsers(ctx context.Context, req *pb.AdminUsersGetRequest) (*pb.AdminUsersGetResponse_Data, error) {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(s.ListUsers))
	defer span.End()
	if f, ok := validate.RequiredFields(req, "XUserId"); !ok {
		err := xerrors.Errorf("%w", e.ErrMissingField(f))
		lib.RecordError(span, err, ctx)
		return nil, err
	}
	limit := lib.ParseInt32Val(req.Limit)
	offset := lib.ParseInt32Val(req.Offset)

	total, err := s.Model.TotalUsers(ctx, &user.Search{
		DefaultSearchModel: database.DefaultSearchModel{},
		User: user.User{
			Mail:  utils.SafeStrPtr(req.Mail),
			Phone: utils.SafeStrPtr(req.Phone),
		},
	})
	if err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err, ctx)
		return nil, err
	}

	users, err := s.Model.ListUsers(ctx, &user.Search{
		DefaultSearchModel: database.DefaultSearchModel{
			Skip:   int(offset),
			Limit:  int(limit),
			Fields: []string{`"users"."id"`, "first_name", "last_name", "phone", "mail", "status", `"Contact"."zipcode"`},
		},
		User: user.User{
			Mail:  utils.SafeStrPtr(req.Mail),
			Phone: utils.SafeStrPtr(req.Phone),
		},
	})

	if err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err, ctx)
		return nil, err
	}

	return &pb.AdminUsersGetResponse_Data{
		Result:     s.Model.ConvertUsersToProtos(users),
		Pagination: lib.Pagination(offset, limit, total),
	}, nil
}

func (s *AdminService) ListCategory(ctx context.Context, req *pb.AdminBusinessesGetRequest) (*pb.AdminBusinessesGetResponse_Data, error) {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(s.ListBusinesses))
	defer span.End()
	if f, ok := validate.RequiredFields(req, "XUserId"); !ok {
		err := xerrors.Errorf("%w", e.ErrMissingField(f))
		lib.RecordError(span, err, ctx)
		return nil, err
	}
	limit := lib.ParseInt32Val(req.Limit)
	offset := lib.ParseInt32Val(req.Offset)

	total, err := s.Model.TotalBusinesses(ctx, &business.Search{
		Business: business.Business{
			Phone: utils.SafeStrPtr(req.Phone),
			Mail:  utils.SafeStrPtr(req.Mail),
		},
	})
	if err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err, ctx)
		return nil, err
	}
	bus, err := s.Model.ListBusinesses(ctx, &business.Search{
		DefaultSearchModel: database.DefaultSearchModel{
			Skip:   int(offset),
			Limit:  int(limit),
			Fields: []string{`"businesses"."id"`, "name", "logo_url", "banner_url", "contact_id", "website", "description", "services", "mail", "status", `"contacts"."zipcode", "phone"`},
		},
		Business: business.Business{
			Phone: utils.SafeStrPtr(req.Phone),
			Mail:  utils.SafeStrPtr(req.Mail),
		},
	})
	if err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err, ctx)
		return nil, err
	}

	return &pb.AdminBusinessesGetResponse_Data{
		Result:     s.Model.ConvertBusinessToProtos(bus),
		Pagination: lib.Pagination(offset, limit, total),
	}, nil
}

func (s *AdminService) ListCategories(ctx context.Context, req *pb.CategoriesGetRequest) (*pb.CategoriesGetResponse_Data, error) {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(s.ListCategories))
	defer span.End()

	limit := lib.ParseInt32Val(req.Limit)
	offset := lib.ParseInt32Val(req.Offset)

	total, err := s.Model.TotalCategory(ctx, &category.Search{
		Query: req.Query,
		Category: category.Category{
			Name: utils.SafeStrPtr(lib.StandardizeSpaces(req.Name)),
		},
	})

	if err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err, ctx)
		return nil, err
	}

	r, err := s.Model.GetCategories(ctx, &category.Search{
		DefaultSearchModel: database.DefaultSearchModel{
			Skip:   int(offset),
			Limit:  int(limit),
			Fields: []string{`"categories"."image_url","categories"."name", count("services"."business_id") as "total_provider", "categories"."id", "groups"."fee"`},
		},
		Query: req.Query,
		Category: category.Category{
			Name: utils.SafeStrPtr(lib.StandardizeSpaces(req.Name)),
		},
	})
	if err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err, ctx)
		return nil, err
	}

	return &pb.CategoriesGetResponse_Data{
		Result:     s.Model.ConvertCategoriesToProtos(r),
		Pagination: lib.Pagination(offset, limit, total),
	}, nil

}

func (s *AdminService) PostCategory(ctx context.Context, req *pb.AdminCategoryPostRequest) (*pb.AdminCategoryPostResponese_Data, error) {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(s.PostCategory))
	defer span.End()

	if f, ok := validate.RequiredFields(req, "XUserId", "Name", "Image"); !ok {
		err := xerrors.Errorf("%w", e.ErrMissingField(f))
		lib.RecordError(span, err, ctx)
		return nil, err
	}

	categoryName := lib.StandardizeSpaces(req.Name)

	//Get category with that name and compare
	_, err := s.Model.GetCategoryByName(ctx, &categoryName)
	if err == nil {
		err = xerrors.Errorf("%w", e.ErrCategoryExisted)
		lib.RecordError(span, err, ctx)
		return nil, err
	}

	// Insert
	err = s.Model.InsertCategory(ctx, &category.Category{
		Name:     &categoryName,
		ImageUrl: utils.SafeStrPtr(req.Image),
	})
	if err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err, ctx)
		return nil, err
	}

	return &pb.AdminCategoryPostResponese_Data{}, nil
}

func (s *AdminService) DeleteCategory(ctx context.Context, req *pb.AdminCategoryPostDeleteRequest) (*pb.AdminCategoryPostDeleteResponese_Data, error) {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(s.DeleteCategory))
	defer span.End()

	if f, ok := validate.RequiredFields(req, "XUserId", "CategoryId"); !ok {
		err := xerrors.Errorf("%w", e.ErrMissingField(f))
		lib.RecordError(span, err, ctx)
		return nil, err
	}

	err := s.Model.DeleteCategory(ctx, req.CategoryId)
	if err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err, ctx)
		return nil, err
	}

	return &pb.AdminCategoryPostDeleteResponese_Data{}, nil
}

func (s *AdminService) EditCategory(ctx context.Context, req *pb.AdminCategoryPostEditRequest) (*pb.AdminCategoryPostEditResponese_Data, error) {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(s.DeleteCategory))
	defer span.End()

	if f, ok := validate.RequiredFields(req, "XUserId"); !ok {
		err := xerrors.Errorf("%w", e.ErrMissingField(f))
		lib.RecordError(span, err, ctx)
		return nil, err
	}

	cid, err := lib.ToUUID(req.Id)
	if err != nil {
		return nil, xerrors.Errorf("%w", err)
	}

	err = s.Model.EditCategory(ctx, &category.Search{
		Category: category.Category{
			BaseModel: database.BaseModel{
				ID: cid,
			},
		},
	}, &category.Category{
		ImageUrl: utils.SafeStrPtr(req.Url),
		Name:     utils.SafeStrPtr(req.Name),
	})
	if err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err, ctx)
		return nil, err
	}

	return &pb.AdminCategoryPostEditResponese_Data{}, nil
}

func (s *AdminService) AddGroup(ctx context.Context, req *pb.AdminGroupPostRequest) (*pb.AdminGroupPostResponse_Data, error) {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(s.AddGroup))
	defer span.End()

	if f, ok := validate.RequiredFields(req, "XUserId", "Name", "Fee", "ServiceIds"); !ok {
		err := xerrors.Errorf("%w", e.ErrMissingField(f))
		lib.RecordError(span, err, ctx)
		return nil, err
	}

	for _, ele := range req.ServiceIds {
		cid, err := lib.ToUUID(ele)
		if err != nil {
			return nil, xerrors.Errorf("%w", err)
		}
		err = s.Model.CheckCategoryIdExisted(ctx, &group.Search{}, cid)
		if err == nil {
			return nil, e.ErrCategoryExisted
		}
	}

	grName := lib.StandardizeSpaces(req.Name)

	if _, err := s.Model.SelectGroup(ctx, &group.Search{
		Group: group.Group{
			Name: &grName,
		},
	}); err == nil {
		return nil, e.ErrGroupExisted
	}

	cIds := make([]uuid.UUID, 0)
	for _, s := range req.ServiceIds {
		cid, err := lib.ToUUID(s)
		if err != nil {
			return nil, xerrors.Errorf("%w", err)
		}
		cIds = append(cIds, cid)
	}

	err := s.Model.InsertGroup(ctx, &group.Group{
		Name:        &grName,
		Fee:         &req.Fee,
		CategoryIds: cIds,
	})
	if err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err, ctx)
		return nil, err
	}

	return &pb.AdminGroupPostResponse_Data{}, nil
}

func (s *AdminService) EditGroup(ctx context.Context, req *pb.AdminGroupPutRequest) (*pb.AdminGroupPutResponse_Data, error) {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(s.EditGroup))
	defer span.End()

	if f, ok := validate.RequiredFields(req, "XUserId", "Id", "ServiceIds"); !ok {
		err := xerrors.Errorf("%w", e.ErrMissingField(f))
		lib.RecordError(span, err, ctx)
		return nil, err
	}

	grName := lib.StandardizeSpaces(req.Name)

	gid, err := lib.ToUUID(req.Id)
	if err != nil {
		return nil, xerrors.Errorf("%w", err)
	}

	for _, ele := range req.ServiceIds {
		cid, err := lib.ToUUID(ele)
		if err != nil {
			return nil, xerrors.Errorf("%w", err)
		}
		err = s.Model.CheckCategoryIdExisted(ctx, &group.Search{
			Group: group.Group{
				BaseModel: database.BaseModel{
					ID: gid,
				},
			},
		}, cid)
		if err == nil {
			return nil, e.ErrCategoryExisted
		}
	}

	cIds := make([]uuid.UUID, 0)
	for _, s := range req.ServiceIds {
		cid, err := lib.ToUUID(s)
		if err != nil {
			return nil, xerrors.Errorf("%w", err)
		}
		cIds = append(cIds, cid)
	}

	err = s.Model.EditGroup(ctx, &group.Search{
		Group: group.Group{
			BaseModel: database.BaseModel{
				ID: gid,
			},
		},
	}, &group.Group{
		Name:        &grName,
		Fee:         &req.Fee,
		CategoryIds: cIds,
	})
	if err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err, ctx)
		return nil, err
	}

	return &pb.AdminGroupPutResponse_Data{}, nil
}

func (s *AdminService) GetGroups(ctx context.Context, req *pb.AdminGroupGetRequest) (*pb.AdminGroupGetResponse_Data, error) {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(s.GetGroups))
	defer span.End()

	if f, ok := validate.RequiredFields(req, "XUserId"); !ok {
		err := xerrors.Errorf("%w", e.ErrMissingField(f))
		lib.RecordError(span, err, ctx)
		return nil, err
	}

	limit := lib.ParseInt32Val(req.Limit)
	offset := lib.ParseInt32Val(req.Offset)

	total, err := s.Model.TotalGroup(ctx, &group.Search{})
	if err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err, ctx)
		return nil, err
	}

	res, err := s.Model.ListGroups(ctx, &group.Search{
		DefaultSearchModel: database.DefaultSearchModel{
			Limit: int(limit),
			Skip:  int(offset),
			Fields: []string{`"groups"."id"`,
				`"groups"."name"`,
				`"groups"."fee"`,
				`array_remove(array_agg("categories"."name"), null) as "service_name"`,
				`array_remove(array_agg("categories"."id"), null) as "service_id"`},
		},
		CategoryId: lib.ParseUUID(req.CategoryId),
	})
	if err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err, ctx)
		return nil, err
	}

	return &pb.AdminGroupGetResponse_Data{
		Result:     s.Model.ConvertGroupToProtos(res),
		Pagination: lib.Pagination(offset, limit, total),
	}, nil
}

func (s *AdminService) CreateAdvertisePackage(ctx context.Context, req *pb.AdminAdvertiseManagementPostRequest) (*pb.AdminAdvertiseManagementPostResponse_Data, error) {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(s.CreateAdvertisePackage))
	defer span.End()

	if f, ok := validate.RequiredFields(req, "XUserId", "Name", "Price"); !ok {
		err := xerrors.Errorf("%w", e.ErrMissingField(f))
		lib.RecordError(span, err, ctx)
		return nil, err
	}

	if len(req.CategoryIds) == 0 {
		err := xerrors.Errorf("%w", e.ErrMissingField("CategoryIds"))
		lib.RecordError(span, err, ctx)
		return nil, err
	}

	for _, ele := range req.CategoryIds {
		cid, err := lib.ToUUID(ele)
		if err != nil {
			return nil, xerrors.Errorf("%w", err)
		}
		err = s.Model.CheckCateIdExistedAdvertise(ctx, &advertise_package.Search{}, cid)
		if err == nil {
			return nil, e.ErrCategoryExisted
		}
	}

	ucids := make([]uuid.UUID, 0)

	for _, s := range req.CategoryIds {
		ucid, err := lib.ToUUID(s)
		if err != nil {
			return nil, xerrors.Errorf("%w", err)
		}
		ucids = append(ucids, ucid)
	}

	pName := lib.StandardizeSpaces(req.Name)

	err := s.Model.InsertPackage(ctx, &advertise_package.AdvertisePackage{
		Name:        &pName,
		Categories:  ucids,
		Price:       &req.Price,
		BannerUrl:   utils.SafeStrPtr(req.BannerUrl),
		Description: &req.Description,
	})

	if err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err, ctx)
		return nil, err
	}

	return &pb.AdminAdvertiseManagementPostResponse_Data{}, nil
}

func (s *AdminService) GetAdvertisePackages(ctx context.Context, req *pb.AdminAdvertiseManagementGetRequest) (*pb.AdminAdvertiseManagementGetResponse_Data, error) {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(s.GetAdvertisePackages))
	defer span.End()

	if f, ok := validate.RequiredFields(req, "XUserId"); !ok {
		err := xerrors.Errorf("%w", e.ErrMissingField(f))
		lib.RecordError(span, err, ctx)
		return nil, err
	}

	limit := lib.ParseInt32Val(req.Limit)
	offset := lib.ParseInt32Val(req.Offset)

	total, err := s.Model.TotalPackage(ctx, &advertise_package.Search{
		DefaultSearchModel: database.DefaultSearchModel{
			Limit: int(limit),
			Skip:  int(offset),
			Fields: []string{`"advertise_packages"."id"`,
				`"advertise_packages"."name"`,
				`"advertise_packages"."price"`,
				`"advertise_packages"."banner_url"`,
				`"advertise_packages"."description"`,
				`array_remove(array_agg("categories"."name"), null) as "service_name"`,
				`array_remove(array_agg("categories"."id"), null) as "service_id"`},
		},

		ServiceName: req.ServiceName,
	})
	if err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err, ctx)
		return nil, err
	}

	res, err := s.Model.ListPackages(ctx, &advertise_package.Search{
		DefaultSearchModel: database.DefaultSearchModel{
			Limit: int(limit),
			Skip:  int(offset),
			Fields: []string{`"advertise_packages"."id"`,
				`"advertise_packages"."name"`,
				`"advertise_packages"."price"`,
				`"advertise_packages"."banner_url"`,
				`"advertise_packages"."description"`,
				`array_remove(array_agg("categories"."name"), null) as "service_name"`,
				`array_remove(array_agg("categories"."id"), null) as "service_id"`},
		},

		ServiceName: req.ServiceName,
	})
	if err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err, ctx)
		return nil, err
	}

	return &pb.AdminAdvertiseManagementGetResponse_Data{
		Result:     s.Model.ConvertPackageToProtos(res),
		Pagination: lib.Pagination(offset, limit, total),
	}, nil
}

func (s *AdminService) EditAdvertisePackage(ctx context.Context, req *pb.AdminAdvertiseManagementPutRequest) (*pb.AdminAdvertiseManagementPutResponse_Data, error) {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(s.EditAdvertisePackage))
	defer span.End()

	if f, ok := validate.RequiredFields(req, "XUserId", "Id", "Name", "Price"); !ok {
		err := xerrors.Errorf("%w", e.ErrMissingField(f))
		lib.RecordError(span, err, ctx)
		return nil, err
	}

	pid, err := lib.ToUUID(req.Id)
	if err != nil {
		return nil, xerrors.Errorf("%w", err)
	}

	if len(req.CategoryIds) == 0 {
		err := xerrors.Errorf("%w", e.ErrMissingField("CategoryIds"))
		lib.RecordError(span, err, ctx)
		return nil, err
	}

	for _, ele := range req.CategoryIds {
		cid, err := lib.ToUUID(ele)
		if err != nil {
			return nil, xerrors.Errorf("%w", err)
		}
		err = s.Model.CheckCateIdExistedAdvertise(ctx, &advertise_package.Search{
			AdvertisePackage: advertise_package.AdvertisePackage{
				BaseModel: database.BaseModel{
					ID: pid,
				},
			},
		}, cid)
		if err == nil {
			return nil, e.ErrCategoryExisted
		}
	}

	ucids := make([]uuid.UUID, 0)

	for _, s := range req.CategoryIds {
		ucid, err := lib.ToUUID(s)
		if err != nil {
			return nil, xerrors.Errorf("%w", err)
		}
		ucids = append(ucids, ucid)
	}

	pName := lib.StandardizeSpaces(req.Name)

	err = s.Model.EditAdvertisePackage(ctx, &advertise_package.Search{
		AdvertisePackage: advertise_package.AdvertisePackage{
			BaseModel: database.BaseModel{ID: pid},
		},
	}, &advertise_package.AdvertisePackage{
		Name:        &pName,
		Price:       &req.Price,
		Categories:  ucids,
		BannerUrl:   utils.SafeStrPtr(req.BannerUrl),
		Description: &req.Description,
	})

	if err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err, ctx)
		return nil, err
	}

	return &pb.AdminAdvertiseManagementPutResponse_Data{}, nil
}

func (s *AdminService) DeleteAdvertisePackage(ctx context.Context, req *pb.AdminAdvertiseManagementDeletePostRequest) (*pb.AdminAdvertiseManagementDeletePostResponse_Data, error) {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(s.DeleteAdvertisePackage))
	defer span.End()

	if f, ok := validate.RequiredFields(req, "Id", "XUserId"); !ok {
		err := xerrors.Errorf("%w", e.ErrMissingField(f))
		lib.RecordError(span, err, ctx)
		return nil, err
	}

	err := s.Model.DeleteAdvertisePackage(ctx, req.Id)
	if err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err, ctx)
		return nil, err
	}

	return &pb.AdminAdvertiseManagementDeletePostResponse_Data{}, nil
}
