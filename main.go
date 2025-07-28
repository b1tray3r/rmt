package main

import (
	"fmt"

	rm "github.com/b1tray3r/rmt/internal/redmine"
	"github.com/fatih/color"
	"github.com/rodaine/table"

	"github.com/mattn/go-redmine"
	"github.com/spf13/viper"
)

func main() {
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
	filter := redmine.NewFilter("user_id", userID, "from", "2025-07-01", "to", "2025-07-31")
	client.GetClient().Limit = 500
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

	/*

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
	/*

		p := tea.NewProgram(tui.NewModel())

			if _, err := p.Run(); err != nil {
				fmt.Printf("An error occured: %v", err)
				os.Exit(1)
			}
	*/
}
