// -------------------------------------------
// @file      : weather_test.go
// @author    : bo cai
// @contact   : caibo923@gmail.com
// @time      : 2024/12/6 下午3:23
// -------------------------------------------

package tool

import (
	"github.com/caibo86/logger"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func initWeatherTest() {
	logger.Init(
		logger.SetStacktrace(logger.FatalLevel),
	)
}

func Test_getClientIP(t *testing.T) {
	initWeatherTest()
	Convey("取本地IP地址", t, func() {
		Convey("success", func() {
			ret, err := getClientIP()
			So(err, ShouldBeNil)
			t.Log(ret)
		})
	})
}

func Test_getIPInfo(t *testing.T) {
	initWeatherTest()
	Convey("取IP地址对应的地址信息", t, func() {
		Convey("success", func() {
			ip, err := getClientIP()
			So(err, ShouldBeNil)
			logger.Info(ip)
			info, err := getIPInfo(ip)
			So(err, ShouldBeNil)
			So(info, ShouldNotBeNil)
			logger.Info(info)
		})
	})
}

func Test_getWeatherInfo(t *testing.T) {
	initWeatherTest()
	Convey("取指定城市的天气信息", t, func() {
		Convey("success", func() {
			ret, err := getWeatherInfo("")
			So(err, ShouldBeNil)
			logger.Info(ret)
		})
	})
}
