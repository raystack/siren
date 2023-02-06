package v1beta1

import (
	"context"

	"github.com/odpf/siren/core/silence"
	sirenv1beta1 "github.com/odpf/siren/proto/odpf/siren/v1beta1"
	"google.golang.org/protobuf/types/known/structpb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (s *GRPCServer) CreateSilence(ctx context.Context, req *sirenv1beta1.CreateSilenceRequest) (*sirenv1beta1.CreateSilenceResponse, error) {
	id, err := s.silenceService.Create(ctx, silence.Silence{
		NamespaceID:      req.GetNamespaceId(),
		Type:             req.GetType(),
		TargetID:         req.GetTargetId(),
		TargetExpression: req.GetTargetExpression().AsMap(),
	})
	if err != nil {
		return nil, s.generateRPCErr(err)
	}

	return &sirenv1beta1.CreateSilenceResponse{
		Id: id,
	}, nil
}

func (s *GRPCServer) ListSilences(ctx context.Context, req *sirenv1beta1.ListSilencesRequest) (*sirenv1beta1.ListSilencesResponse, error) {
	silences, err := s.silenceService.List(ctx, silence.Filter{
		NamespaceID:       req.GetNamespaceId(),
		SubscriptionID:    req.GetSubscriptionId(),
		Match:             req.GetMatch(),
		SubscriptionMatch: req.GetSubscriptionMatch(),
	})
	if err != nil {
		return nil, s.generateRPCErr(err)
	}

	var silencesProto []*sirenv1beta1.Silence
	for _, si := range silences {
		targetExpression, err := structpb.NewStruct(si.TargetExpression)
		if err != nil {
			return nil, s.generateRPCErr(err)
		}

		silencesProto = append(silencesProto, &sirenv1beta1.Silence{
			Id:               si.ID,
			NamespaceId:      si.NamespaceID,
			Type:             si.Type,
			TargetId:         si.TargetID,
			TargetExpression: targetExpression,
			CreatedAt:        timestamppb.New(si.CreatedAt),
			DeletedAt:        timestamppb.New(si.DeletedAt),
		})
	}

	return &sirenv1beta1.ListSilencesResponse{
		Silences: silencesProto,
	}, nil
}

func (s *GRPCServer) GetSilence(ctx context.Context, req *sirenv1beta1.GetSilenceRequest) (*sirenv1beta1.GetSilenceResponse, error) {
	sil, err := s.silenceService.Get(ctx, req.GetId())
	if err != nil {
		return nil, s.generateRPCErr(err)
	}

	targetExpression, err := structpb.NewStruct(sil.TargetExpression)
	if err != nil {
		return nil, s.generateRPCErr(err)
	}

	return &sirenv1beta1.GetSilenceResponse{
		Silence: &sirenv1beta1.Silence{
			Id:               sil.ID,
			NamespaceId:      sil.NamespaceID,
			Type:             sil.Type,
			TargetId:         sil.TargetID,
			TargetExpression: targetExpression,
			CreatedAt:        timestamppb.New(sil.CreatedAt),
			DeletedAt:        timestamppb.New(sil.DeletedAt),
		},
	}, nil
}

func (s *GRPCServer) ExpireSilence(ctx context.Context, req *sirenv1beta1.ExpireSilenceRequest) (*sirenv1beta1.ExpireSilenceResponse, error) {
	if err := s.silenceService.Delete(ctx, req.GetId()); err != nil {
		return nil, s.generateRPCErr(err)
	}

	return &sirenv1beta1.ExpireSilenceResponse{}, nil
}
