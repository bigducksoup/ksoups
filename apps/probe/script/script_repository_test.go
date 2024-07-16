package script

import (
	"context"
	"fmt"
	"testing"
	"time"
)

func TestA(t *testing.T) {

	repo := NewShellScriptRepo("../../scripts")

	id, err := repo.Store([]byte("sleep 5"), "helloworld", nil)

	if err != nil {
		t.Fatal(err)
	}

	script, ok := repo.Get(id)

	if !ok {
		t.Fatal("no script")
	}

	scriptRunner := NewShellScriptRunner()

	ctx, cancel := context.WithTimeout(context.TODO(), 1*time.Second)

	defer cancel()

	res, _, err := scriptRunner.Run(ctx, script)

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
