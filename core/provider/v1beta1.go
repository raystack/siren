package provider

import (
	"fmt"

	sirenv1beta1 "github.com/goto/siren/proto/gotocompany/siren/v1beta1"
	"google.golang.org/protobuf/types/known/structpb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (p *Provider) ToV1beta1Proto() (*sirenv1beta1.Provider, error) {
	credentials, err := structpb.NewStruct(p.Credentials)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch provider credentials: %w", err)
	}

	return &sirenv1beta1.Provider{
		Id:          p.ID,
		Urn:         p.URN,
		Host:        p.Host,
		Type:        p.Type,
		Name:        p.Name,
		Credentials: credentials,
		Labels:      p.Labels,
		CreatedAt:   timestamppb.New(p.CreatedAt),
		UpdatedAt:   timestamppb.New(p.UpdatedAt),
	}, nil
}

func FromV1beta1Proto(proto *sirenv1beta1.Provider) *Provider {
	return &Provider{
		ID:          proto.GetId(),
		Host:        proto.GetHost(),
		URN:         proto.GetUrn(),
		Name:        proto.GetName(),
		Type:        proto.GetType(),
		Credentials: proto.GetCredentials().AsMap(),
		Labels:      proto.GetLabels(),
		CreatedAt:   proto.CreatedAt.AsTime(),
		UpdatedAt:   proto.UpdatedAt.AsTime(),
	}
}
