package main

import (
	"context"
	"github.com/krol44/capy-converter/api"
	"github.com/krol44/capy-converter/pkg"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"io"
	"net/http"
	"os"
	"testing"
	"time"
)

func TestConverter(t *testing.T) {
	// starting server
	go func() {
		_ = os.Setenv("MAX_FILE_SIZE_MB", "100")
		api.Run()
	}()
	time.Sleep(time.Second * 2)

	t.Log("starting client")
	conn, _ := grpc.Dial(":3003", grpc.WithTransportCredentials(insecure.NewCredentials()))
	client := converter.NewConverterClient(conn)

	t.Log("starting GifToWebM")
	get, err := http.Get("https://media4.giphy.com/media/V4NSR1NG2p0KeJJyr5/giphy.gif")
	if err != nil {
		t.Error(err)
	}
	all, err := io.ReadAll(get.Body)
	if err != nil {
		t.Error(err)
	}

	resp, err := client.GifToWebM(context.Background(), &converter.GifToWebMType{File: all})

	if err != nil {
		t.Errorf("could not get answer: %v", err)
	} else {
		if len(resp.File) < 300000 {
			t.Errorf("error length - %d < 300000", len(resp.File))
		}
	}

	t.Log("finish test")
}
