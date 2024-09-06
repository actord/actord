package typedef

type Collection []Type

func (c Collection) Find(name string) *Type {
	for i := range c {
		if c[i].Name == name {
			return &c[i]
		}
	}
	return nil
}
