package main

import (
	"bufio"
	"compress/bzip2"
	"compress/flate"
	"compress/gzip"
	"compress/lzw"
	"compress/zlib"
	"crypto/tls"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"gopkg.in/yaml.v2"
)

var (
	Order   []string
	Cache   = make(map[string]bool)
	Rules   = make(map[string]string)
	Default = &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}
	Download = "https://cdn.jsdelivr.net/gh/DongHuangT1/dry/dry.yaml"
)

func init() {
	log.SetFlags(0)
	log.SetPrefix("\033[1;31m[-]\033[0m ")

	usr, err := os.UserHomeDir()
	if err != nil {
		log.Fatalln(err)
	}

	cfg := filepath.Join(usr, ".config/dry/config.yaml")

	if _, err = os.Stat(cfg); os.IsNotExist(err) {
		os.MkdirAll(filepath.Dir(cfg), 0777)

		res, err := Default.Get(Download)
		if err != nil {
			log.Fatalln(err)
		}
		defer res.Body.Close()

		if res.StatusCode != 200 {
			log.Fatalln("download rules failed !!!")
		}

		rw, err := os.Create(cfg)
		if err != nil {
			log.Fatalln(err)
		}
		defer rw.Close()

		_, err = io.Copy(rw, res.Body)
		if err != nil {
			log.Fatalln(err)
		}

		fmt.Printf("\033[1;32m[+]\033[0m write rules to %s\n", cfg)
	}

	yml, err := os.ReadFile(cfg)
	if err != nil {
		log.Fatalln(err)
	}

	obj := make(map[string]string)

	err = yaml.Unmarshal(yml, obj)
	if err != nil {
		log.Fatalln(err)
	}

	for k, v := range obj {
		Order = append(Order, k)
		Rules[strings.ToLower(k)] = v
	}

	sort.Strings(Order)
}

func main() {
	Prefix := flag.String("0", "", "Set output prefix")
	Suffix := flag.String("9", "", "Set output suffix")
	Compress := flag.String("d", "", "Decompress (lzw, gzip, zlib, bzip2, flate)")
	flag.Usage = func() {
		fmt.Println("Usage: dry [-options] [- | url | file] [rules...]")
		flag.PrintDefaults()
		fmt.Printf("\n\033[0;32mSupport:\033[0m %s", strings.Join(Order, " "))
		os.Exit(0)
	}
	flag.Parse()

	if flag.NArg() < 2 {
		flag.Usage()
	}

	var (
		err error
		r1  io.ReadCloser
		r2  io.ReadCloser
		rdb []*regexp.Regexp
	)

	for _, v := range flag.Args()[1:] {
		if r := Rules[strings.ToLower(v)]; r != "" {
			rdb = append(rdb, regexp.MustCompile(r))
		}
	}

	if n := flag.Arg(0); n == "-" {
		r1 = os.Stdin
	} else if strings.Contains(n, "://") {
		r1, err = func() (io.ReadCloser, error) {
			res, err := Default.Get(n)
			if err != nil {
				return nil, err
			}
			return res.Body, nil
		}()
	} else {
		r1, err = os.Open(n)
	}
	if err != nil {
		log.Fatalln(err)
	}
	defer r1.Close()

	switch *Compress {
	case "lzw":
		r2 = lzw.NewReader(r1, lzw.MSB, 8)
	case "gzip":
		r2, err = gzip.NewReader(r1)
	case "zlib":
		r2, err = zlib.NewReader(r1)
	case "bzip2":
		r2 = io.NopCloser(bzip2.NewReader(r1))
	case "flate", "deflate":
		r2 = flate.NewReader(r1)
	default:
		r2 = r1
	}
	if err != nil {
		log.Fatalln(err)
	}
	defer r2.Close()

	Scanner := bufio.NewScanner(r2)
	for Scanner.Scan() {
		Line := Scanner.Text()
		for _, r := range rdb {
			for _, v := range r.FindAllString(Line, -1) {
				if !Cache[v] {
					Cache[v] = true
					fmt.Println(*Prefix + v + *Suffix)
				}
			}
		}
	}
}
