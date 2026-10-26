package validatex

import (
	"encoding/json"
	"net"
	"net/url"
	"regexp"
	"slices"
	"strconv"
	"strings"
	"time"
	"unicode"
	"unicode/utf8"
)

// Filled 校验字符串非空
func (f *StringField) Filled() *StringField {
	f.check("filled", strings.TrimSpace(f.value) == "")
	return f
}

// RequiredIf 条件成立时校验字符串必填
func (f *StringField) RequiredIf(condition bool) *StringField {
	f.check("required_if", condition && strings.TrimSpace(f.value) == "")
	return f
}

// RequiredUnless 条件不成立时校验字符串必填
func (f *StringField) RequiredUnless(condition bool) *StringField {
	f.check("required_unless", !condition && strings.TrimSpace(f.value) == "")
	return f
}

// Prohibited 校验字符串为空
func (f *StringField) Prohibited() *StringField {
	f.check("prohibited", strings.TrimSpace(f.value) != "")
	return f
}

// ProhibitedIf 条件成立时校验字符串为空
func (f *StringField) ProhibitedIf(condition bool) *StringField {
	f.check("prohibited_if", condition && strings.TrimSpace(f.value) != "")
	return f
}

// ProhibitedUnless 条件不成立时校验字符串为空
func (f *StringField) ProhibitedUnless(condition bool) *StringField {
	f.check("prohibited_unless", !condition && strings.TrimSpace(f.value) != "")
	return f
}

// Max 校验字符串最大字符数
func (f *StringField) Max(n int) *StringField {
	f.check("max", f.value != "" && n >= 0 && utf8.RuneCountInString(f.value) > n, n)
	return f
}

// Between 校验字符串字符数范围
func (f *StringField) Between(min, max int) *StringField {
	length := utf8.RuneCountInString(f.value)
	if f.value != "" && (length < min || length > max) {
		f.check("between", true, arg("min", min), arg("max", max))
	}
	return f
}

// Size 校验字符串字符数等于指定值
func (f *StringField) Size(n int) *StringField {
	length := utf8.RuneCountInString(f.value)
	f.check("size", f.value != "" && length != n, n)
	return f
}

// Alpha 校验字符串仅包含字母
func (f *StringField) Alpha() *StringField {
	f.check("alpha", f.value != "" && !matchRunes(f.value, func(r rune) bool {
		return unicode.IsLetter(r) || unicode.IsMark(r)
	}))
	return f
}

// AlphaDash 校验字符串仅包含字母、数字、短横线和下划线
func (f *StringField) AlphaDash() *StringField {
	f.check("alpha_dash", f.value != "" && !matchRunes(f.value, func(r rune) bool {
		return unicode.IsLetter(r) || unicode.IsMark(r) || unicode.IsDigit(r) || r == '-' || r == '_'
	}))
	return f
}

// AlphaNum 校验字符串仅包含字母和数字
func (f *StringField) AlphaNum() *StringField {
	f.check("alpha_num", f.value != "" && !matchRunes(f.value, func(r rune) bool {
		return unicode.IsLetter(r) || unicode.IsMark(r) || unicode.IsDigit(r)
	}))
	return f
}

// ASCII 校验字符串仅包含ASCII字符
func (f *StringField) ASCII() *StringField {
	valid := true
	for _, r := range f.value {
		if r > unicode.MaxASCII {
			valid = false
			break
		}
	}
	f.check("ascii", f.value != "" && !valid)
	return f
}

// Lowercase 校验字符串为小写
func (f *StringField) Lowercase() *StringField {
	f.check("lowercase", hasCaseDifference(f.value, unicode.ToLower))
	return f
}

// Uppercase 校验字符串为大写
func (f *StringField) Uppercase() *StringField {
	f.check("uppercase", hasCaseDifference(f.value, unicode.ToUpper))
	return f
}

// hasCaseDifference 判断字符串是否包含大小写变化或非法UTF-8
func hasCaseDifference(value string, convert func(rune) rune) bool {
	for index := 0; index < len(value); {
		character, size := utf8.DecodeRuneInString(value[index:])
		if character == utf8.RuneError && size == 1 {
			return true
		}
		if convert(character) != character {
			return true
		}
		index += size
	}
	return false
}

// StartsWith 校验字符串以任一目标开头
func (f *StringField) StartsWith(values ...string) *StringField {
	valid := false
	for _, value := range values {
		if strings.HasPrefix(f.value, value) {
			valid = true
			break
		}
	}
	if f.value != "" && !valid {
		f.check("starts_with", true, arg("values", strings.Join(values, ",")))
	}
	return f
}

// DoesntStartWith 校验字符串不以任何目标开头
func (f *StringField) DoesntStartWith(values ...string) *StringField {
	invalid := false
	for _, value := range values {
		if strings.HasPrefix(f.value, value) {
			invalid = true
			break
		}
	}
	if f.value != "" && invalid {
		f.check("doesnt_start_with", true, arg("values", strings.Join(values, ",")))
	}
	return f
}

// EndsWith 校验字符串以任一目标结尾
func (f *StringField) EndsWith(values ...string) *StringField {
	valid := false
	for _, value := range values {
		if strings.HasSuffix(f.value, value) {
			valid = true
			break
		}
	}
	if f.value != "" && !valid {
		f.check("ends_with", true, arg("values", strings.Join(values, ",")))
	}
	return f
}

// DoesntEndWith 校验字符串不以任何目标结尾
func (f *StringField) DoesntEndWith(values ...string) *StringField {
	invalid := false
	for _, value := range values {
		if strings.HasSuffix(f.value, value) {
			invalid = true
			break
		}
	}
	if f.value != "" && invalid {
		f.check("doesnt_end_with", true, arg("values", strings.Join(values, ",")))
	}
	return f
}

// Contains 校验字符串包含任一目标
func (f *StringField) Contains(values ...string) *StringField {
	valid := false
	for _, value := range values {
		if strings.Contains(f.value, value) {
			valid = true
			break
		}
	}
	if f.value != "" && !valid {
		f.check("contains", true, arg("values", strings.Join(values, ",")))
	}
	return f
}

// DoesntContain 校验字符串不包含任何目标
func (f *StringField) DoesntContain(values ...string) *StringField {
	invalid := false
	for _, value := range values {
		if strings.Contains(f.value, value) {
			invalid = true
			break
		}
	}
	if f.value != "" && invalid {
		f.check("doesnt_contain", true, arg("values", strings.Join(values, ",")))
	}
	return f
}

// In 校验字符串属于指定集合
func (f *StringField) In(values ...string) *StringField {
	if f.value != "" && !slices.Contains(values, f.value) {
		f.check("in", true, arg("values", strings.Join(values, ",")))
	}
	return f
}

// NotIn 校验字符串不属于指定集合
func (f *StringField) NotIn(values ...string) *StringField {
	if f.value != "" && slices.Contains(values, f.value) {
		f.check("not_in", true, arg("values", strings.Join(values, ",")))
	}
	return f
}

// Regex 校验字符串匹配正则
func (f *StringField) Regex(pattern string) *StringField {
	if f.value == "" {
		return f
	}
	matched, err := regexp.MatchString(pattern, f.value)
	if err != nil || !matched {
		f.check("regex", true, arg("pattern", pattern))
	}
	return f
}

// Regexp 使用预编译正则校验字符串
func (f *StringField) Regexp(pattern *regexp.Regexp) *StringField {
	if f.value == "" {
		return f
	}
	f.check("regex", pattern == nil || !pattern.MatchString(f.value))
	return f
}

// NotRegex 校验字符串不匹配正则
func (f *StringField) NotRegex(pattern string) *StringField {
	if f.value == "" {
		return f
	}
	matched, err := regexp.MatchString(pattern, f.value)
	if err != nil || matched {
		f.check("not_regex", true, arg("pattern", pattern))
	}
	return f
}

// NotRegexp 使用预编译正则校验字符串
func (f *StringField) NotRegexp(pattern *regexp.Regexp) *StringField {
	if f.value == "" {
		return f
	}
	f.check("not_regex", pattern == nil || pattern.MatchString(f.value))
	return f
}

// Numeric 校验字符串为数字
func (f *StringField) Numeric() *StringField {
	f.check("numeric", f.value != "" && !isNumericString(f.value))
	return f
}

// Integer 校验字符串为整数
func (f *StringField) Integer() *StringField {
	if f.value == "" {
		return f
	}
	_, err := strconv.ParseInt(f.value, 10, 64)
	f.check("integer", err != nil)
	return f
}

// Decimal 校验字符串小数位数
func (f *StringField) Decimal(places ...int) *StringField {
	decimalPlaces, valid := decimalPlaces(f.value)
	if valid && len(places) > 0 {
		valid = decimalPlaces >= places[0]
	}
	if valid && len(places) > 1 {
		valid = decimalPlaces <= places[1]
	}
	if f.value != "" && !valid {
		f.check("decimal", true,
			arg("min", placesValue(places, 0)),
			arg("max", placesValue(places, 1)),
		)
	}
	return f
}

// decimalPlaces 返回合法十进制数的小数位数
func decimalPlaces(value string) (int, bool) {
	if value == "" {
		return 0, false
	}

	index := 0
	if value[index] == '+' || value[index] == '-' {
		index++
	}
	integerStart := index
	for index < len(value) && value[index] >= '0' && value[index] <= '9' {
		index++
	}
	if index == integerStart {
		return 0, false
	}
	if index == len(value) {
		return 0, true
	}
	if value[index] != '.' {
		return 0, false
	}

	index++
	fractionStart := index
	for index < len(value) && value[index] >= '0' && value[index] <= '9' {
		index++
	}
	if index != len(value) || index == fractionStart {
		return 0, false
	}
	return index - fractionStart, true
}

// Digits 校验数字字符串长度
func (f *StringField) Digits(n int) *StringField {
	valid := isDigits(f.value) && len(f.value) == n
	if f.value != "" && !valid {
		f.check("digits", true, arg("digits", n))
	}
	return f
}

// DigitsBetween 校验数字字符串长度范围
func (f *StringField) DigitsBetween(min, max int) *StringField {
	length := len(f.value)
	valid := isDigits(f.value) && length >= min && length <= max
	if f.value != "" && !valid {
		f.check("digits_between", true, arg("min", min), arg("max", max))
	}
	return f
}

// URL 校验URL格式
func (f *StringField) URL() *StringField {
	if f.value == "" {
		return f
	}
	parsed, err := url.ParseRequestURI(f.value)
	f.check("url", err != nil || parsed.Scheme == "" || parsed.Host == "")
	return f
}

// IP 校验IP地址格式
func (f *StringField) IP() *StringField {
	if f.value == "" {
		return f
	}
	f.check("ip", net.ParseIP(f.value) == nil)
	return f
}

// IPv4 校验IPv4地址格式
func (f *StringField) IPv4() *StringField {
	if f.value == "" {
		return f
	}
	ip := net.ParseIP(f.value)
	f.check("ipv4", ip == nil || ip.To4() == nil)
	return f
}

// IPv6 校验IPv6地址格式
func (f *StringField) IPv6() *StringField {
	if f.value == "" {
		return f
	}
	ip := net.ParseIP(f.value)
	f.check("ipv6", ip == nil || ip.To4() != nil)
	return f
}

// MACAddress 校验MAC地址格式
func (f *StringField) MACAddress() *StringField {
	if f.value == "" {
		return f
	}
	_, err := net.ParseMAC(f.value)
	f.check("mac_address", err != nil)
	return f
}

// UUID 校验UUID格式
func (f *StringField) UUID() *StringField {
	if f.value == "" {
		return f
	}
	matched := uuidPattern.MatchString(f.value)
	f.check("uuid", !matched)
	return f
}

// ULID 校验ULID格式
func (f *StringField) ULID() *StringField {
	if f.value == "" {
		return f
	}
	matched := ulidPattern.MatchString(f.value)
	f.check("ulid", !matched)
	return f
}

// HexColor 校验十六进制颜色格式
func (f *StringField) HexColor() *StringField {
	if f.value == "" {
		return f
	}
	matched := hexColorPattern.MatchString(f.value)
	f.check("hex_color", !matched)
	return f
}

// Timezone 校验时区名称
func (f *StringField) Timezone() *StringField {
	if f.value == "" {
		return f
	}
	_, err := time.LoadLocation(f.value)
	f.check("timezone", err != nil)
	return f
}

// JSON 校验JSON格式
func (f *StringField) JSON() *StringField {
	if f.value == "" {
		return f
	}
	f.check("json", !json.Valid([]byte(f.value)))
	return f
}

// Accepted 校验字符串为接受值
func (f *StringField) Accepted() *StringField {
	f.check("accepted", !isAccepted(f.value))
	return f
}

// AcceptedIf 条件成立时校验字符串为接受值
func (f *StringField) AcceptedIf(condition bool) *StringField {
	f.check("accepted_if", condition && !isAccepted(f.value))
	return f
}

// Declined 校验字符串为拒绝值
func (f *StringField) Declined() *StringField {
	f.check("declined", !isDeclined(f.value))
	return f
}

// DeclinedIf 条件成立时校验字符串为拒绝值
func (f *StringField) DeclinedIf(condition bool) *StringField {
	f.check("declined_if", condition && !isDeclined(f.value))
	return f
}

// Boolean 校验字符串为布尔值
func (f *StringField) Boolean() *StringField {
	f.check("boolean", !isAccepted(f.value) && !isDeclined(f.value))
	return f
}

// Same 校验字符串与目标值相同
func (f *StringField) Same(value string) *StringField {
	f.check("same", f.value != value, value)
	return f
}

// Different 校验字符串与目标值不同
func (f *StringField) Different(value string) *StringField {
	f.check("different", f.value == value, value)
	return f
}

// Confirmed 校验确认值一致
func (f *StringField) Confirmed(value string) *StringField {
	f.check("confirmed", f.value != value)
	return f
}

// Date 校验日期格式
func (f *StringField) Date() *StringField {
	if f.value == "" {
		return f
	}
	_, valid := parseDate(f.value)
	f.check("date", !valid)
	return f
}

// DateFormat 按指定布局校验日期
func (f *StringField) DateFormat(layout string) *StringField {
	if f.value == "" {
		return f
	}
	_, err := time.Parse(layout, f.value)
	if err != nil {
		f.check("date_format", true, arg("format", layout))
	}
	return f
}

// DateEquals 校验日期与目标日期相同
func (f *StringField) DateEquals(value string) *StringField {
	if f.value == "" {
		return f
	}
	date, valid := parseDate(f.value)
	target, targetValid := parseDate(value)
	invalid := !valid || !targetValid || !sameDate(date, target)
	f.check("date_equals", invalid, value)
	return f
}

// Before 校验日期早于目标日期
func (f *StringField) Before(value string) *StringField {
	if f.value == "" {
		return f
	}
	f.check("before", !compareDate(f.value, value, func(date, target time.Time) bool {
		return date.Before(target)
	}), value)
	return f
}

// BeforeOrEqual 校验日期不晚于目标日期
func (f *StringField) BeforeOrEqual(value string) *StringField {
	if f.value == "" {
		return f
	}
	f.check("before_or_equal", !compareDate(f.value, value, func(date, target time.Time) bool {
		return date.Before(target) || date.Equal(target)
	}), value)
	return f
}

// After 校验日期晚于目标日期
func (f *StringField) After(value string) *StringField {
	if f.value == "" {
		return f
	}
	f.check("after", !compareDate(f.value, value, func(date, target time.Time) bool {
		return date.After(target)
	}), value)
	return f
}

// AfterOrEqual 校验日期不早于目标日期
func (f *StringField) AfterOrEqual(value string) *StringField {
	if f.value == "" {
		return f
	}
	f.check("after_or_equal", !compareDate(f.value, value, func(date, target time.Time) bool {
		return date.After(target) || date.Equal(target)
	}), value)
	return f
}

// matchRunes 校验所有字符是否满足条件
func matchRunes(value string, valid func(rune) bool) bool {
	for _, r := range value {
		if !valid(r) {
			return false
		}
	}
	return true
}

// isDigits 判断字符串是否仅包含数字
func isDigits(value string) bool {
	if value == "" {
		return false
	}
	for index := 0; index < len(value); index++ {
		if value[index] < '0' || value[index] > '9' {
			return false
		}
	}
	return true
}

// isNumericString 判断字符串是否为十进制数字,支持小数和科学计数法
func isNumericString(value string) bool {
	if value == "" {
		return false
	}

	index := 0
	if value[index] == '+' || value[index] == '-' {
		index++
		if index == len(value) {
			return false
		}
	}

	digits := 0
	for index < len(value) && value[index] >= '0' && value[index] <= '9' {
		index++
		digits++
	}
	if index < len(value) && value[index] == '.' {
		index++
		for index < len(value) && value[index] >= '0' && value[index] <= '9' {
			index++
			digits++
		}
	}
	if digits == 0 {
		return false
	}

	if index < len(value) && (value[index] == 'e' || value[index] == 'E') {
		index++
		if index < len(value) && (value[index] == '+' || value[index] == '-') {
			index++
		}
		exponentDigits := 0
		for index < len(value) && value[index] >= '0' && value[index] <= '9' {
			index++
			exponentDigits++
		}
		if exponentDigits == 0 {
			return false
		}
	}

	return index == len(value)
}

// isEmail 判断字符串是否为常见 ASCII 邮箱格式
func isEmail(value string) bool {
	if value == "" || len(value) > 254 {
		return false
	}

	at := strings.LastIndexByte(value, '@')
	if at <= 0 || at == len(value)-1 {
		return false
	}
	local, domain := value[:at], value[at+1:]
	if len(local) > 64 || local[0] == '.' || local[len(local)-1] == '.' || strings.Contains(local, "..") {
		return false
	}
	for index := 0; index < len(local); index++ {
		character := local[index]
		if character >= 'a' && character <= 'z' ||
			character >= 'A' && character <= 'Z' ||
			character >= '0' && character <= '9' ||
			character == '.' || character == '_' || character == '%' || character == '+' || character == '-' {
			continue
		}
		return false
	}

	labelStart := 0
	labels := 0
	lastDot := -1
	for index := 0; index < len(domain); index++ {
		character := domain[index]
		if character == '.' {
			labelLength := index - labelStart
			if labelLength == 0 || labelLength > 63 || domain[labelStart] == '-' || domain[index-1] == '-' {
				return false
			}
			labels++
			lastDot = index
			labelStart = index + 1
			continue
		}
		if character >= 'a' && character <= 'z' ||
			character >= 'A' && character <= 'Z' ||
			character >= '0' && character <= '9' ||
			character == '-' {
			continue
		}
		return false
	}

	labelLength := len(domain) - labelStart
	if labelLength == 0 || labelLength > 63 || domain[labelStart] == '-' || domain[len(domain)-1] == '-' {
		return false
	}
	labels++
	if labels < 2 {
		return false
	}

	topLevel := domain[lastDot+1:]
	if len(topLevel) < 2 {
		return false
	}
	for index := 0; index < len(topLevel); index++ {
		character := topLevel[index]
		if character < 'a' || character > 'z' {
			if character < 'A' || character > 'Z' {
				return false
			}
		}
	}
	return true
}

// isAccepted 判断字符串是否为接受值
func isAccepted(value string) bool {
	value = strings.TrimSpace(value)
	return strings.EqualFold(value, "yes") ||
		strings.EqualFold(value, "on") ||
		value == "1" ||
		strings.EqualFold(value, "true")
}

// isDeclined 判断字符串是否为拒绝值
func isDeclined(value string) bool {
	value = strings.TrimSpace(value)
	return strings.EqualFold(value, "no") ||
		strings.EqualFold(value, "off") ||
		value == "0" ||
		strings.EqualFold(value, "false")
}

// parseDate 解析常用日期格式
func parseDate(value string) (time.Time, bool) {
	for _, layout := range dateLayouts {
		if date, err := time.Parse(layout, value); err == nil {
			return date, true
		}
	}
	return time.Time{}, false
}

// compareDate 比较两个日期字符串
func compareDate(value, target string, compare func(time.Time, time.Time) bool) bool {
	date, valid := parseDate(value)
	reference, targetValid := parseDate(target)
	return valid && targetValid && compare(date, reference)
}

// sameDate 判断两个日期是否为同一天
func sameDate(left, right time.Time) bool {
	leftYear, leftMonth, leftDay := left.Date()
	rightYear, rightMonth, rightDay := right.Date()
	return leftYear == rightYear && leftMonth == rightMonth && leftDay == rightDay
}

// placesValue 获取指定位置的小数位数
func placesValue(places []int, index int) any {
	if index >= len(places) {
		return ""
	}
	return places[index]
}
