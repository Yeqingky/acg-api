package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"log"
	"math/rand/v2"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	_ "golang.org/x/image/webp"
)

const (
	addrDefault = ":8080"
	dataDir     = "data"

	widthFile  = "acg_横图.txt"
	heightFile = "acg_竖图.txt"
)

// mobileUA 与 PHP 版本判断移动端的正则保持一致.
var mobileUA = regexp.MustCompile(`(?i)Android|iP(hone|od)|Windows ?CE|Symbian|Mobile|Opera Mobi|BlackBerry|Palm(OS)?|PocketPC|SonyEricsson|Nokia|Vodafone|HTC_|SAMSUNG-SGH|armv[56]l|Go\.Web|J2ME/MIDP`)

// imageClient 用于查询图片尺寸; 超时避免慢速 CDN 拖死请求.
var imageClient = &http.Client{Timeout: 10 * time.Second}

// gallery 是一个图片列表文件, 按文件 mtime/size 自动重载.
type gallery struct {
	path string

	mu    sync.Mutex
	mtime time.Time
	size  int64
	urls  []string
}

type jsonDetail struct {
	Code   string `json:"code"`
	AcgURL string `json:"acgurl"`
	Width  string `json:"width"`
	Height string `json:"height"`
	Size   string `json:"size"`
}

func main() {
	log.SetFlags(log.LstdFlags)

	pcPath := flag.String("pc", filepath.Join(dataDir, widthFile), "横图地址列表文件路径")
	pePath := flag.String("pe", filepath.Join(dataDir, heightFile), "竖图地址列表文件路径")
	addr := flag.String("addr", addrDefault, "监听地址")
	phpCompat := flag.Bool("php-compat", false, "同时注册 .php 后缀的旧路径")
	flag.Parse()

	pc := &gallery{path: *pcPath}
	pe := &gallery{path: *pePath}

	for _, g := range []*gallery{pc, pe} {
		lines, err := g.reload()
		if err != nil {
			log.Fatalf("load %s: %v", g.path, err)
		}
		log.Printf("loaded %d urls from %s", len(lines), g.path)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/api", handleAPI)
	mux.HandleFunc("/pc", galleryHandler(pc))
	mux.HandleFunc("/pe", galleryHandler(pe))
	mux.HandleFunc("/test", handleTest)
	if *phpCompat {
		mux.HandleFunc("/api.php", handleAPI)
		mux.HandleFunc("/pc.php", galleryHandler(pc))
		mux.HandleFunc("/pe.php", galleryHandler(pe))
		mux.HandleFunc("/test.php", handleTest)
	}
	mux.HandleFunc("/", http.NotFound)

	log.Printf("listening on %s", *addr)
	srv := &http.Server{
		Addr:              *addr,
		Handler:           withCORS(mux),
		ReadHeaderTimeout: 10 * time.Second,
	}
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}
}

// handleAPI 按 User-Agent 把调用方送到横图或竖图接口.
func handleAPI(w http.ResponseWriter, r *http.Request) {
	target := "/pc"
	if mobileUA.MatchString(r.UserAgent()) {
		target = "/pe"
	}
	http.Redirect(w, r, target, http.StatusFound)
}

func handleTest(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	fmt.Fprint(w, "Hello World")
}

func withCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func galleryHandler(g *gallery) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		lines, err := g.lines()
		if err != nil {
			log.Printf("load %s: %v", g.path, err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		if len(lines) == 0 {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		img := lines[rand.IntN(len(lines))]

		switch r.URL.Query().Get("return") {
		case "json":
			writeJSON(w, img)
		case "all":
			w.Header().Set("Content-Type", "text/plain; charset=utf-8")
			w.Write([]byte(strings.Join(lines, "\n")))
		default:
			http.Redirect(w, r, img, http.StatusFound)
		}
	}
}

func writeJSON(w http.ResponseWriter, imgURL string) {
	detail := jsonDetail{
		Code:   "200",
		AcgURL: imgURL,
		Width:  "0",
		Height: "0",
		Size:   urlExtension(imgURL),
	}
	if cw, ch, err := imageSize(imgURL); err == nil {
		detail.Width = strconv.Itoa(cw)
		detail.Height = strconv.Itoa(ch)
	} else {
		log.Printf("image size %s: %v", imgURL, err)
	}

	body, err := json.Marshal(detail)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Write(body)
}

// imageSize 拉取远程图片并只解析头部, 得到宽高.
func imageSize(imgURL string) (int, int, error) {
	req, err := http.NewRequest(http.MethodGet, imgURL, nil)
	if err != nil {
		return 0, 0, err
	}
	req.Header.Set("User-Agent", "acg-api")

	resp, err := imageClient.Do(req)
	if err != nil {
		return 0, 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, 0, fmt.Errorf("unexpected status %s", resp.Status)
	}

	cfg, _, err := image.DecodeConfig(resp.Body)
	if err != nil {
		return 0, 0, err
	}
	return cfg.Width, cfg.Height, nil
}

// urlExtension 复刻 PHP pathinfo($url)['extension'].
func urlExtension(raw string) string {
	if u, err := url.Parse(raw); err == nil && u.Path != "" {
		return strings.TrimPrefix(path.Ext(u.Path), ".")
	}
	return strings.TrimPrefix(path.Ext(raw), ".")
}

// lines 返回当前列表; 文件被替换后自动重载.
func (g *gallery) lines() ([]string, error) {
	info, err := os.Stat(g.path)
	if err != nil {
		return nil, err
	}

	g.mu.Lock()
	defer g.mu.Unlock()

	if !info.ModTime().Equal(g.mtime) || info.Size() != g.size {
		if _, err := g.reloadLocked(info); err != nil {
			return nil, err
		}
	}
	return g.urls, nil
}

func (g *gallery) reload() ([]string, error) {
	info, err := os.Stat(g.path)
	if err != nil {
		return nil, err
	}

	g.mu.Lock()
	defer g.mu.Unlock()
	return g.reloadLocked(info)
}

func (g *gallery) reloadLocked(info os.FileInfo) ([]string, error) {
	data, err := os.ReadFile(g.path)
	if err != nil {
		return nil, err
	}

	lines := make([]string, 0, 1024)
	for _, line := range strings.Split(string(data), "\n") {
		if line = cleanURL(line); line != "" {
			lines = append(lines, line)
		}
	}

	g.urls = lines
	g.mtime = info.ModTime()
	g.size = info.Size()
	return lines, nil
}

// cleanURL 去掉所有空白字符, 与 PHP 版本的 cleanUrl 一致.
func cleanURL(s string) string {
	return strings.Map(func(r rune) rune {
		switch r {
		case ' ', '\n', '\r', '\t':
			return -1
		}
		return r
	}, s)
}
