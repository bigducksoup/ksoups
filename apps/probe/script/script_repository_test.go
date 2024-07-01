package script

import (
	"context"
	"fmt"
	"testing"
)

func TestA(t *testing.T) {

	repo := NewShellScriptRepo("../../scripts")

	id, err := repo.Store([]byte("echo helloworld"), "helloworld", nil)

	if err != nil {
		t.Fatal(err)
	}

	script, ok := repo.Get(id)

	if !ok {
		t.Fatal("no script")
	}

	scriptRunner := NewShellScriptRunner()

	res, err := scriptRunner.Run(script)

	if err != nil {
		t.Fatal(err)
	}

	_, outPipe, _, _ := scriptRunner.RunAsync(context.TODO(), script)

	bytes := make([]byte, 1024)

	n, _ := outPipe.Read(bytes)

	fmt.Println(string(bytes[0:n]))

	fmt.Println(string(res))

	fmt.Println(script.Path())
	fmt.Println(script.Args())
	fmt.Println(script.Content())

	err = repo.Delete(id)

	if err != nil {
		t.Fatal(err)
	}

}
