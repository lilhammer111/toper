package command

import (
	"github.com/spf13/cobra"
	"time"
	"to-persist/client/handler"
	"to-persist/client/model"
)

func init() {
	rootCmd.AddCommand(persistCmd, doneCmd)

	persistCmd.Flags().DurationVarP(&toperFlags.Period, "period", "p", 24*time.Hour, "Duration between each perseverance")
	persistCmd.Flags().StringVarP(&toperFlags.DueDate, "due-date", "d", "", "Due date with a specific time (format: 2006-01-02T15:04:05)")
	persistCmd.Flags().StringVarP(&toperFlags.Acronym, "acronym", "a", "", "")
}

var (
	toperFlags model.ToperFlags
	// toper persist "reading excellent open source projects" -a rsc
	persistCmd = &cobra.Command{
		Use:   "persist <doing something>",
		Short: "Set a daily perseverance item",
		Long: `The 'persist' command allows you to set a daily perseverance item, 
				emphasizing consistency and commitment.
				You can specify the duration between each perseverance 
				and a due date for a specific target.`,
		Run: handler.Set,
	}

	// toper done rsc ng ...
	doneCmd = &cobra.Command{
		Use:   "done <item1 acronym> <item2 acronym> ...",
		Short: "Mark a daily perseverance item as completed",
		Long: `The 'done' command allows you to mark a specific daily perseverance item as completed for today. 
				This helps in tracking your consistency and commitment towards the set goals.`,
		Run: handler.Done,
	}

	//     过滤功能：您可能会有很多事项，特别是经过一段时间后。能够基于状态（例如只显示未完成的事项）进行过滤将会很有用。
	//    排序功能：默认情况下，事项可能按添加的顺序或最近完成的时间排序。但提供其他排序选项（如按名称、按重要性等）可能也会很有用。
	//    持续时间显示：为每个事项显示一个统计，例如“已坚持xx天”，这可以作为一种鼓励，让用户看到他们的进度。
	//    彩色高亮：可以使用不同的颜色来高亮已完成和未完成的事项，使其更加直观。
	//    简短的摘要：除了事项名称，可以为每个事项提供一个简短的摘要或描述，帮助用户记住为什么他们选择了这个坚持事项。
	listCmd = &cobra.Command{
		Use:   "list",
		Short: "List all perseverance items",
		Long: `The 'list' command displays all the daily perseverance items you've set. 
				For each item, it will indicate whether it has been completed for today or not. 
				This provides a quick overview of your daily commitments and their statuses.`,
		Run: handler.List,
	}

	//    展示详情
	//    标题：这是该坚持事项的名称或描述。
	//    简称：您提到的简称，以便快速引用。
	//    创建日期：何时开始这个坚持事项。
	//    最后完成日期：上次标记为完成的日期。
	//    总共坚持的天数：从开始到现在，成功坚持的天数。
	//    是否已完成：今天是否已经完成了这个任务。
	//    笔记或注解：与该坚持事项相关的任何额外信息或注释。

	// toper detail rsc
	detailCmd = &cobra.Command{
		Use:   "detail",
		Short: "Display details of a perseverance item",
		Long: `The 'detail' command provides an in-depth view of a specific perseverance item.
				It showcases all the information associated with the item, 
				such as its title, creation date,
				last completed date, total days persisted, and any notes or annotations.`,
		Run: handler.List,
	}
)
