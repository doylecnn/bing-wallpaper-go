package main

import (
	"embed"

	"github.com/getlantern/systray"
	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"github.com/wailsapp/wails/v2/pkg/options/windows"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

//go:embed all:frontend/dist
var assets embed.FS

//go:embed build/windows/icon.ico
var trayIcon []byte

func main() {
	app := NewApp()

	// 启动系统托盘
	go func() {
		systray.Run(onReady(app), func() {})
	}()

	err := wails.Run(&options.App{
		Title:  "Bing Wallpaper",
		Width:  1024,
		Height: 768,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 27, G: 38, B: 54, A: 1},
		OnStartup:        app.startup,
		OnBeforeClose:    app.beforeClose,
		Bind: []interface{}{
			app,
		},
		OnShutdown: app.Quit,
		Windows: &windows.Options{
			WebviewIsTransparent: false,
			WindowIsTranslucent:  false,
			DisableWindowIcon:    false,
			IsZoomControlEnabled: false,
		},
	})

	if err != nil {
		println("Error:", err.Error())
	}
}

func onReady(app *App) func() {
	return func() {
		systray.SetIcon(trayIcon)
		systray.SetTitle("Bing Wallpaper")
		systray.SetTooltip("Bing Wallpaper")

		// 菜单项：显示主窗口
		mShow := systray.AddMenuItem("显示主窗口", "显示主窗口")
		// 菜单项：更换壁纸
		mChange := systray.AddMenuItem("立即更换壁纸", "立即更换壁纸")
		// 分隔线
		systray.AddSeparator()
		// 菜单项：退出
		mQuit := systray.AddMenuItem("退出", "退出应用")

		// 处理菜单事件
		go func() {
			for {
				select {
				case <-mShow.ClickedCh:
					app.ShowWindow()
				case <-mChange.ClickedCh:
					app.ChangeWallpaperNow()
				case <-mQuit.ClickedCh:
					systray.Quit()
					app.isUserClosingWindow = false
					runtime.Quit(app.ctx)
					return
				}
			}
		}()
	}
}
