// Package scaffold contains methods and helpers for code scaffolding
package scaffold

// Option sets options for Scaffold instance
type Option func(s *Scaffold)

// WithOutputLocation allows to specify location where generated files will be located.
// If empty string is provided - then files will be created in current directory.
func WithOutputLocation(location string) Option {
	return func(s *Scaffold) {
		s.location = location
	}
}
