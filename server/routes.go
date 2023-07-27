package server

// routes registers routes and middleware.
func (s server) routes() {
	s.router.Handle("/books", s.bookHandler())
    s.route.Handle("/healthz", s.healthzHandler())
}
