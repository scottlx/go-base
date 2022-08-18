package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

// 定义自定义API对象Custom.
type Custom struct {
	// 继承metav1.TypeMeta和metav1.ObjectMeta才实现了runtime.Object，这样Custom对象
	// 的yaml的格式就像如下：
	// kind: Custom
	// apiVersion: test/v1
	// metadata:
	//   labels:
	//     name: custom
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`
	// 为了演示方便，Custom对象的规格和状态都定义为空类型，读者根据自己的业务进行设计
	Spec   CustomSpec
	Status CustomStatus
}
type CustomSpec struct{}
type CustomStatus struct{}

// DeepCopyObject()是必须要实现的，这是runtime.Objec定义的接口，否则编译就会报错。读者需要注意，
// kubernetes的API对象的DeepCopyObject()函数是代码生成工具生成的，本文的示例是笔者自己写的。
func (in *Custom) DeepCopyObject() runtime.Object {
	if in == nil {
		return nil
	}
	out := new(Custom)
	*out = *in
	return out
}

var _ runtime.Object = &Custom{}

// 定义自定义类型组名，因为是测试例子，所以笔者把组名定义为test，这样test/Custom才是类型全称
const GroupName = "test"

// 定义自定义类型的组名+版本，即test v1
var SchemeGroupVersion = schema.GroupVersion{Group: GroupName, Version: "v1"}

// 程序初始化部分需要调用AddToScheme()来实现自定义类型的注册，具体的实现在addKnownTypes()
var (
	SchemeBuilder = runtime.NewSchemeBuilder(addKnownTypes)
	AddToScheme   = SchemeBuilder.AddToScheme
)

// 把笔者上面定义的Custom添加到scheme中就完成了类型的注册，就是这么简单。读者需要注意，类型注册
// 其实是一个比较复杂的过程，kubernetes把这部分实现全部交给了scheme，把简单的接口留给了使用者。
func addKnownTypes(scheme *runtime.Scheme) error {
	scheme.AddKnownTypes(SchemeGroupVersion,
		&Custom{},
	)

	return nil
}
