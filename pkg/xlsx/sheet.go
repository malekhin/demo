package xlsx

type Sheet struct {
	name string
	data [][]string
}

func (s *Sheet) AddTitle(data []string) {
	s.AddCell(data)
}

func (s *Sheet) AddCell(data []string) {
	s.data = append(s.data, data)
}
