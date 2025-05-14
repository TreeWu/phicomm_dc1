package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"net"
	"strings"
	"time"
)

//go:generate  fyne package -os windows -icon icon.png -name "dc1-client" -appID "dc1-client"

type DeviceInfo struct {
	Action string `json:"action"`
	UUID   string `json:"uuid"`
	Auth   string `json:"auth"`
	Params struct {
		DeviceType string `json:"device_type"`
		Mac        string `json:"mac"`
		DeviceId   string `json:"device_id"`
		Status     int32  `json:"status"`
	} `json:"params"`
}

type StatusUpdate struct {
	UUID   string `json:"uuid"`
	Status int    `json:"status"`
	Result struct {
		Status int `json:"status"`
		I      int `json:"I"`
		V      int `json:"V"`
		P      int `json:"P"`
	} `json:"result"`
	Msg string `json:"msg"`
}

func main() {
	myApp := app.New()
	myWindow := myApp.NewWindow("设备控制器")

	// 输入框
	serverEntry := widget.NewEntry()
	serverEntry.SetPlaceHolder("服务器地址:端口")
	serverEntry.SetText("localhost:8000")

	macEntry := widget.NewEntry()
	macEntry.SetPlaceHolder("MAC地址 (84:F3:EB:BB:5A:19)")
	macEntry.SetText("84:F3:EB:BB:5A:10")

	deviceEntry := widget.NewEntry()
	deviceEntry.SetPlaceHolder("设备ID")
	deviceEntry.SetText("177227183548")
	firstActive := widget.NewCheck("", nil)

	// 使用 Check 代替 Switch
	mainCheck := widget.NewCheck("主开关", nil)
	check1 := widget.NewCheck("开关1", nil)
	check2 := widget.NewCheck("开关2", nil)
	check3 := widget.NewCheck("开关3", nil)

	// 连接按钮和状态标签
	connectBtn := widget.NewButton("连接", nil)
	clearBtn := widget.NewButton("清空消息", nil)

	// 使用 MultiLineEntry 来显示和保留消息
	statusEntry := widget.NewMultiLineEntry()
	statusEntry.SetPlaceHolder("状态消息将显示在这里...")
	statusEntry.Disable() // 设置为只读
	statusScroll := container.NewScroll(statusEntry)
	statusScroll.SetMinSize(fyne.NewSize(0, 300)) // 设置最小高度

	// 连接状态和TCP连接
	var conn net.Conn
	var connected bool

	// 更新状态函数
	updateStatus := func(text string) {
		fyne.Do(func() {
			currentTime := time.Now().Format("15:04:05")
			newText := fmt.Sprintf("[%s] %s\n%s", currentTime, text, statusEntry.Text)
			statusEntry.SetText(newText[:min(len(newText), 10000)])
			statusScroll.ScrollToBottom()
		})

	}
	sendmessage := func(msg any) error {
		if !connected {
			updateStatus("设备未连接，消息发送失败")
			return errors.New("未连接")
		}
		jsonData, err := json.Marshal(msg)
		if err != nil {
			updateStatus("状态JSON编码错误: " + err.Error())
			return err
		}

		_, err = conn.Write(append(jsonData, '\n'))
		if err != nil {
			updateStatus("状态发送失败: " + err.Error())
			connected = false
			connectBtn.SetText("连接")
			return err
		}
		updateStatus("消息发送成功:" + string(jsonData))
		return nil
	}

	// 清空消息
	clearBtn.OnTapped = func() {
		statusEntry.SetText("")
	}

	// 发送设备信息
	sendDeviceInfo := func() {
		if !connected || conn == nil {
			updateStatus("未连接，无法发送")
			return
		}
		if firstActive.Checked {
			updateStatus("设备首次联网")
			deviceInfo := DeviceInfo{
				Action: "activate=",
				UUID:   time.Now().Format("20060102150405"),
				Auth:   "",
			}
			deviceInfo.Params.Mac = macEntry.Text
			deviceInfo.Params.DeviceType = "PLUG_DC1_7"
			sendmessage(deviceInfo)
		} else {
			deviceInfo := DeviceInfo{
				Action: "identify",
				UUID:   time.Now().Format("20060102150405"),
				Auth:   "",
			}
			deviceInfo.Params.Mac = macEntry.Text
			deviceInfo.Params.DeviceId = deviceEntry.Text
			sendmessage(deviceInfo)
		}
	}

	// 发送状态更新
	sendStatusUpdate := func(msgid string) {
		status := 0
		if mainCheck.Checked {
			status += 1
		}
		if check1.Checked {
			status += 1 * 10
		}
		if check2.Checked {
			status += 1 * 100
		}
		if check3.Checked {
			status += 1 * 1000
		}

		if msgid == "" {
			msgid = time.Now().Format("20060102150405")
		}
		update := StatusUpdate{
			UUID:   msgid,
			Status: 200,
			Msg:    "get datapoint success",
		}
		update.Result.Status = status
		update.Result.I = 0
		update.Result.V = 235
		update.Result.P = 0

		sendmessage(update)
	}

	setStatusUpdate := func(deviceInfo DeviceInfo) {
		fyne.Do(func() {
			status := deviceInfo.Params.Status
			check3.SetChecked((status/1000)%10 == 1)
			check2.SetChecked((status/100)%10 == 1)
			check1.SetChecked((status/10)%10 == 1)
			mainCheck.SetChecked(status%10 == 1)
		})

		sendStatusUpdate(deviceInfo.UUID)
	}

	// 开关状态改变处理
	onCheckChanged := func(_ bool) {
		sendStatusUpdate("")
	}

	mainCheck.OnChanged = onCheckChanged
	check1.OnChanged = onCheckChanged
	check2.OnChanged = onCheckChanged
	check3.OnChanged = onCheckChanged

	// 连接按钮点击处理
	connectBtn.OnTapped = func() {
		if connected {
			// 断开连接
			if conn != nil {
				conn.Close()
			}
			connected = false
			connectBtn.SetText("连接")
			updateStatus("已断开连接")
			return
		}

		// 验证输入
		serverAddr := serverEntry.Text
		if serverAddr == "" {
			updateStatus("请输入服务器地址")
			return
		}

		macAddr := macEntry.Text
		if !isValidMAC(macAddr) {
			updateStatus("请输入有效的MAC地址")
			return
		}
		deviceId := deviceEntry.Text
		if deviceId == "" {
			updateStatus("请输入设备ID")
			return
		}

		// 连接服务器
		updateStatus("正在连接...")
		connectBtn.Disable()

		go func() {
			var err error
			conn, err = net.DialTimeout("tcp", serverAddr, 5*time.Second)
			if err != nil {
				updateStatus("连接失败: " + err.Error())
				fyne.Do(func() {
					connectBtn.Enable()
				})
				return
			}

			connected = true
			fyne.Do(func() {
				connectBtn.SetText("断开连接")
				connectBtn.Enable()
			})

			updateStatus("已连接")

			// 发送设备信息
			sendDeviceInfo()

			// 监听服务器响应
			go func() {
				buf := make([]byte, 2048)
				for connected {
					n, err := conn.Read(buf)
					if err != nil {
						if connected { // 检查是否是有意断开
							updateStatus("连接断开: " + err.Error())
						}
						connected = false
						fyne.Do(func() {
							connectBtn.SetText("连接")
						})
						return
					}
					s := string(buf[:n])
					updateStatus("收到消息: " + s)
					var deviceInfo DeviceInfo
					err = json.Unmarshal(buf[:n], &deviceInfo)
					if err != nil {
						updateStatus("状态JSON解码错误: " + err.Error())
					}
					if deviceInfo.Action == "datapoint" {
						sendStatusUpdate(deviceInfo.UUID)
					} else if deviceInfo.Action == "datapoint=" {
						setStatusUpdate(deviceInfo)
					}
				}
			}()
		}()
	}

	// 布局
	form := &widget.Form{
		Items: []*widget.FormItem{
			{Text: "服务器", Widget: serverEntry},
			{Text: "MAC地址", Widget: macEntry},
			{Text: "设备ID", Widget: deviceEntry},
			{Text: "首次联网", Widget: firstActive},
		},
	}

	checks := container.NewHBox(
		widget.NewLabel("开关控制"),
		mainCheck,
		check1,
		check2,
		check3,
	)

	content := container.NewVBox(
		form,
		checks,
		connectBtn,
		clearBtn,
		statusScroll,
	)

	myWindow.SetContent(content)
	myWindow.Resize(fyne.NewSize(400, 400)) // 增加窗口高度以适应消息框
	myWindow.ShowAndRun()
}

// 简单的MAC地址验证
func isValidMAC(mac string) bool {
	parts := strings.Split(mac, ":")
	if len(parts) != 6 {
		return false
	}
	for _, part := range parts {
		if len(part) != 2 {
			return false
		}
	}
	return true
}
