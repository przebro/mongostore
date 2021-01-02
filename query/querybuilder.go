package query

import (
	"fmt"
	"strings"

	"github.com/przebro/databazaar/selector"
)

type mongoFormatter struct {
}

func (f *mongoFormatter) Format(fld, op, val string) string {
	return fmt.Sprintf(`{"%s":{"%s":%s}}`, fld, op, val)
}
func (f *mongoFormatter) FormatArray(op string, val ...string) string {

	result := strings.Join(val, ",")
	return fmt.Sprintf(`{"%s":[%s]}`, op, result)
}

type mongoQueryBuilder struct {
	formatter selector.Formatter
}

func NewBuilder() selector.DataSelectorBuilder {

	return &mongoQueryBuilder{formatter: &mongoFormatter{}}
}

func (qb mongoQueryBuilder) Build(expr selector.Expr) string {

	epxand := expr.Expand(qb.formatter)
	result := fmt.Sprintf(`%s`, epxand)
	return result
}
