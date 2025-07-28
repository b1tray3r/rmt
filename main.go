package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/b1tray3r/rmt/internal/http/handler"
	"github.com/b1tray3r/rmt/internal/http/server"
	"github.com/b1tray3r/rmt/internal/tui"
	tea "github.com/charmbracelet/bubbletea"
)

func getFirstAndLastDayOfMonth() (time.Time, time.Time) {
	now := time.Now()

	year, month, _ := now.Date()
	location := now.Location()

	firstOfMonth := time.Date(year, month, 1, 0, 0, 0, 0, location)

	firstOfNextMonth := firstOfMonth.AddDate(0, 1, 0)
	lastOfMonth := firstOfNextMonth.AddDate(0, 0, -1)

	return firstOfMonth, lastOfMonth
}

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo, // Set log level to Info
	}))

	redmineHandler := handler.NewRedmineIssueHandler(logger)

	server := &http.Server{
		Addr: ":8080",
		Handler: server.New(
			map[string]server.HTTPHandler{
				"/":                 handler.NewHealthHandler(logger),
				"GET /fetch/issues": redmineHandler,
			},
			logger,
		),
	}
	go func() {
		logger.Info("server starting...")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Printf("Error starting server: %v\n", err)
			os.Exit(1)
		}
	}()

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	// Wait for interrupt signal
	<-ctx.Done()
	logger.Info("got interruption signal")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(shutdownCtx); err != nil {
		logger.Info("server shutdown returned an error", "error", err)
	}

	// Print final message and exit cleanly
	logger.Info("shutting down server...")
	os.Exit(0) // Exit with status 0 for a clean shutdown

	/*
		viper.SetConfigName("config")
		viper.SetConfigType("yml")
		viper.AddConfigPath("$HOME/.config/rmt")

		err := viper.ReadInConfig()
		if err != nil {
			panic(fmt.Errorf("fatal error config file: %w", err))
		}

		client := rm.NewRedmineAPI(
			viper.GetString("redmine.url"),
			viper.GetString("redmine.token"),
		)

		userID := viper.GetString("redmine.user_id")
			f, l := getFirstAndLastDayOfMonth()

			fmt.Printf("Fetching time entries for user %s from %s to %s\n", userID, f.Format("2006-01-02"), l.Format("2006-01-02"))
			filter := redmine.NewFilter("user_id", userID, "from", f.Format("2006-01-02"), "to", l.Format("2006-01-02"))
			client.GetClient().Limit = 10
			te, err := client.GetClient().TimeEntriesWithFilter(*filter)
			if err != nil {
				fmt.Printf("Error fetching time entries: %v\n", err)
				return
			}

			headerFmt := color.New(color.FgGreen, color.Underline).SprintfFunc()
			columnFmt := color.New(color.FgYellow).SprintfFunc()

			tbl := table.New("ID", "DateTime", "User", "Hours", "Comment")
			tbl.WithHeaderFormatter(headerFmt).WithFirstColumnFormatter(columnFmt)

			totalHours := float32(0)
			for _, result := range te {
				totalHours += result.Hours
				tbl.AddRow(result.Issue.Id, result.SpentOn, result.User.Name, result.Hours, result.Comments)
			}

			fmt.Printf("Total hours for user %s in July 2025: %.2f\n", userID, totalHours)
			tbl.Print()
	*/
	/*
		l, err := client.Search(redmine.SearchParams{
			Query:      "middleware",
			TitlesOnly: false,
			OpenIssues: false,
			Limit:      0,
		})

		if err != nil {
			fmt.Printf("Error during search: %v\n", err)
		}

		fmt.Printf("Search results found: %d\n", l.TotalCount)
		fmt.Printf("Limit is: %d\n", l.Limit)

		headerFmt := color.New(color.FgGreen, color.Underline).SprintfFunc()
		columnFmt := color.New(color.FgYellow).SprintfFunc()

		tbl := table.New("ID", "Title", "URL", "DateTime")
		tbl.WithHeaderFormatter(headerFmt).WithFirstColumnFormatter(columnFmt)

		for _, result := range l.Results {
			parts := strings.Split(result.Title, "#")
			titleParts := parts[1]
			tparts := strings.Split(titleParts, ": ")
			title := strings.Join(tparts[1:], ": ")

			tbl.AddRow(result.ID, title, result.URL, result.DateTime.Format("2006-01-02 15:04:05"))
		}

		tbl.Print()
	*/

	p := tea.NewProgram(tui.New())

	if _, err := p.Run(); err != nil {
		fmt.Printf("An error occured: %v", err)
		os.Exit(1)
	}
}
