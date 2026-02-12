package service_test

import (
	"context"
	"strconv"
	"strings"
	"testing"

	"{{.Computed.common_module_final}}/proto/params"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"{{ .Computed.module_name_final }}/api/{{ .Computed.service_name_snake }}"
	"{{ .Computed.module_name_final }}/internal/tests/mock"
)

func TestCreate{{ .Computed.service_name_capitalized }}(t *testing.T) {
	s := mock.{{ .Computed.service_name_capitalized }}Service()
	ctx := mock.NewContextWithUserId(context.Background(), uuid.NewString())

	_, err := s.Create{{ .Computed.service_name_capitalized }}(ctx, &{{ .Computed.service_name_snake }}.Create{{ .Computed.service_name_capitalized }}Request{
		Name: "test-{{ .Computed.service_name_final }}-1",
	})
	assert.NoError(t, err)
}

func TestGet{{ .Computed.service_name_capitalized }}(t *testing.T) {
	s := mock.{{ .Computed.service_name_capitalized }}Service()
	ctx := mock.NewContextWithUserId(context.Background(), uuid.NewString())

	// Create a record first
	_, err := s.Create{{ .Computed.service_name_capitalized }}(ctx, &{{ .Computed.service_name_snake }}.Create{{ .Computed.service_name_capitalized }}Request{
		Name: "test-get-{{ .Computed.service_name_final }}",
	})
	assert.NoError(t, err)

	// Find to get the ID
	res, err := s.Find{{ .Computed.service_name_capitalized }}(ctx, &{{ .Computed.service_name_snake }}.Find{{ .Computed.service_name_capitalized }}Request{
		Page: &params.Page{Disable: true},
	})
	assert.NoError(t, err)
	assert.NotEmpty(t, res.List)

	// Get by ID
	item, err := s.Get{{ .Computed.service_name_capitalized }}(ctx, &{{ .Computed.service_name_snake }}.Get{{ .Computed.service_name_capitalized }}Request{
		Id: res.List[0].Id,
	})
	assert.NoError(t, err)
	assert.Equal(t, res.List[0].Name, item.Name)
}

func TestFind{{ .Computed.service_name_capitalized }}(t *testing.T) {
	s := mock.{{ .Computed.service_name_capitalized }}Service()
	ctx := mock.NewContextWithUserId(context.Background(), uuid.NewString())

	// Create test data
	_, _ = s.Create{{ .Computed.service_name_capitalized }}(ctx, &{{ .Computed.service_name_snake }}.Create{{ .Computed.service_name_capitalized }}Request{
		Name: "test-find-{{ .Computed.service_name_final }}-1",
	})
	_, _ = s.Create{{ .Computed.service_name_capitalized }}(ctx, &{{ .Computed.service_name_snake }}.Create{{ .Computed.service_name_capitalized }}Request{
		Name: "test-find-{{ .Computed.service_name_final }}-2",
	})

	// Find with pagination disabled
	res, err := s.Find{{ .Computed.service_name_capitalized }}(ctx, &{{ .Computed.service_name_snake }}.Find{{ .Computed.service_name_capitalized }}Request{
		Page: &params.Page{Disable: true},
	})
	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.NotNil(t, res.Page)
	assert.GreaterOrEqual(t, len(res.List), 2)

	// Find with pagination
	res2, err := s.Find{{ .Computed.service_name_capitalized }}(ctx, &{{ .Computed.service_name_snake }}.Find{{ .Computed.service_name_capitalized }}Request{
		Page: &params.Page{
			Num:  1,
			Size: 10,
		},
	})
	assert.NoError(t, err)
	assert.NotNil(t, res2)
	assert.NotNil(t, res2.Page)
}

func TestUpdate{{ .Computed.service_name_capitalized }}(t *testing.T) {
	s := mock.{{ .Computed.service_name_capitalized }}Service()
	ctx := mock.NewContextWithUserId(context.Background(), uuid.NewString())

	// Create a record first
	_, err := s.Create{{ .Computed.service_name_capitalized }}(ctx, &{{ .Computed.service_name_snake }}.Create{{ .Computed.service_name_capitalized }}Request{
		Name: "test-update-{{ .Computed.service_name_final }}",
	})
	assert.NoError(t, err)

	// Find to get the ID
	res, err := s.Find{{ .Computed.service_name_capitalized }}(ctx, &{{ .Computed.service_name_snake }}.Find{{ .Computed.service_name_capitalized }}Request{
		Page: &params.Page{Disable: true},
	})
	assert.NoError(t, err)
	assert.NotEmpty(t, res.List)

	// Update
	newName := "test-update-{{ .Computed.service_name_final }}-updated"
	_, err = s.Update{{ .Computed.service_name_capitalized }}(ctx, &{{ .Computed.service_name_snake }}.Update{{ .Computed.service_name_capitalized }}Request{
		Id:   res.List[0].Id,
		Name: &newName,
	})
	assert.NoError(t, err)

	// Verify update
	item, err := s.Get{{ .Computed.service_name_capitalized }}(ctx, &{{ .Computed.service_name_snake }}.Get{{ .Computed.service_name_capitalized }}Request{
		Id: res.List[0].Id,
	})
	assert.NoError(t, err)
	assert.Equal(t, newName, item.Name)
}

func TestDelete{{ .Computed.service_name_capitalized }}(t *testing.T) {
	s := mock.{{ .Computed.service_name_capitalized }}Service()
	ctx := mock.NewContextWithUserId(context.Background(), uuid.NewString())

	// Create test data
	_, _ = s.Create{{ .Computed.service_name_capitalized }}(ctx, &{{ .Computed.service_name_snake }}.Create{{ .Computed.service_name_capitalized }}Request{
		Name: "test-delete-{{ .Computed.service_name_final }}-1",
	})
	_, _ = s.Create{{ .Computed.service_name_capitalized }}(ctx, &{{ .Computed.service_name_snake }}.Create{{ .Computed.service_name_capitalized }}Request{
		Name: "test-delete-{{ .Computed.service_name_final }}-2",
	})

	// Find to get IDs
	res, err := s.Find{{ .Computed.service_name_capitalized }}(ctx, &{{ .Computed.service_name_snake }}.Find{{ .Computed.service_name_capitalized }}Request{
		Page: &params.Page{Disable: true},
	})
	assert.NoError(t, err)
	assert.GreaterOrEqual(t, len(res.List), 2)

	// Delete
	ids := make([]string, len(res.List))
	for i, item := range res.List {
		ids[i] = strconv.FormatUint(item.Id, 10)
	}
	_, err = s.Delete{{ .Computed.service_name_capitalized }}(ctx, &params.IdsRequest{
		Ids: strings.Join(ids, ","),
	})
	assert.NoError(t, err)

	// Verify deletion
	res2, _ := s.Find{{ .Computed.service_name_capitalized }}(ctx, &{{ .Computed.service_name_snake }}.Find{{ .Computed.service_name_capitalized }}Request{
		Page: &params.Page{Disable: true},
	})
	assert.Equal(t, 0, len(res2.List))
}
