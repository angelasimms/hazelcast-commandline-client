//go:build std || multimap

package multimap

import (
	"context"
	"fmt"

	"github.com/hazelcast/hazelcast-go-client"

	"github.com/hazelcast/hazelcast-commandline-client/clc"
	. "github.com/hazelcast/hazelcast-commandline-client/internal/check"
	"github.com/hazelcast/hazelcast-commandline-client/internal/output"
	"github.com/hazelcast/hazelcast-commandline-client/internal/plug"
	"github.com/hazelcast/hazelcast-commandline-client/internal/proto/codec"
	"github.com/hazelcast/hazelcast-commandline-client/internal/serialization"
)

type MultiMapRemoveCommand struct{}

func (mc *MultiMapRemoveCommand) Init(cc plug.InitContext) error {
	addKeyTypeFlag(cc)
	help := "Remove values from the given MultiMap"
	cc.SetCommandHelp(help, help)
	cc.SetCommandUsage("remove [key] [flags]")
	cc.SetPositionalArgCount(1, 1)
	return nil
}

func (mc *MultiMapRemoveCommand) Exec(ctx context.Context, ec plug.ExecContext) error {
	mmName := ec.Props().GetString(multiMapFlagName)
	ci, err := ec.ClientInternal(ctx)
	if err != nil {
		return err
	}
	keyStr := ec.Args()[0]
	keyData, err := makeKeyData(ec, ci, keyStr)
	if err != nil {
		return err
	}
	req := codec.EncodeMultiMapRemoveRequest(mmName, keyData, 0)
	rv, stop, err := ec.ExecuteBlocking(ctx, func(ctx context.Context, sp clc.Spinner) (any, error) {
		sp.SetText(fmt.Sprintf("Removing from multimap %s", mmName))
		return ci.InvokeOnKey(ctx, req, keyData, nil)
	})
	if err != nil {
		return err
	}
	stop()
	raw := codec.DecodeMultiMapRemoveResponse(rv.(*hazelcast.ClientMessage))
	var rows []output.Row
	for _, r := range raw {
		vt := r.Type()
		value, err := ci.DecodeData(*r)
		if err != nil {
			ec.Logger().Info("The value for %s was not decoded, due to error: %s", keyStr, err.Error())
			value = serialization.NondecodedType(serialization.TypeToLabel(vt))
		}
		row := output.Row{
			output.Column{
				Name:  output.NameValue,
				Type:  vt,
				Value: value,
			},
		}
		if ec.Props().GetBool(multiMapFlagShowType) {
			row = append(row, output.Column{
				Name:  output.NameValueType,
				Type:  serialization.TypeString,
				Value: serialization.TypeToLabel(vt),
			})
		}
		rows = append(rows, row)
	}
	return ec.AddOutputRows(ctx, rows...)
}

func init() {
	Must(plug.Registry.RegisterCommand("multi-map:remove", &MultiMapRemoveCommand{}))
}
