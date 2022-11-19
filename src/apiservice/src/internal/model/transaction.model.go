package model

import (
	"bytes"
	"context"
	"io"

	"github.com/aqaurius6666/apiservice/src/internal/db/transaction"
	"github.com/aqaurius6666/apiservice/src/internal/lib"
	"github.com/aqaurius6666/apiservice/src/internal/var/c"
	"github.com/aqaurius6666/apiservice/src/pb"
	"github.com/aqaurius6666/go-utils/utils"
	"github.com/google/uuid"
	"github.com/xuri/excelize/v2"
	"go.opentelemetry.io/otel"
	"golang.org/x/xerrors"
)

var (
	_ TransactionModel = (*ServerModel)(nil)
)

type TransactionModel interface {
	TotalTransaction(ctx context.Context, search *transaction.Search) (*int64, error)
	TotalFee(ctx context.Context, search *transaction.Search) (*int64, error)
	GetTransactions(ctx context.Context, search *transaction.Search) ([]*transaction.Transaction, error)
	InsertTransaction(context context.Context, orderId interface{}, businessId interface{}, fee float32, isFree bool) error
	UpdateTransaction(ctx context.Context, search *transaction.Search, value *transaction.Transaction) error
	ExportTransactions(ctx context.Context, transactions []*transaction.Transaction, totalFee *float32) (io.Reader, int64, error)

	ConvertTransactionToProtos(u []*transaction.Transaction) []*pb.Transaction
	ConvertTransactionToProto(u *transaction.Transaction) *pb.Transaction
}

func (s *ServerModel) ConvertTransactionToProto(u *transaction.Transaction) *pb.Transaction {
	upb := new(pb.Transaction)
	if u.StartDate != nil {
		upb.StartDate = *u.StartDate
	}
	if u.EndDate != nil {
		upb.EndDate = *u.EndDate
	}
	if u.ServiceName != nil {
		upb.ServiceName = *u.ServiceName
	}
	if u.CustomerZipcode != nil {
		upb.Zipcode = *u.CustomerZipcode
	}
	if u.Fee != nil {
		upb.Fee = *lib.CentToUsd(u.Fee)
	}
	if u.Status != nil {
		upb.Status = c.ORDER_STATUS(*u.Status)
	}
	if u.CustomerId != uuid.Nil {
		upb.Id = u.CustomerId.String()
	}
	if u.CustomerAvatar != nil {
		upb.Image = *u.CustomerAvatar
	}

	return upb
}

func (s *ServerModel) ConvertTransactionToProtos(u []*transaction.Transaction) []*pb.Transaction {
	arr := make([]*pb.Transaction, 0)
	for _, a := range u {
		arr = append(arr, s.ConvertTransactionToProto(a))
	}
	return arr
}

func (s *ServerModel) GetTransactions(ctx context.Context, search *transaction.Search) ([]*transaction.Transaction, error) {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(s.GetTransactions))
	defer span.End()

	r, err := s.Repo.ListTransactions(ctx, search)
	if err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err, ctx)
		return nil, err
	}

	return r, nil
}

func (s *ServerModel) TotalTransaction(ctx context.Context, search *transaction.Search) (*int64, error) {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(s.TotalTransaction))
	defer span.End()

	r, err := s.Repo.TotalTransaction(ctx, search)
	if err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err, ctx)
		return nil, err
	}
	return r, nil
}

func (s *ServerModel) TotalFee(ctx context.Context, search *transaction.Search) (*int64, error) {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(s.TotalFee))
	defer span.End()

	r, err := s.Repo.TotalFee(ctx, search)
	if err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err, ctx)
		return nil, err
	}
	return r, nil
}

func (s *ServerModel) InsertTransaction(ctx context.Context, orderId interface{}, businessId interface{}, fee float32, isFree bool) error {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(s.TotalFee))
	defer span.End()
	oid, err := lib.ToUUID(orderId)
	if err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err, ctx)
		return err
	}
	bid, err := lib.ToUUID(businessId)
	if err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err, ctx)
		return err
	}
	err = s.Repo.InsertTransaction(ctx, &transaction.Transaction{
		BusinessId: bid,
		OrderId:    oid,
		Fee:        lib.UsdFloatToCent(&fee),
		IsFree:     &isFree,
	})
	if err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err, ctx)
		return err
	}

	return nil
}

func (s *ServerModel) UpdateTransaction(ctx context.Context, search *transaction.Search, value *transaction.Transaction) error {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(s.UpdateTransaction))
	defer span.End()

	err := s.Repo.UpdateTransaction(ctx, search, value)
	if err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err, ctx)
		return err
	}
	return nil
}

func (s *ServerModel) ExportTransactions(ctx context.Context, transactions []*transaction.Transaction, totalFee *float32) (io.Reader, int64, error) {
	ctx, span := otel.GetTracerProvider().Tracer(c.SERVICE_NAME).Start(ctx, lib.GetFunctionName(s.ExportTransactions))
	defer span.End()

	f := excelize.NewFile()

	f.SetCellValue("Sheet1", lib.GenerateCellExcelHeader("B"), "Service requested name")
	f.SetCellValue("Sheet1", lib.GenerateCellExcelHeader("C"), "Requested time")
	f.SetCellValue("Sheet1", lib.GenerateCellExcelHeader("D"), "Expiry time")
	f.SetCellValue("Sheet1", lib.GenerateCellExcelHeader("E"), "Zipcode")
	f.SetCellValue("Sheet1", lib.GenerateCellExcelHeader("F"), "Deal fee")

	for i, s := range transactions {
		f.SetCellValue("Sheet1", lib.GenerateCellExcelData("B", i), utils.StrVal(s.ServiceName))
		f.SetCellValue("Sheet1", lib.GenerateCellExcelData("C", i), lib.FormatMillisecondsToDate(utils.Int64Val(s.StartDate)))
		f.SetCellValue("Sheet1", lib.GenerateCellExcelData("D", i), lib.FormatMillisecondsToDate(utils.Int64Val(s.EndDate)))
		f.SetCellValue("Sheet1", lib.GenerateCellExcelData("E", i), utils.StrVal(s.CustomerZipcode))
		f.SetCellValue("Sheet1", lib.GenerateCellExcelData("F", i), utils.Float32Val(lib.CentToUsd(s.Fee)))
	}

	length := len(transactions)
	f.SetCellValue("Sheet1", lib.GenerateCellExcelData("F", length), utils.Float32Val(totalFee))

	buffer := bytes.Buffer{}
	_, err := f.WriteTo(&buffer)
	len := buffer.Len()
	if err != nil {
		err = xerrors.Errorf("%w", err)
		lib.RecordError(span, err, ctx)
		return nil, 0, err
	}
	return &buffer, int64(len), nil
}
