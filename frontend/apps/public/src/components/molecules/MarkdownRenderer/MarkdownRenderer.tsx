import { useMemo } from "react";

const ALLOWED_TAGS = new Set([
  "p",
  "h1",
  "h2",
  "h3",
  "h4",
  "h5",
  "h6",
  "ul",
  "ol",
  "li",
  "strong",
  "em",
  "a",
  "code",
  "pre",
  "blockquote",
  "img",
  "figure",
  "figcaption"
]);

const ALLOWED_URI_PROTOCOLS = new Set(["http:", "https:", "mailto:"]);

function escapeHtml(value: string): string {
  return value
    .replace(/&/g, "&amp;")
    .replace(/</g, "&lt;")
    .replace(/>/g, "&gt;")
    .replace(/"/g, "&quot;")
    .replace(/'/g, "&#39;");
}

function sanitizeUrl(url: string): string | null {
  try {
    const parsed = new URL(url, "https://placeholder.invalid");
    if (!ALLOWED_URI_PROTOCOLS.has(parsed.protocol)) {
      return null;
    }

    if (parsed.protocol === "mailto:") {
      return parsed.href.replace("https://placeholder.invalid", "");
    }

    return parsed.href;
  } catch {
    return null;
  }

  return null;
}

function convertMarkdownToHtml(markdown: string): string {
  const lines = markdown.split(/\r?\n/);
  let html = "";
  let inList = false;

  const flushList = () => {
    if (inList) {
      html += "</ul>";
      inList = false;
    }
  };

  for (const rawLine of lines) {
    const line = rawLine.trimEnd();

    if (!line.trim()) {
      flushList();
      continue;
    }

    const headingMatch = line.match(/^(#{1,6})\s+(.*)$/);
    if (headingMatch) {
      flushList();
      const level = headingMatch[1].length;
      const content = renderInlineMarkdown(headingMatch[2]);
      html += `<h${level}>${content}</h${level}>`;
      continue;
    }

    const listMatch = line.match(/^[-*+]\s+(.*)$/);
    if (listMatch) {
      if (!inList) {
        html += "<ul>";
        inList = true;
      }
      html += `<li>${renderInlineMarkdown(listMatch[1])}</li>`;
      continue;
    }

    const blockquoteMatch = line.match(/^>\s+(.*)$/);
    if (blockquoteMatch) {
      flushList();
      html += `<blockquote>${renderInlineMarkdown(blockquoteMatch[1])}</blockquote>`;
      continue;
    }

    flushList();
    html += `<p>${renderInlineMarkdown(line)}</p>`;
  }

  flushList();
  return html;
}

function renderInlineMarkdown(value: string): string {
  const escaped = escapeHtml(value);

  const imagePattern = /!\[([^[\]]*)\]\(([^)]+)\)/g;
  const linkPattern = /\[([^[\]]+)\]\(([^)]+)\)/g;

  const withCode = escaped.replace(/`([^`]+)`/g, "<code>$1</code>");
  const withBold = withCode.replace(/\*\*([^*]+)\*\*/g, "<strong>$1</strong>");
  const withItalic = withBold.replace(/\*([^*]+)\*/g, "<em>$1</em>");

  const withImages = withItalic.replace(imagePattern, (_match, alt, url) => {
    const sanitizedUrl = sanitizeUrl(url);
    if (!sanitizedUrl) {
      return alt;
    }

    return `<img src="${sanitizedUrl}" alt="${escapeHtml(alt)}" loading="lazy" />`;
  });

  return withImages.replace(linkPattern, (_match, label, url) => {
    const sanitizedUrl = sanitizeUrl(url);
    if (!sanitizedUrl) {
      return label;
    }

    return `<a href="${sanitizedUrl}" target="_blank" rel="noreferrer">${label}</a>`;
  });
}

function sanitizeHtml(html: string): string {
  const parser = new DOMParser();
  const document = parser.parseFromString(html, "text/html");
  const walker = document.createTreeWalker(document.body, NodeFilter.SHOW_ELEMENT);

  const nodesToRemove: Element[] = [];

  while (walker.nextNode()) {
    const element = walker.currentNode as Element;
    if (!ALLOWED_TAGS.has(element.tagName.toLowerCase())) {
      nodesToRemove.push(element);
      continue;
    }

    [...element.attributes].forEach((attribute) => {
      const name = attribute.name.toLowerCase();

      if (name === "href" || name === "src") {
        const sanitized = sanitizeUrl(attribute.value);
        if (!sanitized) {
          element.removeAttribute(attribute.name);
        } else {
          element.setAttribute(attribute.name, sanitized);
        }
        return;
      }

      if (name === "target") {
        element.setAttribute("target", "_blank");
        element.setAttribute("rel", "noreferrer");
        return;
      }

      if (name.startsWith("on")) {
        element.removeAttribute(attribute.name);
        return;
      }

      if (!["alt", "title", "loading"].includes(name)) {
        element.removeAttribute(attribute.name);
      }
    });
  }

  for (const node of nodesToRemove) {
    const textNode = document.createTextNode(node.textContent ?? "");
    node.replaceWith(textNode);
  }

  return document.body.innerHTML;
}

type MarkdownRendererProps = {
  markdown: string;
  html?: string;
} & React.HTMLAttributes<HTMLDivElement>;

export function MarkdownRenderer({ markdown, html, ...divProps }: MarkdownRendererProps) {
  const renderedHtml = useMemo(() => {
    if (html) {
      return sanitizeHtml(html);
    }

    return sanitizeHtml(convertMarkdownToHtml(markdown));
  }, [html, markdown]);

  return <div {...divProps} dangerouslySetInnerHTML={{ __html: renderedHtml }} />;
}
