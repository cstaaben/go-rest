package database

import (
	"context"
	"errors"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/cstaaben/go-rest/internal/database/dto"

	_ "modernc.org/sqlite"
)

// setupTestDAO initializes a clean DAO instance using a temporary directory and registers cleanup via t.Cleanup.
func setupTestDAO(t *testing.T) *DAO {
	t.Helper()
	tempDir := t.TempDir()
	t.Setenv("HOME", tempDir)
	t.Setenv("XDG_DATA_HOME", tempDir)
	_ = os.Remove("data.db")
	_ = os.Remove("development.db")

	ctx := context.Background()
	dao, err := New(ctx)
	require.NoError(t, err)

	t.Cleanup(func() {
		_ = dao.Close()
		_ = os.Remove("data.db")
		_ = os.Remove("development.db")
	})

	return dao
}

func TestNewDAO(t *testing.T) {
	tests := []struct {
		name    string
		getCtx  func() context.Context
		wantErr bool
	}{
		{
			name: "successful initialization",
			getCtx: func() context.Context {
				return context.Background()
			},
			wantErr: false,
		},
		{
			name: "canceled context failure",
			getCtx: func() context.Context {
				ctx, cancel := context.WithCancel(context.Background())
				cancel()
				return ctx
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tempDir := t.TempDir()
			t.Setenv("HOME", tempDir)
			t.Setenv("XDG_DATA_HOME", tempDir)

			ctx := tt.getCtx()
			dao, err := New(ctx)
			if tt.wantErr {
				assert.Error(t, err)
				if dao != nil {
					t.Cleanup(func() { dao.Close() })
				}
			} else {
				assert.NoError(t, err)
				require.NotNil(t, dao)
				t.Cleanup(func() { dao.Close() })
			}
		})
	}
}

func TestDAO_Close(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
	}{
		{
			name:    "close active database connection",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dao := setupTestDAO(t)

			err := dao.Close()
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestDAO_RequestListData(t *testing.T) {
	tests := []struct {
		name        string
		setup       func(ctx context.Context, dao *DAO)
		limit       int
		offset      int
		getCtx      func() context.Context
		wantCount   int
		checkResult func(t *testing.T, data []*dto.RequestListData)
		wantErr     bool
	}{
		{
			name:      "empty database returns empty list",
			setup:     func(ctx context.Context, dao *DAO) {},
			limit:     10,
			offset:    0,
			getCtx:    func() context.Context { return context.Background() },
			wantCount: 0,
			wantErr:   false,
		},
		{
			name: "unsorted requests and groups pagination",
			setup: func(ctx context.Context, dao *DAO) {
				_ = dao.InsertGroup(ctx, &dto.Group{UUID: "g1", Name: "Auth Service", Desc: "Auth endpoints"})
				_ = dao.InsertRequest(ctx, &dto.Request{
					UUID:      "r1",
					Name:      "Get Token",
					Desc:      "Fetch OAuth token",
					GroupUUID: dto.UnsortedName,
					Data:      &dto.Data{Method: "POST", URL: "https://example.com/token"},
				})
			},
			limit:     10,
			offset:    0,
			getCtx:    func() context.Context { return context.Background() },
			wantCount: 2,
			checkResult: func(t *testing.T, data []*dto.RequestListData) {
				var foundReq, foundGrp bool
				for _, item := range data {
					if item.RequestID == "r1" {
						foundReq = true
						assert.Equal(t, "Get Token", item.RequestName)
						assert.Equal(t, "POST", item.RequestMethod)
						assert.Equal(t, dto.UnsortedName, item.GroupID)
					}
					if item.GroupID == "g1" {
						foundGrp = true
						assert.Equal(t, "Auth Service", item.GroupName)
						assert.Equal(t, "", item.RequestID)
					}
				}
				assert.True(t, foundReq)
				assert.True(t, foundGrp)
			},
			wantErr: false,
		},
		{
			name:   "canceled context returns error",
			setup:  func(ctx context.Context, dao *DAO) {},
			limit:  10,
			offset: 0,
			getCtx: func() context.Context {
				ctx, cancel := context.WithCancel(context.Background())
				cancel()
				return ctx
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dao := setupTestDAO(t)

			ctx := context.Background()
			if tt.setup != nil {
				tt.setup(ctx, dao)
			}

			execCtx := tt.getCtx()
			data, err := dao.RequestListData(execCtx, tt.limit, tt.offset)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Len(t, data, tt.wantCount)
				if tt.checkResult != nil {
					tt.checkResult(t, data)
				}
			}
		})
	}
}

func TestDAO_Request(t *testing.T) {
	tests := []struct {
		name       string
		setup      func(ctx context.Context, dao *DAO)
		requestID  string
		getCtx     func() context.Context
		wantReq    *dto.Request
		wantErrIs  error
		wantErr    bool
	}{
		{
			name: "existing request with headers",
			setup: func(ctx context.Context, dao *DAO) {
				req := &dto.Request{
					UUID:      "req-1",
					Name:      "User Info",
					Desc:      "Get user profile",
					GroupUUID: dto.UnsortedName,
					Data: &dto.Data{
						URL:     "https://api.example.com/user",
						Method:  "GET",
						Proto:   "HTTP/1.1",
						Body:    "",
						Headers: map[string][]string{"Authorization": {"Bearer token"}},
					},
				}
				require.NoError(t, dao.InsertRequest(ctx, req))
			},
			requestID: "req-1",
			getCtx:    func() context.Context { return context.Background() },
			wantReq: &dto.Request{
				UUID:      "req-1",
				Name:      "User Info",
				Desc:      "Get user profile",
				GroupUUID: dto.UnsortedName,
				Data: &dto.Data{
					URL:     "https://api.example.com/user",
					Method:  "GET",
					Proto:   "HTTP/1.1",
					Body:    "",
					Headers: map[string][]string{"Authorization": {"Bearer token"}},
				},
			},
			wantErr: false,
		},
		{
			name:      "non-existent request",
			setup:     func(ctx context.Context, dao *DAO) {},
			requestID: "non-existent",
			getCtx:    func() context.Context { return context.Background() },
			wantErrIs: ErrNotFound,
			wantErr:   true,
		},
		{
			name:      "canceled context returns error",
			setup:     func(ctx context.Context, dao *DAO) {},
			requestID: "req-1",
			getCtx: func() context.Context {
				ctx, cancel := context.WithCancel(context.Background())
				cancel()
				return ctx
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dao := setupTestDAO(t)

			ctx := context.Background()
			if tt.setup != nil {
				tt.setup(ctx, dao)
			}

			execCtx := tt.getCtx()
			req, err := dao.Request(execCtx, tt.requestID)
			if tt.wantErr {
				assert.Error(t, err)
				if tt.wantErrIs != nil {
					assert.ErrorIs(t, err, tt.wantErrIs)
				}
			} else {
				assert.NoError(t, err)
				assert.True(t, tt.wantReq.Equal(req))
			}
		})
	}
}

func TestDAO_GroupRequests(t *testing.T) {
	tests := []struct {
		name       string
		setup      func(ctx context.Context, dao *DAO)
		groupID    string
		getCtx     func() context.Context
		wantCount  int
		wantErrIs  error
		wantErr    bool
	}{
		{
			name: "existing group with requests",
			setup: func(ctx context.Context, dao *DAO) {
				group := &dto.Group{UUID: "g-auth", Name: "Auth", Desc: "Auth group"}
				require.NoError(t, dao.InsertGroup(ctx, group))
				req1 := &dto.Request{
					UUID:      "r-1",
					Name:      "Login",
					GroupUUID: "g-auth",
					Data:      &dto.Data{Method: "POST", URL: "https://example.com/login"},
				}
				req2 := &dto.Request{
					UUID:      "r-2",
					Name:      "Logout",
					GroupUUID: "g-auth",
					Data:      &dto.Data{Method: "POST", URL: "https://example.com/logout"},
				}
				require.NoError(t, dao.InsertRequest(ctx, req1))
				require.NoError(t, dao.InsertRequest(ctx, req2))
			},
			groupID:   "g-auth",
			getCtx:    func() context.Context { return context.Background() },
			wantCount: 2,
			wantErr:   false,
		},
		{
			name: "empty group returns ErrNotFound",
			setup: func(ctx context.Context, dao *DAO) {
				group := &dto.Group{UUID: "g-empty", Name: "Empty", Desc: "Empty group"}
				require.NoError(t, dao.InsertGroup(ctx, group))
			},
			groupID:   "g-empty",
			getCtx:    func() context.Context { return context.Background() },
			wantErrIs: ErrNotFound,
			wantErr:   true,
		},
		{
			name:      "non-existent group returns ErrNotFound",
			setup:     func(ctx context.Context, dao *DAO) {},
			groupID:   "g-missing",
			getCtx:    func() context.Context { return context.Background() },
			wantErrIs: ErrNotFound,
			wantErr:   true,
		},
		{
			name:    "canceled context returns error",
			setup:   func(ctx context.Context, dao *DAO) {},
			groupID: "g-auth",
			getCtx: func() context.Context {
				ctx, cancel := context.WithCancel(context.Background())
				cancel()
				return ctx
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dao := setupTestDAO(t)

			ctx := context.Background()
			if tt.setup != nil {
				tt.setup(ctx, dao)
			}

			execCtx := tt.getCtx()
			reqs, err := dao.GroupRequests(execCtx, tt.groupID)
			if tt.wantErr {
				assert.Error(t, err)
				if tt.wantErrIs != nil {
					assert.ErrorIs(t, err, tt.wantErrIs)
				}
			} else {
				assert.NoError(t, err)
				assert.Len(t, reqs, tt.wantCount)
			}
		})
	}
}

func TestDAO_Groups(t *testing.T) {
	tests := []struct {
		name       string
		setup      func(ctx context.Context, dao *DAO)
		getCtx     func() context.Context
		wantCount  int
		wantUUIDs  []string
		wantErr    bool
	}{
		{
			name:       "initial database has default unsorted group",
			setup:      func(ctx context.Context, dao *DAO) {},
			getCtx:     func() context.Context { return context.Background() },
			wantCount:  1,
			wantUUIDs:  []string{dto.UnsortedName},
			wantErr:    false,
		},
		{
			name: "returns all groups ordered by UUID",
			setup: func(ctx context.Context, dao *DAO) {
				require.NoError(t, dao.InsertGroup(ctx, &dto.Group{UUID: "g-beta", Name: "Beta", Desc: "B"}))
				require.NoError(t, dao.InsertGroup(ctx, &dto.Group{UUID: "g-alpha", Name: "Alpha", Desc: "A"}))
			},
			getCtx:    func() context.Context { return context.Background() },
			wantCount: 3,
			wantUUIDs: []string{"g-alpha", "g-beta", dto.UnsortedName},
			wantErr:   false,
		},
		{
			name:   "canceled context returns error",
			setup:  func(ctx context.Context, dao *DAO) {},
			getCtx: func() context.Context {
				ctx, cancel := context.WithCancel(context.Background())
				cancel()
				return ctx
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dao := setupTestDAO(t)

			ctx := context.Background()
			if tt.setup != nil {
				tt.setup(ctx, dao)
			}

			execCtx := tt.getCtx()
			groups, err := dao.Groups(execCtx)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Len(t, groups, tt.wantCount)
				uuids := make([]string, len(groups))
				for i, g := range groups {
					uuids[i] = g.UUID
				}
				assert.Equal(t, tt.wantUUIDs, uuids)
			}
		})
	}
}

func TestDAO_InsertGroup(t *testing.T) {
	tests := []struct {
		name    string
		setup   func(ctx context.Context, dao *DAO)
		group   *dto.Group
		getCtx  func() context.Context
		wantErr bool
	}{
		{
			name:    "insert valid group",
			setup:   func(ctx context.Context, dao *DAO) {},
			group:   &dto.Group{UUID: "g-10", Name: "Payment Service", Desc: "Payments"},
			getCtx:  func() context.Context { return context.Background() },
			wantErr: false,
		},
		{
			name: "insert group with duplicate UUID",
			setup: func(ctx context.Context, dao *DAO) {
				require.NoError(t, dao.InsertGroup(ctx, &dto.Group{UUID: "g-dup", Name: "Original", Desc: "First"}))
			},
			group:   &dto.Group{UUID: "g-dup", Name: "Duplicate", Desc: "Second"},
			getCtx:  func() context.Context { return context.Background() },
			wantErr: true,
		},
		{
			name:    "canceled context returns error",
			setup:   func(ctx context.Context, dao *DAO) {},
			group:   &dto.Group{UUID: "g-ctx", Name: "Ctx Test"},
			getCtx:  func() context.Context {
				ctx, cancel := context.WithCancel(context.Background())
				cancel()
				return ctx
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dao := setupTestDAO(t)

			ctx := context.Background()
			if tt.setup != nil {
				tt.setup(ctx, dao)
			}

			execCtx := tt.getCtx()
			err := dao.InsertGroup(execCtx, tt.group)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				fetchedGroups, err := dao.Groups(ctx)
				require.NoError(t, err)
				var found bool
				for _, g := range fetchedGroups {
					if g.UUID == tt.group.UUID {
						found = true
						assert.Equal(t, tt.group.Name, g.Name)
						assert.Equal(t, tt.group.Desc, g.Desc)
					}
				}
				assert.True(t, found)
			}
		})
	}
}

func TestDAO_UpdateGroup(t *testing.T) {
	tests := []struct {
		name        string
		setup       func(ctx context.Context, dao *DAO)
		group       *dto.Group
		getCtx      func() context.Context
		checkUpdate func(t *testing.T, dao *DAO)
		wantErr     bool
	}{
		{
			name: "update existing group",
			setup: func(ctx context.Context, dao *DAO) {
				require.NoError(t, dao.InsertGroup(ctx, &dto.Group{UUID: "g-mod", Name: "Old Name", Desc: "Old Desc"}))
			},
			group:  &dto.Group{UUID: "g-mod", Name: "New Name", Desc: "New Desc"},
			getCtx: func() context.Context { return context.Background() },
			checkUpdate: func(t *testing.T, dao *DAO) {
				groups, err := dao.Groups(context.Background())
				require.NoError(t, err)
				var found bool
				for _, g := range groups {
					if g.UUID == "g-mod" {
						found = true
						assert.Equal(t, "New Name", g.Name)
						assert.Equal(t, "New Desc", g.Desc)
					}
				}
				assert.True(t, found)
			},
			wantErr: false,
		},
		{
			name: "exec failure returns error",
			setup: func(ctx context.Context, dao *DAO) {
				require.NoError(t, dao.InsertGroup(ctx, &dto.Group{UUID: "g-fail", Name: "Name", Desc: "Desc"}))
			},
			group: &dto.Group{UUID: "g-fail", Name: "Updated Name", Desc: "Updated Desc"},
			getCtx: func() context.Context {
				ctx, cancel := context.WithCancel(context.Background())
				cancel()
				return ctx
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dao := setupTestDAO(t)

			ctx := context.Background()
			if tt.setup != nil {
				tt.setup(ctx, dao)
			}

			execCtx := tt.getCtx()
			err := dao.UpdateGroup(execCtx, tt.group)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				if tt.checkUpdate != nil {
					tt.checkUpdate(t, dao)
				}
			}
		})
	}
}

func TestDAO_InsertRequest(t *testing.T) {
	tests := []struct {
		name    string
		setup   func(ctx context.Context, dao *DAO)
		req     *dto.Request
		getCtx  func() context.Context
		wantErr bool
	}{
		{
			name:  "insert valid request into unsorted group",
			setup: func(ctx context.Context, dao *DAO) {},
			req: &dto.Request{
				UUID:      "r-ins-1",
				Name:      "Ping",
				Desc:      "Health check",
				GroupUUID: dto.UnsortedName,
				Data: &dto.Data{
					Method:  "GET",
					URL:     "https://example.com/ping",
					Headers: map[string][]string{"Accept": {"application/json"}},
				},
			},
			getCtx:  func() context.Context { return context.Background() },
			wantErr: false,
		},
		{
			name: "insert request into non-existent group fails foreign key constraint",
			setup: func(ctx context.Context, dao *DAO) {},
			req: &dto.Request{
				UUID:      "r-ins-2",
				Name:      "Orphan",
				GroupUUID: "non-existent-group",
				Data:      &dto.Data{Method: "GET", URL: "https://example.com"},
			},
			getCtx:  func() context.Context { return context.Background() },
			wantErr: true,
		},
		{
			name: "insert request with duplicate UUID fails unique constraint",
			setup: func(ctx context.Context, dao *DAO) {
				req := &dto.Request{
					UUID:      "r-dup",
					Name:      "Original Request",
					GroupUUID: dto.UnsortedName,
					Data:      &dto.Data{Method: "GET", URL: "https://example.com"},
				}
				require.NoError(t, dao.InsertRequest(ctx, req))
			},
			req: &dto.Request{
				UUID:      "r-dup",
				Name:      "Duplicate Request",
				GroupUUID: dto.UnsortedName,
				Data:      &dto.Data{Method: "POST", URL: "https://example.com"},
			},
			getCtx:  func() context.Context { return context.Background() },
			wantErr: true,
		},
		{
			name:  "canceled context returns error",
			setup: func(ctx context.Context, dao *DAO) {},
			req: &dto.Request{
				UUID:      "r-ctx",
				Name:      "Cancel Test",
				GroupUUID: dto.UnsortedName,
				Data:      &dto.Data{Method: "GET", URL: "https://example.com"},
			},
			getCtx: func() context.Context {
				ctx, cancel := context.WithCancel(context.Background())
				cancel()
				return ctx
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dao := setupTestDAO(t)

			ctx := context.Background()
			if tt.setup != nil {
				tt.setup(ctx, dao)
			}

			execCtx := tt.getCtx()
			err := dao.InsertRequest(execCtx, tt.req)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				fetched, err := dao.Request(ctx, tt.req.UUID)
				require.NoError(t, err)
				assert.True(t, tt.req.Equal(fetched))
			}
		})
	}
}

func TestDAO_UpdateRequest(t *testing.T) {
	tests := []struct {
		name        string
		setup       func(ctx context.Context, dao *DAO)
		req         *dto.Request
		getCtx      func() context.Context
		checkUpdate func(t *testing.T, dao *DAO)
		wantErr     bool
	}{
		{
			name: "update existing request",
			setup: func(ctx context.Context, dao *DAO) {
				initial := &dto.Request{
					UUID:      "r-upd-1",
					Name:      "Initial Name",
					Desc:      "Initial Desc",
					GroupUUID: dto.UnsortedName,
					Data: &dto.Data{
						Method:  "GET",
						URL:     "https://example.com/v1",
						Proto:   "HTTP/1.1",
						Body:    "",
						Headers: map[string][]string{"Header1": {"Val1"}},
					},
				}
				require.NoError(t, dao.InsertRequest(ctx, initial))
			},
			req: &dto.Request{
				UUID:      "r-upd-1",
				Name:      "Updated Name",
				Desc:      "Updated Desc",
				GroupUUID: dto.UnsortedName,
				Data: &dto.Data{
					Method:  "POST",
					URL:     "https://example.com/v2",
					Proto:   "HTTP/2.0",
					Body:    `{"updated": true}`,
					Headers: map[string][]string{"Header2": {"Val2"}},
				},
			},
			getCtx: func() context.Context { return context.Background() },
			checkUpdate: func(t *testing.T, dao *DAO) {
				fetched, err := dao.Request(context.Background(), "r-upd-1")
				require.NoError(t, err)
				assert.Equal(t, "Updated Name", fetched.Name)
				assert.Equal(t, "Updated Desc", fetched.Desc)
				assert.Equal(t, "POST", fetched.Data.Method)
				assert.Equal(t, "https://example.com/v2", fetched.Data.URL)
				assert.Equal(t, "HTTP/2.0", fetched.Data.Proto)
				assert.Equal(t, `{"updated": true}`, fetched.Data.Body)
			},
			wantErr: false,
		},
		{
			name:  "canceled context returns error",
			setup: func(ctx context.Context, dao *DAO) {},
			req: &dto.Request{
				UUID:      "r-upd-2",
				Name:      "Cancel Update",
				GroupUUID: dto.UnsortedName,
				Data:      &dto.Data{Method: "GET", URL: "https://example.com"},
			},
			getCtx: func() context.Context {
				ctx, cancel := context.WithCancel(context.Background())
				cancel()
				return ctx
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dao := setupTestDAO(t)

			ctx := context.Background()
			if tt.setup != nil {
				tt.setup(ctx, dao)
			}

			execCtx := tt.getCtx()
			err := dao.UpdateRequest(execCtx, tt.req)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				if tt.checkUpdate != nil {
					tt.checkUpdate(t, dao)
				}
			}
		})
	}
}

func TestDAO_DeleteRequest(t *testing.T) {
	tests := []struct {
		name      string
		setup     func(ctx context.Context, dao *DAO)
		requestID string
		getCtx    func() context.Context
		wantErrIs error
		wantErr   bool
	}{
		{
			name: "delete existing request",
			setup: func(ctx context.Context, dao *DAO) {
				req := &dto.Request{
					UUID:      "r-del-1",
					Name:      "Delete Me",
					GroupUUID: dto.UnsortedName,
					Data:      &dto.Data{Method: "DELETE", URL: "https://example.com/item/1"},
				}
				require.NoError(t, dao.InsertRequest(ctx, req))
			},
			requestID: "r-del-1",
			getCtx:    func() context.Context { return context.Background() },
			wantErr:   false,
		},
		{
			name:      "delete non-existent request returns ErrNotFound",
			setup:     func(ctx context.Context, dao *DAO) {},
			requestID: "r-missing",
			getCtx:    func() context.Context { return context.Background() },
			wantErrIs: ErrNotFound,
			wantErr:   true,
		},
		{
			name:      "canceled context returns error",
			setup:     func(ctx context.Context, dao *DAO) {},
			requestID: "r-del-1",
			getCtx: func() context.Context {
				ctx, cancel := context.WithCancel(context.Background())
				cancel()
				return ctx
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dao := setupTestDAO(t)

			ctx := context.Background()
			if tt.setup != nil {
				tt.setup(ctx, dao)
			}

			execCtx := tt.getCtx()
			err := dao.DeleteRequest(execCtx, tt.requestID)
			if tt.wantErr {
				assert.Error(t, err)
				if tt.wantErrIs != nil {
					assert.ErrorIs(t, err, tt.wantErrIs)
				}
			} else {
				assert.NoError(t, err)
				_, fetchErr := dao.Request(ctx, tt.requestID)
				assert.ErrorIs(t, fetchErr, ErrNotFound)
			}
		})
	}
}

func TestDAO_DeleteGroup(t *testing.T) {
	tests := []struct {
		name      string
		setup     func(ctx context.Context, dao *DAO)
		groupID   string
		getCtx    func() context.Context
		checkPost func(t *testing.T, dao *DAO)
		wantErrIs error
		wantErr   bool
	}{
		{
			name: "delete custom group with cascading request deletion",
			setup: func(ctx context.Context, dao *DAO) {
				group := &dto.Group{UUID: "g-cascade", Name: "Cascade Group", Desc: "Will be deleted"}
				require.NoError(t, dao.InsertGroup(ctx, group))
				req := &dto.Request{
					UUID:      "r-cascaded",
					Name:      "Cascaded Request",
					GroupUUID: "g-cascade",
					Data:      &dto.Data{Method: "GET", URL: "https://example.com"},
				}
				require.NoError(t, dao.InsertRequest(ctx, req))
			},
			groupID: "g-cascade",
			getCtx:  func() context.Context { return context.Background() },
			checkPost: func(t *testing.T, dao *DAO) {
				_, err := dao.Request(context.Background(), "r-cascaded")
				assert.ErrorIs(t, err, ErrNotFound)
			},
			wantErr: false,
		},
		{
			name:      "delete unsorted group prohibited",
			setup:     func(ctx context.Context, dao *DAO) {},
			groupID:   dto.UnsortedName,
			getCtx:    func() context.Context { return context.Background() },
			wantErrIs: errors.ErrUnsupported,
			wantErr:   true,
		},
		{
			name:      "delete non-existent group returns ErrNotFound",
			setup:     func(ctx context.Context, dao *DAO) {},
			groupID:   "g-nonexistent",
			getCtx:    func() context.Context { return context.Background() },
			wantErrIs: ErrNotFound,
			wantErr:   true,
		},
		{
			name:    "canceled context returns error",
			setup:   func(ctx context.Context, dao *DAO) {},
			groupID: "g-cascade",
			getCtx: func() context.Context {
				ctx, cancel := context.WithCancel(context.Background())
				cancel()
				return ctx
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dao := setupTestDAO(t)

			ctx := context.Background()
			if tt.setup != nil {
				tt.setup(ctx, dao)
			}

			execCtx := tt.getCtx()
			err := dao.DeleteGroup(execCtx, tt.groupID)
			if tt.wantErr {
				assert.Error(t, err)
				if tt.wantErrIs != nil {
					assert.ErrorIs(t, err, tt.wantErrIs)
				}
			} else {
				assert.NoError(t, err)
				if tt.checkPost != nil {
					tt.checkPost(t, dao)
				}
			}
		})
	}
}

func TestDAO_InsertEnvironment(t *testing.T) {
	tests := []struct {
		name    string
		setup   func(ctx context.Context, dao *DAO)
		env     *dto.Environment
		getCtx  func() context.Context
		wantErr bool
	}{
		{
			name:  "insert environment without variables",
			setup: func(ctx context.Context, dao *DAO) {},
			env: &dto.Environment{
				UUID: "env-1",
				Name: "Staging",
			},
			getCtx:  func() context.Context { return context.Background() },
			wantErr: false,
		},
		{
			name:  "insert environment with variables",
			setup: func(ctx context.Context, dao *DAO) {},
			env: &dto.Environment{
				UUID: "env-2",
				Name: "Production",
				Variables: []*dto.EnvironmentVariable{
					{EnvironmentUUID: "env-2", Name: "HOST", Value: []byte("api.example.com")},
					{EnvironmentUUID: "env-2", Name: "PORT", Value: []byte("443")},
				},
			},
			getCtx:  func() context.Context { return context.Background() },
			wantErr: false,
		},
		{
			name: "insert environment with duplicate UUID",
			setup: func(ctx context.Context, dao *DAO) {
				require.NoError(t, dao.InsertEnvironment(ctx, &dto.Environment{UUID: "env-dup", Name: "First"}))
			},
			env:     &dto.Environment{UUID: "env-dup", Name: "Second"},
			getCtx:  func() context.Context { return context.Background() },
			wantErr: true,
		},
		{
			name:  "canceled context returns error",
			setup: func(ctx context.Context, dao *DAO) {},
			env:   &dto.Environment{UUID: "env-cancel", Name: "Canceled"},
			getCtx: func() context.Context {
				ctx, cancel := context.WithCancel(context.Background())
				cancel()
				return ctx
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dao := setupTestDAO(t)

			ctx := context.Background()
			if tt.setup != nil {
				tt.setup(ctx, dao)
			}

			execCtx := tt.getCtx()
			err := dao.InsertEnvironment(execCtx, tt.env)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				fetched, err := dao.Environment(ctx, tt.env.UUID)
				require.NoError(t, err)
				assert.Equal(t, tt.env.Name, fetched.Name)
				assert.Len(t, fetched.Variables, len(tt.env.Variables))
			}
		})
	}
}

func TestDAO_Environment(t *testing.T) {
	tests := []struct {
		name      string
		setup     func(ctx context.Context, dao *DAO)
		envID     string
		getCtx    func() context.Context
		wantEnv   *dto.Environment
		wantErrIs error
		wantErr   bool
	}{
		{
			name: "existing environment with variables",
			setup: func(ctx context.Context, dao *DAO) {
				env := &dto.Environment{
					UUID: "env-fetch",
					Name: "Development",
					Variables: []*dto.EnvironmentVariable{
						{EnvironmentUUID: "env-fetch", Name: "DEBUG", Value: []byte("true")},
					},
				}
				require.NoError(t, dao.InsertEnvironment(ctx, env))
			},
			envID:  "env-fetch",
			getCtx: func() context.Context { return context.Background() },
			wantEnv: &dto.Environment{
				UUID: "env-fetch",
				Name: "Development",
				Variables: []*dto.EnvironmentVariable{
					{EnvironmentUUID: "env-fetch", Name: "DEBUG", Value: []byte("true")},
				},
			},
			wantErr: false,
		},
		{
			name:      "non-existent environment returns ErrNotFound",
			setup:     func(ctx context.Context, dao *DAO) {},
			envID:     "env-missing",
			getCtx:    func() context.Context { return context.Background() },
			wantErrIs: ErrNotFound,
			wantErr:   true,
		},
		{
			name:  "canceled context returns error",
			setup: func(ctx context.Context, dao *DAO) {},
			envID: "env-fetch",
			getCtx: func() context.Context {
				ctx, cancel := context.WithCancel(context.Background())
				cancel()
				return ctx
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dao := setupTestDAO(t)

			ctx := context.Background()
			if tt.setup != nil {
				tt.setup(ctx, dao)
			}

			execCtx := tt.getCtx()
			env, err := dao.Environment(execCtx, tt.envID)
			if tt.wantErr {
				assert.Error(t, err)
				if tt.wantErrIs != nil {
					assert.ErrorIs(t, err, tt.wantErrIs)
				}
			} else {
				assert.NoError(t, err)
				assert.True(t, tt.wantEnv.Equal(env))
			}
		})
	}
}

func TestDAO_Environments(t *testing.T) {
	tests := []struct {
		name      string
		setup     func(ctx context.Context, dao *DAO)
		getCtx    func() context.Context
		wantCount int
		wantNames []string
		wantErr   bool
	}{
		{
			name:      "empty database returns empty slice",
			setup:     func(ctx context.Context, dao *DAO) {},
			getCtx:    func() context.Context { return context.Background() },
			wantCount: 0,
			wantNames: []string{},
			wantErr:   false,
		},
		{
			name: "returns all environments ordered by name and UUID",
			setup: func(ctx context.Context, dao *DAO) {
				require.NoError(t, dao.InsertEnvironment(ctx, &dto.Environment{UUID: "env-b", Name: "Beta"}))
				require.NoError(t, dao.InsertEnvironment(ctx, &dto.Environment{UUID: "env-a", Name: "Alpha"}))
			},
			getCtx:    func() context.Context { return context.Background() },
			wantCount: 2,
			wantNames: []string{"Alpha", "Beta"},
			wantErr:   false,
		},
		{
			name:   "canceled context returns error",
			setup:  func(ctx context.Context, dao *DAO) {},
			getCtx: func() context.Context {
				ctx, cancel := context.WithCancel(context.Background())
				cancel()
				return ctx
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dao := setupTestDAO(t)

			ctx := context.Background()
			if tt.setup != nil {
				tt.setup(ctx, dao)
			}

			execCtx := tt.getCtx()
			envs, err := dao.Environments(execCtx)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Len(t, envs, tt.wantCount)
				names := make([]string, len(envs))
				for i, e := range envs {
					names[i] = e.Name
				}
				assert.Equal(t, tt.wantNames, names)
			}
		})
	}
}

func TestDAO_UpdateEnvironment(t *testing.T) {
	tests := []struct {
		name      string
		setup     func(ctx context.Context, dao *DAO)
		env       *dto.Environment
		getCtx    func() context.Context
		wantErrIs error
		wantErr   bool
	}{
		{
			name: "update existing environment name",
			setup: func(ctx context.Context, dao *DAO) {
				require.NoError(t, dao.InsertEnvironment(ctx, &dto.Environment{UUID: "env-upd", Name: "Old Name"}))
			},
			env:     &dto.Environment{UUID: "env-upd", Name: "New Name"},
			getCtx:  func() context.Context { return context.Background() },
			wantErr: false,
		},
		{
			name:      "update non-existent environment returns ErrNotFound",
			setup:     func(ctx context.Context, dao *DAO) {},
			env:       &dto.Environment{UUID: "env-missing", Name: "New Name"},
			getCtx:    func() context.Context { return context.Background() },
			wantErrIs: ErrNotFound,
			wantErr:   true,
		},
		{
			name:  "canceled context returns error",
			setup: func(ctx context.Context, dao *DAO) {},
			env:   &dto.Environment{UUID: "env-upd", Name: "New Name"},
			getCtx: func() context.Context {
				ctx, cancel := context.WithCancel(context.Background())
				cancel()
				return ctx
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dao := setupTestDAO(t)

			ctx := context.Background()
			if tt.setup != nil {
				tt.setup(ctx, dao)
			}

			execCtx := tt.getCtx()
			err := dao.UpdateEnvironment(execCtx, tt.env)
			if tt.wantErr {
				assert.Error(t, err)
				if tt.wantErrIs != nil {
					assert.ErrorIs(t, err, tt.wantErrIs)
				}
			} else {
				assert.NoError(t, err)
				fetched, err := dao.Environment(ctx, tt.env.UUID)
				require.NoError(t, err)
				assert.Equal(t, "New Name", fetched.Name)
			}
		})
	}
}

func TestDAO_DeleteEnvironment(t *testing.T) {
	tests := []struct {
		name      string
		setup     func(ctx context.Context, dao *DAO)
		envID     string
		getCtx    func() context.Context
		wantErrIs error
		wantErr   bool
	}{
		{
			name: "delete environment cascades to variables",
			setup: func(ctx context.Context, dao *DAO) {
				env := &dto.Environment{
					UUID: "env-del",
					Name: "To Delete",
					Variables: []*dto.EnvironmentVariable{
						{EnvironmentUUID: "env-del", Name: "VAR1", Value: []byte("val1")},
					},
				}
				require.NoError(t, dao.InsertEnvironment(ctx, env))
			},
			envID:   "env-del",
			getCtx:  func() context.Context { return context.Background() },
			wantErr: false,
		},
		{
			name:      "delete non-existent environment returns ErrNotFound",
			setup:     func(ctx context.Context, dao *DAO) {},
			envID:     "env-missing",
			getCtx:    func() context.Context { return context.Background() },
			wantErrIs: ErrNotFound,
			wantErr:   true,
		},
		{
			name:  "canceled context returns error",
			setup: func(ctx context.Context, dao *DAO) {},
			envID: "env-del",
			getCtx: func() context.Context {
				ctx, cancel := context.WithCancel(context.Background())
				cancel()
				return ctx
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dao := setupTestDAO(t)

			ctx := context.Background()
			if tt.setup != nil {
				tt.setup(ctx, dao)
			}

			execCtx := tt.getCtx()
			err := dao.DeleteEnvironment(execCtx, tt.envID)
			if tt.wantErr {
				assert.Error(t, err)
				if tt.wantErrIs != nil {
					assert.ErrorIs(t, err, tt.wantErrIs)
				}
			} else {
				assert.NoError(t, err)
				_, fetchErr := dao.Environment(ctx, tt.envID)
				assert.ErrorIs(t, fetchErr, ErrNotFound)
			}
		})
	}
}

func TestDAO_InsertEnvironmentVariable(t *testing.T) {
	tests := []struct {
		name    string
		setup   func(ctx context.Context, dao *DAO)
		envVar  *dto.EnvironmentVariable
		getCtx  func() context.Context
		wantErr bool
	}{
		{
			name: "insert variable into existing environment",
			setup: func(ctx context.Context, dao *DAO) {
				require.NoError(t, dao.InsertEnvironment(ctx, &dto.Environment{UUID: "env-v1", Name: "Env 1"}))
			},
			envVar:  &dto.EnvironmentVariable{EnvironmentUUID: "env-v1", Name: "TOKEN", Value: []byte("xyz")},
			getCtx:  func() context.Context { return context.Background() },
			wantErr: false,
		},
		{
			name:    "insert variable into non-existent environment fails foreign key constraint",
			setup:   func(ctx context.Context, dao *DAO) {},
			envVar:  &dto.EnvironmentVariable{EnvironmentUUID: "env-missing", Name: "TOKEN", Value: []byte("xyz")},
			getCtx:  func() context.Context { return context.Background() },
			wantErr: true,
		},
		{
			name: "insert duplicate variable fails primary key constraint",
			setup: func(ctx context.Context, dao *DAO) {
				require.NoError(t, dao.InsertEnvironment(ctx, &dto.Environment{UUID: "env-v2", Name: "Env 2"}))
				require.NoError(t, dao.InsertEnvironmentVariable(ctx, &dto.EnvironmentVariable{EnvironmentUUID: "env-v2", Name: "DUP", Value: []byte("val1")}))
			},
			envVar:  &dto.EnvironmentVariable{EnvironmentUUID: "env-v2", Name: "DUP", Value: []byte("val2")},
			getCtx:  func() context.Context { return context.Background() },
			wantErr: true,
		},
		{
			name:   "canceled context returns error",
			setup:  func(ctx context.Context, dao *DAO) {},
			envVar: &dto.EnvironmentVariable{EnvironmentUUID: "env-v1", Name: "TOKEN", Value: []byte("xyz")},
			getCtx: func() context.Context {
				ctx, cancel := context.WithCancel(context.Background())
				cancel()
				return ctx
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dao := setupTestDAO(t)

			ctx := context.Background()
			if tt.setup != nil {
				tt.setup(ctx, dao)
			}

			execCtx := tt.getCtx()
			err := dao.InsertEnvironmentVariable(execCtx, tt.envVar)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				fetched, err := dao.EnvironmentVariable(ctx, tt.envVar.EnvironmentUUID, tt.envVar.Name)
				require.NoError(t, err)
				assert.True(t, tt.envVar.Equal(fetched))
			}
		})
	}
}

func TestDAO_EnvironmentVariables(t *testing.T) {
	tests := []struct {
		name      string
		setup     func(ctx context.Context, dao *DAO)
		envID     string
		getCtx    func() context.Context
		wantCount int
		wantErrIs error
		wantErr   bool
	}{
		{
			name: "existing variables for environment",
			setup: func(ctx context.Context, dao *DAO) {
				env := &dto.Environment{
					UUID: "env-listvars",
					Name: "Env ListVars",
					Variables: []*dto.EnvironmentVariable{
						{EnvironmentUUID: "env-listvars", Name: "VAR_A", Value: []byte("A")},
						{EnvironmentUUID: "env-listvars", Name: "VAR_B", Value: []byte("B")},
					},
				}
				require.NoError(t, dao.InsertEnvironment(ctx, env))
			},
			envID:     "env-listvars",
			getCtx:    func() context.Context { return context.Background() },
			wantCount: 2,
			wantErr:   false,
		},
		{
			name:      "non-existent environment returns ErrNotFound",
			setup:     func(ctx context.Context, dao *DAO) {},
			envID:     "env-missing",
			getCtx:    func() context.Context { return context.Background() },
			wantErrIs: ErrNotFound,
			wantErr:   true,
		},
		{
			name:  "canceled context returns error",
			setup: func(ctx context.Context, dao *DAO) {},
			envID: "env-listvars",
			getCtx: func() context.Context {
				ctx, cancel := context.WithCancel(context.Background())
				cancel()
				return ctx
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dao := setupTestDAO(t)

			ctx := context.Background()
			if tt.setup != nil {
				tt.setup(ctx, dao)
			}

			execCtx := tt.getCtx()
			vars, err := dao.EnvironmentVariables(execCtx, tt.envID)
			if tt.wantErr {
				assert.Error(t, err)
				if tt.wantErrIs != nil {
					assert.ErrorIs(t, err, tt.wantErrIs)
				}
			} else {
				assert.NoError(t, err)
				assert.Len(t, vars, tt.wantCount)
			}
		})
	}
}

func TestDAO_EnvironmentVariable(t *testing.T) {
	tests := []struct {
		name      string
		setup     func(ctx context.Context, dao *DAO)
		envID     string
		varName   string
		getCtx    func() context.Context
		wantVar   *dto.EnvironmentVariable
		wantErrIs error
		wantErr   bool
	}{
		{
			name: "existing environment variable",
			setup: func(ctx context.Context, dao *DAO) {
				require.NoError(t, dao.InsertEnvironment(ctx, &dto.Environment{UUID: "env-single", Name: "Single"}))
				require.NoError(t, dao.InsertEnvironmentVariable(ctx, &dto.EnvironmentVariable{EnvironmentUUID: "env-single", Name: "KEY", Value: []byte("VAL")}))
			},
			envID:   "env-single",
			varName: "KEY",
			getCtx:  func() context.Context { return context.Background() },
			wantVar: &dto.EnvironmentVariable{EnvironmentUUID: "env-single", Name: "KEY", Value: []byte("VAL")},
			wantErr: false,
		},
		{
			name:      "non-existent variable returns ErrNotFound",
			setup:     func(ctx context.Context, dao *DAO) {},
			envID:     "env-single",
			varName:   "MISSING",
			getCtx:    func() context.Context { return context.Background() },
			wantErrIs: ErrNotFound,
			wantErr:   true,
		},
		{
			name:    "canceled context returns error",
			setup:   func(ctx context.Context, dao *DAO) {},
			envID:   "env-single",
			varName: "KEY",
			getCtx: func() context.Context {
				ctx, cancel := context.WithCancel(context.Background())
				cancel()
				return ctx
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dao := setupTestDAO(t)

			ctx := context.Background()
			if tt.setup != nil {
				tt.setup(ctx, dao)
			}

			execCtx := tt.getCtx()
			v, err := dao.EnvironmentVariable(execCtx, tt.envID, tt.varName)
			if tt.wantErr {
				assert.Error(t, err)
				if tt.wantErrIs != nil {
					assert.ErrorIs(t, err, tt.wantErrIs)
				}
			} else {
				assert.NoError(t, err)
				assert.True(t, tt.wantVar.Equal(v))
			}
		})
	}
}

func TestDAO_UpdateEnvironmentVariable(t *testing.T) {
	tests := []struct {
		name      string
		setup     func(ctx context.Context, dao *DAO)
		envVar    *dto.EnvironmentVariable
		getCtx    func() context.Context
		wantErrIs error
		wantErr   bool
	}{
		{
			name: "update existing environment variable value",
			setup: func(ctx context.Context, dao *DAO) {
				require.NoError(t, dao.InsertEnvironment(ctx, &dto.Environment{UUID: "env-vupd", Name: "Upd"}))
				require.NoError(t, dao.InsertEnvironmentVariable(ctx, &dto.EnvironmentVariable{EnvironmentUUID: "env-vupd", Name: "FOO", Value: []byte("bar")}))
			},
			envVar:  &dto.EnvironmentVariable{EnvironmentUUID: "env-vupd", Name: "FOO", Value: []byte("baz")},
			getCtx:  func() context.Context { return context.Background() },
			wantErr: false,
		},
		{
			name:      "update non-existent variable returns ErrNotFound",
			setup:     func(ctx context.Context, dao *DAO) {},
			envVar:    &dto.EnvironmentVariable{EnvironmentUUID: "env-vupd", Name: "MISSING", Value: []byte("baz")},
			getCtx:    func() context.Context { return context.Background() },
			wantErrIs: ErrNotFound,
			wantErr:   true,
		},
		{
			name:   "canceled context returns error",
			setup:  func(ctx context.Context, dao *DAO) {},
			envVar: &dto.EnvironmentVariable{EnvironmentUUID: "env-vupd", Name: "FOO", Value: []byte("baz")},
			getCtx: func() context.Context {
				ctx, cancel := context.WithCancel(context.Background())
				cancel()
				return ctx
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dao := setupTestDAO(t)

			ctx := context.Background()
			if tt.setup != nil {
				tt.setup(ctx, dao)
			}

			execCtx := tt.getCtx()
			err := dao.UpdateEnvironmentVariable(execCtx, tt.envVar)
			if tt.wantErr {
				assert.Error(t, err)
				if tt.wantErrIs != nil {
					assert.ErrorIs(t, err, tt.wantErrIs)
				}
			} else {
				assert.NoError(t, err)
				fetched, err := dao.EnvironmentVariable(ctx, tt.envVar.EnvironmentUUID, tt.envVar.Name)
				require.NoError(t, err)
				assert.Equal(t, tt.envVar.Value, fetched.Value)
			}
		})
	}
}

func TestDAO_DeleteEnvironmentVariable(t *testing.T) {
	tests := []struct {
		name      string
		setup     func(ctx context.Context, dao *DAO)
		envID     string
		varName   string
		getCtx    func() context.Context
		wantErrIs error
		wantErr   bool
	}{
		{
			name: "delete existing environment variable",
			setup: func(ctx context.Context, dao *DAO) {
				require.NoError(t, dao.InsertEnvironment(ctx, &dto.Environment{UUID: "env-vdel", Name: "Del"}))
				require.NoError(t, dao.InsertEnvironmentVariable(ctx, &dto.EnvironmentVariable{EnvironmentUUID: "env-vdel", Name: "DELETE_ME", Value: []byte("123")}))
			},
			envID:   "env-vdel",
			varName: "DELETE_ME",
			getCtx:  func() context.Context { return context.Background() },
			wantErr: false,
		},
		{
			name:      "delete non-existent variable returns ErrNotFound",
			setup:     func(ctx context.Context, dao *DAO) {},
			envID:     "env-vdel",
			varName:   "NON_EXISTENT",
			getCtx:    func() context.Context { return context.Background() },
			wantErrIs: ErrNotFound,
			wantErr:   true,
		},
		{
			name:    "canceled context returns error",
			setup:   func(ctx context.Context, dao *DAO) {},
			envID:   "env-vdel",
			varName: "DELETE_ME",
			getCtx: func() context.Context {
				ctx, cancel := context.WithCancel(context.Background())
				cancel()
				return ctx
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dao := setupTestDAO(t)

			ctx := context.Background()
			if tt.setup != nil {
				tt.setup(ctx, dao)
			}

			execCtx := tt.getCtx()
			err := dao.DeleteEnvironmentVariable(execCtx, tt.envID, tt.varName)
			if tt.wantErr {
				assert.Error(t, err)
				if tt.wantErrIs != nil {
					assert.ErrorIs(t, err, tt.wantErrIs)
				}
			} else {
				assert.NoError(t, err)
				_, fetchErr := dao.EnvironmentVariable(ctx, tt.envID, tt.varName)
				assert.ErrorIs(t, fetchErr, ErrNotFound)
			}
		})
	}
}

