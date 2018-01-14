package main

import (
	"flag"
	"path"
	"strings"
	"log"
	"net/http"
	"net"
	"io"
	"io/ioutil"
	"os"
	"fmt"
	"time"
	"html/template"
	"strconv"
)

const (
	START = "\x1b[33mStarting up http-server\nAvailable on:\x1b[0m\n  http://127.0.0.1:\x1b[32m%s\x1b[0m\n  http://192.168.0.101:\x1b[32m%[1]s\x1b[0m\nHit CTRL-C to stop the server\n"
	REQUEST = "[%s] \"\x1b[36m%s\x1b[0m\" %q\n"
)

type file struct {
	Name string
	IsDir bool
	Size string
	ModTime string
}

type data struct {
	List []file
	Pathname string
}

var dir string
var port string

func init() {
	flag.StringVar(&port, "port", "", "set the port")
	flag.StringVar(&dir, "dir", "", "set the directory")
}

func main() {
	flag.Parse()
	if dir == "" {
		dir = flag.Arg(0)
	}
	http.HandleFunc("/", handle)
	var err error
	if port == "" {
		port, err = getFreePort()
	}
	if err != nil {
		log.Fatal(err)
		return
	}
	fmt.Printf(START, port)
	log.Fatal(http.ListenAndServe(":" + port, nil))
	if err != nil {
		log.Fatal(err)
		return
	}
}

func handle(w http.ResponseWriter, req *http.Request) {
	if req.Method != "GET" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	pathname := req.URL.Path
	fmt.Printf(REQUEST, time.Now().Format("Mon Jan 02 2006 15:04:05 GMT-0700 (MST)"), req.Method + " " + pathname, req.Header.Get("user-agent"))
	workDir, err := os.Getwd()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	var fullpath string
	if dir != "" {
		if path.IsAbs(dir) {
			fullpath = path.Join(dir, pathname)
		} else {
			fullpath = path.Join(workDir, dir, pathname)
		}
	} else {
		fullpath = path.Join(workDir, pathname)
	}
	fileInfo, err := os.Stat(fullpath)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	if fileInfo.IsDir() {
		files, err := ioutil.ReadDir(fullpath)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		list := make([]file, len(files))
		for i, item := range files {
			name := item.Name()
			if item.IsDir() {
				name += "/"
			}
			list[i] = file{
				Name: name,
				Size: formatSize(item.Size()),
				ModTime: item.ModTime().Format("2006-01-02 15:04:05"),
				IsDir: item.IsDir()}
		}
		tplData := data{
			Pathname: pathname,
			List: list}
		t, err := template.New("tpl").Parse(directory)
		if err != nil {
			fmt.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		err = t.Execute(w, tplData)
		if err != nil {
			fmt.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	} else {
		file, err := os.Open(fullpath)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if _, err := io.Copy(w, file); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
}

func formatSize(size int64) string {
	sizeMap := [...]string{"B", "k", "M", "G", "T"}
	str := ""
	for i := len(sizeMap) - 1; i >= 0; i-- {
		n := int64(1 << (10 * uint(i)))
		if size >= n {
			str = strconv.FormatFloat(float64(size) / float64(n), 'f', 1, 64) + sizeMap[i]
			break
		}
	}
	return str
}

func getFreePort() (string, error) {
	l, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		return "", err
	}
	defer l.Close()
	addr := l.Addr().String()
	return strings.Split(addr, ":")[1], nil
}