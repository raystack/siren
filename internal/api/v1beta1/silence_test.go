package v1beta1_test

import (
	"context"
	"reflect"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/goto/salt/log"
	"github.com/goto/siren/core/silence"
	"github.com/goto/siren/internal/api"
	"github.com/goto/siren/internal/api/mocks"
	"github.com/goto/siren/internal/api/v1beta1"
	"github.com/goto/siren/pkg/errors"
	sirenv1beta1 "github.com/goto/siren/proto/gotocompany/siren/v1beta1"
	"github.com/stretchr/testify/mock"
	"google.golang.org/protobuf/testing/protocmp"
	"google.golang.org/protobuf/types/known/structpb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestGRPCServer_CreateSilence(t *testing.T) {
	mockSilenceData := silence.Silence{
		NamespaceID: 1,
		Type:        silence.TypeMatchers,
		TargetExpression: map[string]any{
			"key1": "value1",
		},
	}

	tests := []struct {
		name    string
		setup   func(*mocks.SilenceService)
		req     *sirenv1beta1.CreateSilenceRequest
		want    *sirenv1beta1.CreateSilenceResponse
		wantErr bool
	}{
		{
			name: "return silence id when successfully created silence",
			setup: func(ss *mocks.SilenceService) {
				ss.EXPECT().Create(mock.AnythingOfType("*context.emptyCtx"), mockSilenceData).Return("123", nil)
			},
			req: &sirenv1beta1.CreateSilenceRequest{
				NamespaceId: mockSilenceData.NamespaceID,
				Type:        mockSilenceData.Type,
				TargetExpression: &structpb.Struct{
					Fields: map[string]*structpb.Value{
						"key1": structpb.NewStringValue("value1"),
					},
				},
			},
			want: &sirenv1beta1.CreateSilenceResponse{
				Id: "123",
			},
		},
		{
			name: "return error if service create return error",
			setup: func(ss *mocks.SilenceService) {
				ss.EXPECT().Create(mock.AnythingOfType("*context.emptyCtx"), mockSilenceData).Return("", errors.New("some error"))
			},
			req: &sirenv1beta1.CreateSilenceRequest{
				NamespaceId: mockSilenceData.NamespaceID,
				Type:        mockSilenceData.Type,
				TargetExpression: &structpb.Struct{
					Fields: map[string]*structpb.Value{
						"key1": structpb.NewStringValue("value1"),
					},
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		ctx := context.TODO()
		t.Run(tt.name, func(t *testing.T) {
			mockSilenceService := new(mocks.SilenceService)

			if tt.setup != nil {
				tt.setup(mockSilenceService)
			}

			s := v1beta1.NewGRPCServer(nil, log.NewNoop(), api.HeadersConfig{}, &api.Deps{SilenceService: mockSilenceService})
			got, err := s.CreateSilence(ctx, tt.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("GRPCServer.CreateSilence() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GRPCServer.CreateSilence() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGRPCServer_ListSilences(t *testing.T) {
	mockSilenceData := silence.Silence{
		NamespaceID: 1,
		Type:        silence.TypeMatchers,
		TargetExpression: map[string]any{
			"key1": "value1",
		},
	}

	tests := []struct {
		name    string
		setup   func(*mocks.SilenceService)
		want    []*sirenv1beta1.Silence
		wantErr bool
	}{
		{
			name: "return silences when successfully list silences",
			setup: func(ss *mocks.SilenceService) {
				ss.EXPECT().List(mock.AnythingOfType("*context.emptyCtx"), mock.AnythingOfType("silence.Filter")).Return([]silence.Silence{
					mockSilenceData,
				}, nil)
			},
			want: []*sirenv1beta1.Silence{
				{
					NamespaceId: mockSilenceData.NamespaceID,
					Type:        mockSilenceData.Type,
					TargetExpression: &structpb.Struct{
						Fields: map[string]*structpb.Value{
							"key1": structpb.NewStringValue("value1"),
						},
					},
					CreatedAt: timestamppb.New(time.Time{}),
				},
			},
		},
		{
			name: "return error if service list return error",
			setup: func(ss *mocks.SilenceService) {
				ss.EXPECT().List(mock.AnythingOfType("*context.emptyCtx"), mock.AnythingOfType("silence.Filter")).Return(nil, errors.New("some error"))
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		ctx := context.TODO()
		t.Run(tt.name, func(t *testing.T) {
			mockSilenceService := new(mocks.SilenceService)

			if tt.setup != nil {
				tt.setup(mockSilenceService)
			}

			s := v1beta1.NewGRPCServer(nil, log.NewNoop(), api.HeadersConfig{}, &api.Deps{SilenceService: mockSilenceService})
			got, err := s.ListSilences(ctx, &sirenv1beta1.ListSilencesRequest{})
			if (err != nil) != tt.wantErr {
				t.Errorf("GRPCServer.ListSilences() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if diff := cmp.Diff(got.GetSilences(), tt.want, protocmp.Transform(), protocmp.IgnoreFields(&sirenv1beta1.Silence{}, "id", "deleted_at")); diff != "" {
				t.Errorf("GRPCServer.ListSilences() diff = %v", diff)
			}
		})
	}
}

func TestGRPCServer_GetSilence(t *testing.T) {
	mockSilenceData := silence.Silence{
		ID:          "silence-id",
		NamespaceID: 1,
		Type:        silence.TypeMatchers,
		TargetExpression: map[string]any{
			"key1": "value1",
		},
	}

	tests := []struct {
		name    string
		setup   func(*mocks.SilenceService)
		req     *sirenv1beta1.GetSilenceRequest
		want    *sirenv1beta1.Silence
		wantErr bool
	}{
		{
			name: "return silence when successfully get silence",
			setup: func(ss *mocks.SilenceService) {
				ss.EXPECT().Get(mock.AnythingOfType("*context.emptyCtx"), mockSilenceData.ID).Return(mockSilenceData, nil)
			},
			req: &sirenv1beta1.GetSilenceRequest{
				Id: mockSilenceData.ID,
			},
			want: &sirenv1beta1.Silence{
				NamespaceId: mockSilenceData.NamespaceID,
				Type:        mockSilenceData.Type,
				TargetExpression: &structpb.Struct{
					Fields: map[string]*structpb.Value{
						"key1": structpb.NewStringValue("value1"),
					},
				},
				CreatedAt: timestamppb.New(time.Time{}),
			},
		},
		{
			name: "return error if service get return error",
			setup: func(ss *mocks.SilenceService) {
				ss.EXPECT().Get(mock.AnythingOfType("*context.emptyCtx"), mockSilenceData.ID).Return(silence.Silence{}, errors.New("some error"))
			},
			req: &sirenv1beta1.GetSilenceRequest{
				Id: mockSilenceData.ID,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		ctx := context.TODO()
		t.Run(tt.name, func(t *testing.T) {
			mockSilenceService := new(mocks.SilenceService)

			if tt.setup != nil {
				tt.setup(mockSilenceService)
			}

			s := v1beta1.NewGRPCServer(nil, log.NewNoop(), api.HeadersConfig{}, &api.Deps{SilenceService: mockSilenceService})
			got, err := s.GetSilence(ctx, tt.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("GRPCServer.GetSilence() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if diff := cmp.Diff(got.GetSilence(), tt.want, protocmp.Transform(), protocmp.IgnoreFields(&sirenv1beta1.Silence{}, "id", "deleted_at")); diff != "" {
				t.Errorf("GRPCServer.GetSilence() diff = %v", diff)
			}
		})
	}
}

func TestGRPCServer_ExpireSilence(t *testing.T) {
	mockSilenceData := silence.Silence{
		ID:          "silence-id",
		NamespaceID: 1,
		Type:        silence.TypeMatchers,
		TargetExpression: map[string]any{
			"key1": "value1",
		},
	}

	tests := []struct {
		name    string
		setup   func(*mocks.SilenceService)
		req     *sirenv1beta1.ExpireSilenceRequest
		want    *sirenv1beta1.ExpireSilenceResponse
		wantErr bool
	}{
		{
			name: "return success when successfully deleted silence",
			setup: func(ss *mocks.SilenceService) {
				ss.EXPECT().Delete(mock.AnythingOfType("*context.emptyCtx"), mockSilenceData.ID).Return(nil)
			},
			req: &sirenv1beta1.ExpireSilenceRequest{
				Id: mockSilenceData.ID,
			},
			want: &sirenv1beta1.ExpireSilenceResponse{},
		},
		{
			name: "return error if service delete return error",
			setup: func(ss *mocks.SilenceService) {
				ss.EXPECT().Delete(mock.AnythingOfType("*context.emptyCtx"), mockSilenceData.ID).Return(errors.New("some error"))
			},
			req: &sirenv1beta1.ExpireSilenceRequest{
				Id: mockSilenceData.ID,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		ctx := context.TODO()
		t.Run(tt.name, func(t *testing.T) {
			mockSilenceService := new(mocks.SilenceService)

			if tt.setup != nil {
				tt.setup(mockSilenceService)
			}

			s := v1beta1.NewGRPCServer(nil, log.NewNoop(), api.HeadersConfig{}, &api.Deps{SilenceService: mockSilenceService})
			got, err := s.ExpireSilence(ctx, tt.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("GRPCServer.ExpireSilence() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GRPCServer.ExpireSilence() = %v, want %v", got, tt.want)
			}
		})
	}
}
