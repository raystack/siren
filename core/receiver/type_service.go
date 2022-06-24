package receiver

import "context"

//go:generate mockery --name=TypeService -r --case underscore --with-expecter --structname TypeService --filename receiver_type_service.go --output=./mocks
type TypeService interface {
	Encrypt(r *Receiver) error
	Decrypt(r *Receiver) error
	PopulateReceiver(ctx context.Context, rcv *Receiver) (*Receiver, error)
	Notify(ctx context.Context, rcv *Receiver, payloadMessage NotificationMessage) error
	ValidateConfiguration(rcv *Receiver) error
	GetSubscriptionConfig(subsConfs map[string]string, receiverConfs Configurations) (map[string]string, error)
}
