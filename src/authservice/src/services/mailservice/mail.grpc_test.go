package mailservice

import "context"

func NewMailService() (Service, error) {
	client, err := ConnectClient(context.Background(), MailServiceAddr("localhost:50052"))
	if err != nil {
		return nil, err
	}
	return ServiceGRPC{
		Ctx:    context.Background(),
		Client: client,
	}, err
}
