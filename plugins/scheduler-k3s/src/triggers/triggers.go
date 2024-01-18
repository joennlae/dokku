package main

import (
	"fmt"
	"os"
	"strings"

	flag "github.com/spf13/pflag"
	"k8s.io/utils/ptr"

	"github.com/dokku/dokku/plugins/common"
	scheduler_k3s "github.com/dokku/dokku/plugins/scheduler-k3s"
)

// main entrypoint to all triggers
func main() {
	parts := strings.Split(os.Args[0], "/")
	trigger := parts[len(parts)-1]
	podIdentifier := flag.String("container-id", "", "--container-id: A pod identifier")
	flag.Parse()

	var err error
	switch trigger {
	case "install":
		err = scheduler_k3s.TriggerInstall()
	case "post-app-clone-setup":
		oldAppName := flag.Arg(0)
		newAppName := flag.Arg(1)
		err = scheduler_k3s.TriggerPostAppCloneSetup(oldAppName, newAppName)
	case "post-app-rename-setup":
		oldAppName := flag.Arg(0)
		newAppName := flag.Arg(1)
		err = scheduler_k3s.TriggerPostAppRenameSetup(oldAppName, newAppName)
	case "post-delete":
		appName := flag.Arg(0)
		err = scheduler_k3s.TriggerPostDelete(appName)
	case "post-registry-login":
		server := flag.Arg(0)
		username := flag.Arg(1)
		password := flag.Arg(2)
		err = scheduler_k3s.TriggerPostRegistryLogin(server, username, password)
	case "report":
		appName := flag.Arg(0)
		err = scheduler_k3s.ReportSingleApp(appName, "", "")
	case "scheduler-deploy":
		scheduler := flag.Arg(0)
		appName := flag.Arg(1)
		imageTag := flag.Arg(2)
		err = scheduler_k3s.TriggerSchedulerDeploy(scheduler, appName, imageTag)
	case "scheduler-enter":
		scheduler := flag.Arg(0)
		appName := flag.Arg(1)
		containerType := flag.Arg(2)
		args := flag.Args()
		if len(args) > 3 {
			_, args = common.ShiftString(args)
			_, args = common.ShiftString(args)
			_, args = common.ShiftString(args)
		}

		err = scheduler_k3s.TriggerSchedulerEnter(scheduler, appName, containerType, ptr.Deref(podIdentifier, ""), args)
	default:
		err = fmt.Errorf("Invalid plugin trigger call: %s", trigger)
	}

	if err != nil {
		common.LogFailWithError(err)
	}
}
