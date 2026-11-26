// คำสั่ง sync ดึงข้อมูลจาก API-Football มาลงฐานข้อมูล
//
//	go run ./cmd/sync leagues ดูรายการลีกของประเทศ (ค่าเริ่มต้น Thailand)
//	go run ./cmd/sync link -name "Thai League 1" -id 296
//	go run ./cmd/sync fixtures -from 2026-11-05 -to 2026-12-09
package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"text/tabwriter"
	"time"
	_ "time/tzdata" // ให้โหลดเขตเวลา Asia/Bangkok ได้แม้เครื่องไม่มีฐานข้อมูลเขตเวลา

	"football/internal/config"
	"football/internal/database"
	"football/internal/ingest"
	"football/internal/models"
	"football/internal/provider/apifootball"
)

func main() {
	if len(os.Args) < 2 {
		usage()
		os.Exit(2)
	}
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	if err := run(ctx, os.Args[1], os.Args[2:]); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}

func usage() {
	fmt.Fprintln(os.Stderr, `usage:
	sync leagues  [-country Thailand]
	sync link     -name "Thai League 1" -id 296
	sync fixtures [-season 2024 หรือ 2022-2026] [-from YYYY-MM-DD] [-to YYYY-MM-DD]`)
}

func run(ctx context.Context, cmd string, args []string) error {
	cfg, err := config.Load()
	if err != nil {
		return err
	}
	if cfg.APIFootballKey == "" && cmd != "link" {
		return errors.New("FOOTBALL_API_KEY is not set (ใส่ใน .env)")
	}
	client := apifootball.New(cfg.APIFootballKey, apifootball.WithBaseURL(cfg.APIFootballBaseURL))

	switch cmd {
	case "leagues":
		return runLeagues(ctx, client, args)
	case "link":
		return runLink(cfg, args)
	case "fixtures":
		return runFixtures(ctx, cfg, client, args)
	default:
		usage()
		return fmt.Errorf("unknown command %q", cmd)
	}
}

func runLeagues(ctx context.Context, c *apifootball.Client, args []string) error {
	fs := flag.NewFlagSet("leagues", flag.ContinueOnError)
	country := fs.String("country", "Thailand", "ชื่อประเทศ (ภาษาอังกฤษ)")
	if err := fs.Parse(args); err != nil {
		return err
	}
	leagues, err := c.Leagues(ctx, *country)
	if err != nil {
		return err
	}
	tw := tabwriter.NewWriter(os.Stdout, 0, 4, 2, ' ', 0)
	fmt.Fprintln(tw, "ID\tTYPE\tCURRENT SEASON\tNAME")
	for _, l := range leagues {
		fmt.Fprintf(tw, "%s\t%s\t%d\t%s\n", l.ExternalID, l.Type, l.CurrentSeason, l.Name)
	}
	if err := tw.Flush(); err != nil {
		return err
	}
	printQuota(c)
	return nil
}

func runLink(cfg *config.Config, args []string) error {
	fs := flag.NewFlagSet("link", flag.ContinueOnError)
	name := fs.String("name", "", "ชื่อลีกในฐานข้อมูลของเรา (ตรงตัวอักษร)")
	id := fs.String("id", "", "league id ของ API-Football")
	if err := fs.Parse(args); err != nil {
		return err
	}
	if *name == "" || *id == "" {
		return errors.New("ต้องระบุทั้ง -name และ -id")
	}
	db, err := database.Connect(cfg.DatabaseURL)
	if err != nil {
		return err
	}
	res := db.Model(&models.League{}).Where("name = ?", *name).Update("external_api_id", *id)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return fmt.Errorf("ไม่พบลีกชื่อ %q ในฐานข้อมูล", *name)
	}
	fmt.Printf("ผูก %q กับ API-Football league id %s แล้ว\n", *name, *id)
	return nil
}

func parseSeasons(s string, defaultSeason int) ([]int, error) {
	if s == "" {
		if defaultSeason != 0 {
			return []int{defaultSeason}, nil
		}
		return nil, errors.New("ต้องระบุ -season (เช่น -season 2022-2026 หรือ -season 2024) หรือตั้ง FOOTBALL_SEASON ใน .env")
	}

	// รองรับรูปแบบช่วงปี เช่น "2022-2026"
	if strings.Contains(s, "-") {
		parts := strings.Split(s, "-")
		if len(parts) == 2 {
			start, err1 := strconv.Atoi(strings.TrimSpace(parts[0]))
			end, err2 := strconv.Atoi(strings.TrimSpace(parts[1]))
			if err1 == nil && err2 == nil && start <= end {
				var out []int
				for yr := start; yr <= end; yr++ {
					out = append(out, yr)
				}
				return out, nil
			}
		}
	}

	// รองรับรูปแบบจุลภาค เช่น "2022,2023,2024" หรือปีเดี่ยว "2024"
	parts := strings.Split(s, ",")
	var out []int
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p == "" {
			continue
		}
		yr, err := strconv.Atoi(p)
		if err != nil {
			return nil, fmt.Errorf("รูปแบบปีฤดูกาลไม่ถูกต้อง %q: %w", s, err)
		}
		out = append(out, yr)
	}
	if len(out) == 0 {
		return nil, errors.New("ไม่พบปีฤดูกาลที่ถูกต้อง")
	}
	return out, nil
}

func runFixtures(ctx context.Context, cfg *config.Config, c *apifootball.Client, args []string) error {
	loc, err := time.LoadLocation(apifootball.DefaultTimezone)
	if err != nil {
		return err
	}

	fs := flag.NewFlagSet("fixtures", flag.ContinueOnError)
	fromS := fs.String("from", "", "วันเริ่ม YYYY-MM-DD (ถ้าไม่ระบุ จะดึงทั้งฤดูกาล)")
	toS := fs.String("to", "", "วันสิ้นสุด YYYY-MM-DD (ถ้าไม่ระบุ จะดึงทั้งฤดูกาล)")
	seasonS := fs.String("season", "", "ปีฤดูกาล เช่น 2024 หรือช่วงปี 2022-2026 (ค่าเริ่มต้นจาก FOOTBALL_SEASON ใน .env)")
	if err := fs.Parse(args); err != nil {
		return err
	}

	seasons, err := parseSeasons(*seasonS, cfg.Season)
	if err != nil {
		return err
	}

	var from, to time.Time
	if *fromS != "" {
		from, err = time.ParseInLocation("2006-01-02", *fromS, loc)
		if err != nil {
			return fmt.Errorf("-from: %w", err)
		}
	}
	if *toS != "" {
		to, err = time.ParseInLocation("2006-01-02", *toS, loc)
		if err != nil {
			return fmt.Errorf("-to: %w", err)
		}
	}

	db, err := database.Connect(cfg.DatabaseURL)
	if err != nil {
		return err
	}
	svc := ingest.New(db, c, slog.Default())

	var (
		totalRes   ingest.Result
		syncErrors []error
	)

	for _, season := range seasons {
		fmt.Printf("===> กำลังซิงก์ข้อมูลฤดูกาล %d...\n", season)

		res, err := svc.SyncFixtures(ctx, season, from, to)
		if err != nil {
			syncErrors = append(
				syncErrors,
				fmt.Errorf("season %d: %w", season, err),
			)
		}

		totalRes.Leagues += res.Leagues
		totalRes.TeamsCreated += res.TeamsCreated
		totalRes.MatchesCreated += res.MatchesCreated
		totalRes.MatchesUpdated += res.MatchesUpdated
		totalRes.Skipped += res.Skipped

		fmt.Printf(
			"   [ฤดูกาล %d] ลีก %d | ทีมใหม่ %d | แมตช์ใหม่ %d | แมตช์ที่อัปเดต %d | ข้าม %d\n",
			season,
			res.Leagues,
			res.TeamsCreated,
			res.MatchesCreated,
			res.MatchesUpdated,
			res.Skipped,
		)
	}

	fmt.Printf(
		"\nสรุปผลรวมทั้งหมด: ทีมใหม่ %d | แมตช์ใหม่ %d | แมตช์ที่อัปเดต %d | ข้าม %d\n",
		totalRes.TeamsCreated,
		totalRes.MatchesCreated,
		totalRes.MatchesUpdated,
		totalRes.Skipped,
	)
	printQuota(c)
	return errors.Join(syncErrors...)
}

func printQuota(c *apifootball.Client) {
	if day, min := c.Quota(); day >= 0 || min >= 0 {
		fmt.Printf("โควตาที่เหลือ: รายวัน %d | รายนาที %d (-1 = ไม่ทราบ)\n", day, min)
	}
}
