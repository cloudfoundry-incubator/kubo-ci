package test_helpers

import (
	"fmt"
	"time"

	"k8s.io/apimachinery/pkg/watch"

	. "github.com/onsi/ginkgo"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	tappsv1 "k8s.io/client-go/kubernetes/typed/apps/v1"
)

func GetNginxDeploymentSpec() appsv1.DeploymentSpec {
	nginxPodSpec := corev1.PodSpec{
		Containers: []corev1.Container{{
			Name:  "nginx",
			Image: "nginx",
			Ports: []corev1.ContainerPort{{ContainerPort: 80}},
		}},
	}
	var replicas int32
	replicas = 1
	labelMap := make(map[string]string)
	labelMap["app"] = "nginx"

	return appsv1.DeploymentSpec{
		Replicas: &replicas,
		Selector: &metav1.LabelSelector{MatchLabels: labelMap},
		Template: corev1.PodTemplateSpec{ObjectMeta: metav1.ObjectMeta{Labels: labelMap}, Spec: nginxPodSpec}}
}

func WaitForDeployment(deploymentAPI tappsv1.DeploymentInterface, namespace string, deploymentName string) error {
	w, err := deploymentAPI.Watch(metav1.ListOptions{})
	if err != nil {
		return err
	}

	_, err = watch.Until(1*time.Minute, w, func(event watch.Event) (bool, error) {
		deployment, ok := event.Object.(*appsv1.Deployment)
		if !ok {
			return false, fmt.Errorf("Expected `%#v` to be of type appsv1.Deployment", event.Object)
		}

		if deployment.Name == deploymentName {
			if deployment.Status.AvailableReplicas == deployment.Status.UpdatedReplicas {
				return true, nil
			}
			fmt.Fprintf(GinkgoWriter, "Expected %d to be equal to %d\n", deployment.Status.AvailableReplicas, deployment.Status.UpdatedReplicas)
		}

		return false, nil
	})

	if err != nil {
		return fmt.Errorf("Deployment `%s` did not finish rolling out with error: %s", deploymentName, err)
	}

	return nil
}

func NewDeployment(name string, spec appsv1.DeploymentSpec) *appsv1.Deployment {
	return &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{Name: name},
		Spec:       spec,
	}
}
