package main

import (
	"log"
	"net/http"
	"os/exec"
	"runtime"
	"time"

	clipboard "github.com/atotto/clipboard"
	qrcode "github.com/skip2/go-qrcode"
	hotkey "golang.design/x/hotkey"
)

func generator(url string) ([]byte, error) {
	var png []byte
	png, err := qrcode.Encode(url, qrcode.High, 256)
	return png, err
}

func show(qr []byte) error {
	http.HandleFunc("/qr", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "image/png")
		w.Write(qr)
	})
	server := &http.Server{Addr: "localhost:8080"}
	go func() {
		if err := server.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatal("Server error:", err)
		}
	}()
	if err := openBrowser("http://localhost:8080/qr"); err != nil {
		server.Close()
		log.Println("Impossible to open a browser: ", err)
		return err
	}
	time.Sleep(30 * time.Second)
	if err := server.Close(); err != nil {
		log.Println("Impossible to close server correctly: ", err)
	}
	return nil
}

func openBrowser(url string) error {
	var (
		cmd  string
		args []string
	)
	switch runtime.GOOS {
	case "windows":
		cmd = "cmd"
		args = []string{"/c", "start", url}
	case "darwin":
		cmd = "open"
		args = []string{url}
	default: // linux, freebsd и др.
		cmd = "xdg-open"
		args = []string{url}
	}
	return exec.Command(cmd, args...).Start()
}

func main() {
	hk := hotkey.New([]hotkey.Modifier{hotkey.ModCtrl}, hotkey.KeyF1)
	if err := hk.Register(); err != nil {
		log.Println("HotKey error:", err)
	}
	for range hk.Keydown() {
		data, err := clipboard.ReadAll()
		if err != nil {
			log.Println(err)
		}
		qr, err := generator(data)
		if err != nil {
			log.Println(err)
		}
		if err := show(qr); err != nil {
			log.Println(err)
		}
		time.Sleep(100 * time.Millisecond)
	}
}
