package cmd

import (
	"fmt"

	"disco.bot/src/framework"
)

func AddKeywordCommand(ctx framework.Context) {

	if len(ctx.Args) < 3 {
		ctx.Reply("Add command usage: `keyword add <keyword> <song>`\nValid inputs: `youtube url`, `soundcloud url`, " +
			"`youtube id`, `soundcloud id`")
		return
	}
	fmt.Println(ctx.Guild.ID)
	sess := ctx.Sessions.GetByGuild(ctx.Guild.ID)
	if sess == nil {
		ctx.Reply("Not in a voice channel! To make the bot join one, use `music join`.")
		return
	}
	fmt.Println(ctx.Args)
	err := ctx.Db.Insert(ctx.Guild.ID, ctx.Args[1], ctx.Args[2])
	if !err {
		ctx.Reply("An error occured!")
		fmt.Println(err)
		fmt.Println("error getting input, %s", err)
		return
	}
}
