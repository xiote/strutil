package httputil

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"

	"github.com/pkg/errors"

	log "github.com/xiote/go-utils/chanlog"
	"golang.org/x/sync/errgroup"
	"golang.org/x/text/encoding/korean"
	tf "golang.org/x/text/transform"
)

func GetCookieValue(jar http.CookieJar, rawUrl string, cookieName string) (value string, err error) {

	u, err := url.Parse(rawUrl)
	if err != nil {
		err = errors.Wrap(err, "Parse failed")
		return
	}

	for _, cookie := range jar.Cookies(u) {
		if cookie.Name == cookieName {
			value = cookie.Value
			return
		}
	}

	err = errors.Wrap(err, fmt.Sprintf("%s is  not found!", cookieName))
	return

}

func EuckrDo2(client *http.Client, req *http.Request, nameforlog string) (respdate time.Time, src string, err error) {
	// var starttime time.Time
	var body []byte
	// var elaspedHttpDuration time.Duration
	// go func() {
	// 	starttime = time.Now()
	// 	go log.Printf("[%s] [START]\n", nameforlog)
	// }()
	// defer func() {
	// 	go func() {
	// 		l := int64(len(body) * 8)
	// 		log.Printf("[%s] [END] [ %s ] [ %s ] [ %d Bits ] [ %d Kbps ]\n", nameforlog, time.Since(starttime), elaspedHttpDuration, l, l*1000000/elaspedHttpDuration.Nanoseconds())
	// 	}()
	// }()

	if respdate, body, err = do2(client, req); err != nil {
		return
	}
	// elaspedHttpDuration = time.Since(starttime)

	{
		var bufs bytes.Buffer
		wr := tf.NewWriter(&bufs, korean.EUCKR.NewDecoder())
		defer wr.Close()
		wr.Write(body)
		src = bufs.String()

	}
	return
}

func EuckrDo(client *http.Client, req *http.Request, nameforlog string) (src string, err error) {
	var starttime time.Time
	var body []byte
	var elaspedHttpDuration time.Duration
	go func() {
		starttime = time.Now()
		go log.Printf("[%s] [START]\n", nameforlog)
	}()
	defer func() {
		go func() {
			l := int64(len(body) * 8)
			log.Printf("[%s] [END] [ %s ] [ %s ] [ %d Bits ]\n", nameforlog, time.Since(starttime), elaspedHttpDuration, l)
		}()
	}()

	if body, err = do(client, req); err != nil {
		return
	}
	go func() { elaspedHttpDuration = time.Since(starttime) }()

	var bufs bytes.Buffer
	wr := tf.NewWriter(&bufs, korean.EUCKR.NewDecoder())
	defer wr.Close()
	wr.Write(body)

	src = bufs.String()
	return
}

func DoWithoutLog(client *http.Client, req *http.Request) (src string, err error) {
	var body []byte
	if body, err = do(client, req); err != nil {
		return
	}
	src = string(body)
	return
}

func Do(client *http.Client, req *http.Request, nameforlog string) (src string, err error) {
	var starttime time.Time
	var body []byte
	var elaspedHttpDuration time.Duration
	go func() {
		starttime = time.Now()
		go log.Printf("[%s] [START]\n", nameforlog)
	}()
	defer func() {
		go func() {
			l := int64(len(body) * 8)
			elaspedHttpDuration = time.Since(starttime)
			log.Printf("[%s] [END] [ %s ] [ %s ] [ %d Bits ]\n", nameforlog, time.Since(starttime), elaspedHttpDuration, l)
		}()
	}()

	if body, err = do(client, req); err != nil {
		return
	}
	// go func() { elaspedHttpDuration = time.Since(starttime) }()

	src = string(body)
	return
}

func do(client *http.Client, req *http.Request) (body []byte, err error) {
	var resp *http.Response
	if resp, err = client.Do(req); err != nil {
		return
	}
	defer resp.Body.Close()

	var reader io.ReadCloser
	if reader, err = ContentDecodingReader(resp.Header.Get("Content-Encoding"), resp.Body); err != nil {
		return
	}

	if body, err = ioutil.ReadAll(reader); err != nil {
		return
	}

	return
}

func do2(client *http.Client, req *http.Request) (respdate time.Time, body []byte, err error) {
	var resp *http.Response
	if resp, err = client.Do(req); err != nil {
		return
	}
	defer resp.Body.Close()

	g := new(errgroup.Group)
	g.Go(func() (err error) {
		var reader io.ReadCloser
		if reader, err = ContentDecodingReader(resp.Header.Get("Content-Encoding"), resp.Body); err != nil {
			return
		}

		if body, err = ioutil.ReadAll(reader); err != nil {
			return
		}
		return
	})
	g.Go(func() (err error) {
		if respdate, err = http.ParseTime(http.Header.Get(resp.Header, "Date")); err != nil {
			return
		}
		return
	})
	err = g.Wait()

	return
}

func ContentDecodingReader(contentEncoding string, body io.ReadCloser) (reader io.ReadCloser, err error) {

	switch contentEncoding {
	case "gzip":
		reader, err = gzip.NewReader(body)
		if err != nil {
			return
		}
		defer reader.Close()
	default:
		reader = body
	}
	return
}
