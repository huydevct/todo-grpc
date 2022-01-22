package v1

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/golang/protobuf/ptypes"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/huydevct/todo-grpc/pkg/api/v1"
)

const (
	apiVersion = "v1"
)

type toDoServiceServer struct {
	db *sql.DB
}

func NewToDoServiceServer(db *sql.DB) v1.ToDoServiceServer {
	return &toDoServiceServer{db: db}
}

func (s *toDoServiceServer) checkAPI(api string) error {
	if len(api) > 0 {
		if apiVersion != api {
			return status.Errorf(codes.Unimplemented, "Unsupported API version: service implement API version '%s', but asked for '%s'", apiVersion, api)
		}
	}
	return nil
}

func (s *toDoServiceServer) connect(ctx context.Context) (*sql.Conn, error) {
	c, err := s.db.Conn(ctx)
	if err != nil {
		return nil, status.Error(codes.Unknown, "failed to Connect to Database -> "+err.Error())
	}
	return c, nil
}

// Create new todo
func (s *toDoServiceServer) Create(ctx context.Context, req *v1.CreateRequest) (*v1.CreateResponse, error) {
	// check if the API version requested by client is supported by server
	if err := s.checkAPI(req.Api); err != nil {
		return nil, err
	}

	// get SQL connection 
	c, err := s.connect(ctx)
	if err != nil {
		return nil, err
	}
	defer c.Close()

	insert_at, err := ptypes.Timestamp(req.ToDo.InsertAt)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "insert_at field has invalid format ->"+err.Error())
	}

	update_at, err := ptypes.Timestamp(req.ToDo.UpdateAt)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "update_at filed has invalid format ->"+err.Error())
	}

	res, err := c.ExecContext(ctx, "INSERT INTO ToDo(`Title`,`Description`, `InsertAt`, `UpdateAt`) VALUES(?,?,?,?)",
		req.ToDo.Title, req.ToDo.Description, insert_at, update_at)
	if err != nil {
		return nil, status.Error(codes.Unknown, "Failed to insert to ToDo ->"+err.Error())
	}

	id,err := res.LastInsertId()
	if err != nil {
		return nil, status.Error(codes.Unknown, "Failed to retrieve id for created ToDo ->"+err.Error())
	}

	return &v1.CreateResponse{
		Api: apiVersion,
		Id: id,
	}, nil
}

func (s *toDoServiceServer) Read(ctx context.Context, req *v1.ReadRequest) (*v1.ReadResponse, error){
	// check if the API version requested by client is supported by server
	if err := s.checkAPI(req.Api); err != nil {
		return nil, err
	}

	// get SQL connection 
	c, err := s.connect(ctx)
	if err != nil {
		return nil, err
	}
	defer c.Close()

	rows, err := c.QueryContext(ctx, "SELECT * FROM ToDo WHERE `Id`=?", req.Id)
	if err != nil {
		return nil, status.Error(codes.Unknown, "Failed to select from ToDo -> "+err.Error())
	}
	defer rows.Close()

	if !rows.Next() {
		if err := rows.Err(); err != nil {
			return nil, status.Error(codes.Unknown, "failed to retrieves data from ToDo->"+err.Error())
		}
		return nil, status.Error(codes.NotFound, fmt.Sprintf("ToDo with ID='%d' is not found", req.Id))
	}

	var todo v1.ToDo
	var insert_at time.Time
	var update_at time.Time
	if err := rows.Scan(&todo.Id, &todo.Title, &todo.Description, &insert_at, &update_at); err != nil {
		return nil, status.Error(codes.Unknown, "Failed to retrieve field values from ToDo row-> "+err.Error())
	}

	todo.InsertAt, err = ptypes.TimestampProto(insert_at)
	if err != nil {
		return nil, status.Error(codes.Unknown, "insert_at field has invalid format-> "+err.Error())
	}

	todo.UpdateAt, err = ptypes.TimestampProto(update_at)
	if err != nil {
		return nil, status.Error(codes.Unknown, "update_at field has invalid format-> "+err.Error())
	}

	if rows.Next() {
		return nil, status.Error(codes.Unknown, fmt.Sprintf("found multiple ToDo rows with ID='%d'", req.Id))
	}

	return &v1.ReadResponse{
		Api: apiVersion,
		ToDo: &todo,
	}, nil
}

func (s *toDoServiceServer) Update(ctx context.Context, req *v1.UpdateRequest) (*v1.UpdateResponse, error) {
	// check if the API version requested by client is supported by server
	if err := s.checkAPI(req.Api); err != nil {
		return nil, err
	}

	// get SQL connection 
	c, err := s.connect(ctx)
	if err != nil {
		return nil, err
	}
	defer c.Close()

	update_at, err := ptypes.Timestamp(req.ToDo.UpdateAt)
	if err != nil {
		return nil, status.Error(codes.Unknown, "update_at field has invalid format-> "+err.Error())
	}

	res, err := c.ExecContext(ctx, "UPDATE ToDo SET `Title`=?, `Description`=?, `UpdateAt`=? WHERE `ID`=?", 
		req.ToDo.Title, req.ToDo.Description, update_at, req.ToDo.Id)

	if err != nil {
		return nil,status.Error(codes.Unknown, "failed to update ToDo-> "+err.Error())
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return nil, status.Error(codes.Unknown, "failed to retrieve rows affected value-> "+err.Error())
	}

	if rows == 0 {
		return nil, status.Error(codes.NotFound, fmt.Sprintf("ToDo with Id='%d' is not found", req.ToDo.Id))
	}

	return &v1.UpdateResponse{
		Api: apiVersion,
		Updated: rows,
	},nil
}


func (s *toDoServiceServer) Delete(ctx context.Context, req *v1.DeleteRequest) (*v1.DeleteResponse, error) {
	// check if the API version requested by client is supported by server
	if err := s.checkAPI(req.Api); err != nil {
		return nil, err
	}

	// get SQL connection 
	c, err := s.connect(ctx)
	if err != nil {
		return nil, err
	}
	defer c.Close()

	// delete ToDo
	res, err := c.ExecContext(ctx, "DELETE FROM ToDo WHERE `Id`=?", req.Id)
	if err != nil {
		return nil, status.Error(codes.Unknown, "failed to delete ToDo-> "+err.Error())
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return nil, status.Error(codes.Unknown, "failed to retrieve rows affected value-> "+err.Error())
	}

	if rows == 0 {
		return nil, status.Error(codes.NotFound, fmt.Sprintf("ToDo with Id='%d' is not found", req.Id))
	}

	return &v1.DeleteResponse{
		Api: apiVersion,
		Deleted: rows,
	}, nil
}

func (s *toDoServiceServer) ReadAll(ctx context.Context, req *v1.ReadAllRequest) (*v1.ReadAllResponse, error) {
	// check if the API version requested by client is supported by server
	if err := s.checkAPI(req.Api); err != nil {
		return nil, err
	}

	// get SQl connect
	c, err := s.connect(ctx)
	if err != nil {
		return nil, err
	}
	defer c.Close()

	rows, err := c.QueryContext(ctx, "SELECT * FROM ToDo")
	if err != nil {
		return nil, status.Error(codes.Unknown, "failed to select from ToDo-> "+err.Error())
	}
	defer rows.Close()

	var insert_at time.Time
	var update_at time.Time
	list := []*v1.ToDo{}
	for rows.Next() {
		todo := new(v1.ToDo)
		if err := rows.Scan(&todo.Id, &todo.Title, &todo.Description, &insert_at, &update_at); err != nil {
			return nil, status.Error(codes.Unknown, "failed to retrieve field values from ToDo row-> "+err.Error())
		}
		todo.InsertAt, err = ptypes.TimestampProto(insert_at)
		if err != nil {
			return nil, status.Error(codes.Unknown, "insert_at field has invalid format-> "+err.Error())
		}

		todo.UpdateAt, err = ptypes.TimestampProto(update_at)
		if err != nil {
			return nil, status.Error(codes.Unknown, "update_at field has invalid format-> "+err.Error())
		}
		list = append(list, todo)
	}

	if err := rows.Err(); err != nil {
		return nil, status.Error(codes.Unknown, "failed to retrieve data from ToDo-> "+err.Error())
	}

	return &v1.ReadAllResponse{
		Api: apiVersion,
		ToDos: list,
	}, nil
}