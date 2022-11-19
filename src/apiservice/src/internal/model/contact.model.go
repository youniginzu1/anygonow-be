package model

import (
	"context"

	"github.com/aqaurius6666/apiservice/src/internal/db/contact"
	"github.com/aqaurius6666/apiservice/src/internal/db/state"
	"github.com/aqaurius6666/apiservice/src/internal/lib"
	"github.com/aqaurius6666/apiservice/src/internal/var/c"
	"github.com/aqaurius6666/apiservice/src/pb"
	"github.com/aqaurius6666/go-utils/database"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel"
	"golang.org/x/xerrors"
)

var (
	_ ContactModel = (*ServerModel)(nil)
)

type ContactModel interface {
	GetContactById(ctx context.Context, id interface{}) (*contact.Contact, error)
	GetStateByContactId(ctx context.Context, id interface{}) (*contact.Contact, error)
	ConvertContactToProto(*contact.Contact) *pb.Contact
	GetStateById(ctx context.Context, id interface{}) (*state.State, error)
	ListStates(ctx context.Context) ([]*state.State, error)
	UpdateContact(ctx context.Context, id interface{}, b *contact.Contact) error

	ConvertStateToProto(*state.State) *pb.State
	ConvertStatesToProto([]*state.State) []*pb.State
}

func (s *ServerModel) GetStateById(ctx context.Context, id interface{}) (*state.State, error) {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(s.GetStateById))
	defer span.End()

	uid, err := lib.ToUUID(id)
	if err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err, ctx)
		return nil, err
	}
	ct, err := s.Repo.SelectState(ctx, &state.Search{
		State: state.State{BaseModel: database.BaseModel{ID: uid}},
	})
	if err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err, ctx)
		return nil, err
	}
	return ct, nil
}

func (s *ServerModel) ConvertStateToProto(u *state.State) *pb.State {
	return &pb.State{
		Id:   u.ID.String(),
		Name: *u.Name,
	}
}

func (s *ServerModel) ConvertStatesToProto(u []*state.State) []*pb.State {
	arr := make([]*pb.State, 0)
	for _, a := range u {
		arr = append(arr, s.ConvertStateToProto(a))
	}
	return arr
}
func (s *ServerModel) ListStates(ctx context.Context) ([]*state.State, error) {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(s.ListStates))
	defer span.End()

	ct, err := s.Repo.ListStates(ctx, &state.Search{})
	if err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err, ctx)
		return nil, err
	}
	return ct, nil
}

func (s *ServerModel) GetContactById(ctx context.Context, id interface{}) (*contact.Contact, error) {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(s.GetContactById))
	defer span.End()

	uid, err := lib.ToUUID(id)
	if err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err, ctx)
		return nil, err
	}
	ct, err := s.Repo.SelectContact(ctx, &contact.Search{
		Contact: contact.Contact{BaseModel: database.BaseModel{ID: uid}},
	})
	if err == nil {
		return ct, nil
	}
	ct, err = s.Repo.InsertContact(ctx, &contact.Contact{
		BaseModel: database.BaseModel{ID: uid},
	})
	if err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err, ctx)
		return nil, err
	}
	return ct, nil
}

func (s *ServerModel) ConvertContactToProto(u *contact.Contact) *pb.Contact {
	upb := new(pb.Contact)
	if u.ID != uuid.Nil {
		upb.Id = u.ID.String()
	}
	if u.Address1 != nil {
		upb.Address1 = *u.Address1
	}
	if u.Address2 != nil {
		upb.Address2 = *u.Address2
	}
	if u.City != nil {
		upb.City = *u.City
	}
	if u.StateId != uuid.Nil {
		stt, err := s.GetStateById(context.TODO(), u.StateId)
		if err != nil {
			upb.State = "unknown"
		} else {
			upb.State = *stt.Name
			upb.StateId = stt.ID.String()
		}
	}
	if u.Zipcode != nil {
		upb.Zipcode = *u.Zipcode
	}
	return upb
}

func (s *ServerModel) UpdateContact(ctx context.Context, id interface{}, u *contact.Contact) error {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(s.UpdateContact))
	defer span.End()

	uid, err := lib.ToUUID(id)
	if err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err, ctx)
		return err
	}
	err = s.Repo.UpdateContact(ctx, &contact.Search{
		Contact: contact.Contact{
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

func (s *ServerModel) GetStateByContactId(ctx context.Context, id interface{}) (*contact.Contact, error) {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(s.GetStateByContactId))
	defer span.End()

	uid, err := lib.ToUUID(id)
	if err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err, ctx)
		return nil, err
	}
	ct, err := s.Repo.SelectContact(ctx, &contact.Search{
		DefaultSearchModel: database.DefaultSearchModel{
			Fields: []string{`"states"."id" as "state_id"`},
		},
		Contact: contact.Contact{
			BaseModel: database.BaseModel{ID: uid},
		},
	})
	if err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err, ctx)
		return nil, err
	}
	return ct, nil
}
