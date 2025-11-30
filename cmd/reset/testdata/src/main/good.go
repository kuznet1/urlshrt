package main

// generate:reset
type Test struct {
	I      int
	F32    float32
	Str    string
	Flag   bool
	Custom SomeType

	PI      *int
	PStr    *string
	PBool   *bool
	PCustom *SomeType

	Sl  []int
	Arr [5]int
	M   map[string]int

	Any     interface{}
	SomeStr *struct{}
}

type SomeType struct{}

func (s *SomeType) Reset() {}
