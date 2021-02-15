package cmd

import (
	"fmt"
	"strings"

	"disco.bot/src/framework"
)

func CheersCommand(ctx framework.Context) {
	if len(ctx.Args) < 2 {
		ctx.Reply("Add command usage: `music add <song>`\nValid inputs: `youtube url`, `soundcloud url`, " +
			"`youtube id`, `soundcloud id`")
		return
	}
	fmt.Println(ctx.Guild.ID)
	sess := ctx.Sessions.GetByGuild(ctx.Guild.ID)
	if sess == nil {
		ctx.Reply("Not in a voice channel! To make the bot join one, use `music join`.")
		return
	}
	msg := ctx.Reply("Adding songs to queue...")
	arg := ctx.Args[0]
	cheers := strings.Join(ctx.Args[1:], " ")
	t, inp, err := ctx.Youtube.Get(arg)
	check, res := ctx.Db.GetKeywordValue(ctx.Guild.ID, arg)
	if check {
		t, inp, err = ctx.Youtube.Get(res)
	}
	if err != nil {
		ctx.Reply("Please use keyword or link to youtube!")
		fmt.Println(err)
		fmt.Println("error getting input, %s", err)
		return
	}

	switch t {
	case framework.ERROR_TYPE:
		ctx.Reply("An error occured!")
		fmt.Println("error type", t)
		return
	case framework.VIDEO_TYPE:
		{
			video, err := ctx.Youtube.Video(*inp)
			if err != nil {
				ctx.Reply("An error occured!")
				fmt.Println("error getting video1,", err)
				return
			}
			song := framework.NewSongCheers(video.Media, video.Title, arg, cheers)
			sess.Queue.Add(*song)
			ctx.Discord.ChannelMessageEdit(cheers+"\n"+ctx.TextChannel.ID, msg.ID, "Added `"+song.Title+"` to the song queue."+
				" Use `music play` to start playing the songs! To see the song queue, use `music queue`.")
			break
		}
	case framework.PLAYLIST_TYPE:
		{
			videos, err := ctx.Youtube.Playlist(*inp)
			if err != nil {
				ctx.Reply("An error occured!")
				fmt.Println("error getting playlist,", err)
				return
			}
			for _, v := range *videos {
				id := v.Id
				_, i, err := ctx.Youtube.Get(id)
				if err != nil {
					ctx.Reply("An error occured!")
					fmt.Println("error getting video2,", err)
					continue
				}
				video, err := ctx.Youtube.Video(*i)
				if err != nil {
					ctx.Reply("An error occured!")
					fmt.Println("error getting video3,", err)
					return
				}
				song := framework.NewSong(video.Media, video.Title, arg)
				sess.Queue.Add(*song)
			}
			ctx.Reply("Finished adding songs to the playlist. Use `music play` to start playing the songs! " +
				"To see the song queue, use `music queue`.")
			break
		}
	}

}
