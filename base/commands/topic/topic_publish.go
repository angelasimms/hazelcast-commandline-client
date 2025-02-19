//go:build std || topic

package topic

import (
	"context"
	"fmt"
	"math"

	"github.com/hazelcast/hazelcast-go-client"

	"github.com/hazelcast/hazelcast-commandline-client/clc"
	"github.com/hazelcast/hazelcast-commandline-client/internal/plug"
	"github.com/hazelcast/hazelcast-commandline-client/internal/topic"

	. "github.com/hazelcast/hazelcast-commandline-client/internal/check"
)

type topicPublishCommand struct{}

func (tpc *topicPublishCommand) Init(cc plug.InitContext) error {
	addValueTypeFlag(cc)
	help := "Publish new messages for a Topic."
	cc.SetCommandHelp(help, help)
	cc.SetPositionalArgCount(1, math.MaxInt)
	cc.SetCommandUsage("publish [values] [flags]")
	return nil
}

func (tpc *topicPublishCommand) Exec(ctx context.Context, ec plug.ExecContext) error {
	topicName := ec.Props().GetString(topicFlagName)
	// get the topic just to ensure the corresponding proxy is created
	_, err := ec.Props().GetBlocking(topicPropertyName)
	if err != nil {
		return err
	}
	ci, err := ec.ClientInternal(ctx)
	if err != nil {
		return err
	}

	_, stop, err := ec.ExecuteBlocking(ctx, func(ctx context.Context, sp clc.Spinner) (any, error) {
		sp.SetText(fmt.Sprintf("Publishing values into topic %s", topicName))
		var vals []hazelcast.Data
		for _, valStr := range ec.Args() {
			val, err := makeValueData(ec, ci, valStr)
			if err != nil {
				return nil, err
			}
			vals = append(vals, val)
		}
		return nil, topic.PublishAll(ctx, ci, topicName, vals)
	})
	if err != nil {
		return err
	}
	stop()
	return nil
}

func init() {
	Must(plug.Registry.RegisterCommand("topic:publish", &topicPublishCommand{}))
}
