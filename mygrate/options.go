package mygrate

type Option func(s *Service)

// WithStore will set a custom Store implementation.
func WithStore(store Store) Option {
	return func(s *Service) {
		s.store = store
	}
}
