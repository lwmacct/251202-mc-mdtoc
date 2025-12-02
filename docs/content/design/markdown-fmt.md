# Markdown 格式化调研

参考 [Prettier](https://prettier.io/) 的设计理念，调研 Markdown 格式化的规则和选项。

## 1. Prettier 格式化选项

### 1.1 核心选项

| 选项 | 默认值 | 可选值 | 说明 |
|------|--------|--------|------|
| `proseWrap` | `"preserve"` | `always`, `never`, `preserve` | 段落换行策略 |
| `printWidth` | `80` | 数字 | 目标行宽 |
| `tabWidth` | `2` | 数字 | 缩进空格数 |
| `useTabs` | `false` | 布尔 | 使用 Tab 缩进 |
| `endOfLine` | `"lf"` | `lf`, `crlf`, `cr`, `auto` | 行尾符 |
| `embeddedLanguageFormatting` | `"auto"` | `auto`, `off` | 代码块格式化 |

### 1.2 proseWrap 详解

```
proseWrap: "preserve" (默认)
  保持原有换行，不做改动
  适合: GitHub comments, BitBucket 等行敏感渲染器

proseWrap: "always"
  超过 printWidth 自动换行
  适合: 需要固定行宽的场景

proseWrap: "never"
  移除段落内换行，每段一行
  适合: 依赖编辑器软换行的场景
```

### 1.3 配置示例

```json
// .prettierrc
{
  "printWidth": 80,
  "tabWidth": 2,
  "proseWrap": "preserve",
  "overrides": [
    {
      "files": ["*.md"],
      "options": {
        "proseWrap": "always"
      }
    }
  ]
}
```

## 2. Markdown 格式化规则

### 2.1 标题格式化

```markdown
# Before
#Title without space
##  Extra spaces

# After
# Title without space
## Extra spaces
```

规则:
- `#` 后必须有且只有一个空格
- 移除标题前后多余空行

### 2.2 列表格式化

```markdown
# Before
-item1
-  item2
  - nested
1.first
2.second

# After
- item1
- item2
  - nested
1. first
2. second
```

规则:
- 列表标记后一个空格
- 嵌套列表缩进 2 空格
- 有序列表数字后加点和空格

### 2.3 代码块格式化

```markdown
# Before
```js
const x=1
```

# After
```js
const x = 1;
```
```

规则:
- 代码块内容可调用外部格式化器
- 保持语言标识符

### 2.4 表格格式化

```markdown
# Before
|a|b|c|
|-|-|-|
|1|2|3|

# After
| a | b | c |
|---|---|---|
| 1 | 2 | 3 |
```

规则:
- 单元格内容两侧加空格
- 分隔符对齐
- 可选: 列宽对齐

### 2.5 链接格式化

```markdown
# Before
[text](  url  )
[text]( url "title" )

# After
[text](url)
[text](url "title")
```

规则:
- 移除 URL 前后空格
- 保留 title 属性

### 2.6 空行规范

```markdown
# Before
# Title


paragraph1


paragraph2

# After
# Title

paragraph1

paragraph2
```

规则:
- 标题后一个空行
- 段落间一个空行
- 文件末尾一个换行

## 3. mc-mdtool fmt 设计

### 3.1 命令行接口

```shell
mc-mdtool fmt [options] <file>

Options:
  -i, --in-place       原地更新文件
  -w, --print-width    目标行宽 (默认 80, 0=不限制)
  --prose-wrap         段落换行: always|never|preserve (默认 preserve)
  --tab-width          缩进空格数 (默认 2)
  --use-tabs           使用 Tab 缩进
  --end-of-line        行尾符: lf|crlf|auto (默认 lf)
  --code               格式化代码块 (需要外部格式化器)
  -c, --config         配置文件路径
```

### 3.2 配置文件

```yaml
# .mdtool.yaml
fmt:
  print-width: 80
  prose-wrap: preserve
  tab-width: 2
  use-tabs: false
  end-of-line: lf

  # 代码块格式化器 (可选)
  code-formatters:
    go: gofmt
    js: prettier --parser babel
    python: black -
```

### 3.3 格式化优先级

```
1. 命令行参数
2. 项目配置文件 (.mdtool.yaml)
3. 用户配置文件 (~/.config/mdtool/config.yaml)
4. 默认值
```

## 4. 实现策略

### 4.1 基于 goldmark AST

```go
// 伪代码
func Format(content []byte, opts Options) ([]byte, error) {
    // 1. 解析为 AST
    doc := goldmark.Parse(content)

    // 2. 遍历 AST 节点
    ast.Walk(doc, func(n ast.Node, entering bool) {
        switch node := n.(type) {
        case *ast.Heading:
            formatHeading(node, opts)
        case *ast.List:
            formatList(node, opts)
        case *ast.FencedCodeBlock:
            formatCodeBlock(node, opts)
        case *ast.Table:
            formatTable(node, opts)
        // ...
        }
    })

    // 3. 渲染回 Markdown
    return Render(doc, opts)
}
```

### 4.2 难点

| 难点 | 说明 | 解决方案 |
|------|------|----------|
| 保留注释 | goldmark 不保留 HTML 注释位置 | 预处理标记 |
| 代码块格式化 | 需调用外部工具 | 类似 mdsf 设计 |
| 表格对齐 | 需计算字符宽度 (CJK) | runewidth 库 |
| 原文保留 | 某些场景需保持原样 | preserve 模式 |

## 5. 参考项目

| 项目 | 语言 | 特点 |
|------|------|------|
| [Prettier](https://prettier.io/) | JS | 业界标准，opinionated |
| [mdformat](https://github.com/hukkin/mdformat) | Python | CommonMark 兼容，可扩展 |
| [mdsf](https://github.com/hougesen/mdsf) | Rust | 专注代码块格式化 |
| [remark-stringify](https://github.com/remarkjs/remark) | JS | AST to Markdown |

## 6. 参考资料

- [Prettier Options](https://prettier.io/docs/options)
- [Prettier Markdown Support (v1.8)](https://prettier.io/blog/2017/11/07/1.8.0.html)
- [mdformat Documentation](https://mdformat.readthedocs.io/)
- [mdsf GitHub](https://github.com/hougesen/mdsf)
