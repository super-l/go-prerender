package internal

import (
	"context"
	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
	"go-prerender/internal/config"
	"log"
	"time"
)

func buildOption() []chromedp.ExecAllocatorOption {
	agent := "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/102.0.0.0 Safari/537.36"
	options := []chromedp.ExecAllocatorOption{
		chromedp.NoFirstRun,            // 第一次不运行
		chromedp.NoDefaultBrowserCheck, // 不检查默认浏览器
		chromedp.Flag("disable-background-networking", true),
		chromedp.Flag("disable-background-timer-throttling", true),
		chromedp.Flag("disable-backgrounding-occluded-windows", true),
		chromedp.Flag("disable-breakpad", true),
		chromedp.Flag("disable-client-side-phishing-detection", true),
		chromedp.Flag("disable-default-apps", true),
		chromedp.Flag("disable-dev-shm-usage", true),
		chromedp.Flag("disable-extensions", true), // 禁用扩展
		chromedp.Flag("disable-plugins", true),    // 禁用插件
		chromedp.Flag("disable-hang-monitor", true),
		chromedp.Flag("disable-ipc-flooding-protection", true),
		chromedp.Flag("disable-popup-blocking", true),
		chromedp.Flag("disable-prompt-on-repost", true),
		chromedp.Flag("disable-renderer-backgrounding", true),
		chromedp.Flag("disable-sync", true),
		chromedp.Flag("disable-web-security", true),
		chromedp.Flag("safebrowsing-disable-auto-update", true),
		chromedp.Flag("disable-features", "site-per-process,TranslateUI,BlinkGenPropertyTrees"),
		chromedp.Flag("enable-automation", false), // 隐藏调试

		//chromedp.Flag("blink-settings", "imagesEnabled=false"), // 禁用图片
		chromedp.Flag("excludeSwitches", "enable-automation"),
		chromedp.Flag("mute-audio", false), // 关闭音频

		// 远程调试地址 0.0.0.0 可以外网调用但是安全性低,建议使用默认值 127.0.0.1
		chromedp.Flag("remote-debugging-address", "127.0.0.1"), // 限制IP

		chromedp.Flag("disable-infobars", true),
		chromedp.Flag("disable-gpu", true),                 // 禁用谷歌浏览器GPU加速-配置1(linux上用)
		chromedp.Flag("disable-software-rasterizer", true), // 禁用谷歌浏览器GPU加速-配置2(linux上用)
		chromedp.Flag("incognito", true),
		chromedp.Flag("ignore-certificate-errors", true), // 屏蔽--ignore-certificate-errors提示信息
		chromedp.Flag("page-load-strategy", "eager"),     // https://blog.csdn.net/yexiaomodemo/article/details/99958509
		chromedp.UserAgent(agent),
		//chromedp.NoSandbox,             // 不开启沙盒模式可以减少对服务器的资源消耗,但是服务器安全性降低,配和参数
	}

	if config.GetConfig().Browser.Show == 0 {
		options = append(options, chromedp.Flag("headless", true))
	} else {
		options = append(options, chromedp.Flag("headless", false))
	}
	return options
}

func GetSource(url string) (string, error) {
	var source string
	allocCtx, cancelWindows := chromedp.NewExecAllocator(context.Background(), buildOption()...)
	defer cancelWindows()

	// 创建标签页
	ctxTab, cancelTab := chromedp.NewContext(allocCtx, chromedp.WithLogf(log.Printf))

	// 设置超时时间，秒
	ctxTab, cancelTab = context.WithTimeout(ctxTab, time.Duration(config.GetConfig().Browser.Timeout)*time.Second)

	defer cancelTab()

	err := chromedp.Run(ctxTab, runTasks(url, &source))
	if err != nil {
		SLogger.GetStdoutLogger().Error(err)
		return "", err
	}
	return source, nil
}

func runTasks(url string, source *string) chromedp.Tasks {
	// 任务组
	task := chromedp.Tasks{
		// 设置webdriver检测反爬
		chromedp.ActionFunc(func(cxt context.Context) error {
			_, err := page.AddScriptToEvaluateOnNewDocument("Object.defineProperty(navigator, 'webdriver', { get: () => undefined, });").Do(cxt)
			return err
		}),

		chromedp.Navigate(url),

		// 获取源码
		chromedp.ActionFunc(func(ctx context.Context) error {
			chromedp.OuterHTML("html", source).Do(ctx)
			return nil
		}),
	}
	return task
}
