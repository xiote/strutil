package namevalue

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/pkg/errors"

	"github.com/antchfx/htmlquery"
	iconv "github.com/djimenez/iconv-go"
)

type NameValues struct {
	Names  []string
	Values url.Values
}

func NewNameValues() NameValues {
	return NameValues{
		[]string{},
		url.Values{},
	}
}

type NameValue struct {
	Name  string
	Value string
}

type IdValue struct {
	Id    string
	Value string
}

func (c *NameValues) Set(name string, value string) {
	c.Values.Set(name, value)
	return
}

func (c *NameValues) Add(name string, value string) {
	c.Names = append(c.Names, name)
	c.Values.Add(name, value)
	return
}

func (c *NameValues) Get(name string) string {
	return c.Values.Get(name)
}

func (c *NameValues) UrlEncode() string {

	var paramsString string
	for _, name := range c.Names {
		paramsString += fmt.Sprintf("%s=%s&", name, url.QueryEscape(c.Values.Get(name)))
	}
	paramsString = strings.TrimSuffix(paramsString, "&")
	// paramsString = strings.ReplaceAll(paramsString, "%2B", "+")
	paramsString = strings.ReplaceAll(paramsString, "~", "%7E")
	return paramsString
}

func ExtractNameValues(html string, xpath string) (values NameValues, err error) {
	doc, err := htmlquery.Parse(strings.NewReader(html))
	if err != nil {
		err = errors.Wrap(err, "Parse failed")
		return
	}

	list, err := htmlquery.QueryAll(doc, xpath)
	if err != nil {
		err = errors.Wrap(err, "QueryAll failed")
		return
	}

	if len(list) == 0 {
		err = errors.Wrap(err, fmt.Sprintf("Cannot find node : %s", xpath))
		return
	}

	values = NewNameValues()
	for _, input := range list {
		var v string
		v, err = iconv.ConvertString(htmlquery.SelectAttr(input, "value"), "utf-8", "euc-kr")
		if err != nil {
			err = errors.Wrap(err, "ConvertString failed")
			return
		}
		values.Add(htmlquery.SelectAttr(input, "name"), v)
	}

	return
}

func GetNameValue(nameValues []NameValue, name string) (nameValue NameValue) {
	for _, item := range nameValues {
		if item.Name == name {
			nameValue = item
			return
		}
	}
	return
}

func GetValue(nameValues []NameValue, name string) (value string) {
	for _, item := range nameValues {
		if item.Name == name {
			value = item.Value
			return
		}
	}
	return
}

func UpdateNameValue(ptNameValues *[]NameValue, name string, value string) (err error) {
	for idx, item := range *ptNameValues {
		if item.Name == name {
			(*ptNameValues)[idx].Value = value
			break
		}
	}

	err = errors.Wrap(err, "NameValue not found!")
	return
}

func NameValuesToUrlValues(nameValues []NameValue) (v url.Values) {
	for _, item := range nameValues {
		v.Add(item.Name, item.Value)
	}

	return
}

func ExtractNameValue(html string, xpath string) (nameValues []NameValue, err error) {
	doc, err := htmlquery.Parse(strings.NewReader(html))
	if err != nil {
		err = errors.Wrap(err, "Parse failed")
		return
	}

	list, err := htmlquery.QueryAll(doc, xpath)
	if err != nil {
		err = errors.Wrap(err, "QueryAll failed")
		return
	}

	if len(list) == 0 {
		err = errors.Wrap(err, fmt.Sprintf("Cannot find node : %s", xpath))
		return
	}

	for _, input := range list {
		var v string
		v, err = iconv.ConvertString(htmlquery.SelectAttr(input, "value"), "utf-8", "euc-kr")
		if err != nil {
			err = errors.Wrap(err, "ConvertString failed")
			return
		}
		nameValues = append(nameValues, NameValue{htmlquery.SelectAttr(input, "name"), v})
	}

	return
}

func ExtractIdValueUtf8(html string, xpath string) (idValues []IdValue, err error) {
	doc, err := htmlquery.Parse(strings.NewReader(html))
	if err != nil {
		err = err
		return
	}

	list, err := htmlquery.QueryAll(doc, xpath)
	if err != nil {
		err = err
		return
	}

	if len(list) == 0 {
		err = fmt.Errorf("Cannot find node : %s", xpath)
		return
	}

	for _, input := range list {
		idValues = append(idValues, IdValue{htmlquery.SelectAttr(input, "id"), htmlquery.SelectAttr(input, "value")})
	}
	return

}

func GetQueryString(params []NameValue) string {

	var paramsString string
	for _, input := range params {
		paramsString += fmt.Sprintf("%s=%s&", input.Name, url.QueryEscape(input.Value))
	}
	paramsString = strings.TrimSuffix(paramsString, "&")
	// paramsString = strings.ReplaceAll(paramsString, "%2B", "+")
	paramsString = strings.ReplaceAll(paramsString, "~", "%7E")
	return paramsString
}
