/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/client-go/kubernetes"
)

var configFlags = genericclioptions.NewConfigFlags(true)
var clientK8s *kubernetes.Clientset
var pvcName string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "kpod-mount-pvc",
	Short: "Create a pod in k8s and mount a PVC into it",
	Long: `Create a pod in k8s and mount a PVC into it
	in the directory /data`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		clientK8s = configureK8sClient(configFlags)
	},

	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {
		ns := *configFlags.Namespace
		if ns == "" {
			ns = "default"
			fmt.Printf("The namespace option is empty, setting the namespace to %q\n", ns)
		}
		pvc, err := clientK8s.CoreV1().PersistentVolumeClaims(ns).Get(context.TODO(), pvcName, metav1.GetOptions{})
		if err != nil {
			fmt.Printf("Failed to get PVC name. Error stack:\n %+v\n", err)
			os.Exit(1)
		}
		if pvc.Status.Phase == corev1.ClaimPending {
			fmt.Println("Cannot attach pending PVC")
			os.Exit(1)
		}
		// create volume in memory
		fmt.Println("Creating the volume in memory")
		volumes := make([]corev1.Volume, 1)
		volumes[0] = corev1.Volume{
			Name: "volume1",
			VolumeSource: corev1.VolumeSource{
				PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{ClaimName: pvcName},
			},
		}

		// create pod in memory
		//TODO: test if the pod already exists in k8s, print msg and stop
		pod := corev1.Pod{}
		pod.ObjectMeta.Name = "mount"
		pod.ObjectMeta.Namespace = ns
		fmt.Printf("Creating the pod %q in memory\n", pod.ObjectMeta.Name)
		containers := make([]corev1.Container, 1)
		containers[0] = corev1.Container{
			Name:    "container1",
			Image:   "busybox",
			Command: []string{"/bin/sleep", "infinity"}, // with infinity, do not restart
			VolumeMounts: []corev1.VolumeMount{
				{Name: volumes[0].Name, MountPath: "/data"},
			},
		}
		pod.Spec.Containers = containers
		pod.Spec.Volumes = volumes

		// submit the pod to k8s
		fmt.Printf("Submitting the pod to k8s\n")
		_, err = clientK8s.CoreV1().Pods(ns).Create(context.TODO(), &pod, metav1.CreateOptions{})
		if err != nil {
			fmt.Printf("Failed to create POD. Error stack:\n %+v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Execution succeded\n")
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	configFlags.AddFlags(rootCmd.PersistentFlags())
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.kpod-mount-pvc.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	rootCmd.Flags().StringVarP(&pvcName, "pvc-name", "p", "", "Name of the PVC to mount (required)")
	rootCmd.MarkFlagRequired("pvc-name")
}

func configureK8sClient(flags *genericclioptions.ConfigFlags) *kubernetes.Clientset {
	restConfig, err := flags.ToRESTConfig()
	if err != nil {
		fmt.Printf("Failed to configure access to k8s: %+v\n", err)
		os.Exit(1)
	}
	client, err := kubernetes.NewForConfig(restConfig)
	if err != nil {
		fmt.Printf("Failed to build a K8S client: %+v\n", err)
		os.Exit(1)
	}
	return client
}
