/*
Copyright 2016 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package cmd

import (
	"fmt"
	"io"

	"github.com/renstrom/dedent"
	"github.com/spf13/cobra"

	"k8s.io/apiserver/pkg/util/flag"
	"k8s.io/kubernetes/cmd/kubeadm/app/cmd/phases"
	cmdutil "k8s.io/kubernetes/pkg/kubectl/cmd/util"
)

func NewKubeadmCommand(f cmdutil.Factory, in io.Reader, out, err io.Writer) *cobra.Command {
	cmds := &cobra.Command{
		Use:   "kubeadm",
		Short: "kubeadm: easily bootstrap a secure Kubernetes cluster",
		Long: dedent.Dedent(`
			kubeadm: easily bootstrap a secure Kubernetes cluster.

			    ┌──────────────────────────────────────────────────────────┐
			    │ KUBEADM IS BETA, DO NOT USE IT FOR PRODUCTION CLUSTERS!  │
			    │                                                          │
			    │ But, please try it out! Give us feedback at:             │
			    │ https://github.com/kubernetes/kubeadm/issues             │
			    │ and at-mention @kubernetes/sig-cluster-lifecycle-misc    │
			    └──────────────────────────────────────────────────────────┘

			Example usage:

			    Create a two-machine cluster with one master (which controls the cluster),
			    and one node (where your workloads, like Pods and ReplicaSets run).

			    ┌──────────────────────────────────────────────────────────┐
			    │ On the first machine                                     │
			    ├──────────────────────────────────────────────────────────┤
			    │ master# kubeadm init                                     │
			    └──────────────────────────────────────────────────────────┘

			    ┌──────────────────────────────────────────────────────────┐
			    │ On the second machine                                    │
			    ├──────────────────────────────────────────────────────────┤
			    │ node# kubeadm join --token=<token> <ip-of-master>:<port> │
			    └──────────────────────────────────────────────────────────┘

			    You can then repeat the second step on as many other machines as you like.

		`),
	}

	cmds.ResetFlags()
	cmds.SetGlobalNormalizationFunc(flag.WarnWordSepNormalizeFunc)

	cmds.AddCommand(NewCmdCompletion(out, ""))
	cmds.AddCommand(NewCmdConfig(out))
	cmds.AddCommand(NewCmdInit(out))
	cmds.AddCommand(NewCmdJoin(out))
	cmds.AddCommand(NewCmdReset(out))
	cmds.AddCommand(NewCmdVersion(out))
	cmds.AddCommand(NewCmdToken(out, err))

	// Wrap not yet fully supported commands in an alpha subcommand
	experimentalCmd := &cobra.Command{
		Use:   "alpha",
		Short: "Experimental sub-commands not yet fully functional.",
	}
	experimentalCmd.AddCommand(phases.NewCmdPhase(out))
	cmds.AddCommand(experimentalCmd)

	return cmds
}

// subCmdRunE returns a function that handles a case where a subcommand must be specified
// Without this callback, if a user runs just the command without a subcommand,
// or with an invalid subcommand, cobra will print usage information, but still exit cleanly.
// We want to return an error code in these cases so that the
// user knows that their command was invalid.
func subCmdRunE(name string) func(*cobra.Command, []string) error {
	return func(_ *cobra.Command, args []string) error {
		if len(args) < 1 {
			return fmt.Errorf("missing subcommand; %q is not meant to be run on its own", name)
		}

		return fmt.Errorf("invalid subcommand: %q", args[0])
	}
}