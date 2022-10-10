package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/aqaurius6666/authservice/src/internal/lib"
	"github.com/aqaurius6666/authservice/src/internal/var/c"
	"github.com/aqaurius6666/authservice/src/internal/var/e"
	"github.com/aqaurius6666/authservice/src/pb/authpb"
	"github.com/aqaurius6666/go-utils/cryptography"
	"github.com/aqaurius6666/go-utils/utils"
	"go.opentelemetry.io/otel"
	"golang.org/x/xerrors"
)

type HeaderStruct struct {
	Signature       string `json:"signature"`
	CertificateInfo struct {
		ID        string `json:"id"`
		Timestamp int64  `json:"timestamp"`
		Exp       int64  `json:"exp"`
	} `json:"certificateInfo"`
	PublicKey string `json:"publicKey"`
}

type BodyStruct struct {
	Data      json.RawMessage `json:"data"`
	Signature string          `json:"_signature"`
}

type BaseBodyStruct struct {
	ActionType string `json:"_actionType"`
	Timestamp  int64  `json:"_timestamp"`
}

var ADD_EXPIRE_DURATION time.Duration = 7 * 24 * time.Hour

func (s *ApiServer) CheckAuth(ctx context.Context, req *authpb.CheckAuthRequest) (*authpb.CheckAuthResponse, error) {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(s.CheckAuth))
	defer span.End()
	header := HeaderStruct{}
	if err := json.Unmarshal(req.Header, &header); err != nil {
		err = xerrors.Errorf("%w", e.ErrAuthParseModelFail)
		lib.RecordError(span, err)
		panic(err)
	}
	if header.CertificateInfo.Timestamp+header.CertificateInfo.Exp < time.Now().UnixMilli() {
		err := xerrors.Errorf("%w", e.ErrAuthExpired)
		lib.RecordError(span, err)
		panic(err)
	}
	bHeaderSig, err := cryptography.Base64ToBytes(header.Signature)
	if err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err)
		panic(err)
	}
	bPublicKey, err := cryptography.Base64ToBytes(header.PublicKey)
	if err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err)
		panic(err)
	}
	if ok, err := cryptography.VerifySignature(header.CertificateInfo, bHeaderSig, bPublicKey); err != nil || !ok {
		err = xerrors.Errorf("%w", e.ErrAuthVerifySignatureFail)
		lib.RecordError(span, err)
		panic(err)
	}

	if req.Method != http.MethodGet {
		body := BodyStruct{}
		if err := json.Unmarshal(req.Body, &body); err != nil {
			err = xerrors.Errorf("%w", e.ErrAuthParseModelFail)
			lib.RecordError(span, err)
			panic(err)
		}
		bData, err := body.Data.MarshalJSON()
		if err != nil {
			err = xerrors.Errorf("%w", err)
			lib.RecordError(span, err)
			panic(err)
		}
		baseBody := BaseBodyStruct{}
		err = json.Unmarshal(bData, &baseBody)
		if err != nil {
			err = xerrors.Errorf("%w", e.ErrAuthParseModelFail)
			lib.RecordError(span, err)
			panic(err)
		}
		if baseBody.ActionType == "" {
			err = xerrors.Errorf("%w", e.ErrInvalidActionType)
			lib.RecordError(span, err)
			panic(err)
		}
		if baseBody.Timestamp < time.Now().Add(-ADD_EXPIRE_DURATION).UnixMilli() {
			fmt.Println(baseBody.Timestamp, time.Now().UnixMilli())
			err = xerrors.Errorf("%w", e.ErrInvalidTimestamp)
			lib.RecordError(span, err)
			panic(err)
		}
		bBodySig, err := cryptography.Base64ToBytes(body.Signature)
		if err != nil {
			err = xerrors.Errorf("%w", err)
			lib.RecordError(span, err)
			panic(err)
		}
		if ok, err := cryptography.VerifySignature(body.Data, bBodySig, bPublicKey); err != nil || !ok {
			err = xerrors.Errorf("%w", e.ErrAuthVerifySignatureFail)
			lib.RecordError(span, err)
			panic(err)
		}
	}

	usr, err := s.Model.GetUserByIdPk(ctx, header.CertificateInfo.ID, header.PublicKey)
	if err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err)
		panic(err)
	}
	if !*usr.IsActive {
		err = xerrors.Errorf("%w", e.ErrUserInactive)
		lib.RecordError(span, err)
		panic(err)
	}
	return &authpb.CheckAuthResponse{
		Id:   usr.ID.String(),
		Role: c.ROLE(utils.IntVal(usr.Role.Code)),
	}, nil
}
