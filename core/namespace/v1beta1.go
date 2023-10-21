package namespace

import (
	"fmt"

	"github.com/goto/siren/core/provider"
	sirenv1beta1 "github.com/goto/siren/proto/gotocompany/siren/v1beta1"
	"google.golang.org/protobuf/types/known/structpb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (n *Namespace) ToV1beta1Proto() (*sirenv1beta1.Namespace, error) {
	creds, err := structpb.NewStruct(n.Credentials)
	if err != nil {
		return nil, fmt.Errorf("cannot transform credentials to proto: %w", err)
	}
	return &sirenv1beta1.Namespace{
		Id:          n.ID,
		Urn:         n.URN,
		Name:        n.Name,
		Provider:    n.Provider.ID,
		Credentials: creds,
		Labels:      n.Labels,
		CreatedAt:   timestamppb.New(n.CreatedAt),
		UpdatedAt:   timestamppb.New(n.UpdatedAt),
	}, nil
}

func FromV1beta1Proto(proto *sirenv1beta1.Namespace) Namespace {
	return Namespace{
		ID:   proto.GetId(),
		URN:  proto.GetUrn(),
		Name: proto.GetName(),
		Provider: provider.Provider{
			ID: proto.GetProvider(),
		},
		Credentials: proto.GetCredentials().AsMap(),
		Labels:      proto.GetLabels(),
		CreatedAt:   proto.GetCreatedAt().AsTime(),
		UpdatedAt:   proto.GetUpdatedAt().AsTime(),
	}
}
