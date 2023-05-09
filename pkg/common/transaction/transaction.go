package transaction

type Service struct {
}

func New() *Service {
	return &Service{}
}

func (s *Service) Send() error {

	return nil
}
