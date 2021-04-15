package expression

import (
	cesql "github.com/cloudevents/sdk-go/sql/v2"
	"github.com/cloudevents/sdk-go/sql/v2/runtime"
	cloudevents "github.com/cloudevents/sdk-go/v2"
)

type negateExpression struct {
	child cesql.Expression
}

func (l negateExpression) Evaluate(event cloudevents.Event) (interface{}, error) {
	val, err := l.child.Evaluate(event)
	if err != nil {
		return nil, err
	}

	val, err = runtime.Cast(val, runtime.IntegerType)
	if err != nil {
		return nil, err
	}

	return -(val.(int32)), nil
}

func NewNegateExpression(child cesql.Expression) cesql.Expression {
	return negateExpression{child: child}
}
