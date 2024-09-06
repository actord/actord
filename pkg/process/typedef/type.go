package typedef

type Type struct {
	Name   string  `hcl:"name,label"`
	Any    bool    `hcl:"any,optional"`
	Fields []Field `hcl:"field,block"`
}

type Field struct {
	Name        string  `hcl:"name,label"`
	Type        string  `hcl:"type"`
	Label       string  `hcl:"label"`
	Optional    *bool   `hcl:"optional"`
	Description *string `hcl:"description"`
}
