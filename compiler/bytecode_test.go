package compiler_test

import (
	"bytes"
	"testing"

	"github.com/d5/tengo/assert"
	"github.com/d5/tengo/compiler"
	"github.com/d5/tengo/objects"
)

func TestBytecode(t *testing.T) {
	testBytecodeSerialization(t, &compiler.Bytecode{})

	testBytecodeSerialization(t, bytecode(
		concat(), objectsArray(
			&objects.Array{
				Value: objectsArray(
					&objects.Int{Value: 12},
					&objects.String{Value: "foo"},
					&objects.Bool{Value: true},
					&objects.Float{Value: 93.11},
					&objects.Char{Value: 'x'},
				),
			},
			&objects.Bool{Value: false},
			&objects.Char{Value: 'y'},
			&objects.Float{Value: 93.11},
			compiledFunction(1, 0,
				compiler.MakeInstruction(compiler.OpConstant, 3),
				compiler.MakeInstruction(compiler.OpSetLocal, 0),
				compiler.MakeInstruction(compiler.OpGetGlobal, 0),
				compiler.MakeInstruction(compiler.OpGetFree, 0)),
			&objects.Float{Value: 39.2},
			&objects.Int{Value: 192},
			&objects.Map{
				Value: map[string]objects.Object{
					"a": &objects.Float{Value: -93.1},
					"b": &objects.Bool{Value: false},
				},
			},
			&objects.String{Value: "bar"},
			&objects.Undefined{})))

	testBytecodeSerialization(t, bytecode(
		concat(
			compiler.MakeInstruction(compiler.OpConstant, 0),
			compiler.MakeInstruction(compiler.OpSetGlobal, 0),
			compiler.MakeInstruction(compiler.OpConstant, 6),
			compiler.MakeInstruction(compiler.OpPop)),
		objectsArray(
			intObject(55),
			intObject(66),
			intObject(77),
			intObject(88),
			compiledFunction(1, 0,
				compiler.MakeInstruction(compiler.OpConstant, 3),
				compiler.MakeInstruction(compiler.OpSetLocal, 0),
				compiler.MakeInstruction(compiler.OpGetGlobal, 0),
				compiler.MakeInstruction(compiler.OpGetFree, 0),
				compiler.MakeInstruction(compiler.OpAdd),
				compiler.MakeInstruction(compiler.OpGetFree, 1),
				compiler.MakeInstruction(compiler.OpAdd),
				compiler.MakeInstruction(compiler.OpGetLocal, 0),
				compiler.MakeInstruction(compiler.OpAdd),
				compiler.MakeInstruction(compiler.OpReturnValue, 1)),
			compiledFunction(1, 0,
				compiler.MakeInstruction(compiler.OpConstant, 2),
				compiler.MakeInstruction(compiler.OpSetLocal, 0),
				compiler.MakeInstruction(compiler.OpGetFree, 0),
				compiler.MakeInstruction(compiler.OpGetLocal, 0),
				compiler.MakeInstruction(compiler.OpClosure, 4, 2),
				compiler.MakeInstruction(compiler.OpReturnValue, 1)),
			compiledFunction(1, 0,
				compiler.MakeInstruction(compiler.OpConstant, 1),
				compiler.MakeInstruction(compiler.OpSetLocal, 0),
				compiler.MakeInstruction(compiler.OpGetLocal, 0),
				compiler.MakeInstruction(compiler.OpClosure, 5, 1),
				compiler.MakeInstruction(compiler.OpReturnValue, 1)))))

	bytecode, _, err := traceCompile(`
/*
    Tengo Language
*/

// variable definition and primitive types
a := "foo"   // string
b := -19.84  // floating point
c := 5       // integer
d := true    // boolean
e := '九'    // char

// assignment
b = "bar"    // can assign value of different type

// map and array
m := {a: {b: {c: [1, 2, 3]}}}
two := m.a.b.c[1] == m["a"]["b"]["c"][1]
m.d = "dee"
m.a.b.e = "eee"

// slicing
str := "hello world"
substr := str[1:5]    // "ello"
arr := [1, 2, 3, 4, 5]
subarr := arr[2:4]    // [3, 4]

// functions
each := func(seq, fn) {
    // array iteration
    for x in seq {  
        fn(x) 
    } 
}

sum := func(seq) {
   s := 0
   each(seq, func(x) { 
       s += x    // closure: capturing variable 's'
   })
   return s
}

six := sum([1, 2, 3]) // array: [1, 2, 3]

map_to_array := func(m) {
    arr := []
    // map iteration
    for key, value in m { 
        arr = append(arr, key, value)  // builtin function 'append'
    }
    return arr
}

m_arr := map_to_array(m)
m_arr_len := len(m_arr)   // builtin function 'len'

// tail-call optimization: faster and enables loop via recursion
count_odds := func(n, c) {
	if n == 0 {
		return c
	} else if n % 2 == 1 {
	    c++
	}
	return count_odds(n-1, c)
}
num_odds := count_odds(100000, 0)

// type coercion
s1 := string(1984)  // "1984"
i2 := int("-999")   // -999
f3 := float(-51)    // -51.0
b4 := bool(1)       // true
c5 := char("X")     // 'X'


// if statement
three_is := ""
if three := 3; three > 2 {  // optional init statement
    three_is = "> 2"
} else if three == 2 {
    three_is = "= 2"
} else {
    three_is = "< 2"
}

// for statement
seven := 0
arr2 := [1, 2, 3, 1]
for i:=0; i<len(arr2); i++ {
    seven += arr[1]
}
`, nil)
	assert.NoError(t, err)
	testBytecodeSerialization(t, bytecode)
}

func testBytecodeSerialization(t *testing.T, b *compiler.Bytecode) {
	var buf bytes.Buffer
	err := b.Encode(&buf)
	assert.NoError(t, err)

	r := &compiler.Bytecode{}
	err = r.Decode(bytes.NewReader(buf.Bytes()))
	assert.NoError(t, err)

	assert.Equal(t, b.Instructions, r.Instructions)
	assert.Equal(t, b.Constants, r.Constants)
}
