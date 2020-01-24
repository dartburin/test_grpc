package server

import (
	"context"
	"errors"
	"fmt"
	"net"
	"strconv"

	bk "test_grpc/api/proto"
	mid "test_grpc/internal/api/middleware"
	pdb "test_grpc/internal/books"

	"google.golang.org/grpc"
)

// Data for handlers
type supportGRPC struct {
	db    *pdb.ParamDB
	GPort string
	GConn string
}

// New creates new server struct
func New(db *pdb.ParamDB, gport string) *supportGRPC {
	gstr := fmt.Sprintf(":%s", gport)
	return &supportGRPC{
		db:    db,
		GPort: gport,
		GConn: gstr,
	}
}

func (s *supportGRPC) Start() error {
	s.db.Log.Println("Server gRPC init.")
	lis, err := net.Listen("tcp", s.GConn)
	if err != nil {
		return err
	}

	ll := &mid.LogData{
		Log: s.db.Log,
		Name: "Server gRPC",
	}

	srv := grpc.NewServer(ll.WithServerUnaryInterceptors())
	bk.RegisterLibraryServer(srv, s)

	s.db.Log.Println("Server gRPC start.")
	err = srv.Serve(lis)
	if err != nil {
		return err
	}
	return nil
}

func toTransport(b *pdb.BookRecord) *bk.OneBook {
	return &bk.OneBook{
		Id:     b.Id,
		Author: b.Author,
		Title:  b.Title,
	}
}

func toDB(id int64, bb *bk.BookData) *pdb.BookRecord {
	return &pdb.BookRecord{
		Id:     id,
		Author: bb.Author,
		Title:  bb.Title,
	}
}

func (s *supportGRPC) GetBooks(ctx context.Context, imsg *bk.GetBooksRequest) (*bk.Books, error) {
	var bb pdb.BookRecord
	books, err := bb.SelectBooks(s.db)
	if err != nil {
		return nil, errors.New("Error get books")
	}

	bs := bk.Books{}
	bs.Books = make([]*bk.OneBook, 0, 50)
	for i := range *books {
		b := toTransport(&(*books)[i])
		bs.Books = append(bs.Books, b)
	}

	return &bs, nil
}

func (s *supportGRPC) GetBook(ctx context.Context, imsg *bk.GetBookRequest) (*bk.OneBook, error) {
	id, err := strconv.Atoi(imsg.BookId)
	if err != nil {
		return nil, errors.New("Error get book (bad id)")
	}
	b := &pdb.BookRecord{}
	b.Id = int64(id)

	bk, err := b.SelectBook(s.db)
	if err != nil {
		return nil, errors.New("Error get book")
	}

	bb := toTransport(bk)
	return bb, nil
}

// Handler for post request
func (s *supportGRPC) PostBook(ctx context.Context, imsg *bk.PostBookRequest) (*bk.OneBook, error) {
	b := toDB(0, imsg.Msg)

	id, err := b.InsertBook(s.db)
	if err != nil {
		return nil, errors.New("Error post (insert) book")
	}

	b.Id = id
	// Maybe select not need
	_, err = b.SelectBook(s.db)
	if err != nil {
		return nil, errors.New("Error post (select) book")
	}

	bb := toTransport(b)
	return bb, nil
}

// Handler for delete request
func (s *supportGRPC) DeleteBook(ctx context.Context, imsg *bk.DeleteBookRequest) (*bk.Result, error) {
	id, err := strconv.Atoi(imsg.BookId)
	if err != nil {
		return nil, errors.New("Error delete book (bad id)")
	}
	b := &pdb.BookRecord{}
	b.Id = int64(id)

	err = b.DeleteBook(s.db)
	if err != nil {
		return nil, errors.New("Error delete book")
	}

	bb := bk.Result{}
	bb.Rez = fmt.Sprintf("Delete book with id = %v", b.Id)
	return &bb, nil
}

// Handler for put request
func (s *supportGRPC) UpdateBook(ctx context.Context, imsg *bk.UpdateBookRequest) (*bk.OneBook, error) {
	id, err := strconv.Atoi(imsg.BookId)
	if err != nil {
		return nil, errors.New("Error update book (bad id)")
	}

	b := toDB(int64(id), imsg.Msg)
	if b.Author == "" || b.Title == "" {
		return nil, errors.New("Error some parameters not set for PUT request")
	}

	err = b.UpdateBook(s.db)
	if err != nil {
		return nil, errors.New("Error update book")
	}

	bb := toTransport(b)
	return bb, nil
}

// Handler for patch request
func (s *supportGRPC) PatchBook(ctx context.Context, imsg *bk.UpdateBookRequest) (*bk.OneBook, error) {
	id, err := strconv.Atoi(imsg.BookId)
	if err != nil {
		return nil, errors.New("Error patch book (bad id)")
	}

	b := toDB(int64(id), imsg.Msg)
	err = b.UpdateBook(s.db)
	if err != nil {
		return nil, errors.New("Error patch book")
	}

	bb := toTransport(b)
	return bb, nil
}
