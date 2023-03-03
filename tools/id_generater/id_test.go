package id_generater

import (
	"fmt"
	"testing"
)

func Test_ID(t *testing.T) {
	idGenerate, _ := New( 0)
	id := idGenerate.Generate()

	fmt.Println("snowflake id int64:", id.Int64())
	fmt.Println("snowflake id string:", id.String())
}

func Test_tw_id(t *testing.T) {
	s, _ := NewTwitterSF(0)
	s.Generate()
}
