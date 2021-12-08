package main

import (
	"context"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
)

type ElectronicPre struct {
	Rid               int    //处方编号
	GetTitleName      string //
	CanLiuzhuan       int
	Patient           Patient
	Diagnosis         Diagnosis //门诊ID号/住院号
	Owner             Owner     //科别
	GetOpinionNames   string    //临床诊断
	CreatedTime       string    //开具日期
	Items             []map[string]interface{}
	Type              string
	GetZhongyaoUsages string
	Mode              string
}

type Patient struct {
	Name             string //
	GetGenderDisplay string
	GetAgeByIdno     int
	Address          string //住址
	Mobile           int    //电话

}
type Diagnosis struct {
	Rid int
}
type Owner struct {
	Department Department
}
type Department struct {
	Department int
	Name       string
}

func prescription(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	//接收处方id
	id, err := strconv.Atoi(query.Get("id"))
	//根据处方id去查数据库获取处方信息




	//定义好模板文件的路径
	htmlByte, err := ioutil.ReadFile("./prescription.html")
	if err != nil {
		fmt.Println("read html failed, err:", err)
		return
	}
	// 采用链式操作在Parse之前调用Funcs
	tmpl, err := template.New("prescription").Parse(string(htmlByte))
	if err != nil {
		fmt.Println("create template failed, err:", err)
		return
	}

	//数据库查询到的结果赋值给结构体渲染html
	user := ElectronicPre{
		Rid:          id,
		GetTitleName: "测试",
		CanLiuzhuan:  1,
		Patient: Patient{
			Name:             "测试1111",
			GetGenderDisplay: "男",
			GetAgeByIdno:     18,
			Address:          "广州市",
			Mobile:           15845698745,
		},
		Diagnosis: Diagnosis{
			Rid: 1111,
		},
		Owner: Owner{
			Department: Department{
				Department: 1,
				Name:       "部门名称",
			},
		},
		GetOpinionNames: "测试临床诊断",
		CreatedTime:     "2021.12.32",
		Items: []map[string]interface{}{
			{"Name": "noob", "Usage": 18, "Frequency": 1, "GetDosage": 10, "GetDosageUnit": 1},
			{"Name": "noob1", "Usage": 118, "Frequency": 11, "GetDosage": 110, "GetDosageUnit": 11},
		},

		Type:              "1",
		GetZhongyaoUsages: "测试1312",
		Mode:              "wx",
	}
	// 使用user渲染模板，并将结果写入w
	tmpl.Execute(w, user)
}

func ChromedpPrintPdf(url string, to string) error {
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()
	var buf []byte
	err := chromedp.Run(ctx, chromedp.Tasks{
		chromedp.Navigate(url),
		chromedp.WaitReady("body"),
		chromedp.ActionFunc(func(ctx context.Context) error {
			var err error
			buf, _, err = page.PrintToPDF().WithPaperHeight(15).WithPaperWidth(8.5).//这里调整纸张大小
				Do(ctx)
			return err
		}),
	})
	if err != nil {
		return fmt.Errorf("chromedp Run failed,err:%+v", err)
	}
	if err := ioutil.WriteFile(to, buf, 0644); err != nil {
		return fmt.Errorf("write to file failed,err:%+v", err)
	}
	return nil
}


func ToPdf(id string) (pdfpath string, err error) {

	timeUnixNano := time.Now().UnixNano()
	pdfpath = string(timeUnixNano)+id

	errs := ChromedpPrintPdf("http://127.0.0.1:8080/prescription?id="+id, pdfpath+".pdf")
	if errs != nil {
		fmt.Println(errs)
	}
	return
}
func main() {
	server := http.Server{
		Addr: "127.0.0.1:8080",
	}
	//渲染处方html路由
	http.HandleFunc("/prescription", prescription)
	log.Fatal(server.ListenAndServe())

}
