package main

import (
	"encoding/json"
	"fmt"

	"github.com/ghodss/yaml"
)


// events: [
// 	event1,
// ]
// source: [
// 	source1,
// ]

func list() {
	获取某项目某分支缓存中的事件列表
	如果缓存中不存在，利用该用户发送api请求获取整个仓库的文件
	解析文件保存事件定义到缓存中
	如果前端请求不到事件列表需要进行轮询
	events := listFile(projectId, branch)
	if len(events) == 0 {
		// send_task(get gitlab achieve)
	}
	// font polling
}

func main() {
	j := []byte(`{"name": "John", "age": 30}`)
	y, err := yaml.JSONToYAML(j)
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return
	}
	fmt.Println(string(y))

	content := `
apiVersion: v1
kind: PersistentVolume
metadata:
  name: mysql-pv-volume
  labels:
    type: local
spec:
  storageClassName: manual
  capacity:
    storage: 20Gi
  accessModes:
    - ReadWriteOnce
  hostPath:
    path: "/mnt/data"
---
apiVersion: v1
kind: PersistentVolumeClaim
metadata:
  name: mysql-pv-claim
spec:
  storageClassName: manual
  accessModes:
    - ReadWriteOnce
  resources:
    requests:
      storage: 20Gi
`

	/* Output:
	name: John
	age: 30
	*/
	j2, err := yaml.YAMLToJSON([]byte(content))
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return
	}
	fmt.Println(string(j2))
	data := make(map[string]interface{})
	json.Unmarshal(j2, &data)
	data2 := convert(data)
	fmt.Println(data2["metadata"])

	/* Output:
	{"age":30,"name":"John"}
	*/
}

type Event struct {
	ApiVersion `"json:apiVersion"`
	Kind       `"json:kind"`
	// 	metadata
	//   name: mysql-pv-volume
	//   labels:
	//     type: local
	// spec:
	//   storageClassName: manual
	//   capacity:
	//     storage: 20Gi
	//   accessModes:
	//     - ReadWriteOnce
	//   hostPath:
	//     path: "/mnt/data"
	// ---
	// apiVersion: v1
	// kind: PersistentVolumeClaim
	// metadata:
	//   name: mysql-pv-claim
	// spec:
	//   storageClassName: manual
	//   accessModes:
	//     - ReadWriteOnce
	//   resources:
	//     requests:
	//       storage: 20Gi
}

func convert(i interface{}) interface{} {
	switch x := i.(type) {
	case map[interface{}]interface{}:
		m2 := map[string]interface{}{}
		for k, v := range x {
			m2[k.(string)] = convert(v)
		}
		return m2
	case []interface{}:
		for i, v := range x {
			x[i] = convert(v)
		}
	}
	return i
}
