package mdtool

import (
	"github.com/lwmacct/251202-mc-mdtool/internal/commands/mdtool/toc"
	"github.com/urfave/cli/v3"
)

// Command 返回 mc-mdtool 主命令
func Command(version string) *cli.Command {
	return &cli.Command{
		Name:    "mc-mdtool",
		Usage:   "Markdown CLI 工具集",
		Version: version,
		Commands: []*cli.Command{
			toc.Command(),
			// 预留其他子命令
			// fmt.Command(),    // 格式化
			// lint.Command(),   // 检查
			// convert.Command(), // 转换
		},
	}
}
