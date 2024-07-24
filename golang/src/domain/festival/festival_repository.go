package festival

type FestivalRepository interface {
	FindUnPosted() ([]Festival, error)
	FindUnQuoted() ([]Festival, error)
	Save(festival Festival) error
}
