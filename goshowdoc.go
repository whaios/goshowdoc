package main

import (
	"os"

	"github.com/urfave/cli/v2"
	"github.com/whaios/goshowdoc/log"
	"github.com/whaios/goshowdoc/runapi"
)

const (
	GOSHOWDOC_HOST     = "GOSHOWDOC_HOST"
	GOSHOWDOC_APIKEY   = "GOSHOWDOC_APIKEY"
	GOSHOWDOC_APITOKEN = "GOSHOWDOC_APITOKEN"
)

const (
	flagDir    = "dir"
	flagItemId = "itemid"
	flagFile   = "file"
)

func main() {
	app := cli.NewApp()
	app.Name = "goshowdoc"
	app.Usage = "ShowDoc API 接口文档工具"
	app.Description = `项目地址： https://github.com/whaios/goshowdoc
支持以下功能：
1. 通过代码注释生成 API 接口文档。
2. 导出和导入 ShowDoc 项目。`
	app.Version = Version

	app.Flags = []cli.Flag{
		&cli.StringFlag{
			Name:        "host",
			Value:       "https://www.showdoc.com.cn",
			Usage:       "ShowDoc 地址。",
			EnvVars:     []string{GOSHOWDOC_HOST},
			Destination: &runapi.Host,
		},
		&cli.BoolFlag{
			Name:        "debug",
			Value:       false,
			Usage:       "开启调试模式。",
			Destination: &log.IsDebug,
		},
	}
	app.Commands = []*cli.Command{
		{
			Name:  "flags",
			Usage: "查询应用全局相关参数。",
			Action: func(c *cli.Context) error {
				log.Info("showdoc.host=%s", runapi.Host)
				return nil
			},
		},
		{
			Name:  "item",
			Usage: "ShowDoc 项目导出和导入。",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:        "ssid",
					Value:       "",
					Usage:       "ShowDoc 用户的 SessionId，浏览器登录后从 Cookie 中获取 PHPSESSID。",
					Destination: &runapi.SSID,
				},
				&cli.StringFlag{
					Name:     flagItemId,
					Value:    "",
					Usage:    "ShowDoc 项目 id。",
					Required: true,
				},
			},
			Subcommands: []*cli.Command{
				{
					Name:  "info",
					Usage: "ShowDoc 项目详情",
					Action: func(c *cli.Context) error {
						ItemInfo(c.String(flagItemId))
						return nil
					},
				},
				{
					Name:  "export",
					Usage: "导出 ShowDoc 项目到目标文件夹。",
					Flags: []cli.Flag{
						&cli.StringFlag{
							Name:     flagDir,
							Value:    "",
							Usage:    "目标文件夹。",
							Required: true,
						},
					},
					Action: func(c *cli.Context) error {
						ItemExport(c.String(flagItemId), c.String(flagDir))
						return nil
					},
				},
				{
					Name:  "import",
					Usage: "导入 ShowDoc 文档。",
					Flags: []cli.Flag{
						&cli.StringFlag{
							Name:     flagFile,
							Value:    "",
							Usage:    "文件名",
							Required: true,
						},
					},
					Action: func(c *cli.Context) error {
						ItemImport(c.String(flagItemId), c.String(flagFile))
						return nil
					},
				},
			},
		},
		{
			Name:        "update",
			Aliases:     []string{"u"},
			Usage:       "解析 Go 源码中的注释，生成并更新 ShowDoc 文档。",
			Description: "为了确保能正确解析类型，请在目标项目路径下使用该工具。",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:        "apikey",
					Value:       "",
					Usage:       "ShowDoc 开放 API 认证凭证。",
					EnvVars:     []string{GOSHOWDOC_APIKEY},
					Destination: &runapi.ApiKey,
				},
				&cli.StringFlag{
					Name:        "apitoken",
					Value:       "",
					Usage:       "ShowDoc 开放 API 认证凭证。",
					EnvVars:     []string{GOSHOWDOC_APITOKEN},
					Destination: &runapi.ApiToken,
				},
				&cli.StringFlag{
					Name:     flagDir,
					Value:    "",
					Usage:    "搜索 Go 源码文件的目录，该目录下必须有 Go 源码文件。",
					Required: true,
				},
			},
			Action: func(c *cli.Context) error {
				Update(c.String(flagDir))
				return nil
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
