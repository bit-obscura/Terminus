package main

import (
	"context"
	"fmt"
	"strconv"

	"github.com/algorand/go-algorand-sdk/v2/client/v2/algod"
	"github.com/algorand/go-algorand-sdk/v2/crypto"
	"github.com/algorand/go-algorand-sdk/v2/types"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	appStyle = lipgloss.NewStyle().Padding(1, 2)

	titleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FAFAFA")).
			Background(lipgloss.Color("#1E90FF")).
			Padding(0, 1)

	itemStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FAFAFA")).
			PaddingLeft(4)

	selectedItemStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#1E90FF")).
				PaddingLeft(2).
				Border(lipgloss.NormalBorder(), false, false, false, true).
				BorderForeground(lipgloss.Color("#1E90FF"))

	errorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FF0000")).
			PaddingLeft(4)
)

type item struct {
	title, desc string
}

func (i item) Title() string       { return i.title }
func (i item) Description() string { return i.desc }
func (i item) FilterValue() string { return i.title }

type model struct {
	list         list.Model
	accounts     []Account
	state        string
	txFromInput  textinput.Model
	txToInput    textinput.Model
	txAmtInput   textinput.Model
	selectedFrom int
	selectedTo   int
	algodClient  *algod.Client
	error        string
}

type Account struct {
	address types.Address
	sk      crypto.PrivateKey
}

func createAlgorandAccount() Account {
	account := crypto.GenerateAccount()
	return Account{
		address: account.Address,
		sk:      account.PrivateKey,
	}
}

func initialModel() model {
	// Connect to the Algorand network
	algodClient, err := algod.MakeClient(
		"http://localhost:4001", // Replace with your Algorand node address
		"aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa", // Replace with your API key
	)
	if err != nil {
		panic(fmt.Sprintf("Failed to create Algod client: %s", err))
	}

	accounts := []Account{
		createAlgorandAccount(),
		createAlgorandAccount(),
	}

	items := make([]list.Item, len(accounts))
	for i, acc := range accounts {
		items[i] = item{
			title: acc.address.String(),
			desc:  "Fetching balance...",
		}
	}

	l := list.New(items, list.NewDefaultDelegate(), 0, 0)
	l.Title = "Algorand Wallet Management"
	l.Styles.Title = titleStyle

	txFromInput := textinput.New()
	txFromInput.Placeholder = "From account index"
	txToInput := textinput.New()
	txToInput.Placeholder = "To account index"
	txAmtInput := textinput.New()
	txAmtInput.Placeholder = "Amount to transfer (in microAlgos)"

	m := model{
		list:        l,
		accounts:    accounts,
		state:       "list",
		txFromInput: txFromInput,
		txToInput:   txToInput,
		txAmtInput:  txAmtInput,
		algodClient: algodClient,
	}

	return m
}

func (m model) Init() tea.Cmd {
	return m.updateBalances
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch m.state {
		case "list":
			switch msg.String() {
			case "q", "ctrl+c":
				return m, tea.Quit
			case "a":
				newAccount := createAlgorandAccount()
				m.accounts = append(m.accounts, newAccount)
				m.list.InsertItem(len(m.accounts)-1, item{
					title: newAccount.address.String(),
					desc:  "Fetching balance...",
				})
				return m, m.updateBalances
			case "t":
				m.state = "transaction"
				return m, textinput.Blink
			case "r":
				return m, m.updateBalances
			}
		case "transaction":
			switch msg.String() {
			case "enter":
				fromIdx, _ := strconv.Atoi(m.txFromInput.Value())
				toIdx, _ := strconv.Atoi(m.txToInput.Value())
				amount, _ := strconv.ParseUint(m.txAmtInput.Value(), 10, 64)
				if fromIdx >= 0 && fromIdx < len(m.accounts) && toIdx >= 0 && toIdx < len(m.accounts) && amount > 0 {
					return m, m.sendTransaction(fromIdx, toIdx, amount)
				}
				m.error = "Invalid transaction parameters"
				m.state = "list"
				m.txFromInput.SetValue("")
				m.txToInput.SetValue("")
				m.txAmtInput.SetValue("")
			case "esc":
				m.state = "list"
				m.txFromInput.SetValue("")
				m.txToInput.SetValue("")
				m.txAmtInput.SetValue("")
			}
		}
	case tea.WindowSizeMsg:
		h, v := appStyle.GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v)
	case balanceUpdateMsg:
		m.updateListItems()
		m.error = ""
	case errorMsg:
		m.error = string(msg)
	}

	var cmd tea.Cmd
	switch m.state {
	case "list":
		m.list, cmd = m.list.Update(msg)
	case "transaction":
		m.txFromInput, cmd = m.txFromInput.Update(msg)
		m.txToInput, cmd = m.txToInput.Update(msg)
		m.txAmtInput, cmd = m.txAmtInput.Update(msg)
	}
	return m, cmd
}

type balanceUpdateMsg struct{}
type errorMsg string

func (m model) updateBalances() tea.Msg {
	for i, acc := range m.accounts {
		accountInfo, err := m.algodClient.AccountInformation(acc.address.String()).Do(context.Background())
		if err != nil {
			return errorMsg(fmt.Sprintf("Failed to fetch account info: %s", err))
		}
		m.accounts[i].balance = accountInfo.Amount
	}
	return balanceUpdateMsg{}
}

func (m *model) updateListItems() {
	items := make([]list.Item, len(m.accounts))
	for i, acc := range m.accounts {
		items[i] = item{
			title: acc.address.String(),
			desc:  fmt.Sprintf("Balance: %d microAlgos", acc.balance),
		}
	}
	m.list.SetItems(items)
}

func (m model) sendTransaction(fromIdx, toIdx int, amount uint64) tea.Cmd {
	return func() tea.Msg {
		fromAddr := m.accounts[fromIdx].address.String()
		toAddr := m.accounts[toIdx].address.String()

		sp, err := m.algodClient.SuggestedParams().Do(context.Background())
		if err != nil {
			return errorMsg(fmt.Sprintf("Failed to get suggested params: %s", err))
		}

		tx, err := types.MakePaymentTxn(fromAddr, toAddr, amount, nil, "", sp)
		if err != nil {
			return errorMsg(fmt.Sprintf("Failed to create transaction: %s", err))
		}

		signedTx, err := tx.Sign(m.accounts[fromIdx].sk)
		if err != nil {
			return errorMsg(fmt.Sprintf("Failed to sign transaction: %s", err))
		}

		txID, err := m.algodClient.SendRawTransaction(signedTx).Do(context.Background())
		if err != nil {
			return errorMsg(fmt.Sprintf("Failed to send transaction: %s", err))
		}

		_, err = algod.WaitForConfirmation(m.algodClient, txID, 4, context.Background())
		if err != nil {
			return errorMsg(fmt.Sprintf("Failed to confirm transaction: %s", err))
		}

		return m.updateBalances()
	}
}

func (m model) View() string {
	var content string
	switch m.state {
	case "list":
		content = m.list.View() + "\nPress (a) to add account, (t) to make transaction, (r) to refresh balances, (q) to quit"
	case "transaction":
		content = fmt.Sprintf(
			"Make transaction:\n\nFrom account index: %s\nTo account index: %s\nAmount (microAlgos): %s\n\n(Enter) to confirm, (Esc) to cancel",
			m.txFromInput.View(),
			m.txToInput.View(),
			m.txAmtInput.View(),
		)
	}
	if m.error != "" {
		content += "\n" + errorStyle.Render(m.error)
	}
	return appStyle.Render(content)
}

func main() {
	p := tea.NewProgram(initialModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("There was an error: %v", err)
	}
}