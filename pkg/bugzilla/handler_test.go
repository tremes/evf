package bugzilla

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

type mockClient struct {
	bugs []Bug
}

func (m mockClient) GetAllBugs(ctx context.Context, param SearchParams) ([]Bug, error) {

	return m.bugs, nil
}

func (m mockClient) GetComments(ctx context.Context, bzID int64) (*Comments, error) {
	switch bzID {
	case 1:
		return &Comments{Comments: []Comment{
			{
				Creator: "Joe Dev",
				Text:    "hey this is really annoying bug",
			},
			{
				Creator: CommentCreator,
				Text:    "I was working on this bug\n, but Bug report changed to RELEASE_PENDING status. advisory/11112345",
			},
		}}, nil
	case 2:
		return &Comments{Comments: []Comment{
			{
				Creator: "Joe Dev",
				Text:    "hey this is really annoying bug",
			},
			{
				Creator: CommentCreator,
				Text:    "I was working on this bug\n, but Bug report changed to RELEASE_PENDING status. advisory/11112345",
			},
		}}, nil
	}
	return nil, nil
}
func Test_Handler_Mapping(t *testing.T) {
	mockClient := mockClient{
		bugs: []Bug{
			{
				ID:      1,
				Summary: "testing bug 1",
			},
			{
				ID:      2,
				Summary: "testing bug 2",
			},
		},
	}
	h := NewHandler(mockClient)
	bugs, err := mockClient.GetAllBugs(context.Background(), SearchParams{})
	assert.NoError(t, err)
	assert.Greater(t, len(bugs), 0)
	m := h.CreateBZToErrataMap(context.Background(), bugs)
	assert.Len(t, m, 1)
	errataBugs, ok := m["11112345"]
	assert.True(t, ok)
	assert.Len(t, errataBugs, 2)
}
