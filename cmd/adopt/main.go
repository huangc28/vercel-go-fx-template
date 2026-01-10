package main

import (
	"embed"
	"flag"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"time"
)

//go:embed assets/AGENTS.md
//go:embed assets/architecture/go-vercel-reusable-template-plan.md
//go:embed assets/codex/skills/adopt/SKILL.md
var assets embed.FS

type options struct {
	dir            string
	force          bool
	dryRun         bool
	noAgents       bool
	noArchitecture bool
	noSkill        bool
}

type fileToWrite struct {
	assetPath string
	outPath   string
}

func main() {
	var opt options
	flag.StringVar(&opt.dir, "dir", ".", "target repository directory")
	flag.BoolVar(&opt.force, "force", false, "overwrite existing files (backs up originals)")
	flag.BoolVar(&opt.dryRun, "dry-run", false, "print actions without writing")
	flag.BoolVar(&opt.noAgents, "no-agents", false, "skip writing AGENTS.md")
	flag.BoolVar(&opt.noArchitecture, "no-architecture", false, "skip writing architecture plan doc")
	flag.BoolVar(&opt.noSkill, "no-skill", false, "skip writing codex /adopt skill into codex/skills/")
	flag.Parse()

	targetRoot, err := filepath.Abs(opt.dir)
	must(err)

	if info, err := os.Stat(targetRoot); err != nil || !info.IsDir() {
		fatalf("target dir does not exist or is not a directory: %s", targetRoot)
	}

	var files []fileToWrite
	if !opt.noAgents {
		files = append(files, fileToWrite{
			assetPath: "assets/AGENTS.md",
			outPath:   "AGENTS.md",
		})
	}
	if !opt.noArchitecture {
		files = append(files, fileToWrite{
			assetPath: "assets/architecture/go-vercel-reusable-template-plan.md",
			outPath:   "architecture/go-vercel-reusable-template-plan.md",
		})
	}
	if !opt.noSkill {
		files = append(files, fileToWrite{
			assetPath: "assets/codex/skills/adopt/SKILL.md",
			outPath:   "codex/skills/adopt/SKILL.md",
		})
	}

	backupRoot := ""
	if opt.force {
		backupRoot = filepath.Join(targetRoot, ".adopt-backup-"+time.Now().Format("20060102_150405"))
	}

	for _, f := range files {
		data, err := fs.ReadFile(assets, f.assetPath)
		must(err)

		dest := filepath.Join(targetRoot, filepath.FromSlash(f.outPath))
		exists := fileExists(dest)

		switch {
		case exists && !opt.force:
			fmt.Printf("skip (exists): %s\n", f.outPath)
			continue
		case exists && opt.force:
			backup := filepath.Join(backupRoot, filepath.FromSlash(f.outPath))
			if opt.dryRun {
				fmt.Printf("backup: %s -> %s\n", f.outPath, filepath.ToSlash(strings.TrimPrefix(backup, targetRoot+string(os.PathSeparator))))
			} else {
				must(os.MkdirAll(filepath.Dir(backup), 0o755))
				must(os.Rename(dest, backup))
			}
		}

		if opt.dryRun {
			fmt.Printf("write: %s\n", f.outPath)
			continue
		}

		must(os.MkdirAll(filepath.Dir(dest), 0o755))
		must(os.WriteFile(dest, data, 0o644))
		fmt.Printf("write: %s\n", f.outPath)
	}

	if opt.force && !opt.dryRun {
		fmt.Printf("backup dir: %s\n", backupRoot)
	}
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func must(err error) {
	if err != nil {
		fatalf("%v", err)
	}
}

func fatalf(format string, args ...any) {
	fmt.Fprintf(os.Stderr, "adopt: "+format+"\n", args...)
	os.Exit(1)
}
