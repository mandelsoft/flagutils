package output

type FieldNameProvider interface {
	GetFieldNames() []string
}
