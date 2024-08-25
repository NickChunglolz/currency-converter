package main

import (
	"errors"
	"fmt"
	"log"
	"strconv"

	"github.com/NickChunglolz/currency-converter/client"
	"github.com/charmbracelet/huh"
	"github.com/joho/godotenv"
)

type Currency struct {
	Code   string
	Symbol string
}

func main() {

	godotenv.Load()

	client := client.NewClient()

	currencyReplies, _ := client.GetCurrencies()
	currencyOptions := generateCurrencyOptions(currencyReplies)

	var baseCurrency *Currency
	var optCurrency *Currency
	var amount string
	var result float64

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[*Currency]().
				Title("What is you base currency?").
				Value(&baseCurrency).
				Options(currencyOptions...),

			huh.NewSelect[*Currency]().
				Title("What currency do you want to convert into?").
				Value(&optCurrency).
				Options(currencyOptions...),

			huh.NewInput().
				Title("How much do you want to convert?").
				Value(&amount).
				Placeholder("Enter amount of money...").
				Validate(func(str string) error {
					if _, err := strconv.ParseFloat(amount, 64); err != nil {
						return errors.New("wrong format input")
					}
					return nil
				}),
		),
	)

	err := form.Run()

	if err != nil {
		log.Fatal("Trouble int converting currency convert:", err)
	}

	rateReplies, _ := client.GetRates(baseCurrency.Code, optCurrency.Code)
	amountValue, _ := strconv.ParseFloat(amount, 64)

	result = rateReplies[0].Rate * amountValue

	fmt.Printf("The %s%s equal to %s%.2f \n", baseCurrency.Symbol, amount, optCurrency.Symbol, result)
}

func generateCurrencyOptions(currencyReplies []client.CurrencyReply) []huh.Option[*Currency] {
	var options []huh.Option[*Currency]

	for _, r := range currencyReplies {
		options = append(options, huh.NewOption(fmt.Sprintf("%s(%s)", r.Code, r.Symbol), &Currency{Code: r.Code, Symbol: r.Symbol}))
	}

	return options
}
