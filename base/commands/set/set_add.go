//go:build std || set

package set

import (
	"context"
	"fmt"
	"math"

	"github.com/hazelcast/hazelcast-go-client"

	"github.com/hazelcast/hazelcast-commandline-client/clc"
	. "github.com/hazelcast/hazelcast-commandline-client/internal/check"
	"github.com/hazelcast/hazelcast-commandline-client/internal/plug"
)

type SetAddCommand struct{}

func (sc *SetAddCommand) Init(cc plug.InitContext) error {
	addValueTypeFlag(cc)
	cc.SetPositionalArgCount(1, math.MaxInt)
	help := "Add values to the given Set"
	cc.SetCommandHelp(help, help)
	cc.SetCommandUsage("add [values] [flags]")
	return nil
}

func (sc *SetAddCommand) Exec(ctx context.Context, ec plug.ExecContext) error {
	name := ec.Props().GetString(setFlagName)
	sv, err := ec.Props().GetBlocking(setPropertyName)
	if err != nil {
		return err
	}
	ci, err := ec.ClientInternal(ctx)
	if err != nil {
		return err
	}
	s := sv.(*hazelcast.Set)
	_, stop, err := ec.ExecuteBlocking(ctx, func(ctx context.Context, sp clc.Spinner) (any, error) {
		sp.SetText(fmt.Sprintf("Adding values into set %s", name))
		for _, arg := range ec.Args() {
			vd, err := makeValueData(ec, ci, arg)
			if err != nil {
				return nil, err
			}
			_, err = s.Add(ctx, vd)
			if err != nil {
				return nil, err
			}
		}
		return nil, nil
	})
	if err != nil {
		return err
	}
	stop()
	return nil
}

func init() {
	Must(plug.Registry.RegisterCommand("set:add", &SetAddCommand{}))
}
