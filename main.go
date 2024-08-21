package main

import (
	"errors"
	"fmt"
	"os"
	"strconv"

	"github.com/charmbracelet/huh"
)

const (
	twd = "TWD"
	cad = "CAD"
	usd = "USD"
	gbp = "GBP"
	eur = "EUR"
	jpy = "JPY"
)

var currenyOptions = []huh.Option[string]{
	huh.NewOption("NTD", twd),
	huh.NewOption("CAD", cad),
	huh.NewOption("USD", usd),
	huh.NewOption("GBP", gbp),
	huh.NewOption("EUR", eur),
	huh.NewOption("JYP", jpy),
}

func main() {

	var baseCurrency string
	var optCurrency string
	var amount string
	var result float64

	err :=
		huh.NewForm(
			huh.NewGroup(
				huh.NewSelect[string]().
					Title("What is you base currency?").
					Value(&baseCurrency).
					Options(currenyOptions...),

				huh.NewSelect[string]().
					Title("What currency do you want to convert into?").
					Value(&optCurrency).
					Options(currenyOptions...),

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
		).Run()

	if err != nil {
		fmt.Println("Trouble int converting currency convert:", err)
		os.Exit(1)
	}

	amountValue, _ := strconv.ParseFloat(amount, 64)
	result = amountValue * 32

	fmt.Printf("The %s %s equal to %s %.2f \n", baseCurrency, amount, optCurrency, result)
}
