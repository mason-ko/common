package converter

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"
)

func TestConvert25(t *testing.T) {
	type a struct {
		A string
		B string
		C string
	}

	type b struct {
		A string
		B string
		C string
		D string
	}

	a1 := a{A: "1", B: "2", C: "3"}
	b1 := b{}

	toCopyOnOtherType(&a1, &b1)

	fmt.Println(a1)
	fmt.Println(b1)
}

func TestConvert_Arr(t *testing.T) {
	type dbModel struct {
		A string
	}

	type pbModel struct {
		A string
	}

	var db []*dbModel
	db = nil
	var pb []*pbModel

	toCopyOnOtherType(db, &pb)

	fmt.Println(db)

	fmt.Println("=====================")

	fmt.Println(pb)

	type response struct {
		Code int
		Data interface{}
	}
	da := response{
		Code: 222,
		Data: pb,
	}
	b, _ := json.Marshal(da)
	fmt.Println(string(b))
}

func TestConvert_Split(t *testing.T) {
	type dbModel struct {
		A string
		B string
		C string
		D string `seperator:","`
		E []string
	}

	type pbModel struct {
		A string
		B string
		C time.Time
		D []string
		E []string
	}

	db := dbModel{A: "ABCD", B: "CCCC", C: "2019-01-01 00:00:00", D: "1,2,3,4,5"}
	pb := pbModel{}

	toCopyOnOtherType(&db, &pb)

	fmt.Println(db)

	fmt.Println("=====================")

	fmt.Println(pb)
}

func TestConvert_Map(t *testing.T) {
	type dbModel3 struct {
		A string
		B string
		C string
	}

	type dbModel2 struct {
		A string
		B map[string]string
		C map[string]dbModel3
	}

	type dbModel struct {
		A string
		B string
		C map[string]dbModel2
	}

	db := map[string]dbModel{
		"zz": {
			A: "a",
			B: "b",
			C: map[string]dbModel2{
				"wow": {
					A: "T",
					B: map[string]string{"z": "X"},
					C: map[string]dbModel3{
						"wow": {
							A: "zT",
							B: "zT2",
							C: "zT3",
						},
					},
				},
			},
		},
	}
	pb := map[string]dbModel{}

	ToCopyObject(&db, &pb)

	fmt.Println(db)

	fmt.Println("=====================")

	fmt.Println(pb)

	//for _, v := range pb {
	//	fmt.Println(v)
	//}
}

func TestConvert235(t *testing.T) {
	type dbModel struct {
		B  *string
		C  *time.Time
		D  int64
		E  *int64
		G  *int64
		A1 string
		A2 string
		A3 *string
		A4 []string
		A5 int
		A6 int
		A7 *int32
	}
	type pbModel struct {
		B  *string
		C  time.Time
		D  *int64
		E  int64
		G  *int64
		A1 []string
		A2 *string
		A3 string
		A4 string
		A5 *int64
		A6 int32
		A8 int
	}

	v := "zz"
	tt := time.Now()
	var e int64 = 1
	db := dbModel{
		B:  &v,
		C:  &tt,
		D:  0,
		E:  &e,
		G:  nil,
		A1: "z",
		A2: "x",
		A3: nil,
		A4: []string{"aa", "bb", "cc"},
		A5: 19,
		A6: 91,
	}
	var pb pbModel

	ToCopyObject(&db, &pb)

	b, _ := json.Marshal(db)
	fmt.Println(string(b))

	fmt.Println("=====================")

	b, _ = json.Marshal(pb)
	fmt.Println(string(b))
}


func TestConvertPBBBB2(t *testing.T) {

	type C struct {
		A string
		B string
		C string
	}
	type B struct {
		A string
		B *C
		C C
	}

	type A struct {
		//A string
		//B string
		//C *B
		//D B
		A2 []*B
		A3 []*B
		E  *[]*B

		E3 []B
		G  *B
	}

	type PBModel struct {
		AA *A
	}

	type DBModel struct {
		AA *A
	}

	//pb := PBModel{AA: &A{
	//	A:"z1",B:"bb2",C: &B{
	//		A:"zz",B:"WOW~!#@!",
	//	}, E: []*B{{A:"z"}},
	//}}
	pb := PBModel{AA: &A{
		E:  &[]*B{{A: "z"}, {B: &C{C: "z"}}},
		G:  &B{A: "ziw"},
		A2: []*B{},
	}}
	var db DBModel

	ToCopyObject(&pb, &db)

	fmt.Println("=====================")

	if b, _ := json.Marshal(pb); true {
		fmt.Println(string(b))
	}

	fmt.Println("=====================")

	if b, _ := json.Marshal(db); true {
		fmt.Println(string(b))
	}
}
