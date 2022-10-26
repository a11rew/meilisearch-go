package meilisearch

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/valyala/fasthttp"
)

func TestClient_Version(t *testing.T) {
	tests := []struct {
		name   string
		client *Client
	}{
		{
			name:   "TestVersion",
			client: defaultClient,
		},
		{
			name:   "TestVersionWithCustomClient",
			client: customClient,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotResp, err := tt.client.GetVersion()
			require.NoError(t, err)
			require.NotNil(t, gotResp, "Version() should not return nil value")
		})
	}
}

func TestClient_TimeoutError(t *testing.T) {
	tests := []struct {
		name          string
		client        *Client
		expectedError Error
	}{
		{
			name:   "TestTimeoutError",
			client: timeoutClient,
			expectedError: Error{
				MeilisearchApiError: meilisearchApiError{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotResp, err := tt.client.GetVersion()
			require.Error(t, err)
			require.Nil(t, gotResp)
			require.Equal(t, tt.expectedError.MeilisearchApiError.Code,
				err.(*Error).MeilisearchApiError.Code)
		})
	}
}

func TestClient_GetStats(t *testing.T) {
	tests := []struct {
		name   string
		client *Client
	}{
		{
			name:   "TestGetStats",
			client: defaultClient,
		},
		{
			name:   "TestGetStatsWithCustomClient",
			client: customClient,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotResp, err := tt.client.GetStats()
			require.NoError(t, err)
			require.NotNil(t, gotResp, "GetStats() should not return nil value")
		})
	}
}

func TestClient_GetKey(t *testing.T) {
	tests := []struct {
		name   string
		client *Client
	}{
		{
			name:   "TestGetKey",
			client: defaultClient,
		},
		{
			name:   "TestGetKeyWithCustomClient",
			client: customClient,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotResp, err := tt.client.GetKeys(nil)
			require.NoError(t, err)

			gotKey, err := tt.client.GetKey(gotResp.Results[0].Key)
			require.NoError(t, err)
			require.NotNil(t, gotKey.ExpiresAt)
			require.NotNil(t, gotKey.CreatedAt)
			require.NotNil(t, gotKey.UpdatedAt)
		})
	}
}

func TestClient_GetKeys(t *testing.T) {
	type args struct {
		client  *Client
		request *KeysQuery
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "TestBasicGetKeys",
			args: args{
				client:  defaultClient,
				request: nil,
			},
		},
		{
			name: "TestGetKeysWithCustomClient",
			args: args{
				client:  customClient,
				request: nil,
			},
		},
		{
			name: "TestGetKeysWithEmptyParam",
			args: args{
				client:  defaultClient,
				request: &KeysQuery{},
			},
		},
		{
			name: "TestGetKeysWithLimit",
			args: args{
				client: defaultClient,
				request: &KeysQuery{
					Limit: 1,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotResp, err := tt.args.client.GetKeys(tt.args.request)

			require.NoError(t, err)
			require.NotNil(t, gotResp, "GetKeys() should not return nil value")
			if tt.args.request != nil && tt.args.request.Limit != 0 {
				require.Equal(t, tt.args.request.Limit, int64(len(gotResp.Results)))
			} else {
				require.GreaterOrEqual(t, len(gotResp.Results), 2)
			}
		})
	}
}

func TestClient_CreateKey(t *testing.T) {
	tests := []struct {
		name   string
		client *Client
		key    Key
	}{
		{
			name:   "TestCreateBasicKey",
			client: defaultClient,
			key: Key{
				Actions: []string{"*"},
				Indexes: []string{"*"},
			},
		},
		{
			name:   "TestCreateKeyWithCustomClient",
			client: customClient,
			key: Key{
				Actions: []string{"*"},
				Indexes: []string{"*"},
			},
		},
		{
			name:   "TestCreateKeyWithExpirationAt",
			client: defaultClient,
			key: Key{
				Actions:   []string{"*"},
				Indexes:   []string{"*"},
				ExpiresAt: time.Now().Add(time.Hour * 10),
			},
		},
		{
			name:   "TestCreateKeyWithDescription",
			client: defaultClient,
			key: Key{
				Name:        "TestCreateKeyWithDescription",
				Description: "TestCreateKeyWithDescription",
				Actions:     []string{"*"},
				Indexes:     []string{"*"},
			},
		},
		{
			name:   "TestCreateKeyWithActions",
			client: defaultClient,
			key: Key{
				Name:        "TestCreateKeyWithActions",
				Description: "TestCreateKeyWithActions",
				Actions:     []string{"documents.add", "documents.delete"},
				Indexes:     []string{"*"},
			},
		},
		{
			name:   "TestCreateKeyWithIndexes",
			client: defaultClient,
			key: Key{
				Name:        "TestCreateKeyWithIndexes",
				Description: "TestCreateKeyWithIndexes",
				Actions:     []string{"*"},
				Indexes:     []string{"movies", "games"},
			},
		},
		{
			name:   "TestCreateKeyWithWildcardedAction",
			client: defaultClient,
			key: Key{
				Name:        "TestCreateKeyWithWildcardedAction",
				Description: "TestCreateKeyWithWildcardedAction",
				Actions:     []string{"documents.*"},
				Indexes:     []string{"movies", "games"},
			},
		},
		{
			name:   "TestCreateKeyWithUID",
			client: defaultClient,
			key: Key{
				Name:    "TestCreateKeyWithUID",
				UID:     "9aec34f4-e44c-4917-86c2-9c9403abb3b6",
				Actions: []string{"*"},
				Indexes: []string{"*"},
			},
		},
		{
			name:   "TestCreateKeyWithAllOptions",
			client: defaultClient,
			key: Key{
				Name:        "TestCreateKeyWithAllOptions",
				Description: "TestCreateKeyWithAllOptions",
				UID:         "9aec34f4-e44c-4917-86c2-9c9403abb3b6",
				Actions:     []string{"documents.add", "documents.delete"},
				Indexes:     []string{"movies", "games"},
				ExpiresAt:   time.Now().Add(time.Hour * 10),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			const Format = "2006-01-02T15:04:05"
			c := tt.client
			t.Cleanup(cleanup(c))

			gotResp, err := c.CreateKey(&tt.key)
			require.NoError(t, err)

			gotKey, err := c.GetKey(gotResp.Key)
			require.NoError(t, err)
			require.Equal(t, tt.key.Name, gotKey.Name)
			require.Equal(t, tt.key.Description, gotKey.Description)
			if tt.key.UID != "" {
				require.Equal(t, tt.key.UID, gotKey.UID)
			}
			require.Equal(t, tt.key.Actions, gotKey.Actions)
			require.Equal(t, tt.key.Indexes, gotKey.Indexes)
			if !tt.key.ExpiresAt.IsZero() {
				require.Equal(t, tt.key.ExpiresAt.Format(Format), gotKey.ExpiresAt.Format(Format))
			}
		})
	}
}

func TestClient_UpdateKey(t *testing.T) {
	tests := []struct {
		name        string
		client      *Client
		keyToCreate Key
		keyToUpdate Key
	}{
		{
			name:   "TestUpdateKeyWithDescription",
			client: defaultClient,
			keyToCreate: Key{
				Actions: []string{"*"},
				Indexes: []string{"*"},
			},
			keyToUpdate: Key{
				Description: "TestUpdateKeyWithDescription",
			},
		},
		{
			name:   "TestUpdateKeyWithCustomClientWithDescription",
			client: customClient,
			keyToCreate: Key{
				Actions: []string{"*"},
				Indexes: []string{"TestUpdateKeyWithCustomClientWithDescription"},
			},
			keyToUpdate: Key{
				Description: "TestUpdateKeyWithCustomClientWithDescription",
			},
		},
		{
			name:   "TestUpdateKeyWithName",
			client: defaultClient,
			keyToCreate: Key{
				Actions: []string{"*"},
				Indexes: []string{"TestUpdateKeyWithName"},
			},
			keyToUpdate: Key{
				Name: "TestUpdateKeyWithName",
			},
		},
		{
			name:   "TestUpdateKeyWithNameAndAction",
			client: defaultClient,
			keyToCreate: Key{
				Actions: []string{"search"},
				Indexes: []string{"*"},
			},
			keyToUpdate: Key{
				Name: "TestUpdateKeyWithName",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			const Format = "2006-01-02T15:04:05"
			c := tt.client
			t.Cleanup(cleanup(c))

			gotResp, err := c.CreateKey(&tt.keyToCreate)
			require.NoError(t, err)

			if tt.keyToCreate.Description != "" {
				require.Equal(t, tt.keyToCreate.Description, gotResp.Description)
			}
			if len(tt.keyToCreate.Actions) != 0 {
				require.Equal(t, tt.keyToCreate.Actions, gotResp.Actions)
			}
			if len(tt.keyToCreate.Indexes) != 0 {
				require.Equal(t, tt.keyToCreate.Indexes, gotResp.Indexes)
			}
			if !tt.keyToCreate.ExpiresAt.IsZero() {
				require.Equal(t, tt.keyToCreate.ExpiresAt.Format(Format), gotResp.ExpiresAt.Format(Format))
			}

			gotKey, err := c.UpdateKey(gotResp.Key, &tt.keyToUpdate)
			require.NoError(t, err)

			if tt.keyToUpdate.Description != "" {
				require.Equal(t, tt.keyToUpdate.Description, gotKey.Description)
			}
			if len(tt.keyToUpdate.Actions) != 0 {
				require.Equal(t, tt.keyToUpdate.Actions, gotKey.Actions)
			}
			if len(tt.keyToUpdate.Indexes) != 0 {
				require.Equal(t, tt.keyToUpdate.Indexes, gotKey.Indexes)
			}
			if tt.keyToUpdate.Description != "" {
				require.Equal(t, tt.keyToUpdate.Name, gotKey.Name)
			}
		})
	}
}

func TestClient_DeleteKey(t *testing.T) {
	tests := []struct {
		name   string
		client *Client
		key    Key
	}{
		{
			name:   "TestDeleteBasicKey",
			client: defaultClient,
			key: Key{
				Actions: []string{"*"},
				Indexes: []string{"*"},
			},
		},
		{
			name:   "TestDeleteKeyWithCustomClient",
			client: customClient,
			key: Key{
				Actions: []string{"*"},
				Indexes: []string{"*"},
			},
		},
		{
			name:   "TestDeleteKeyWithExpirationAt",
			client: defaultClient,
			key: Key{
				Actions:   []string{"*"},
				Indexes:   []string{"*"},
				ExpiresAt: time.Now().Add(time.Hour * 10),
			},
		},
		{
			name:   "TestDeleteKeyWithDescription",
			client: defaultClient,
			key: Key{
				Description: "TestDeleteKeyWithDescription",
				Actions:     []string{"*"},
				Indexes:     []string{"*"},
			},
		},
		{
			name:   "TestDeleteKeyWithActions",
			client: defaultClient,
			key: Key{
				Description: "TestDeleteKeyWithActions",
				Actions:     []string{"documents.add", "documents.delete"},
				Indexes:     []string{"*"},
			},
		},
		{
			name:   "TestDeleteKeyWithIndexes",
			client: defaultClient,
			key: Key{
				Description: "TestDeleteKeyWithIndexes",
				Actions:     []string{"*"},
				Indexes:     []string{"movies", "games"},
			},
		},
		{
			name:   "TestDeleteKeyWithAllOptions",
			client: defaultClient,
			key: Key{
				Description: "TestDeleteKeyWithAllOptions",
				Actions:     []string{"documents.add", "documents.delete"},
				Indexes:     []string{"movies", "games"},
				ExpiresAt:   time.Now().Add(time.Hour * 10),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := tt.client

			gotKey, err := c.CreateKey(&tt.key)
			require.NoError(t, err)

			gotResp, err := c.DeleteKey(gotKey.Key)
			require.NoError(t, err)
			require.True(t, gotResp)

			gotResp, err = c.DeleteKey(gotKey.Key)
			require.Error(t, err)
			require.False(t, gotResp)
		})
	}
}

func TestClient_Health(t *testing.T) {
	tests := []struct {
		name     string
		client   *Client
		wantResp *Health
		wantErr  bool
	}{
		{
			name:   "TestHealth",
			client: defaultClient,
			wantResp: &Health{
				Status: "available",
			},
			wantErr: false,
		},
		{
			name:   "TestHealthWithCustomClient",
			client: customClient,
			wantResp: &Health{
				Status: "available",
			},
			wantErr: false,
		},
		{
			name: "TestHealthWithBadUrl",
			client: &Client{
				config: ClientConfig{
					Host:   "http://wrongurl:1234",
					APIKey: masterKey,
				},
				httpClient: &fasthttp.Client{
					Name: "meilisearch-client",
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotResp, err := tt.client.Health()
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.wantResp, gotResp, "Health() got response %v, want %v", gotResp, tt.wantResp)
			}
		})
	}
}

func TestClient_IsHealthy(t *testing.T) {
	tests := []struct {
		name   string
		client *Client
		want   bool
	}{
		{
			name:   "TestIsHealthy",
			client: defaultClient,
			want:   true,
		},
		{
			name:   "TestIsHealthyWithCustomClient",
			client: customClient,
			want:   true,
		},
		{
			name: "TestIsHealthyWIthBadUrl",
			client: &Client{
				config: ClientConfig{
					Host:   "http://wrongurl:1234",
					APIKey: masterKey,
				},
				httpClient: &fasthttp.Client{
					Name: "meilisearch-client",
				},
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.client.IsHealthy()
			require.Equal(t, tt.want, got, "IsHealthy() got response %v, want %v", got, tt.want)
		})
	}
}

func TestClient_CreateDump(t *testing.T) {
	tests := []struct {
		name     string
		client   *Client
		wantResp *Task
	}{
		{
			name:   "TestCreateDump",
			client: defaultClient,
			wantResp: &Task{
				Status: "enqueued",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := tt.client

			gotResp, err := c.CreateDump()
			require.NoError(t, err)
			if assert.NotNil(t, gotResp, "CreateDump() should not return nil value") {
				require.Equal(t, tt.wantResp.Status, gotResp.Status, "CreateDump() got response status %v, want: %v", gotResp.Status, tt.wantResp.Status)
			}

			// Waiting for CreateDump() to finished
			for {
				gotResp, _ := c.GetTask(gotResp.TaskUID)
				if gotResp.Status == "succeeded" {
					break
				}
			}
		})
	}
}

func TestClient_GetTask(t *testing.T) {
	type args struct {
		UID      string
		client   *Client
		taskUID  int64
		document []docTest
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "TestBasicGetTask",
			args: args{
				UID:     "TestBasicGetTask",
				client:  defaultClient,
				taskUID: 0,
				document: []docTest{
					{ID: "123", Name: "Pride and Prejudice"},
				},
			},
		},
		{
			name: "TestGetTaskWithCustomClient",
			args: args{
				UID:     "TestGetTaskWithCustomClient",
				client:  customClient,
				taskUID: 1,
				document: []docTest{
					{ID: "123", Name: "Pride and Prejudice"},
				},
			},
		},
		{
			name: "TestGetTask",
			args: args{
				UID:     "TestGetTask",
				client:  defaultClient,
				taskUID: 2,
				document: []docTest{
					{ID: "456", Name: "Le Petit Prince"},
					{ID: "1", Name: "Alice In Wonderland"},
				},
			},
		},
	}

	t.Cleanup(cleanup(defaultClient))

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := tt.args.client
			i := c.Index(tt.args.UID)
			t.Cleanup(cleanup(c))

			task, err := i.AddDocuments(tt.args.document)
			require.NoError(t, err)

			_, err = c.WaitForTask(task.TaskUID)
			require.NoError(t, err)

			gotResp, err := c.GetTask(task.TaskUID)
			require.NoError(t, err)
			require.NotNil(t, gotResp)
			require.NotNil(t, gotResp.Details)
			require.GreaterOrEqual(t, gotResp.UID, tt.args.taskUID)
			require.Equal(t, tt.args.UID, gotResp.IndexUID)
			require.Equal(t, TaskStatusSucceeded, gotResp.Status)
			require.Equal(t, len(tt.args.document), gotResp.Details.ReceivedDocuments)
			require.Equal(t, len(tt.args.document), gotResp.Details.IndexedDocuments)

			// Make sure that timestamps are also retrieved
			require.NotZero(t, gotResp.EnqueuedAt)
			require.NotZero(t, gotResp.StartedAt)
			require.NotZero(t, gotResp.FinishedAt)
		})
	}
}

func TestClient_GetTasks(t *testing.T) {
	type args struct {
		UID      string
		client   *Client
		document []docTest
		query    *TasksQuery
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "TestBasicGetTasks",
			args: args{
				UID:    "indexUID",
				client: defaultClient,
				document: []docTest{
					{ID: "123", Name: "Pride and Prejudice"},
				},
				query: nil,
			},
		},
		{
			name: "TestGetTasksWithCustomClient",
			args: args{
				UID:    "indexUID",
				client: customClient,
				document: []docTest{
					{ID: "123", Name: "Pride and Prejudice"},
				},
				query: nil,
			},
		},
		{
			name: "TestGetTasksWithLimit",
			args: args{
				UID:    "indexUID",
				client: defaultClient,
				document: []docTest{
					{ID: "123", Name: "Pride and Prejudice"},
				},
				query: &TasksQuery{
					Limit: 1,
				},
			},
		},
		{
			name: "TestGetTasksWithLimit",
			args: args{
				UID:    "indexUID",
				client: defaultClient,
				document: []docTest{
					{ID: "123", Name: "Pride and Prejudice"},
				},
				query: &TasksQuery{
					Limit: 1,
				},
			},
		},
		{
			name: "TestGetTasksWithFrom",
			args: args{
				UID:    "indexUID",
				client: defaultClient,
				document: []docTest{
					{ID: "123", Name: "Pride and Prejudice"},
				},
				query: &TasksQuery{
					From: 0,
				},
			},
		},
		{
			name: "TestGetTasksWithParameters",
			args: args{
				UID:    "indexUID",
				client: defaultClient,
				document: []docTest{
					{ID: "123", Name: "Pride and Prejudice"},
				},
				query: &TasksQuery{
					Limit:    1,
					From:     0,
					IndexUID: []string{"indexUID"},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := tt.args.client
			i := c.Index(tt.args.UID)
			t.Cleanup(cleanup(c))

			task, err := i.AddDocuments(tt.args.document)
			require.NoError(t, err)

			_, err = c.WaitForTask(task.TaskUID)
			require.NoError(t, err)

			gotResp, err := i.GetTasks(tt.args.query)
			require.NoError(t, err)
			require.NotNil(t, (*gotResp).Results[0].Status)
			require.NotZero(t, (*gotResp).Results[0].UID)
			require.NotNil(t, (*gotResp).Results[0].Type)
			if tt.args.query != nil {
				if tt.args.query.Limit != 0 {
					require.Equal(t, tt.args.query.Limit, (*gotResp).Limit)
				} else {
					require.Equal(t, int64(20), (*gotResp).Limit)
				}
				if tt.args.query.From != 0 {
					require.Equal(t, tt.args.query.From, (*gotResp).From)
				}
			}
		})
	}
}

func TestClient_DefaultWaitForTask(t *testing.T) {
	type args struct {
		UID      string
		client   *Client
		taskUID  *Task
		document []docTest
	}
	tests := []struct {
		name string
		args args
		want TaskStatus
	}{
		{
			name: "TestDefaultWaitForTask",
			args: args{
				UID:    "TestDefaultWaitForTask",
				client: defaultClient,
				taskUID: &Task{
					UID: 0,
				},
				document: []docTest{
					{ID: "123", Name: "Pride and Prejudice"},
					{ID: "456", Name: "Le Petit Prince"},
					{ID: "1", Name: "Alice In Wonderland"},
				},
			},
			want: "succeeded",
		},
		{
			name: "TestDefaultWaitForTaskWithCustomClient",
			args: args{
				UID:    "TestDefaultWaitForTaskWithCustomClient",
				client: customClient,
				taskUID: &Task{
					UID: 0,
				},
				document: []docTest{
					{ID: "123", Name: "Pride and Prejudice"},
					{ID: "456", Name: "Le Petit Prince"},
					{ID: "1", Name: "Alice In Wonderland"},
				},
			},
			want: "succeeded",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := tt.args.client
			t.Cleanup(cleanup(c))

			task, err := c.Index(tt.args.UID).AddDocuments(tt.args.document)
			require.NoError(t, err)

			gotTask, err := c.WaitForTask(task.TaskUID)
			require.NoError(t, err)
			require.Equal(t, tt.want, gotTask.Status)
		})
	}
}

func TestClient_WaitForTaskWithContext(t *testing.T) {
	type args struct {
		UID      string
		client   *Client
		interval time.Duration
		timeout  time.Duration
		taskUID  *Task
		document []docTest
	}
	tests := []struct {
		name string
		args args
		want TaskStatus
	}{
		{
			name: "TestWaitForTask50",
			args: args{
				UID:      "TestWaitForTask50",
				client:   defaultClient,
				interval: time.Millisecond * 50,
				timeout:  time.Second * 5,
				taskUID: &Task{
					UID: 0,
				},
				document: []docTest{
					{ID: "123", Name: "Pride and Prejudice"},
					{ID: "456", Name: "Le Petit Prince"},
					{ID: "1", Name: "Alice In Wonderland"},
				},
			},
			want: "succeeded",
		},
		{
			name: "TestWaitForTask50WithCustomClient",
			args: args{
				UID:      "TestWaitForTask50WithCustomClient",
				client:   customClient,
				interval: time.Millisecond * 50,
				timeout:  time.Second * 5,
				taskUID: &Task{
					UID: 0,
				},
				document: []docTest{
					{ID: "123", Name: "Pride and Prejudice"},
					{ID: "456", Name: "Le Petit Prince"},
					{ID: "1", Name: "Alice In Wonderland"},
				},
			},
			want: "succeeded",
		},
		{
			name: "TestWaitForTask10",
			args: args{
				UID:      "TestWaitForTask10",
				client:   defaultClient,
				interval: time.Millisecond * 10,
				timeout:  time.Second * 5,
				taskUID: &Task{
					UID: 1,
				},
				document: []docTest{
					{ID: "123", Name: "Pride and Prejudice"},
					{ID: "456", Name: "Le Petit Prince"},
					{ID: "1", Name: "Alice In Wonderland"},
				},
			},
			want: "succeeded",
		},
		{
			name: "TestWaitForTaskWithTimeout",
			args: args{
				UID:      "TestWaitForTaskWithTimeout",
				client:   defaultClient,
				interval: time.Millisecond * 50,
				timeout:  time.Millisecond * 10,
				taskUID: &Task{
					UID: 1,
				},
				document: []docTest{
					{ID: "123", Name: "Pride and Prejudice"},
					{ID: "456", Name: "Le Petit Prince"},
					{ID: "1", Name: "Alice In Wonderland"},
				},
			},
			want: "succeeded",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := tt.args.client
			t.Cleanup(cleanup(c))

			task, err := c.Index(tt.args.UID).AddDocuments(tt.args.document)
			require.NoError(t, err)

			ctx, cancelFunc := context.WithTimeout(context.Background(), tt.args.timeout)
			defer cancelFunc()

			gotTask, err := c.WaitForTask(task.TaskUID, WaitParams{Context: ctx, Interval: tt.args.interval})
			if tt.args.timeout < tt.args.interval {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.want, gotTask.Status)
			}
		})
	}
}

func TestClient_ConnectionCloseByServer(t *testing.T) {
	t.Skip("Skip until <https://github.com/meilisearch/meilisearch/pull/2471> merged.")

	meili := NewClient(ClientConfig{Host: "http://localhost:7700"})

	// Simulate 10 clients sending requests.
	g := sync.WaitGroup{}
	for i := 0; i < 10; i++ {
		g.Add(1)
		go func() {
			defer g.Done()

			_, _ = meili.Index("foo").Search("bar", &SearchRequest{})
			time.Sleep(5 * time.Second)
			_, err := meili.Index("foo").Search("bar", &SearchRequest{})
			if e, ok := err.(*Error); ok && e.ErrCode == MeilisearchCommunicationError {
				require.NoErrorf(t, e, "unexpected error")
			}
		}()
	}
	g.Wait()
}

func TestClient_GenerateTenantToken(t *testing.T) {
	type args struct {
		IndexUID    string
		client      *Client
		APIKeyUID   string
		searchRules Unknown
		options     *TenantTokenOptions
		filter      []string
	}
	tests := []struct {
		name       string
		args       args
		wantErr    bool
		wantFilter bool
	}{
		{
			name: "TestDefaultGenerateTenantToken",
			args: args{
				IndexUID:  "TestDefaultGenerateTenantToken",
				client:    privateClient,
				APIKeyUID: GetPrivateUIDKey(),
				searchRules: Unknown{
					"*": map[string]string{},
				},
				options: nil,
				filter:  nil,
			},
			wantErr:    false,
			wantFilter: false,
		},
		{
			name: "TestGenerateTenantTokenWithApiKey",
			args: args{
				IndexUID:  "TestGenerateTenantTokenWithApiKey",
				client:    defaultClient,
				APIKeyUID: GetPrivateUIDKey(),
				searchRules: Unknown{
					"*": map[string]string{},
				},
				options: &TenantTokenOptions{
					APIKey: GetPrivateKey(),
				},
				filter: nil,
			},
			wantErr:    false,
			wantFilter: false,
		},
		{
			name: "TestGenerateTenantTokenWithOnlyExpiresAt",
			args: args{
				IndexUID:  "TestGenerateTenantTokenWithOnlyExpiresAt",
				client:    privateClient,
				APIKeyUID: GetPrivateUIDKey(),
				searchRules: Unknown{
					"*": map[string]string{},
				},
				options: &TenantTokenOptions{
					ExpiresAt: time.Now().Add(time.Hour * 10),
				},
				filter: nil,
			},
			wantErr:    false,
			wantFilter: false,
		},
		{
			name: "TestGenerateTenantTokenWithApiKeyAndExpiresAt",
			args: args{
				IndexUID:  "TestGenerateTenantTokenWithApiKeyAndExpiresAt",
				client:    defaultClient,
				APIKeyUID: GetPrivateUIDKey(),
				searchRules: Unknown{
					"*": map[string]string{},
				},
				options: &TenantTokenOptions{
					APIKey:    GetPrivateKey(),
					ExpiresAt: time.Now().Add(time.Hour * 10),
				},
				filter: nil,
			},
			wantErr:    false,
			wantFilter: false,
		},
		{
			name: "TestGenerateTenantTokenWithFilters",
			args: args{
				IndexUID:  "indexUID",
				client:    privateClient,
				APIKeyUID: GetPrivateUIDKey(),
				searchRules: Unknown{
					"*": map[string]string{
						"filter": "book_id > 1000",
					},
				},
				options: nil,
				filter: []string{
					"book_id",
				},
			},
			wantErr:    false,
			wantFilter: true,
		},
		{
			name: "TestGenerateTenantTokenWithFilterOnOneINdex",
			args: args{
				IndexUID:  "indexUID",
				client:    privateClient,
				APIKeyUID: GetPrivateUIDKey(),
				searchRules: Unknown{
					"indexUID": map[string]string{
						"filter": "year > 2000",
					},
				},
				options: nil,
				filter: []string{
					"year",
				},
			},
			wantErr:    false,
			wantFilter: true,
		},
		{
			name: "TestGenerateTenantTokenWithoutSearchRules",
			args: args{
				IndexUID:    "TestGenerateTenantTokenWithoutSearchRules",
				client:      privateClient,
				APIKeyUID:   GetPrivateUIDKey(),
				searchRules: nil,
				options:     nil,
				filter:      nil,
			},
			wantErr:    true,
			wantFilter: false,
		},
		{
			name: "TestGenerateTenantTokenWithoutApiKey",
			args: args{
				IndexUID: "TestGenerateTenantTokenWithoutApiKey",
				client: NewClient(ClientConfig{
					Host:   getenv("MEILISEARCH_URL", "http://localhost:7700"),
					APIKey: "",
				}),
				APIKeyUID: GetPrivateUIDKey(),
				searchRules: Unknown{
					"*": map[string]string{},
				},
				options: nil,
				filter:  nil,
			},
			wantErr:    true,
			wantFilter: false,
		},
		{
			name: "TestGenerateTenantTokenWithBadExpiresAt",
			args: args{
				IndexUID:  "TestGenerateTenantTokenWithBadExpiresAt",
				client:    defaultClient,
				APIKeyUID: GetPrivateUIDKey(),
				searchRules: Unknown{
					"*": map[string]string{},
				},
				options: &TenantTokenOptions{
					ExpiresAt: time.Now().Add(-time.Hour * 10),
				},
				filter: nil,
			},
			wantErr:    true,
			wantFilter: false,
		},
		{
			name: "TestGenerateTenantTokenWithBadAPIKeyUID",
			args: args{
				IndexUID:  "TestGenerateTenantTokenWithBadAPIKeyUID",
				client:    defaultClient,
				APIKeyUID: GetPrivateUIDKey() + "1234",
				searchRules: Unknown{
					"*": map[string]string{},
				},
				options: nil,
				filter:  nil,
			},
			wantErr:    true,
			wantFilter: false,
		},
		{
			name: "TestGenerateTenantTokenWithEmptyAPIKeyUID",
			args: args{
				IndexUID:  "TestGenerateTenantTokenWithEmptyAPIKeyUID",
				client:    defaultClient,
				APIKeyUID: "",
				searchRules: Unknown{
					"*": map[string]string{},
				},
				options: nil,
				filter:  nil,
			},
			wantErr:    true,
			wantFilter: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := tt.args.client
			t.Cleanup(cleanup(c))

			token, err := c.GenerateTenantToken(tt.args.APIKeyUID, tt.args.searchRules, tt.args.options)

			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)

				if tt.wantFilter {
					gotTask, err := c.Index(tt.args.IndexUID).UpdateFilterableAttributes(&tt.args.filter)
					require.NoError(t, err, "UpdateFilterableAttributes() in TestGenerateTenantToken error should be nil")
					testWaitForTask(t, c.Index(tt.args.IndexUID), gotTask)
				} else {
					_, err := SetUpEmptyIndex(&IndexConfig{Uid: tt.args.IndexUID})
					require.NoError(t, err, "CreateIndex() in TestGenerateTenantToken error should be nil")
				}

				client := NewClient(ClientConfig{
					Host:   getenv("MEILISEARCH_URL", "http://localhost:7700"),
					APIKey: token,
				})

				_, err = client.Index(tt.args.IndexUID).Search("", &SearchRequest{})

				require.NoError(t, err)
			}
		})
	}
}
