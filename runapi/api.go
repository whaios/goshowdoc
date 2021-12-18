package runapi

import (
	"fmt"

	"github.com/levigross/grequests"
)

var SSID = "" // ShowDoc用户的SessionId，浏览器登录后从Cookie中获取PHPSESSID。

// PageSave 保存接口文档。
func PageSave(page *Page) error {
	url := fmt.Sprintf("%s/server/?s=/api/page/save", Host)

	result := ErrResult{}
	if err := post(url, page.ToMap(), &result); err != nil {
		return err
	}
	return result.Error()
}

// PageInfo 接口文档详情。
func PageInfo(pageId string) (*Page, error) {
	url := fmt.Sprintf("%s/server/?s=/api/page/info", Host)

	result := struct {
		ErrResult
		Data *Page `json:"data"`
	}{}
	if err := post(url, map[string]string{"page_id": pageId}, &result); err != nil {
		return nil, err
	}
	return result.Data, result.Error()
}

// ItemInfo 项目目录列表。
func ItemInfo(itemId string) (*Item, error) {
	url := fmt.Sprintf("%s/server/?s=/api/item/info", Host)

	result := struct {
		ErrResult
		Data *Item `json:"data"`
	}{}
	if err := post(url, map[string]string{"item_id": itemId}, &result); err != nil {
		return nil, err
	}
	return result.Data, result.Error()
}

// CatalogSave 保存目录。
func CatalogSave(itemId, parentCatId, catName string) (string, error) {
	url := fmt.Sprintf("%s/server/?s=/api/catalog/save", Host)

	result := struct {
		ErrResult
		Id string `json:"data"`
	}{}
	if err := post(url, map[string]string{
		"item_id":       itemId,
		"parent_cat_id": parentCatId,
		"cat_name":      catName,
	}, &result); err != nil {
		return "", err
	}
	return result.Id, result.Error()
}

func post(url string, data map[string]string, result interface{}) error {
	resp, err := grequests.Post(url, &grequests.RequestOptions{
		Data: data,
		Headers: map[string]string{
			"Cookie": "PHPSESSID=" + SSID,
		},
	})
	if err != nil {
		return err
	}
	return resp.JSON(result)
}
