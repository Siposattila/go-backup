package server

func (s *server) alertSystem(message string) {
	if s.Discord != nil {
		s.Discord.Send(message)
	}

	if s.Email != nil {
		s.Email.Send(message)
	}
}
