#!/usr/bin/env node
/**
 * Validates all mermaid code blocks in the given markdown files (default: specs/c4).
 *
 * Usage:
 *   pnpm validate:mermaid            # specs/c4 (all .md files)
 *   pnpm validate:mermaid:all        # all project specs (specs glob)
 *   node scripts/validate-mermaid.mjs path/to/file.md
 *
 * Exits non-zero if any block fails to parse.
 */
import { readFileSync, readdirSync } from 'node:fs';
import { join, dirname, resolve, relative } from 'node:path';
import { fileURLToPath } from 'node:url';
import { globSync } from 'glob';
import { JSDOM } from 'jsdom';

// mermaid v10 requires a DOM environment (DOMPurify/sanitize) — provide jsdom globals.
const dom = new JSDOM('<!DOCTYPE html><html><body></body></html>');
global.window = dom.window;
global.document = dom.window.document;
Object.defineProperty(global, 'navigator', { value: dom.window.navigator, configurable: true });

const mermaid = (await import('mermaid')).default;
mermaid.initialize({ startOnLoad: false, securityLevel: 'loose' });

const root = dirname(dirname(fileURLToPath(import.meta.url)));
const defaultDir = join(root, 'specs', 'c4');
const args = process.argv.slice(2);

let files;
if (args.length > 0) {
  // normalize to forward slashes for glob (Windows path separators break glob patterns)
  files = globSync(args.map((p) => resolve(root, p).split('\\').join('/')));
} else {
  files = readdirSync(defaultDir).filter((f) => f.endsWith('.md')).map((f) => join(defaultDir, f));
}

if (files.length === 0) {
  console.error('No markdown files matched.');
  process.exit(1);
}

let failed = 0;
let total = 0;

for (const file of files) {
  const content = readFileSync(file, 'utf8').replace(/\r\n/g, '\n');
  const regex = /```mermaid\n([\s\S]*?)```/g;
  let match;
  let idx = 0;
  while ((match = regex.exec(content)) !== null) {
    idx++;
    total++;
    const code = match[1];
    try {
      await mermaid.parse(code);
      console.log(`OK   ${relative(root, file)} [block ${idx}] (${code.trim().split('\n')[0]})`);
    } catch (e) {
      failed++;
      console.log(`FAIL ${relative(root, file)} [block ${idx}]: ${e.message?.split('\n')[0] ?? e}`);
    }
  }
}

console.log(`\nChecked ${total} block(s) in ${files.length} file(s), ${failed} failed.`);
process.exit(failed === 0 ? 0 : 1);
