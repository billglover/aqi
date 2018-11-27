package aqi

import (
	"context"
	"fmt"
	"net/http"
	"reflect"
	"testing"
)

var measurements = []struct {
	name string
	mock string
	want Measurement
}{
	{
		name: "successful reading",
		mock: `{"status":"ok","data":{"aqi":90,"idx":459,"attributions":[{"url":"http://www.bjmemc.com.cn/","name":"Beijing Environmental Protection Monitoring Center (北京市环境保护监测中心)"},{"url":"https://waqi.info/","name":"World Air Quality Index Project"}],"city":{"geo":[39.718,116.404],"name":"Huangcunzhen, Daxing, Beijing (北京大兴黄村镇)","url":"https://aqicn.org/city/beijing/daxinghuangcunzhen"},"dominentpol":"pm10","iaqi":{"co":{"v":9.1},"no2":{"v":15.1},"o3":{"v":4.5},"pm10":{"v":90},"pm25":{"v":78},"so2":{"v":1.6},"w":{"v":2.5},"wg":{"v":10.8}},"time":{"s":"2018-11-27 18:00:00","tz":"+08:00","v":1543341600},"debug":{"sync":"2018-11-27T19:18:40+09:00"}}}`,
		want: Measurement{AQI: 90},
	},
}

func TestLatest(t *testing.T) {
	for _, tc := range measurements {
		t.Run(tc.name, func(st *testing.T) {
			testLatest(t, tc.name, tc.mock, tc.want)
		})
	}
}

func testLatest(t *testing.T, name string, mock string, want Measurement) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/feed/", func(w http.ResponseWriter, r *http.Request) {
		got, want := r.URL.Query().Get("token"), "dummytoken"
		if got != want {
			t.Errorf("unexpected token: got %s, want %s", got, want)
		}
		fmt.Fprint(w, mock)
	})

	got, err := client.Latest(context.Background(), "beijing/daxing")
	if err != nil {
		t.Error("error returned:", err)
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("unexpected Measurement Returned: %+v", got)
	}
}
