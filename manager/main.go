package main

func main() {
	// initialize a new server
	s := NewServer()
	r := s.Routes()

	// run the server
	s.Run(r)
}
