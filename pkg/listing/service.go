package listing

type RepositoryMySQL interface {
	ReadUser(User) (User, error)
}

type Service interface {
	GetUser(User) (User, error)
}

type service struct {
	rmy RepositoryMySQL
}

func NewService(rmy RepositoryMySQL) Service {
	return &service{rmy}
}

func (s *service) GetUser(lu User) (User, error) {
	var err error
	lu, err = s.rmy.ReadUser(lu)
	if err != nil {
		return lu, err
	}

	return lu, nil
}