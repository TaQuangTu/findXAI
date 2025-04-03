package server

import (
	"errors"
	"findx/pkg/protogen"
	"fmt"
	"reflect"
	"regexp"
	"strings"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
)

var (
	validFilterValues       = map[string]bool{"0": true, "1": true}
	validSafeValues         = map[string]bool{"active": true, "off": true}
	validSearchTypeValues   = map[string]bool{"image": true}
	validSiteSearchFilter   = map[string]bool{"e": true, "i": true}
	validImgColorTypeValues = map[string]bool{"color": true, "gray": true, "mono": true, "trans": true}
	validImgDominantColors  = map[string]bool{"black": true, "blue": true, "brown": true, "gray": true, "green": true, "orange": true, "pink": true, "purple": true, "red": true, "teal": true, "white": true, "yellow": true}
	validImgSizeValues      = map[string]bool{"huge": true, "icon": true, "large": true, "medium": true, "small": true, "xlarge": true, "xxlarge": true}
	validImgTypeValues      = map[string]bool{"clipart": true, "face": true, "lineart": true, "stock": true, "photo": true, "animated": true}
	validLrValues           = map[string]bool{
		"lang_ar": true, "lang_bg": true, "lang_ca": true, "lang_cs": true, "lang_da": true, "lang_de": true,
		"lang_el": true, "lang_en": true, "lang_es": true, "lang_et": true, "lang_fi": true, "lang_fr": true,
		"lang_hr": true, "lang_hu": true, "lang_id": true, "lang_is": true, "lang_it": true, "lang_iw": true,
		"lang_ja": true, "lang_ko": true, "lang_lt": true, "lang_lv": true, "lang_nl": true, "lang_no": true,
		"lang_pl": true, "lang_pt": true, "lang_ro": true, "lang_ru": true, "lang_sk": true, "lang_sl": true,
		"lang_sr": true, "lang_sv": true, "lang_tr": true, "lang_zh-CN": true, "lang_zh-TW": true,
	}
	validRightsValues = map[string]bool{"cc_publicdomain": true, "cc_attribute": true, "cc_sharealike": true, "cc_noncommercial": true, "cc_nonderived": true}
	validSearchType   = map[string]bool{"image": true}
)

func ValidateSearchRequest(req *protogen.SearchRequest) error {
	if req.SearchType != "" && !validSearchType[req.SearchType] {
		return errors.New("invalid search type value, can be null or image")
	}
	if req.C2Coff != "" && req.C2Coff != "0" && req.C2Coff != "1" {
		return errors.New("invalid c2coff value")
	}
	if req.Filter != "" && !validFilterValues[req.Filter] {
		return errors.New("invalid filter value")
	}
	if req.Safe != "" && !validSafeValues[req.Safe] {
		return errors.New("invalid safe value")
	}
	if req.SearchType != "" && !validSearchTypeValues[req.SearchType] {
		return errors.New("invalid searchType value")
	}
	if req.SiteSearchFilter != "" && !validSiteSearchFilter[req.SiteSearchFilter] {
		return errors.New("invalid siteSearchFilter value")
	}
	if req.ImgColorType != "" && !validImgColorTypeValues[req.ImgColorType] {
		return errors.New("invalid imgColorType value")
	}
	if req.ImgDominantColor != "" && !validImgDominantColors[req.ImgDominantColor] {
		return errors.New("invalid imgDominantColor value")
	}
	if req.ImgSize != "" && !validImgSizeValues[req.ImgSize] {
		return errors.New("invalid imgSize value")
	}
	if req.ImgType != "" && !validImgTypeValues[req.ImgType] {
		return errors.New("invalid imgType value")
	}
	if req.Lr != "" && !validLrValues[req.Lr] {
		return errors.New("invalid lr value")
	}
	if req.Num < 1 || req.Num > 10 {
		return errors.New("num must be between 1 and 10")
	}
	if req.Start < 0 || req.Start+req.Num > 100 {
		return errors.New("start index must not exceed 100 in combination with num")
	}
	if req.DateRestrict != "" {
		matched, _ := regexp.MatchString(`^[dwmy][0-9]+$`, req.DateRestrict)
		if !matched {
			return errors.New("invalid dateRestrict format")
		}
	}
	if req.Rights != "" && !validRightsValues[req.Rights] {
		return errors.New("invalid 'rights' value")
	}
	if req.Gl != "" && len(req.Gl) != 2 {
		return errors.New("invalid gl value, should be two lowercase letters")
	}
	return nil
}

// ProtoMessageToMap converts a protobuf message to map[string]string
func ProtoMessageToMap(msg proto.Message) map[string]string {
	result := make(map[string]string)

	m := msg.ProtoReflect()
	m.Range(func(fd protoreflect.FieldDescriptor, v protoreflect.Value) bool {
		// Skip empty values
		if reflect.ValueOf(v.Interface()).IsZero() {
			return true
		}

		// Get the field name (in snake_case)
		fieldName := string(fd.Name())

		// Convert to camelCase for Google API params
		fieldName = toCamelCase(fieldName)

		// Convert the value to string based on its type
		var strValue string
		switch fd.Kind() {
		case protoreflect.StringKind:
			strValue = v.String()
		case protoreflect.Int32Kind, protoreflect.Int64Kind,
			protoreflect.Uint32Kind, protoreflect.Uint64Kind:
			strValue = fmt.Sprintf("%d", v.Interface())
		case protoreflect.FloatKind, protoreflect.DoubleKind:
			strValue = fmt.Sprintf("%f", v.Interface())
		case protoreflect.BoolKind:
			strValue = fmt.Sprintf("%t", v.Interface())
		default:
			// Skip complex types or handle them specifically if needed
			return true
		}

		result[fieldName] = strValue
		return true
	})

	return result
}

// toCamelCase converts snake_case to camelCase
func toCamelCase(s string) string {
	// Handle special cases for Google API params
	switch s {
	case "c2coff":
		return "c2coff"
	case "cr":
		return "cr"
	case "cx":
		return "cx"
	case "gl":
		return "gl"
	case "hl":
		return "hl"
	case "hq":
		return "hq"
	case "lr":
		return "lr"
	case "num":
		return "num"
	case "q":
		return "q"
	}

	// General snake_case to camelCase conversion
	parts := strings.Split(s, "_")
	for i := 1; i < len(parts); i++ {
		if len(parts[i]) > 0 {
			parts[i] = strings.ToUpper(parts[i][:1]) + parts[i][1:]
		}
	}
	return strings.Join(parts, "")
}
