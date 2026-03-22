package cli

import (
	"context"
	"testing"

	portainer "github.com/portainer/portainer/api"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	kfake "k8s.io/client-go/kubernetes/fake"
)

func Test_GetServiceAccount(t *testing.T) {

	t.Run("returns error if non-existent", func(t *testing.T) {
		k := &KubeClient{
			cli:        kfake.NewSimpleClientset(),
			instanceID: "test",
		}
		tokenData := &portainer.TokenData{ID: 1}
		_, err := k.GetPortainerUserServiceAccount(tokenData)
		if err == nil {
			t.Error("GetPortainerUserServiceAccount should fail with service account not found")
		}
	})

	t.Run("succeeds for cluster admin role", func(t *testing.T) {
		k := &KubeClient{
			cli:        kfake.NewSimpleClientset(),
			instanceID: "test",
		}

		tokenData := &portainer.TokenData{
			ID:       1,
			Role:     portainer.AdministratorRole,
			Username: portainerClusterAdminServiceAccountName,
		}
		serviceAccount := &v1.ServiceAccount{
			ObjectMeta: metav1.ObjectMeta{
				Name: tokenData.Username,
			},
		}
		_, err := k.cli.CoreV1().ServiceAccounts(portainerNamespace).Create(context.Background(), serviceAccount, metav1.CreateOptions{})
		if err != nil {
			t.Errorf("failed to create service acount; err=%s", err)
		}
		defer func() {
			err := k.cli.CoreV1().ServiceAccounts(portainerNamespace).Delete(context.Background(), serviceAccount.Name, metav1.DeleteOptions{})
			require.NoError(t, err)
		}()

		sa, err := k.GetPortainerUserServiceAccount(tokenData)
		if err != nil {
			t.Errorf("GetPortainerUserServiceAccount should succeed; err=%s", err)
		}

		want := "portainer-sa-clusteradmin"
		if sa.Name != want {
			t.Errorf("GetServiceAccount should succeed and return correct sa name; got=%s want=%s", sa.Name, want)
		}
	})

	t.Run("succeeds for standard user role", func(t *testing.T) {
		k := &KubeClient{
			cli:        kfake.NewSimpleClientset(),
			instanceID: "test",
		}

		tokenData := &portainer.TokenData{
			ID:   1,
			Role: portainer.StandardUserRole,
		}
		serviceAccountName := UserServiceAccountName(int(tokenData.ID), k.instanceID)
		serviceAccount := &v1.ServiceAccount{
			ObjectMeta: metav1.ObjectMeta{
				Name: serviceAccountName,
			},
		}
		_, err := k.cli.CoreV1().ServiceAccounts(portainerNamespace).Create(context.Background(), serviceAccount, metav1.CreateOptions{})
		if err != nil {
			t.Errorf("failed to create service acount; err=%s", err)
		}
		defer func() {
			err := k.cli.CoreV1().ServiceAccounts(portainerNamespace).Delete(context.Background(), serviceAccount.Name, metav1.DeleteOptions{})
			require.NoError(t, err)
		}()

		sa, err := k.GetPortainerUserServiceAccount(tokenData)
		if err != nil {
			t.Errorf("GetPortainerUserServiceAccount should succeed; err=%s", err)
		}

		want := "portainer-sa-user-test-1"
		if sa.Name != want {
			t.Errorf("GetPortainerUserServiceAccount should succeed and return correct sa name; got=%s want=%s", sa.Name, want)
		}
	})

}

func TestGetServiceAccountDetails(t *testing.T) {
	t.Run("returns service account details", func(t *testing.T) {
		automount := false
		sa := &v1.ServiceAccount{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "my-sa",
				Namespace: "default",
				Labels:    map[string]string{"app": "web"},
			},
			AutomountServiceAccountToken: &automount,
			ImagePullSecrets: []v1.LocalObjectReference{
				{Name: "registry-secret"},
			},
		}
		kcl := &KubeClient{
			cli:        kfake.NewSimpleClientset(sa),
			instanceID: "test",
		}

		result, err := kcl.GetServiceAccount("default", "my-sa")
		require.NoError(t, err)

		assert.Equal(t, "my-sa", result.Name)
		assert.Equal(t, "default", result.Namespace)
		assert.Equal(t, &automount, result.AutomountServiceAccountToken)
		assert.Len(t, result.ImagePullSecrets, 1)
		assert.Equal(t, "registry-secret", result.ImagePullSecrets[0].Name)
		assert.Equal(t, map[string]string{"app": "web"}, result.Labels)
	})

	t.Run("returns error when service account not found", func(t *testing.T) {
		kcl := &KubeClient{
			cli:        kfake.NewSimpleClientset(),
			instanceID: "test",
		}

		_, err := kcl.GetServiceAccount("default", "does-not-exist")
		require.Error(t, err)
	})

	t.Run("marks system namespace accounts as system", func(t *testing.T) {
		sa := &v1.ServiceAccount{
			ObjectMeta: metav1.ObjectMeta{Name: "default", Namespace: "kube-system"},
		}
		ns := &v1.Namespace{
			ObjectMeta: metav1.ObjectMeta{Name: "kube-system"},
		}
		kcl := &KubeClient{
			cli:        kfake.NewSimpleClientset(sa, ns),
			instanceID: "test",
		}

		result, err := kcl.GetServiceAccount("kube-system", "default")
		require.NoError(t, err)
		assert.True(t, result.IsSystem)
	})

	t.Run("returns nil automount when not set", func(t *testing.T) {
		sa := &v1.ServiceAccount{
			ObjectMeta: metav1.ObjectMeta{Name: "my-sa", Namespace: "default"},
		}
		kcl := &KubeClient{
			cli:        kfake.NewSimpleClientset(sa),
			instanceID: "test",
		}

		result, err := kcl.GetServiceAccount("default", "my-sa")
		require.NoError(t, err)
		assert.Nil(t, result.AutomountServiceAccountToken)
	})
}

func TestGetServiceAccount_CreatesAndFetches(t *testing.T) {
	t.Run("returns annotations when set", func(t *testing.T) {
		sa := &v1.ServiceAccount{
			ObjectMeta: metav1.ObjectMeta{
				Name:        "annotated-sa",
				Namespace:   "default",
				Annotations: map[string]string{"example.com/key": "value"},
			},
		}
		kcl := &KubeClient{
			cli:        kfake.NewSimpleClientset(sa),
			instanceID: "test",
		}

		result, err := kcl.GetServiceAccount("default", "annotated-sa")
		require.NoError(t, err)
		assert.Equal(t, map[string]string{"example.com/key": "value"}, result.Annotations)
	})

	t.Run("round-trips UID correctly", func(t *testing.T) {
		sa := &v1.ServiceAccount{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "uid-sa",
				Namespace: "default",
				UID:       "abc-123-def",
			},
		}
		kcl := &KubeClient{
			cli:        kfake.NewSimpleClientset(sa),
			instanceID: "test",
		}

		result, err := kcl.GetServiceAccount("default", "uid-sa")
		require.NoError(t, err)
		assert.Equal(t, "abc-123-def", string(result.UID))
	})

	t.Run("creates service account and fetches it back", func(t *testing.T) {
		kcl := &KubeClient{
			cli:        kfake.NewSimpleClientset(),
			instanceID: "test",
		}

		sa := &v1.ServiceAccount{
			ObjectMeta: metav1.ObjectMeta{Name: "fresh-sa", Namespace: "staging"},
		}
		_, err := kcl.cli.CoreV1().ServiceAccounts("staging").Create(context.Background(), sa, metav1.CreateOptions{})
		require.NoError(t, err)

		result, err := kcl.GetServiceAccount("staging", "fresh-sa")
		require.NoError(t, err)
		assert.Equal(t, "fresh-sa", result.Name)
		assert.Equal(t, "staging", result.Namespace)
	})
}
