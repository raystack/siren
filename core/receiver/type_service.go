package receiver

//go:generate mockery --name=TypeService -r --case underscore --with-expecter --structname TypeService --filename receiver_type_service.go --output=./mocks
type TypeService interface {
	Encrypt(r *Receiver) error
	Decrypt(r *Receiver) error
	PopulateReceiver(rcv *Receiver) (*Receiver, error)
	Notify(rcv *Receiver, payloadMessage NotificationMessage) error
	ValidateConfiguration(configurations Configurations) error
}
