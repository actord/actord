package main

import (
	"log"

	"github.com/actord/actord/pkg/actor"
	"github.com/actord/actord/pkg/executor"
	"github.com/actord/actord/pkg/process"
)

func main() {
	packagePath := "example/auth"
	p, err := process.Parse(packagePath)
	if err != nil {
		log.Fatal(err)
	}

	exec := executor.New()

	log.Println("======== create")
	act, err := exec.Trigger(p, "create", []byte(`{"email":"test@test.com", "password": "123", "password_confirmation": "123"}`))
	if err != nil {
		log.Fatal(err)
	}
	printActor(act)

	log.Println("======== block")
	if err := exec.SendEvent(p, act, "block", []byte(`{"reason": "just reason"}`)); err != nil {
		log.Fatal(err)
	}
	printActor(act)

	log.Println("======== unblock")
	if err := exec.SendEvent(p, act, "unblock", []byte(`{}`)); err != nil {
		log.Fatal(err)
	}
	printActor(act)

	log.Println("======== login invalid")
	if err := exec.SendEvent(p, act, "login", []byte(`{"password": "1233"}`)); err != nil {
		log.Fatal(err)
	}
	printActor(act)

	log.Println("======== login valid")
	if err := exec.SendEvent(p, act, "login", []byte(`{"password": "123"}`)); err != nil {
		log.Fatal(err)
	}
	printActor(act)
}

func printActor(act *actor.Actor) {
	log.Println("ID:", act.ID)
	log.Println("ACT:", string(act.Data))
	log.Println("STATE:", act.State)
	log.Println("LogicIndex:", act.LogicIndex)
	log.Println("AwaitEvent:", act.AwaitEvent)
}
