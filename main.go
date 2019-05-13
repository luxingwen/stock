package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/axgle/mahonia"
)

type SinaStock struct {
	Name     string
	Price    float64
	Percent  float64
	TickSize string
}

var format = flag.String("f", "term", "formatting: term or vim")

func init() {
	flag.Parse()
}

func main() {
	var (
		sname string
		err   error
	)
	if len(flag.Args()) == 1 {
		sname = flag.Args()[0]
	} else {
		sname, err = getStockList(GetCurPath() + "/" + "stock.list")
		if err != nil {
			fmt.Println("Usage: stock [symbol], e.g. stock sh600271")
			os.Exit(1)
		}
	}
	if sname == "" {
		fmt.Println("Usage: stock [symbol], e.g. stock sh600271")
		os.Exit(1)
	}

	rList, err := GetSinaStock(sname)
	if err != nil {
		log.Fatal(err)
	}

	res := ""
	if *format == "term" {
		for _, item := range rList {
			res = formatForTerminal(item.Name, item.Price, item.Percent, item.TickSize)
			fmt.Println(res)
		}
	} else {
		item := rList[0]
		res = formatForVim(item.Price, item.Percent, item.TickSize)
		fmt.Println(res)
	}
}

func getStockList(filename string) (r string, err error) {
	f, err := os.Open(filename)
	if err != nil {
		return
	}
	defer f.Close()
	buf := bufio.NewReader(f)
	lines := make([]string, 0)
	for {
		line, err := buf.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				if strings.TrimSpace(line) != "" {
					lines = append(lines, strings.TrimSpace(line))
				}
				break
			}
			return "", err
		}
		lines = append(lines, strings.TrimSpace(line))
	}
	r = strings.Join(lines, ",")
	return
}

func GetSinaStock(sname string) (list []*SinaStock, err error) {
	urlAdress := fmt.Sprintf("http://hq.sinajs.cn/list=%s", sname)
	req, err := http.Get(urlAdress)
	if err != nil {
		return
	}
	defer req.Body.Close()
	b, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return
	}
	enc := mahonia.NewDecoder("gbk")
	content := enc.ConvertString(string(b))
	list = make([]*SinaStock, 0)
	for _, item := range strings.Split(content, ";") {
		if strings.TrimSpace(item) == "" {
			continue
		}
		index := strings.Index(item, "\"")
		rs := strings.Split(item[index+1:len(item)-1], ",")
		price, _ := strconv.ParseFloat(rs[3], 10)
		oldPrice, _ := strconv.ParseFloat(rs[2], 10)
		zdfv := price - oldPrice
		zdf := zdfv / oldPrice * 100
		sinaStock := &SinaStock{
			Name:     rs[0],
			Price:    price,
			Percent:  zdfv,
			TickSize: fmt.Sprintf("(%.2f%%)", zdf),
		}
		list = append(list, sinaStock)
	}
	return
}

func GetCurPath() string {
	file, _ := exec.LookPath(os.Args[0])
	path, _ := filepath.Abs(file)
	rst := filepath.Dir(path)
	return rst
}

func formatForTerminal(name string, price float64, delta float64, percentage string) string {
	var deltaFormatted string
	var percentageFormatted = percentage

	if delta > 0 {
		deltaFormatted = fmt.Sprintf("\x1b[31m+%.2f\x1b[0m", delta)
		percentageFormatted = fmt.Sprintf("\x1b[31m%s\x1b[0m", percentage)
	} else if delta < 0 {
		deltaFormatted = fmt.Sprintf("\x1b[32m%.2f\x1b[0m", delta)
		percentageFormatted = fmt.Sprintf("\x1b[32m%s\x1b[0m", percentage)
	} else {
		deltaFormatted = fmt.Sprintf("%.2f", delta)
	}
	return fmt.Sprintf("\x1b[32m%s\x1b[0m %.2f %s %s", name, price, deltaFormatted, percentageFormatted)
}

func formatForVim(price float64, delta float64, percentage string) string {
	var result string
	if delta > 0 {
		result = fmt.Sprintf(`echohl Normal
echo "%.2f"
echohl MoreMsg
echon " +%.2f %s"
echohl Normal`, price, delta, percentage)

	} else if delta < 0 {
		result = fmt.Sprintf(`echohl Normal
echo "%.2f"
echohl WarningMsg
echon " %.2f %s"
echohl Normal`, price, delta, percentage)
	} else {
		result = fmt.Sprintf(`echohl Normal
echo "%.2f"
echon " %.2f %s"
echohl Normal`, price, delta, percentage)
	}
	return result
}
