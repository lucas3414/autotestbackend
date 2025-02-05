package utils

import (
	"fmt"
	"github.com/spf13/viper"
	"github.com/tebeka/selenium"
	"github.com/tebeka/selenium/chrome"
	"go-gin-demo/global"
	"reflect"
	"runtime"
	"strconv"
	"testing"
	"time"
)

func GetOS() (SysType, ChromeDriverPath string, Port int, headless bool) {
	switch runtime.GOOS {
	case "windows":
		SysType = "windows"
		ChromeDriverPath = viper.GetString("selenium.winChromeDriverPath")
		Port = viper.GetInt("selenium.port")
	case "linux":
		SysType = "linux"
		ChromeDriverPath = viper.GetString("selenium.LinuxChromeDriverPath")
		Port = viper.GetInt("selenium.port")
	case "darwin":
		SysType = "darwin"
		ChromeDriverPath = viper.GetString("selenium.maxChromeDriverPath")
		Port = viper.GetInt("selenium.port")
	default:
		SysType = "other"
		ChromeDriverPath = viper.GetString("selenium.winChromeDriverPath")
		Port = viper.GetInt("selenium.port")
	}
	headless = viper.GetBool("selenium.isHeadless")
	global.Logger.Info("运行平台 ", SysType)
	global.Logger.Info("ChromeDriver 存放地址:", ChromeDriverPath)
	global.Logger.Info("ChromeDriver 端口:", Port)
	global.Logger.Info("是否启用无头浏览模式:", headless)
	return
}

func SetIsDisplayedElementStyle(value string, style string) string {
	//fmt.Println(fmt.Sprintf("document.evaluate(\"%s\", document).iterateNext().style.border = \"%s\"", value, style))
	return fmt.Sprintf("document.evaluate(\"%s\", document).iterateNext().style.border = \"%s\"", value, style)
}

type ServiceConfig struct {
	ChromeDriverPath string
	Port             int
	Config           []selenium.ServiceOption
	SetDebug         bool
}

// NewChromeDriverService 构造函数 初始化配置
func NewChromeDriverService(seleniumPath string, Port int, setDebug bool) ServiceConfig {
	return ServiceConfig{
		ChromeDriverPath: seleniumPath,
		Port:             Port,
		Config:           append(([]selenium.ServiceOption)(nil)),
		SetDebug:         setDebug,
	}
}

// NewService 构造函数 初始化服务
func NewService(config ServiceConfig) *selenium.Service {
	selenium.SetDebug(config.SetDebug)
	service, err := selenium.NewChromeDriverService(config.ChromeDriverPath, config.Port, config.Config...)
	if err != nil {
		panic(err)
	}
	return service
}

// NewWebDriver 构造函数 初始化driver,和服务
func NewWebDriver(config ServiceConfig, isHeadless bool, browser string) (service *selenium.Service, w WebDriver, err error) {
	service = NewService(config)
	caps := selenium.Capabilities{"browserName": browser}
	chromeOpts := chrome.Capabilities{
		Args: []string{
			"--disable-gpu",          // 谷歌文档提到需要加上这个属性来规避bug
			"--no-sandbox",           // 谷歌文档提到需要加上这个属性来规避bug
			"--single-process",       // 单线程运行
			"--disable-default-apps", // 禁用默认安装提示
			"--incognito",            // 启用隐身模式

		},
		ExcludeSwitches: []string{
			"enable-automation", // 解决chrome正在受自动化控制提示
		},
	}

	if isHeadless {
		chromeOpts.Args = append(chromeOpts.Args, "--headless")
	}
	caps.AddChrome(chromeOpts)
	w.Driver, err = selenium.NewRemote(caps, fmt.Sprintf("http://localhost:%d/wd/hub", config.Port))
	if err != nil {
		panic(err)
	}
	return service, w, nil
}

type WebDriver struct {
	Driver  selenium.WebDriver
	Element selenium.WebElement
	Result  string
	Msg     string
}

func (s *ServiceConfig) AddServiceOption(opt ...selenium.ServiceOption) {
	s.Config = append(s.Config, opt...)
}

func (d *WebDriver) WaitElementTimeout(value string) *WebDriver {
	var ele selenium.WebElement
	var err error
	//eleX, _ := d.Driver.PageSource()
	//fmt.Println(eleX)
	var IsElementDisplay = func(wd selenium.WebDriver) (bool, error) {
		ele, err = wd.FindElement(selenium.ByXPATH, value)

		if err != nil {
			errInfo := fmt.Sprintf("没找到此元素:%s", value)
			global.Logger.Error(errInfo)
			return false, nil
		}
		b, errs := ele.IsDisplayed()
		if errs != nil {
			errInfo := fmt.Sprintf("元素不可见:%s", value)
			global.Logger.Error(errInfo)
			return false, nil
		}
		return b, nil

	}
	timeout := viper.GetDuration("selenium.elementDefaultWaitTime") * time.Second
	if err = d.Driver.WaitWithTimeoutAndInterval(IsElementDisplay, timeout, 500*time.Millisecond); err != nil {
		errInfo := fmt.Sprintf("元素: %s, 等待超时, 默认超时时间%s", value, timeout)
		global.Logger.Error(errInfo)
		return &WebDriver{
			Driver:  d.Driver,
			Element: nil,
			Result:  "失败",
			Msg:     errInfo,
		}
	}
	d.SetElementStyle(value, "2px red dashed")
	time.Sleep(100 * time.Millisecond)
	d.SetElementStyle(value, "")
	return &WebDriver{
		Driver:  d.Driver,
		Element: ele,
		Result:  "成功",
		Msg:     "执行成功",
	}
}

func (d *WebDriver) WebDriverOpenUrl(url string) *WebDriver {
	err := d.Driver.Get(url)
	if err != nil {
		errInfo := fmt.Sprintf("无法访问到此地址, :%s", url)
		global.Logger.Error(errInfo)
		return &WebDriver{
			Driver: nil,
			Result: "失败",
			Msg:    errInfo,
		}
	}
	return &WebDriver{
		Driver: d.Driver,
		Result: "成功",
		Msg:    "执行成功",
	}
}

func (d *WebDriver) ElementSleep(num string) *WebDriver {
	sleepNum, err := strconv.ParseFloat(num, 64)
	if err != nil {
		errInfo := fmt.Sprintf("输入值不是数值类型%s", num)
		global.Logger.Error(errInfo)
		return &WebDriver{
			Driver: nil,
			Result: "失败",
			Msg:    errInfo,
		}
	}
	time.Sleep(time.Duration(sleepNum) * time.Second)
	return &WebDriver{
		Driver: d.Driver,
		Result: "成功",
		Msg:    "执行成功",
	}
}

func (d *WebDriver) GetElementValueByAttribute(xpath, name string) *WebDriver {
	var err error
	var resultString string
	wd := d.WaitElementTimeout(xpath)
	if wd.Result == "失败" {
		return &WebDriver{
			Driver: d.Driver,
			Result: "失败",
			Msg:    wd.Msg,
		}
	}
	resultString, err = wd.Element.GetAttribute(name)
	if err != nil {
		errInfo := fmt.Sprintf("获取元素:%s,属性值:%s失败 ", xpath, name)
		global.Logger.Error(errInfo)
		return &WebDriver{
			Driver: nil,
			Result: "失败",
			Msg:    errInfo,
		}
	}
	return &WebDriver{
		Driver: d.Driver,
		Result: "成功",
		Msg:    resultString,
	}
}

//// MoveX xxxxx
//func (d *WebDriver) MoveX() *WebDriver {
//	elem, _ := d.Driver.FindElement(selenium.ByXPATH, `//div[@class="details_sheet"]//div[@class="el-scrollbar__bar is-horizontal"]//div[@class="el-scrollbar__thumb"]`)
//	fmt.Println("滚动元素...")
//
//	for i := 0; i < 5; i++ {
//		ss := fmt.Sprintf("arguments[0].scrollLeft += 100", []interface{}{elem})
//		d.Driver.ExecuteScript(ss)
//		fmt.Println("%s,%s...", i, elem)
//		time.Sleep(500 * time.Millisecond)
//	}
//	return &WebDriver{
//		Driver: d.Driver,
//		Result: "成功",
//		Msg:    "滚动成功",
//	}
//}

// ButtonPermission 列表页面上普通按钮
func (d *WebDriver) ButtonPermission(btn string) *WebDriver {
	xpath := fmt.Sprintf(`//div[@class='btn-body']/span[contains(text(), '%s')]`, btn)
	return d.ClickByXpath(xpath)
}

// ButtonPermissionSelect 列表页面上相关操作 点击有下拉按钮
func (d *WebDriver) ButtonPermissionSelect(btn1, btn2 string) *WebDriver {
	xpath1 := fmt.Sprintf(`//div[@class='btn-line']//span[contains(text(), '%s')]`, btn1)
	xpath2 := fmt.Sprintf(`//li[@class='el-dropdown-menu__item' and text()='%s']`, btn2)
	return d.ClickByXpath(xpath1).ClickByXpath(xpath2)
}

// GlobalSearch 全域查询
func (d *WebDriver) GlobalSearch(value string) *WebDriver {
	xpath1 := fmt.Sprintf(`//div[@class='screen-card']//div[@class='screen-search']//input[@placeholder='请搜索']`)
	xpath2 := fmt.Sprintf(`//div[@class='screen-search']/button/span[text()='查询']`)
	return d.SendKeysByXpath(xpath1, value).ClickByXpath(xpath2)
}

// GeneralResetSearch 列表页面上重置按钮
func (d *WebDriver) GeneralResetSearch(value string) *WebDriver {
	xpath := fmt.Sprintf(`//div[@class='screen-search']/button/span[text()='%s']`, value)
	return d.ClickByXpath(xpath)
}

// OrderListSearchByValueAndSelect 列表搜索和选中
func (d *WebDriver) OrderListSearchByValueAndSelect(value string) *WebDriver {
	xpath := fmt.Sprintf(`//span[contains(text(), '%s')]/../../../../td//span[@class='vxe-cell--checkbox']/span`, value)
	return d.ClickByXpath(xpath)
}

// DigSelectWithOutDefaultDictValue 这个是对话框的下拉选择，下拉值来源于数据字典
func (d *WebDriver) DigSelectWithOutDefaultDictValue(digName, labelName, value string) *WebDriver {
	xpath1 := fmt.Sprintf(`//div[@role='dialog' and @aria-label='%s']//label[text()='%s']/..//div[@class='el-select__selected-item el-select__placeholder is-transparent']`, digName, labelName)
	xpath2 := fmt.Sprintf(`//div[@aria-hidden='false']//div[@class='el-scrollbar']//li//span[text()='%s']`, value)
	return d.ClickByXpath(xpath1).ClickByXpath(xpath2)
}

// DigSelectWithOutDefaultValue 这个是对话框的下拉选择，下拉值来源于其他数据
func (d *WebDriver) DigSelectWithOutDefaultValue(digName, labelName, value string) *WebDriver {
	xpath1 := fmt.Sprintf(`//div[@role='dialog' and @aria-label='%s']//label[text()='%s']/..//div[@class='el-select__selected-item el-select__placeholder is-transparent']`, digName, labelName)
	xpath2 := fmt.Sprintf(`//div[@aria-hidden='false']//div[@class='el-scrollbar']//span[contains(text(), '%s')]`, value)
	return d.ClickByXpath(xpath1).ClickByXpath(xpath2)
}

// DigButtonPermission 详情页面上的按钮
func (d *WebDriver) DigButtonPermission(digName, btn string) *WebDriver {
	xpath := fmt.Sprintf(`//div[@role='dialog' and @aria-label='%s']//button/span[text()='%s']`, digName, btn)
	return d.ClickByXpath(xpath)
}

// OrderDetailInput 详情页面上普通输入框
func (d *WebDriver) OrderDetailInput(labelName, value string) *WebDriver {
	xpath := fmt.Sprintf(`//div[@class='el-card__body']//label[text()='%s']/../div//input`, labelName)
	return d.ClickByXpath(xpath).SendKeysByXpath(xpath, value)
}

// OrderDetailTextarea 详情页面上富文本输入框
func (d *WebDriver) OrderDetailTextarea(labelName, value string) *WebDriver {
	xpath := fmt.Sprintf(`//div[@class='el-card__body']//label[text()='%s']/../div//textarea`, labelName)
	return d.SendKeysByXpath(xpath, value)
}

// OrderDetailSelectWithOutDefaultValue 详情页面上下拉没有默认值
func (d *WebDriver) OrderDetailSelectWithOutDefaultValue(labelName, value string) *WebDriver {
	xpath1 := fmt.Sprintf(`//div[@class='el-card__body']//label[text()='%s']/..//div[@class='el-select__selected-item el-select__placeholder is-transparent']`, labelName)
	xpath2 := fmt.Sprintf(`//div[@aria-hidden='false' and @role='tooltip']//span[contains(text(),'%s')]`, value)
	return d.ClickByXpath(xpath1).ClickByXpath(xpath2)
}

// OrderDetailSelectWithDefaultValue 详情页面上下拉有默认值
func (d *WebDriver) OrderDetailSelectWithDefaultValue(labelName, value string) *WebDriver {
	xpath1 := fmt.Sprintf(`//div[@class='el-card__body']//label[text()='%s']/..//div[@class='el-select__selected-item el-select__placeholder']`, labelName)
	xpath2 := fmt.Sprintf(`//div[@aria-hidden='false' and @role='tooltip']//span[contains(text(),'%s')]`, value)
	return d.ClickByXpath(xpath1).ClickByXpath(xpath2)
}

func (d *WebDriver) ElementAssert(xpath, key string) *WebDriver {
	var err error
	var resultString string
	wd := d.WaitElementTimeout(xpath)
	if wd.Result == "失败" {
		return &WebDriver{
			Driver: d.Driver,
			Result: "失败",
			Msg:    wd.Msg,
		}
	}
	resultString, err = wd.Element.Text()
	if err != nil {
		errInfo := fmt.Sprintf("元素不可输入:%s, %s", xpath, key)
		global.Logger.Error(errInfo)
		return &WebDriver{
			Driver: nil,
			Result: "失败",
			Msg:    errInfo,
		}
	}
	if resultString != key {
		errInfo := fmt.Sprintf("断言失败, 预期值:%s, 实际值:%s, %s != %s", resultString, key, resultString, key)
		global.Logger.Error(errInfo)
		return &WebDriver{
			Driver: nil,
			Result: "失败",
			Msg:    errInfo,
		}
	}
	successInfo := fmt.Sprintf("预期值:%s, 实际值:%s", resultString, key)
	return &WebDriver{
		Driver: d.Driver,
		Result: "成功",
		Msg:    successInfo,
	}
}

func (d *WebDriver) WebDriverQuit() *WebDriver {
	err := d.Driver.Quit()
	if err != nil {
		errInfo := fmt.Sprintf("关闭浏览器失败,")
		global.Logger.Error(errInfo)
		return &WebDriver{
			Driver: nil,
			Result: "失败",
			Msg:    errInfo,
		}
	}
	return &WebDriver{
		Driver: d.Driver,
		Result: "成功",
		Msg:    "执行成功",
	}
}

func (d *WebDriver) SetElementStyle(value string, style string) *WebDriver {
	_, err := d.Driver.ExecuteScript(SetIsDisplayedElementStyle(value, style), nil)
	if err != nil {
		errInfo := fmt.Sprintf("设置元素:%s,样式:%s", value, style)
		global.Logger.Error(errInfo)
		return &WebDriver{
			Driver: nil,
			Result: "失败",
			Msg:    errInfo,
		}
	}
	return &WebDriver{
		Driver: d.Driver,
		Result: "成功",
		Msg:    "执行成功",
	}

}

func (d *WebDriver) WebDriverMaximizeWindow(size string) *WebDriver {
	err := d.Driver.MaximizeWindow(size)
	if err != nil {
		errInfo := fmt.Sprintf("浏览器设置大小失败:%s", size)
		global.Logger.Error(errInfo)
		return &WebDriver{
			Driver: nil,
			Result: "失败",
			Msg:    errInfo,
		}
	}
	return &WebDriver{
		Driver: d.Driver,
		Result: "成功",
		Msg:    "执行成功",
	}
}

func (d *WebDriver) ClearXpath(xpath string) *WebDriver {
	var err error
	//var ele selenium.WebElement
	wd := d.WaitElementTimeout(xpath)
	if wd.Result == "失败" {
		return &WebDriver{
			Driver: d.Driver,
			Result: "失败",
			Msg:    wd.Msg,
		}
	}
	err = wd.Element.Clear()
	if err != nil {
		errInfo := fmt.Sprintf("元素不可清除:%s", xpath)
		global.Logger.Error(errInfo)
		return &WebDriver{
			Driver: nil,
			Result: "失败",
			Msg:    errInfo,
		}
	}
	return &WebDriver{
		Driver: d.Driver,
		Result: "成功",
		Msg:    "执行成功",
	}
}

func (d *WebDriver) ClickByXpath(xpath string) *WebDriver {
	var err error
	wd := d.WaitElementTimeout(xpath)
	if wd.Result == "失败" {
		return &WebDriver{
			Driver: d.Driver,
			Result: "失败",
			Msg:    wd.Msg,
		}
	}
	err = wd.Element.Click()
	if err != nil {
		errInfo := fmt.Sprintf("元素不可点击:%s", xpath)
		global.Logger.Error(errInfo)
		return &WebDriver{
			Driver: nil,
			Result: "失败",
			Msg:    errInfo,
		}
	}
	return &WebDriver{
		Driver: d.Driver,
		Result: "成功",
		Msg:    "执行成功",
	}
}

func (d *WebDriver) SendKeysByXpath(xpath, key string) *WebDriver {
	var err error
	wd := d.WaitElementTimeout(xpath)
	if wd.Result == "失败" {
		return &WebDriver{
			Driver: d.Driver,
			Result: "失败",
			Msg:    wd.Msg,
		}
	}
	err = wd.Element.SendKeys(key)
	if err != nil {
		errInfo := fmt.Sprintf("元素不可输入:%s, %s", xpath, key)
		global.Logger.Error(errInfo)
		return &WebDriver{
			Driver: nil,
			Result: "失败",
			Msg:    errInfo,
		}
	}
	return &WebDriver{
		Driver: d.Driver,
		Result: "成功",
		Msg:    "执行成功",
	}
}

func (d *WebDriver) ClickAndSendKeysByXpath(xpath, key string) *WebDriver {
	d.ClickByXpath(xpath).SendKeysByXpath(xpath, key)
	return &WebDriver{
		Driver: d.Driver,
		Result: "成功",
		Msg:    "执行成功",
	}
}

func (d *WebDriver) ClickAndClearAndSendKeysByXpath(xpath, key string) *WebDriver {
	d.ClickByXpath(xpath).ClearXpath(xpath).SendKeysByXpath(xpath, key)
	return &WebDriver{
		Driver: d.Driver,
		Result: "成功",
		Msg:    "执行成功",
	}
}

func (d *WebDriver) SelectValueByXpath(inputXpath, valueXpath string) *WebDriver {
	d.ClickByXpath(inputXpath).ClickByXpath(valueXpath)
	return &WebDriver{
		Driver: d.Driver,
		Result: "成功",
		Msg:    "执行成功",
	}
}

func CallMethod(receiver any, methodName string, argList []any) (result []reflect.Value) {
	method := reflect.ValueOf(receiver).MethodByName(methodName)
	if method.Kind() == reflect.Invalid {
		global.Logger.Error(fmt.Sprintf(" 没有找到:%s此方法 \n", methodName))
		return
	}
	in := make([]reflect.Value, len(argList))
	for i, arg := range argList {
		in[i] = reflect.ValueOf(arg)
	}
	defer func() {
		if r := recover(); r != nil {
			errMsg := fmt.Sprintf("报错来源: %s, %s, %s", methodName, in, r)
			fmt.Println(errMsg)
			global.Logger.Error(errMsg)
		}
	}()
	results := method.Call(in)
	return results
}

func RunCase(caseListMap []map[string]any) []map[string]any {
	_, ChromeDriverPath, Port, flag := GetOS()
	// 测试一下. 首先配置服务选项
	opt := []selenium.ServiceOption{
		//selenium.Output(os.Stderr),
		selenium.ChromeDriver(ChromeDriverPath),
	}
	config := NewChromeDriverService(ChromeDriverPath, Port, false)
	config.AddServiceOption(opt...)

	_, driver, _ := NewWebDriver(config, flag, "chrome")

	for _, caseMap := range caseListMap {
		methodName, _ := caseMap["method_name"].(string)
		args, _ := caseMap["args"].([]any)
		global.Logger.Info(fmt.Sprintf("%s(%s)", methodName, args))
		sTime := time.Now()
		caseMap["start_time"] = sTime.Format("2006-01-02 15:04:05")
		res := CallMethod(&driver, methodName, args)
		eTime := time.Now()
		caseMap["end_time"] = eTime.Format("2006-01-02 15:04:05")
		costTime := eTime.Sub(sTime).Seconds()
		caseMap["cost_time"] = costTime
		caseMap["result"] = res[0].Elem().Field(2).String()
		caseMap["msg"] = res[0].Elem().Field(3).String()
		global.Logger.Info("耗时:", costTime)

		// 控制执行的开关 caseBreak为true 用例有一个执行保持则自动停止
		//caseBreak := caseMap["case_break"].(bool)
		caseBreak, _ := strconv.ParseBool(caseMap["case_break"].(string))
		//caseBreak := caseMap["case_break"].(bool)
		if caseBreak && res[0].Elem().Field(2).String() == "失败" {
			CallMethod(&driver, "WebDriverQuit", []any{})
			break
		}

	}
	return caseListMap
}

func TestSelenium(t *testing.T) {
	url := viper.GetString("selenium.base_url")

	//defer service.Stop()
	//const url = "http://218.94.154.54:46400"
	//const url = "http://10.0.2.155:6400"
	//责任链调用
	//driver.WebDriverOpenUrl(url).WebDriverMaximizeWindow().
	//	SendKeysByXpath("//input[@placeholder='租户']", "test").
	//	SendKeysByXpath("//input[@placeholder='用户名']", "test").
	//	SendKeysByXpath("//input[@placeholder='密码']", "test").
	//	ClickByXpath("//button//span[text()='登 录']/..").
	//	WebDriverQuit()
	//反射调用
	//callMethod(&driver, "WebDriverOpenUrl", url)
	//callMethod(&driver, "WebDriverMaximizeWindow", "")
	//callMethod(&driver, "SendKeysByXpath", "//input[@placeholder='租户']", "test")
	//callMethod(&driver, "SendKeysByXpath", "//input[@placeholder='用户名']", "test")
	//callMethod(&driver, "SendKeysByXpath", "//input[@placeholder='密码']", "test")
	//callMethod(&driver, "ClickByXpath", "//button//span[text()='登 录']/..")
	//callMethod(&driver, "WebDriverQuit")

	caseList := make([]map[string]any, 10)

	caseList = append(caseList,
		map[string]any{"WebDriverOpenUrl": []string{url}},
		map[string]any{"WebDriverMaximizeWindow": []string{""}},
		map[string]any{"SendKeysByXpath": []string{"//input[@placeholder='租户']", "test"}},
		map[string]any{"SendKeysByXpath": []string{"//input[@placeholder='用户名']", "test"}},
		map[string]any{"SendKeysByXpath": []string{"//input[@placeholder='密码']", "test"}},
		map[string]any{"ClickByXpath": []string{"//button//span[text()='登 录']/.."}},
		//map[string]any{"WebDriverQuit": []string{}},
		map[string]any{"WebDriverQuit": []string{}},
	)

	RunCase(caseList)

}
