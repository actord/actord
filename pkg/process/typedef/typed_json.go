package typedef

import (
	"github.com/valyala/fastjson"
)

// TypedJSON is a json data with a type
type TypedJSON struct {
	Type  Type
	State *fastjson.Value
}

func NewTypedJSON(t Type, state *fastjson.Value) *TypedJSON {
	return &TypedJSON{
		Type:  t,
		State: state,
	}
}

func NewTypedJSONFromBytes(t Type, data []byte) (*TypedJSON, error) {
	state, err := fastjson.ParseBytes(data)
	if err != nil {
		return nil, err
	}
	tj := NewTypedJSON(t, state)
	return tj, tj.Validate()
}

func (t *TypedJSON) Validate() error {
	if t.Type.Any {
		// type.Any means that the type is not specified
		return nil
	}

	// TODO: validate json data according to the type

	return nil
}

func (t *TypedJSON) Marshal() []byte {
	return t.State.MarshalTo(nil)
}

func (t *TypedJSON) Get(key string) (*fastjson.Value, error) {
	if key == "" {
		return t.State, nil
	}
	// TODO: validate key according to the type

	// TODO: maybe validate value according to the type

	return t.State.Get(key), nil
}

func (t *TypedJSON) Set(key string, value *fastjson.Value) error {
	// TODO: validate key/value according to the type
	t.State.Set(key, value)
	return nil
}
