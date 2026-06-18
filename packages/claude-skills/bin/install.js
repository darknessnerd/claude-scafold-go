#!/usr/bin/env node
// claude-skills install
// Copies skills from:
//   1. packages/claude-skills/skills/     (team-authored skills)
//   2. node_modules/@team/caveman-skill/  (versioned external skill via GitHub Packages)
// into .claude/skills/ in the consumer repo root.
//
// Usage:
//   npm run setup:claude           normal install
//   npm run setup:claude -- --dry-run   preview without writing

const fs = require('fs')
const path = require('path')

const dryRun = process.argv.includes('--dry-run')

// Target: .claude/skills/ at the repo root (where npm install was run)
const targetDir = path.join(process.cwd(), '.claude', 'skills')

if (!dryRun) {
  fs.mkdirSync(targetDir, { recursive: true })
}

// ── Sources ────────────────────────────────────────────────────────────────

const sources = [
  // 1. Team skills bundled in this package
  {
    label: '@team/claude-skills',
    dir: path.join(__dirname, '..', 'skills'),
  },
  // 2. caveman-skill — published separately to GitHub Packages (@team/caveman-skill)
  //    Resolved from node_modules of THIS package (not the consumer repo root)
  //    so the consumer never needs to install it directly.
  {
    label: '@team/caveman-skill',
    dir: resolveSkillDir('@team/caveman-skill'),
  },
]

function resolveSkillDir(pkgName) {
  // Walk up from this file to find node_modules
  const candidates = [
    path.join(__dirname, '..', 'node_modules', pkgName, 'skills'),
    path.join(__dirname, '..', '..', 'node_modules', pkgName, 'skills'),
    path.join(process.cwd(), 'node_modules', pkgName, 'skills'),
  ]
  return candidates.find(fs.existsSync) || null
}

// ── Copy ───────────────────────────────────────────────────────────────────

let total = 0

for (const source of sources) {
  if (!source.dir || !fs.existsSync(source.dir)) {
    console.warn(`WARN: skills dir not found for ${source.label} — skipping`)
    continue
  }

  const files = fs.readdirSync(source.dir).filter(f => f.endsWith('.md'))

  for (const file of files) {
    const src = path.join(source.dir, file)
    const dest = path.join(targetDir, file)
    if (dryRun) {
      console.log(`[dry-run] ${source.label}/${file} → .claude/skills/${file}`)
    } else {
      fs.copyFileSync(src, dest)
      console.log(`installed [${source.label}]: .claude/skills/${file}`)
    }
    total++
  }
}

console.log(dryRun ? `\nDry run — ${total} file(s) would be installed.` : `\nDone — ${total} skill(s) installed.`)
