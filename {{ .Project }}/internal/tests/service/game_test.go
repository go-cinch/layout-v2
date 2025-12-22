package service_test

import (
	"context"
	"strconv"
	"strings"
	"testing"

	"github.com/go-cinch/common/proto/params"
	"github.com/google/uuid"

	"{{ .Computed.module_name_final }}/api/{{ .Computed.service_name_final }}"
	"{{ .Computed.module_name_final }}/internal/tests/mock"
)

func Test{{ .Computed.service_name_capitalized }}Service_Create{{ .Computed.service_name_capitalized }}(t *testing.T) {
	s := mock.{{ .Computed.service_name_capitalized }}Service()
	ctx := context.Background()
	userID := uuid.NewString()
	ctx = mock.NewContextWithUserId(ctx, userID)

	_, err := s.Create{{ .Computed.service_name_capitalized }}(ctx, &{{ .Computed.service_name_final }}.Create{{ .Computed.service_name_capitalized }}Request{
		Name: "{{ .Computed.service_name_final }}1",
	})
	if err != nil {
		t.Error(err)
		return
	}
	_, err = s.Create{{ .Computed.service_name_capitalized }}(ctx, &{{ .Computed.service_name_final }}.Create{{ .Computed.service_name_capitalized }}Request{
		Name: "{{ .Computed.service_name_final }}2",
	})
	if err != nil {
		t.Error(err)
		return
	}
	res1, _ := s.Find{{ .Computed.service_name_capitalized }}(ctx, &{{ .Computed.service_name_final }}.Find{{ .Computed.service_name_capitalized }}Request{
		Page: &params.Page{
			Disable: true,
		},
	})
	if res1 == nil || len(res1.List) != 2 {
		t.Error("res1 len must be 2")
		return
	}
	res2, err := s.Get{{ .Computed.service_name_capitalized }}(ctx, &{{ .Computed.service_name_final }}.Get{{ .Computed.service_name_capitalized }}Request{
		Id: res1.List[0].Id,
	})
	if err != nil {
		t.Error(err)
		return
	}
	if res2.Name != res1.List[0].Name {
		t.Errorf("res2.Name must be %s", res1.List[0].Name)
		return
	}
	_, err = s.Delete{{ .Computed.service_name_capitalized }}(ctx, &params.IdsRequest{
		Ids: strings.Join([]string{
			strconv.FormatUint(res1.List[0].Id, 10),
			strconv.FormatUint(res1.List[1].Id, 10),
		}, ","),
	})
	if err != nil {
		t.Error(err)
		return
	}
	res3, _ := s.Find{{ .Computed.service_name_capitalized }}(ctx, &{{ .Computed.service_name_final }}.Find{{ .Computed.service_name_capitalized }}Request{
		Page: &params.Page{
			Disable: true,
		},
	})
	if res3 == nil || len(res3.List) != 0 {
		t.Error("res3 len must be 0")
		return
	}
}
