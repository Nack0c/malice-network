package file

import (
	"os"

	"github.com/chainreactors/grumble"
	"github.com/chainreactors/malice-network/client/console"
	"github.com/chainreactors/malice-network/proto/client/clientpb"
)

func sync(ctx *grumble.Context, con *console.Console) {
	tid := ctx.Flags.String("taskID")
	sid := con.GetInteractive().SessionId
	syncTask, err := con.Rpc.Sync(con.ActiveTarget.Context(), &clientpb.Sync{
		FileId: sid + "-" + tid,
	})
	if err != nil {
		console.Log.Errorf("Can't sync file: %s", err)
		return
	}
	file, err := os.OpenFile(syncTask.Name, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		console.Log.Errorf("Can't Open file: %s", err)
		return
	}
	defer file.Close()
	_, err = file.Write(syncTask.Content)
	if err != nil {
		con.SessionLog(sid).Errorf("Can't write file: %s", err)
		return
	}
}
