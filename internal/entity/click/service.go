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

func (s *Service) GetStats(
	linkID int,
) (*LinkStats, error) {
	return s.repo.GetStats(linkID)
}
