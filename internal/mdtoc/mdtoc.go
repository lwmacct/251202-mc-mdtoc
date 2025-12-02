package mdtoc

import (
	"bytes"
	"os"
	"strings"
)

// TOC 是主要的门面结构，封装所有 TOC 生成功能
type TOC struct {
	parser    *Parser
	generator *Generator
	marker    *MarkerHandler
	options   Options
}

// New 创建新的 TOC 实例
func New(opts Options) *TOC {
	return &TOC{
		parser:    NewParser(opts),
		generator: NewGenerator(opts),
		marker:    NewMarkerHandler(DefaultMarker),
		options:   opts,
	}
}

// GenerateFromFile 从文件生成 TOC 字符串
func (t *TOC) GenerateFromFile(filename string) (string, error) {
	content, err := os.ReadFile(filename)
	if err != nil {
		return "", err
	}
	return t.GenerateFromContent(content)
}

// GenerateFromContent 从内容生成 TOC 字符串
func (t *TOC) GenerateFromContent(content []byte) (string, error) {
	headers, err := t.parser.Parse(content)
	if err != nil {
		return "", err
	}
	return t.generator.Generate(headers), nil
}

// GenerateSectionTOCs 生成章节模式的 TOC (每个 H1 有独立的子目录)
func (t *TOC) GenerateSectionTOCs(content []byte) ([]SectionTOC, error) {
	// 解析所有标题
	headers, err := t.parser.ParseAllHeaders(content)
	if err != nil {
		return nil, err
	}

	// 按 H1 分割成章节
	sections := SplitSections(headers)

	// 为每个章节生成 TOC
	var sectionTOCs []SectionTOC
	for _, section := range sections {
		toc := t.generator.GenerateSection(section)
		if toc != "" {
			sectionTOCs = append(sectionTOCs, SectionTOC{
				H1Line: section.Title.Line - 1, // 转换为 0-based
				TOC:    toc,
			})
		}
	}

	return sectionTOCs, nil
}

// GenerateSectionTOCsWithOffset 生成章节模式的 TOC，预计算偏移量使行号一次正确
// 这个方法解决了需要执行两次 toc 命令才能得到正确行号的问题
func (t *TOC) GenerateSectionTOCsWithOffset(cleanContent []byte) ([]SectionTOC, error) {
	// 在干净内容上解析所有标题（基准行号）
	headers, err := t.parser.ParseAllHeaders(cleanContent)
	if err != nil {
		return nil, err
	}

	// 按 H1 分割成章节
	sections := SplitSections(headers)

	// 第一遍：计算每个章节的 TOC 内容行数
	type sectionInfo struct {
		section   *Section
		tocLines  int    // TOC 内容行数
		tocString string // TOC 字符串（不带行号的临时版本）
	}
	var infos []sectionInfo

	// 临时禁用行号生成，只计算 TOC 结构
	origLineNumber := t.options.LineNumber
	t.options.LineNumber = false
	t.generator = NewGenerator(t.options)

	for _, section := range sections {
		toc := t.generator.GenerateSection(section)
		if toc != "" {
			tocLines := strings.Count(toc, "\n") + 1
			infos = append(infos, sectionInfo{
				section:   section,
				tocLines:  tocLines,
				tocString: toc,
			})
		}
	}

	// 恢复行号设置
	t.options.LineNumber = origLineNumber
	t.generator = NewGenerator(t.options)

	// 第二遍：计算累积偏移量并应用到标题行号
	var sectionTOCs []SectionTOC
	cumulativeOffset := 0

	for _, info := range infos {
		// 计算这个 TOC 块会增加的行数
		tocBlockLines := CalcTOCBlockLines(info.tocString)

		// 保存原始 H1 行号（在干净内容中的位置，用于插入定位）
		originalH1Line := info.section.Title.Line

		// 调整子标题的行号（这些行号显示在 TOC 中）
		// 需要加上：1) 之前所有 TOC 块的累积偏移 2) 当前 TOC 块的行数
		for _, h := range info.section.SubHeaders {
			h.Line += cumulativeOffset + tocBlockLines
			h.EndLine += cumulativeOffset + tocBlockLines
		}

		// 生成带正确行号的 TOC
		toc := t.generator.GenerateSection(info.section)
		if toc != "" {
			sectionTOCs = append(sectionTOCs, SectionTOC{
				H1Line: originalH1Line - 1, // 使用干净内容中的行号（0-based），用于定位插入位置
				TOC:    toc,
			})
		}

		// 累加偏移量
		cumulativeOffset += tocBlockLines
	}

	return sectionTOCs, nil
}

// GenerateSectionTOCsPreview 生成章节模式的 TOC 预览 (用于 stdout 输出)
func (t *TOC) GenerateSectionTOCsPreview(content []byte) (string, error) {
	// 解析所有标题
	headers, err := t.parser.ParseAllHeaders(content)
	if err != nil {
		return "", err
	}

	// 按 H1 分割成章节
	sections := SplitSections(headers)

	var sb strings.Builder
	for i, section := range sections {
		toc := t.generator.GenerateSection(section)
		if toc != "" {
			sb.WriteString("### ")
			sb.WriteString(section.Title.Text)
			sb.WriteString("\n\n")
			sb.WriteString(toc)
			if i < len(sections)-1 {
				sb.WriteString("\n\n")
			}
		}
	}

	return sb.String(), nil
}

// UpdateFile 原地更新文件中的 TOC
// 如果文件没有 TOC 标记，会自动在第一个标题后插入
func (t *TOC) UpdateFile(filename string) error {
	content, err := os.ReadFile(filename)
	if err != nil {
		return err
	}

	var newContent []byte

	if t.options.SectionTOC {
		// 章节模式：在每个 H1 后插入独立的子目录
		// 先清理现有 TOC 块，获取干净内容
		cleanContent, _ := t.marker.CleanTOCBlocks(content)

		// 使用预计算偏移量的方法生成 TOC
		sectionTOCs, err := t.GenerateSectionTOCsWithOffset(cleanContent)
		if err != nil {
			return err
		}

		// 在干净内容上插入新的 TOC
		newContent = t.marker.InsertSectionTOCs(cleanContent, sectionTOCs)
	} else {
		// 普通模式：在 <!--TOC--> 标记处插入完整 TOC
		toc, err := t.GenerateFromContent(content)
		if err != nil {
			return err
		}

		markers := t.marker.FindMarkers(content)
		if markers.Found {
			newContent = t.marker.InsertTOC(content, toc)
		} else {
			newContent = t.marker.InsertTOCAfterFirstHeading(content, toc)
		}
	}

	return os.WriteFile(filename, newContent, 0644)
}

// CheckDiff 检查 TOC 是否需要更新
// 返回 true 表示需要更新 (有差异)
func (t *TOC) CheckDiff(filename string) (bool, error) {
	content, err := os.ReadFile(filename)
	if err != nil {
		return false, err
	}

	if t.options.SectionTOC {
		// 章节模式：生成新内容并与原内容比较
		cleanContent, _ := t.marker.CleanTOCBlocks(content)
		sectionTOCs, err := t.GenerateSectionTOCsWithOffset(cleanContent)
		if err != nil {
			return false, err
		}
		newContent := t.marker.InsertSectionTOCs(cleanContent, sectionTOCs)
		return !bytes.Equal(content, newContent), nil
	}

	// 普通模式：比较 TOC 内容
	newTOC, err := t.GenerateFromContent(content)
	if err != nil {
		return false, err
	}
	existingTOC := t.marker.ExtractExistingTOC(content)
	return newTOC != existingTOC, nil
}

// HasMarker 检查文件是否包含 TOC 标记
func (t *TOC) HasMarker(filename string) (bool, error) {
	content, err := os.ReadFile(filename)
	if err != nil {
		return false, err
	}
	markers := t.marker.FindMarkers(content)
	return markers.Found, nil
}
