package toc

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"github.com/urfave/cli/v3"
)

func action(ctx context.Context, cmd *cli.Command) error {
	// 解析命令行参数
	minLevel := cmd.Int("min-level")
	maxLevel := cmd.Int("max-level")
	inPlace := cmd.Bool("in-place")
	diff := cmd.Bool("diff")
	ordered := cmd.Bool("ordered")

	// 获取文件参数
	file := cmd.Args().First()
	if file == "" {
		return fmt.Errorf("请指定要处理的 Markdown 文件")
	}

	// 检查文件是否存在
	if _, err := os.Stat(file); os.IsNotExist(err) {
		return fmt.Errorf("文件不存在: %s", file)
	}

	// 验证层级参数
	if minLevel < 1 || minLevel > 6 {
		return fmt.Errorf("min-level 必须在 1-6 之间")
	}
	if maxLevel < 1 || maxLevel > 6 {
		return fmt.Errorf("max-level 必须在 1-6 之间")
	}
	if minLevel > maxLevel {
		return fmt.Errorf("min-level 不能大于 max-level")
	}

	slog.Debug("处理 Markdown 文件",
		"file", file,
		"min_level", minLevel,
		"max_level", maxLevel,
		"in_place", inPlace,
		"diff", diff,
		"ordered", ordered,
	)

	// TODO: 实现 Markdown 目录生成逻辑
	// 1. 读取文件内容
	// 2. 解析标题结构 (GitHub 风格 ATX 标题)
	// 3. 生成 TOC (anchor link 规则参考 vendor/md-toc/md_toc/api.py)
	// 4. 查找 <!--TOC--> 标记并替换
	// 5. 根据 inPlace/diff 参数决定输出方式

	fmt.Println("md-toc 功能开发中...")

	return nil
}
