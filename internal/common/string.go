package common

type String struct {
	*string
}

func NewString(v string) String {
	return String{string: &v}
}

func (j String) PtrString() *string {
	return j.string
}
