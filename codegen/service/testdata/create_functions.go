package testdata

var CreateStringCode = `// CreateFromStringT initializes t from the fields of v
func (t *StringType) CreateFromStringT(v *testdata.StringT) {
	temp := &StringType{
		String: &v.String,
	}
	*t = *temp
}
`

var CreateStringRequiredCode = `// CreateFromStringT initializes t from the fields of v
func (t *StringType) CreateFromStringT(v *testdata.StringT) {
	temp := &StringType{
		String: v.String,
	}
	*t = *temp
}
`

var CreateStringPointerCode = `// CreateFromStringPointerT initializes t from the fields of v
func (t *StringType) CreateFromStringPointerT(v *testdata.StringPointerT) {
	temp := &StringType{
		String: v.String,
	}
	*t = *temp
}
`

var CreateStringPointerRequiredCode = `// CreateFromStringPointerT initializes t from the fields of v
func (t *StringType) CreateFromStringPointerT(v *testdata.StringPointerT) {
	temp := &StringType{}
	if v.String != nil {
		temp.String = *v.String
	}
	*t = *temp
}
`

var CreateArrayStringCode = `// CreateFromArrayStringT initializes t from the fields of v
func (t *ArrayStringType) CreateFromArrayStringT(v *testdata.ArrayStringT) {
	temp := &ArrayStringType{}
	if v.ArrayString != nil {
		temp.ArrayString = make([]string, len(v.ArrayString))
		for j, val := range v.ArrayString {
			temp.ArrayString[j] = val
		}
	}
	*t = *temp
}
`

var CreateArrayStringRequiredCode = `// CreateFromArrayStringT initializes t from the fields of v
func (t *ArrayStringType) CreateFromArrayStringT(v *testdata.ArrayStringT) {
	temp := &ArrayStringType{}
	if v.ArrayString != nil {
		temp.ArrayString = make([]string, len(v.ArrayString))
		for j, val := range v.ArrayString {
			temp.ArrayString[j] = val
		}
	}
	*t = *temp
}
`

var CreateObjectCode = `// CreateFromObjectT initializes t from the fields of v
func (t *ObjectType) CreateFromObjectT(v *testdata.ObjectT) {
	temp := &ObjectType{}
	if v.Object != nil {
		temp.Object = marshalObjectFieldTToObjectField(v.Object)
	}
	*t = *temp
}
`

var CreateObjectRequiredCode = `// CreateFromObjectT initializes t from the fields of v
func (t *ObjectType) CreateFromObjectT(v *testdata.ObjectT) {
	temp := &ObjectType{}
	if v.Object != nil {
		temp.Object = marshalObjectFieldTToObjectField(v.Object)
	}
	*t = *temp
}
`

var CreateObjectExtraCode = `// CreateFromObjectExtraT initializes t from the fields of v
func (t *ObjectType) CreateFromObjectExtraT(v *testdata.ObjectExtraT) {
	temp := &ObjectType{}
	if v.Object != nil {
		temp.Object = marshalObjectFieldTToObjectField(v.Object)
	}
	*t = *temp
}
`