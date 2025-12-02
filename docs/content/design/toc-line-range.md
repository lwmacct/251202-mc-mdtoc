# TOC 行号功能开发计划

## 1. 功能需求

在 TOC 输出中显示每个标题的行号，帮助用户快速定位内容。

### 1.1 输出示例

使用 **VS Code 兼容格式** `:line`，终端中可直接点击跳转：

```markdown
- [第一章 介绍](#第一章-介绍) `:10`
  - [1.1 背景](#11-背景) `:12`
  - [1.2 目标](#12-目标) `:19`
- [第二章 设计](#第二章-设计) `:26`
```

### 1.2 命令行接口

```shell
# 显示行号
mc-mdtool toc -L README.md
mc-mdtool toc --line-number README.md
```

### 1.3 格式说明

| 格式 | 说明 |
|------|------|
| `:line` | VS Code 终端可点击跳转 |
| 只显示起始行 | 简洁，结束行可从下一标题推断 |

## 2. 技术方案

### 2.1 核心修改

```
internal/mdtoc/
├── types.go      # Header 增加 Line 字段，Options 增加 LineNumber
├── parser.go     # 解析时记录行号
└── generator.go  # 支持行号输出
```

### 2.2 数据结构变更

```go
// types.go
type Header struct {
    Level      int
    Text       string
    AnchorLink string
    Line       int  // 新增：标题所在行 (1-based)
}

type Options struct {
    MinLevel   int
    MaxLevel   int
    Ordered    bool
    LineNumber bool // 新增：是否显示行号
}
```

### 2.3 实现步骤

1. **修改 Parser** - 从 AST 节点获取行号
2. **修改 Generator** - 添加 `:line` 格式输出
3. **修改 CLI** - 添加 `-L/--line-number` flag

## 3. 实现细节

### 3.1 行号获取

goldmark AST 通过 `text.Segment` 提供位置信息，需要将 byte offset 转换为行号：

```go
func byteOffsetToLine(src []byte, offset int) int {
    line := 1
    for i := 0; i < offset && i < len(src); i++ {
        if src[i] == '\n' {
            line++
        }
    }
    return line
}
```

### 3.2 输出格式

```go
// generator.go - VS Code 兼容格式
func (g *Generator) formatLink(h *Header) string {
    link := fmt.Sprintf("[%s](#%s)", h.Text, h.AnchorLink)
    if g.options.LineNumber && h.Line > 0 {
        link += fmt.Sprintf(" `:%d`", h.Line)
    }
    return link
}
```

## 4. 测试用例

输入：
```markdown
# Title

## Section 1
Content...

## Section 2
More content...
```

输出 (`mc-mdtool toc -L`):
```markdown
- [Title](#title) `:1`
  - [Section 1](#section-1) `:3`
  - [Section 2](#section-2) `:6`
```

## 5. 任务清单

- [ ] 修改 `types.go` - Header 增加 Line，Options 增加 LineNumber
- [ ] 修改 `parser.go` - 解析行号
- [ ] 修改 `generator.go` - 输出 `:line` 格式
- [ ] 修改 `toc/command.go` - 添加 `-L` flag
- [ ] 修改 `toc/action.go` - 传递选项
- [ ] 添加单元测试
