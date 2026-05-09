package click

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) Create(click *Click) error {
	return s.repo.Create(click)
}

func (s *Service) CountByLink(shortLinkID int) (int, error) {
	return s.repo.CountByLink(shortLinkID)
}
