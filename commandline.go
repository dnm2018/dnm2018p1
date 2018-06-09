package main

import (
	"fmt"

	"github.com/dnm2018/dnm2018p1/tree/master/modules/marketplace"
	"github.com/dnm2018/dnm2018p1/tree/master/modules/util"
)

func ManageRole(username, action, role string) {
	user, _ := marketplace.FindUserByUsername(username)
	if user == nil {
		fmt.Println("No such user")
		return
	}
	if action == "grant" && role == "seller" {
		user.IsSeller = !user.IsSeller
	} else if action == "grant" && role == "admin" {
		user.IsAdmin = !user.IsAdmin
	} else {
		fmt.Println("Wrong action")
		return
	}
	user.Save()
}

func RemoveUser(username string) {
	user, _ := marketplace.FindUserByUsername(username)
	if user == nil {
		panic("user not found")
	}
	user.Remove()
}

func IndexItems() {
	util.Log.Debug("[Index] Indexing items...")
	for _, item := range marketplace.GetAllItems() {
		util.Log.Debug("[Index] %s", item.Name)
		err := item.Index()
		if err != nil {
			util.Log.Error("%s", err)
		}
	}
}

func SearchItems(text string) {
	util.Log.Debug("[Index] Searching items...")
	marketplace.SearchItems(text)
}

func syncModels() {
	marketplace.SyncModels()

}
func updateStalledTransactions() {
	util.Log.Debug("[Transactions] UpdatingStalledTransactions")
	marketplace.TaskUpdateBalancesOrRecentlyReleasedAndCancelledTransactions()
	marketplace.TaskFinalizeReleasedAndCancelledTransactionsWithNonZeroAmount()
}

func updateOldAndPending() {
	marketplace.UpdateOldFailedAndPendingTransactions()
}

func resendReleasedTransactions() {
	marketplace.ResendReleasedTransactions()
}

func testStats() {
	interval := "7 days"
	sItems, err := marketplace.StaffSupportTicketsResolutionStats(interval)
	if err != nil {
		return
	}

	var (
		text = fmt.Sprintf(`

	Support Agent | Ticket Status | Number Of Tickets
	--- | --- | ---
	`, interval)
	)
	for _, si := range sItems {
		text += fmt.Sprintf("%s | %s | %d\n", si.ResolverUsername, si.CurrentStatus, si.TicketCount)
	}

	println(text)
}
