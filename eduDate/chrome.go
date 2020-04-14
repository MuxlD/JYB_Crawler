package eduDate

import (
	"context"

	"github.com/chromedp/chromedp"
)

type Opts struct {
	IgnoreCertificateErrors bool
	UserAgent               string
	Proxy                   string
}

type ChromeBrowser struct {
	allocCtx context.Context
	cancel   context.CancelFunc
}

func (chrome *ChromeBrowser) NewTab() (context.Context, context.CancelFunc) {
	return chromedp.NewContext(chrome.allocCtx)
}

func (chrome *ChromeBrowser) Close() {
	chrome.cancel()
}

//OpenBrowser 打开浏览器
func OpenBrowser(ctx context.Context, opts *Opts) *ChromeBrowser {

	ops := initOpts(opts)

	allocCtx, cancel := chromedp.NewExecAllocator(ctx, ops...)

	return &ChromeBrowser{
		allocCtx: allocCtx,
		cancel:   cancel,
	}
}

func initOpts(opts *Opts) []chromedp.ExecAllocatorOption {

	rets := []chromedp.ExecAllocatorOption{
		chromedp.NoFirstRun,
		chromedp.NoDefaultBrowserCheck,
		chromedp.Headless,
		chromedp.DisableGPU,
		chromedp.NoSandbox,
		chromedp.UserAgent(opts.UserAgent),
		// runner.Flag("--show-paint-rects", true),
	}

	if opts.IgnoreCertificateErrors {
		rets = append(rets, func(a *chromedp.ExecAllocator) {
			chromedp.Flag("ignore-certificate-errors", true)(a)
		})
	}

	if opts.Proxy != "" {
		rets = append(rets, chromedp.ProxyServer(opts.Proxy))
	}

	return rets
}

func NewChromedp(ctx context.Context) *ChromeBrowser {

	opts := &Opts{
		true,
		`Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/78.0.3904.97 Safari/537.36`,
		"",
	}
	return OpenBrowser(ctx,opts)
}
