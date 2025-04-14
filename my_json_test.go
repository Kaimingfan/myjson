package myjson

import (
	"testing"

	"code.byted.org/gopkg/mockito"
	"github.com/smartystreets/goconvey/convey"
)

type simpleStruct1 struct {
	Int1    int           `json:"int1"`
	Float2  float32       `json:"float2"`
	Name    string        `json:"name"`
	Unknown []interface{} `json:"unknown"`
}
type simpleStruct2 struct {
	Name string         `json:"name"`
	S1   *simpleStruct1 `json:"s1"`
}

func TestMarshalStruct(t *testing.T) {
	type args struct {
		v interface{}
	}
	type testConfig struct {
		args    args
		want    string
		wantErr bool
	}
	mockito.PatchConvey("test", t, func() {
		// your mock code...
		mockito.PatchConvey("simpleStruct test case1", func() {
			tt := testConfig{
				args: args{
					v: simpleStruct2{
						Name: "s2",
						S1: &simpleStruct1{
							Int1:    12,
							Float2:  1.23,
							Name:    "s1",
							Unknown: []interface{}{},
						},
					},
				},
				wantErr: false,
				want:    "{\"name\":\"s2\",\"s1\":{\"int1\":12,\"float2\":1.23,\"name\":\"s1\",\"unknown\":[]}}",
			}
			got, err := MarshalString(tt.args.v)
			convey.So(err != nil, convey.ShouldEqual, tt.wantErr)
			convey.So(got, convey.ShouldEqual, tt.want)
		})

		// 未兼容unicode
		// mockito.PatchConvey("twitter test case", func() {
		// 	data := testdata.GetTwiterJson()
		// 	bs := []byte(data)
		// 	var st testdata.Twitter
		// 	json.Unmarshal(bs, &st)
		// 	data = strings.Replace(data, " ", "", -1)
		// 	data = strings.Replace(data, "\n", "", -1)
		// 	tt := testConfig{
		// 		args: args{
		// 			v: st,
		// 		},
		// 		wantErr: false,
		// 		want:    data,
		// 	}
		// 	got, err := MarshalString(tt.args.v)
		// 	convey.So(err != nil, convey.ShouldEqual, tt.wantErr)
		//
		// 	res := testdata.Twitter{}
		// 	err = json.Unmarshal([]byte(got), &res)
		// 	convey.So(err != nil, convey.ShouldEqual, tt.wantErr)
		// 	equal := reflect.DeepEqual(res, st)
		// 	convey.So(equal, convey.ShouldEqual, true)
		// })
	})
}
