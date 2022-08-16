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

const flagDir = "dir"

func main() {
	cli.HelpFlag = &cli.BoolFlag{
		Name:  "help",
		Usage: "显示帮助",
	}
	app := cli.NewApp()
	app.Name = "goshowdoc"
	app.Usage = "ShowDoc API 接口文档工具"
	app.Description = `通过代码注释生成 API 接口文档`
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
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
