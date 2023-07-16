package server

// routes registers routes and middleware.
func (s server) routes() {
	s.router.Handle("/books", s.mw.Authenticator.Authenticate(s.booksCreateHandler()))
}
