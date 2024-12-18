// -------------------------------------------
// @file      : weather.go
// @author    : bo cai
// @contact   : caibo923@gmail.com
// @time      : 2024/12/6 下午3:09
// -------------------------------------------

package tool

import (
	"encoding/json"
	"fmt"
	"github.com/caibo86/logger"
	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
	"io"
	"net/http"
	"strings"
)

type Weather struct {
}

type IPInfo struct {
	Status string `json:"status"`
	Data   []IP   `json:"data"`
}

type IP struct {
	Location string `json:"location"`
}

type WeatherRes struct {
	Forecasts []*Forecast `json:"forecasts"`
}

// Forecast 天气预告
type Forecast struct {
	City     string  `json:"city"`     // 城市
	AdCode   string  `json:"adcode"`   // 邮编
	Province string  `json:"province"` // 省份
	Casts    []*Cast `json:"casts"`    // 预告数据
}

// Cast 单日预告数据
type Cast struct {
	Date         string `json:"date"`         // 日期
	Week         string `json:"week"`         // 星期
	DayWeather   string `json:"dayweather"`   // 白天天气
	NightWeather string `json:"nightweather"` // 夜晚天气
	DayTemp      string `json:"daytemp"`      // 白天温度
	NightTemp    string `json:"nighttemp"`    // 夜晚温度
	DayWind      string `json:"daywind"`      // 白天风向
	NightWind    string `json:"nightwind"`    // 夜晚风向
	DayPower     string `json:"daypower"`     // 白天风力
	NightPower   string `json:"nightpower"`   // 夜晚风力
}

func (cast *Cast) Temperature() string {
	return fmt.Sprintf("%s-%s度", cast.NightTemp, cast.DayTemp)
}

func (cast *Cast) Weather() string {
	return fmt.Sprintf("日间:%s\n夜晚:%s", cast.DayWeather, cast.NightWeather)
}

func (cast *Cast) Wind() string {
	return fmt.Sprintf("日间:%s风%s级\n夜晚:%s风%s级",
		cast.DayWind, cast.DayPower, cast.NightWind, cast.NightPower)
}

type WeatherInfo struct {
	Date    string `json:"date"`
	Weather string `json:"weather"`
	Temp    string `json:"temp"`
	W       string `json:"w"`
	Wind    string `json:"wind"`
}

func (w *Weather) GetWeather(city string) error {
	var err error
	if city == "" {
		city, err = getLocalCity()
		if err != nil {
			return err
		}
	}
	res, err := getWeatherInfo(city)
	if err != nil {
		return err
	}
	if len(res.Forecasts) == 0 {
		return fmt.Errorf("未找到%s的天气信息", city)
	}
	forecast := res.Forecasts[0]
	if err = ui.Init(); err != nil {
		return err
	}
	defer ui.Close()
	table := widgets.NewTable()
	table.Title = forecast.City + "天气"
	table.Rows = [][]string{
		{"日期", "天气", "风向", "温度"},
	}
	for _, cast := range forecast.Casts {
		table.Rows = append(table.Rows, []string{cast.Date, cast.Weather(), cast.Wind(), cast.Temperature()})
	}
	table.TextStyle = ui.NewStyle(ui.ColorGreen)
	table.TitleStyle = ui.NewStyle(ui.ColorBlue)
	table.SetRect(0, 0, 90, 10)
	table.ColumnWidths = []int{15, 25, 35, 15}
	ui.Render(table)
	uiEvents := ui.PollEvents()
	for {
		e := <-uiEvents
		switch e.ID {
		case "q", "<C-c>":
			return nil
		case "c":
		}
	}
}

func getLocalCity() (string, error) {
	ip, err := getClientIP()
	if err != nil {
		return ip, err
	}
	ipInfo, err := getIPInfo(ip)
	if err != nil {
		return "", err
	}
	ret := extractCity(ipInfo.Data[0].Location)
	return ret, nil
}

// 获取本机IP
func getClientIP() (string, error) {
	resp, err := http.Get("https://ifconfig.me/")
	if err != nil {
		return "", err
	}
	defer func() {
		_ = resp.Body.Close()
	}()
	content, err := io.ReadAll(resp.Body)
	return string(content), err
}

// 获取IP相关数据
func getIPInfo(ip string) (*IPInfo, error) {
	// url := "https://ip.taobao.com/outGetIpInfo?ip=" + ip + "&accessKey=alibaba-inc"
	url := "https://opendata.baidu.com/api.php?query=" + ip + "&co=&resource_id=6006&ie=utf8&oe=utf8&format=json"
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = resp.Body.Close()
	}()
	out, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	logger.Info("ip info: ", string(out))
	var ret IPInfo
	if err = json.Unmarshal(out, &ret); err != nil {
		logger.Errorf("unmarshal failed, err: %v", err)
		return nil, err
	}
	return &ret, nil
}

func extractCity(s string) string {
	cityKey := "市"
	provinceKey := "省"
	cityIndex := strings.Index(s, cityKey)
	if cityIndex != -1 {
		provinceIndex := strings.Index(s[:cityIndex], provinceKey)
		if provinceIndex != -1 {
			return s[provinceIndex+len(provinceKey) : cityIndex+len(cityKey)]
		}
		return s[:cityIndex+len(cityKey)]
	}
	return "未找到市"
}

func getWeatherInfo(city string) (*WeatherRes, error) {
	api := "https://query.asilu.com/weather/gaode?address=" + city
	resp, err := http.Get(api)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = resp.Body.Close()
	}()
	out, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	ret := &WeatherRes{}
	err = json.Unmarshal(out, &ret)
	if err != nil {
		return nil, err
	}
	return ret, err
}
