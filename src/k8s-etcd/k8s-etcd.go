package main

import (
	"context"
	"fmt"

	testv1 "go-base/src/k8s-etcd/v1"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"k8s.io/apiserver/pkg/storage/storagebackend"
	"k8s.io/apiserver/pkg/storage/storagebackend/factory"
)

func main() {
	// 构造scheme，然后调用刚刚实现的注册函数实现自定义API对象的注册。
	scheme := runtime.NewScheme()
	testv1.AddToScheme(scheme)

	// 注册了自定义API对象，就可以构造codec工厂了，在通过codec工厂构造codec。所谓的codec就是
	// 自定义API对象的序列化与反序列化的实现。
	cf := serializer.NewCodecFactory(scheme)
	codec := cf.LegacyCodec(testv1.SchemeGroupVersion)

	// 有了codec就可以创建storagebackend.Config，他是构造storage.Interfaces的必要参数
	config := storagebackend.NewDefaultConfig("", codec)

	// 笔者在本机装了etcd，所以把etcd的地址填写为本机地址，这样方便测试
	config.Transport.ServerList = append(config.Transport.ServerList, "127.0.0.1:2379")
	// 创建storage.Interfaces对象，storage.Interfaces是apiserver对存储的抽象，这样我们
	// 就可以像apiserver一样在etcd上操作自己定义的对象了。

	// 构造Custom对象
	custom := testv1.Custom{}
	storage, destroy, err := factory.Create(*config.ForResource(schema.ParseGroupResource("")), custom.DeepCopyObject)
	if nil != err {
		fmt.Printf("%v\n", err)
	}

	// 把Custom对象写入etcd
	if err = storage.Create(context.Background(), "test", &custom, &custom, 0); nil != err {
		fmt.Println(err)
	}
	// 把写入的Custom对象打印出来看看结果
	if data, err := runtime.Encode(codec, &custom); nil == err {
		fmt.Printf("%s\n", string(data))
	}
	// 必要的析构函数
	destroy()
}
