package apiutils

import "net/url"

var DefaultURLValuesGenerator = func(values url.Values) URLValuesGenerator {
	return URLValues(values)
}

type URLValues url.Values

func (u URLValues) URLValues(signGenerators ...func(url.Values) string) url.Values {
	formData := url.Values(u)
	var signGenerator func(url.Values) string
	if len(signGenerators) > 0 {
		signGenerator = signGenerators[0]
	} else {
		signGenerator = GenSign
	}
	if signGenerator != nil {
		sign := signGenerator(formData)
		formData.Set(`sign`, sign)
	}
	return formData
}
