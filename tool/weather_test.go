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
	Convey("TestWeather_getClientIP", t, func() {
		Convey("success", func() {
			ret, err := getClientIP()
			So(err, ShouldBeNil)
			t.Log(ret)
		})
	})
}

func Test_getIPInfo(t *testing.T) {
	initWeatherTest()
	Convey("TestWeather_getClientIP", t, func() {
		Convey("success", func() {
			ip, err := getClientIP()
			So(err, ShouldBeNil)
			logger.Info(ip)
			info := getIPInfo(ip)
			So(info, ShouldNotBeNil)
			logger.Info(info)
		})
	})
}
