package images

import (
	"context"
	"fmt"
        "io/ioutil"
	urlP "net/url"
	"os"
	"os/exec"
	"os/user"
	"strconv"
	"strings"
        "time"
	"github.com/containers/common/pkg/config"
	"github.com/containers/podman/v3/cmd/podman/common"
	"github.com/containers/podman/v3/cmd/podman/parse"
	"github.com/containers/podman/v3/cmd/podman/registry"
	"github.com/containers/podman/v3/cmd/podman/system/connection"
	"github.com/containers/podman/v3/libpod/define"
	"github.com/containers/podman/v3/pkg/domain/entities"
	"github.com/containers/podman/v3/pkg/rootless"
	"github.com/docker/distribution/reference"
	scpD "github.com/dtylman/scp"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"golang.org/x/crypto/ssh"
)

var (
	saveScpDescription = `Securely copy an image from one host to another.`
	imageScpCommand    = &cobra.Command{
		Use: "scp [options] IMAGE [HOST::]",
		Annotations: map[string]string{
			registry.UnshareNSRequired: "",
			registry.ParentNSRequired:  "",
			registry.EngineMode:        registry.ABIMode,
		},
		Long:              saveScpDescription,
		Short:             "securely copy images",
		RunE:              scp,
		Args:              cobra.RangeArgs(1, 2),
		ValidArgsFunction: common.AutocompleteScp,
		Example:           `podman image scp myimage:latest otherhost::`,
	}
)

var (
	parentFlags []string
	scpOpts     entities.ImageScpOptions
)

func init() {
	registry.Commands = append(registry.Commands, registry.CliCommand{
		Command: imageScpCommand,
		Parent:  imageCmd,
	})
	scpFlags(imageScpCommand)
}

func scpFlags(cmd *cobra.Command) {
	flags := cmd.Flags()
	flags.BoolVarP(&scpOpts.Save.Quiet, "quiet", "q", false, "Suppress the output")
}

func scp(cmd *cobra.Command, args []string) (finalErr error) {
	var (
		// TODO add tag support for images
		err error
	)
	for i, val := range os.Args {
		if val == "image" {
			break
		} else if i == 0 {
			continue
		} else if strings.Contains(val, "CIRRUS") {
			continue
		}
		parentFlags = append(parentFlags, val)
	}
	podman, err := os.Executable()
	if err != nil {
		return err
	}
	if scpOpts.Save.Quiet { // set quiet for both load and save
		scpOpts.Load.Quiet = true
	}
	f, err := ioutil.TempFile("", "podman") // open temp file for load/save output
	if err != nil {
		return err
	}
	scpOpts.Save.Output = f.Name()
	scpOpts.Save.Format = "oci-archive"
	err = os.Remove(f.Name()) // remove the file and simply use its name so podman creates the file upon save. avoids umask errors
	if err != nil {
		return err
	}
	_, err = os.Open(scpOpts.Save.Output)
	if err == nil {
		return errors.New("temporary file already exists. If this issue persists please clear out /var/tmp")
	}
        time.Sleep(5 * time.Second)
     	scpOpts.Load.Input = scpOpts.Save.Output
	if err := parse.ValidateFileName(saveOpts.Output); err != nil {
		return err
	}
	confR, err := config.NewConfig("") // create a hand made config for the remote engine since we might use remote and native at once
	if err != nil {
		return errors.Wrapf(err, "could not make config")
	}

	abiEng, err := registry.NewImageEngine(cmd, args) // abi native engine
	if err != nil {
		return err
	}

	cfg, err := config.ReadCustomConfig() // get ready to set ssh destination if necessary
	if err != nil {
		return err
	}
	serv, err := parseArgs(args, cfg) // parses connection data and "which way" we are loading and saving
	if err != nil {
		return err
	}
	// TODO: Add podman remote support
	confR.Engine = config.EngineConfig{Remote: true, CgroupManager: "cgroupfs", ServiceDestinations: serv} // pass the service dest (either remote or something else) to engine
	switch {
	case scpOpts.FromRemote: // if we want to load FROM the remote
		err = saveToRemote(scpOpts.SourceImageName, scpOpts.Save.Output, "", scpOpts.URI[0], scpOpts.Iden[0])
		if err != nil {
			return err
		}
		if scpOpts.ToRemote { // we want to load remote -> remote
			rep, err := loadToRemote(scpOpts.Save.Output, "", scpOpts.URI[1], scpOpts.Iden[1])
			if err != nil {
				return err
			}
			fmt.Println(rep)
			break
		}
		err = execLoad(podman)
		if err != nil {
			return err
		}
	case scpOpts.ToRemote: // remote host load
		err = execSave(podman)
		if err != nil {
			return err
		}
		rep, err := loadToRemote(scpOpts.Save.Output, "", scpOpts.URI[0], scpOpts.Iden[0])
		if err != nil {
			return err
		}
		fmt.Println(rep)
		err = os.Remove(scpOpts.Save.Output)
		if err != nil {
			return err
		}
	// TODO: Add podman remote support
	default: // else native load
		if scpOpts.Tag != "" {
			return errors.Wrapf(define.ErrInvalidArg, "Renaming of an image is currently not supported")
		}
		if scpOpts.Transfer { // if we ae tranferring between users...
			var u *user.User
			if scpOpts.User != "" {
				scpUser := os.Getenv("USER")
				if !rootless.IsRootless() {
					scpUser = os.Getenv("SUDO_USER")
				}
				half := strings.Split(scpOpts.User, ":")[0] //if we are given a uid:gid we need to convert to a username before checking if the user is proper
				_, err := strconv.Atoi(half)
				if err == nil {
					u, err = user.LookupId(half)
				} else {
					u, err = user.Lookup(scpOpts.User)
				}
				if err != nil {
					return err
				}
				if scpUser != "" && scpUser != u.Username && u.Username != "root" {
					return errors.New("the given user must be the default rootless user or root")
				}
			} else if scpOpts.User == "" {
				scpOpts.User = os.Getenv("USER")
				if !rootless.IsRootless() {
					scpOpts.User = os.Getenv("SUDO_USER")
				}
				if scpOpts.User == "" {
					return errors.New("could not obtain user, make sure the environmental variable $USER is set")
				}
			}
			err := abiEng.Transfer(context.Background(), scpOpts, parentFlags)
			if err != nil {
				return err
			}
		} else { // else do the default (save and load, no affect)
			err = execSave(podman)
			if err != nil {
				return err
			}
			err = execLoad(podman)
			if err != nil {
				return err
			}
			err = os.Remove(scpOpts.Save.Output)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// loadToRemote takes image and remote connection information. it connects to the specified client
// and copies the saved image dir over to the remote host and then loads it onto the machine
// returns a string containing output or an error
func loadToRemote(localFile string, tag string, url *urlP.URL, iden string) (string, error) {
	dial, remoteFile, err := createConnection(url, iden)
	if err != nil {
		return "", err
	}
	defer dial.Close()

	n, err := scpD.CopyTo(dial, localFile, remoteFile)
	if err != nil {
		errOut := strconv.Itoa(int(n)) + " Bytes copied before error"
		return " ", errors.Wrapf(err, errOut)
	}
	run := ""
	if tag != "" {
		return "", errors.Wrapf(define.ErrInvalidArg, "Renaming of an image is currently not supported")
	}
	podman := os.Args[0]
	run = podman + " image load --input=" + remoteFile + ";rm " + remoteFile // run ssh image load of the file copied via scp
	out, err := connection.ExecRemoteCommand(dial, run)
	if err != nil {
		return "", err
	}
	return strings.TrimSuffix(string(out), "\n"), nil
}

// saveToRemote takes image information and remote connection information. it connects to the specified client
// and saves the specified image on the remote machine and then copies it to the specified local location
// returns an error if one occurs.
func saveToRemote(image, localFile string, tag string, uri *urlP.URL, iden string) error {
	dial, remoteFile, err := createConnection(uri, iden)

	if err != nil {
		return err
	}
	defer dial.Close()

	if tag != "" {
		return errors.Wrapf(define.ErrInvalidArg, "Renaming of an image is currently not supported")
	}
	podman := os.Args[0]
	run := podman + " image save " + image + " --format=oci-archive --output=" + remoteFile // run ssh image load of the file copied via scp. Files are reverse in this case...
	_, err = connection.ExecRemoteCommand(dial, run)
	if err != nil {
		return nil
	}
	n, err := scpD.CopyFrom(dial, remoteFile, localFile)
	connection.ExecRemoteCommand(dial, "rm "+remoteFile)
	if err != nil {
		errOut := strconv.Itoa(int(n)) + " Bytes copied before error"
		return errors.Wrapf(err, errOut)
	}
	return nil
}

// makeRemoteFile creates the necessary remote file on the host to
// save or load the image to. returns a string with the file name or an error
func makeRemoteFile(dial *ssh.Client) (string, error) {
	run := "mktemp"
	remoteFile, err := connection.ExecRemoteCommand(dial, run)
	if err != nil {
		return "", err
	}
	return strings.TrimSuffix(string(remoteFile), "\n"), nil
}

// createConnections takes a boolean determining which ssh client to dial
// and returns the dials client, its newly opened remote file, and an error if applicable.
func createConnection(url *urlP.URL, iden string) (*ssh.Client, string, error) {
	cfg, err := connection.ValidateAndConfigure(url, iden)
	if err != nil {
		return nil, "", err
	}
	dialAdd, err := ssh.Dial("tcp", url.Host, cfg) // dial the client
	if err != nil {
		return nil, "", errors.Wrapf(err, "failed to connect")
	}
	file, err := makeRemoteFile(dialAdd)
	if err != nil {
		return nil, "", err
	}

	return dialAdd, file, nil
}

// validateImageName makes sure that the image given is valid and no injections are occurring
// we simply use this for error checking, bot setting the image
func validateImageName(input string) error {
	// ParseNormalizedNamed transforms a shortname image into its
	// full name reference so busybox => docker.io/library/busybox
	// we want to keep our shortnames, so only return an error if
	// we cannot parse what th euser has given us
	_, err := reference.ParseNormalizedNamed(input)
	return err
}

// remoteArgLength is a helper function to simplify the extracting of host argument data
// returns an int which contains the length of a specified index in a host::image string
func remoteArgLength(input string, side int) int {
	return len((strings.Split(input, "::"))[side])
}

// parseArgs returns the valid connection data based off of the information provided by the user
// args is an array of the command arguments and cfg is tooling configuration used to get service destinations
// returned is serv and an error if applicable. serv is a map of service destinations with the connection name as the index
// this connection name is intended to be used as EngineConfig.ServiceDestinations
// this function modifies the global scpOpt entities: FromRemote, ToRemote, Connections, and SourceImageName
func parseArgs(args []string, cfg *config.Config) (map[string]config.Destination, error) {
	serv := map[string]config.Destination{}
	cliConnections := []string{}
	switch len(args) {
	case 1:
		if strings.Contains(args[0], "@localhost") { // image transfer between users
			scpOpts.User = strings.Split(args[0], "@")[0]
			scpOpts.Transfer = true
			if len(strings.Split(args[0], "::")) > 1 {
				scpOpts.SourceImageName = strings.Split(args[0], "::")[1]
			} else {
				return nil, errors.New("no image provided")
			}
		} else if strings.Contains(args[0], "::") {
			scpOpts.FromRemote = true
			cliConnections = append(cliConnections, args[0])
		} else {
			err := validateImageName(args[0])
			if err != nil {
				return nil, err
			}
			scpOpts.SourceImageName = args[0]
		}
	case 2:
		if strings.Contains(args[0], "@localhost") || strings.Contains(args[1], "@localhost") { // image transfer between users
			scpOpts.Transfer = true
			if len(strings.Split(args[0], "::")) > 1 && len(strings.Split(args[0], "::")[1]) > 0 { // first argument contains image and therefore is our user
				scpOpts.User = strings.Split(args[0], "@")[0]
				scpOpts.SourceImageName = strings.Split(args[0], "::")[1]
			} else if len(strings.Split(args[1], "::")) > 1 && len(strings.Split(args[1], "::")[1]) > 0 { // second argument contains image and therefore is our user
				scpOpts.User = strings.Split(args[1], "@")[0]
				scpOpts.SourceImageName = strings.Split(args[1], "::")[1]
			} else {
				return nil, errors.New("no image provided")
			}
		} else if strings.Contains(args[0], "::") {
			if !(strings.Contains(args[1], "::")) && remoteArgLength(args[0], 1) == 0 { // if an image is specified, this mean we are loading to our client
				cliConnections = append(cliConnections, args[0])
				scpOpts.ToRemote = true
				scpOpts.SourceImageName = args[1]
			} else if strings.Contains(args[1], "::") { // both remote clients
				scpOpts.FromRemote = true
				scpOpts.ToRemote = true
				if remoteArgLength(args[0], 1) == 0 { // is save->load w/ one image name
					cliConnections = append(cliConnections, args[0])
					cliConnections = append(cliConnections, args[1])
				} else if remoteArgLength(args[0], 1) > 0 && remoteArgLength(args[1], 1) > 0 {
					//in the future, this function could, instead of rejecting renames, also set a DestImageName field
					return nil, errors.Wrapf(define.ErrInvalidArg, "cannot specify an image rename")
				} else { // else its a load save (order of args)
					cliConnections = append(cliConnections, args[1])
					cliConnections = append(cliConnections, args[0])
				}
			} else {
				//in the future, this function could, instead of rejecting renames, also set a DestImageName field
				return nil, errors.Wrapf(define.ErrInvalidArg, "cannot specify an image rename")
			}
		} else if strings.Contains(args[1], "::") { // if we are given image host::
			if remoteArgLength(args[1], 1) > 0 {
				//in the future, this function could, instead of rejecting renames, also set a DestImageName field
				return nil, errors.Wrapf(define.ErrInvalidArg, "cannot specify an image rename")
			}
			err := validateImageName(args[0])
			if err != nil {
				return nil, err
			}
			scpOpts.SourceImageName = args[0]
			scpOpts.ToRemote = true
			cliConnections = append(cliConnections, args[1])
		} else {
			//in the future, this function could, instead of rejecting renames, also set a DestImageName field
			return nil, errors.Wrapf(define.ErrInvalidArg, "cannot specify an image rename")
		}
	}
	var url string
	var iden string
	for i, val := range cliConnections {
		splitEnv := strings.SplitN(val, "::", 2)
		scpOpts.Connections = append(scpOpts.Connections, splitEnv[0])
		if len(splitEnv[1]) != 0 {
			err := validateImageName(splitEnv[1])
			if err != nil {
				return nil, err
			}
			scpOpts.SourceImageName = splitEnv[1]
			//TODO: actually use the new name given by the user
		}
		conn, found := cfg.Engine.ServiceDestinations[scpOpts.Connections[i]]
		if found {
			url = conn.URI
			iden = conn.Identity
		} else { // no match, warn user and do a manual connection.
			url = "ssh://" + scpOpts.Connections[i]
			iden = ""
			logrus.Warnf("Unknown connection name given. Please use system connection add to specify the default remote socket location")
		}
		urlT, err := urlP.Parse(url) // create an actual url to pass to exec command
		if err != nil {
			return nil, err
		}
		if urlT.User.Username() == "" {
			if urlT.User, err = connection.GetUserInfo(urlT); err != nil {
				return nil, err
			}
		}
		scpOpts.URI = append(scpOpts.URI, urlT)
		scpOpts.Iden = append(scpOpts.Iden, iden)
	}
	return serv, nil
}

// execLoad executes the podman load command given the podman binary
func execLoad(podman string) error {
	var loadCommand []string
	if scpOpts.Save.Quiet || scpOpts.Load.Quiet {
		loadCommand = []string{"load", "-q", "--input", scpOpts.Save.Output}
	} else {
		loadCommand = []string{"load", "--input", scpOpts.Save.Output}
	}
	if rootless.IsRootless() {
		cmd := exec.Command(podman)
		cmd.Args = append(cmd.Args, parentFlags...)
		cmd.Args = append(cmd.Args, loadCommand...)
		cmd.Env = os.Environ()
		cmd.Stderr = os.Stderr
		cmd.Stdout = os.Stdout
		logrus.Debug("Executing load command")
		err := cmd.Run()
		if err != nil {
			return err
		}
		return nil
	}
	machinectl, err := exec.LookPath("machinectl")
	if err != nil {
		cmd := exec.Command("su", "-l", "root", "--command")
		fullCommand := podman
		for _, val := range parentFlags {
			fullCommand += (" " + val)
		}
		for _, val := range loadCommand {
			fullCommand += (" " + val)
		}
		cmd.Args = append(cmd.Args, fullCommand)
		cmd.Env = os.Environ()
		cmd.Stderr = os.Stderr
		cmd.Stdout = os.Stdout
		logrus.Debug("Executing load command su")
		err = cmd.Run()
		if err != nil {
			return err
		}
	} else {
		cmd := exec.Command(machinectl, "shell", "-q", "root@.host")
		cmd.Args = append(cmd.Args, podman)
		cmd.Args = append(cmd.Args, parentFlags...)
		cmd.Args = append(cmd.Args, loadCommand...)
		cmd.Env = os.Environ()
		cmd.Stderr = os.Stderr
		cmd.Stdout = os.Stdout
		logrus.Debug("Executing load command machinectl")
		err = cmd.Run()
		if err != nil {
			return err
		}
	}
	return nil
}

// execSave executes the podman save command given the podman binary
func execSave(podman string) error {
	var saveCommand []string
	if scpOpts.Save.Quiet || scpOpts.Load.Quiet {
		saveCommand = []string{"save", "-q", "--output", scpOpts.Save.Output, scpOpts.SourceImageName}
	} else {
		saveCommand = []string{"save", "--output", scpOpts.Save.Output, scpOpts.SourceImageName}
	}
	if rootless.IsRootless() {
		cmd := exec.Command(podman)
		cmd.Args = append(cmd.Args, parentFlags...)
		cmd.Args = append(cmd.Args, saveCommand...)
		cmd.Env = os.Environ()
		cmd.Stderr = os.Stderr
		cmd.Stdout = os.Stdout
		logrus.Debug("Executing save command")
		err := cmd.Run()
		if err != nil {
			return err
		}
		return nil
	}
	machinectl, err := exec.LookPath("machinectl")
	if err != nil {
		cmd := exec.Command("su", "-l", "root", "--command")
		fullCommand := podman
		for _, val := range parentFlags {
			fullCommand += (" " + val)
		}
		for _, val := range saveCommand {
			fullCommand += (" " + val)
		}
		cmd.Args = append(cmd.Args, fullCommand)
		cmd.Env = os.Environ()
		cmd.Stderr = os.Stderr
		cmd.Stdout = os.Stdout
		logrus.Debug("Executing save command su")
		err = cmd.Run()
		if err != nil {
			return err
		}
	} else {
		cmd := exec.Command(machinectl, "shell", "-q", "root@.host")
		cmd.Args = append(cmd.Args, podman)
		cmd.Args = append(cmd.Args, parentFlags...)
		cmd.Args = append(cmd.Args, saveCommand...)
		cmd.Env = os.Environ()
		cmd.Stderr = os.Stderr
		cmd.Stdout = os.Stdout
		logrus.Debug("Executing save command machinectl")
		err = cmd.Run()
		if err != nil {
			return err
		}
	}
	return nil
}
