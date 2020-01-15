package server

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net"
	"strconv"

	lg "github.com/sirupsen/logrus"
	bk "test_grpc/api/proto"
	pdb "test_grpc/internal/books"

	"google.golang.org/grpc"
)

// Data for handlers
type supportGRPC struct {
	base  *sql.DB
	log   lg.FieldLogger
	GPort string
	GConn string
}

// New creates new server struct
func New(b *sql.DB, gport string, log lg.FieldLogger) *supportGRPC {
	gstr := fmt.Sprintf(":%s", gport)
	return &supportGRPC{
		base:  b,
		log:   log,
		GPort: gport,
		GConn: gstr,
	}
}

func (s *supportGRPC) Start() error {
	s.log.Println("Server gRPC init.")
	lis, err := net.Listen("tcp", s.GConn)
	if err != nil {
		return err
	}
	srv := grpc.NewServer()
	bk.RegisterLibraryServer(srv, s)

	s.log.Println("Server gRPC start.")
	err = srv.Serve(lis)
	if err != nil {
		return err
	}
	return nil
}

func (s *supportGRPC) GetBooks(ctx context.Context, imsg *bk.GetBookRequest) (*bk.Books, error) {
	var bb pdb.Book
	books, err := bb.SelectBook(s.base)
	if err != nil {
		return nil, errors.New("Error get book")
	}

	bs := bk.Books{}
	bs.Books = make([]*bk.OneBook, 0, 50)
	for i := range *books {
		b := bk.OneBook{}
		b.Id = (*books)[i].Id
		b.Author = (*books)[i].Author
		b.Title = (*books)[i].Title
		bs.Books = append(bs.Books, &b)
	}

	return &bs, nil
}

// Handler for post request
func (s *supportGRPC) PostBook(ctx context.Context, imsg *bk.PostBookRequest) (*bk.OneBook, error) {

	b := &pdb.Book{}
	b.Author = imsg.Message.Author
	b.Title = imsg.Message.Title

	s.log.Printf("gRPC post a book %v.", b)

	id, err := b.InsertBook(s.base)
	if err != nil {
		return nil, errors.New("Error post (insert) book")
	}

	b.Id = id
	// Maybe select not need
	_, err = b.SelectBook(s.base)
	if err != nil {
		return nil, errors.New("Error post (select) book")
	}

	bb := bk.OneBook{}
	bb.Id = id
	bb.Author = b.Author
	bb.Title = b.Title
	s.log.Printf("gRPC post a book ret %v.", bb)

	return &bb, nil
}

// Handler for delete request
func (s *supportGRPC) DeleteBook(ctx context.Context, imsg *bk.DeleteBookRequest) (*bk.Result, error) {
	id, err := strconv.Atoi(imsg.BookId)
	if err != nil {
		return nil, errors.New("Error delete book (bad id)")
	}
	b := &pdb.Book{}
	b.Id = int64(id)

	err = b.DeleteBook(s.base)
	if err != nil {
		return nil, errors.New("Error delete book")
	}

	bb := bk.Result{}
	bb.Rez = fmt.Sprintf("Delete book with id = %v", b.Id)

	s.log.Printf("gRPC delete a book %v.", bb)
	return &bb, nil
}

// Handler for put request
func (s *supportGRPC) UpdateBook(ctx context.Context, imsg *bk.UpdateBookRequest) (*bk.OneBook, error) {
	id, err := strconv.Atoi(imsg.BookId)
	if err != nil {
		return nil, errors.New("Error delete book (bad id)")
	}
	b := &pdb.Book{}
	b.Id = int64(id)
	b.Author = imsg.Message.Author
	b.Title = imsg.Message.Title

	if b.Author == "" || b.Title == "" {
		return nil, errors.New("Error some parameters not set for PUT request")
	}

	err = b.UpdateBook(s.base)
	if err != nil {
		return nil, errors.New("Error update book")
	}

	bb := bk.OneBook{}
	bb.Id = b.Id
	bb.Author = b.Author
	bb.Title = b.Title

	s.log.Printf("gRPC updata a book %v.", bb)
	return &bb, nil
}

// Handler for patch request
func (s *supportGRPC) PathBook(ctx context.Context, imsg *bk.UpdateBookRequest) (*bk.OneBook, error) {
	id, err := strconv.Atoi(imsg.BookId)
	if err != nil {
		return nil, errors.New("Error delete book (bad id)")
	}
	b := &pdb.Book{}
	b.Id = int64(id)
	b.Author = imsg.Message.Author
	b.Title = imsg.Message.Title

	err = b.UpdateBook(s.base)
	if err != nil {
		return nil, errors.New("Error update book")
	}

	bb := bk.OneBook{}
	bb.Id = b.Id
	bb.Author = b.Author
	bb.Title = b.Title

	s.log.Printf("gRPC path a book %v.", bb)
	return &bb, nil
}
