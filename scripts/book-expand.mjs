#!/usr/bin/env node
/**
 * Builds PAGE.md per chapter: README + stub .go + tests (for Honkit / GitHub Pages).
 * PAGE.md is gitignored; run before book:build / book:serve.
 */
import fs from "fs";
import path from "path";

const root = path.resolve(import.meta.dirname, "..");

function listChapterDirs() {
  const out = [];
  for (const name of fs.readdirSync(root, { withFileTypes: true })) {
    if (!name.isDirectory()) continue;
    if (/^\d{2}-/.test(name.name)) out.push(name.name);
  }
  out.sort();
  for (const name of fs.readdirSync(path.join(root, "capstones"), {
    withFileTypes: true,
  })) {
    if (!name.isDirectory()) continue;
    if (/^\d{2}-/.test(name.name)) {
      out.push(path.join("capstones", name.name));
    }
  }
  return out;
}

function goSources(dir) {
  return fs
    .readdirSync(dir)
    .filter(
      (f) =>
        f.endsWith(".go") &&
        f !== "solution.go" &&
        !f.endsWith("_test.go") &&
        !f.startsWith("solution_")
    )
    .sort();
}

function goTests(dir) {
  return fs
    .readdirSync(dir)
    .filter((f) => f.endsWith("_test.go"))
    .sort();
}

function goCodeBlock(code) {
  // Honkit pages are Nunjucks; Go composite literals like {{-1, 0}} break the parser.
  return `{% raw %}\n\`\`\`go\n${code.trimEnd()}\n\`\`\`\n{% endraw %}\n`;
}

function section(title, filename, code) {
  return `### \`${filename}\`\n\n${goCodeBlock(code)}`;
}

/** Honkit passes <details> through to HTML — spoiler for solution.go */
function solutionSpoiler(filename, code) {
  return `<details>
<summary><strong>Эталон (\`${filename}\`)</strong> — нажмите, чтобы показать</summary>

${goCodeBlock(code)}

</details>
`;
}

function expandChapter(relDir) {
  const abs = path.join(root, relDir);
  const readmePath = path.join(abs, "README.md");
  if (!fs.existsSync(readmePath)) {
    console.warn(`skip ${relDir}: no README.md`);
    return;
  }

  const readme = fs.readFileSync(readmePath, "utf8");
  const pkg = path.basename(relDir).replace(/^\d{2}-/, "") || relDir;
  const module = "assignments-go";

  let body = readme.trimEnd();
  body += "\n\n---\n\n";
  body += "## Исходники на этой странице\n\n";
  body +=
    "Ниже — те же файлы, что в репозитории. Локально: клонируйте курс и правьте код в этой папке.\n\n";
  body += "```bash\n";
  body += `git clone https://github.com/moonmouse11/assignments-go.git\ncd assignments-go\n`;
  body += `go test ./${relDir}/\n`;
  body += "```\n\n";

  const stubs = goSources(abs);
  const tests = goTests(abs);

  if (stubs.length) {
    body += "## Заготовки (реализуйте)\n\n";
    for (const f of stubs) {
      const code = fs.readFileSync(path.join(abs, f), "utf8");
      body += section(f, f, code) + "\n";
    }
  }

  if (tests.length) {
    body += "## Тесты\n\n";
    for (const f of tests) {
      const code = fs.readFileSync(path.join(abs, f), "utf8");
      body += section(f, f, code) + "\n";
    }
  }

  const solutionPath = path.join(abs, "solution.go");
  if (fs.existsSync(solutionPath)) {
    body += "## Эталонное решение\n\n";
    body +=
      "Сначала попробуйте сами. Решение скрыто — раскрывайте только для сверки.\n\n";
    const code = fs.readFileSync(solutionPath, "utf8");
    body += solutionSpoiler("solution.go", code) + "\n";
  }

  body += `_Пакет \`${module}/${relDir}\`, упражнение **${pkg}**. Локально: \`go test -tags solution ./${relDir}/\`._\n`;

  const outPath = path.join(abs, "PAGE.md");
  fs.writeFileSync(outPath, body);
  console.log(`wrote ${relDir}/PAGE.md`);
}

for (const dir of listChapterDirs()) {
  expandChapter(dir);
}
