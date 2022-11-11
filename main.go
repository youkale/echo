package main

import (
	"encoding/json"
	"flag"
	"io"
	"log"
	"math/rand"
	"net/http"
	"strings"
	"time"
)

type Request struct {
	URL    string              `json:"url"`
	Query  string              `json:"query"`
	Header map[string][]string `json:"header"`
	Method string              `json:"method"`
	Host   string              `json:"host"`
	Body   string              `json:"body"`
}

var logger = log.Default()

var addr string

func init() {
	flag.StringVar(&addr, "l", ":6111", "listen address")
}

var (
	b = [][]byte{
		[]byte(`{"rand":"error","message":"error"}`),
		[]byte(`<soapenv:Envelope xmlns:xsi=\"http://www.w3.org/2001/XMLSchema-instance\" xmlns:xsd=\"http://www.w3.org/2001/XMLSchema\" xmlns:soapenv=\"http://schemas.xmlsoap.org/soap/envelope/\" xmlns:web=\"http://webservice.ws.util.srt.eas.kingdee.com\"><soapenv:Header/><soapenv:Body><web:execute soapenv:encodingStyle=\"http://schemas.xmlsoap.org/soap/encoding/\"><param xsi:type=\"xsd:string\">{\"bizType\":\"0\",\"bill\":{\"columns\":[\"FID\",\"FLASTUPDATETIME\",\"FNUMBER\",\"FNAME_L2\",\"FSTATUS\",\"FMATERIALGROUPID\",\"FSHORTNAME\",\"FMODEL\",\"FPRICEPRECISION\",\"FHELPCODE\",\"FBARCODE\",\"FGROSSWEIGHT\",\"FNETWEIGHT\",\"FLENGTH\",\"FWIDTH\",\"FHEIGHT\",\"FVOLUME\",\"FBASEUNIT\",\"FWEIGHTUNIT\",\"FLENGTHUNIT\",\"FVOLUMNUNIT\",\"FALIAS\",\"FFOREIGNNAME\",\"FREGISTEREDMARK\",\"FWARRANTNUMBER\",\"FASSISTUNIT\",\"FOLDNUMBER\",\"CFTEAMID\",\"CFISCUSTOMIZED\",\"CFISDANGEROUS\",\"CFISBIGMATERIAL\",\"CFCOSTPRICE\",\"CFMAXPACKAGEQTY\",\"CFPRODUCTATTRIBUTE\",\"CFASCRIPTIONORGID\"],\"dataType\":[\"String\",\"String\",\"String\",\"String\",\"String\",\"String\",\"String\",\"String\",\"String\",\"String\",\"String\",\"String\",\"String\",\"String\",\"String\",\"String\",\"String\",\"String\",\"String\",\"String\",\"String\",\"String\",\"String\",\"String\",\"String\",\"String\",\"String\",\"String\",\"String\",\"String\",\"String\",\"String\",\"String\",\"String\",\"String\"],\"sql\":\"/*dialect*/select FID,FLASTUPDATETIME,FNUMBER,FNAME_L2,FSTATUS,FMATERIALGROUPID,FSHORTNAME,FMODEL,FPRICEPRECISION,FHELPCODE,FBARCODE,FGROSSWEIGHT,FNETWEIGHT,FLENGTH,FWIDTH,FHEIGHT,FVOLUME,FBASEUNIT,FWEIGHTUNIT,FLENGTHUNIT,FVOLUMNUNIT,FALIAS,FFOREIGNNAME,FREGISTEREDMARK,FWARRANTNUMBER,FASSISTUNIT,FOLDNUMBER,CFTEAMID,CFISCUSTOMIZED,CFISDANGEROUS,CFISBIGMATERIAL,CFCOSTPRICE,CFMAXPACKAGEQTY,CFPRODUCTATTRIBUTE,CFASCRIPTIONORGID from ( select  a.FID || to_char (a.fcreatetime,\\u0027yyyymmddhh24miss\\u0027) AS FID,a.FLASTUPDATETIME,a.FNUMBER,a.FNAME_L2,a.FSTATUS,b.fid || to_char (b.fcreatetime,\\u0027yyyymmddhh24miss\\u0027) AS FMATERIALGROUPID,a.FSHORTNAME,a.FMODEL,a.FPRICEPRECISION,a.FHELPCODE,a.FBARCODE,a.FGROSSWEIGHT,a.FNETWEIGHT,a.FLENGTH,a.FWIDTH,a.FHEIGHT,a.FVOLUME,a.FBASEUNIT,a.FWEIGHTUNIT,a.FLENGTHUNIT,a.FVOLUMNUNIT,a.FALIAS,a.FFOREIGNNAME,a.FREGISTEREDMARK,a.FWARRANTNUMBER,a.FASSISTUNIT,a.FOLDNUMBER,a.CFTEAMID,a.CFISCUSTOMIZED,a.CFISDANGEROUS,a.CFISBIGMATERIAL,a.CFCOSTPRICE,a.CFMAXPACKAGEQTY,a.CFPRODUCTATTRIBUTE,a.CFASCRIPTIONORGID, rownum as rn from T_BD_Material a LEFT JOIN t_bd_materialgroup b ON a.fmaterialgroupid \\u003d b.fid where  a.FLastUpdateTime between to_date(\\u00272022-02-28 10:45:53\\u0027,\\u0027yyyy-mm-dd hh24:mi:ss\\u0027) and to_date(\\u00272022-07-20 13:45:53\\u0027,\\u0027yyyy-mm-dd hh24:mi:ss\\u0027)  and rownum \\u003c\\u003d  200 order by a.FID||to_char(a.fcreatetime,\\u0027yyyymmddhh24miss\\u0027) asc ) t WHERE t.rn \\u003e\\u003d 1 order by fid asc\"}}</param><isDebug xsi:type=\"xsd:int\">0</isDebug></web:execute></soapenv:Body></soapenv:Envelope>`),
		[]byte(`rand/error`),
	}
)

func main() {
	flag.Parse()
	logger.Print("listen port:", addr)
	rand.Seed(time.Now().UnixNano())
	http.ListenAndServe(addr, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		p := r.RequestURI
		if strings.HasPrefix(p, "/rand") {
			i := rand.Intn(len(b))
			w.Write(b[i])
		} else {
			req := Request{
				URL:    p,
				Host:   r.Host,
				Method: r.Method,
				Header: r.Header,
				Query:  r.URL.RawQuery,
			}
			if nil != r.Body {
				if body, re := io.ReadAll(r.Body); nil == re {
					req.Body = string(body)
				}
			}

			if res, err := json.Marshal(req); nil == err {
				logger.Print("rec: " + string(res))
				w.Write(res)
			}
		}
	}))
}
