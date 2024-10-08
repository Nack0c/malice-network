package sys

import (
	"github.com/chainreactors/grumble"
	"github.com/chainreactors/malice-network/client/console"
	"github.com/chainreactors/malice-network/helper/consts"
	"github.com/chainreactors/malice-network/proto/implant/implantpb"
	"google.golang.org/protobuf/proto"
)

func EnvCmd(ctx *grumble.Context, con *console.Console) {
	sid := con.GetInteractive().SessionId
	envTask, err := con.Rpc.Env(con.ActiveTarget.Context(), &implantpb.Request{
		Name: consts.ModuleEnv,
	})
	if err != nil {
		console.Log.Errorf("Env error: %v", err)
		return
	}
	con.AddCallback(envTask.TaskId, func(msg proto.Message) {
		env := msg.(*implantpb.Spite).GetResponse().GetKv()
		for k, v := range env {
			con.SessionLog(sid).Consolef("export %s = %s\n", k, v)
		}
	})
}

func SetEnvCmd(ctx *grumble.Context, con *console.Console) {
	sid := con.GetInteractive().SessionId
	env := ctx.Flags.String("env")
	value := ctx.Flags.String("value")
	args := []string{env, value}
	setEnvTask, err := con.Rpc.SetEnv(con.ActiveTarget.Context(), &implantpb.Request{
		Name: consts.ModuleSetEnv,
		Args: args,
	})
	if err != nil {
		console.Log.Errorf("SetEnv error: %v", err)
		return
	}
	con.AddCallback(setEnvTask.TaskId, func(msg proto.Message) {
		con.SessionLog(sid).Consolef("Set environment variable success\n")
	})
}

func UnsetEnvCmd(ctx *grumble.Context, con *console.Console) {
	sid := con.GetInteractive().SessionId
	env := ctx.Flags.String("env")
	unsetEnvTask, err := con.Rpc.UnsetEnv(con.ActiveTarget.Context(), &implantpb.Request{
		Name:  consts.ModuleUnsetEnv,
		Input: env,
	})
	if err != nil {
		console.Log.Errorf("UnsetEnv error: %v", err)
		return
	}
	con.AddCallback(unsetEnvTask.TaskId, func(msg proto.Message) {
		con.SessionLog(sid).Consolef("Unset environment variable success\n")
	})
}
