package services

import "context"

func (s *Service) UserExistsByEmail(ctx context.Context, email string) (bool, error) {
	return s.repository.UserExistsByEmail(ctx, email)
}
