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
	"to-persist/client/constants"
	"to-persist/client/form"
	"to-persist/client/global"
	"to-persist/client/util"
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
		fmt.Println(constants.InternalErrorReply)
		os.Exit(1)
	}

	toper.Period, err = cmd.Flags().GetString("period")
	if err != nil {
		zap.S().Error("failed to get period flag")
		fmt.Println(constants.InternalErrorReply)
		os.Exit(1)
	}

	// 10:30 or 18:30
	toper.DueDate, err = cmd.Flags().GetString("due-date")
	if err != nil {
		zap.S().Error("failed to get due-date flag")
		fmt.Println(constants.InternalErrorReply)
		os.Exit(1)
	}

	resp, err := util.Request2(http.MethodPost, global.ClientConfig.Url.Root+global.ClientConfig.Url.Toper, toper, nil, true)
	if err != nil {
		zap.S().Error(err)
		fmt.Println(constants.InternalErrorReply)
		os.Exit(1)
	}

	switch resp.StatusCode {
	case http.StatusOK:
		fmt.Println(constants.OKReply)
	case http.StatusBadRequest:
		fmt.Println(constants.BadRequestReply)
	case http.StatusUnauthorized:
		fmt.Println(constants.UnauthorizedReply)
	default:
		fmt.Println(constants.DefaultErrorReply)
	}

}

// toper done rsc ng ...

func Done(cmd *cobra.Command, args []string) {
	var request []form.DoneRequest
	for _, arg := range args {
		request = append(request, form.DoneRequest{Acronym: arg})
	}

	resp, err := util.Request2(http.MethodPost, global.ClientConfig.Url.Root+global.ClientConfig.Url.Done, request, nil, true)
	if err != nil {
		zap.S().Errorf("failed to get response from done api: %v", err)
		fmt.Println(constants.InternalErrorReply)
		os.Exit(1)
	}

	switch resp.StatusCode {
	case http.StatusOK:
		fmt.Println(constants.OKReply)
	case http.StatusBadRequest:
		fmt.Println(constants.BadRequestReply)
	case http.StatusUnauthorized:
		fmt.Println(constants.UnauthorizedReply)
	default:
		fmt.Println(constants.DefaultErrorReply)
	}

}

func List(cmd *cobra.Command, args []string) {
	resp, err := util.Request2(http.MethodGet, global.ClientConfig.Url.Root+global.ClientConfig.Url.Toper, nil, nil, true)
	if err != nil {
		zap.S().Errorf("failed to invoke list api: %v", err)
		fmt.Println(constants.InternalErrorReply)
		os.Exit(1)
	}

	switch resp.StatusCode {
	case http.StatusOK:
		respJson, err := io.ReadAll(resp.Body)
		if err != nil {
			zap.S().Error(err)
			fmt.Println(constants.InternalErrorReply)
			os.Exit(1)
		}

		listResponses := make([]ListResp, 0)
		err = json.Unmarshal(respJson, &listResponses)
		if err != nil {
			zap.S().Error(err)
			fmt.Println(constants.InternalErrorReply)
			os.Exit(1)
		}
		showInTable(listResponses)
	case http.StatusUnauthorized:
		fmt.Println(constants.UnauthorizedReply)
	case http.StatusBadRequest:
		fmt.Println(constants.BadRequestReply)
	default:
		fmt.Println(constants.DefaultErrorReply)
	}
}

// toper history ng -n 20

type historyResp struct {
	ID       int    `json:"-"`
	DoneTime string `json:"done-time"`
	ToperID  string `json:"toper-id"`
	Acronym  string `json:"acronym"`
	Done     string `json:"done"`
}

func History(cmd *cobra.Command, args []string) {
	historyUrl := global.ClientConfig.Url.History
	limit, err := cmd.Flags().GetString("limit")
	if err != nil {
		zap.S().Error("failed to get limit flag")
		fmt.Println(constants.InternalErrorReply)
		os.Exit(1)
	}
	queryParams := map[string]string{
		"acronym": args[0],
		"limit":   limit,
	}
	resp, err := util.Request2(http.MethodGet, global.ClientConfig.Url.Root+historyUrl, nil, queryParams, true)
	if err != nil {
		zap.S().Error(err)
		fmt.Println(constants.InternalErrorReply)
		os.Exit(1)
	}

	switch resp.StatusCode {
	case http.StatusOK:
		respJson, err := io.ReadAll(resp.Body)
		if err != nil {
			zap.S().Error(err)
			fmt.Println(constants.InternalErrorReply)
			os.Exit(1)
		}

		historyResponses := make([]historyResp, 0)
		err = json.Unmarshal(respJson, &historyResponses)
		if err != nil {
			zap.S().Error(err)
			fmt.Println(constants.InternalErrorReply)
			os.Exit(1)
		}
		table := tablewriter.NewWriter(os.Stdout)

		if len(historyResponses) == 0 {
			fmt.Println(constants.NoDataReply)
			return
		}

		headers := getStructFieldNames(historyResponses[0])
		table.SetHeader(headers)
		for i, data := range historyResponses {
			id := strconv.FormatInt(int64(i+1), 10)
			table.Append([]string{id, data.DoneTime, data.ToperID, data.Acronym, data.Done})
		}
		table.Render()
	case http.StatusUnauthorized:
		fmt.Println(constants.UnauthorizedReply)
	case http.StatusBadRequest:
		fmt.Println(constants.BadRequestReply)
	default:
		fmt.Println(constants.DefaultErrorReply)
	}
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
