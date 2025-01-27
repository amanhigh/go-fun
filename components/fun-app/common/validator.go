package common

// Custom Name Validator for alphanum and _,-
import (
	"regexp"

	"github.com/go-playground/validator/v10"
	"github.com/rs/zerolog/log"
)

const NAME_REGEX = `^[0-9a-zA-Z- ]+$`

var nameMatcher, _ = regexp.Compile(NAME_REGEX)
var matcherMap = map[string]*regexp.Regexp{
	"person": nameMatcher,
}

/*
	Validate Person Name

Checks if Name is AlphaNum.
Allowed Chars -,_
*/
func NameValidator(fl validator.FieldLevel) (check bool) {
	/* Extract Entity Name */
	entityName := fl.Param()

	//Finding Corresponding Matcher
	if matcher, ok := matcherMap[entityName]; ok {
		name := fl.Field().String()
		//Apply Regex
		check = matcher.MatchString(name)

		log.Debug().Str("Entity", entityName).Str("Name", name).Bool("Check", check).Msg("NameValidator Error")
	} else {
		log.Error().Str("MatcherName", entityName).Msg("NameValidator Error")
	}
	return
}
