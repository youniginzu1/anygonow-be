package api

import (
	"context"

	"github.com/aqaurius6666/apiservice/src/internal/db/contact"
	"github.com/aqaurius6666/apiservice/src/internal/lib"
	"github.com/aqaurius6666/apiservice/src/internal/lib/validate"
	"github.com/aqaurius6666/apiservice/src/internal/model"
	"github.com/aqaurius6666/apiservice/src/internal/var/c"
	"github.com/aqaurius6666/apiservice/src/internal/var/e"
	"github.com/aqaurius6666/apiservice/src/pb"
	"github.com/aqaurius6666/go-utils/utils"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel"
	"golang.org/x/xerrors"
)

type ContactService struct {
	Model model.Server
}

func (s ContactService) GetById(ctx context.Context, req *pb.ContactGetRequest) (*pb.ContactGetResponse_Data, error) {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(s.GetById))
	defer span.End()

	if f, ok := validate.RequiredFields(req, "Id"); !ok {
		err := xerrors.Errorf("%w", e.ErrMissingField(f))
		lib.RecordError(span, err, ctx)
		return nil, err
	}
	ct, err := s.Model.GetContactById(ctx, req.Id)
	if err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err, ctx)
		return nil, err
	}

	return &pb.ContactGetResponse_Data{
		Contact: s.Model.ConvertContactToProto(ct),
	}, nil
}

func (s ContactService) Update(ctx context.Context, req *pb.ContactPutRequest) (*pb.ContactPutResponse_Data, error) {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(s.Update))
	defer span.End()

	if f, ok := validate.RequiredFields(req, "Id", "XUserId"); !ok {
		err := xerrors.Errorf("%w", e.ErrMissingField(f))
		lib.RecordError(span, err, ctx)
		return nil, err
	}
	if req.Id != req.XUserId {
		err := xerrors.Errorf("%w", e.ErrNoPermission)
		lib.RecordError(span, err, ctx)
		return nil, err
	}
	_, err := s.Model.GetStateById(ctx, req.StateId)
	if err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err, ctx)
		return nil, err
	}
	b := &contact.Contact{
		Zipcode:  utils.StrPtr(req.Zipcode),
		Address1: utils.StrPtr(req.Address1),
		Address2: utils.StrPtr(req.Address2),
		StateId:  uuid.MustParse(req.StateId),
		City:     utils.StrPtr(req.City),
	}
	if err := validate.Validate(b); err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err, ctx)
		return nil, err
	}
	if err := s.Model.UpdateContact(ctx, req.Id, b); err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err, ctx)
		return nil, err
	}

	return &pb.ContactPutResponse_Data{
		Contact: s.Model.ConvertContactToProto(b),
	}, nil
}

func (s ContactService) ListStates(ctx context.Context, req *pb.StatesGetRequest) (*pb.StatesGetResponse_Data, error) {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(s.ListStates))
	defer span.End()

	ct, err := s.Model.ListStates(ctx)
	if err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err, ctx)
		return nil, err
	}

	return &pb.StatesGetResponse_Data{
		States: s.Model.ConvertStatesToProto(ct),
	}, nil
}


func (s ContactService) GetStateById(ctx context.Context, req *pb.ContactGetRequest) (*pb.ContactGetResponse_Data, error) {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(s.GetStateById))
	defer span.End()

	if f, ok := validate.RequiredFields(req, "Id"); !ok {
		err := xerrors.Errorf("%w", e.ErrMissingField(f))
		lib.RecordError(span, err, ctx)
		return nil, err
	}
	ct, err := s.Model.GetStateByContactId(ctx, req.Id)
	if err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err, ctx)
		return nil, err
	}

	return &pb.ContactGetResponse_Data{
		Contact: s.Model.ConvertContactToProto(ct),
	}, nil
}