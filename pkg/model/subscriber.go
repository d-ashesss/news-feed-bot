package model

type Subscriber struct {
	ID         string
	UserID     string
	Categories []Category
}

func NewSubscriber(userID string) *Subscriber {
	return &Subscriber{UserID: userID}
}

func (s *Subscriber) AddCategory(c Category) {
	s.Categories = append(s.Categories, c)
}

func (s *Subscriber) RemoveCategory(c Category) {
	subs := make([]Category, 0, len(s.Categories))
	for _, cat := range s.Categories {
		if cat.ID != c.ID {
			subs = append(subs, cat)
		}
	}
	s.Categories = subs
}
