package main

import (
	"context"
	"encoding/csv"
	"flag"
	"fmt"
	"io/fs"
	"log/slog"
	"os"
	"path/filepath"
	"t-sumisaki/svnchecker"

	"github.com/google/subcommands"
)

func findFiles(root string) ([]string, error) {
	findList := []string{}

	err := filepath.WalkDir(root, func(path string, info fs.DirEntry, err error) error {
		if err != nil {
			return fmt.Errorf("failed filepath.WalkDir: %w", err)
		}

		if info.IsDir() {
			return nil
		}

		findList = append(findList, path)
		return nil
	})

	return findList, err
}

type locklistCmd struct {
	path   string
	output string
	limit  int
}

func (*locklistCmd) Name() string     { return "locklist" }
func (*locklistCmd) Synopsis() string { return "Print svn locked file list" }
func (*locklistCmd) Usage() string {
	return `locklist -p path [-o output_path]:
	Output svn locked file list
`
}

func (c *locklistCmd) SetFlags(f *flag.FlagSet) {
	f.StringVar(&c.path, "p", "", "Search root path")
	f.StringVar(&c.output, "o", "", "Report output path")
	f.IntVar(&c.limit, "l", 0, "Query limit")
}

func (c *locklistCmd) Execute(ctx context.Context, f *flag.FlagSet, _ ...any) subcommands.ExitStatus {

	files, err := findFiles(c.path)

	if err != nil {
		slog.Error("target file search error.", slog.Any("error", err), slog.String("root", c.path))
		return subcommands.ExitFailure
	}

	fout, err := os.Create(c.output)

	if err != nil {
		slog.Error("output file cannot create.", slog.Any("error", err), slog.String("output", c.output))
		return subcommands.ExitFailure
	}

	defer fout.Close()

	writer := csv.NewWriter(fout)

	writer.Write([]string{
		"path",
		"user",
		"last_changed_date",
	})

	slog.Info("start query", slog.Int("length", len(files)))

	count := 0

	for i, fp := range files {

		slog.Debug("getinfo", slog.String("path", fp))

		if c.limit > 0 && i > c.limit {
			break
		}

		if i%100 == 0 {
			slog.Info(fmt.Sprintf("%d file checked, %d file remaining", i, len(files)-i), slog.Int("current", i), slog.Int("remain", len(files)-i))
		}

		svninfo, err := svnchecker.GetInfo(fp)
		if err != nil {
			slog.Error("svninfo command error", slog.Any("error", err), slog.String("path", fp))
			return subcommands.ExitFailure
		}

		if svninfo == nil {
			continue
		}

		if svninfo.IsLocked() {
			// ロックされているので出力する
			if err := writer.Write([]string{
				svninfo.Path,
				svninfo.LockOwner,
				svninfo.LastChangedDate,
			}); err != nil {
				slog.Error("lock info output error", slog.Any("error", err))
				return subcommands.ExitFailure
			}

			count++
			if count%100 == 0 {
				writer.Flush()
			}
		}
	}
	writer.Flush()

	return subcommands.ExitSuccess
}

func main() {

	slog.SetLogLoggerLevel(slog.LevelDebug)

	subcommands.Register(subcommands.HelpCommand(), "")
	subcommands.Register(subcommands.FlagsCommand(), "")
	subcommands.Register(subcommands.CommandsCommand(), "")
	subcommands.Register(&locklistCmd{}, "")

	flag.Parse()

	ctx := context.Background()
	os.Exit(int(subcommands.Execute(ctx)))
}
