package query

type Validator interface {
    Validate(interface{}) error
}

type StructValidator interface {
    Validator
    RegisterValidation(tag string, fn func(fl FieldLevel) bool) error
    RegisterStructValidation(fn func(sl StructLevel), types ...interface{})
}

type FieldLevel interface {
    Field() reflect.Value
    Param() string
}

type StructLevel interface {
    Current() interface{}
    Parent() interface{}
    Field(name string) reflect.Value
} 