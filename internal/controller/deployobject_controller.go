/*
Copyright 2024.

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

package controller

import (
	"context"
	"fmt"
	"github.com/go-logr/logr"
	appsv1 "k8s.io/api/apps/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/types"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	appv1 "kubebuilder.io/api/v1"
)

// DeployObjectReconciler reconciles a DeployObject object
type DeployObjectReconciler struct {
	client.Client
	Scheme *runtime.Scheme
	Log    logr.Logger
}

//大多数控制器最终会在群集上运行，因此它们需要 RBAC 权限，我们使用控制器工具 RBAC 标记指定这些权限。这些是运行所需的最低权限
// +kubebuilder:rbac:groups=app.kubebuilder.io,resources=deployobjects,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=app.kubebuilder.io,resources=deployobjects/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=app.kubebuilder.io,resources=deployobjects/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the DeployObject object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.19.0/pkg/reconcile

// 大多数控制器都需要 log 和 context，因此我们在此处设置它们。context 允许我们传递上下文，log 允许我们记录日志
func (r *DeployObjectReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := log.FromContext(ctx) //定义日志信息
	_ = r.Log.WithValues("deployobject", req.NamespacedName)

	// TODO(user): your logic here
	//appv1 "kubebuilder.io/api/v1"---appv1是上面import自动生成的
	//appv1.DeployObject{}是appv1组合上api/v1/xxx.types.go中的结构体定义 type DeployObject struct
	memcached := &appv1.DeployObject{}
	err := r.Get(ctx, req.NamespacedName, memcached)
	if err != nil {
		if apierrors.IsNotFound(err) {
			// 在资源未找到时候则停止reconciliation过程---If the custom resource is not found then it usually means that it was deleted or not created
			// In this way, we will stop the reconciliation
			log.Info("memcached 资源未找到resource not found，Ignoring since object must be deleted")
			return ctrl.Result{}, nil
		}
		// 如果读取出错则打印日志---Error reading the object - requeue the request.
		log.Error(err, "读取这个memcache资源错误")
		//并且返回err
		return ctrl.Result{}, err
	}

	// Let's just set the status as Unknown when no status is available
	// 当获取不到状态时候，将状态设置成unknown
	if memcached.Status.Conditions == nil || len(memcached.Status.Conditions) == 0 {
		//这个meta记得在import导入
		//	apierrors "k8s.io/apimachinery/pkg/api/errors"
		//	"k8s.io/apimachinery/pkg/api/meta"
		meta.SetStatusCondition(&memcached.Status.Conditions, metav1.Condition{
			Type:               "typeAvailableMemcached",
			Status:             "metav1.ConditionUnknown",
			ObservedGeneration: 0,
			LastTransitionTime: metav1.Time{},
			Reason:             "Reconciling",
			Message:            "Starting reconciliation",
		})
		// 在更新状态后，让我们重新获取 memcached 自定义资源
		// 以便我们能够获得集群中资源的最新状态，并避免
		// 引发错误“对象已被修改，请将您的更改应用到最新版本并重试”，
		// 如果我们在接下来的操作中再次尝试更新，将会重新触发调解
		if err = r.Status().Update(ctx, memcached); err != nil {
			log.Error(err, "Failed to re-fetch memcached")
			return ctrl.Result{}, err
		}
	}

	//确认应用resource是否存在,不存在则新建 Check if the deployment already exists, if not create a new one
	found := &appsv1.Deployment{}
	err = r.Get(ctx, types.NamespacedName{
		Namespace: memcached.Namespace,
		Name:      memcached.Name,
	}, found)
	if err != nil && apierrors.IsNotFound(err) {
		//到这了。error
		//dep, err := r.deploymentForMemcached(memcached)
		if err != nil{
			log.Error(err, "定义Deployment resource的时候失败Failed to define new Deployment resource for Memcached")
			meta.SetStatusCondition(&memcached.Status.Conditions, metav1.Condition{
				Type:               typeAvailableMemcached,
				Status:             metav1.ConditionFalse,
				ObservedGeneration: 0,
				LastTransitionTime: metav1.Time{},
				Reason:             "Reconciling",
				Message: fmt.Sprintf("Failed to create Deployment for the custom resource (%s): (%s)", memcached.Name, err)
			}
       })
	}
	if err := r.Status().Update(ctx, memcached); err != nil {
		log.Error(err, "Failed to update Memcached status")
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, err

}

// 最后，我们将此协调器添加到管理器，以便在管理器启动时启动它。
// SetupWithManager sets up the controller with the Manager.
func (r *DeployObjectReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&appv1.DeployObject{}).
		Complete(r)
}
