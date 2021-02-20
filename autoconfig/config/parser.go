package config

import (
	"github.com/sirupsen/logrus"
	"github.com/thingworks/common/utils/strings2"
	"os"
	"reflect"
	"regexp"
)

type configParserMap map[reflect.Type]ConfigurationParser

type ConfigurationParser interface {
	ParseConfig(interface{})
}

type DefaultConfigParser struct {
	configParserMap
}

type MongoConfigParser struct {
	reg *regexp.Regexp
}

func NewMongoConfigParser() *MongoConfigParser {
	return &MongoConfigParser{
		reg: regexp.MustCompile("\\$\\{(?P<env>[^}]*)\\}"),
	}
}

func (mongoParser *MongoConfigParser) ParseConfig(config interface{}) {
	groupNames := mongoParser.reg.SubexpNames()
	value := reflect.ValueOf(config).Elem()

	for i := 0; i < value.NumField(); i++ {
		field := value.Field(i)
		mongoParser.parseValue(field, groupNames)
	}
}

func (mongoParser *MongoConfigParser) parseValue(field reflect.Value, groupNames []string) {
	switch field.Kind() {
	case reflect.String:
		parseString(mongoParser.reg, field, groupNames, "env")
	case reflect.Struct:
		parseStruct(field, groupNames, mongoParser.parseValue)
	}
}

func NewDefaultConfigParser() DefaultConfigParser {
	return DefaultConfigParser{
		configParserMap{
			reflect.ValueOf(&MongoConfig{}).Type(): NewMongoConfigParser(),
		},
	}
}

func (defaultParser *DefaultConfigParser) ParseConfig(config interface{}) {

	if reflect.ValueOf(config).Kind() != reflect.Ptr {
		logrus.Debug("Only config in pointer form can be parsed")
		return
	}

	defaultParser.
		configParserMap[reflect.ValueOf(config).Type()].
		ParseConfig(config)
}

func parseString(reg *regexp.Regexp, field reflect.Value, groupNames []string, name string) {
	if reg.MatchString(field.String()) {
		submatch := reg.FindStringSubmatch(field.String())
		for index, groupName := range groupNames {
			if strings2.IsNotBlank(groupNames[index]) && strings2.Equals(groupName, name) {
				field.SetString(os.Getenv(submatch[index]))
			}
		}
	}
}

func parseStruct(field reflect.Value, groupNames []string, valueParser func(value reflect.Value, groupNames []string)) {
	for i := 0; i < field.NumField(); i++ {
		valueParser(field.Elem().Field(i), groupNames)
	}
}

type AppConfigParser struct {
	DefaultConfigParser
}

func NewAppConfigParser() *AppConfigParser {
	return &AppConfigParser{
		DefaultConfigParser: NewDefaultConfigParser(),
	}
}

func (appParser *AppConfigParser) ParseConfig(config *AppConfig) {
	appParser.DefaultConfigParser.ParseConfig(&config.Mongodb)
}
