package main

import (
	"fmt"
	"os"

	"github.com/whaios/goshowdoc/datadict"

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

const (
	flagDriver = "driver"
	flagHost   = "host"
	flagUser   = "user"
	flagPwd    = "pwd"
	flagDb     = "db"
	flagSchema = "schema"
	flagCat    = "cat"
)

func main() {
	cli.HelpFlag = &cli.BoolFlag{
		Name:  "help",
		Usage: "显示帮助",
	}
	app := cli.NewApp()
	app.Name = "goshowdoc"
	app.Usage = "ShowDoc API 接口文档工具"
	app.Description = `项目地址： https://github.com/whaios/goshowdoc
支持以下功能：
1. 通过代码注释生成 API 接口文档。
2. 自动化生成数据字典。
3. 导出和导入 ShowDoc 项目。`
	app.Version = Version

	app.Flags = []cli.Flag{
		&cli.StringFlag{
			Name:        "host",
			Usage:       "ShowDoc 地址。",
			Value:       runapi.Host,
			Destination: &runapi.Host,
			EnvVars:     []string{GOSHOWDOC_HOST},
		},
		&cli.StringFlag{
			Name:        "apikey",
			Usage:       "ShowDoc 开放 API 认证凭证。",
			Value:       runapi.ApiKey,
			Destination: &runapi.ApiKey,
			EnvVars:     []string{GOSHOWDOC_APIKEY},
		},
		&cli.StringFlag{
			Name:        "apitoken",
			Usage:       "ShowDoc 开放 API 认证凭证。",
			Value:       runapi.ApiToken,
			Destination: &runapi.ApiToken,
			EnvVars:     []string{GOSHOWDOC_APITOKEN},
		},
		&cli.BoolFlag{
			Name:        "debug",
			Usage:       "开启调试模式。",
			Value:       log.IsDebug,
			Destination: &log.IsDebug,
		},
	}
	app.Commands = []*cli.Command{
		{
			Name:  "flags",
			Usage: "查询应用全局相关参数。",
			Action: func(c *cli.Context) error {
				log.Info("showdoc.host=%s", runapi.Host)
				log.Info("showdoc.apiKey=%s", runapi.ApiKey)
				log.Info("showdoc.apiToken=%s", runapi.ApiToken)
				return nil
			},
		},
		{
			Name:        "update",
			Aliases:     []string{"u"},
			Usage:       "解析 Go 源码中的注释，生成并更新 ShowDoc 文档。",
			Description: "为了确保能正确解析类型，请在目标项目路径下使用该工具。",
			Flags: []cli.Flag{
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
		{
			Name:    "datadict",
			Aliases: []string{"dd"},
			Usage:   "自动生成数据字典。",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:  flagDriver,
					Usage: fmt.Sprintf("数据库类型，支持：%s, %s, %s, %s", datadict.MySQL, datadict.PostgreSQL, datadict.SQLServer, datadict.SQlite),
					Value: datadict.MySQL,
				},
				&cli.StringFlag{
					Name:    flagHost,
					Aliases: []string{"h"},
					Usage:   "数据库地址和端口，如果是SQlite数据库则为文件",
					Value:   "127.0.0.1:3306",
				},
				&cli.StringFlag{
					Name:    flagUser,
					Aliases: []string{"u"},
					Usage:   "数据库用户名",
				},
				&cli.StringFlag{
					Name:    flagPwd,
					Aliases: []string{"p"},
					Usage:   "数据库密码",
				},
				&cli.StringFlag{
					Name:  flagDb,
					Usage: "要同步的数据库名",
				},
				&cli.StringFlag{
					Name:  flagSchema,
					Usage: "PostgreSQL 数据库模式",
				},
				&cli.StringFlag{
					Name:  flagCat,
					Usage: "文档所在目录，如果需要多层目录请用斜杠隔开，例如：“一层/二层/三层”",
				},
			},
			Action: func(c *cli.Context) error {
				UpdateDataDict(
					c.String(flagDriver),
					c.String(flagHost),
					c.String(flagUser),
					c.String(flagPwd),
					c.String(flagDb),
					c.String(flagSchema),
					c.String(flagCat),
				)
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
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
