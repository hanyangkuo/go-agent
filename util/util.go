package util

import (
	"errors"
	"github.com/go-ini/ini"
	log "github.com/sirupsen/logrus"
	"reflect"
	"strconv"

	//"log"
	"time"
)

func TimeTrack(start time.Time, name string) {
	elapsed := time.Since(start)
	log.WithFields(log.Fields{
		"category": "timeTrack",
		"key"     : name,
		"value"   : elapsed.Milliseconds(),
	}).Infof("%s took %s", name, elapsed)
	//log.Printf("%s took %s", name, elapsed)
}

func LoadConfig(filepath string, iface interface{}) error {
	// set interface default value by tags.
	err := SetDefault(iface)
	if err != nil {
		return err
	}
	// open config.ini.
	cfg, err := ini.Load(filepath)
	if err != nil {
		return err
	}
	// try load config from ini file.
	t := reflect.TypeOf(iface)
	if cfg.HasSection(t.Elem().Name()){
		if err = cfg.Section(t.Elem().Name()).MapTo(iface); err != nil {
			return err
		}
		return nil
	}
	if err = cfg.MapTo(iface); err != nil {
		return err
	}
	return nil
}

func SetDefault(v interface{}) error {
	typ := reflect.TypeOf(v)
	val := reflect.ValueOf(v)
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
		val = val.Elem()
	} else {
		return errors.New("not a pointer to a struct")
	}

	for i := 0; i < val.NumField(); i++ {
		if defaultValue := typ.Field(i).Tag.Get("default"); defaultValue != "" {
			switch val.Field(i).Kind() {
			case reflect.Bool:
				if value, err := strconv.ParseBool(defaultValue); err == nil {
					val.Field(i).SetBool(value)
				}
			case reflect.Int:
			case reflect.Int32:
			case reflect.Int64:
				if value, err := strconv.ParseInt(defaultValue,10,64); err == nil {
					val.Field(i).SetInt(value)
				}
			case reflect.Uint:
			case reflect.Uint32:
			case reflect.Uint64:
				if value, err := strconv.ParseUint(defaultValue,10,64); err == nil {
					val.Field(i).SetUint(value)
				}
			case reflect.Float32:
			case reflect.Float64:
				if value, err := strconv.ParseFloat(defaultValue,64); err == nil {
					val.Field(i).SetFloat(value)
				}
			case reflect.String:
				val.Field(i).SetString(defaultValue)
			}
		}
	}
	return nil
}