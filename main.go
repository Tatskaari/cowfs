package main

import (
	"fmt"
	"github.com/peterebden/go-cli-init"
	"github.com/tatskaari/cowfs/fs"
	"os"
	"os/signal"
	"syscall"
)

var opts = struct {
	Usage     string
	Mount struct {
		Args     struct {
			MountPoint string `positional-arg-name:"mount-point" description:"The directory to mount the filesystem"`
			Files []string `positional-arg-name:"files" description:"A list of files to add to the filesystem"`
		} `positional-args:"true"`
	} `command:"mount" description:"Mount a cowfs file system in the specified directory."`
	Unmount struct {
		Args     struct {
			MountPoint string `positional-arg-name:"files" description:"The mount point to unmount"`
		} `positional-args:"true"`
	} `command:"unmount" description:"Unmount a cowfs file system from the specified directory."`

}{
	Usage: `cowfs can be used to mount a copy-on-write (cow) FUSE file-system.`,
}



func main() {
	command := cli.ParseFlagsOrDie("cowfs", &opts)
	if command == "mount" {
		go handleSignals()
		fs.MountAndServe(opts.Mount.Args.MountPoint, opts.Mount.Args.Files)
	} else {
		if err := fs.Unmount(opts.Unmount.Args.MountPoint); err != nil {
			fmt.Println("Failed to unmount fs ", err)
		}
	}
}

func handleSignals() {
	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGHUP, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGABRT, syscall.SIGTERM)
	sig := <-ch
	fmt.Printf("Recived signal %v, shutting down\n", sig)

	if err := fs.Unmount(opts.Mount.Args.MountPoint); err != nil {
		fmt.Println("Failed to unmount fs ", err)
	}
}
