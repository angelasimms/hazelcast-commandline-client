//go:build std || atomiclong

package atomiclong

import (
	"context"
	"fmt"

	"github.com/hazelcast/hazelcast-go-client"

	"github.com/hazelcast/hazelcast-commandline-client/clc"
	. "github.com/hazelcast/hazelcast-commandline-client/internal/check"
	"github.com/hazelcast/hazelcast-commandline-client/internal/plug"
)

const (
	atomicLongFlagName     = "name"
	atomicLongPropertyName = "atomic-long"
)

type AtomicLongCommand struct {
}

func (mc *AtomicLongCommand) Init(cc plug.InitContext) error {
	cc.AddCommandGroup(clc.GroupDDSID, clc.GroupDDSTitle)
	cc.SetCommandGroup(clc.GroupDDSID)
	cc.AddStringFlag(atomicLongFlagName, "n", defaultAtomicLongName, false, "atomic long name")
	cc.SetTopLevel(true)
	cc.SetCommandUsage("atomic-long [command] [flags]")
	help := "Atomic long operations"
	cc.SetCommandHelp(help, help)
	return nil
}

func (mc *AtomicLongCommand) Exec(context.Context, plug.ExecContext) error {
	return nil
}

func (mc *AtomicLongCommand) Augment(ec plug.ExecContext, props *plug.Properties) error {
	ctx := context.TODO()
	props.SetBlocking(atomicLongPropertyName, func() (any, error) {
		name := ec.Props().GetString(atomicLongFlagName)
		// empty atomic long name is allowed
		ci, err := ec.ClientInternal(ctx)
		if err != nil {
			return nil, err
		}
		mv, stop, err := ec.ExecuteBlocking(ctx, func(ctx context.Context, sp clc.Spinner) (any, error) {
			sp.SetText(fmt.Sprintf("Getting atomic long %s", name))
			m, err := ci.Client().CPSubsystem().GetAtomicLong(ctx, name)
			if err != nil {
				return nil, err
			}
			return m, nil
		})
		if err != nil {
			return nil, err
		}
		stop()
		return mv.(*hazelcast.AtomicLong), nil
	})
	return nil
}

func init() {
	cmd := &AtomicLongCommand{}
	Must(plug.Registry.RegisterCommand("atomic-long", cmd))
	plug.Registry.RegisterAugmentor("20-atomic-long", cmd)
}
