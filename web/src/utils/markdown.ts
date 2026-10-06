/**
 * Markdown 文章导入解析器
 *
 * 支持两种来源：
 *   1. 带 YAML 前置声明（front matter）的文件：
 *        ---
 *        title: 我的文章
 *        category: 技术
 *        tags: [Go, Redis]
 *        abstract: 一句话简介
 *        cover: /uploads/image/xxx.png
 *        ---
 *        正文……
 *   2. 普通 .md 文件：标题取正文第一个一级标题（没有就用文件名），简介取第一段正文。
 *
 * 这里刻意不引入 yaml 依赖，而是手写一个「够用就好」的子集解析器，
 * 只支持以下写法：
 *        key: value
 *        key: "带引号的值"
 *        key: [a, b, c]
 *        key: a, b, c
 *        key:
 *          - a
 *          - b
 */

export interface ParsedArticle {
  title: string;
  content: string;
  category: string;
  tags: string[];
  abstract: string;
  cover: string;
}

/** 自动提取简介时的截断长度 */
const ABSTRACT_MAX_LENGTH = 150;

/** front matter 的字段别名：中英文都认，避免用户从别的工具导出后用不了 */
const FIELD_ALIASES: Record<string, keyof ParsedArticle> = {
  title: "title",
  "标题": "title",

  category: "category",
  categories: "category",
  "类别": "category",
  "分类": "category",

  tags: "tags",
  tag: "tags",
  "标签": "tags",

  abstract: "abstract",
  description: "abstract",
  summary: "abstract",
  excerpt: "abstract",
  "简介": "abstract",
  "摘要": "abstract",

  cover: "cover",
  image: "cover",
  thumbnail: "cover",
  "封面": "cover",
};

/** 去掉一层包裹的引号 */
function stripQuotes(value: string): string {
  const trimmed = value.trim();

  if (trimmed.length >= 2) {
    const first = trimmed[0];
    const last = trimmed[trimmed.length - 1];

    if ((first === '"' && last === '"') || (first === "'" && last === "'")) {
      return trimmed.slice(1, -1);
    }
  }

  return trimmed;
}

/** 解析 `[a, b]` / `a, b, c` / `单个` 这几种标签写法 */
function parseTags(value: string): string[] {
  const inner =
    value.startsWith("[") && value.endsWith("]")
      ? value.slice(1, -1)
      : value;

  return inner
    .split(",")
    .map((item) => stripQuotes(item))
    .filter(Boolean);
}

/** 拆出前置声明与正文 */
function splitFrontMatter(raw: string): { meta: string; body: string } {
  // 必须从文件最开头开始，避免把正文里的 --- 分隔线误判成前置声明
  const matched = raw.match(/^---\n([\s\S]*?)\n---[ \t]*(?:\n|$)/);

  if (!matched) {
    return { meta: "", body: raw };
  }

  return { meta: matched[1], body: raw.slice(matched[0].length) };
}

/** 解析前置声明块 */
function parseFrontMatter(meta: string): Partial<ParsedArticle> {
  const result: Record<string, unknown> = {};

  if (!meta) {
    return result as Partial<ParsedArticle>;
  }

  const lines = meta.split("\n");

  for (let i = 0; i < lines.length; i++) {
    const line = lines[i];

    // 跳过空行与注释
    if (!line.trim() || line.trim().startsWith("#")) {
      continue;
    }

    const colonIndex = line.indexOf(":");
    if (colonIndex === -1) {
      continue;
    }

    const rawKey = line.slice(0, colonIndex).trim();
    const key = FIELD_ALIASES[rawKey] ?? FIELD_ALIASES[rawKey.toLowerCase()];

    // 不认识的字段直接忽略
    if (!key) {
      continue;
    }

    const value = line.slice(colonIndex + 1).trim();

    // 块状列表：key: 后面跟着若干行 "- xxx"
    if (value === "") {
      const items: string[] = [];
      let j = i + 1;

      while (j < lines.length && /^\s*-\s+/.test(lines[j])) {
        items.push(stripQuotes(lines[j].replace(/^\s*-\s+/, "").trim()));
        j++;
      }

      if (items.length) {
        result[key] = items;
        i = j - 1;
      }

      continue;
    }

    result[key] = key === "tags" ? parseTags(value) : stripQuotes(value);
  }

  return result as Partial<ParsedArticle>;
}

/** 去掉行内 markdown 标记，让简介是纯文本 */
function stripInlineMarkdown(text: string): string {
  return text
    .replace(/!\[[^\]]*\]\([^)]*\)/g, "") // 图片整体去掉
    .replace(/\[([^\]]*)\]\([^)]*\)/g, "$1") // 链接只保留文字
    .replace(/[*_~`]+/g, "") // 强调、行内代码的标记符
    .replace(/\s+/g, " ")
    .trim();
}

/** 从正文里推断简介：第一段不是标题/图片/列表/代码块/表格/HTML 的文字 */
function extractAbstract(body: string): string {
  let inCodeBlock = false;

  for (const line of body.split("\n")) {
    const text = line.trim();

    // 代码围栏：切换状态；围栏本身与它内部的所有行都要跳过
    if (/^(```|~~~)/.test(text)) {
      inCodeBlock = !inCodeBlock;
      continue;
    }

    if (inCodeBlock) continue;

    if (!text) continue;
    if (/^#{1,6}\s/.test(text)) continue;
    if (/^!\[/.test(text)) continue;
    if (/^>/.test(text)) continue;
    if (/^([-*+]|\d+\.)\s/.test(text)) continue;
    if (/^\|/.test(text)) continue;
    if (/^</.test(text)) continue;

    const plain = stripInlineMarkdown(text);

    if (!plain) continue;

    return plain.length > ABSTRACT_MAX_LENGTH
      ? plain.slice(0, ABSTRACT_MAX_LENGTH) + "…"
      : plain;
  }

  return "";
}

/**
 * 解析一个 Markdown 文件的内容，得到可直接填进表单的文章信息
 *
 * @param raw      文件的原始文本
 * @param fileName 文件名，用于在没有其它线索时兜底成标题
 */
export function parseMarkdownArticle(raw: string, fileName = ""): ParsedArticle {
  // 统一换行符并去掉 BOM，否则 Windows 上编辑的文件会让所有正则失配
  const normalized = raw.replace(/^\uFEFF/, "").replace(/\r\n/g, "\n");

  const { meta, body } = splitFrontMatter(normalized);
  const frontMatter = parseFrontMatter(meta);
  const trimmedBody = body.trim();

  // 标题优先级：前置声明 > 正文首个一级标题 > 文件名
  let title = (frontMatter.title ?? "").trim();
  let content = trimmedBody;

  if (!title) {
    const heading = trimmedBody.match(/^#[ \t]+(.+)$/m);

    if (heading) {
      title = heading[1].trim();

      // 标题已经单独填进「文章标题」，正文里就别再重复一遍了
      content = trimmedBody.replace(/^#[ \t]+.+\n?/, "").trim();
    } else {
      title = fileName.replace(/\.(md|markdown)$/i, "").trim();
    }
  }

  return {
    title,
    content,
    category: frontMatter.category ?? "",
    tags: Array.isArray(frontMatter.tags) ? frontMatter.tags : [],
    abstract: (frontMatter.abstract ?? "").trim() || extractAbstract(content),
    cover: frontMatter.cover ?? "",
  };
}
