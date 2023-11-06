package handler

import (
	"encoding/json"
	"fmt"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"io"
	"net/http"
	"os"
	"reflect"
	"strconv"
	"strings"
	"to-persist/client/constant"
	"to-persist/client/global"
	"to-persist/client/util"
)

const (
	Sunday    = 1 << iota // 1 (binary: 0000001)
	Monday                // 2 (binary: 0000010)
	Tuesday               // 4 (binary: 0000100)
	Wednesday             // 8 (binary: 0001000)
	Thursday              // 16 (binary: 0010000)
	Friday                // 32 (binary: 0100000)
	Saturday              // 64 (binary: 1000000)
	Everyday
)

type Toper struct {
	Description string `json:"description,omitempty"`
	Acronym     string `json:"acronym,omitempty"`
	DueDate     string `json:"due-date"`
	Period      string `json:"period,omitempty"`
}

type ListResp struct {
	ID      int    `json:"id,omitempty"`
	Acronym string `json:"acronym,omitempty"`
	Desc    string `json:"desc,omitempty"`
	DueDate string `json:"due-date,omitempty"`
	Period  string `json:"period,omitempty"`
	Done    string `json:"done,omitempty"`
}

// toper persist "reading excellent open source projects" -a rsc

func Create(cmd *cobra.Command, args []string) {
	toper := Toper{}
	var err error

	toper.Description = strings.Join(args, " ")

	toper.Acronym, err = cmd.Flags().GetString("acronym")
	if err != nil {
		zap.S().Error("failed to get acronym flag")
		fmt.Println(constant.InternalError)
		os.Exit(1)
	}

	toper.Period, err = cmd.Flags().GetString("period")
	if err != nil {
		zap.S().Error("failed to get period flag")
		fmt.Println(constant.InternalError)
		os.Exit(1)
	}

	// 10:30 or 18:30
	toper.DueDate, err = cmd.Flags().GetString("due-date")
	if err != nil {
		zap.S().Error("failed to get due-date flag")
		fmt.Println(constant.InternalError)
		os.Exit(1)
	}

	resp, err := util.Request2(http.MethodPost, global.Config.Url.Base+global.Config.Url.Toper, toper, true)
	if err != nil {
		zap.S().Error(err)
		fmt.Println(constant.InternalError)
		os.Exit(1)
	}

	switch resp.StatusCode {
	case http.StatusOK:
		fmt.Println("The toper was created successfully, and you can view your toper with the list command.")
	case http.StatusBadRequest:
		fmt.Println("Bad request. Please check the information you provided.")
	case http.StatusUnauthorized:
		fmt.Println("Authentication failed. Please login first.")
	default:
		fmt.Println("An unexpected error occurred. Please try again later.")
	}

}

func Done(cmd *cobra.Command, args []string) {

}

func List(cmd *cobra.Command, args []string) {
	resp, err := util.Request2(http.MethodGet, global.Config.Url.Base+global.Config.Url.Toper, nil, true)
	if err != nil {
		zap.S().Error(err)
		fmt.Println(constant.InternalError)
		os.Exit(1)
	}

	respJson, err := io.ReadAll(resp.Body)
	if err != nil {
		zap.S().Error(err)
		fmt.Println(constant.InternalError)
		os.Exit(1)
	}

	listResponses := make([]ListResp, 0)
	err = json.Unmarshal(respJson, &listResponses)
	if err != nil {
		zap.S().Error(err)
		fmt.Println(constant.InternalError)
		os.Exit(1)
	}

	switch resp.StatusCode {
	case http.StatusOK:
		showInTable(listResponses)
	case http.StatusUnauthorized:
		fmt.Println("Authentication failed, please login first.")
	default:
		fmt.Println("An unexpected error occurred. Please try again later.")
	}
}

func History(cmd *cobra.Command, args []string) {

}

func showInTable(listResponses []ListResp) {
	table := tablewriter.NewWriter(os.Stdout)
	headers := getStructFieldNames(listResponses[0])
	table.SetHeader(headers)
	for _, data := range listResponses {
		id := strconv.FormatInt(int64(data.ID), 10)
		table.Append([]string{id, data.Acronym, data.Desc, data.Period, data.DueDate, data.Done})
	}

	table.Render()
}

func getStructFieldNames(item interface{}) []string {
	t := reflect.TypeOf(item)
	var names []string
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		names = append(names, field.Name)
	}

	return names
}

var periodMeanings = map[int]string{
	Sunday:    "Sunday",
	Monday:    "Monday",
	Tuesday:   "Tuesday",
	Wednesday: "Wednesday",
	Thursday:  "Thursday",
	Friday:    "Friday",
	Saturday:  "Saturday",
}

var periodFlag = map[string]int{
	"sunday":    Sunday,
	"monday":    Monday,
	"tuesday":   Tuesday,
	"wednesday": Wednesday,
	"thursday":  Thursday,
	"friday":    Friday,
	"saturday":  Saturday,
}

func periodToString(period int) string {
	var periods []string
	if period&Sunday != 0 {
		periods = append(periods, periodMeanings[Sunday])
	}
	if period&Monday != 0 {
		periods = append(periods, periodMeanings[Monday])
	}
	if period&Tuesday != 0 {
		periods = append(periods, periodMeanings[Tuesday])
	}
	if period&Wednesday != 0 {
		periods = append(periods, periodMeanings[Wednesday])
	}
	if period&Thursday != 0 {
		periods = append(periods, periodMeanings[Thursday])
	}
	if period&Friday != 0 {
		periods = append(periods, periodMeanings[Friday])
	}
	if period&Saturday != 0 {
		periods = append(periods, periodMeanings[Saturday])
	}

	return strings.Join(periods, ", ")
}
