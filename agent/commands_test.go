package main

import (
	"log"
	"os"
	"reflect"
	"slices"
	"strings"
	"testing"

	"howett.net/plist"
)

func TestQuotedStringSplit(t *testing.T) {
	out := quotedStringSplit("stuff and things")
	log.Printf("Output: %+v", out)
	if !slices.Equal(out, []string{"stuff", "and", "things"}) {
		t.FailNow()
	}
}

func TestDecodePlist(t *testing.T) {
	xmlFile, err := os.Open("systemprofiler.xml")
	if checkError(err) {
		return
	}

	pd := plist.NewDecoder(xmlFile)

	var this interface{}
	err = pd.Decode(&this)
	if checkError(err) {
		return
	}

	var aspo AppleSystemProfilerOutput

	for i, el := range this.([]any) {
		_ = i
		//log.Printf("Index: %#v, Element: %#v\n\n", i, el)

		items := el.(map[string]any)

		if _, ok := items["_dataType"]; ok {
			log.Printf("==================================")
			log.Printf("Data Type: %s", items["_dataType"])
			log.Printf("==================================")
		}

		//rt := reflect.TypeOf(aspo)
		//field, _ := rt.FieldByName(items["_dataType"].(string))
		rv := reflect.ValueOf(aspo)
		subfieldT := reflect.TypeOf(rv.FieldByName(items["_dataType"].(string)))
		subfield, _ := subfieldT.FieldByName(items["_name"].(string))
		subfield.Tag.Get("plist")
		sf := rv.FieldByName(items["_name"].(string))

		if _, ok := items["_name"]; ok {
			log.Printf("----------------------------------")
			log.Printf("Name: %s", items["_name"])
			log.Printf("----------------------------------")
		}

		for k, v := range items {
			//log.Printf("Key: %#v, Val: %#v", k, v)

			if k == "_items" {

				switch v.(type) {
				case string:
					log.Printf("Item string: %s", v.(string))
				case []any:
					unrollPlistSlice(rv, v.([]any), 0)
				}
			}
		}
	}

}

func unrollPlistSlice(rv reflect.Value, in []any, depth int) {
	for i, item := range in {
		_ = i
		switch item.(type) {
		case string:
			log.Printf("Item string: %#v", item.(string))
		case []any:
			unrollPlistSlice(rv, item.([]any), depth+1)
		case map[string]any:
			unrollPlistMap(rv, item.(map[string]any), depth+1)
		}
	}

}

func unrollPlistMap(rv reflect.Value, in map[string]any, depth int) {

	if _, ok := in["_name"]; ok {
		log.Printf("%s ----------------------------------", strings.Repeat(">", depth))
		log.Printf("%s Name: %s", strings.Repeat(">", depth), in["_name"])
		log.Printf("%s -----------------------------------", strings.Repeat(">", depth))
	}

	for k, v := range in {
		if _, ok := in["_items"]; !ok {
			if k != "_name" {
				switch v.(type) {
				case string:

					t := reflect.TypeOf(rv)
					f, _ := rv.FieldByName(t)

					ft := f.Tag
					ft.Get("plist")

					sf.Set(v.(reflect.Value))
					log.Printf("%s Sub%dKey: %#v, Sub%dVal: %#v", strings.Repeat(">", depth), depth, k, depth, v)
				case map[string]any:
					unrollPlistMap(rv, v.(map[string]any), depth+1)
				}
			}
		}

		if k == "_items" {
			switch v.(type) {
			case string:
				log.Printf("Item name: %s", v.(string))
			case []any:
				unrollPlistSlice(rv, v.([]any), depth+1)
			}
		}

	}
}
