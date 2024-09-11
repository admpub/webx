package xcommon

import (
	"fmt"
	"html/template"

	"github.com/coscms/webcore/library/config"
	"github.com/webx-top/com"
	"github.com/webx-top/echo"
)

// CurrencySymbols 货币符号
var CurrencySymbols = map[string]string{
	`CNY`: `¥`,   //人民币
	`USD`: `$`,   //美元
	`JPY`: `¥`,   //日元
	`EUR`: `€`,   //欧元
	`GBP`: `£`,   //英镑
	`FRF`: `₣`,   //法郎
	`KRW`: `₩`,   //韩元
	`RUB`: `₽`,   //俄罗斯卢布
	`HKD`: `HK$`, //港元
	`AUD`: `A$`,  //澳元
	`CAD`: `C$`,  //加元
	`INR`: `₹`,   //印度卢比
}

// SetCurrencySymbol 登记货币符号
func SetCurrencySymbol(currency string, symbol string) {
	CurrencySymbols[currency] = symbol
}

// DefaultCurrency 默认币种
func DefaultCurrency() string {
	defaultCurrency := config.Setting(`base`).String(`defaultCurrency`)
	if len(defaultCurrency) == 0 {
		defaultCurrency = `CNY`
	}
	return defaultCurrency
}

// DefaultCurrencySymbol 默认币种符号
func DefaultCurrencySymbol() string {
	currencySymbol, ok := CurrencySymbols[DefaultCurrency()]
	if !ok || len(currencySymbol) == 0 {
		currencySymbol = `¥`
	}
	return currencySymbol
}

// HTMLCurrency HTML模板函数：币种
func HTMLCurrency(ctx echo.Context, v float64, withFlags ...bool) interface{} {
	currencySymbol := DefaultCurrencySymbol()
	if currency := ctx.Internal().String(`currency`); len(currency) > 0 {
		if symbol, ok := CurrencySymbols[currency]; ok {
			currencySymbol = symbol
		}
	}
	var numberFormatted string
	if len(withFlags) > 0 {
		if len(withFlags) > 1 && withFlags[1] {
			numberFormatted = com.NumberFormat(v, 3)
			numberFormatted = com.NumberTrimZero(numberFormatted)
		} else {
			numberFormatted = fmt.Sprintf(`%v`, v)
		}
		if withFlags[0] {
			return template.HTML(currencySymbol + numberFormatted)
		}
	} else {
		numberFormatted = fmt.Sprintf(`%v`, v)
	}
	return numberFormatted
}

// HTMLCurrencySymbol HTML模板函数：币种符号
func HTMLCurrencySymbol(ctx echo.Context) template.HTML {
	currencySymbol := DefaultCurrencySymbol()
	if currency := ctx.Internal().String(`currency`); len(currency) > 0 {
		if symbol, ok := CurrencySymbols[currency]; ok {
			currencySymbol = symbol
		}
	}
	return template.HTML(currencySymbol)
}
