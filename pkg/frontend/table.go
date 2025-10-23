package frontend

import (
	"fmt"
	"slices"
	"strconv"
	"strings"

	"github.com/aknea001/goDoList/pkg"
)

func line(lgstID int, lgstT int, lgstD int) {
	fmt.Printf("+%s+%s+%s+\n",
		strings.Repeat("-", lgstID+2),
		strings.Repeat("-", lgstT+2),
		strings.Repeat("-", lgstD+2),
	)
}

func DrawTable(tasks []pkg.Task) {
	IDlens := make([]int, 0)
	titleLens := make([]int, 0)
	descLens := make([]int, 0)

	for i := range tasks {
		titleLen := len(tasks[i].Title)
		descLen := len(tasks[i].Description)

		if descLen > 50 {
			descLen = 53
		}

		titleLens = append(titleLens, titleLen)
		descLens = append(descLens, descLen)
	}

	// add the len of column titles
	IDlens = append(IDlens, 2)
	titleLens = append(titleLens, 5)
	descLens = append(descLens, 11)

	digitsOfTasksLen := len(strconv.Itoa(len(tasks)))
	IDlens = append(IDlens, digitsOfTasksLen)

	longestID := slices.Max(IDlens)
	longestTitle := slices.Max(titleLens)
	longestDesc := slices.Max(descLens)

	line(longestID, longestTitle, longestDesc)

	fmt.Printf("| ID%s | Title%s | Description%s |\n",
		strings.Repeat(" ", longestID-2),
		strings.Repeat(" ", longestTitle-5),
		strings.Repeat(" ", longestDesc-11),
	)

	line(longestID, longestTitle, longestDesc)

	if len(tasks) == 0 {
		fmt.Println("Type 'new' to make a task")
	}

	for i := range tasks {
		currentTitle := tasks[i].Title
		currentDesc := tasks[i].Description

		if longestDesc == 53 && len(currentDesc) > 50 {
			newCurrentDesc := currentDesc[:50]
			currentDesc = newCurrentDesc + "..."
		}

		fmt.Printf("| %s%s | %s%s | %s%s |\n",
			strconv.Itoa(i+1), strings.Repeat(" ", longestID-len(strconv.Itoa(i+1))),
			currentTitle, strings.Repeat(" ", longestTitle-len(currentTitle)),
			currentDesc, strings.Repeat(" ", longestDesc-len(currentDesc)),
		)
	}

	line(longestID, longestTitle, longestDesc)
}

func DrawOneTask(id int, task pkg.Task) {
	longestID := max(len("ID"), len(strconv.Itoa(id)))
	longestTitle := max(len("title"), len(task.Title))

	longestDesc := max(len("Description"), len(task.Description))
	longestDesc = min(longestDesc, 50)

	line(longestID, longestTitle, longestDesc)
	fmt.Printf("| ID%s | Title%s | Description%s |\n",
		strings.Repeat(" ", longestID-2),
		strings.Repeat(" ", longestTitle-5),
		strings.Repeat(" ", longestDesc-11),
	)
	line(longestID, longestTitle, longestDesc)

	descSlice := strings.SplitSeq(task.Description, " ")
	descLines := make([]string, 0)

	var newLine string
	for element := range descSlice {
		if len(newLine)+len(element)+1 > 51 {
			descLines = append(descLines, newLine)

			newLine = element + " "
			continue
		}

		newLine += element

		if len(newLine) != 50 {
			newLine += " "
		}
	}

	if newLine != "" {
		descLines = append(descLines, newLine)
	}

	mid := len(descLines) / 2

	for i, element := range descLines {
		idToPrint := ""
		titleToPrint := ""

		if i == mid {
			idToPrint = strconv.Itoa(id)
			titleToPrint = task.Title
		}

		fmt.Printf("| %s%s | %s%s | %s%s |\n",
			idToPrint, strings.Repeat(" ", longestID-len(idToPrint)),
			titleToPrint, strings.Repeat(" ", longestTitle-len(titleToPrint)),
			element, strings.Repeat(" ", longestDesc-len(element)),
		)
	}

	line(longestID, longestTitle, longestDesc)
}
