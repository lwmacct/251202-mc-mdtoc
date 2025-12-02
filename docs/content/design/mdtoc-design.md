# mc-mdtool 开发计划

**mc-mdtool** 是一个 Markdown CLI 工具集，`toc` 是其中的一个子命令。

TOC 功能基于 [frnmst/md-toc](https://github.com/frnmst/md-toc) 的设计思路重写。

## 0. 规范与参考项目

### 0.1 Markdown 规范版本

| 规范                                              | 版本     | 发布日期   | 说明        |
| ------------------------------------------------- | -------- | ---------- | ----------- |
| [CommonMark](https://spec.commonmark.org/0.31.2/) | 0.31.2   | 2024-01-28 | 最新标准    |
| [GFM](https://github.github.com/gfm/)             | 0.29-gfm | 2019-04-06 | GitHub 扩展 |

**CommonMark 0.31 新变化**：

- Unicode 符号现在被视为 Unicode 标点（影响强调判断）
- 内联 HTML 注释规则与 HTML 规范对齐
- 新增 `<search>` 元素到已知块元素列表
- 代码块闭合栅栏可以跟随制表符

### 0.2 参考项目 (vendor/)

| 项目                                                               | 语言   | Stars | 说明                          |
| ------------------------------------------------------------------ | ------ | ----- | ----------------------------- |
| [md-toc](https://github.com/frnmst/md-toc)                         | Python | ~200  | TOC 生成，多解析器支持        |
| [gh-md-toc-go](https://github.com/ekalinin/github-markdown-toc.go) | Go     | ~600  | Go TOC 生成器                 |
| [goldmark](https://github.com/yuin/goldmark)                       | Go     | ~3.5k | CommonMark 解析器 (Hugo 使用) |
| [goldmark-toc](https://github.com/abhinav/goldmark-toc)            | Go     | ~100  | goldmark TOC 扩展             |
| [glamour](https://github.com/charmbracelet/glamour)                | Go     | ~2k   | Markdown 渲染 (使用 goldmark) |
| [mdsf](https://github.com/hougesen/mdsf)                           | Rust   | ~300  | 代码块格式化 (372+ 格式器)    |

### 0.3 其他参考工具

| 工具                                                       | 语言    | 功能            |
| ---------------------------------------------------------- | ------- | --------------- |
| [markdownlint](https://github.com/DavidAnson/markdownlint) | Node.js | Markdown 检查   |
| [mdformat](https://github.com/hukkin/mdformat)             | Python  | Markdown 格式化 |
| [lychee](https://github.com/lycheeverse/lychee)            | Rust    | 链接检查        |
| [liche](https://github.com/raviqqe/liche)                  | Go      | 链接检查        |

### 0.4 实现策略选择

**方案 A: 基于 goldmark 扩展** ✅ 推荐

- 复用成熟的 CommonMark 解析器
- AST 结构清晰，易于扩展
- Hugo/VitePress 兼容性好

**方案 B: 独立实现解析器**

- 更轻量，无外部依赖
- 需要处理所有边界情况
- 适合简单场景

**决策**: 采用方案 A，基于 goldmark 实现，同时参考 gh-md-toc-go 的 anchor link 生成逻辑。

## 1. 项目目标

为 Markdown 文件自动生成符合规范的目录（Table of Contents），支持多种文档框架的语法规则。

### 1.1 核心功能

| 功能       | 说明                               | 优先级 |
| ---------- | ---------------------------------- | ------ |
| 标题解析   | 解析 ATX 风格标题 (`# ~ ######`)   | P0     |
| 锚点生成   | 生成符合 GitHub 规范的 anchor link | P0     |
| TOC 标记   | 支持 `<!--TOC-->` 标记定位         | P0     |
| 原地更新   | `-i` 直接修改文件                  | P0     |
| 差异检测   | `-d` 检查 TOC 是否需要更新         | P1     |
| 有序列表   | `-o` 生成 `1. 2. 3.` 格式          | P1     |
| 多框架支持 | VitePress、Hugo 等                 | P2     |

### 1.2 命令行接口

```shell
# 查看帮助
mc-mdtool --help
mc-mdtool toc --help

# 输出 TOC 到 stdout
mc-mdtool toc README.md

# 原地更新文件
mc-mdtool toc -i README.md

# 检查差异（CI 场景）
mc-mdtool toc -d README.md && echo "TOC is up to date"

# 指定标题层级
mc-mdtool toc -m 2 -M 4 README.md

# 有序列表
mc-mdtool toc -o README.md

# 使用别名
mc-mdtool t -i README.md
```

### 1.3 子命令规划

| 子命令    | 说明               | 优先级    | 参考项目             |
| --------- | ------------------ | --------- | -------------------- |
| `toc`     | 生成目录           | P0 (当前) | md-toc, goldmark-toc |
| `fmt`     | 格式化 Markdown    | P2        | mdformat, mdsf       |
| `lint`    | 检查 Markdown 规范 | P2        | markdownlint         |
| `links`   | 检查链接有效性     | P2        | lychee, liche        |
| `convert` | 格式转换           | P3        | pandoc               |

#### toc 子命令 (P0)

```shell
mc-mdtool toc [options] <file>
  -m, --min-level   最小标题层级 (默认 1)
  -M, --max-level   最大标题层级 (默认 3)
  -i, --in-place    原地更新文件
  -d, --diff        检查差异
  -o, --ordered     有序列表
```

#### fmt 子命令 (P2 预留)

```shell
mc-mdtool fmt [options] <file>
  -i, --in-place    原地更新文件
  -w, --wrap        行宽 (默认 80, 0=不换行)
  --code            格式化代码块 (需要外部格式化器)
```

#### lint 子命令 (P2 预留)

```shell
mc-mdtool lint [options] <file>
  -c, --config      配置文件路径
  --fix             自动修复
  --format          输出格式 (text/json)
```

#### links 子命令 (P2 预留)

```shell
mc-mdtool links [options] <file>
  --external        只检查外部链接
  --internal        只检查内部链接
  --timeout         请求超时 (默认 10s)
```

## 2. 架构设计

### 2.1 目录结构

```
cmd/
└── mc-mdtool/
    └── main.go              # 入口点

internal/
├── commands/
│   └── mdtool/
│       ├── command.go       # 主命令 (子命令注册)
│       └── toc/
│           ├── command.go   # toc 子命令定义
│           └── action.go    # toc 命令处理逻辑
│       └── fmt/             # (预留) 格式化子命令
│       └── lint/            # (预留) 检查子命令
│
└── mdtoc/                   # TOC 核心库
    ├── types.go             # 核心类型 (Header, TOCLine, Options)
    ├── parser.go            # Parser 接口 + goldmark 封装
    ├── anchor.go            # Anchor link 生成 (GitHub 风格)
    ├── generator.go         # TOC 字符串生成器
    ├── marker.go            # <!--TOC--> 标记处理
    └── parsers/
        ├── github.go        # GitHub 解析器 (默认)
        ├── vitepress.go     # VitePress 解析器 (预留)
        └── hugo.go          # Hugo 解析器 (预留)
```

### 2.2 依赖关系

```
github.com/yuin/goldmark           # CommonMark 解析器
github.com/yuin/goldmark/parser    # AST 解析
github.com/yuin/goldmark/ast       # AST 节点类型
```

### 2.3 核心接口

```go
// Parser 解析器接口
type Parser interface {
    Name() string
    ParseHeaders(content []byte, opts ParserOptions) ([]*Header, error)
    BuildAnchorLink(text string, duplicateCounter map[string]int) string
    SkipLine(line string, state *ParseState) bool
}

// Generator 生成器接口
type Generator interface {
    Generate(headers []*Header, opts ParserOptions) string
    FindMarkers(content []byte, marker string) (*TOCMarker, error)
    InsertTOC(content []byte, toc string, marker string) ([]byte, error)
}
```

### 2.4 解析器差异对比

| 特性        | GitHub      | VitePress           | Hugo                |
| ----------- | ----------- | ------------------- | ------------------- |
| Frontmatter | 不跳过      | 跳过 YAML           | 跳过 YAML/TOML      |
| 代码块      | ``` ~~~     | ``` ~~~             | ``` ~~~             |
| 锚点规则    | 小写+连字符 | 小写+连字符         | 可配置              |
| 特殊容器    | 无          | `:::` 块            | Shortcodes          |
| 标题 ID     | 自动生成    | 支持 `{#custom-id}` | 支持 `{#custom-id}` |

## 3. GitHub 解析器规范

### 3.1 ATX 标题规则

参考 [CommonMark Spec 0.31.2](https://spec.commonmark.org/0.31.2/#atx-headings)：

1. 以 1-6 个 `#` 开头
2. `#` 后必须有空格或行尾
3. 前面最多 3 个空格缩进
4. 忽略代码块内的标题

```markdown
# Valid H1

## Valid H2 (1 space indent)

### Valid H3 (2 spaces indent)

#### Valid H4 (3 spaces indent)

    ##### Invalid (4 spaces = code block)

#Invalid (no space after #)
```

### 3.2 Anchor Link 生成规则

参考 [github/html-pipeline](https://github.com/gjtorikian/html-pipeline)：

```
输入: "Hello, World! 你好"
步骤:
1. 转小写        → "hello, world! 你好"
2. 移除 HTML     → "hello, world! 你好"
3. 移除强调符号  → "hello, world! 你好"
4. 保留 \w\- 空格 → "hello world 你好"
5. 空格转连字符  → "hello-world-你好"
输出: "hello-world-你好"
```

### 3.3 重复标题处理

```markdown
# Title → #title

# Title → #title-1

# Title → #title-2
```

### 3.4 代码块检测

`````markdown
````python # 开始代码块
# Not a header   # 忽略
```              # 结束代码块

~~~             # 开始代码块
# Not a header  # 忽略
~~~             # 结束代码块
````
`````

`````

## 4. TOC 标记规范

### 4.1 标记格式

使用 HTML 注释作为标记，渲染后不可见：

```markdown
<!--TOC-->

- [Section 1](#section-1)
- [Section 2](#section-2)

<!--TOC-->
```

### 4.2 更新逻辑

1. 查找第一个 `<!--TOC-->` 标记
2. 查找第二个 `<!--TOC-->` 标记（可选）
3. 替换两个标记之间的内容
4. 如果只有一个标记，在其后插入

## 5. 实现计划

### Phase 1: 核心功能 (P0)

- [ ] GitHub 解析器实现
  - [ ] ATX 标题解析
  - [ ] 代码块跳过
  - [ ] Anchor link 生成
  - [ ] 重复标题处理
- [ ] TOC 生成器
  - [ ] 无序列表格式
  - [ ] 缩进计算
- [ ] 文件操作
  - [ ] 标记查找
  - [ ] 内容替换
  - [ ] stdout 输出

### Phase 2: 增强功能 (P1)

- [ ] 有序列表支持
- [ ] 差异检测 (`--diff`)
- [ ] 多文件批量处理
- [ ] 自定义标记字符串

### Phase 3: 多框架支持 (P2)

- [ ] VitePress 解析器
  - [ ] Frontmatter 跳过
  - [ ] `:::` 容器跳过
  - [ ] 自定义标题 ID `{#id}`
- [ ] Hugo 解析器
- [ ] Docusaurus 解析器

## 6. 测试用例

### 6.1 标题解析测试

```markdown
# H1 Title

## H2 Title

### H3 Title

#### H4 Title

##### H5 Title

###### H6 Title

####### Not a header (7 #)
#No space after hash
```

### 6.2 特殊字符测试

```markdown
# Hello, World!

# C++ Programming

# 中文标题

# Title with `code`

# Title with **bold** and _italic_

# Title with [link](url)
```

### 6.3 代码块测试

````markdown
# Real Header

​```go

# Not a header

​```

# Another Real Header
`````

## 7. 参考资料

### 7.1 规范文档

- [CommonMark Spec 0.31.2](https://spec.commonmark.org/0.31.2/) - 最新 CommonMark 标准
- [GitHub Flavored Markdown Spec](https://github.github.com/gfm/) - GFM 扩展规范
- [HTML5 Spec - Named Character References](https://html.spec.whatwg.org/multipage/named-characters.html)

### 7.2 参考实现 (vendor/)

| 目录                   | 项目                            | 重点参考内容                    |
| ---------------------- | ------------------------------- | ------------------------------- |
| `vendor/md-toc/`       | frnmst/md-toc                   | `api.py` - anchor link 生成算法 |
| `vendor/gh-md-toc-go/` | ekalinin/github-markdown-toc.go | `ghtoc.go` - GitHub 风格锚点    |
| `vendor/goldmark/`     | yuin/goldmark                   | `parser/` - AST 解析结构        |
| `vendor/goldmark-toc/` | abhinav/goldmark-toc            | `toc.go` - TOC 生成模式         |

### 7.3 文档框架

- [VitePress Markdown Extensions](https://vitepress.dev/guide/markdown) - 自定义容器、frontmatter
- [Hugo Markup Configuration](https://gohugo.io/getting-started/configuration-markup/) - goldmark 配置
- [Docusaurus Markdown Features](https://docusaurus.io/docs/markdown-features)

### 7.4 相关工具

- [github/cmark-gfm](https://github.com/github/cmark-gfm) - GitHub 官方 CommonMark 实现
- [gjtorikian/html-pipeline](https://github.com/gjtorikian/html-pipeline) - GitHub anchor link 生成源码
