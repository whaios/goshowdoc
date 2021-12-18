package main

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	"github.com/whaios/goshowdoc/log"
	"github.com/whaios/goshowdoc/runapi"
)

// ItemInfo 显示项目详情。
func ItemInfo(itemId string) {
	item, err := runapi.ItemInfo(itemId)
	if err != nil {
		log.Error(err.Error())
		return
	}
	log.Info("item.id=%s", item.ItemId)
	log.Info("item.name=%s", item.ItemName)
	log.Info("item.pages=%d", len(item.MenuPages()))
}

type ExportData struct {
	Item  *runapi.Item   `json:"item"`
	Pages []*runapi.Page `json:"pages"`
}

// ItemExport 导出项目到指定目录。
func ItemExport(itemId, out string) {
	// 获取项目信息，包含所有目录
	item, err := runapi.ItemInfo(itemId)
	if err != nil {
		log.Error(err.Error())
		return
	}
	menuPages := item.MenuPages()
	log.Info("开始导出 %s", item.ItemName)

	// 获取项目下的所有文档
	pages := make([]*runapi.Page, 0, len(menuPages))
	{
		for i, mp := range menuPages {
			page, err := runapi.PageInfo(mp.PageId)
			if err != nil {
				log.Error(err.Error())
				return
			}
			pages = append(pages, page)
			log.DrawProgressBar("获取文档", i+1, len(menuPages))
		}
	}

	bytes, _ := json.Marshal(&ExportData{Item: item, Pages: pages})

	fileName := fmt.Sprintf("%s_goshowdoc_%s.json", item.ItemName, time.Now().Format("20060102150405"))
	fileName = filepath.Join(out, fileName)

	log.Info("导出到文件 %s", fileName)
	if err := mkdir(out); err != nil {
		log.Error(err.Error())
		return
	}
	if err = ioutil.WriteFile(fileName, bytes, fs.ModePerm); err != nil {
		log.Error(err.Error())
		return
	}
	log.Success("导出完成")
}

// ItemImport 导入项目
func ItemImport(itemId, fileName string) {
	// 获取项目信息，包含所有目录
	item, err := runapi.ItemInfo(itemId)
	if err != nil {
		log.Error(err.Error())
		return
	}

	bytes, err := ioutil.ReadFile(fileName)
	if err != nil {
		log.Error(err.Error())
		return
	}

	data := &ExportData{}
	json.Unmarshal(bytes, data)
	log.Info("开始导入 %s > %s", data.Item.ItemName, item.ItemName)

	catMap := make(map[string]*runapi.Catalog) // key=旧的目录id
	{
		for _, cat := range data.Item.Catalogs() {
			catMap[cat.CatId] = cat
		}
	}
	importCatTotal = len(catMap)
	importedCatCount = 0
	if ok := importCatalog(itemId, "", data.Item.Menu.Catalogs); !ok {
		log.Error("导入失败")
		return
	}

	pagesTotal := len(data.Pages)
	for i, page := range data.Pages {
		page.ItemId = itemId
		if cat, ok := catMap[page.CatId]; ok {
			page.CatId = cat.CatId
		} else {
			page.CatId = "" // 没有找到所属目录时，添加到根部
		}
		page.PageId = "" // 新建不需要id
		if err = runapi.PageSave(page); err != nil {
			log.Error("导入文档[%s]失败: %s", page.PageTitle, err.Error())
			return
		}
		log.DrawProgressBar("导入文档", i+1, pagesTotal)
	}
	log.Success("导入完成")
}

var (
	importCatTotal   int
	importedCatCount int
)

// 导入目录，并更新目录id为新的id
func importCatalog(itemId, parentCatId string, catalogs []*runapi.Catalog) bool {
	for _, cat := range catalogs {
		newid, err := runapi.CatalogSave(itemId, parentCatId, cat.CatName)
		if err != nil {
			log.Error("导入目录[%s]失败: %s", cat.CatName, err.Error())
			return false
		}
		if newid == "" {
			log.Error("导入目录[%s]失败: 没有获取到服务端返回的目录id", cat.CatName)
		}
		importedCatCount++
		log.DrawProgressBar("导入目录", importedCatCount, importCatTotal)

		cat.CatId = newid

		if len(cat.Catalogs) > 0 {
			importCatalog(itemId, cat.CatId, cat.Catalogs)
		}
	}
	return true
}

func mkdir(dir string) error {
	_, err := os.Stat(dir)
	if os.IsNotExist(err) {
		return os.MkdirAll(dir, os.ModePerm)
	}
	return err
}
