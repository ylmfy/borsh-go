package borsh

import (
	"reflect"
	strings2 "strings"
	"testing"
)

type A struct {
	A int
	B int32
}

func TestSimple(t *testing.T) {
	x := A{
		A: 1,
		B: 32,
	}
	data, err := Serialize(x)
	if err != nil {
		t.Error(err)
	}
	y := new(A)
	err = Deserialize(y, data)
	if err != nil {
		t.Error(err)
	}
	if !reflect.DeepEqual(x, *y) {
		t.Error(x, y)
	}
}

type B struct {
	I8  int8
	I16 int16
	I32 int32
	I64 int64
	U8  uint8
	U16 uint16
	U32 uint32
	U64 uint64
	F32 float32
	F64 float64
}

func TestBasic(t *testing.T) {
	x := B{
		I8:  12,
		I16: -1,
		I32: 124,
		I64: 1243,
		U8:  1,
		U16: 979,
		U32: 123124,
		U64: 1135351135,
		F32: -231.23,
		F64: 3121221.232,
	}
	data, err := Serialize(x)
	if err != nil {
		t.Error(err)
	}
	y := new(B)
	err = Deserialize(y, data)
	if err != nil {
		t.Error(err)
	}
	if !reflect.DeepEqual(x, *y) {
		t.Error(x, y)
	}
}

type C struct {
	A3 [3]int
	S  []int
	P  *int
	M  map[string]string
}

func TestBasicContainer(t *testing.T) {
	ip := new(int)
	*ip = 213
	x := C{
		A3: [3]int{234, -123, 123},
		S:  []int{21442, 421241241, 2424},
		P:  ip,
		M:  make(map[string]string),
	}
	data, err := Serialize(x)
	if err != nil {
		t.Error(err)
	}
	y := new(C)
	err = Deserialize(y, data)
	if err != nil {
		t.Error(err)
	}
	if !reflect.DeepEqual(x, *y) {
		t.Error(x, y)
	}
}

type N struct {
	B B
	C C
}

func TestNested(t *testing.T) {
	ip := new(int)
	*ip = 213
	x := N{
		B: B{
			I8:  12,
			I16: -1,
			I32: 124,
			I64: 1243,
			U8:  1,
			U16: 979,
			U32: 123124,
			U64: 1135351135,
			F32: -231.23,
			F64: 3121221.232,
		},
		C: C{
			A3: [3]int{234, -123, 123},
			S:  []int{21442, 421241241, 2424},
			P:  ip,
			M:  make(map[string]string),
		},
	}
	data, err := Serialize(x)
	if err != nil {
		t.Error(err)
	}
	y := new(N)
	err = Deserialize(y, data)
	if err != nil {
		t.Error(err)
	}
	if !reflect.DeepEqual(x, *y) {
		t.Error(x, y)
	}
}

type Dummy Enum

const (
	x Dummy = iota
	y
	z
)

type D struct {
	D Dummy
}

func TestSimpleEnum(t *testing.T) {
	x := D{
		D: y,
	}
	data, err := Serialize(x)
	if err != nil {
		t.Error(err)
	}
	y := new(D)
	err = Deserialize(y, data)
	if err != nil {
		t.Error(err)
	}
	if !reflect.DeepEqual(x, *y) {
		t.Error(x, y)
	}
}

type ComplexEnum struct {
	Enum Enum `borsh_enum:"true"`
	Foo  Foo
	Bar  Bar
}

type Foo struct {
	FooA int32
	FooB string
}

type Bar struct {
	BarA int64
	BarB string
}

func TestComplexEnum(t *testing.T) {
	x := ComplexEnum{
		Enum: 0,
		Foo: Foo{
			FooA: 23,
			FooB: "baz",
		},
	}
	data, err := Serialize(x)
	if err != nil {
		t.Fatal(err)
	}
	y := new(ComplexEnum)
	err = Deserialize(y, data)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(x, *y) {
		t.Fatal(x, y)
	}
}

type S struct {
	S map[int]struct{}
}

func TestSet(t *testing.T) {
	x := S{
		S: map[int]struct{}{124: struct{}{}, 214: struct{}{}, 24: struct{}{}, 53: struct{}{}},
	}
	data, err := Serialize(x)
	if err != nil {
		t.Error(err)
	}
	y := new(S)
	err = Deserialize(y, data)
	if err != nil {
		t.Error(err)
	}
	if !reflect.DeepEqual(x, *y) {
		t.Error(x, y)
	}
}

type Skipped struct {
	A int
	B int `borsh_skip:"true"`
	C int
}

func TestSkipped(t *testing.T) {
	x := Skipped{
		A: 32,
		B: 535,
		C: 123,
	}
	data, err := Serialize(x)
	if err != nil {
		t.Error(err)
	}
	y := new(Skipped)
	err = Deserialize(y, data)
	if err != nil {
		t.Error(err)
	}
	if x.A != y.A || x.C != y.C {
		t.Errorf("%v fields not equal to %v", x, y)
	}
	if y.B == x.B {
		t.Errorf("didn't skip field B")
	}
}

type E struct{}

func TestEmpty(t *testing.T) {
	x := E{}
	data, err := Serialize(x)
	if err != nil {
		t.Error(err)
	}
	if len(data) != 0 {
		t.Error("not empty")
	}
	y := new(E)
	err = Deserialize(y, data)
	if err != nil {
		t.Error(err)
	}
	if !reflect.DeepEqual(x, *y) {
		t.Error(x, y)
	}
}

func testValue(t *testing.T, v interface{}) {
	data, err := Serialize(v)
	if err != nil {
		t.Error(err)
	}
	parsed := reflect.New(reflect.TypeOf(v))
	err = Deserialize(parsed.Interface(), data)
	if err != nil {
		t.Error(err)
	}
	reflect.DeepEqual(v, parsed.Elem().Interface())
}

func TestStrings(t *testing.T) {
	tests := []struct {
		in string
	}{
		{""},
		{"a"},
		{"hellow world"},
		{strings2.Repeat("x", 1024)},
		{strings2.Repeat("x", 4096)},
		{strings2.Repeat("x", 65535)},
		{strings2.Repeat("hello world!", 1000)},
		{"💩"},
	}

	for _, tt := range tests {
		testValue(t, tt.in)
	}
}

func makeInt32Slice(val int32, len int) []int32 {
	s := make([]int32, len)
	for i := 0; i < len; i++ {
		s[i] = val
	}
	return s
}

func TestSlices(t *testing.T) {
	tests := []struct {
		in []int32
	}{
		{makeInt32Slice(1000000000, 0)},
		{makeInt32Slice(1000000000, 1)},
		{makeInt32Slice(1000000000, 2)},
		{makeInt32Slice(1000000000, 3)},
		{makeInt32Slice(1000000000, 4)},
		{makeInt32Slice(1000000000, 8)},
		{makeInt32Slice(1000000000, 16)},
		{makeInt32Slice(1000000000, 32)},
		{makeInt32Slice(1000000000, 64)},
		{makeInt32Slice(1000000000, 65)},
	}

	for _, tt := range tests {
		testValue(t, tt.in)
	}
}
