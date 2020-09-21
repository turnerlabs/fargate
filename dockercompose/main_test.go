package dockercompose

import (
	"fmt"
	"testing"
)

func doTest(t *testing.T, f string) (ComposeFile, error) {
	//round-trip unmarshal and marshal
	file, e := Read(f)
	if e != nil {
		return file, e
	}
	fmt.Println(file.Data)
	b, e := file.Yaml()
	if e != nil {
		return file, e
	}
	t.Log(string(b))
	return file, nil
}

const (
	image         = "1234567890.dkr.ecr.us-east-1.amazonaws.com/my-service:0.1.0"
	publishedPort = 80
	targetPort    = 8080
	labelKey      = "aws.ecs.fargate.deploy"
	labelValue    = "1"
	secretKey     = "QUX"
	secretValue   = "arn:key:ssm:us-east-1:000000000000:parameter/path/to/my_parameter"
)

// 2
func TestComposeV2(t *testing.T) {
	f, e := doTest(t, "v2.yml")
	if e != nil {
		t.Error(e)
		return
	}
	svc := f.Data.Services["web"]
	if svc.Image != image {
		t.Error("expecting image")
	}
	if svc.Ports[0].Published != publishedPort {
		t.Error("expecting published port")
	}
	if svc.Ports[0].Target != targetPort {
		t.Error("expecting published port")
	}
	if svc.Labels[labelKey] != labelValue {
		t.Error("expecting label")
	}
}

// 2.4
func TestComposeV24(t *testing.T) {
	f, e := doTest(t, "v2.4.yml")
	if e != nil {
		t.Error(e)
		return
	}
	svc := f.Data.Services["web"]
	if svc.Image != image {
		t.Error("expecting image")
	}
	if svc.Ports[0].Published != publishedPort {
		t.Error("expecting published port")
	}
	if svc.Ports[0].Target != targetPort {
		t.Error("expecting published port")
	}
	if svc.Labels[labelKey] != labelValue {
		t.Error("expecting label")
	}
	if svc.Secrets[secretKey] != secretValue {
		t.Error("expecting secret")
	}
}

// 3.2 short
func TestComposeV32Short(t *testing.T) {
	f, e := doTest(t, "v3.2-short.yml")
	if e != nil {
		t.Error(e)
		return
	}
	fmt.Println(f.Data)
	svc := f.Data.Services["web"]
	if svc.Image != image {
		t.Error("expecting image")
	}
	if svc.Ports[0].Published != publishedPort {
		t.Error("expecting published port")
	}
	if svc.Ports[0].Target != targetPort {
		t.Error("expecting published port")
	}
	if svc.Labels[labelKey] != labelValue {
		t.Error("expecting label")
	}
}

// 3.2 long
func TestComposeV32Long(t *testing.T) {
	f, e := doTest(t, "v3.2-long.yml")
	if e != nil {
		t.Error(e)
		return
	}
	svc := f.Data.Services["web"]
	if svc.Image != image {
		t.Error("expecting image")
	}
	if svc.Ports[0].Published != publishedPort {
		t.Error("expecting published port")
	}
	if svc.Ports[0].Target != targetPort {
		t.Error("expecting published port")
	}
	if svc.Labels[labelKey] != labelValue {
		t.Error("expecting label")
	}
}

// 3.7
func TestComposeV37(t *testing.T) {
	f, e := doTest(t, "v3.7.yml")
	if e != nil {
		t.Error(e)
		return
	}
	svc := f.Data.Services["web"]
	if svc.Image != image {
		t.Error("expecting image")
	}
	if svc.Ports[0].Published != publishedPort {
		t.Error("expecting published port")
	}
	if svc.Ports[0].Target != targetPort {
		t.Error("expecting published port")
	}
	if svc.Labels[labelKey] != labelValue {
		t.Error("expecting label")
	}
	if svc.Secrets[secretKey] != secretValue {
		t.Error("expecting secret")
	}
}
