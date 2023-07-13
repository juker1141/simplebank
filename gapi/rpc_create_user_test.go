package gapi

import (
	"context"
	"database/sql"
	"fmt"
	"reflect"
	"testing"

	"github.com/golang/mock/gomock"
	mockdb "github.com/juker1141/simplebank/db/mock"
	db "github.com/juker1141/simplebank/db/sqlc"
	"github.com/juker1141/simplebank/pb"
	"github.com/juker1141/simplebank/util"
	"github.com/juker1141/simplebank/worker"
	mockwk "github.com/juker1141/simplebank/worker/mock"
	"github.com/lib/pq"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)


type eqCreateUserTxParamsMatcher struct {
	arg db.CreateUserTxParams
	password string
	user db.User
}

func (expected eqCreateUserTxParamsMatcher) Matches(x interface{}) bool {
	actualArg, ok := x.(db.CreateUserTxParams)
	if !ok {
		return false
	}

	err := util.CheckPassword(expected.password, actualArg.HashedPassword)
	if err != nil {
		return false
	}

	expected.arg.HashedPassword = actualArg.HashedPassword
	if !reflect.DeepEqual(expected.arg.CreateUserParams, actualArg.CreateUserParams) {
		return false
	}

	// call the AfterCreate function here
	err = actualArg.AfterCreate(expected.user)

	return err == nil
}

func (e eqCreateUserTxParamsMatcher) String() string {
	return fmt.Sprintf("matches arg %v and password %v", e.arg, e.password)
}

func eqCreateUserTxParams(arg db.CreateUserTxParams, password string, user db.User) gomock.Matcher {
	return eqCreateUserTxParamsMatcher{arg, password, user}
}

func randomUser(t *testing.T) (user db.User, password string) {
	password = util.RandomString(6)
	hashedPassword, err := util.HashPassword(password)
	require.NoError(t, err)

	user = db.User{
		Username: util.RandomOwner(),
		HashedPassword: hashedPassword,
		FullName: util.RandomOwner(),
		Email: util.RandomEmail(),
	}
	return
} 

func TestCreateUser(t *testing.T) {
	user, password := randomUser(t)

	testCases := []struct{
		name string
		req  *pb.CreateUserRequest
		buildStubs func(store *mockdb.MockStore, taskDistributor *mockwk.MockTaskDistribtor)
		checkResponse func(t *testing.T, res *pb.CreateUserResponse,err error)
	}{ 
		{
			name: "OK",
			req: &pb.CreateUserRequest{
				Username: user.Username,
				Password: password,
				FullName: user.FullName,
				Email: user.Email,
			},
			buildStubs: func(store *mockdb.MockStore, taskDistributor *mockwk.MockTaskDistribtor) {
				arg := db.CreateUserTxParams{
					CreateUserParams: db.CreateUserParams{
						Username: user.Username,
						FullName: user.FullName,
						Email: user.Email,
					},
				}

				store.EXPECT().
					CreateUserTx(gomock.Any(), eqCreateUserTxParams(arg, password, user)).
					Times(1).
					Return(db.CreateUserTxResult{User: user}, nil)
					
				taskPayload := &worker.PayloadSendVerifyEmail{
					Username: user.Username,
				}
				
				taskDistributor.EXPECT().
					DistributeTaskSendVerifyEmail(gomock.Any(), taskPayload, gomock.Any()).
					Times(1).
					Return(nil)
			},
			checkResponse: func(t *testing.T, res *pb.CreateUserResponse,err error) {
				require.NoError(t, err)
				require.NotNil(t, res)
				createUser := res.GetUser()
				require.Equal(t, user.Username, createUser.Username)
				require.Equal(t, user.FullName, createUser.FullName)
				require.Equal(t, user.Email, createUser.Email)
			},
		},
		{
			name: "InternalError",
			req: &pb.CreateUserRequest{
				Username: user.Username,
				Password: password,
				FullName: user.FullName,
				Email: user.Email,
			},
			buildStubs: func(store *mockdb.MockStore, taskDistributor *mockwk.MockTaskDistribtor) {
				store.EXPECT().
					CreateUserTx(gomock.Any(), gomock.Any()).
					Times(1).
					Return(db.CreateUserTxResult{}, sql.ErrConnDone)
					
				taskDistributor.EXPECT().
					DistributeTaskSendVerifyEmail(gomock.Any(), gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, res *pb.CreateUserResponse,err error) {
				require.Error(t, err)
				// if the err is from grpc error type, returns ok is true
				st, ok := status.FromError(err)
				require.True(t, ok)
				require.Equal(t, codes.Internal, st.Code())
			},
		},
		{
			name: "DuplicateUsername",
			req: &pb.CreateUserRequest{
				Username: user.Username,
				Password: password,
				FullName: user.FullName,
				Email: user.Email,
			},
			buildStubs: func(store *mockdb.MockStore, taskDistributor *mockwk.MockTaskDistribtor) {
				store.EXPECT().
					CreateUserTx(gomock.Any(), gomock.Any()).
					Times(1).
					Return(db.CreateUserTxResult{}, &pq.Error{Code: "23505"})
					
				taskDistributor.EXPECT().
					DistributeTaskSendVerifyEmail(gomock.Any(), gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, res *pb.CreateUserResponse,err error) {
				require.Error(t, err)
				// if the err is from grpc error type, returns ok is true
				st, ok := status.FromError(err)
				require.True(t, ok)
				require.Equal(t, codes.AlreadyExists, st.Code())
			},
		},
		{
			name: "InvalidUsername",
			req: &pb.CreateUserRequest{
				Username: "invalid-user#1",
				Password: password,
				FullName: user.FullName,
				Email: user.Email,
			},
			buildStubs: func(store *mockdb.MockStore, taskDistributor *mockwk.MockTaskDistribtor) {
				store.EXPECT().
					CreateUserTx(gomock.Any(), gomock.Any()).
					Times(0)
				
				taskDistributor.EXPECT().
					DistributeTaskSendVerifyEmail(gomock.Any(), gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, res *pb.CreateUserResponse,err error) {
				require.Error(t, err)
				// if the err is from grpc error type, returns ok is true
				st, ok := status.FromError(err)
				require.True(t, ok)
				require.Equal(t, codes.InvalidArgument, st.Code())
			},
		},
		{
			name: "InvalidEmail",
			req: &pb.CreateUserRequest{
				Username: user.Username,
				Password: password,
				FullName: user.FullName,
				Email: "invalid-email",
			},
			buildStubs: func(store *mockdb.MockStore, taskDistributor *mockwk.MockTaskDistribtor) {
				store.EXPECT().
					CreateUserTx(gomock.Any(), gomock.Any()).
					Times(0)
				
				taskDistributor.EXPECT().
					DistributeTaskSendVerifyEmail(gomock.Any(), gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, res *pb.CreateUserResponse,err error) {
				require.Error(t, err)
				// if the err is from grpc error type, returns ok is true
				st, ok := status.FromError(err)
				require.True(t, ok)
				require.Equal(t, codes.InvalidArgument, st.Code())
			},
		},
		{
			name: "TooShortPassword",
			req: &pb.CreateUserRequest{
				Username: user.Username,
				Password: "pwd",
				FullName: user.FullName,
				Email: user.Email,
			},
			buildStubs: func(store *mockdb.MockStore, taskDistributor *mockwk.MockTaskDistribtor) {
				store.EXPECT().
					CreateUserTx(gomock.Any(), gomock.Any()).
					Times(0)
				
				taskDistributor.EXPECT().
					DistributeTaskSendVerifyEmail(gomock.Any(), gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, res *pb.CreateUserResponse,err error) {
				require.Error(t, err)
				// if the err is from grpc error type, returns ok is true
				st, ok := status.FromError(err)
				require.True(t, ok)
				require.Equal(t, codes.InvalidArgument, st.Code())
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T){
			// The controller will be locked after it checks if the CreateUserTx argument matches in the CreateUseTx.
			storeCtrl := gomock.NewController(t)
			defer storeCtrl.Finish()
			store := mockdb.NewMockStore(storeCtrl)

			taskCtrl := gomock.NewController(t)
			defer taskCtrl.Finish()
			taskDistributor := mockwk.NewMockTaskDistribtor(taskCtrl)

			tc.buildStubs(store, taskDistributor)

			server := newTestServer(t, store, taskDistributor)
			
			res, err := server.CreateUser(context.Background(), tc.req)

			tc.checkResponse(t, res, err)
		})
	}
}