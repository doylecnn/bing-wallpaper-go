package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"math/rand"

	"github.com/reujab/wallpaper"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// App struct
type App struct {
	ctx                 context.Context
	cancel              context.CancelFunc
	isUserClosingWindow bool
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{isUserClosingWindow: true}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx, a.cancel = context.WithCancel(ctx)
	// 每天凌晨更新壁纸
	go func() {
		ticker := time.NewTicker(24 * time.Hour)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				if err := a.AutoSetWallpaper(); err != nil {
					fmt.Printf("Failed to auto set wallpaper: %v\n", err)
				}
			case <-a.ctx.Done():
				return
			}
		}
	}()
}

// GetAndSetWallpaper 获取并设置壁纸
func (a *App) GetAndSetWallpaper(url string) (string, error) {
	// 下载图片到本地
	homeDir, _ := os.UserHomeDir()
	imagePath := filepath.Join(homeDir, "Pictures", "bing_wallpaper.jpg")
	if err := DownloadImage(url, imagePath); err != nil {
		return "", err
	}

	// 设置壁纸
	err := wallpaper.SetFromFile(imagePath)
	return imagePath, err
}

// 壁纸信息结构体
type BingImage struct {
	StartDate     string `json:"startdate"`
	FullStartDate string `json:"fullstartdate"`
	EndDate       string `json:"enddate"`
	URL           string `json:"url"`
	URLBase       string `json:"urlbase"`
	Copyright     string `json:"copyright"`
	Title         string `json:"title"`
	Quiz          string `json:"quiz"`
}

// 获取Bing每日壁纸URL
func GetBingImageURL() (string, error) {
	resp, err := http.Get("https://www.bing.com/HPImageArchive.aspx?format=js&idx=0&n=1")
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var data struct {
		Images []BingImage `json:"images"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return "", err
	}

	if len(data.Images) > 0 {
		return "https://www.bing.com" + data.Images[0].URL, nil
	}

	return "", nil
}

// 下载图片到本地
func DownloadImage(url, dest string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// 创建文件
	out, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer out.Close()

	// 保存文件
	_, err = io.Copy(out, resp.Body)
	return err
}

// 添加这个方法到 App struct
type BingWallpaper struct {
	URL       string `json:"url"`
	Title     string `json:"title"`
	Copyright string `json:"copyright"`
	StartDate string `json:"startdate"`
}

func (a *App) GetBingWallpapers() ([]BingWallpaper, error) {
	resp, err := http.Get("https://www.bing.com/HPImageArchive.aspx?format=js&idx=0&n=8")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var data struct {
		Images []BingImage `json:"images"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, err
	}

	wallpapers := make([]BingWallpaper, len(data.Images))
	for i, img := range data.Images {
		wallpapers[i] = BingWallpaper{
			URL:       "https://www.bing.com" + img.URL,
			Title:     img.Title,
			Copyright: img.Copyright,
			StartDate: img.StartDate,
		}
	}

	return wallpapers, nil
}

func (a *App) AutoSetWallpaper() error {
	wallpapers, err := a.GetBingWallpapers()
	if err != nil {
		return err
	}

	if len(wallpapers) > 0 {
		_, err = a.GetAndSetWallpaper(wallpapers[0].URL)
		return err
	}

	return fmt.Errorf("no wallpapers available")
}

// 隐藏窗口
func (a *App) beforeClose(ctx context.Context) (prevent bool) {
	// 检查是否是用户点击窗口关闭按钮
	if a.isUserClosingWindow {
		runtime.Hide(ctx)
		return true // 阻止关闭，仅隐藏窗口
	}

	// 如果是通过 runtime.Quit 调用，允许程序退出
	return false
}

// 退出应用
func (a *App) Quit(ctx context.Context) {
	if a.cancel != nil {
		a.cancel() // Cancel any running goroutines
	}
}

func (a *App) ShowWindow() {
	runtime.WindowShow(a.ctx)
}

// Add this method to your App struct
func (a *App) ChangeWallpaperNow() {
	wps, err := a.GetBingWallpapers()
	if err != nil {
		log.Println("Error getting wallpapers:", err)
		return
	}

	if len(wps) == 0 {
		log.Println("No wallpapers available")
		return
	}

	randomIndex := rand.Intn(len(wps))
	_, err = a.GetAndSetWallpaper(wps[randomIndex].URL)
	if err != nil {
		log.Println("Error setting wallpaper:", err)
	}
}
