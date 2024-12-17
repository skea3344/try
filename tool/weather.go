// -------------------------------------------
// @file      : weather.go
// @author    : bo cai
// @contact   : caibo923@gmail.com
// @time      : 2024/12/6 下午3:09
// -------------------------------------------

package tool

import (
	"encoding/json"
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
	City    string        `json:"city"`
	Weather []WeatherInfo `json:"weather"`
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
	if err = ui.Init(); err != nil {
		return err
	}
	defer ui.Close()
	table := widgets.NewTable()
	table.Title = res.City + "天气"
	table.Rows = [][]string{
		{"日期", "天气", "风向", "体感温度"},
	}
	for _, v := range res.Weather {
		table.Rows = append(table.Rows, []string{v.Date, v.Weather, v.Wind, v.Temp})
	}
	table.TextStyle = ui.NewStyle(ui.ColorGreen)
	table.TitleStyle = ui.NewStyle(ui.ColorBlue)
	table.SetRect(0, 0, 60, 10)
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
	api := "https://api.asilu.com/weather/?city=" + city
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
	err = json.Unmarshal(out, ret)
	return ret, err
}
